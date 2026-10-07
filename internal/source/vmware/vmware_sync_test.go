package vmware

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// newTestSource returns a VmwareSource with empty vSphere data, ready to be synced
// against inventory.MockInventory in dry-run mode.
func newTestSource(t *testing.T, sourceConfig *parser.SourceConfig) *VmwareSource {
	t.Helper()
	lg, err := logger.New("", logger.ERROR)
	if err != nil {
		t.Fatal(err)
	}
	if sourceConfig.Name == "" {
		sourceConfig.Name = "test-vmware"
	}
	return &VmwareSource{
		Config: common.Config{
			Logger:        lg,
			Ctx:           context.WithValue(context.Background(), constants.CtxSourceKey, sourceConfig.Name),
			SourceConfig:  sourceConfig,
			SourceNameTag: &objects.Tag{ID: 9101, Name: "Source: " + sourceConfig.Name, Slug: "source-test-vmware"},
			SourceTypeTag: &objects.Tag{ID: 9102, Name: "vmware", Slug: "vmware"},
		},
		Networks: NetworkData{
			DistributedVirtualPortgroups: map[string]*DistributedPortgroupData{},
			Vid2Name:                     map[int]string{},
			HostVirtualSwitches:          map[string]map[string]*HostVirtualSwitchData{},
			HostProxySwitches:            map[string]map[string]*HostProxySwitchData{},
			HostPortgroups:               map[string]map[string]*HostPortgroupData{},
		},
	}
}

// newTestInventory returns the dry-run mock inventory, with the default site that
// NetboxInventory.Init creates in a real run.
func newTestInventory(t *testing.T) (*inventory.NetboxInventory, *objects.Site) {
	t.Helper()
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	defaultSite, err := nbi.AddSite(context.Background(), &objects.Site{
		Name: constants.DefaultSite,
		Slug: "defaultsite",
	})
	if err != nil {
		t.Fatal(err)
	}
	return nbi, defaultSite
}

// newTestHost returns an ESXi host with the summary fields vCenter reports for a
// connected host.
func newTestHost(id, name string) mo.HostSystem {
	host := mo.HostSystem{
		Summary: types.HostListSummary{
			Hardware: &types.HostHardwareSummary{Vendor: "Dell Inc.", Model: "PowerEdge R650", NumCpuCores: 32},
			Runtime:  &types.HostRuntimeInfo{ConnectionState: types.HostSystemConnectionStateConnected},
			Config: types.HostConfigSummary{
				Product: &types.AboutInfo{Name: "VMware ESXi", Version: "8.0.3"},
			},
		},
	}
	host.Self = types.ManagedObjectReference{Type: "HostSystem", Value: id}
	host.Name = name
	return host
}

// newTestVMOnHost returns a powered-on VM registered on the given host, with one disk.
func newTestVMOnHost(vc *VmwareSource, id, name, hostID string) mo.VirtualMachine {
	vm := mo.VirtualMachine{
		Config: &types.VirtualMachineConfigInfo{
			GuestFullName: "Debian GNU/Linux 12 (64-bit)",
			Hardware: types.VirtualHardware{
				NumCPU:   2,
				MemoryMB: 4096,
				Device: []types.BaseVirtualDevice{&types.VirtualDisk{
					VirtualDevice:   types.VirtualDevice{DeviceInfo: &types.Description{Label: "Hard disk 1"}},
					CapacityInBytes: 8 * constants.GiB,
				}},
			},
		},
		Runtime: types.VirtualMachineRuntimeInfo{PowerState: types.VirtualMachinePowerStatePoweredOn},
		Guest:   &types.GuestInfo{},
	}
	vm.Self = types.ManagedObjectReference{Type: "VirtualMachine", Value: id}
	vm.Name = name
	if vc.VM2Host == nil {
		vc.VM2Host = map[string]string{}
	}
	vc.VM2Host[id] = hostID
	return vm
}

func TestSyncVMsSkipsVMsWithoutConfig(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	vc.Hosts = map[string]mo.HostSystem{"host-2": newTestHost("host-2", "esx-no-config")}
	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v", err)
	}
	inaccessible := newTestVMOnHost(vc, "vm-20", "inaccessible-vm", "host-2")
	inaccessible.Config = nil
	vc.Vms = map[string]mo.VirtualMachine{"vm-20": inaccessible}

	if err := vc.syncVMs(nbi); err != nil {
		t.Errorf("syncVMs() error = %v, want the VM without config skipped", err)
	}
}

func TestCollectHostPhysicalNicDataSpeed(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	nbHost := &objects.Device{Name: "esx-nic-speed"}
	tests := []struct {
		name      string
		pnic      types.PhysicalNic
		wantSpeed objects.InterfaceSpeed
	}{
		{
			name:      "link up",
			pnic:      types.PhysicalNic{Device: "vmnic0", LinkSpeed: &types.PhysicalNicLinkInfo{SpeedMb: 10000}},
			wantSpeed: 10000000,
		},
		{
			name: "link down with a forced speed",
			pnic: types.PhysicalNic{
				Device: "vmnic1",
				Spec:   types.PhysicalNicSpec{LinkSpeed: &types.PhysicalNicLinkInfo{SpeedMb: 1000}},
			},
			wantSpeed: 1000000,
		},
		{
			name:      "link down without any speed",
			pnic:      types.PhysicalNic{Device: "vmnic2", ValidLinkSpecification: []types.PhysicalNicLinkInfo{}},
			wantSpeed: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iface, _, err := vc.collectHostPhysicalNicData(nbi, nbHost, tt.pnic, nil)
			if err != nil {
				t.Fatalf("collectHostPhysicalNicData() error = %v", err)
			}
			if iface.Speed != tt.wantSpeed {
				t.Errorf("speed = %d, want %d", iface.Speed, tt.wantSpeed)
			}
		})
	}
}

func TestSyncHostsWithoutHostSiteRelationsUsesDefaultSite(t *testing.T) {
	nbi, defaultSite := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	vc.Hosts = map[string]mo.HostSystem{"host-1": newTestHost("host-1", "esx-default-site")}

	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v, want the host synced to the default site", err)
	}
	if _, ok := nbi.GetDevice("esx-default-site", defaultSite.ID); !ok {
		t.Errorf("host esx-default-site not found in site %s", constants.DefaultSite)
	}
}

func TestParseDvpgVlanSpec(t *testing.T) {
	trunk := func(ranges ...types.NumericRange) *types.VmwareDistributedVirtualSwitchTrunkVlanSpec {
		return &types.VmwareDistributedVirtualSwitchTrunkVlanSpec{VlanId: ranges}
	}
	tests := []struct {
		name        string
		spec        types.BaseVmwareDistributedVirtualSwitchVlanSpec
		wantIDs     []int
		wantRanges  []string
		wantPrivate bool
	}{
		{
			name:    "access vlan",
			spec:    &types.VmwareDistributedVirtualSwitchVlanIdSpec{VlanId: 42},
			wantIDs: []int{42},
		},
		{
			name:        "private vlan",
			spec:        &types.VmwareDistributedVirtualSwitchPvlanSpec{PvlanId: 7},
			wantIDs:     []int{7},
			wantPrivate: true,
		},
		{
			name:       "trunk of single vlans and a range",
			spec:       trunk(types.NumericRange{Start: 10, End: 10}, types.NumericRange{Start: 100, End: 103}),
			wantIDs:    []int{10, 100, 101, 102, 103},
			wantRanges: []string{"10", "100-103"},
		},
		{
			name:       "trunk of every vlan",
			spec:       trunk(types.NumericRange{Start: 0, End: 4094}),
			wantIDs:    []int{4095},
			wantRanges: []string{"0-4094"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, ranges, private, err := parseDvpgVlanSpec(tt.spec)
			if err != nil {
				t.Fatalf("parseDvpgVlanSpec() error = %v", err)
			}
			if !reflect.DeepEqual(ids, tt.wantIDs) || !reflect.DeepEqual(ranges, tt.wantRanges) ||
				private != tt.wantPrivate {
				t.Errorf("parseDvpgVlanSpec() = (%v, %v, %t), want (%v, %v, %t)",
					ids, ranges, private, tt.wantIDs, tt.wantRanges, tt.wantPrivate)
			}
		})
	}
}

func TestCollectVMInterfaceDataBoundsDescriptionOfWideTrunks(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	ids, ranges, _, err := parseDvpgVlanSpec(&types.VmwareDistributedVirtualSwitchTrunkVlanSpec{
		VlanId: []types.NumericRange{{Start: 100, End: 199}, {Start: 300, End: 399}},
	})
	if err != nil {
		t.Fatal(err)
	}
	vc.Networks.DistributedVirtualPortgroups["dvportgroup-1"] = &DistributedPortgroupData{
		Name: "trunk-wide", VlanIDs: ids, VlanIDRanges: ranges,
	}
	card := &types.VirtualEthernetCard{
		VirtualDevice: types.VirtualDevice{
			DeviceInfo: &types.Description{Label: "Network adapter 1"},
			Backing: &types.VirtualEthernetCardDistributedVirtualPortBackingInfo{
				Port: types.DistributedVirtualSwitchPortConnection{PortgroupKey: "dvportgroup-1"},
			},
			Connectable: &types.VirtualDeviceConnectInfo{Connected: true},
		},
	}
	netboxVM := &objects.VM{Name: "trunk-vm", Host: &objects.Device{Name: "esx-trunk"}}

	_, _, iface, _, err := vc.collectVMInterfaceData(nbi, netboxVM, mo.VirtualMachine{Guest: &types.GuestInfo{}}, card)
	if err != nil {
		t.Fatalf("collectVMInterfaceData() error = %v", err)
	}
	want := "Network adapter 1 (vlan ID: 100-199, 300-399)"
	if iface.Description != want {
		t.Errorf("description = %q, want %q", iface.Description, want)
	}
}

func TestCollectVMInterfaceDataTruncatesLongDescriptions(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	vlanRanges := make([]types.NumericRange, 0, 60)
	for vid := int32(1000); vid < 1060; vid++ {
		vlanRanges = append(vlanRanges, types.NumericRange{Start: vid, End: vid})
	}
	ids, ranges, _, err := parseDvpgVlanSpec(&types.VmwareDistributedVirtualSwitchTrunkVlanSpec{VlanId: vlanRanges})
	if err != nil {
		t.Fatal(err)
	}
	vc.Networks.DistributedVirtualPortgroups["dvportgroup-2"] = &DistributedPortgroupData{
		Name: "trunk-many", VlanIDs: ids, VlanIDRanges: ranges,
	}
	card := &types.VirtualEthernetCard{
		VirtualDevice: types.VirtualDevice{
			DeviceInfo: &types.Description{Label: "Network adapter 2"},
			Backing: &types.VirtualEthernetCardDistributedVirtualPortBackingInfo{
				Port: types.DistributedVirtualSwitchPortConnection{PortgroupKey: "dvportgroup-2"},
			},
			Connectable: &types.VirtualDeviceConnectInfo{Connected: true},
		},
	}
	netboxVM := &objects.VM{Name: "trunk-many-vm", Host: &objects.Device{Name: "esx-trunk"}}

	_, _, iface, _, err := vc.collectVMInterfaceData(nbi, netboxVM, mo.VirtualMachine{Guest: &types.GuestInfo{}}, card)
	if err != nil {
		t.Fatalf("collectVMInterfaceData() error = %v", err)
	}
	if n := len([]rune(iface.Description)); n > objects.MaxDescriptionLength {
		t.Errorf("description has %d characters, want at most %d", n, objects.MaxDescriptionLength)
	}
	if !strings.HasPrefix(iface.Description, "Network adapter 2 (vlan ID: 1000, 1001,") {
		t.Errorf("description = %q, want it to start with the adapter and its first VLANs", iface.Description)
	}
}

func TestSyncHostVirtualNicsCollectsAddressesForThePrimaryIP(t *testing.T) {
	nbi, defaultSite := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	host := newTestHost("host-3", "esx-primary-ip")
	vc.Hosts = map[string]mo.HostSystem{"host-3": host}
	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v", err)
	}
	nbHost, ok := nbi.GetDevice("esx-primary-ip", defaultSite.ID)
	if !ok {
		t.Fatal("host esx-primary-ip not synced")
	}
	host.Config = &types.HostConfigInfo{Network: &types.HostNetworkInfo{Vnic: []types.HostVirtualNic{{
		Device: "vmk0",
		Spec: types.HostVirtualNicSpec{
			Ip:  &types.HostIpConfig{IpAddress: "192.0.2.10", SubnetMask: "255.255.255.0"},
			Mac: "00:50:56:aa:bb:cc",
		},
	}}}}

	hostIPv4Addresses, _, err := vc.syncHostVirtualNics(nbi, host, nbHost)
	if err != nil {
		t.Fatalf("syncHostVirtualNics() error = %v", err)
	}
	if len(hostIPv4Addresses) != 1 || hostIPv4Addresses[0].Address != "192.0.2.10/24" {
		t.Errorf("collected IPv4 addresses = %v, want [192.0.2.10/24]", hostIPv4Addresses)
	}
}

func TestSyncVMsMatchesVMRoleRelationsOnTheVMName(t *testing.T) {
	nbi, defaultSite := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{VMRoleRelations: map[string]string{"^db-": "Database"}})
	vc.Hosts = map[string]mo.HostSystem{"host-4": newTestHost("host-4", "db-esx")}
	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v", err)
	}
	nbHost, _ := nbi.GetDevice("db-esx", defaultSite.ID)
	vc.Vms = map[string]mo.VirtualMachine{
		"vm-40": newTestVMOnHost(vc, "vm-40", "db-role-vm", "host-4"),
		"vm-41": newTestVMOnHost(vc, "vm-41", "web-role-vm", "host-4"),
	}

	if err := vc.syncVMs(nbi); err != nil {
		t.Fatalf("syncVMs() error = %v", err)
	}
	tests := []struct{ vmName, wantRole string }{
		{vmName: "db-role-vm", wantRole: "Database"},
		{vmName: "web-role-vm", wantRole: constants.DeviceRoleVM},
	}
	for _, tt := range tests {
		vm, ok := nbi.GetVM(tt.vmName, nbHost.Cluster.ID)
		if !ok {
			t.Errorf("vm %s not synced", tt.vmName)
			continue
		}
		if vm.Role == nil || vm.Role.Name != tt.wantRole {
			t.Errorf("vm %s role = %v, want %s", tt.vmName, vm.Role, tt.wantRole)
		}
	}
}

func TestSyncHostsWithoutOptionalSummaryFields(t *testing.T) {
	nbi, defaultSite := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	host := newTestHost("host-5", "esx-bare-summary")
	host.Summary = types.HostListSummary{}
	vc.Hosts = map[string]mo.HostSystem{"host-5": host}

	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v", err)
	}
	nbHost, ok := nbi.GetDevice("esx-bare-summary", defaultSite.ID)
	if !ok {
		t.Fatal("host esx-bare-summary not synced")
	}
	if nbHost.Status == nil || *nbHost.Status != objects.DeviceStatusOffline {
		t.Errorf("host status = %v, want %v without runtime info", nbHost.Status, objects.DeviceStatusOffline)
	}
}

func TestSyncVMsWithoutOptionalGuestFields(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	vc.Hosts = map[string]mo.HostSystem{"host-6": newTestHost("host-6", "esx-guest-fields")}
	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v", err)
	}
	withoutGuest := newTestVMOnHost(vc, "vm-60", "no-guest-vm", "host-6")
	withoutGuest.Guest = nil
	withoutRoutes := newTestVMOnHost(vc, "vm-61", "no-route-vm", "host-6")
	withoutRoutes.Guest = &types.GuestInfo{IpStack: []types.GuestStackInfo{{}}}
	vc.Vms = map[string]mo.VirtualMachine{"vm-60": withoutGuest, "vm-61": withoutRoutes}

	if err := vc.syncVMs(nbi); err != nil {
		t.Errorf("syncVMs() error = %v", err)
	}
}

func TestSyncVMsKeepsLongDescriptionsInComments(t *testing.T) {
	nbi, defaultSite := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{CustomFieldMappings: map[string]string{"Notes": "description"}})
	vc.CustomFieldID2Name = map[int32]string{1: "Notes"}
	vc.Hosts = map[string]mo.HostSystem{"host-7": newTestHost("host-7", "esx-description")}
	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v", err)
	}
	nbHost, _ := nbi.GetDevice("esx-description", defaultSite.ID)
	longText := strings.Repeat("long description ", 15)
	accentedText := strings.Repeat("é", 150)
	withDescription := func(id, name, description string) mo.VirtualMachine {
		vm := newTestVMOnHost(vc, id, name, "host-7")
		vm.Summary.CustomValue = []types.BaseCustomFieldValue{&types.CustomFieldStringValue{
			CustomFieldValue: types.CustomFieldValue{Key: 1},
			Value:            description,
		}}
		return vm
	}
	vc.Vms = map[string]mo.VirtualMachine{
		"vm-70": withDescription("vm-70", "long-description-vm", longText),
		"vm-71": withDescription("vm-71", "accented-description-vm", accentedText),
	}

	if err := vc.syncVMs(nbi); err != nil {
		t.Fatalf("syncVMs() error = %v", err)
	}
	tests := []struct{ vmName, wantDescription, wantComments string }{
		{vmName: "long-description-vm", wantDescription: "See comments.", wantComments: strings.TrimSpace(longText)},
		{vmName: "accented-description-vm", wantDescription: accentedText},
	}
	for _, tt := range tests {
		vm, ok := nbi.GetVM(tt.vmName, nbHost.Cluster.ID)
		if !ok {
			t.Errorf("vm %s not synced", tt.vmName)
			continue
		}
		if vm.Description != tt.wantDescription || vm.Comments != tt.wantComments {
			t.Errorf("vm %s description, comments = %q, %q, want %q, %q",
				tt.vmName, vm.Description, vm.Comments, tt.wantDescription, tt.wantComments)
		}
	}
}

func TestCollectVMInterfaceDataAssignsTaggedVlans(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	vc.Networks.DistributedVirtualPortgroups = map[string]*DistributedPortgroupData{
		"dvportgroup-10": {Name: "tagged-vlan-10", VlanIDs: []int{3010}},
		"dvportgroup-20": {Name: "tagged-vlan-20", VlanIDs: []int{3020}},
	}
	vc.Networks.Vid2Name = map[int]string{3010: "tagged-vlan-10", 3020: "tagged-vlan-20"}
	if err := vc.syncNetworks(nbi); err != nil {
		t.Fatalf("syncNetworks() error = %v", err)
	}
	vc.Networks.DistributedVirtualPortgroups["dvportgroup-30"] = &DistributedPortgroupData{
		Name: "tagged-trunk", VlanIDs: []int{3010, 3020}, VlanIDRanges: []string{"3010", "3020"},
	}
	card := &types.VirtualEthernetCard{
		VirtualDevice: types.VirtualDevice{
			DeviceInfo: &types.Description{Label: "Network adapter 1"},
			Backing: &types.VirtualEthernetCardDistributedVirtualPortBackingInfo{
				Port: types.DistributedVirtualSwitchPortConnection{PortgroupKey: "dvportgroup-30"},
			},
			Connectable: &types.VirtualDeviceConnectInfo{Connected: true},
		},
	}
	netboxVM := &objects.VM{Name: "tagged-vm", Host: &objects.Device{Name: "esx-tagged"}}

	_, _, iface, _, err := vc.collectVMInterfaceData(nbi, netboxVM, mo.VirtualMachine{Guest: &types.GuestInfo{}}, card)
	if err != nil {
		t.Fatalf("collectVMInterfaceData() error = %v", err)
	}
	gotVids := []int{}
	for _, vlan := range iface.TaggedVlans {
		gotVids = append(gotVids, vlan.Vid)
	}
	sort.Ints(gotVids)
	if !reflect.DeepEqual(gotVids, []int{3010, 3020}) {
		t.Errorf("tagged vlans = %v, want [3010 3020]", gotVids)
	}
}

func BenchmarkCollectVMInterfaceDataWideTrunk(b *testing.B) {
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	lg, err := logger.New("", logger.ERROR)
	if err != nil {
		b.Fatal(err)
	}
	vc := &VmwareSource{Config: common.Config{Logger: lg, Ctx: context.Background(), SourceConfig: &parser.SourceConfig{}}}
	ids, ranges, _, err := parseDvpgVlanSpec(&types.VmwareDistributedVirtualSwitchTrunkVlanSpec{
		VlanId: []types.NumericRange{{Start: 2, End: 4094}},
	})
	if err != nil {
		b.Fatal(err)
	}
	vc.Networks = NetworkData{
		DistributedVirtualPortgroups: map[string]*DistributedPortgroupData{
			"dvportgroup-wide": {Name: "trunk-2-4094", VlanIDs: ids, VlanIDRanges: ranges},
		},
		Vid2Name: map[int]string{},
	}
	card := &types.VirtualEthernetCard{
		VirtualDevice: types.VirtualDevice{
			DeviceInfo: &types.Description{Label: "Network adapter 1"},
			Backing: &types.VirtualEthernetCardDistributedVirtualPortBackingInfo{
				Port: types.DistributedVirtualSwitchPortConnection{PortgroupKey: "dvportgroup-wide"},
			},
			Connectable: &types.VirtualDeviceConnectInfo{},
		},
	}
	netboxVM := &objects.VM{Name: "wide-trunk-vm", Host: &objects.Device{Name: "esx-wide"}}
	b.ResetTimer()
	for range b.N {
		_, _, _, _, err := vc.collectVMInterfaceData(nbi, netboxVM, mo.VirtualMachine{Guest: &types.GuestInfo{}}, card)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestSyncClustersQualifiesClusterNamesWithTheirDatacenter(t *testing.T) {
	nbi, defaultSite := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	newCluster := func(id, name string) mo.ClusterComputeResource {
		cluster := mo.ClusterComputeResource{}
		cluster.Self = types.ManagedObjectReference{Type: "ClusterComputeResource", Value: id}
		cluster.Name = name
		return cluster
	}
	newDatacenter := func(id, name string) mo.Datacenter {
		datacenter := mo.Datacenter{}
		datacenter.Self = types.ManagedObjectReference{Type: "Datacenter", Value: id}
		datacenter.Name = name
		return datacenter
	}
	vc.DataCenters = map[string]mo.Datacenter{
		"datacenter-1": newDatacenter("datacenter-1", "dc-homonym-a"),
		"datacenter-2": newDatacenter("datacenter-2", "dc-homonym-b"),
	}
	vc.Clusters = map[string]mo.ClusterComputeResource{
		"domain-c1": newCluster("domain-c1", "Homonym01"),
		"domain-c2": newCluster("domain-c2", "Homonym01"),
	}
	vc.Cluster2Datacenter = map[string]string{"domain-c1": "datacenter-1", "domain-c2": "datacenter-2"}
	vc.Hosts = map[string]mo.HostSystem{
		"host-81": newTestHost("host-81", "esx-homonym-a"),
		"host-82": newTestHost("host-82", "esx-homonym-b"),
	}
	vc.Host2Cluster = map[string]string{"host-81": "domain-c1", "host-82": "domain-c2"}

	if err := vc.syncDatacenters(nbi); err != nil {
		t.Fatalf("syncDatacenters() error = %v", err)
	}
	if err := vc.syncClusters(nbi); err != nil {
		t.Fatalf("syncClusters() error = %v", err)
	}
	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v", err)
	}
	tests := []struct{ hostName, wantCluster string }{
		{hostName: "esx-homonym-a", wantCluster: "Homonym01 (dc-homonym-a)"},
		{hostName: "esx-homonym-b", wantCluster: "Homonym01 (dc-homonym-b)"},
	}
	for _, tt := range tests {
		host, ok := nbi.GetDevice(tt.hostName, defaultSite.ID)
		if !ok {
			t.Errorf("host %s not synced", tt.hostName)
			continue
		}
		if host.Cluster == nil || host.Cluster.Name != tt.wantCluster {
			t.Errorf("host %s cluster = %v, want %s", tt.hostName, host.Cluster, tt.wantCluster)
		}
	}
}

func TestSyncClustersRenamesClustersNamedBeforeQualification(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	legacy, err := nbi.AddCluster(vc.Ctx, &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			CustomFields: map[string]interface{}{constants.CustomFieldSourceIDName: "domain-c9"},
		},
		Name: "Legacy01",
	})
	if err != nil {
		t.Fatal(err)
	}
	cluster := mo.ClusterComputeResource{}
	cluster.Self = types.ManagedObjectReference{Type: "ClusterComputeResource", Value: "domain-c9"}
	cluster.Name = "Legacy01"
	datacenter := mo.Datacenter{}
	datacenter.Self = types.ManagedObjectReference{Type: "Datacenter", Value: "datacenter-9"}
	datacenter.Name = "dc-legacy"
	vc.Clusters = map[string]mo.ClusterComputeResource{"domain-c9": cluster}
	vc.DataCenters = map[string]mo.Datacenter{"datacenter-9": datacenter}
	vc.Cluster2Datacenter = map[string]string{"domain-c9": "datacenter-9"}

	if err := vc.syncClusters(nbi); err != nil {
		t.Fatalf("syncClusters() error = %v", err)
	}
	renamed, ok := nbi.GetCluster("Legacy01 (dc-legacy)")
	if !ok || renamed.ID != legacy.ID {
		t.Errorf("cluster Legacy01 (dc-legacy) = %v, want cluster %d renamed in place", renamed, legacy.ID)
	}
	if _, ok := nbi.GetCluster("Legacy01"); ok {
		t.Errorf("cluster Legacy01 still indexed under its former name")
	}
}

func TestInterfacesOnPortGroupsWithoutVlan(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	vc.Networks.DistributedVirtualPortgroups["dvportgroup-empty"] = &DistributedPortgroupData{Name: "no-vlan"}
	nbHost := &objects.Device{Name: "esx-no-vlan"}

	card := &types.VirtualEthernetCard{
		VirtualDevice: types.VirtualDevice{
			DeviceInfo: &types.Description{Label: "Network adapter 1"},
			Backing: &types.VirtualEthernetCardDistributedVirtualPortBackingInfo{
				Port: types.DistributedVirtualSwitchPortConnection{PortgroupKey: "dvportgroup-empty"},
			},
			Connectable: &types.VirtualDeviceConnectInfo{},
		},
	}
	netboxVM := &objects.VM{Name: "no-vlan-vm", Host: nbHost}
	_, _, vmIface, _, err := vc.collectVMInterfaceData(nbi, netboxVM, mo.VirtualMachine{Guest: &types.GuestInfo{}}, card)
	if err != nil {
		t.Fatalf("collectVMInterfaceData() error = %v", err)
	}
	if vmIface.Mode != nil {
		t.Errorf("vm interface mode = %v, want none without VLAN", vmIface.Mode)
	}

	vnic := types.HostVirtualNic{
		Device: "vmk1",
		Spec: types.HostVirtualNicSpec{
			DistributedVirtualPort: &types.DistributedVirtualSwitchPortConnection{PortgroupKey: "dvportgroup-empty"},
		},
	}
	host := newTestHost("host-9", "esx-no-vlan")
	if _, _, err := vc.collectHostVirtualNicData(nbi, nbHost, host, vnic); err != nil {
		t.Fatalf("collectHostVirtualNicData() error = %v", err)
	}
}
