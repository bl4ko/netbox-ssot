package inventory

import (
	"fmt"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

func (nbi *NetboxInventory) DeleteOrphans(hard bool) error {
	for i := 0; i < len(nbi.OrphanManager.OrphanObjectPriority); i++ {
		deleteTypeStr := "soft"
		if hard {
			deleteTypeStr = "hard"
		}
		objectAPIPath := nbi.OrphanManager.OrphanObjectPriority[i]
		id2orphanItem := nbi.OrphanManager.Items[objectAPIPath]
		if len(id2orphanItem) == 0 {
			continue
		}

		nbi.OrphanManager.Logger.Infof(nbi.Ctx, "Performing %s deletion of orphaned objects of type %s", deleteTypeStr, objectAPIPath)
		nbi.OrphanManager.Logger.Debugf(nbi.Ctx, "IDs of objects to be %s deleted: %v", deleteTypeStr, id2orphanItem)

		for _, orphanItem := range id2orphanItem {
			if hard {
				// Perform hard deletion
				err := nbi.hardDelete(objectAPIPath, orphanItem)
				if err != nil {
					nbi.OrphanManager.Logger.Errorf(nbi.Ctx, "hard delete object: %s", err)
					continue
				}
			} else {
				err := nbi.softDelete(objectAPIPath, orphanItem)
				if err != nil {
					nbi.OrphanManager.Logger.Errorf(nbi.Ctx, "soft delete object: %s", err)
				}
			}
		}
	}

	return nil
}

func (nbi *NetboxInventory) hardDelete(apiPath string, orphanItem objects.OrphanItem) error {
	// Perform hard deletion
	err := nbi.NetboxAPI.DeleteObject(nbi.Ctx, apiPath, orphanItem.GetID())
	if err != nil {
		return fmt.Errorf("Failed deleting %s object: %s", orphanItem, err)
	}
	return nil
}

func (nbi *NetboxInventory) softDelete(apiPath string, orphanItem objects.OrphanItem) error {
	// Perform soft deletion
	// Add tag to the object to mark it as orphaned
	todayDate := time.Now().Format(constants.CustomFieldOrphanLastSeenFormat)
	if !orphanItem.GetNetboxObject().HasTag(nbi.OrphanManager.Tag) {
		// This OrphanItem has been marked as orphan for the first time
		orphanItem.GetNetboxObject().AddTag(nbi.OrphanManager.Tag)
		orphanItem.GetNetboxObject().SetCustomField(constants.CustomFieldOrphanLastSeenName, todayDate)
		diffMap := utils.ExtractFieldsFromDiffMap(utils.StructToNetboxJSONMap(orphanItem.GetNetboxObject()), []string{"tags", "custom_fields"})
		// Update object on the API
		var err error
		switch orphanItem.(type) {
		case *objects.VlanGroup:
			_, err = service.Patch[objects.VlanGroup](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.Prefix:
			_, err = service.Patch[objects.Prefix](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.Vlan:
			_, err = service.Patch[objects.Vlan](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.IPAddress:
			_, err = service.Patch[objects.IPAddress](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.VirtualDeviceContext:
			_, err = service.Patch[objects.VirtualDeviceContext](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.Interface:
			_, err = service.Patch[objects.Interface](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.VMInterface:
			_, err = service.Patch[objects.VMInterface](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.VM:
			_, err = service.Patch[objects.VM](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.Device:
			_, err = service.Patch[objects.Device](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.Platform:
			_, err = service.Patch[objects.Platform](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.DeviceType:
			_, err = service.Patch[objects.DeviceType](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.Manufacturer:
			_, err = service.Patch[objects.Manufacturer](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.DeviceRole:
			_, err = service.Patch[objects.DeviceRole](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.ClusterType:
			_, err = service.Patch[objects.ClusterType](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.Cluster:
			_, err = service.Patch[objects.Cluster](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.ClusterGroup:
			_, err = service.Patch[objects.ClusterGroup](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.ContactAssignment:
			_, err = service.Patch[objects.ContactAssignment](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.Contact:
			_, err = service.Patch[objects.Contact](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.WirelessLAN:
			_, err = service.Patch[objects.WirelessLAN](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		case *objects.WirelessLANGroup:
			_, err = service.Patch[objects.WirelessLANGroup](nbi.OrphanManager.Ctx, nbi.NetboxAPI, orphanItem.GetID(), diffMap)
		default:
			return fmt.Errorf("unsupported type for orphan item%T", orphanItem)
		}
		if err != nil {
			return fmt.Errorf("failed updating %s object with orphan tag: %s", orphanItem, err)
		}
	} else {
		nbi.Logger.Debugf(nbi.Ctx, "%s is already marked as orphan", orphanItem)
		if orphanItem.GetNetboxObject().GetCustomField(constants.CustomFieldOrphanLastSeenName) == "" {
			orphanItem.GetNetboxObject().SetCustomField(constants.CustomFieldOrphanLastSeenName, todayDate)
		} else {
			// We check if the object has been orphaned for more than nbi.Config.OrphanRemoveAfterDays
			lastSeen, err := time.Parse(constants.CustomFieldOrphanLastSeenFormat, orphanItem.GetNetboxObject().GetCustomField(constants.CustomFieldOrphanLastSeenName).(string))
			if err != nil {
				return fmt.Errorf("failed parsing last seen date: %s", err)
			}
			if int((time.Since(lastSeen).Hours())/24) > nbi.NetboxConfig.RemoveOrphansAfterDays { //nolint:mnd
				// Perform hard deletion
				err := nbi.hardDelete(apiPath, orphanItem)
				if err != nil {
					return fmt.Errorf("failed deleting %s object: %s", orphanItem, err)
				}
			}
		}
	}
	return nil
}
