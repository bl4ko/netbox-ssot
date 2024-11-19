package inventory

import (
	"context"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

type OrphanManager struct {
	// Items is a map of objectAPIPath to a set of managed ids for that object type.
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
	Items map[string]map[int]bool
	// OrphanObjectPriority is a map that stores priorities for each object. This is necessary
	// because map order is non deterministic and if we delete dependent object first we will
	// get the dependency error.
	//
	// {
	//   0: service.TagApiPath,
	//   1: service.CustomFieldApiPath,
	//   ...
	// }
	OrphanObjectPriority map[int]string
	// Tag for orphaned objects. Initialized in initTags.
	Tag *objects.Tag
	// Logger for orphan manager
	Logger *logger.Logger
	// Context for orphan manager
	Ctx context.Context
}

func NewOrphanManager(logger *logger.Logger) *OrphanManager {
	// Starts with 0 for easier integration with for loops
	orphanObjectPriority := map[int]string{
		0:  constants.VlanGroupsAPIPath,
		1:  constants.PrefixesAPIPath,
		2:  constants.VlansAPIPath,
		3:  constants.IPAddressesAPIPath,
		4:  constants.VirtualDeviceContextsAPIPath,
		5:  constants.InterfacesAPIPath,
		6:  constants.VMInterfacesAPIPath,
		7:  constants.VirtualMachinesAPIPath,
		8:  constants.DevicesAPIPath,
		9:  constants.PlatformsAPIPath,
		10: constants.DeviceTypesAPIPath,
		11: constants.ManufacturersAPIPath,
		12: constants.DeviceRolesAPIPath,
		13: constants.ClustersAPIPath,
		14: constants.ClusterTypesAPIPath,
		15: constants.ClusterGroupsAPIPath,
		16: constants.ContactAssignmentsAPIPath,
		17: constants.ContactsAPIPath,
		18: constants.WirelessLANsAPIPath,
		19: constants.WirelessLANGroupsAPIPath,
	}
	orphanCtx := context.WithValue(context.Background(), constants.CtxSourceKey, "orphanManager")

	return &OrphanManager{
		Items:                map[string]map[int]bool{},
		OrphanObjectPriority: orphanObjectPriority,
		Logger:               logger,
		Ctx:                  orphanCtx,
	}
}

func (orphanManager *OrphanManager) AddItem(itemAPIPath string, netboxObject *objects.NetboxObject) {
	// Manage only objects created with netbox-ssot tag
	if netboxObject.HasTagByName(constants.SsotTagName) {
		if orphanManager.Items[itemAPIPath] == nil {
			orphanManager.Items[itemAPIPath] = map[int]bool{}
		}
		// Add the id to the Items map
		orphanManager.Items[itemAPIPath][netboxObject.ID] = true
		// Add orphan tag to object if it doesn't have it already
		netboxObject.AddTag(orphanManager.Tag)
	}
}

func (orphanManager *OrphanManager) RemoveItem(itemAPIPath string, netboxObject *objects.NetboxObject) {
	delete(orphanManager.Items[itemAPIPath], netboxObject.ID)
	netboxObject.RemoveTag(orphanManager.Tag)
}

// DeleteOrphans deletes orphaned objects from the Netbox API.
// Orphaned object are objects that were collected by OrphanManager.
func (orphanManager *OrphanManager) DeleteOrphans(ctx context.Context, nbi *NetboxInventory) error {
	// Ensure OrphanObjectPriority and OrphanManager lengths are the same,
	// if not, there are missing entries somewhere and need to be fixed.
	if len(orphanManager.Items) != len(orphanManager.OrphanObjectPriority) {
		panic("len(orphanManager) != len(orphanObjectPriority). This should not happen. Every orphan managed object must have its corresponding priority")
	}
	for i := 0; i < len(orphanManager.OrphanObjectPriority); i++ {
		objectAPIPath := orphanManager.OrphanObjectPriority[i]
		ids := orphanManager.Items[objectAPIPath]
		if len(ids) != 0 {
			orphanManager.Logger.Infof(ctx, "Deleting orphaned objects of type %s", objectAPIPath)
			orphanManager.Logger.Debugf(ctx, "Ids of objects to be deleted: %v", ids)
			for id := range ids {
				err := nbi.NetboxAPI.DeleteObject(ctx, objectAPIPath, id)
				if err != nil {
					orphanManager.Logger.Errorf(orphanManager.Ctx, "delete objects: %s", err)
				}
			}
		}
	}
	return nil
}
