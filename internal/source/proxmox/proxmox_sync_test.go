package proxmox

import (
	"context"
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
