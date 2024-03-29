package constants

type SourceType string

const (
	Ovirt    SourceType = "ovirt"
	Vmware   SourceType = "vmware"
	Dnac     SourceType = "dnac"
	Proxmox  SourceType = "proxmox"
	PaloAlto SourceType = "paloalto"
)

const DefaultSourceName = "netbox-ssot"

const (
	DefaultOSName       string = "Generic OS"
	DefaultOSVersion    string = "Generic Version"
	DefaultManufacturer string = "Generic Manufacturer"
	DefaultModel        string = "Generic Model"
	DefaultSite         string = "DefaultSite"
)

type Color string

const (
	ColorDarkRed    = "aa1409"
	ColorRed        = "f44336"
	ColorPink       = "e91e63"
	ColorRose       = "ffe4e1"
	ColorFuchsia    = "ff66ff"
	ColorPurple     = "9c27b0"
	ColorDarkPurple = "673ab7"
	ColorIndigo     = "3f51b5"
	ColorBlue       = "2196f3"
	ColorLightBlue  = "03a9f4"
	ColorCyan       = "00bcd4"
	ColorTeal       = "009688"
	ColorAqua       = "00ffff"
	ColorDarkGreen  = "2f6a31"
	ColorGreen      = "4caf50"
	ColorLightGreen = "8bc34a"
	ColorLime       = "cddc39"
	ColorYellow     = "ffeb3b"
	ColorAmber      = "ffc107"
	ColorOrange     = "ff9800"
	ColorDarkOrange = "ff5722"
	ColorBrown      = "795548"
	ColorLightGrey  = "c0c0c0"
	ColorGrey       = "9e9e9e"
	ColorDarkGrey   = "607d8b"
	ColorBlack      = "111111"
	ColorWhite      = "ffffff"
)

// Default mappings of sources to colors (for tags).
var DefaultSourceToTagColorMap = map[SourceType]string{
	Ovirt:    ColorDarkRed,
	Vmware:   ColorLightGreen,
	Dnac:     ColorLightBlue,
	PaloAlto: ColorDarkOrange,
}

// Object for mapping source type to tag color.
var SourceTypeToTagColorMap = map[SourceType]string{
	Ovirt:    ColorRed,
	Vmware:   ColorGreen,
	Dnac:     ColorBlue,
	PaloAlto: ColorOrange,
}

const (
	// API timeout in seconds.
	DefaultAPITimeout = 30
)

// Magic numbers for dealing with bytes.
const (
	B   = 1
	KB  = 1000 * B
	MB  = 1000 * KB
	GB  = 1000 * MB
	TB  = 1000 * GB
	KiB = 1024 * B
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
)

// Magic numbers for dealing with IP addresses.
const (
	IPv4 = 4
	IPv6 = 6
)

const (
	HTTPSDefaultPort = 443
)

// Names used for netbox objects custom fields attribute.
const (
	// Custom Field for matching object with a source. This custom field is important
	// for priority diff.
	CustomFieldSourceName        = "source"
	CustomFieldSourceLabel       = "Source"
	CustomFieldSourceDescription = "Name of the source from which the object was collected"

	// Custom field for adding source ID for each object.
	CustomFieldSourceIDName        = "source_id"
	CustomFieldSourceIDLabel       = "Source ID"
	CustomFieldSourceIDDescription = "ID of the object on the source API"

	// Custom field dcim.device, so we can add number of cpu cores for each server.
	CustomFieldHostCPUCoresName        = "host_cpu_cores"
	CustomFieldHostCPUCoresLabel       = "Host CPU cores"
	CustomFieldHostCPUCoresDescription = "Number of CPU cores on the host"

	// Custom field for dcim.device, so we can add number of ram for each server.
	CustomFieldHostMemoryName        = "host_memory"
	CustomFieldHostMemoryLabel       = "Host memory"
	CustomFieldHostMemoryDescription = "Amount of memory on the host"
)

// Device Role constants.
const (
	DeviceRoleFirewall      = "Firewall"
	DeviceRoleFirewallColor = "f57842"

	DeviceRoleServer      = "Server"
	DeviceRoleServerColor = "00add8"

	DeviceRoleContainer      = "Container"
	DeviceRoleContainerColor = "0db7ed"
)

// Constants used for variables in our contexts.
type CtxKey int

const (
	CtxSourceKey CtxKey = iota
)

const (
	UntaggedVID = 0
	DefaultVID  = 1
	MaxVID      = 4094
	TaggedVID   = 4095
)

// Here all mappings are defined so we don't hardcode api paths of objects
// in our code.
const (
	// Tenancy paths.
	ContactGroupsAPIPath      = "/api/tenancy/contact-groups/"
	ContactRolesAPIPath       = "/api/tenancy/contact-roles/"
	ContactsAPIPath           = "/api/tenancy/contacts/"
	TenantsAPIPath            = "/api/tenancy/tenants/"
	ContactAssignmentsAPIPath = "/api/tenancy/contact-assignments/"

	// IPAM paths.
	PrefixesAPIPath    = "/api/ipam/prefixes/"
	VlanGroupsAPIPath  = "/api/ipam/vlan-groups/"
	VlansAPIPath       = "/api/ipam/vlans/"
	IPAddressesAPIPath = "/api/ipam/ip-addresses/"

	// Virtualization paths.
	ClusterTypesAPIPath    = "/api/virtualization/cluster-types/"
	ClusterGroupsAPIPath   = "/api/virtualization/cluster-groups/"
	ClustersAPIPath        = "/api/virtualization/clusters/"
	VirtualMachinesAPIPath = "/api/virtualization/virtual-machines/"
	VMInterfacesAPIPath    = "/api/virtualization/interfaces/"

	// DCIM paths.
	DevicesAPIPath               = "/api/dcim/devices/"
	DeviceRolesAPIPath           = "/api/dcim/device-roles/"
	DeviceTypesAPIPath           = "/api/dcim/device-types/"
	InterfacesAPIPath            = "/api/dcim/interfaces/"
	SitesAPIPath                 = "/api/dcim/sites/"
	ManufacturersAPIPath         = "/api/dcim/manufacturers/"
	PlatformsAPIPath             = "/api/dcim/platforms/"
	VirtualDeviceContextsAPIPath = "/api/dcim/virtual-device-contexts/"

	// Extras paths.
	CustomFieldsAPIPath = "/api/extras/custom-fields/"
	TagsAPIPath         = "/api/extras/tags/"
)
