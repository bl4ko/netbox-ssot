package constants

type SourceType string

const (
	Ovirt     SourceType = "ovirt"
	Vmware    SourceType = "vmware"
	Dnac      SourceType = "dnac"
	Proxmox   SourceType = "proxmox"
	PaloAlto  SourceType = "paloalto"
	Fortigate SourceType = "fortigate"
	FMC       SourceType = "fmc"
	IOSXE     SourceType = "ios-xe"
)

const DefaultNetboxTagColor = "00add8"
const DefaultSourceName = "netbox-ssot"

const DefaultArpTagName = "arp-entry"
const DefaultArpTagColor = ColorRed
const ArpLastSeenFormat = "2006-01-02 15:04:05"

const DefaultArpDataLifeSpan = 60 * 60 * 24 * 2 // 2 days in seconds

const (
	DefaultOSName       string = "Unknown"
	DefaultOSVersion    string = "X"
	DefaultCPUArch      string = "unknown"
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

// Default mappings of sources to colors (for tags), fallback mechanism.
// E.g. we name a source "prodvmware", tag "Source: prodvmware" is created
// with our color.
var SourceTagColorMap = map[SourceType]string{
	Ovirt:     ColorDarkRed,
	Vmware:    ColorLightGreen,
	Dnac:      ColorLightBlue,
	PaloAlto:  ColorDarkOrange,
	Fortigate: ColorDarkGreen,
	FMC:       ColorLightBlue,
	IOSXE:     "0d294f",
}

// Each source Mapping for source type tag. E.g. tag "paloalto" -> color orange.
var SourceTypeTagColorMap = map[SourceType]string{
	Ovirt:     ColorRed,
	Vmware:    ColorGreen,
	Dnac:      ColorBlue,
	PaloAlto:  ColorOrange,
	Fortigate: ColorDarkGreen,
	FMC:       ColorBlue,
	IOSXE:     "0d294f",
}

const (
	// API timeout in seconds.
	DefaultAPITimeout = 15
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
	IPv4            = 4
	IPv6            = 6
	MaxIPv4MaskBits = 32
	MaxIPv6MaskBits = 128
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

	// Custom field for dcim.device, so we can store uuid for it.
	CustomFieldDeviceUUIDName        = "uuid"
	CustomFieldDeviceUUIDLabel       = "uuid"
	CustomFieldDeviceUUIDDescription = "Universally Unique Identifier for a device"

	// Custom field for ObjectTypeIPAddress, so we can determine if an ip is part of an arp table or not.
	CustomFieldArpEntryName        = "arp_entry"
	CustomFieldArpEntryLabel       = "Arp Entry"
	CustomFieldArpEntryDescription = "Was this IP collected from ARP table"

	// Custom field for ipam.ipaddress, so we can track when was arp entry last found.
	CustomFieldArpIPLastSeenName        = "last_seen"
	CustomFieldArpIPLastSeenLabel       = "Last seen"
	CustomFieldArpIPLastSeenDescription = "Last time the IP was found in the arp table"
)

// Device Role constants.
const (
	DeviceRoleFirewall            = "Firewall"
	DeviceRoleFirewallDescription = "Device role for marking firewall device."
	DeviceRoleFirewallColor       = "f57842"

	DeviceRoleSwitch            = "Switch"
	DeviceRoleSwitchDescription = "Device role for marking switch device."
	DeviceRoleSwitchColor       = "7aefea"

	DeviceRoleServer            = "Server"
	DeviceRoleServerDescription = "Device role for marking server."
	DeviceRoleServerColor       = "00add8"

	DeviceRoleContainer            = "Container"
	DeviceRoleContainerDescription = "VM role for separating containers from VMs."
	DeviceRoleContainerColor       = "0db7ed"

	DeviceRoleVM            = "VM"
	DeviceRoleVMDescription = "Role for representing VMs."
	DeviceRoleVMColor       = "81eaea"

	DeviceRoleVMTemplate            = "VM Template"
	DeviceRoleVMTemplateDescription = "VM role for separating VM templates from VMs."
	DeviceRoleVMTemplateColor       = "82c1ea"
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

var Arch2Bit = map[string]string{
	"x86_64":  "64-bit",
	"i386":    "32-bit",
	"i486":    "32-bit",
	"i586":    "32-bit",
	"i686":    "32-bit",
	"aarch64": "64-bit",
	"arm64":   "64-bit",
	"arm":     "32-bit",
	"arm32":   "32-bit",
	"ppc64le": "64-bit",
	"s390x":   "64-bit",
	"mips64":  "64-bit",
	"riscv64": "64-bit",
	"unknown": "unknown",
}
