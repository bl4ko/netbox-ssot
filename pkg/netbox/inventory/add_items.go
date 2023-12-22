package inventory

import (
	"slices"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/objects"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
)

// Function for adding a new tag to NetBoxInventory
func (ni *NetBoxInventory) AddTag(newTag *objects.Tag) (*objects.Tag, error) {
	existingTagIndex := slices.IndexFunc(ni.Tags, func(t *objects.Tag) bool {
		return t.Name == newTag.Name
	})
	if existingTagIndex == -1 {
		ni.Logger.Debug("Tag ", newTag.Name, " does not exist in NetBox. Creating it...")
		createdTag, err := ni.NetboxApi.CreateTag(newTag)
		if err != nil {
			return nil, err
		}
		ni.Tags = append(ni.Tags, createdTag)
		return createdTag, nil
	} else {
		ni.Logger.Debug("Tag ", newTag.Name, " already exists in NetBox...")
		existingTag := ni.Tags[existingTagIndex]
		diffMap, err := utils.JsonDiffMapExceptId(newTag, existingTag)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			patchedTag, err := ni.NetboxApi.PatchTag(diffMap, existingTag.ID)
			if err != nil {
				return nil, err
			}
			ni.Tags[existingTagIndex] = patchedTag
			return patchedTag, nil
		} else {
			return existingTag, nil
		}
	}
}

// Adding new custom-field to inventory
func (ni *NetBoxInventory) AddCustomField(newCf *objects.CustomField) error {
	if _, ok := ni.CustomFieldsIndexByName[newCf.Name]; ok {
		existingCf := ni.CustomFieldsIndexByName[newCf.Name]
		diffMap, err := utils.JsonDiffMapExceptId(newCf, existingCf)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Custom field ", newCf.Name, " already exists in NetBox but is out of date. Patching it... ")
			patchedCf, err := ni.NetboxApi.PatchCustomField(diffMap, existingCf.ID)
			if err != nil {
				return err
			}
			ni.CustomFieldsIndexByName[newCf.Name] = patchedCf
		} else {
			ni.Logger.Debug("Custom field ", newCf.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Custom field ", newCf.Name, " does not exist in NetBox. Creating it...")
		newCf, err := ni.NetboxApi.CreateCustomField(newCf)
		if err != nil {
			return err
		}
		ni.CustomFieldsIndexByName[newCf.Name] = newCf
	}
	return nil
}

// Add Cluster to NetBoxInventory
func (ni *NetBoxInventory) AddClusterGroup(newCg *objects.ClusterGroup, newTags []*objects.Tag) error {
	newCg.Tags = append(newCg.Tags, ni.SsotTag)
	if _, ok := ni.ClusterGroupsIndexByName[newCg.Name]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newCg, ni.ClusterGroupsIndexByName[newCg.Name])
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster group ", newCg.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedCg, err := ni.NetboxApi.PatchClusterGroup(diffMap, ni.ClusterGroupsIndexByName[newCg.Name].ID)
			if err != nil {
				return err
			}
			ni.ClusterGroupsIndexByName[newCg.Name] = patchedCg
		} else {
			ni.Logger.Debug("Cluster group ", newCg.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Cluster group ", newCg.Name, " does not exist in NetBox. Creating it...")
		newCg, err := ni.NetboxApi.CreateClusterGroup(newCg)
		if err != nil {
			return err
		}
		ni.ClusterGroupsIndexByName[newCg.Name] = newCg
	}
	return nil
}

// Add ClusterType to NetBoxInventory
func (ni *NetBoxInventory) AddClusterType(newClusterType *objects.ClusterType) (*objects.ClusterType, error) {
	newClusterType.Tags = append(newClusterType.Tags, ni.SsotTag)
	if _, ok := ni.ClusterTypesIndexByName[newClusterType.Name]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newClusterType, ni.ClusterTypesIndexByName[newClusterType.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedClusterType, err := ni.NetboxApi.PatchClusterType(diffMap, ni.ClusterTypesIndexByName[newClusterType.Name].ID)
			if err != nil {
				return nil, err
			}
			ni.ClusterTypesIndexByName[newClusterType.Name] = patchedClusterType
			return patchedClusterType, nil
		} else {
			ni.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in NetBox and is up to date...")
			existingClusterType := ni.ClusterTypesIndexByName[newClusterType.Name]
			return existingClusterType, nil
		}
	} else {
		ni.Logger.Debug("Cluster type ", newClusterType.Name, " does not exist in NetBox. Creating it...")
		newClusterType, err := ni.NetboxApi.CreateClusterType(newClusterType)
		if err != nil {
			return nil, err
		}
		ni.ClusterTypesIndexByName[newClusterType.Name] = newClusterType
		return newClusterType, nil
	}
}

func (ni *NetBoxInventory) AddCluster(newCluster *objects.Cluster) error {
	newCluster.Tags = append(newCluster.Tags, ni.SsotTag)
	if _, ok := ni.ClustersIndexByName[newCluster.Name]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newCluster, ni.ClustersIndexByName[newCluster.Name])
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster ", newCluster.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedCluster, err := ni.NetboxApi.PatchCluster(diffMap, ni.ClustersIndexByName[newCluster.Name].ID)
			if err != nil {
				return err
			}
			ni.ClustersIndexByName[newCluster.Name] = patchedCluster
		} else {
			ni.Logger.Debug("Cluster ", newCluster.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Cluster ", newCluster.Name, " does not exist in NetBox. Creating it...")
		newCluster, err := ni.NetboxApi.CreateCluster(newCluster)
		if err != nil {
			return err
		}
		ni.ClustersIndexByName[newCluster.Name] = newCluster
	}
	return nil
}

func (ni *NetBoxInventory) AddDeviceRole(newDeviceRole *objects.DeviceRole) error {
	newDeviceRole.Tags = append(newDeviceRole.Tags, ni.SsotTag)
	if _, ok := ni.DeviceRolesIndexByName[newDeviceRole.Name]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newDeviceRole, ni.DeviceRolesIndexByName[newDeviceRole.Name])
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedDeviceRole, err := ni.NetboxApi.PatchDeviceRole(diffMap, ni.DeviceRolesIndexByName[newDeviceRole.Name].ID)
			if err != nil {
				return err
			}
			ni.DeviceRolesIndexByName[newDeviceRole.Name] = patchedDeviceRole
		} else {
			ni.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Device role ", newDeviceRole.Name, " does not exist in NetBox. Creating it...")
		newDeviceRole, err := ni.NetboxApi.CreateDeviceRole(newDeviceRole)
		if err != nil {
			return err
		}
		ni.DeviceRolesIndexByName[newDeviceRole.Name] = newDeviceRole
	}
	return nil
}

func (ni *NetBoxInventory) AddManufacturer(newManufacturer *objects.Manufacturer) (*objects.Manufacturer, error) {
	newManufacturer.Tags = append(newManufacturer.Tags, ni.SsotTag)
	if _, ok := ni.ManufacturersIndexByName[newManufacturer.Name]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newManufacturer, ni.ManufacturersIndexByName[newManufacturer.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Manufacturer ", newManufacturer.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedManufacturer, err := ni.NetboxApi.PatchManufacturer(diffMap, ni.ManufacturersIndexByName[newManufacturer.Name].ID)
			if err != nil {
				return nil, err
			}
			ni.ManufacturersIndexByName[newManufacturer.Name] = patchedManufacturer
		} else {
			ni.Logger.Debug("Manufacturer ", newManufacturer.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Manufacturer ", newManufacturer.Name, " does not exist in NetBox. Creating it...")
		newManufacturer, err := ni.NetboxApi.CreateManufacturer(newManufacturer)
		if err != nil {
			return nil, err
		}
		ni.ManufacturersIndexByName[newManufacturer.Name] = newManufacturer
	}
	return ni.ManufacturersIndexByName[newManufacturer.Name], nil
}

func (ni *NetBoxInventory) AddDeviceType(newDeviceType *objects.DeviceType) (*objects.DeviceType, error) {
	newDeviceType.Tags = append(newDeviceType.Tags, ni.SsotTag)
	if _, ok := ni.DeviceTypesIndexByModel[newDeviceType.Model]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newDeviceType, ni.DeviceTypesIndexByModel[newDeviceType.Model])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Device type ", newDeviceType.Model, " already exists in NetBox but is out of date. Patching it...")
			patchedDeviceType, err := ni.NetboxApi.PatchDeviceType(diffMap, ni.DeviceTypesIndexByModel[newDeviceType.Model].ID)
			if err != nil {
				return nil, err
			}
			ni.DeviceTypesIndexByModel[newDeviceType.Model] = patchedDeviceType
		} else {
			ni.Logger.Debug("Device type ", newDeviceType.Model, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Device type ", newDeviceType.Model, " does not exist in NetBox. Creating it...")
		newDeviceType, err := ni.NetboxApi.CreateDeviceType(newDeviceType)
		if err != nil {
			return nil, err
		}
		ni.DeviceTypesIndexByModel[newDeviceType.Model] = newDeviceType
	}
	return ni.DeviceTypesIndexByModel[newDeviceType.Model], nil
}

func (ni *NetBoxInventory) AddPlatform(newPlatform *objects.Platform) (*objects.Platform, error) {
	newPlatform.Tags = append(newPlatform.Tags, ni.SsotTag)
	if _, ok := ni.PlatformsIndexByName[newPlatform.Name]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newPlatform, ni.PlatformsIndexByName[newPlatform.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Platform ", newPlatform.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedPlatform, err := ni.NetboxApi.PatchPlatform(diffMap, ni.PlatformsIndexByName[newPlatform.Name].ID)
			if err != nil {
				return nil, err
			}
			ni.PlatformsIndexByName[newPlatform.Name] = patchedPlatform
		} else {
			ni.Logger.Debug("Platform ", newPlatform.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Platform ", newPlatform.Name, " does not exist in NetBox. Creating it...")
		newPlatform, err := ni.NetboxApi.CreatePlatform(newPlatform)
		if err != nil {
			return nil, err
		}
		ni.PlatformsIndexByName[newPlatform.Name] = newPlatform
	}
	return ni.PlatformsIndexByName[newPlatform.Name], nil
}

func (ni *NetBoxInventory) AddDevice(newDevice *objects.Device) (*objects.Device, error) {
	newDevice.Tags = append(newDevice.Tags, ni.SsotTag)
	if _, ok := ni.DevicesIndexByUuid[newDevice.AssetTag]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newDevice, ni.DevicesIndexByUuid[newDevice.AssetTag])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Device ", newDevice.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedDevice, err := ni.NetboxApi.PatchDevice(diffMap, ni.DevicesIndexByUuid[newDevice.AssetTag].ID)
			if err != nil {
				return nil, err
			}
			ni.DevicesIndexByUuid[newDevice.AssetTag] = patchedDevice
		} else {
			ni.Logger.Debug("Device ", newDevice.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Device ", newDevice.Name, " does not exist in NetBox. Creating it...")
		newDevice, err := ni.NetboxApi.CreateDevice(newDevice)
		if err != nil {
			return nil, err
		}
		ni.DevicesIndexByUuid[newDevice.AssetTag] = newDevice
	}
	return ni.DevicesIndexByUuid[newDevice.AssetTag], nil
}

func (ni *NetBoxInventory) AddVlan(newVlan *objects.Vlan) (*objects.Vlan, error) {
	newVlan.Tags = append(newVlan.Tags, ni.SsotTag)
	if _, ok := ni.VlansIndexByName[newVlan.Name]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newVlan, ni.VlansIndexByName[newVlan.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Vlan ", newVlan.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedVlan, err := ni.NetboxApi.PatchVlan(diffMap, ni.VlansIndexByName[newVlan.Name].ID)
			if err != nil {
				return nil, err
			}
			ni.VlansIndexByName[newVlan.Name] = patchedVlan
		} else {
			ni.Logger.Debug("Vlan ", newVlan.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Vlan ", newVlan.Name, " does not exist in NetBox. Creating it...")
		newVlan, err := ni.NetboxApi.CreateVlan(newVlan)
		if err != nil {
			return nil, err
		}
		ni.VlansIndexByName[newVlan.Name] = newVlan
	}
	return ni.VlansIndexByName[newVlan.Name], nil
}

func (ni *NetBoxInventory) AddInterface(newInterface *objects.Interface) (*objects.Interface, error) {
	newInterface.Tags = append(newInterface.Tags, ni.SsotTag)
	if _, ok := ni.InterfacesIndexByDeviceAndName[newInterface.Device.ID][newInterface.Name]; ok {
		diffMap, err := utils.JsonDiffMapExceptId(newInterface, ni.InterfacesIndexByDeviceAndName[newInterface.Device.ID][newInterface.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Interface ", newInterface.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedInterface, err := ni.NetboxApi.PatchInterface(diffMap, ni.InterfacesIndexByDeviceAndName[newInterface.Device.ID][newInterface.Name].ID)
			if err != nil {
				return nil, err
			}
			ni.InterfacesIndexByDeviceAndName[newInterface.Device.ID][newInterface.Name] = patchedInterface
		} else {
			ni.Logger.Debug("Interface ", newInterface.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Interface ", newInterface.Name, " does not exist in NetBox. Creating it...")
		newInterface, err := ni.NetboxApi.CreateInterface(newInterface)
		if err != nil {
			return nil, err
		}
		ni.InterfacesIndexByDeviceAndName[newInterface.Device.ID][newInterface.Name] = newInterface
	}
	return ni.InterfacesIndexByDeviceAndName[newInterface.Device.ID][newInterface.Name], nil
}
