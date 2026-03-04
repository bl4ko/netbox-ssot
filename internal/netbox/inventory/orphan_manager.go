package inventory

import (
	"context"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/logger"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
)

type OrphanManager struct {
	// ItemsSet is a map of objectAPIPath to a set of managed ids for that object type.
	//
	// {
	//		"/api/dcim/devices/": {22: true, 3: true, ...},
	//		"/api/dcim/interface/": {15: true, 36: true, ...},
	//  	"/api/virtualization/clusters/": {121: true, 122: true, ...},
	//  	"...": [...]
	// }
	//
	// It stores which objects have been created by netbox-ssot and can be deleted
	// because they are not available in the sources anymore
	Items map[constants.APIPath]map[int]objects.OrphanItem
	// OrphanObjectPriority is a map that stores priorities for each object. This is necessary
	// because map order is non deterministic and if we delete dependent object first we will
	// get the dependency error.
	//
	// {
	//   0: service.TagApiPath,
	//   1: service.CustomFieldApiPath,
	//   ...
	// }
	OrphanObjectPriority map[int]constants.APIPath
	// Tag for orphaned objects. Initialized in initTags.
	Tag *objects.Tag
	// Logger for orphan manager
	Logger *logger.Logger
	// Context for orphan manager
	Ctx context.Context
}

func NewOrphanManager(logger *logger.Logger) *OrphanManager {
	// Starts with 0 for easier integration with for loops
	orphanObjectPriority := map[int]constants.APIPath{
		0:  constants.VlanGroupsAPIPath,
		1:  constants.PrefixesAPIPath,
		2:  constants.VlansAPIPath,
		3:  constants.IPAddressesAPIPath,
		4:  constants.VirtualDeviceContextsAPIPath,
		5:  constants.InterfacesAPIPath,
		6:  constants.VMInterfacesAPIPath,
		7:  constants.VirtualDisksAPIPath,
		8:  constants.VirtualMachinesAPIPath,
		9:  constants.DevicesAPIPath,
		10: constants.PlatformsAPIPath,
		11: constants.DeviceTypesAPIPath,
		12: constants.ManufacturersAPIPath,
		13: constants.DeviceRolesAPIPath,
		14: constants.ClustersAPIPath,
		15: constants.ClusterTypesAPIPath,
		16: constants.ClusterGroupsAPIPath,
		17: constants.ContactAssignmentsAPIPath,
		18: constants.ContactsAPIPath,
		19: constants.WirelessLANsAPIPath,
		20: constants.WirelessLANGroupsAPIPath,
		21: constants.MACAddressesAPIPath,
		22: constants.VRFsAPIPath,
	}
	orphanCtx := context.WithValue(context.Background(), constants.CtxSourceKey, "orphanManager")

	return &OrphanManager{
		Items:                map[constants.APIPath]map[int]objects.OrphanItem{},
		OrphanObjectPriority: orphanObjectPriority,
		Logger:               logger,
		Ctx:                  orphanCtx,
	}
}

func (orphanManager *OrphanManager) AddItem(orphanItem objects.OrphanItem) {
	// Manage only objects created with netbox-ssot tag
	netboxObject := orphanItem.GetNetboxObject()
	if netboxObject.HasTagByName(constants.SsotTagName) {
		if orphanManager.Items[orphanItem.GetAPIPath()] == nil {
			orphanManager.Items[orphanItem.GetAPIPath()] = map[int]objects.OrphanItem{}
		}
		orphanManager.Items[orphanItem.GetAPIPath()][netboxObject.ID] = orphanItem
	}
}

func (orphanManager *OrphanManager) RemoveItem(obj objects.OrphanItem) {
	delete(orphanManager.Items[obj.GetAPIPath()], obj.GetID())
}
