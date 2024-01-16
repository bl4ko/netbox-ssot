package inventory

import (
	"slices"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

func (ni *NetBoxInventory) AddTag(newTag *objects.Tag) (*objects.Tag, error) {
	existingTagIndex := slices.IndexFunc(ni.Tags, func(t *objects.Tag) bool {
		return t.Name == newTag.Name
	})
	if existingTagIndex == -1 {
		ni.Logger.Debug("Tag ", newTag.Name, " does not exist in Netbox. Creating it...")
		createdTag, err := ni.NetboxApi.CreateTag(newTag)
		if err != nil {
			return nil, err
		}
		ni.Tags = append(ni.Tags, createdTag)
		return createdTag, nil
	} else {
		ni.Logger.Debug("Tag ", newTag.Name, " already exists in Netbox...")
		existingTag := ni.Tags[existingTagIndex]
		diffMap, err := utils.JsonDiffMapExceptId(newTag, existingTag)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			patchedTag, err := ni.NetboxApi.PatchTag(diffMap, existingTag.Id)
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

func (ni *NetBoxInventory) AddCustomField(newCf *objects.CustomField) error {
	if _, ok := ni.CustomFieldsIndexByName[newCf.Name]; ok {
		existingCf := ni.CustomFieldsIndexByName[newCf.Name]
		diffMap, err := utils.JsonDiffMapExceptId(newCf, existingCf)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Custom field ", newCf.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedCf, err := ni.NetboxApi.PatchCustomField(diffMap, existingCf.Id)
			if err != nil {
				return err
			}
			ni.CustomFieldsIndexByName[newCf.Name] = patchedCf
		} else {
			ni.Logger.Debug("Custom field ", newCf.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Custom field ", newCf.Name, " does not exist in Netbox. Creating it...")
		newCf, err := ni.NetboxApi.CreateCustomField(newCf)
		if err != nil {
			return err
		}
		ni.CustomFieldsIndexByName[newCf.Name] = newCf
	}
	return nil
}

func (ni *NetBoxInventory) AddClusterGroup(newCg *objects.ClusterGroup) (*objects.ClusterGroup, error) {
	newCg.Tags = append(newCg.Tags, ni.SsotTag)
	if _, ok := ni.ClusterGroupsIndexByName[newCg.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/virtualization/cluster-groups/"], ni.ClusterGroupsIndexByName[newCg.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newCg, ni.ClusterGroupsIndexByName[newCg.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster group ", newCg.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedCg, err := ni.NetboxApi.PatchClusterGroup(diffMap, ni.ClusterGroupsIndexByName[newCg.Name].Id)
			if err != nil {
				return nil, err
			}
			ni.ClusterGroupsIndexByName[newCg.Name] = patchedCg
		} else {
			ni.Logger.Debug("Cluster group ", newCg.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Cluster group ", newCg.Name, " does not exist in Netbox. Creating it...")
		newCg, err := ni.NetboxApi.CreateClusterGroup(newCg)
		if err != nil {
			return nil, err
		}
		ni.ClusterGroupsIndexByName[newCg.Name] = newCg
	}
	// Delete id from orphan manager
	return ni.ClusterGroupsIndexByName[newCg.Name], nil
}

func (ni *NetBoxInventory) AddClusterType(newClusterType *objects.ClusterType) (*objects.ClusterType, error) {
	newClusterType.Tags = append(newClusterType.Tags, ni.SsotTag)
	if _, ok := ni.ClusterTypesIndexByName[newClusterType.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/virtualization/cluster-types/"], ni.ClusterTypesIndexByName[newClusterType.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newClusterType, ni.ClusterTypesIndexByName[newClusterType.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedClusterType, err := ni.NetboxApi.PatchClusterType(diffMap, ni.ClusterTypesIndexByName[newClusterType.Name].Id)
			if err != nil {
				return nil, err
			}
			ni.ClusterTypesIndexByName[newClusterType.Name] = patchedClusterType
			return patchedClusterType, nil
		} else {
			ni.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in Netbox and is up to date...")
			existingClusterType := ni.ClusterTypesIndexByName[newClusterType.Name]
			return existingClusterType, nil
		}
	} else {
		ni.Logger.Debug("Cluster type ", newClusterType.Name, " does not exist in Netbox. Creating it...")
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
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/virtualization/clusters/"], ni.ClustersIndexByName[newCluster.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newCluster, ni.ClustersIndexByName[newCluster.Name])
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster ", newCluster.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedCluster, err := ni.NetboxApi.PatchCluster(diffMap, ni.ClustersIndexByName[newCluster.Name].Id)
			if err != nil {
				return err
			}
			ni.ClustersIndexByName[newCluster.Name] = patchedCluster
		} else {
			ni.Logger.Debug("Cluster ", newCluster.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Cluster ", newCluster.Name, " does not exist in Netbox. Creating it...")
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
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/dcim/device-roles/"], ni.DeviceRolesIndexByName[newDeviceRole.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newDeviceRole, ni.DeviceRolesIndexByName[newDeviceRole.Name])
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedDeviceRole, err := ni.NetboxApi.PatchDeviceRole(diffMap, ni.DeviceRolesIndexByName[newDeviceRole.Name].Id)
			if err != nil {
				return err
			}
			ni.DeviceRolesIndexByName[newDeviceRole.Name] = patchedDeviceRole
		} else {
			ni.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Device role ", newDeviceRole.Name, " does not exist in Netbox. Creating it...")
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
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/dcim/manufacturers/"], ni.ManufacturersIndexByName[newManufacturer.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newManufacturer, ni.ManufacturersIndexByName[newManufacturer.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Manufacturer ", newManufacturer.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedManufacturer, err := ni.NetboxApi.PatchManufacturer(diffMap, ni.ManufacturersIndexByName[newManufacturer.Name].Id)
			if err != nil {
				return nil, err
			}
			ni.ManufacturersIndexByName[newManufacturer.Name] = patchedManufacturer
		} else {
			ni.Logger.Debug("Manufacturer ", newManufacturer.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Manufacturer ", newManufacturer.Name, " does not exist in Netbox. Creating it...")
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
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/dcim/device-types/"], ni.DeviceTypesIndexByModel[newDeviceType.Model].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newDeviceType, ni.DeviceTypesIndexByModel[newDeviceType.Model])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Device type ", newDeviceType.Model, " already exists in Netbox but is out of date. Patching it...")
			patchedDeviceType, err := ni.NetboxApi.PatchDeviceType(diffMap, ni.DeviceTypesIndexByModel[newDeviceType.Model].Id)
			if err != nil {
				return nil, err
			}
			ni.DeviceTypesIndexByModel[newDeviceType.Model] = patchedDeviceType
		} else {
			ni.Logger.Debug("Device type ", newDeviceType.Model, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Device type ", newDeviceType.Model, " does not exist in Netbox. Creating it...")
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
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/dcim/platforms/"], ni.PlatformsIndexByName[newPlatform.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newPlatform, ni.PlatformsIndexByName[newPlatform.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Platform ", newPlatform.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedPlatform, err := ni.NetboxApi.PatchPlatform(diffMap, ni.PlatformsIndexByName[newPlatform.Name].Id)
			if err != nil {
				return nil, err
			}
			ni.PlatformsIndexByName[newPlatform.Name] = patchedPlatform
		} else {
			ni.Logger.Debug("Platform ", newPlatform.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Platform ", newPlatform.Name, " does not exist in Netbox. Creating it...")
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
	if _, ok := ni.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/dcim/devices/"], ni.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newDevice, ni.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Device ", newDevice.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedDevice, err := ni.NetboxApi.PatchDevice(diffMap, ni.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id].Id)
			if err != nil {
				return nil, err
			}
			ni.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id] = patchedDevice
		} else {
			ni.Logger.Debug("Device ", newDevice.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Device ", newDevice.Name, " does not exist in Netbox. Creating it...")
		newDevice, err := ni.NetboxApi.CreateDevice(newDevice)
		if err != nil {
			return nil, err
		}
		ni.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id] = newDevice
	}
	return ni.DevicesIndexByNameAndSiteId[newDevice.Name][newDevice.Site.Id], nil
}

func (ni *NetBoxInventory) AddVlan(newVlan *objects.Vlan) (*objects.Vlan, error) {
	newVlan.Tags = append(newVlan.Tags, ni.SsotTag)
	if _, ok := ni.VlansIndexByName[newVlan.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/ipam/vlans/"], ni.VlansIndexByName[newVlan.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newVlan, ni.VlansIndexByName[newVlan.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Vlan ", newVlan.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVlan, err := ni.NetboxApi.PatchVlan(diffMap, ni.VlansIndexByName[newVlan.Name].Id)
			if err != nil {
				return nil, err
			}
			ni.VlansIndexByName[newVlan.Name] = patchedVlan
		} else {
			ni.Logger.Debug("Vlan ", newVlan.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Vlan ", newVlan.Name, " does not exist in Netbox. Creating it...")
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
	if _, ok := ni.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/dcim/interfaces/"], ni.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newInterface, ni.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Interface ", newInterface.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedInterface, err := ni.NetboxApi.PatchInterface(diffMap, ni.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name].Id)
			if err != nil {
				return nil, err
			}
			ni.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name] = patchedInterface
		} else {
			ni.Logger.Debug("Interface ", newInterface.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Interface ", newInterface.Name, " does not exist in Netbox. Creating it...")
		newInterface, err := ni.NetboxApi.CreateInterface(newInterface)
		if err != nil {
			return nil, err
		}
		ni.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name] = newInterface
	}
	return ni.InterfacesIndexByDeviceIdAndName[newInterface.Device.Id][newInterface.Name], nil
}

func (ni *NetBoxInventory) AddVM(newVm *objects.VM) (*objects.VM, error) {
	newVm.Tags = append(newVm.Tags, ni.SsotTag)
	if _, ok := ni.VMsIndexByName[newVm.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/virtualization/virtual-machines/"], ni.VMsIndexByName[newVm.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newVm, ni.VMsIndexByName[newVm.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("VM ", newVm.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVm, err := ni.NetboxApi.PatchVM(diffMap, ni.VMsIndexByName[newVm.Name].Id)
			if err != nil {
				return nil, err
			}
			ni.VMsIndexByName[newVm.Name] = patchedVm
		} else {
			ni.Logger.Debug("VM ", newVm.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("VM ", newVm.Name, " does not exist in Netbox. Creating it...")
		newVm, err := ni.NetboxApi.CreateVM(newVm)
		if err != nil {
			return nil, err
		}
		ni.VMsIndexByName[newVm.Name] = newVm
		return newVm, nil
	}
	return ni.VMsIndexByName[newVm.Name], nil
}

func (ni *NetBoxInventory) AddVMInterface(newVMInterface *objects.VMInterface) (*objects.VMInterface, error) {
	newVMInterface.Tags = append(newVMInterface.Tags, ni.SsotTag)
	if _, ok := ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/virtualization/interfaces/"], ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newVMInterface, ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("VM interface ", newVMInterface.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVMInterface, err := ni.NetboxApi.PatchVMInterface(diffMap, ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name].Id)
			if err != nil {
				return nil, err
			}
			ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name] = patchedVMInterface
		} else {
			ni.Logger.Debug("VM interface ", newVMInterface.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("VM interface ", newVMInterface.Name, " does not exist in Netbox. Creating it...")
		newVMInterface, err := ni.NetboxApi.CreateVMInterface(newVMInterface)
		if err != nil {
			return nil, err
		}
		if ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id] == nil {
			ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id] = make(map[string]*objects.VMInterface)
		}
		ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name] = newVMInterface
	}
	return ni.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.Id][newVMInterface.Name], nil
}

func (ni *NetBoxInventory) AddIPAddress(newIPAddress *objects.IPAddress) (*objects.IPAddress, error) {
	newIPAddress.Tags = append(newIPAddress.Tags, ni.SsotTag)
	if _, ok := ni.IPAdressesIndexByAddress[newIPAddress.Address]; ok {
		// Delete id from orphan manager, because it still exists in the sources
		delete(ni.OrphanManager["/api/ipam/ip-addresses/"], ni.IPAdressesIndexByAddress[newIPAddress.Address].Id)
		diffMap, err := utils.JsonDiffMapExceptId(newIPAddress, ni.IPAdressesIndexByAddress[newIPAddress.Address])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("IP address ", newIPAddress.Address, " already exists in Netbox but is out of date. Patching it...")
			patchedIPAddress, err := ni.NetboxApi.PatchIPAddress(diffMap, ni.IPAdressesIndexByAddress[newIPAddress.Address].Id)
			if err != nil {
				return nil, err
			}
			ni.IPAdressesIndexByAddress[newIPAddress.Address] = patchedIPAddress
		} else {
			ni.Logger.Debug("IP address ", newIPAddress.Address, " already exists in Netbox and is up to date...")
		}
	} else {
		ni.Logger.Debug("IP address ", newIPAddress.Address, " does not exist in Netbox. Creating it...")
		newIPAddress, err := ni.NetboxApi.CreateIPAddress(newIPAddress)
		if err != nil {
			return nil, err
		}
		ni.IPAdressesIndexByAddress[newIPAddress.Address] = newIPAddress
	}
	return ni.IPAdressesIndexByAddress[newIPAddress.Address], nil
}
