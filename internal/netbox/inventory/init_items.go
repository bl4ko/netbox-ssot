package inventory

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Collect all tags from Netbox API and store them in the NetBoxInventory
func (nbi *NetBoxInventory) InitTags() error {
	nbTags, err := nbi.NetboxApi.GetAllTags()
	if err != nil {
		return err
	}
	nbi.Tags = nbTags
	nbi.Logger.Debug("Successfully collected tags from Netbox: ", nbi.Tags)

	// Custom tag for all netbox objects
	ssotTag, err := nbi.NetboxApi.GetTagByName("netbox-ssot")
	if err != nil {
		return err
	}
	if ssotTag == nil {
		nbi.Logger.Info("Tag netbox-ssot not found in Netbox. Creating it now...")
		newTag := objects.Tag{Name: "netbox-ssot", Slug: "netbox-ssot", Description: "Tag used by netbox-ssot to mark devices that are managed by it", Color: "00add8"}
		ssotTag, err = nbi.NetboxApi.CreateTag(&newTag)
		if err != nil {
			return err
		}
	}
	nbi.SsotTag = ssotTag
	return nil
}

// Collects all tenants from Netbox API and store them in the NetBoxInventory
func (nbi *NetBoxInventory) InitTenants() error {
	nbTenants, err := nbi.NetboxApi.GetAllTenants()
	if err != nil {
		return err
	}
	// We also create an index of tenants by name for easier access
	nbi.TenantsIndexByName = make(map[string]*objects.Tenant)
	for _, tenant := range nbTenants {
		nbi.TenantsIndexByName[tenant.Name] = tenant
	}
	nbi.Logger.Debug("Successfully collected tenants from Netbox: ", nbi.TenantsIndexByName)
	return nil
}

// Collects all sites from Netbox API and store them in the NetBoxInventory
func (nbi *NetBoxInventory) InitSites() error {
	nbSites, err := nbi.NetboxApi.GetAllSites()
	if err != nil {
		return err
	}
	// We also create an index of sites by name for easier access
	nbi.SitesIndexByName = make(map[string]*objects.Site)
	for _, site := range nbSites {
		nbi.SitesIndexByName[site.Name] = site
	}
	nbi.Logger.Debug("Successfully collected sites from Netbox: ", nbi.SitesIndexByName)
	return nil
}

// Collects all manufacturers from Netbox API and store them in NetBoxInventory
func (nbi *NetBoxInventory) InitManufacturers() error {
	nbManufacturers, err := nbi.NetboxApi.GetAllManufacturers()
	if err != nil {
		return err
	}
	// Initialize internal index of manufacturers by name
	nbi.ManufacturersIndexByName = make(map[string]*objects.Manufacturer)
	// OrphanManager takes care of all manufacturers created by netbox-ssot
	nbi.OrphanManager["/api/dcim/manufacturers/"] = make(map[int]bool, 0)

	for _, manufacturer := range nbManufacturers {
		nbi.ManufacturersIndexByName[manufacturer.Name] = manufacturer
		nbi.OrphanManager["/api/dcim/manufacturers/"][manufacturer.Id] = true
	}

	nbi.Logger.Debug("Successfully collected manufacturers from Netbox: ", nbi.ManufacturersIndexByName)
	return nil
}

// Collects all platforms from Netbox API and store them in the NetBoxInventory
func (nbi *NetBoxInventory) InitPlatforms() error {
	nbPlatforms, err := nbi.NetboxApi.GetAllPlatforms()
	if err != nil {
		return err
	}
	// Initialize internal index of platforms by name
	nbi.PlatformsIndexByName = make(map[string]*objects.Platform)
	// OrphanManager takes care of all platforms created by netbox-ssot
	nbi.OrphanManager["/api/dcim/platforms/"] = make(map[int]bool, 0)

	for _, platform := range nbPlatforms {
		nbi.PlatformsIndexByName[platform.Name] = platform
		nbi.OrphanManager["/api/dcim/platforms/"][platform.Id] = true
	}

	nbi.Logger.Debug("Successfully collected platforms from Netbox: ", nbi.PlatformsIndexByName)
	return nil
}

// Collect all devices from Netbox API and store them in the NetBoxInventory.
func (nbi *NetBoxInventory) InitDevices() error {
	nbDevices, err := nbi.NetboxApi.GetAllDevices()
	if err != nil {
		return err
	}
	// Initialize internal index of devices by Name and SiteId
	nbi.DevicesIndexByNameAndSiteId = make(map[string]map[int]*objects.Device)
	// OrphanManager takes care of all devices created by netbox-ssot
	nbi.OrphanManager["/api/dcim/devices/"] = make(map[int]bool, 0)

	for _, device := range nbDevices {
		if nbi.DevicesIndexByNameAndSiteId[device.Name] == nil {
			nbi.DevicesIndexByNameAndSiteId[device.Name] = make(map[int]*objects.Device)
		}
		nbi.DevicesIndexByNameAndSiteId[device.Name][device.Site.Id] = device
		nbi.OrphanManager["/api/dcim/devices/"][device.Id] = true
	}

	nbi.Logger.Debug("Successfully collected devices from Netbox: ", nbi.DevicesIndexByNameAndSiteId)
	return nil
}

// Collects all deviceRoles from Netbox API and store them in the
// NetBoxInventory
func (nbi *NetBoxInventory) InitDeviceRoles() error {
	nbDeviceRoles, err := nbi.NetboxApi.GetAllDeviceRoles()
	if err != nil {
		return err
	}
	// We also create an index of device roles by name for easier access
	nbi.DeviceRolesIndexByName = make(map[string]*objects.DeviceRole)
	// OrphanManager takes care of all device roles created by netbox-ssot
	nbi.OrphanManager["/api/dcim/device-roles/"] = make(map[int]bool, 0)

	for _, deviceRole := range nbDeviceRoles {
		nbi.DeviceRolesIndexByName[deviceRole.Name] = deviceRole
		nbi.OrphanManager["/api/dcim/device-roles/"][deviceRole.Id] = true
	}

	nbi.Logger.Debug("Successfully collected device roles from Netbox: ", nbi.DeviceRolesIndexByName)
	return nil
}

// Ensures that attribute ServerDeviceRole is proper initialized
func (nbi *NetBoxInventory) InitServerDeviceRole() error {
	err := nbi.AddDeviceRole(&objects.DeviceRole{Name: "Server", Slug: "server", Color: "00add8", VMRole: true})
	if err != nil {
		return err
	}
	return nil
}

func (nbi *NetBoxInventory) InitCustomFields() error {
	customFields, err := nbi.NetboxApi.GetAllCustomFields()
	if err != nil {
		return err
	}
	// Initialize internal index of custom fields by name
	nbi.CustomFieldsIndexByName = make(map[string]*objects.CustomField)
	for _, customField := range customFields {
		nbi.CustomFieldsIndexByName[customField.Name] = customField
	}
	nbi.Logger.Debug("Successfully collected custom fields from Netbox: ", nbi.CustomFieldsIndexByName)
	return nil
}

// This function Initializes all custom fields required for servers and other objects
// Currently these are two:
// - host_cpu_cores
// - host_memory
// - sourceId - this is used to store the ID of the source object in Netbox (interfaces)
func (netboxInventory *NetBoxInventory) InitSsotCustomFields() error {
	err := netboxInventory.AddCustomField(&objects.CustomField{
		Name:                  "host_cpu_cores",
		Label:                 "Host CPU cores",
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         100,
		Description:           "Number of CPU cores on the host",
		SearchWeight:          1000,
		ContentTypes:          []string{"dcim.device"},
	})
	if err != nil {
		return err
	}
	err = netboxInventory.AddCustomField(&objects.CustomField{
		Name:                  "host_memory",
		Label:                 "Host memory",
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         100,
		Description:           "Amount of memory on the host",
		SearchWeight:          1000,
		ContentTypes:          []string{"dcim.device"},
	})
	if err != nil {
		return err
	}
	err = netboxInventory.AddCustomField(&objects.CustomField{
		Name:                  "source_id",
		Label:                 "Source ID",
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         100,
		Description:           "ID of the object on the source API",
		SearchWeight:          1000,
		ContentTypes:          []string{"dcim.interface"},
	})
	if err != nil {
		return err
	}

	return nil
}

// Collects all nbClusters from Netbox API and stores them in the NetBoxInventory
func (nbi *NetBoxInventory) InitClusterGroups() error {
	nbClusters, err := nbi.NetboxApi.GetAllClusterGroups()
	if err != nil {
		return err
	}
	// Initialize internal index of cluster groups by name
	nbi.ClusterGroupsIndexByName = make(map[string]*objects.ClusterGroup)
	// OrphanManager takes care of all cluster groups created by netbox-ssot
	nbi.OrphanManager["/api/virtualization/cluster-groups/"] = make(map[int]bool, 0)

	for _, clusterGroup := range nbClusters {
		nbi.ClusterGroupsIndexByName[clusterGroup.Name] = clusterGroup
		nbi.OrphanManager["/api/virtualization/cluster-groups/"][clusterGroup.Id] = true
	}
	nbi.Logger.Debug("Successfully collected cluster groups from Netbox: ", nbi.ClusterGroupsIndexByName)
	return nil
}

// Collects all ClusterTypes from Netbox API and stores them in the NetBoxInventory
func (nbi *NetBoxInventory) InitClusterTypes() error {
	nbClusterTypes, err := nbi.NetboxApi.GetAllClusterTypes()
	if err != nil {
		return err
	}

	// Initialize internal index of cluster types by name
	nbi.ClusterTypesIndexByName = make(map[string]*objects.ClusterType)
	// OrphanManager takes care of all cluster types created by netbox-ssot
	nbi.OrphanManager["/api/virtualization/cluster-types/"] = make(map[int]bool, 0)

	for _, clusterType := range nbClusterTypes {
		nbi.ClusterTypesIndexByName[clusterType.Name] = clusterType
		nbi.OrphanManager["/api/virtualization/cluster-types/"][clusterType.Id] = true
	}

	nbi.Logger.Debug("Successfully collected cluster types from Netbox: ", nbi.ClusterTypesIndexByName)
	return nil
}

// Collects all clusters from Netbox API and stores them to local inventory
func (nbi *NetBoxInventory) InitClusters() error {
	nbClusters, err := nbi.NetboxApi.GetAllClusters()
	if err != nil {
		return err
	}

	// Initialize internal index of clusters by name
	nbi.ClustersIndexByName = make(map[string]*objects.Cluster)
	// OrphanManager takes care of all clusters created by netbox-ssot
	nbi.OrphanManager["/api/virtualization/clusters/"] = make(map[int]bool, 0)

	for _, cluster := range nbClusters {
		nbi.ClustersIndexByName[cluster.Name] = cluster
		nbi.OrphanManager["/api/virtualization/clusters/"][cluster.Id] = true
	}

	nbi.Logger.Debug("Successfully collected clusters from Netbox: ", nbi.ClustersIndexByName)
	return nil
}

func (ni *NetBoxInventory) InitDeviceTypes() error {
	nbDeviceTypes, err := ni.NetboxApi.GetAllDeviceTypes()
	if err != nil {
		return err
	}

	// Initialize internal index of device types by model
	ni.DeviceTypesIndexByModel = make(map[string]*objects.DeviceType)
	// OrphanManager takes care of all device types created by netbox-ssot
	ni.OrphanManager["/api/dcim/device-types/"] = make(map[int]bool, 0)

	for _, deviceType := range nbDeviceTypes {
		ni.DeviceTypesIndexByModel[deviceType.Model] = deviceType
		ni.OrphanManager["/api/dcim/device-types/"][deviceType.Id] = true
	}

	ni.Logger.Debug("Successfully collected device types from Netbox: ", ni.DeviceTypesIndexByModel)
	return nil
}

// Collects all interfaces from Netbox API and stores them to local inventory
func (ni *NetBoxInventory) InitInterfaces() error {
	nbInterfaces, err := ni.NetboxApi.GetAllInterfaces()
	if err != nil {
		return err
	}

	// Initialize internal index of interfaces by device id and name
	ni.InterfacesIndexByDeviceIdAndName = make(map[int]map[string]*objects.Interface)
	// OrphanManager takes care of all interfaces created by netbox-ssot
	ni.OrphanManager["/api/dcim/interfaces/"] = make(map[int]bool, 0)

	for _, intf := range nbInterfaces {
		if ni.InterfacesIndexByDeviceIdAndName[intf.Device.Id] == nil {
			ni.InterfacesIndexByDeviceIdAndName[intf.Device.Id] = make(map[string]*objects.Interface)
		}
		ni.InterfacesIndexByDeviceIdAndName[intf.Device.Id][intf.Name] = intf
		ni.OrphanManager["/api/dcim/interfaces/"][intf.Id] = true
	}

	ni.Logger.Debug("Successfully collected interfaces from Netbox: ", ni.InterfacesIndexByDeviceIdAndName)
	return nil
}

// Inits default VlanGroup, which is required to group all Vlans that are not part of other
// vlangroups into it. Each vlan is indexed by their (vlanGroup, vid).
func (ni *NetBoxInventory) InitDefaultVlanGroup() error {
	_, err := ni.AddVlanGroup(&objects.VlanGroup{
		NetboxObject: objects.NetboxObject{
			Tags:        []*objects.Tag{ni.SsotTag},
			Description: "Default netbox-ssot VlanGroup for all vlans that are not part of any other vlanGroup. This group is required for netbox-ssot vlan index to work.",
		},
		Name:   objects.DefaultVlanGroupName,
		Slug:   utils.Slugify(objects.DefaultVlanGroupName),
		MinVid: 1,
		MaxVid: 4094,
	})
	if err != nil {
		return fmt.Errorf("init default vlan group: %s", err)
	}
	return nil
}

// Collects all vlans from Netbox API and stores them to local inventory
func (ni *NetBoxInventory) InitVlanGroups() error {
	nbVlanGroups, err := ni.NetboxApi.GetAllVlanGroups()
	if err != nil {
		return err
	}

	// Initialize internal index of vlans by name
	ni.VlanGroupsIndexByName = make(map[string]*objects.VlanGroup)
	// Add VlanGroups to orphan manager
	ni.OrphanManager["/api/ipam/vlan-groups/"] = make(map[int]bool, 0)

	for _, vlanGroup := range nbVlanGroups {
		ni.VlanGroupsIndexByName[vlanGroup.Name] = vlanGroup
		ni.OrphanManager["/api/ipam/vlan-groups/"][vlanGroup.Id] = true
	}

	ni.Logger.Debug("Successfully collected vlans from Netbox: ", ni.VlanGroupsIndexByName)
	return nil
}

// Collects all vlans from Netbox API and stores them to local inventory
func (ni *NetBoxInventory) InitVlans() error {
	nbVlans, err := ni.NetboxApi.GetAllVlans()
	if err != nil {
		return err
	}

	// Initialize internal index of vlans by VlanGroupId and Vid
	ni.VlansIndexByVlanGroupIdAndVid = make(map[int]map[int]*objects.Vlan)
	// Add vlans to orphan manager
	ni.OrphanManager["/api/ipam/vlans/"] = make(map[int]bool, 0)

	for _, vlan := range nbVlans {
		if vlan.Group == nil {
			// Update all existing vlans with default vlanGroup. This only happens
			// when there are predefined vlans in netbox.
			vlan.Group = ni.VlanGroupsIndexByName[objects.DefaultVlanGroupName] // This should not fail, because InitDefaultVlanGroup executes before InitVlans
			vlan, err = ni.AddVlan(vlan)
			if err != nil {
				return err
			}
		}
		if ni.VlansIndexByVlanGroupIdAndVid[vlan.Group.Id] == nil {
			ni.VlansIndexByVlanGroupIdAndVid[vlan.Group.Id] = make(map[int]*objects.Vlan)
		}
		ni.VlansIndexByVlanGroupIdAndVid[vlan.Group.Id][vlan.Vid] = vlan
		ni.OrphanManager["/api/ipam/vlans/"][vlan.Id] = true
	}

	ni.Logger.Debug("Successfully collected vlans from Netbox: ", ni.VlansIndexByVlanGroupIdAndVid)
	return nil
}

// Collects all vms from Netbox API and stores them to local inventory
func (ni *NetBoxInventory) InitVMs() error {
	nbVMs, err := ni.NetboxApi.GetAllVMs()
	if err != nil {
		return err
	}

	// Initialize internal index of VMs by name
	ni.VMsIndexByName = make(map[string]*objects.VM)
	// Add VMs to orphan manager
	ni.OrphanManager["/api/virtualization/virtual-machines/"] = make(map[int]bool, 0)

	for _, vm := range nbVMs {
		ni.VMsIndexByName[vm.Name] = vm
		ni.OrphanManager["/api/virtualization/virtual-machines/"][vm.Id] = true
	}

	ni.Logger.Debug("Successfully collected VMs from Netbox: ", ni.VMsIndexByName)
	return nil
}

// Collects all VMInterfaces from Netbox API and stores them to local inventory
func (ni *NetBoxInventory) InitVMInterfaces() error {
	nbVMInterfaces, err := ni.NetboxApi.GetAllVMInterfaces()
	if err != nil {
		return err
	}

	// Initialize internal index of VM interfaces by VM id and name
	ni.VMInterfacesIndexByVMIdAndName = make(map[int]map[string]*objects.VMInterface)
	// Add VMInterfaces to orphan manager
	ni.OrphanManager["/api/virtualization/interfaces/"] = make(map[int]bool, 0)

	for _, vmIntf := range nbVMInterfaces {
		if ni.VMInterfacesIndexByVMIdAndName[vmIntf.VM.Id] == nil {
			ni.VMInterfacesIndexByVMIdAndName[vmIntf.VM.Id] = make(map[string]*objects.VMInterface)
		}
		ni.VMInterfacesIndexByVMIdAndName[vmIntf.VM.Id][vmIntf.Name] = vmIntf
		ni.OrphanManager["/api/virtualization/interfaces/"][vmIntf.Id] = true
	}

	ni.Logger.Debug("Successfully collected VM interfaces from Netbox: ", ni.VMInterfacesIndexByVMIdAndName)
	return nil
}

// Collects all IP addresses from Netbox API and stores them to local inventory
func (ni *NetBoxInventory) InitIPAddresses() error {
	ipAddresses, err := ni.NetboxApi.GetAllIPAddresses()
	if err != nil {
		return err
	}

	// Initializes internal index of IP addresses by address
	ni.IPAdressesIndexByAddress = make(map[string]*objects.IPAddress)
	// Add IP addresses to orphan manager
	ni.OrphanManager["/api/ipam/ip-addresses/"] = make(map[int]bool, 0)

	for _, ipAddr := range ipAddresses {
		ni.IPAdressesIndexByAddress[ipAddr.Address] = ipAddr
		ni.OrphanManager["/api/ipam/ip-addresses/"][ipAddr.Id] = true
	}

	ni.Logger.Debug("Successfully collected IP addresses from Netbox: ", ni.IPAdressesIndexByAddress)
	return nil
}
