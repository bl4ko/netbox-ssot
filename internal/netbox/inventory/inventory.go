package inventory

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/parser"
)

// NetBoxInventory is a singleton class to manage a inventory of NetBoxObject objects
type NetBoxInventory struct {
	// Logger is the logger used for logging messages
	Logger *logger.Logger
	// NetboxConfig is the Netbox configuration
	NetboxConfig *parser.NetboxConfig
	// NetboxApi is the Netbox API object, for communicating with the Netbox API
	NetboxApi *service.NetboxAPI
	// Tags is a list of all tags in the netbox inventory
	Tags []*objects.Tag
	// SitesIndexByName is a map of all sites in the Netbox's inventory, indexed by their name
	SitesIndexByName map[string]*objects.Site
	// ManufacturersIndexByName is a map of all manufacturers in the Netbox's inventory, indexed by their name
	ManufacturersIndexByName map[string]*objects.Manufacturer
	// PlatformsIndexByName is a map of all platforms in the Netbox's inventory, indexed by their name
	PlatformsIndexByName map[string]*objects.Platform
	// TenantsIndexByName is a map of all tenants in the Netbox's inventory, indexed by their name
	TenantsIndexByName map[string]*objects.Tenant
	// DeviceTypesIndexByModel is a map of all device types in the Netbox's inventory, indexed by their model
	DeviceTypesIndexByModel map[string]*objects.DeviceType
	// DevicesIndexByUuid is a map of all devices in the Netbox's inventory, indexed by uuid (unique identifier)
	DevicesIndexByUuid map[string]*objects.Device
	// VlansIndexByName is a map of all vlans in the Netbox's inventory, indexed by their name
	VlansIndexByName map[string]*objects.Vlan
	// ClusterGroupsIndexByName is a map of all cluster groups in the Netbox's inventory, indexed by their name
	ClusterGroupsIndexByName map[string]*objects.ClusterGroup
	// ClusterTypesIndexByName is a map of all cluster types in the Netbox's inventory, indexed by their name
	ClusterTypesIndexByName map[string]*objects.ClusterType
	// ClustersIndexByName is a map of all clusters in the Netbox's inventory, indexed by their name
	ClustersIndexByName map[string]*objects.Cluster
	// Netbox's Device Roles is a map of all device roles in the inventory, indexed by name
	DeviceRolesIndexByName map[string]*objects.DeviceRole
	// CustomFieldsIndexByName is a map of all custom fields in the inventory, indexed by name
	CustomFieldsIndexByName map[string]*objects.CustomField
	// InterfacesIndexByDeviceAnName is a map of all interfaces in the inventory, indexed by their's
	// device id and their name.
	InterfacesIndexByDeviceIdAndName map[int]map[string]*objects.Interface
	// VirtualMachinedIndexByName is a map of all virtual machines in the inventory, indexed by their name
	VMsIndexByName map[string]*objects.VM
	// VirtualMachineInterfacesIndexByVMAndName is a map of all virtual machine interfaces in the inventory, indexed by their's virtual machine id and their name
	VMInterfacesIndexByVMIdAndName map[int]map[string]*objects.VMInterface
	// IPAdressesIndexByAddress is a map of all IP addresses in the inventory, indexed by their address
	IPAdressesIndexByAddress map[string]*objects.IPAddress

	// Orphan manager is a map of objectAPIPath to a set of managed ids for that object type.
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
	OrphanManager map[string]map[int]bool

	// Tag used by netbox-ssot to mark devices that are managed by it
	SsotTag *objects.Tag
}

// Func string representation
func (nbi NetBoxInventory) String() string {
	return fmt.Sprintf("NetBoxInventory{Logger: %+v, NetboxConfig: %+v...}", nbi.Logger, nbi.NetboxConfig)
}

// NewNetboxInventory creates a new NetBoxInventory object.
// It takes a logger and a NetboxConfig as parameters, and returns a pointer to the newly created NetBoxInventory.
// The logger is used for logging messages, and the NetboxConfig is used to configure the NetBoxInventory.
func NewNetboxInventory(logger *logger.Logger, nbConfig *parser.NetboxConfig) *NetBoxInventory {
	nbi := &NetBoxInventory{Logger: logger, NetboxConfig: nbConfig, OrphanManager: make(map[string]map[int]bool)}
	return nbi
}

// Init function that initialises the NetBoxInventory object with objects from Netbox
func (netboxInventory *NetBoxInventory) Init() error {
	baseURL := fmt.Sprintf("%s://%s:%d", netboxInventory.NetboxConfig.HTTPScheme, netboxInventory.NetboxConfig.Hostname, netboxInventory.NetboxConfig.Port)

	netboxInventory.Logger.Debug("Initialising Netbox API with baseURL: ", baseURL)
	netboxInventory.NetboxApi = service.NewNetBoxAPI(netboxInventory.Logger, baseURL, netboxInventory.NetboxConfig.ApiToken, netboxInventory.NetboxConfig.ValidateCert)

	err := netboxInventory.InitTags()
	if err != nil {
		return err
	}
	err = netboxInventory.InitTenants()
	if err != nil {
		return err
	}
	err = netboxInventory.InitSites()
	if err != nil {
		return err
	}
	err = netboxInventory.InitManufacturers()
	if err != nil {
		return err
	}
	err = netboxInventory.InitPlatforms()
	if err != nil {
		return err
	}
	err = netboxInventory.InitDevices()
	if err != nil {
		return err
	}
	err = netboxInventory.InitInterfaces()
	if err != nil {
		return err
	}
	err = netboxInventory.InitIPAddresses()
	if err != nil {
		return err
	}
	err = netboxInventory.InitVlans()
	if err != nil {
		return err
	}
	err = netboxInventory.InitDeviceRoles()
	if err != nil {
		return err
	}
	// init server device role which is required for separation of device object into servers
	err = netboxInventory.InitServerDeviceRole()
	if err != nil {
		return err
	}
	err = netboxInventory.InitDeviceTypes()
	if err != nil {
		return err
	}
	// init custom fields. Custom fields can be used for devices to add physical cores and memory to each device representing server.
	err = netboxInventory.InitCustomFields()
	if err != nil {
		return err
	}
	err = netboxInventory.InitSsotCustomFields()
	if err != nil {
		return err
	}
	err = netboxInventory.InitClusterGroups()
	if err != nil {
		return err
	}
	err = netboxInventory.InitClusterTypes()
	if err != nil {
		return err
	}
	err = netboxInventory.InitClusters()
	if err != nil {
		return err
	}
	err = netboxInventory.InitVMs()
	if err != nil {
		return err
	}
	err = netboxInventory.InitVMInterfaces()
	if err != nil {
		return err
	}
	return nil
}
