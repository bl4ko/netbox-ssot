package objects

import (
	"fmt"
)

type SiteStatus struct {
	Choice
}

var (
	SiteStatusActive          = SiteStatus{Choice{Value: "active", Label: "Active"}}
	SiteStatusOffline         = SiteStatus{Choice{Value: "offline", Label: "Offline"}}
	SiteStatusPlanned         = SiteStatus{Choice{Value: "planned", Label: "Planned"}}
	SiteStatusStaged          = SiteStatus{Choice{Value: "staged", Label: "Staged"}}
	SiteStatusFailed          = SiteStatus{Choice{Value: "failed", Label: "Failed"}}
	SiteStatusInventory       = SiteStatus{Choice{Value: "inventory", Label: "Inventory"}}
	SiteStatusDecommissioning = SiteStatus{Choice{Value: "decommissioning", Label: "Decommissioning"}}
)

// Site ares used for functional groupings.
// A site usually represents a building within a region.
type Site struct {
	NetboxObject
	// Full name for the site. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Site status. This field is required.
	Status *SiteStatus `json:"status,omitempty"`
}

// Platform represents an operating system or other software platform which may be running on a device.
type Platform struct {
	NetboxObject
	// Name of the platform. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Manufacturer is the manufacturer of the platform.
	Manufacturer *Manufacturer `json:"manufacturer,omitempty"`
}

type Region struct {
	NetboxObject
	// Name is the name of the region. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is a URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
}

// Location represents a physical location, such as a floor or room in a building.
type Location struct {
	NetboxObject
	// Site is the site to which the location belongs. This field is required.
	Site *Site
	// Name is the name of the location. This field is required.
	Name string
	// URL-friendly unique shorthand. This field is required.
	Slug string
	// Status is the status of the location. This field is required.
	Status *SiteStatus
}

// Manufacturer represents a hardware manufacturer (e.g. Cisco, HP, ...).
type Manufacturer struct {
	NetboxObject
	// Name of the manufacturer (e.g. Cisco). This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
}

var ManufacturerMap = map[string]string{
	"^AMD$":               "AMD",
	".*Broadcom.*":        "Broadcom",
	".*Cisco.*":           "Cisco",
	".*Dell.*":            "Dell",
	"FTS Corp":            "Fujitsu",
	".*Fujitsu.*":         "Fujitsu",
	".*HiSilicon.*":       "HiSilicon",
	"^HP$":                "HPE",
	"^HPE$":               "HPE",
	".*Huawei.*":          "Huawei",
	".*Hynix.*":           "Hynix",
	".*Inspur.*":          "Inspur",
	".*Intel.*":           "Intel",
	"LEN":                 "Lenovo",
	".*Lenovo.*":          "Lenovo",
	".*Micron.*":          "Micron",
	".*Nvidea.*":          "Nvidia",
	".*Samsung.*":         "Samsung",
	".*Supermicro.*":      "Supermicro",
	".*Toshiba.*":         "Toshiba",
	"^WD$":                "Western Digital",
	".*Western Digital.*": "Western Digital",
}

// DeviceType represents the physical and operational characteristics of a device.
// For example, a device type may represent a Cisco C2960 switch running IOS 15.2.
type DeviceType struct {
	NetboxObject
	// Manufacturer is the manufacturer of the device type. This field is required.
	Manufacturer *Manufacturer `json:"manufacturer,omitempty"`
	// Model is the model of the device type. This field is required.
	Model string `json:"model,omitempty"`
	// Slug is a URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
}

// DeviceRole represents the functional role of a device.
// For example, a device may play the role of a router, a switch, a firewall, etc.
type DeviceRole struct {
	NetboxObject
	// Name is the name of the device role. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Color of the device role. This field is required.
	Color string `json:"color,omitempty"`
	// VMRole is whether this device role is used to represent virtual machines.
	VMRole bool `json:"vm_role,omitempty"`
}

// https://github.com/netbox-community/netbox/blob/b93735861d3bde0354c855a8bbd2a2311e8eb920/netbox/dcim/choices.py#L182
// Predefined airflow types from netbox
type DeviceAirFlowType struct {
	Choice
}

var (
	FrontToRear = DeviceAirFlowType{Choice{Value: "front-to-rear", Label: "Front to rear"}}
	RearToFront = DeviceAirFlowType{Choice{Value: "rear-to-front", Label: "Rear to front"}}
	LeftToRight = DeviceAirFlowType{Choice{Value: "left-to-right", Label: "Left to right"}}
	RightToLeft = DeviceAirFlowType{Choice{Value: "right-to-left", Label: "Right to left"}}
	SideToRear  = DeviceAirFlowType{Choice{Value: "side-to-rear", Label: "Side to rear"}}
	Passive     = DeviceAirFlowType{Choice{Value: "passive", Label: "Passive"}}
	Mixed       = DeviceAirFlowType{Choice{Value: "mixed", Label: "Mixed"}}
)

type DeviceStatus struct {
	Choice
}

var (
	DeviceStatusOffline         = DeviceStatus{Choice{Value: "offline", Label: "Offline"}}
	DeviceStatusActive          = DeviceStatus{Choice{Value: "active", Label: "Active"}}
	DeviceStatusPlanned         = DeviceStatus{Choice{Value: "planned", Label: "Planned"}}
	DeviceStatusStaged          = DeviceStatus{Choice{Value: "staged", Label: "Staged"}}
	DeviceStatusFailed          = DeviceStatus{Choice{Value: "failed", Label: "Failed"}}
	DeviceStatusInventory       = DeviceStatus{Choice{Value: "inventory", Label: "Inventory"}}
	DeviceStatusDecommissioning = DeviceStatus{Choice{Value: "decommissioning", Label: "Decommissioning"}}
)

// Device can be any piece of physical hardware, such as a server, router, or switch.
type Device struct {
	NetboxObject

	// Device
	// Name is the name of the device.
	Name string `json:"name,omitempty"`
	// DeviceRole is the functional role of the device. This field is required.
	DeviceRole *DeviceRole `json:"role,omitempty"`

	// Hardware
	// DeviceType is the type of device. This field is required.
	DeviceType *DeviceType `json:"device_type,omitempty"`
	// Airflow is the airflow pattern of the device.
	Airflow *DeviceAirFlowType `json:"airflow,omitempty"`
	// Status is the status of the device.
	SerialNumber string `json:"serial,omitempty"`
	// AssetTag is an unique tag for identifying the device.
	AssetTag string `json:"asset_tag,omitempty"`

	// Location
	// Site is the site to which the device belongs. This field is required.
	Site *Site `json:"site,omitempty"`
	// Location is the location of the device.
	Location *Location `json:"location,omitempty"`

	// Management
	// Status of the device (e.g. active, offline, planned, etc.). This field is required.
	Status *DeviceStatus `json:"status,omitempty"`
	// Platform of the device (e.g. Cisco IOS, Dell OS9, etc.).
	Platform *Platform `json:"platform,omitempty"`

	// Primary IPv4

	// Primary IPv6
	// Out-of-band IP

	// Virtualization
	// Cluster is the cluster to which the device belongs. (e.g. VMWare server belonging to a specific cluster).
	Cluster *Cluster `json:"cluster,omitempty"`

	// Tenancy
	// Tenant group
	Tenant *Tenant `json:"tenant,omitempty"`

	// Virtual Chassis

	// The position in the virtual chassis this device is identified by
	// Position

	// The priority of the device in the virtual chassis
	// Priority
	// Additional comments.
	Comments string `json:"comments,omitempty"`

	// CustomFields is a dictionary of custom fields defined for the device type. map[customFieldName]: valueStr
	CustomFields map[string]string `json:"custom_fields,omitempty"`
}

type InterfaceType struct {
	Choice
}

// Predefined types: https://github.com/netbox-community/netbox/blob/ec245b968f50bdbafaadd5d6b885832d858fa167/netbox/dcim/choices.py#L800
var (
	// VirtualInterfaces
	VirtualInterfaceType = InterfaceType{Choice{Value: "virtual", Label: "Virtual"}}
	BridgeInterfaceType  = InterfaceType{Choice{Value: "bridge", Label: "Bridge"}}
	LAGInterfaceType     = InterfaceType{Choice{Value: "lag", Label: "Link Aggregation Group (LAG)"}}

	// Ethernet (Fixed)
	BASEFXInterfaceType       = InterfaceType{Choice{Value: "100base-fx", Label: "100BASE-FX (10/100ME FIBER)"}}
	BASELFXInterfaceType      = InterfaceType{Choice{Value: "100base-lfx", Label: "100BASE-LFX (10/100ME FIBER)"}}
	BASETXInterfaceType       = InterfaceType{Choice{Value: "100base-tx", Label: "100BASE-TX (10/100ME)"}}
	BASET1InterfaceType       = InterfaceType{Choice{Value: "100base-t1", Label: "100BASE-T1 (10/100ME Single Pair)"}}
	GE1FixedInterfaceType     = InterfaceType{Choice{Value: "1000base-t", Label: "1000BASE-T (1GE)"}}
	GE1GBICInterfaceType      = InterfaceType{Choice{Value: "1000base-x-gbic", Label: "GBIC (1GE)"}}
	GE1SFPInterfaceType       = InterfaceType{Choice{Value: "1000base-x-sfp", Label: "SFP (1GE)"}}
	GE2FixedInterfaceType     = InterfaceType{Choice{Value: "2.5gbase-t", Label: "2.5GBASE-T (2.5GE)"}}
	GE5FixedInterfaceType     = InterfaceType{Choice{Value: "5gbase-t", Label: "5GBASE-T (5GE)"}}
	GE10FixedInterfaceType    = InterfaceType{Choice{Value: "10gbase-t", Label: "10GBASE-T (10GE)"}}
	GE10CX4InterfaceType      = InterfaceType{Choice{Value: "10gbase-cx4", Label: "10GBASE-CX4 (10GE)"}}
	GE10SFPPInterfaceType     = InterfaceType{Choice{Value: "10gbase-x-sfpp", Label: "SFP+ (10GE)"}}
	GE10XFPInterfaceType      = InterfaceType{Choice{Value: "10gbase-x-xfp", Label: "XFP (10GE)"}}
	GE10XENPAKInterfaceType   = InterfaceType{Choice{Value: "10gbase-x-xenpak", Label: "XENPAK (10GE)"}}
	GE10X2InterfaceType       = InterfaceType{Choice{Value: "10gbase-x-x2", Label: "X2 (10GE)"}}
	GE25SFP28InterfaceType    = InterfaceType{Choice{Value: "25gbase-x-sfp28", Label: "SFP28 (25GE)"}}
	GE50SFP56InterfaceType    = InterfaceType{Choice{Value: "50gbase-x-sfp56", Label: "SFP56 (50GE)"}}
	GE40QSFPPlusInterfaceType = InterfaceType{Choice{Value: "40gbase-x-qsfpp", Label: "QSFP+ (40GE)"}}
	GE50QSFP28InterfaceType   = InterfaceType{Choice{Value: "50gbase-x-sfp28", Label: "QSFP28 (50GE)"}}
	GE100CFPInterfaceType     = InterfaceType{Choice{Value: "100gbase-x-cfp", Label: "CFP (100GE)"}}
	GE100CFP2InterfaceType    = InterfaceType{Choice{Value: "100gbase-x-cfp2", Label: "CFP2 (100GE)"}}
	GE100CFP4InterfaceType    = InterfaceType{Choice{Value: "100gbase-x-cfp4", Label: "CFP4 (100GE)"}}
	GE100CXPInterfaceType     = InterfaceType{Choice{Value: "100gbase-x-cxp", Label: "CXP (100GE)"}}
	GE100CPAKInterfaceType    = InterfaceType{Choice{Value: "100gbase-x-cpak", Label: "Cisco CPAK (100GE)"}}
	GE100DSFPInterfaceType    = InterfaceType{Choice{Value: "100gbase-x-dsfp", Label: "DSFP (100GE)"}}
	GE100SFPDDInterfaceType   = InterfaceType{Choice{Value: "100gbase-x-sfpdd", Label: "SFP-DD (100GE)"}}
	GE100QSFP28InterfaceType  = InterfaceType{Choice{Value: "100gbase-x-qsfp28", Label: "QSFP28 (100GE)"}}
	GE100QSFPDDInterfaceType  = InterfaceType{Choice{Value: "100gbase-x-qsfpdd", Label: "QSFP-DD (100GE)"}}
	GE200CFP2InterfaceType    = InterfaceType{Choice{Value: "200gbase-x-cfp2", Label: "CFP2 (200GE)"}}
	GE200QSFP56InterfaceType  = InterfaceType{Choice{Value: "200gbase-x-qsfp56", Label: "QSFP56 (200GE)"}}
	GE200QSFPDDInterfaceType  = InterfaceType{Choice{Value: "200gbase-x-qsfpdd", Label: "QSFP-DD (200GE)"}}
	GE400CFP2InterfaceType    = InterfaceType{Choice{Value: "400gbase-x-cfp2", Label: "CFP2 (400GE)"}}
	GE400QSFP112InterfaceType = InterfaceType{Choice{Value: "400gbase-x-qsfp112", Label: "QSFP112 (400GE)"}}
	GE400QSFPDDInterfaceType  = InterfaceType{Choice{Value: "400gbase-x-qsfpdd", Label: "QSFP-DD (400GE)"}}
	GE400OSFPInterfaceType    = InterfaceType{Choice{Value: "400gbase-x-osfp", Label: "OSFP (400GE)"}}
	GE400OSFPRHSInterfaceType = InterfaceType{Choice{Value: "400gbase-x-osfp-rhs", Label: "OSFP-RHS (400GE)"}}
	GE400CDFPInterfaceType    = InterfaceType{Choice{Value: "400gbase-x-cdfp", Label: "CDFP (400GE)"}}
	GE400CFP8InterfaceType    = InterfaceType{Choice{Value: "400gbase-x-cfp8", Label: "CPF8 (400GE)"}}
	GE800QSFPDDInterfaceType  = InterfaceType{Choice{Value: "800gbase-x-qsfpdd", Label: "QSFP-DD (800GE)"}}
	GE800OSFPInterfaceType    = InterfaceType{Choice{Value: "800gbase-x-osfp", Label: "OSFP (800GE)"}}

	// Wireless
	IEEE80211AInterfaceType  = InterfaceType{Choice{Value: "ieee802.11a", Label: "IEEE 802.11a"}}
	IEEE80211GInterfaceType  = InterfaceType{Choice{Value: "ieee802.11g", Label: "IEEE 802.11b/g"}}
	IEEE80211NInterfaceType  = InterfaceType{Choice{Value: "ieee802.11n", Label: "IEEE 802.11n"}}
	IEEE80211ACInterfaceType = InterfaceType{Choice{Value: "ieee802.11ac", Label: "IEEE 802.11ac"}}
	IEEE80211ADInterfaceType = InterfaceType{Choice{Value: "ieee802.11ad", Label: "IEEE 802.11ad"}}
	IEEE80211AXInterfaceType = InterfaceType{Choice{Value: "ieee802.11ax", Label: "IEEE 802.11ax"}}

	// PON
	GPONInterfaceType     = InterfaceType{Choice{Value: "gpon", Label: "GPON (2.5 Gbps / 1.25 Gps)"}}
	XGPONInterfaceType    = InterfaceType{Choice{Value: "xg-pon", Label: "XG-PON (10 Gbps / 2.5 Gbps)"}}
	XGSPONInterfaceType   = InterfaceType{Choice{Value: "xgs-pon", Label: "XGS-PON (10 Gbps)"}}
	NGPON2InterfaceType   = InterfaceType{Choice{Value: "ng-pon2", Label: "NG-PON2 (TWDM-PON) (4x10 Gbps)"}}
	EPONInterfaceType     = InterfaceType{Choice{Value: "epon", Label: "EPON (1 Gbps)"}}
	TenGEPONInterfaceType = InterfaceType{Choice{Value: "10g-epon", Label: "10G-EPON (10 Gbps)"}}

	// Stacking

	// Cellular
	GSMInterfaceType  = InterfaceType{Choice{Value: "gsm", Label: "GSM"}}
	CDMAInterfaceType = InterfaceType{Choice{Value: "cdma", Label: "CDMA"}}
	LTEInterfaceType  = InterfaceType{Choice{Value: "lte", Label: "LTE"}}

	// Other type
	OtherInterfaceType = InterfaceType{Choice{Value: "other", Label: "Other"}}
)

// // Maps interface speed to interface type
// InterfaceTypeMap := map[int]InterfaceType{
// 	100: BASETXInterfaceType,
// 	1000: BASET1InterfaceType,
// }

// Interface speed in kbps
type InterfaceSpeed int64

// Available interface speeds
const (
	MBPS10  InterfaceSpeed = 10000
	MBPS100 InterfaceSpeed = 100000
	GBPS1   InterfaceSpeed = 1000000
	GBPS10  InterfaceSpeed = 10000000
	GBPS25  InterfaceSpeed = 25000000
	GBPS40  InterfaceSpeed = 40000000
	GBPS100 InterfaceSpeed = 100000000
	GBPS200 InterfaceSpeed = 200000000
	GBPS400 InterfaceSpeed = 400000000
)

// const speed2interfaceType = map[int]InterfaceType{
// 	100: BASEFX,
// }

// Interface represents a physical data interface within a device.
type Interface struct {
	NetboxObject
	// Device is the device to which the interface belongs. This field is required.
	Device *Device `json:"device,omitempty"`
	// Name is the name of the interface. This field is required.
	Name string `json:"name,omitempty"`
	// Status whether the interface is enabled or not.
	Status bool `json:"enabled,omitempty"`
	// Type is the type of interface. This field is required. Can only be one of the predetermined values.
	Type *InterfaceType `json:"type,omitempty"`
	// Interface speed in kbps
	Speed InterfaceSpeed `json:"speed,omitempty"`
	// Related Interfaces
	// Parent is the parent interface, if any.
	ParentInterface *Interface `json:"parent,omitempty"`
	// BridgedInterface is the bridged interface, if any.
	BridgedInterface *Interface `json:"bridge,omitempty"`
	// LAG is the LAG to which the interface belongs, if any.
	LAG *Interface `json:"lag,omitempty"`
	// MTU is the maximum transmission unit (MTU) configured for the interface.
	MTU int64 `json:"mtu,omitempty"`
	// TaggedVlans is a list of all the VLANs to which the interface is tagged.
	TaggedVlans []*Vlan `json:"tagged_vlans,omitempty"`
	// CustomFields that can be added to a device. We use source_id custom field to store the id of the interface in the source system.
	CustomFields map[string]string `json:"custom_fields,omitempty"`
}

func (i Interface) String() string {
	return fmt.Sprintf("Interface{Id: %d, Device: %s, Name: %s}", i.Id, i.Device.Name, i.Name)
}
