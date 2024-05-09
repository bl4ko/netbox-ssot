package inventory

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// NetboxInventory is a singleton class to manage a inventory of NetBoxObject objects.
type NetboxInventory struct {
	// Logger is the logger used for logging messages
	Logger *logger.Logger
	// NetboxConfig is the Netbox configuration
	NetboxConfig *parser.NetboxConfig
	// NetboxAPI is the Netbox API object, for communicating with the Netbox API
	NetboxAPI *service.NetboxClient
	// SourcePriority: if object is found on multiple sources, which source has the priority for the object attributes.
	SourcePriority map[string]int
	// TagsIndexByName is a map of all tags in the Netbox's inventory, indexed by their name
	TagsIndexByName map[string]*objects.Tag
	// ContactGroupsIndexByName is a map of all contact groups indexed by their names.
	ContactGroupsIndexByName map[string]*objects.ContactGroup
	// ContactRolesIndexByName is a map of all contact roles indexed by their names.
	ContactRolesIndexByName map[string]*objects.ContactRole
	// ContactsIndexByName is a map of all contacts in the Netbox's inventory, indexed by their names
	ContactsIndexByName map[string]*objects.Contact
	// ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID is a map of all contact assignments indexed by their content type, object id, contact id and role id.
	ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID map[objects.ObjectType]map[int]map[int]map[int]*objects.ContactAssignment
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
	// DevicesIndexByNameAndSiteID is a map of all devices in the Netbox's inventory, indexed by their name, and
	// site ID (This is because, netbox constraints: https://github.com/netbox-community/netbox/blob/3d941411d438f77b66d2036edf690c14b459af58/netbox/dcim/models/devices.py#L775)
	DevicesIndexByNameAndSiteID map[string]map[int]*objects.Device
	// VirtualDeviceContextsIndexByNameAndDeviceID is a map of all virtual device contexts in the Netbox's inventory indexed by their name and device ID.
	VirtualDeviceContextsIndexByNameAndDeviceID map[string]map[int]*objects.VirtualDeviceContext
	// PrefixesIndexByPrefix is a map of all prefixes in the Netbox's inventory, indexed by their prefix
	PrefixesIndexByPrefix map[string]*objects.Prefix
	// VlanGroupsIndexByName is a map of all VlanGroups in the Netbox's inventory, indexed by their name
	VlanGroupsIndexByName map[string]*objects.VlanGroup
	// VlansIndexByVlanGroupIDAndVID is a map of all vlans in the Netbox's inventory, indexed by their VlanGroup and vid.
	VlansIndexByVlanGroupIDAndVID map[int]map[int]*objects.Vlan
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
	InterfacesIndexByDeviceIDAndName map[int]map[string]*objects.Interface
	// VMsIndexByNameAndClusterID is a map of all virtual machines in the inventory, indexed by their name and their cluster id
	VMsIndexByNameAndClusterID map[string]map[int]*objects.VM
	// VirtualMachineInterfacesIndexByVMAndName is a map of all virtual machine interfaces in the inventory, indexed by their's virtual machine id and their name
	VMInterfacesIndexByVMIdAndName map[int]map[string]*objects.VMInterface
	// IPAdressesIndexByAddress is a map of all IP addresses in the inventory, indexed by their address
	IPAdressesIndexByAddress map[string]*objects.IPAddress

	// We also store locks for all objects, so inventory can be updated by multiple parallel goroutines
	TenantsLock            sync.Mutex
	TagsLock               sync.Mutex
	SitesLock              sync.Mutex
	ContactRolesLock       sync.Mutex
	ContactGroupsLock      sync.Mutex
	ContactsLock           sync.Mutex
	ContactAssignmentsLock sync.Mutex
	CustomFieldsLock       sync.Mutex
	ClusterGroupsLock      sync.Mutex
	ClusterTypesLock       sync.Mutex
	ClustersLock           sync.Mutex
	DeviceRolesLock        sync.Mutex
	ManufacturersLock      sync.Mutex
	DeviceTypesLock        sync.Mutex
	PlatformsLock          sync.Mutex
	DevicesLock            sync.Mutex
	VlanGroupsLock         sync.Mutex
	VlansLock              sync.Mutex
	InterfacesLock         sync.Mutex
	VMsLock                sync.Mutex
	VMInterfacesLock       sync.Mutex
	IPAddressesLock        sync.Mutex
	PrefixesLock           sync.Mutex

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

	// ArpDataLifeSpan determines the lifespan of arp entries in seconds.
	ArpDataLifeSpan int
	// Tag used by netbox-ssot to mark devices that are managed by it.
	SsotTag *objects.Tag
	// Default context for the inventory, we use it to pass sourcename to functions for logging.
	Ctx context.Context //nolint:containedctx
}

// Func string representation.
func (nbi *NetboxInventory) String() string {
	return fmt.Sprintf("NetBoxInventory{Logger: %+v, NetboxConfig: %+v...}", nbi.Logger, nbi.NetboxConfig)
}

// NewNetboxInventory creates a new NetBoxInventory object.
// It takes a logger and a NetboxConfig as parameters, and returns a pointer to the newly created NetBoxInventory.
// The logger is used for logging messages, and the NetboxConfig is used to configure the NetBoxInventory.
func NewNetboxInventory(ctx context.Context, logger *logger.Logger, nbConfig *parser.NetboxConfig) *NetboxInventory {
	sourcePriority := make(map[string]int, len(nbConfig.SourcePriority))
	for i, sourceName := range nbConfig.SourcePriority {
		sourcePriority[sourceName] = i
	}
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
	}
	nbi := &NetboxInventory{Ctx: ctx, Logger: logger, NetboxConfig: nbConfig, SourcePriority: sourcePriority, OrphanManager: make(map[string]map[int]bool), OrphanObjectPriority: orphanObjectPriority}
	return nbi
}

// Init function that initializes the NetBoxInventory object with objects from Netbox.
func (nbi *NetboxInventory) Init() error {
	baseURL := fmt.Sprintf("%s://%s:%d", nbi.NetboxConfig.HTTPScheme, nbi.NetboxConfig.Hostname, nbi.NetboxConfig.Port)

	nbi.Logger.Debug(nbi.Ctx, "Initializing Netbox API with baseURL: ", baseURL)
	var err error
	nbi.NetboxAPI, err = service.NewNetboxClient(nbi.Logger, baseURL, nbi.NetboxConfig.APIToken, nbi.NetboxConfig.ValidateCert, nbi.NetboxConfig.Timeout, nbi.NetboxConfig.CAFile)
	if err != nil {
		return fmt.Errorf("create new netbox client: %s", err)
	}

	err = nbi.checkVersion()
	if err != nil {
		return err
	}

	// Order matters. TODO: use parallelization in the future, on the init functions that can be parallelized
	initFunctions := []func(context.Context) error{
		nbi.InitCustomFields,
		nbi.InitSsotCustomFields,
		nbi.InitTags,
		nbi.InitContactGroups,
		nbi.InitContactRoles,
		nbi.InitAdminContactRole,
		nbi.InitContacts,
		nbi.InitContactAssignments,
		nbi.InitTenants,
		nbi.InitSites,
		nbi.InitDefaultSite,
		nbi.InitManufacturers,
		nbi.InitPlatforms,
		nbi.InitDevices,
		nbi.InitVirtualDeviceContexts,
		nbi.InitInterfaces,
		nbi.InitIPAddresses,
		nbi.InitVlanGroups,
		nbi.InitDefaultVlanGroup,
		nbi.InitPrefixes,
		nbi.InitVlans,
		nbi.InitDeviceRoles,
		nbi.InitDeviceTypes,
		nbi.InitClusterGroups,
		nbi.InitClusterTypes,
		nbi.InitClusters,
		nbi.InitVMs,
		nbi.InitVMInterfaces,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(nbi.Ctx); err != nil {
			return fmt.Errorf("%s: %s", err, utils.ExtractFunctionName(initFunc))
		}
		duration := time.Since(startTime)
		nbi.Logger.Infof(nbi.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionName(initFunc), duration.Seconds())
	}

	return nil
}

func (nbi *NetboxInventory) checkVersion() error {
	version, err := service.GetVersion(nbi.Ctx, nbi.NetboxAPI)
	if err != nil {
		return fmt.Errorf("get version: %s", err)
	}
	supportedVersion := 4
	versionComponents := strings.Split(version, ".")
	majorVersion, err := strconv.Atoi(versionComponents[0])
	if err != nil {
		return fmt.Errorf("parse major version: %s", err)
	}
	if majorVersion < supportedVersion {
		return fmt.Errorf("this version of netbox-ssot works only with netbox version > 4.x.x, but received version: %s", version)
	}
	return nil
}
