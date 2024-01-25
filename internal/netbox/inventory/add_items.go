package inventory

import (
	"fmt"
	"slices"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// AddTag adds the newTag to the local netbox inventory.
func (nbi *NetBoxInventory) AddTag(newTag *objects.Tag) (*objects.Tag, error) {
	existingTagIndex := slices.IndexFunc(nbi.Tags, func(t *objects.Tag) bool {
		return t.Name == newTag.Name
	})
	if existingTagIndex == -1 {
		nbi.Logger.Debug("Tag ", newTag.Name, " does not exist in Netbox. Creating it...")
		createdTag, err := service.Create[objects.Tag](nbi.NetboxApi, fmt.Sprintf("/api/extras/tags/%d/", nbi.Tags[existingTagIndex].Id), newTag)
		if err != nil {
			return nil, err
		}
		nbi.Tags = append(nbi.Tags, createdTag)
		return createdTag, nil
	} else {
		nbi.Logger.Debug("Tag ", newTag.Name, " already exists in Netbox...")
		existingTag := nbi.Tags[existingTagIndex]
		diffMap, err := utils.JsonDiffMapExceptId(newTag, existingTag)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			patchedTag, err := service.Patch[objects.Tag](nbi.NetboxApi, fmt.Sprintf("/api/extras/tags/%d/", existingTag.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.Tags[existingTagIndex] = patchedTag
			return patchedTag, nil
		} else {
			return existingTag, nil
		}
	}
}

// AddContactRole adds the newContactRole to the local netbox inventory.
func (nbi *NetBoxInventory) AddContactRole(newContactRole *objects.ContactRole) error {
	newContactRole.NetboxObject.Tags = []*objects.Tag{nbi.SsotTag}
	if _, ok := nbi.ContactRolesIndexByName[newContactRole.Name]; ok {
		existingContactRole := nbi.ContactRolesIndexByName[newContactRole.Name]
		diffMap, err := utils.JsonDiffMapExceptId(newContactRole, existingContactRole)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Contact role ", newContactRole.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContactRole, err := service.Patch[objects.ContactRole](nbi.NetboxApi, fmt.Sprintf("/api/tenancy/contact-roles/%d/", existingContactRole.Id), diffMap)
			if err != nil {
				return err
			}
			nbi.ContactRolesIndexByName[newContactRole.Name] = patchedContactRole
		} else {
			nbi.Logger.Debug("Contact role ", newContactRole.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Contact role ", newContactRole.Name, " does not exist in Netbox. Creating it...")
		newContactRole, err := service.Create[objects.ContactRole](nbi.NetboxApi, "/api/tenancy/contact-roles/", newContactRole)
		if err != nil {
			return err
		}
		nbi.ContactRolesIndexByName[newContactRole.Name] = newContactRole
	}
	return nil
}

// AddContactGroup adds contact group to the local netbox inventory.
func (nbi *NetBoxInventory) AddContactGroup(newContactGroup *objects.ContactGroup) error {
	if _, ok := nbi.ContactGroupsIndexByName[newContactGroup.Name]; ok {
		existingContactGroup := nbi.ContactGroupsIndexByName[newContactGroup.Name]
		diffMap, err := utils.JsonDiffMapExceptId(newContactGroup, existingContactGroup)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Contact group ", newContactGroup.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContactGroup, err := service.Patch[objects.ContactGroup](nbi.NetboxApi, fmt.Sprintf("/api/tenancy/contact-groups/%d/", existingContactGroup.Id), diffMap)
			if err != nil {
				return err
			}
			nbi.ContactGroupsIndexByName[newContactGroup.Name] = patchedContactGroup
		} else {
			nbi.Logger.Debug("Contact group ", newContactGroup.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Contact group ", newContactGroup.Name, " does not exist in Netbox. Creating it...")
		newContactGroup, err := service.Create[objects.ContactGroup](nbi.NetboxApi, "/api/tenancy/contact-groups/", newContactGroup)
		if err != nil {
			return err
		}
		nbi.ContactGroupsIndexByName[newContactGroup.Name] = newContactGroup
	}
	return nil
}

// AddContact adds a contact to the local netbox inventory.
func (nbi *NetBoxInventory) AddContact(newContact *objects.Contact) error {
	if _, ok := nbi.ContactsIndexByName[newContact.Name]; ok {
		existingContact := nbi.ContactsIndexByName[newContact.Name]
		delete(nbi.OrphanManager["/api/tenancy/contacts/"], existingContact.Id)
		diffMap, err := utils.JsonDiffMapExceptId(newContact, existingContact)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Contact ", newContact.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContact, err := service.Patch[objects.Contact](nbi.NetboxApi, fmt.Sprintf("/api/tenancy/contacts/%d/", existingContact.Id), diffMap)
			if err != nil {
				return err
			}
			nbi.ContactsIndexByName[newContact.Name] = patchedContact
		} else {
			nbi.Logger.Debug("Contact ", newContact.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Contact ", newContact.Name, " does not exist in Netbox. Creating it...")
		newContact, err := service.Create[objects.Contact](nbi.NetboxApi, "/api/tenancy/contacts/", newContact)
		if err != nil {
			return err
		}
		nbi.ContactsIndexByName[newContact.Name] = newContact
	}
	return nil
}

func (nbi *NetBoxInventory) AddCustomField(newCf *objects.CustomField) error {
	if _, ok := nbi.CustomFieldsIndexByName[newCf.Name]; ok {
		existingCf := nbi.CustomFieldsIndexByName[newCf.Name]
		diffMap, err := utils.JsonDiffMapExceptId(newCf, existingCf)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Custom field ", newCf.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedCf, err := service.Patch[objects.CustomField](nbi.NetboxApi, fmt.Sprintf("/api/extras/custom-fields/%d/", existingCf.Id), diffMap)
			if err != nil {
				return err
			}
			nbi.CustomFieldsIndexByName[newCf.Name] = patchedCf
		} else {
			nbi.Logger.Debug("Custom field ", newCf.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Custom field ", newCf.Name, " does not exist in Netbox. Creating it...")
		newCf, err := service.Create[objects.CustomField](nbi.NetboxApi, "/api/extras/custom-fields/", newCf)
		if err != nil {
			return err
		}
		nbi.CustomFieldsIndexByName[newCf.Name] = newCf
	}
	return nil
}

func (nbi *NetBoxInventory) AddClusterGroup(newCg *objects.ClusterGroup) (*objects.ClusterGroup, error) {
	newCg.Tags = append(newCg.Tags, nbi.SsotTag)
	if _, ok := nbi.ClusterGroupsIndexByName[newCg.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager["/api/virtualization/cluster-groups/"], nbi.ClusterGroupsIndexByName[newCg.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newCg, nbi.ClusterGroupsIndexByName[newCg.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Cluster group ", newCg.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedCg, err := service.Patch[objects.ClusterGroup](nbi.NetboxApi, fmt.Sprintf("/api/virtualization/cluster-groups/%d/", nbi.ClusterGroupsIndexByName[newCg.Name].Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ClusterGroupsIndexByName[newCg.Name] = patchedCg
		} else {
			nbi.Logger.Debug("Cluster group ", newCg.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Cluster group ", newCg.Name, " does not exist in Netbox. Creating it...")
		newCg, err := service.Create[objects.ClusterGroup](nbi.NetboxApi, "/api/virtualization/cluster-groups/", newCg)
		if err != nil {
			return nil, err
		}
		nbi.ClusterGroupsIndexByName[newCg.Name] = newCg
	}
	// Delete id from orphan manager
	return nbi.ClusterGroupsIndexByName[newCg.Name], nil
}

func (nbi *NetBoxInventory) AddClusterType(newClusterType *objects.ClusterType) (*objects.ClusterType, error) {
	newClusterType.Tags = append(newClusterType.Tags, nbi.SsotTag)
	if _, ok := nbi.ClusterTypesIndexByName[newClusterType.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldClusterType := nbi.ClusterTypesIndexByName[newClusterType.Name]
		delete(nbi.OrphanManager["/api/virtualization/cluster-types/"], nbi.ClusterTypesIndexByName[newClusterType.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newClusterType, nbi.ClusterTypesIndexByName[newClusterType.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedClusterType, err := service.Patch[objects.ClusterType](nbi.NetboxApi, fmt.Sprintf("/api/virtualization/cluster-types/%d/", oldClusterType.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ClusterTypesIndexByName[newClusterType.Name] = patchedClusterType
			return patchedClusterType, nil
		} else {
			nbi.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in Netbox and is up to date...")
			existingClusterType := nbi.ClusterTypesIndexByName[newClusterType.Name]
			return existingClusterType, nil
		}
	} else {
		nbi.Logger.Debug("Cluster type ", newClusterType.Name, " does not exist in Netbox. Creating it...")
		newClusterType, err := service.Create[objects.ClusterType](nbi.NetboxApi, "/api/virtualization/cluster-types/", newClusterType)
		if err != nil {
			return nil, err
		}
		nbi.ClusterTypesIndexByName[newClusterType.Name] = newClusterType
		return newClusterType, nil
	}
}

func (nbi *NetBoxInventory) AddCluster(newCluster *objects.Cluster) error {
	newCluster.Tags = append(newCluster.Tags, nbi.SsotTag)
	if _, ok := nbi.ClustersIndexByName[newCluster.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldCluster := nbi.ClustersIndexByName[newCluster.Name]
		delete(nbi.OrphanManager["/api/virtualization/clusters/"], oldCluster.Id)
		diffMap, err := utils.JsonDiffMapExceptId(newCluster, oldCluster)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Cluster ", newCluster.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedCluster, err := service.Patch[objects.Cluster](nbi.NetboxApi, fmt.Sprintf("/api/virtualization/clusters/%d/", oldCluster.Id), diffMap)
			if err != nil {
				return err
			}
			nbi.ClustersIndexByName[newCluster.Name] = patchedCluster
		} else {
			nbi.Logger.Debug("Cluster ", newCluster.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Cluster ", newCluster.Name, " does not exist in Netbox. Creating it...")
		newCluster, err := service.Create[objects.Cluster](nbi.NetboxApi, "/api/virtualization/clusters/", newCluster)
		if err != nil {
			return err
		}
		nbi.ClustersIndexByName[newCluster.Name] = newCluster
	}
	return nil
}

func (nbi *NetBoxInventory) AddDeviceRole(newDeviceRole *objects.DeviceRole) error {
	newDeviceRole.Tags = append(newDeviceRole.Tags, nbi.SsotTag)
	if _, ok := nbi.DeviceRolesIndexByName[newDeviceRole.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldDeviceRole := nbi.DeviceRolesIndexByName[newDeviceRole.Name]
		delete(nbi.OrphanManager["/api/dcim/device-roles/"], nbi.DeviceRolesIndexByName[newDeviceRole.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newDeviceRole, oldDeviceRole)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedDeviceRole, err := service.Patch[objects.DeviceRole](nbi.NetboxApi, fmt.Sprintf("/api/dcim/device-roles/%d/", oldDeviceRole.Id), diffMap)
			if err != nil {
				return err
			}
			nbi.DeviceRolesIndexByName[newDeviceRole.Name] = patchedDeviceRole
		} else {
			nbi.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Device role ", newDeviceRole.Name, " does not exist in Netbox. Creating it...")
		newDeviceRole, err := service.Create[objects.DeviceRole](nbi.NetboxApi, "/api/dcim/device-roles/", newDeviceRole)
		if err != nil {
			return err
		}
		nbi.DeviceRolesIndexByName[newDeviceRole.Name] = newDeviceRole
	}
	return nil
}

func (nbi *NetBoxInventory) AddManufacturer(newManufacturer *objects.Manufacturer) (*objects.Manufacturer, error) {
	newManufacturer.Tags = append(newManufacturer.Tags, nbi.SsotTag)
	if _, ok := nbi.ManufacturersIndexByName[newManufacturer.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldManufacturer := nbi.ManufacturersIndexByName[newManufacturer.Name]
		delete(nbi.OrphanManager["/api/dcim/manufacturers/"], oldManufacturer.Id)
		diffMap, err := utils.JsonDiffMapExceptId(newManufacturer, oldManufacturer)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Manufacturer ", newManufacturer.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedManufacturer, err := service.Patch[objects.Manufacturer](nbi.NetboxApi, fmt.Sprintf("/api/dcim/manufacturers/%d/", oldManufacturer.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ManufacturersIndexByName[newManufacturer.Name] = patchedManufacturer
		} else {
			nbi.Logger.Debug("Manufacturer ", newManufacturer.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Manufacturer ", newManufacturer.Name, " does not exist in Netbox. Creating it...")
		newManufacturer, err := service.Create[objects.Manufacturer](nbi.NetboxApi, "/api/dcim/manufacturers/", newManufacturer)
		if err != nil {
			return nil, err
		}
		nbi.ManufacturersIndexByName[newManufacturer.Name] = newManufacturer
	}
	return nbi.ManufacturersIndexByName[newManufacturer.Name], nil
}

func (nbi *NetBoxInventory) AddDeviceType(newDeviceType *objects.DeviceType) (*objects.DeviceType, error) {
	newDeviceType.Tags = append(newDeviceType.Tags, nbi.SsotTag)
	if _, ok := nbi.DeviceTypesIndexByModel[newDeviceType.Model]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldDeviceType := nbi.DeviceTypesIndexByModel[newDeviceType.Model]
		delete(nbi.OrphanManager["/api/dcim/device-types/"], oldDeviceType.Id)
		diffMap, err := utils.JsonDiffMapExceptId(newDeviceType, oldDeviceType)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Device type ", newDeviceType.Model, " already exists in Netbox but is out of date. Patching it...")
			patchedDeviceType, err := service.Patch[objects.DeviceType](nbi.NetboxApi, fmt.Sprintf("/api/dcim/device-types/%d/", oldDeviceType.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.DeviceTypesIndexByModel[newDeviceType.Model] = patchedDeviceType
		} else {
			nbi.Logger.Debug("Device type ", newDeviceType.Model, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Device type ", newDeviceType.Model, " does not exist in Netbox. Creating it...")
		newDeviceType, err := service.Create[objects.DeviceType](nbi.NetboxApi, "/api/dcim/device-types/", newDeviceType)
		if err != nil {
			return nil, err
		}
		nbi.DeviceTypesIndexByModel[newDeviceType.Model] = newDeviceType
	}
	return nbi.DeviceTypesIndexByModel[newDeviceType.Model], nil
}

func (nbi *NetBoxInventory) AddPlatform(newPlatform *objects.Platform) (*objects.Platform, error) {
	newPlatform.Tags = append(newPlatform.Tags, nbi.SsotTag)
	if _, ok := nbi.PlatformsIndexByName[newPlatform.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager["/api/dcim/platforms/"], nbi.PlatformsIndexByName[newPlatform.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newPlatform, nbi.PlatformsIndexByName[newPlatform.Name])
		oldPlatform := nbi.PlatformsIndexByName[newPlatform.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Platform ", newPlatform.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedPlatform, err := service.Patch[objects.Platform](nbi.NetboxApi, fmt.Sprintf("/api/dcim/platforms/%d/", oldPlatform.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.PlatformsIndexByName[newPlatform.Name] = patchedPlatform
		} else {
			nbi.Logger.Debug("Platform ", newPlatform.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Platform ", newPlatform.Name, " does not exist in Netbox. Creating it...")
		newPlatform, err := service.Create[objects.Platform](nbi.NetboxApi, "/api/dcim/platforms/", newPlatform)
		if err != nil {
			return nil, err
		}
		nbi.PlatformsIndexByName[newPlatform.Name] = newPlatform
	}
	return nbi.PlatformsIndexByName[newPlatform.Name], nil
}

func (nbi *NetBoxInventory) AddDevice(newDevice *objects.Device) (*objects.Device, error) {
	newDevice.Tags = append(newDevice.Tags, nbi.SsotTag)
	if _, ok := nbi.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id]; ok {
		oldDevice := nbi.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id]
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager["/api/dcim/devices/"], nbi.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newDevice, nbi.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Device ", newDevice.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedDevice, err := service.Patch[objects.Device](nbi.NetboxApi, fmt.Sprintf("/api/dcim/devices/%d/", oldDevice.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id] = patchedDevice
		} else {
			nbi.Logger.Debug("Device ", newDevice.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Device ", newDevice.Name, " does not exist in Netbox. Creating it...")
		newDevice, err := service.Create[objects.Device](nbi.NetboxApi, "/api/dcim/devices/", newDevice)
		if err != nil {
			return nil, err
		}
		if nbi.DevicesIndexByNameAndSiteId[newDevice.Name] == nil {
			nbi.DevicesIndexByNameAndSiteId[newDevice.Name] = make(map[int]*objects.Device)
		}
		nbi.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id] = newDevice
	}
	return nbi.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id], nil
}

func (nbi *NetBoxInventory) AddVlanGroup(newVlanGroup *objects.VlanGroup) (*objects.VlanGroup, error) {
	newVlanGroup.Tags = append(newVlanGroup.Tags, nbi.SsotTag)
	if _, ok := nbi.VlanGroupsIndexByName[newVlanGroup.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldVlanGroup := nbi.VlanGroupsIndexByName[newVlanGroup.Name]
		delete(nbi.OrphanManager["/api/ipam/vlan-groups/"], oldVlanGroup.Id)
		diffMap, err := utils.JsonDiffMapExceptId(newVlanGroup, oldVlanGroup)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("VlanGroup ", newVlanGroup.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVlanGroup, err := service.Patch[objects.VlanGroup](nbi.NetboxApi, fmt.Sprintf("/api/ipam/vlan-groups/%d/", oldVlanGroup.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VlanGroupsIndexByName[newVlanGroup.Name] = patchedVlanGroup
		} else {
			nbi.Logger.Debug("Vlan ", newVlanGroup.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Vlan ", newVlanGroup.Name, " does not exist in Netbox. Creating it...")
		newVlan, err := service.Create[objects.VlanGroup](nbi.NetboxApi, "/api/ipam/vlan-groups/", newVlanGroup)
		if err != nil {
			return nil, err
		}
		nbi.VlanGroupsIndexByName[newVlan.Name] = newVlan
	}
	return nbi.VlanGroupsIndexByName[newVlanGroup.Name], nil
}

func (nbi *NetBoxInventory) AddVlan(newVlan *objects.Vlan) (*objects.Vlan, error) {
	newVlan.Tags = append(newVlan.Tags, nbi.SsotTag)
	if _, ok := nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id][newVlan.Vid]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldVlan := nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id][newVlan.Vid]
		delete(nbi.OrphanManager["/api/ipam/vlans/"], nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id][newVlan.Vid].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newVlan, nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id][newVlan.Vid])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Vlan ", newVlan.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVlan, err := service.Patch[objects.Vlan](nbi.NetboxApi, fmt.Sprintf("/api/ipam/vlans/%d/", oldVlan.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id][newVlan.Vid] = patchedVlan
		} else {
			nbi.Logger.Debug("Vlan ", newVlan.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Vlan ", newVlan.Name, " does not exist in Netbox. Creating it...")
		newVlan, err := service.Create[objects.Vlan](nbi.NetboxApi, "/api/ipam/vlans/", newVlan)
		if err != nil {
			return nil, err
		}
		if nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id] == nil {
			nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id] = make(map[int]*objects.Vlan)
		}
		nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id][newVlan.Vid] = newVlan
	}
	return nbi.VlansIndexByVlanGroupIdAndVid[newVlan.Group.Id][newVlan.Vid], nil
}

func (nbi *NetBoxInventory) AddInterface(newInterface *objects.Interface) (*objects.Interface, error) {
	newInterface.Tags = append(newInterface.Tags, nbi.SsotTag)
	if _, ok := nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager["/api/dcim/interfaces/"], nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newInterface, nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name])
		oldIntf := nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Interface ", newInterface.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedInterface, err := service.Patch[objects.Interface](nbi.NetboxApi, fmt.Sprintf("/api/dcim/interfaces/%d/", oldIntf.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name] = patchedInterface
		} else {
			nbi.Logger.Debug("Interface ", newInterface.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Interface ", newInterface.Name, " does not exist in Netbox. Creating it...")
		newInterface, err := service.Create[objects.Interface](nbi.NetboxApi, "/api/dcim/interfaces/", newInterface)
		if err != nil {
			return nil, err
		}
		if nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id] == nil {
			nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id] = make(map[string]*objects.Interface)
		}
		nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name] = newInterface
	}
	return nbi.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name], nil
}

func (nbi *NetBoxInventory) AddVM(newVm *objects.VM) (*objects.VM, error) {
	newVm.Tags = append(newVm.Tags, nbi.SsotTag)
	if _, ok := nbi.VMsIndexByName[newVm.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager["/api/virtualization/virtual-machines/"], nbi.VMsIndexByName[newVm.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newVm, nbi.VMsIndexByName[newVm.Name])
		oldVm := nbi.VMsIndexByName[newVm.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("VM ", newVm.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVm, err := service.Patch[objects.VM](nbi.NetboxApi, fmt.Sprintf("/api/virtualization/virtual-machines/%d/", oldVm.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VMsIndexByName[newVm.Name] = patchedVm
		} else {
			nbi.Logger.Debug("VM ", newVm.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("VM ", newVm.Name, " does not exist in Netbox. Creating it...")
		newVm, err := service.Create[objects.VM](nbi.NetboxApi, "/api/virtualization/virtual-machines/", newVm)
		if err != nil {
			return nil, err
		}
		nbi.VMsIndexByName[newVm.Name] = newVm
		return newVm, nil
	}
	return nbi.VMsIndexByName[newVm.Name], nil
}

func (nbi *NetBoxInventory) AddVMInterface(newVMInterface *objects.VMInterface) (*objects.VMInterface, error) {
	newVMInterface.Tags = append(newVMInterface.Tags, nbi.SsotTag)
	if _, ok := nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager["/api/virtualization/interfaces/"], nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newVMInterface, nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name])
		oldVmIntf := nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("VM interface ", newVMInterface.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVMInterface, err := service.Patch[objects.VMInterface](nbi.NetboxApi, fmt.Sprintf("/api/virtualization/interfaces/%d/", oldVmIntf.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name] = patchedVMInterface
		} else {
			nbi.Logger.Debug("VM interface ", newVMInterface.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("VM interface ", newVMInterface.Name, " does not exist in Netbox. Creating it...")
		newVMInterface, err := service.Create[objects.VMInterface](nbi.NetboxApi, "/api/virtualization/interfaces/", newVMInterface)
		if err != nil {
			return nil, err
		}
		if nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id] == nil {
			nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id] = make(map[string]*objects.VMInterface)
		}
		nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name] = newVMInterface
	}
	return nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name], nil
}

func (nbi *NetBoxInventory) AddIPAddress(newIPAddress *objects.IPAddress) (*objects.IPAddress, error) {
	newIPAddress.Tags = append(newIPAddress.Tags, nbi.SsotTag)
	if _, ok := nbi.IPAdressesIndexByAddress[newIPAddress.Address]; ok {
		// Delete id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager["/api/ipam/ip-addresses/"], nbi.IPAdressesIndexByAddress[newIPAddress.Address].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newIPAddress, nbi.IPAdressesIndexByAddress[newIPAddress.Address])
		oldIpAddress := nbi.IPAdressesIndexByAddress[newIPAddress.Address]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("IP address ", newIPAddress.Address, " already exists in Netbox but is out of date. Patching it...")
			patchedIPAddress, err := service.Patch[objects.IPAddress](nbi.NetboxApi, fmt.Sprintf("/api/ipam/ip-addresses/%d/", oldIpAddress.Id), diffMap)
			if err != nil {
				return nil, err
			}
			nbi.IPAdressesIndexByAddress[newIPAddress.Address] = patchedIPAddress
		} else {
			nbi.Logger.Debug("IP address ", newIPAddress.Address, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("IP address ", newIPAddress.Address, " does not exist in Netbox. Creating it...")
		newIPAddress, err := service.Create[objects.IPAddress](nbi.NetboxApi, "/api/ipam/ip-addresses/", newIPAddress)
		if err != nil {
			return nil, err
		}
		nbi.IPAdressesIndexByAddress[newIPAddress.Address] = newIPAddress
	}
	return nbi.IPAdressesIndexByAddress[newIPAddress.Address], nil
}
