package proxmox

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/luthermonson/go-proxmox"
)

// newTestSource returns a ProxmoxSource wired to a site-scoped NetBox cluster,
// ready to be synced against inventory.MockInventory.
func newTestSource(t *testing.T, sourceConfig *parser.SourceConfig) *ProxmoxSource {
	t.Helper()
	lg, err := logger.New("", logger.ERROR)
	if err != nil {
		t.Fatal(err)
	}
	if sourceConfig.Name == "" {
		sourceConfig.Name = "test-proxmox"
	}
	return &ProxmoxSource{
		Config: common.Config{
			Logger:        lg,
			Ctx:           context.WithValue(context.Background(), constants.CtxSourceKey, sourceConfig.Name),
			SourceConfig:  sourceConfig,
			SourceNameTag: &objects.Tag{ID: 9001, Name: "Source: " + sourceConfig.Name, Slug: "source-test-proxmox"},
			SourceTypeTag: &objects.Tag{ID: 9002, Name: "proxmox", Slug: "proxmox"},
		},
		NetboxCluster: &objects.Cluster{
			NetboxObject: objects.NetboxObject{ID: 1},
			Name:         "test-cluster",
			ScopeType:    constants.ContentTypeDcimSite,
			ScopeID:      1,
		},
	}
}

func newTestVM(name string, vmid uint64) *proxmox.VirtualMachine {
	osType := "win11"
	return &proxmox.VirtualMachine{
		Name:                 name,
		VMID:                 proxmox.StringOrUint64(vmid),
		VirtualMachineConfig: &proxmox.VirtualMachineConfig{OSType: &osType},
	}
}

func TestSyncVMsReturnsWhenMoreVMsThanNodesFail(t *testing.T) {
	nbi := inventory.MockInventory
	saved := nbi.NetboxAPI
	nbi.NetboxAPI = service.FailingMockNetboxClient
	defer func() { nbi.NetboxAPI = saved }()

	ps := newTestSource(t, &parser.SourceConfig{})
	ps.NetboxNodes = map[string]*objects.Device{
		"n1": {
			NetboxObject: objects.NetboxObject{ID: 1},
			Name:         "n1",
			Site:         &objects.Site{NetboxObject: objects.NetboxObject{ID: 1}},
		},
	}
	// One node and two VMs whose platform creation fails against the failing client.
	ps.Vms = map[string][]*proxmox.VirtualMachine{
		"n1": {newTestVM("failing-vm1", 101), newTestVM("failing-vm2", 102)},
	}

	done := make(chan error, 1)
	go func() { done <- ps.syncVMs(nbi) }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatalf("syncVMs returned nil, want an error")
		}
		for _, vmName := range []string{"failing-vm1", "failing-vm2"} {
			if !strings.Contains(err.Error(), vmName) {
				t.Errorf("syncVMs error does not report %s: %s", vmName, err)
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("syncVMs did not return within 5s (1 node, 2 failing VMs)")
	}
}

func TestNodeInterfaceType(t *testing.T) {
	tests := []struct {
		proxmoxType string
		want        *objects.InterfaceType
	}{
		{proxmoxType: "bond", want: &objects.LAGInterfaceType},
		{proxmoxType: "OVSBond", want: &objects.LAGInterfaceType},
		{proxmoxType: "bridge", want: &objects.BridgeInterfaceType},
		{proxmoxType: "OVSBridge", want: &objects.BridgeInterfaceType},
		{proxmoxType: "vlan", want: &objects.VirtualInterfaceType},
		{proxmoxType: "OVSIntPort", want: &objects.VirtualInterfaceType},
		{proxmoxType: "alias", want: &objects.VirtualInterfaceType},
		// Physical NICs: Proxmox does not expose their real type, so it is left to NetBox.
		{proxmoxType: "eth", want: nil},
		{proxmoxType: "", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.proxmoxType, func(t *testing.T) {
			got := nodeInterfaceType(tt.proxmoxType)
			if (got == nil) != (tt.want == nil) || (got != nil && got.Value != tt.want.Value) {
				t.Errorf("nodeInterfaceType(%q) = %v, want %v", tt.proxmoxType, got, tt.want)
			}
		})
	}
}

func TestSyncNodeNetworksTypesInterfacesAndAttachesBondMembers(t *testing.T) {
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	ps := newTestSource(t, &parser.SourceConfig{})
	node := &proxmox.Node{Name: "pve-lag"}
	host := &objects.Device{NetboxObject: objects.NetboxObject{ID: 4242}, Name: "pve-lag"}
	ps.NetboxNodes = map[string]*objects.Device{node.Name: host}
	// Members are listed before their bond, as Proxmox may return them.
	ps.NodeIfaces = map[string][]*proxmox.NodeNetwork{node.Name: {
		{Iface: "eno1", Type: "eth"},
		{Iface: "eno2", Type: "eth"},
		{Iface: "bond0", Type: "bond", Slaves: "eno1 eno2"},
		{Iface: "vmbr0", Type: "bridge", BridgePorts: "bond0"},
		{Iface: "vmbr0.100", Type: "vlan"},
	}}

	if err := ps.syncNodeNetworks(nbi, node); err != nil {
		t.Fatalf("syncNodeNetworks() error = %v", err)
	}

	wantTypes := map[string]string{
		"eno1": "other", "eno2": "other", "bond0": "lag", "vmbr0": "bridge", "vmbr0.100": "virtual",
	}
	for name, wantType := range wantTypes {
		iface, ok := nbi.GetInterface(name, host.ID)
		if !ok {
			t.Errorf("interface %s not synced", name)
			continue
		}
		if iface.Type == nil || iface.Type.Value != wantType {
			t.Errorf("interface %s type = %v, want %s", name, iface.Type, wantType)
		}
	}
	for _, member := range []string{"eno1", "eno2"} {
		iface, ok := nbi.GetInterface(member, host.ID)
		if !ok {
			continue
		}
		if iface.LAG == nil || iface.LAG.Name != "bond0" {
			t.Errorf("interface %s LAG = %v, want bond0", member, iface.LAG)
		}
	}
}

func TestSyncNodesWithDomainSuffixSyncsHostInterfaces(t *testing.T) {
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	nbi.IgnoreDeviceTypeTag = &objects.Tag{ID: 9003, Name: constants.IgnoreDeviceTypeTagName, Slug: "ignore"}
	tests := []struct {
		node, suffix, wantHostName string
	}{
		{node: "suffix-node-a", suffix: "", wantHostName: "suffix-node-a"},
		{node: "suffix-node-b", suffix: ".lab01", wantHostName: "suffix-node-b.lab01"},
	}
	for _, tt := range tests {
		t.Run(tt.wantHostName, func(t *testing.T) {
			ps := newTestSource(t, &parser.SourceConfig{AssignDomainName: tt.suffix})
			ps.Nodes = []*proxmox.Node{{Name: tt.node}}
			ps.NodeIfaces = map[string][]*proxmox.NodeNetwork{tt.node: {{Iface: "eno1", Type: "eth"}}}
			if err := ps.syncNodes(nbi); err != nil {
				t.Fatalf("syncNodes() error = %v", err)
			}
			// syncVMs and syncContainers look hosts up by their Proxmox node name.
			host := ps.NetboxNodes[tt.node]
			if host == nil {
				t.Fatalf("no NetBox host for Proxmox node %s", tt.node)
			}
			if host.Name != tt.wantHostName {
				t.Errorf("host name = %q, want %q", host.Name, tt.wantHostName)
			}
			if _, ok := nbi.GetInterface("eno1", host.ID); !ok {
				t.Errorf("host %s has no interface eno1", host.Name)
			}
		})
	}
}

func TestVMPlatformName(t *testing.T) {
	l26 := "l26"
	existingWithPlatform := &objects.VM{Platform: &objects.Platform{NetboxObject: objects.NetboxObject{ID: 6}, Name: "Debian 12"}}
	tests := []struct {
		name         string
		agentOsInfo  *proxmox.AgentOsInfo
		osType       *string
		existingVM   *objects.VM
		wantName     string
		wantKeepCurr bool
	}{
		{
			name:        "agent reports the OS",
			agentOsInfo: &proxmox.AgentOsInfo{PrettyName: "Debian GNU/Linux 13 (trixie)"},
			osType:      &l26,
			existingVM:  existingWithPlatform,
			wantName:    "Debian GNU/Linux 13 (trixie)",
		},
		{
			name:         "agent silent keeps the current platform",
			osType:       &l26,
			existingVM:   existingWithPlatform,
			wantKeepCurr: true,
		},
		{
			name:         "agent without OS name keeps the current platform",
			agentOsInfo:  &proxmox.AgentOsInfo{},
			osType:       &l26,
			existingVM:   existingWithPlatform,
			wantKeepCurr: true,
		},
		{
			name:     "agent silent on a new VM falls back to ostype",
			osType:   &l26,
			wantName: "Other 2.6.x Linux (64-bit)",
		},
		{
			name:       "agent silent on a VM without platform falls back to ostype",
			osType:     &l26,
			existingVM: &objects.VM{},
			wantName:   "Other 2.6.x Linux (64-bit)",
		},
		{
			name:     "nothing known",
			wantName: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotKeep := vmPlatformName(tt.agentOsInfo, tt.osType, tt.existingVM)
			if gotKeep != tt.wantKeepCurr || (!gotKeep && gotName != tt.wantName) {
				t.Errorf("vmPlatformName() = (%q, %t), want (%q, %t)", gotName, gotKeep, tt.wantName, tt.wantKeepCurr)
			}
		})
	}
}

func TestSyncVMKeepsKnownNetworkObjectsWhenAgentDataIsUnknown(t *testing.T) {
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	host := &objects.Device{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "n1",
		Site:         &objects.Site{NetboxObject: objects.NetboxObject{ID: 1}},
	}
	tests := []struct {
		name       string
		vmIfaces   map[uint64][]*proxmox.AgentNetworkIface
		wantOrphan bool
	}{
		{name: "agent data unknown", vmIfaces: map[uint64][]*proxmox.AgentNetworkIface{}, wantOrphan: false},
		{
			name:       "agent reports no interface",
			vmIfaces:   map[uint64][]*proxmox.AgentNetworkIface{1: {}},
			wantOrphan: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existingIface := nbi.GetVMInterfaceByVMIDAndName(1, "vmeth0")
			if existingIface == nil {
				t.Fatalf("mock inventory has no vmeth0 on VM 1")
			}
			nbi.OrphanManager.AddItem(existingIface)
			t.Cleanup(func() { nbi.OrphanManager.RemoveItem(existingIface) })

			ps := newTestSource(t, &parser.SourceConfig{})
			ps.NetboxCluster.ID = 1
			ps.VMIfaces = tt.vmIfaces
			vm := newTestVM("existing_vm1", 1)
			vm.Status = "stopped"
			if err := ps.syncVM(nbi, vm, host); err != nil {
				t.Fatalf("syncVM() error = %v", err)
			}
			_, orphan := nbi.OrphanManager.Items[constants.VMInterfacesAPIPath][existingIface.ID]
			if orphan != tt.wantOrphan {
				t.Errorf("vmeth0 orphan = %t, want %t", orphan, tt.wantOrphan)
			}
		})
	}
}

func TestParseDiskSizeMiB(t *testing.T) {
	tests := []struct {
		item string
		want int
	}{
		{item: "size=32G", want: 32768},
		{item: "size=2T", want: 2097152},
		{item: "size=2252M", want: 2252},
		{item: "size=1.5T", want: 1572864},
		{item: "size=34359738368", want: 32768},
		{item: "size=4194304K", want: 4096},
		{item: "size=32GiB", want: 32768},
		{item: "size=", want: 0},
		{item: "size=abc", want: 0},
		{item: "discard=on", want: 0},
		{item: "cache=writeback", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			if got := parseDiskSizeMiB(tt.item); got != tt.want {
				t.Errorf("parseDiskSizeMiB(%q) = %d, want %d", tt.item, got, tt.want)
			}
		})
	}
}

func TestSyncStopsWhenClusterFailsEvenWithContinueOnError(t *testing.T) {
	nbi := inventory.MockInventory
	saved := nbi.NetboxAPI
	nbi.NetboxAPI = service.FailingMockNetboxClient
	defer func() { nbi.NetboxAPI = saved }()

	ps := newTestSource(t, &parser.SourceConfig{ContinueOnError: true})
	ps.NetboxCluster = nil
	ps.Cluster = &proxmox.Cluster{Name: "failing-cluster-type"}
	ps.Nodes = []*proxmox.Node{{Name: "n1"}}

	if err := ps.Sync(nbi); err == nil {
		t.Errorf("Sync() = nil, want the cluster error")
	}
}

func TestSyncVMsAndContainersSkipGuestsWithoutHost(t *testing.T) {
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	ps := newTestSource(t, &parser.SourceConfig{})
	ps.NetboxNodes = map[string]*objects.Device{}
	ps.Vms = map[string][]*proxmox.VirtualMachine{"missing-node": {newTestVM("orphan-guest-vm", 301)}}
	ps.Containers = map[string][]*proxmox.Container{"missing-node": {{Name: "orphan-guest-ct", VMID: 302}}}

	if err := ps.syncVMs(nbi); err != nil {
		t.Errorf("syncVMs() error = %v, want the VM skipped", err)
	}
	if err := ps.syncContainers(nbi); err != nil {
		t.Errorf("syncContainers() error = %v, want the container skipped", err)
	}
}

func TestSyncContainersReportsDiskInMiB(t *testing.T) {
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	ps := newTestSource(t, &parser.SourceConfig{})
	ps.NetboxNodes = map[string]*objects.Device{
		"n1": {
			NetboxObject: objects.NetboxObject{ID: 1},
			Name:         "n1",
			Site:         &objects.Site{NetboxObject: objects.NetboxObject{ID: 1}},
		},
	}
	ps.Containers = map[string][]*proxmox.Container{
		"n1": {{Name: "disk-unit-ct", VMID: 401, MaxDisk: 8 * constants.GiB, MaxMem: 512 * constants.MiB}},
	}
	ps.ContainerIfaces = map[uint64][]*proxmox.ContainerInterface{}

	if err := ps.syncContainers(nbi); err != nil {
		t.Fatalf("syncContainers() error = %v", err)
	}
	nbContainer, ok := nbi.GetVM("disk-unit-ct", ps.NetboxCluster.ID)
	if !ok {
		t.Fatalf("container disk-unit-ct not synced")
	}
	if nbContainer.Disk != 8192 {
		t.Errorf("container disk = %d, want 8192 (MiB)", nbContainer.Disk)
	}
}

func TestSplitProxmoxTags(t *testing.T) {
	tests := []struct {
		raw  string
		want []string
	}{
		{raw: "", want: []string{}},
		{raw: " ", want: []string{}},
		{raw: "prod", want: []string{"prod"}},
		{raw: "prod;", want: []string{"prod"}},
		{raw: " prod ; web ", want: []string{"prod", "web"}},
		{raw: "prod;;web", want: []string{"prod", "web"}},
	}
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			if got := splitProxmoxTags(tt.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitProxmoxTags(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestSyncSkipsHomonymGuestsOfACluster(t *testing.T) {
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	ps := newTestSource(t, &parser.SourceConfig{})
	ps.NetboxCluster.ID = 77
	site := &objects.Site{NetboxObject: objects.NetboxObject{ID: 1}}
	ps.NetboxNodes = map[string]*objects.Device{
		"n1": {NetboxObject: objects.NetboxObject{ID: 1}, Name: "n1", Site: site},
		"n2": {NetboxObject: objects.NetboxObject{ID: 2}, Name: "n2", Site: site},
	}
	// The lowest VMID of homonym guests is synced, VMs and containers alike.
	ps.Vms = map[string][]*proxmox.VirtualMachine{
		"n1": {newTestVM("dup-test", 101), newTestVM("dup-guest", 120)},
		"n2": {newTestVM("dup-test", 205)},
	}
	ps.Containers = map[string][]*proxmox.Container{"n2": {{Name: "dup-guest", VMID: 150}}}
	ps.VMIfaces = map[uint64][]*proxmox.AgentNetworkIface{
		101: {{Name: "eth0"}},
		120: {{Name: "eth0"}},
		205: {{Name: "eth1"}},
	}
	ps.ContainerIfaces = map[uint64][]*proxmox.ContainerInterface{150: {{Name: "veth150"}}}

	if err := ps.syncVMs(nbi); err != nil {
		t.Fatalf("syncVMs() error = %v", err)
	}
	if err := ps.syncContainers(nbi); err != nil {
		t.Fatalf("syncContainers() error = %v", err)
	}

	tests := []struct {
		vmName, keptIface, skippedIface string
	}{
		{vmName: "dup-test", keptIface: "eth0", skippedIface: "eth1"},
		{vmName: "dup-guest", keptIface: "eth0", skippedIface: "veth150"},
	}
	for _, tt := range tests {
		nbVM, ok := nbi.GetVM(tt.vmName, ps.NetboxCluster.ID)
		if !ok {
			t.Errorf("vm %s not synced", tt.vmName)
			continue
		}
		if nbi.GetVMInterfaceByVMIDAndName(nbVM.ID, tt.keptIface) == nil {
			t.Errorf("vm %s has no interface %s of the lowest VMID", tt.vmName, tt.keptIface)
		}
		if nbi.GetVMInterfaceByVMIDAndName(nbVM.ID, tt.skippedIface) != nil {
			t.Errorf("vm %s got interface %s of a skipped homonym", tt.vmName, tt.skippedIface)
		}
	}
}

func TestKeptGuestIDsIgnoresGuestsThatAreNotSynced(t *testing.T) {
	ps := newTestSource(t, &parser.SourceConfig{IgnoreVMTemplates: true})
	ps.NetboxNodes = map[string]*objects.Device{"n1": {Name: "n1"}}
	template := newTestVM("debian", 100)
	template.Template = true
	ps.Vms = map[string][]*proxmox.VirtualMachine{
		"n1":        {template, newTestVM("debian", 200), newTestVM("app", 60)},
		"lost-node": {newTestVM("app", 50)},
	}
	ps.Containers = map[string][]*proxmox.Container{"lost-node": {{Name: "web", VMID: 10}}}

	got := ps.keptGuestIDs()

	want := map[string]uint64{"debian": 200, "app": 60}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("keptGuestIDs() = %v, want %v", got, want)
	}
}
