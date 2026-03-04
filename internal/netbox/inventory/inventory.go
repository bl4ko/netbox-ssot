package inventory

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/logger"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
	"github.com/src-doo/netbox-ssot/internal/netbox/service"
	"github.com/src-doo/netbox-ssot/internal/parser"
	"github.com/src-doo/netbox-ssot/internal/utils"
)

// NetboxInventory is a singleton class to manage a inventory of NetBoxObject objects.
type NetboxInventory struct {
	// Logger is the logger used for logging messages
	Logger *logger.Logger
	// NetboxConfig is the Netbox configuration
	NetboxConfig *parser.NetboxConfig
	// NetboxAPI is the Netbox API object, for communicating with the Netbox API
	NetboxAPI *service.NetboxClient
	// SourcePriority: if object is found on multiple sources, which source has
	// the priority for the object attributes.
	SourcePriority map[string]int
	// ArpDataLifeSpan determines the lifespan of arp entries in seconds.
	ArpDataLifeSpan int
	// OrphanManager object that manages orphaned objects.
	OrphanManager *OrphanManager
	// Tag used by netbox-ssot to mark devices that are managed by it.
	SsotTag *objects.Tag
	// Tag used by netbox-ssot to preserve manually set device type
	IgnoreDeviceTypeTag *objects.Tag
	// Default context for the inventory, we use it to pass sourcename
	// to functions for logging.
	Ctx context.Context //nolint:containedctx

	// tagsIndexByName is a map of all tags in the Netbox's inventory,
	// indexed by their name
	tagsIndexByName map[string]*objects.Tag
	tagsLock        sync.Mutex

	// contactGroupsIndexByName is a map of all contact groups
	// indexed by their names.
	contactGroupsIndexByName map[string]*objects.ContactGroup
	contactGroupsLock        sync.Mutex

	// contactRolesIndexByName is a map of all contact roles
	// indexed by their names.
	contactRolesIndexByName map[string]*objects.ContactRole
	contactRolesLock        sync.Mutex

	// contactsIndexByName is a map of all contacts in the Netbox's inventory,
	// indexed by their names
	contactsIndexByName map[string]*objects.Contact
	contactsLock        sync.Mutex

	// contactAssignmentsIndex is a
	// map of all contact assignments indexed by their
	// content type, object id, contact id and role id.
	contactAssignmentsIndex map[constants.ContentType]map[int]map[int]map[int]*objects.ContactAssignment
	contactAssignmentsLock  sync.Mutex

	// sitesIndexByName is a map of all sites in the Netbox's inventory,
	// indexed by their name
	sitesIndexByName map[string]*objects.Site
	sitesLock        sync.Mutex

	// siteGroupsIndexByName is a map of all site groups in the Netbox's inventory,
	// indexed by their name
	siteGroupsIndexByName map[string]*objects.SiteGroup
	siteGroupsLock        sync.Mutex

	// manufacturersIndexByName is a map of all manufacturers in the Netbox's inventory,
	// indexed by their name
	manufacturersIndexByName map[string]*objects.Manufacturer
	manufacturersLock        sync.Mutex

	// platformsIndexByName is a map of all platforms in the Netbox's inventory, indexed by their name
	platformsIndexByName map[string]*objects.Platform
	platformsLock        sync.Mutex

	// tenantsIndexByName is a map of all tenants in the Netbox's inventory,
	// indexed by their name
	tenantsIndexByName map[string]*objects.Tenant
	tenantsLock        sync.Mutex

	// deviceTypesIndexByModel is a map of all device types in the Netbox's inventory,
	// indexed by their model
	deviceTypesIndexByModel map[string]*objects.DeviceType
	deviceTypesLock         sync.Mutex

	// devicesIndexByNameAndSiteID is a map of all devices in the Netbox's inventory,
	// indexed by their name and SiteID
	devicesIndexByNameAndSiteID map[string]map[int]*objects.Device
	// devicesIndexByID is a helper index, that we use in init functions
	// to create relationships between objects.
	devicesIndexByID map[int]*objects.Device
	devicesLock      sync.Mutex

	// virtualDeviceContextsIndex is a map of all virtual device contexts
	// in the Netbox's inventory indexed by their name and device ID.
	virtualDeviceContextsIndex map[string]map[int]*objects.VirtualDeviceContext
	virtualDeviceContextsLock  sync.Mutex

	// prefixesIndexByPrefix is a map of all prefixes in the Netbox's inventory,
	// indexed by their prefix and VRF
	prefixesIndexByPrefix map[string]map[int]*objects.Prefix
	prefixesLock          sync.Mutex

	// vrfsIndexByName is a map of all VRFs in the Netbox's inventory,
	// indexed by their name.
	vrfsIndexByName map[string]*objects.VRF
	vrfsLock        sync.Mutex

	// vlanGroupsIndexByName is a map of all VlanGroups in the Netbox's inventory,
	// indexed by their name.
	vlanGroupsIndexByName map[string]*objects.VlanGroup
	vlanGroupsLock        sync.Mutex

	// vlansIndexByVlanGroupIDAndVID is a map of all vlans in the Netbox's inventory,
	// indexed by their VlanGroup and vid.
	vlansIndexByVlanGroupIDAndVID map[int]map[int]*objects.Vlan
	vlansLock                     sync.Mutex

	// clusterGroupsIndexByName is a map of all cluster groups in the Netbox's
	// inventory indexed by their name
	clusterGroupsIndexByName map[string]*objects.ClusterGroup
	clusterGroupsLock        sync.Mutex

	// clusterTypesIndexByName is a map of all cluster types in the Netbox's
	// inventory, indexed by their name
	clusterTypesIndexByName map[string]*objects.ClusterType
	clusterTypesLock        sync.Mutex

	// clustersIndexByName is a map of all clusters in the Netbox's inventory,
	// indexed by their name
	clustersIndexByName map[string]*objects.Cluster
	clustersLock        sync.Mutex

	// Netbox's Device Roles is a map of all device roles in the inventory,
	// indexed by name.
	deviceRolesIndexByName map[string]*objects.DeviceRole
	deviceRolesLock        sync.Mutex

	// customFieldsIndexByName is a map of all custom fields in the inventory,
	// indexed by name.
	customFieldsIndexByName map[string]*objects.CustomField
	customFieldsLock        sync.Mutex

	// InterfacesIndexByDeviceAnName is a map of all interfaces in the inventory,
	// indexed by their's device id and their name.
	interfacesIndexByDeviceIDAndName map[int]map[string]*objects.Interface
	// interfacesIndexByID is a helper index, that we use in init functions
	// to create relationships between objects.
	interfacesIndexByID map[int]*objects.Interface
	interfacesLock      sync.Mutex

	// vmsIndexByNameAndClusterID is a map of all virtual machines in the inventory,
	// indexed by their name and their cluster id
	vmsIndexByNameAndClusterID map[string]map[int]*objects.VM
	// vmsIndexByID is a helper index, that we use in init functions
	// to create relationships between objects.
	vmsIndexByID map[int]*objects.VM
	vmsLock      sync.Mutex

	// vmInterfacesIndexByVMAndName is a map of all virtual machine interfaces in the
	// inventory, indexed by their's virtual machine id and their name
	vmInterfacesIndexByVMIdAndName map[int]map[string]*objects.VMInterface
	// vmInterfacesIndexByID is a helper index, that we use in init functions
	// to create relationships between objects.
	vmInterfacesIndexByID map[int]*objects.VMInterface
	vmInterfacesLock      sync.Mutex

	// ipAdressesIndex is a map of all IP addresses in the inventory,
	// indexed:
	//   * iface type (vmiface or device iface)
	//   * iface name (name of the vminterface/deviceinterface)
	//   * iface parent name (name of the vm/device that the interface belongs to)
	//   * ip address
	ipAddressesIndex map[constants.ContentType]map[string]map[string]map[string]*objects.IPAddress
	ipAddressesLock  sync.Mutex

	// macAddressesIndex is a map of all MAC addresses in the inventory,
	// indexed with these levels:
	//  * iface type (vmiface or device iface)
	//  * iface name (name of the vminterface/deviceinterface)
	//  * iface parent name (name of the vm/device that the interface belongs to)
	//  * mac address
	macAddressesIndex map[constants.ContentType]map[string]map[string]map[string]*objects.MACAddress
	macAddressesLock  sync.Mutex

	// wirelessLANGroupsIndexByName is a map of all wireless lan groups in the Netbox's
	// inventory, indexed by their name
	wirelessLANGroupsIndexByName map[string]*objects.WirelessLANGroup
	wirelessLANGroupsLock        sync.Mutex

	// wirelessLANsIndexBySSID is a map of all wireless lans in the Netbox's inventory,
	// indexed by their ssid
	wirelessLANsIndexBySSID map[string]*objects.WirelessLAN
	wirelessLANsLock        sync.Mutex

	// virtualDisksIndex is a map of all virtual disks in the inventory,
	// indexed by their vm's id and vm's name
	virtualDisksIndexByVMIDAndName map[int]map[string]*objects.VirtualDisk
	virtualDisksLock               sync.Mutex
}

// Func string representation.
func (nbi *NetboxInventory) String() string {
	return fmt.Sprintf(
		"NetBoxInventory{Logger: %+v, NetboxConfig: %+v...}",
		nbi.Logger,
		nbi.NetboxConfig,
	)
}

// NewNetboxInventory creates a new NetBoxInventory object.
// It takes a logger and a NetboxConfig as parameters, and returns a pointer to the newly created NetBoxInventory.
// The logger is used for logging messages, and the NetboxConfig is used to configure the NetBoxInventory.
func NewNetboxInventory(
	ctx context.Context,
	logger *logger.Logger,
	nbConfig *parser.NetboxConfig,
) *NetboxInventory {
	sourcePriority := make(map[string]int, len(nbConfig.SourcePriority))
	for i, sourceName := range nbConfig.SourcePriority {
		sourcePriority[sourceName] = i
	}
	orphanManager := NewOrphanManager(logger)

	nbi := &NetboxInventory{
		Ctx:            ctx,
		Logger:         logger,
		NetboxConfig:   nbConfig,
		SourcePriority: sourcePriority,
		OrphanManager:  orphanManager,
	}
	return nbi
}

// Init function that initializes the NetBoxInventory object with objects from Netbox.
func (nbi *NetboxInventory) Init() error {
	baseURL := fmt.Sprintf(
		"%s://%s:%d",
		nbi.NetboxConfig.HTTPScheme,
		nbi.NetboxConfig.Hostname,
		nbi.NetboxConfig.Port,
	)

	nbi.Logger.Debug(nbi.Ctx, "Initializing Netbox API with baseURL: ", baseURL)
	var err error
	nbi.NetboxAPI, err = service.NewNetboxClient(
		nbi.Logger,
		baseURL,
		nbi.NetboxConfig.APIToken,
		nbi.NetboxConfig.ValidateCert,
		nbi.NetboxConfig.Timeout,
		nbi.NetboxConfig.CAFile,
	)
	if err != nil {
		return fmt.Errorf("create new netbox client: %s", err)
	}

	err = nbi.checkVersion()
	if err != nil {
		return err
	}

	// WARNING: Order matters
	initFunctions := []func(context.Context) error{
		nbi.initCustomFields,
		nbi.initSsotCustomFields,
		nbi.initTags,
		nbi.initContactGroups,
		nbi.initContactRoles,
		nbi.initAdminContactRole,
		nbi.initContacts,
		nbi.initContactAssignments,
		nbi.initTenants,
		nbi.initSiteGroups,
		nbi.initSites,
		nbi.initDefaultSite,
		nbi.initManufacturers,
		nbi.initPlatforms,
		nbi.initVMs,
		nbi.initVirtualDisks,
		nbi.initVMInterfaces,
		nbi.initDevices,
		nbi.initInterfaces,
		nbi.initIPAddresses,
		nbi.initMACAddresses,
		nbi.initVlanGroups,
		nbi.initPrefixes,
		nbi.initVRFs,
		nbi.initVlans,
		nbi.initDeviceRoles,
		nbi.initDeviceTypes,
		nbi.initClusterGroups,
		nbi.initClusterTypes,
		nbi.initClusters,
		nbi.initVirtualDeviceContexts,
		nbi.initWirelessLANs,
		nbi.initWirelessLANGroups,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(nbi.Ctx); err != nil {
			return fmt.Errorf("%s: %s", err, utils.ExtractFunctionName(initFunc))
		}
		duration := time.Since(startTime)
		nbi.Logger.Infof(
			nbi.Ctx,
			"Successfully initialized %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"),
			duration.Seconds(),
		)
	}

	return nil
}

func (nbi *NetboxInventory) checkVersion() error {
	version, err := service.GetVersion(nbi.Ctx, nbi.NetboxAPI)
	if err != nil {
		return fmt.Errorf("get version: %s", err)
	}
	supportedVersion := 4
	supportedMinorVersion := 2
	versionComponents := strings.Split(version, ".")
	majorVersion, err := strconv.Atoi(versionComponents[0])
	if err != nil {
		return fmt.Errorf("parse major version: %s", err)
	}
	minorVersion, err := strconv.Atoi(versionComponents[1])
	if err != nil {
		return fmt.Errorf("parse minor version: %s", err)
	}
	if majorVersion != supportedVersion || minorVersion < supportedMinorVersion {
		return fmt.Errorf(
			"this version of netbox-ssot works only with netbox version >= 4.2.0, but received version: %s",
			version,
		)
	}
	return nil
}
