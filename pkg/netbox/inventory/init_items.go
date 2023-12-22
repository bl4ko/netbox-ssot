package inventory

import (
	"github.com/bl4ko/netbox-ssot/pkg/netbox/objects"
)

// Collect all tags from NetBox API and store them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitTags() error {
	nbTags, err := netboxInventory.NetboxApi.GetAllTags()
	if err != nil {
		return err
	}
	netboxInventory.Tags = nbTags
	netboxInventory.Logger.Debug("Successfully collected tags from NetBox: ", netboxInventory.Tags)

	// Custom tag for all netbox objects
	ssotTag, err := netboxInventory.NetboxApi.GetTagByName("netbox-ssot")
	if err != nil {
		return err
	}
	if ssotTag == nil {
		netboxInventory.Logger.Info("Tag netbox-ssot not found in NetBox. Creating it now...")
		newTag := objects.Tag{Name: "netbox-ssot", Slug: "netbox-ssot", Description: "Tag used by netbox-ssot to mark devices that are managed by it", Color: "00add8"}
		ssotTag, err = netboxInventory.NetboxApi.CreateTag(&newTag)
		if err != nil {
			return err
		}
	}
	netboxInventory.SsotTag = ssotTag
	return nil
}

// Collects all tenants from NetBox API and store them in the NetBoxInventory
func (NetBoxInventory *NetBoxInventory) InitTenants() error {
	nbTenants, err := NetBoxInventory.NetboxApi.GetAllTenants()
	if err != nil {
		return err
	}
	// We also create an index of tenants by name for easier access
	NetBoxInventory.TenantsIndexByName = make(map[string]*objects.Tenant)
	for _, tenant := range nbTenants {
		NetBoxInventory.TenantsIndexByName[tenant.Name] = tenant
	}
	NetBoxInventory.Logger.Debug("Successfully collected tenants from NetBox: ", NetBoxInventory.TenantsIndexByName)
	return nil
}

// Collects all sites from NetBox API and store them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitSites() error {
	nbSites, err := netboxInventory.NetboxApi.GetAllSites()
	if err != nil {
		return err
	}
	// We also create an index of sites by name for easier access
	netboxInventory.SitesIndexByName = make(map[string]*objects.Site)
	for _, site := range nbSites {
		netboxInventory.SitesIndexByName[site.Name] = site
	}
	netboxInventory.Logger.Debug("Successfully collected sites from NetBox: ", netboxInventory.SitesIndexByName)
	return nil
}

// Collects all manufacturesrs from NetBox API and store them in NetBoxInventory
func (netboxInventory *NetBoxInventory) InitManufacturers() error {
	nbManufacturers, err := netboxInventory.NetboxApi.GetAllManufacturers()
	if err != nil {
		return err
	}
	// We also create an index of manufacturers by name for easier access
	netboxInventory.ManufacturersIndexByName = make(map[string]*objects.Manufacturer)
	for _, manufacturer := range nbManufacturers {
		netboxInventory.ManufacturersIndexByName[manufacturer.Name] = manufacturer
	}
	netboxInventory.Logger.Debug("Successfully collected manufacturers from NetBox: ", netboxInventory.ManufacturersIndexByName)
	return nil
}

// Collects all platforms from NetBox API and store them in the NetBoxInventory
func (NetBoxInventory *NetBoxInventory) InitPlatforms() error {
	nbPlatforms, err := NetBoxInventory.NetboxApi.GetAllPlatforms()
	if err != nil {
		return err
	}
	// We also create an index of platforms by name for easier access
	NetBoxInventory.PlatformsIndexByName = make(map[string]*objects.Platform)
	for _, platform := range nbPlatforms {
		NetBoxInventory.PlatformsIndexByName[platform.Name] = platform
	}
	NetBoxInventory.Logger.Debug("Successfully collected platforms from NetBox: ", NetBoxInventory.PlatformsIndexByName)
	return nil
}

// Collect all devices from NetBox API and store them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitDevices() error {
	nbDevices, err := netboxInventory.NetboxApi.GetAllDevices()
	if err != nil {
		return err
	}
	// We also create an index of devices by name for easier access
	netboxInventory.DevicesIndexByUuid = make(map[string]*objects.Device)
	for _, device := range nbDevices {
		netboxInventory.DevicesIndexByUuid[device.AssetTag] = device
	}
	netboxInventory.Logger.Debug("Successfully collected devices from NetBox: ", netboxInventory.DevicesIndexByUuid)
	return nil
}

// Collects all deviceRoles from NetBox API and store them in the
// NetBoxInventory
func (netboxInventory *NetBoxInventory) InitDeviceRoles() error {
	nbDeviceRoles, err := netboxInventory.NetboxApi.GetAllDeviceRoles()
	if err != nil {
		return err
	}
	// We also create an index of device roles by name for easier access
	netboxInventory.DeviceRolesIndexByName = make(map[string]*objects.DeviceRole)
	for _, deviceRole := range nbDeviceRoles {
		netboxInventory.DeviceRolesIndexByName[deviceRole.Name] = deviceRole
	}

	netboxInventory.Logger.Debug("Successfully collected device roles from NetBox: ", netboxInventory.DeviceRolesIndexByName)
	return nil
}

// Ensures that attribute ServerDeviceRole is proper initialized
func (netboxInventory *NetBoxInventory) InitServerDeviceRole() error {
	err := netboxInventory.AddDeviceRole(&objects.DeviceRole{Name: "Server", Slug: "server", Color: "00add8", VMRole: true})
	if err != nil {
		return err
	}
	return nil
}

func (netboxInventory *NetBoxInventory) InitCustomFields() error {
	customFields, err := netboxInventory.NetboxApi.GetAllCustomFields()
	if err != nil {
		return err
	}
	netboxInventory.CustomFieldsIndexByName = make(map[string]*objects.CustomField)
	for _, customField := range customFields {
		netboxInventory.CustomFieldsIndexByName[customField.Name] = customField
	}
	netboxInventory.Logger.Debug("Successfully collected custom fields from NetBox: ", netboxInventory.CustomFieldsIndexByName)
	return nil
}

// This function initialises all custom fields required for servers and other objects
// Currently these are two:
// - host_cpu_cores
// - host_memory
// - sourceId - this is used to store the ID of the source object in NetBox (interfaces)
func (netboxInventory *NetBoxInventory) InitSsotCustomFields() error {
	err := netboxInventory.AddCustomField(&objects.CustomField{
		Name:          "host_cpu_cores",
		Label:         "Host CPU cores",
		Type:          objects.CustomFieldTypeText,
		FilterLogic:   objects.FilterLogicLoose,
		UIVisibility:  objects.UIVisibilityReadWrite,
		DisplayWeight: 100,
		Description:   "Number of CPU cores on the host",
		SearchWeight:  1000,
		ContentTypes:  []string{"dcim.device"},
	})
	if err != nil {
		return err
	}
	err = netboxInventory.AddCustomField(&objects.CustomField{
		Name:          "host_memory",
		Label:         "Host memory",
		Type:          objects.CustomFieldTypeText,
		FilterLogic:   objects.FilterLogicLoose,
		UIVisibility:  objects.UIVisibilityReadWrite,
		DisplayWeight: 100,
		Description:   "Amount of memory on the host",
		SearchWeight:  1000,
		ContentTypes:  []string{"dcim.device"},
	})
	if err != nil {
		return err
	}
	err = netboxInventory.AddCustomField(&objects.CustomField{
		Name:          "source_id",
		Label:         "Source ID",
		Type:          objects.CustomFieldTypeText,
		FilterLogic:   objects.FilterLogicLoose,
		UIVisibility:  objects.UIVisibilityReadWrite,
		DisplayWeight: 100,
		Description:   "ID of the object on the source API",
		SearchWeight:  1000,
		ContentTypes:  []string{"dcim.interface"},
	})
	if err != nil {
		return err
	}

	return nil
}

// Collects all nbClusters from NetBox API and stores them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitClusterGroups() error {
	nbClusters, err := netboxInventory.NetboxApi.GetAllClusterGroups()
	if err != nil {
		return err
	}
	// We also create an index of cluster groups by name for easier access
	netboxInventory.ClusterGroupsIndexByName = make(map[string]*objects.ClusterGroup)
	for _, clusterGroup := range nbClusters {
		netboxInventory.ClusterGroupsIndexByName[clusterGroup.Name] = clusterGroup
	}
	netboxInventory.Logger.Debug("Successfully collected cluster groups from NetBox: ", netboxInventory.ClusterGroupsIndexByName)
	return nil
}

// Collects all ClusterTypes from NetBox API and stores them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitClusterTypes() error {
	nbClusterTypes, err := netboxInventory.NetboxApi.GetAllClusterTypes()
	if err != nil {
		return err
	}
	netboxInventory.ClusterTypesIndexByName = make(map[string]*objects.ClusterType)
	for _, clusterType := range nbClusterTypes {
		netboxInventory.ClusterTypesIndexByName[clusterType.Name] = clusterType
	}
	netboxInventory.Logger.Debug("Successfully collected cluster types from NetBox: ", netboxInventory.ClusterTypesIndexByName)
	return nil
}

// Collects all clusters from NetBox API and stores them to local inventory
func (netboxInventory *NetBoxInventory) InitClusters() error {
	nbClusters, err := netboxInventory.NetboxApi.GetAllClusters()
	if err != nil {
		return err
	}
	netboxInventory.ClustersIndexByName = make(map[string]*objects.Cluster)
	for _, cluster := range nbClusters {
		netboxInventory.ClustersIndexByName[cluster.Name] = cluster
	}
	netboxInventory.Logger.Debug("Successfully collected clusters from NetBox: ", netboxInventory.ClustersIndexByName)
	return nil
}

func (ni *NetBoxInventory) InitDeviceTypes() error {
	nbDeviceTypes, err := ni.NetboxApi.GetAllDeviceTypes()
	if err != nil {
		return err
	}
	ni.DeviceTypesIndexByModel = make(map[string]*objects.DeviceType)
	for _, deviceType := range nbDeviceTypes {
		ni.DeviceTypesIndexByModel[deviceType.Model] = deviceType
	}
	ni.Logger.Debug("Successfully collected device types from NetBox: ", ni.DeviceTypesIndexByModel)
	return nil
}

func (ni *NetBoxInventory) InitInterfaces() error {
	nbInterfaces, err := ni.NetboxApi.GetAllInterfaces()
	if err != nil {
		return err
	}
	ni.InterfacesIndexByDeviceIdAndName = make(map[int]map[string]*objects.Interface)
	for _, intf := range nbInterfaces {
		if ni.InterfacesIndexByDeviceIdAndName[intf.Device.Id] == nil {
			ni.InterfacesIndexByDeviceIdAndName[intf.Device.Id] = make(map[string]*objects.Interface)
		}
		ni.InterfacesIndexByDeviceIdAndName[intf.Device.Id][intf.Name] = intf
	}
	ni.Logger.Debug("Successfully collected interfaces from NetBox: ", ni.InterfacesIndexByDeviceIdAndName)
	return nil
}

func (ni *NetBoxInventory) InitVlans() error {
	nbVlans, err := ni.NetboxApi.GetAllVlans()
	if err != nil {
		return err
	}
	ni.VlansIndexByName = make(map[string]*objects.Vlan)
	for _, vlan := range nbVlans {
		ni.VlansIndexByName[vlan.Name] = vlan
	}
	ni.Logger.Debug("Successfully collected vlans from NetBox: ", ni.VlansIndexByName)
	return nil
}

func (ni *NetBoxInventory) InitVMs() error {
	nbVMs, err := ni.NetboxApi.GetAllVMs()
	if err != nil {
		return err
	}
	ni.VMsIndexByName = make(map[string]*objects.VM)
	for _, vm := range nbVMs {
		ni.VMsIndexByName[vm.Name] = vm
	}
	ni.Logger.Debug("Successfully collected VMs from NetBox: ", ni.VMsIndexByName)
	return nil
}

func (ni *NetBoxInventory) InitVMInterfaces() error {
	nbVMInterfaces, err := ni.NetboxApi.GetAllVMInterfaces()
	if err != nil {
		return err
	}
	ni.VMInterfacesIndexByVmIdAndName = make(map[int]map[string]*objects.VMInterface)
	for _, vmIntf := range nbVMInterfaces {
		if ni.VMInterfacesIndexByVmIdAndName[vmIntf.VM.Id] == nil {
			ni.VMInterfacesIndexByVmIdAndName[vmIntf.VM.Id] = make(map[string]*objects.VMInterface)
		}
		ni.VMInterfacesIndexByVmIdAndName[vmIntf.VM.Id][vmIntf.Name] = vmIntf
	}
	ni.Logger.Debug("Successfully collected VM interfaces from NetBox: ", ni.VMInterfacesIndexByVmIdAndName)
	return nil
}
