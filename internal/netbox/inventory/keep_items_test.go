package inventory

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestKeepVMNetworkObjects(t *testing.T) {
	ssotTag := &objects.Tag{ID: 1, Name: constants.SsotTagName}
	vm := &objects.VM{NetboxObject: objects.NetboxObject{ID: 50, Tags: []*objects.Tag{ssotTag}}, Name: "backup-proxy"}
	otherVM := &objects.VM{NetboxObject: objects.NetboxObject{ID: 51, Tags: []*objects.Tag{ssotTag}}, Name: "other-vm"}
	iface := &objects.VMInterface{
		NetboxObject: objects.NetboxObject{ID: 60, Tags: []*objects.Tag{ssotTag}}, Name: "eth0", VM: vm,
	}
	otherIface := &objects.VMInterface{
		NetboxObject: objects.NetboxObject{ID: 61, Tags: []*objects.Tag{ssotTag}}, Name: "eth0", VM: otherVM,
	}
	ip := &objects.IPAddress{
		NetboxObject:       objects.NetboxObject{ID: 70, Tags: []*objects.Tag{ssotTag}},
		Address:            "198.51.100.5/27",
		AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
		AssignedObjectID:   60,
	}
	mac := &objects.MACAddress{
		NetboxObject:       objects.NetboxObject{ID: 80, Tags: []*objects.Tag{ssotTag}},
		MAC:                "BC:24:11:00:00:01",
		AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
		AssignedObjectID:   60,
	}
	prefix := &objects.Prefix{
		NetboxObject: objects.NetboxObject{ID: 90, Tags: []*objects.Tag{ssotTag}}, Prefix: "198.51.100.0/27",
	}
	vmType := constants.ContentTypeVirtualizationVirtualMachine

	nbi := &NetboxInventory{
		Logger:        mockLogger,
		OrphanManager: NewOrphanManager(mockLogger),
		vmInterfacesIndexByVMIdAndName: map[int]map[string]*objects.VMInterface{
			50: {"eth0": iface},
			51: {"eth0": otherIface},
		},
		ipAddressesIndex: map[constants.ContentType]map[string]map[string]map[string]*objects.IPAddress{
			vmType: {"eth0": {"backup-proxy": {"198.51.100.5/27": ip}}},
		},
		macAddressesIndex: map[constants.ContentType]map[string]map[string]map[string]*objects.MACAddress{
			vmType: {"eth0": {"backup-proxy": {"BC:24:11:00:00:01": mac}}},
		},
		prefixesIndexByPrefix: map[string]map[int]*objects.Prefix{"198.51.100.0/27": {0: prefix}},
	}
	for _, item := range []objects.OrphanItem{iface, otherIface, ip, mac, prefix} {
		nbi.OrphanManager.AddItem(item)
	}

	nbi.KeepVMNetworkObjects(vm)

	kept := []struct {
		path constants.APIPath
		id   int
	}{
		{constants.VMInterfacesAPIPath, 60},
		{constants.IPAddressesAPIPath, 70},
		{constants.MACAddressesAPIPath, 80},
		{constants.PrefixesAPIPath, 90},
	}
	for _, k := range kept {
		if _, orphan := nbi.OrphanManager.Items[k.path][k.id]; orphan {
			t.Errorf("%s %d is still an orphan, want it kept", k.path, k.id)
		}
	}
	if _, orphan := nbi.OrphanManager.Items[constants.VMInterfacesAPIPath][61]; !orphan {
		t.Errorf("interface 61 of another VM was kept, want it left as an orphan")
	}
}

func TestKeepPlatform(t *testing.T) {
	ssotTag := &objects.Tag{ID: 1, Name: constants.SsotTagName}
	platform := &objects.Platform{
		NetboxObject: objects.NetboxObject{ID: 6, Tags: []*objects.Tag{ssotTag}}, Name: "Debian 12",
	}
	nbi := &NetboxInventory{Logger: mockLogger, OrphanManager: NewOrphanManager(mockLogger)}
	nbi.OrphanManager.AddItem(platform)

	nbi.KeepPlatform(platform)

	if _, orphan := nbi.OrphanManager.Items[constants.PlatformsAPIPath][6]; orphan {
		t.Errorf("platform 6 is still an orphan, want it kept")
	}
}

func TestGetVM(t *testing.T) {
	vm := &objects.VM{NetboxObject: objects.NetboxObject{ID: 50}, Name: "backup-proxy"}
	nbi := &NetboxInventory{vmsIndexByNameAndClusterID: map[string]map[int]*objects.VM{"backup-proxy": {3: vm}}}
	if got, ok := nbi.GetVM("backup-proxy", 3); !ok || got != vm {
		t.Errorf("GetVM(backup-proxy, 3) = (%v, %t), want the indexed VM", got, ok)
	}
	if _, ok := nbi.GetVM("backup-proxy", 4); ok {
		t.Errorf("GetVM(backup-proxy, 4) found a VM in another cluster")
	}
}
