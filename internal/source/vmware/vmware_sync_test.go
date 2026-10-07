package vmware

import (
	"context"
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
