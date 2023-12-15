package dcim

import (
	"github.com/bl4ko/netbox-ssot/pkg/netbox/common"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/tenancy"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/virtualization"
)

type Region struct {
	common.NetboxObject
	// Name is the name of the region. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is a URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
}

// Location represents a physical location, such as a floor or room in a building.
type Location struct {
	common.NetboxObject
	// Site is the site to which the location belongs. This field is required.
	Site *common.Site
	// Name is the name of the location. This field is required.
	Name string
	// URL-friendly unique shorthand. This field is required.
	Slug string
	// Status is the status of the location. This field is required.
	Status *common.SiteStatus
}

// DeviceType represents the physical and operational characteristics of a device.
// For example, a device type may represent a Cisco C2960 switch running IOS 15.2.
type DeviceType struct {
	common.NetboxObject
	// Manufacturer is the manufacturer of the device type. This field is required.
	Manafacturer *common.Manafacturer `json:"manufacturer,omitempty"`
	// Model is the model of the device type. This field is required.
	Model string `json:"model,omitempty"`
	// Slug is a URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// UHeight is the height of the device type in rack units. This field is required.
	UHeight float32 `json:"u_height,omitempty"`
}

// DeviceRole represents the functional role of a device.
// For example, a device may play the role of a router, a switch, a firewall, etc.
type DeviceRole struct {
	common.NetboxObject
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
	common.Choice
}

var (
	FrontToRear = DeviceAirFlowType{common.Choice{Value: "front-to-rear", Label: "Front to rear"}}
	RearToFront = DeviceAirFlowType{common.Choice{Value: "rear-to-front", Label: "Rear to front"}}
	LeftToRight = DeviceAirFlowType{common.Choice{Value: "left-to-right", Label: "Left to right"}}
	RightToLeft = DeviceAirFlowType{common.Choice{Value: "right-to-left", Label: "Right to left"}}
	SideToRear  = DeviceAirFlowType{common.Choice{Value: "side-to-rear", Label: "Side to rear"}}
	Passive     = DeviceAirFlowType{common.Choice{Value: "passive", Label: "Passive"}}
	Mixed       = DeviceAirFlowType{common.Choice{Value: "mixed", Label: "Mixed"}}
)

type DeviceStatus struct {
	common.Choice
}

var (
	DeviceStatusOffline         = DeviceStatus{common.Choice{Value: "offline", Label: "Offline"}}
	DeviceStatusActive          = DeviceStatus{common.Choice{Value: "active", Label: "Active"}}
	DeviceStatusPlanned         = DeviceStatus{common.Choice{Value: "planned", Label: "Planned"}}
	DeviceStatusStaged          = DeviceStatus{common.Choice{Value: "staged", Label: "Staged"}}
	DeviceStatusFailed          = DeviceStatus{common.Choice{Value: "failed", Label: "Failed"}}
	DeviceStatusInventory       = DeviceStatus{common.Choice{Value: "inventory", Label: "Inventory"}}
	DeviceStatusDecommissioning = DeviceStatus{common.Choice{Value: "decommissioning", Label: "Decommissioning"}}
)

// Device can be any piece of physical hardware, such as a server, router, or switch.
type Device struct {
	common.NetboxObject

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
	Site *common.Site `json:"site,omitempty"`
	// Location is the location of the device.
	Location *Location `json:"location,omitempty"`

	// Management
	// Status of the device (e.g. active, offline, planned, etc.). This field is required.
	Status *DeviceStatus `json:"status,omitempty"`
	// Platform of the device (e.g. Cisco IOS, Dell OS9, etc.).
	Platform *common.Platform `json:"platform,omitempty"`

	// Primary IPv4

	// Primary IPv6
	// Out-of-band IP

	// Virtualization
	// Cluster is the cluster to which the device belongs. (e.g. VMWare server belonging to a specific cluster).
	Cluster *virtualization.Cluster `json:"cluster,omitempty"`

	// Tenancy
	// Tenant group
	Tenant *tenancy.Tenant `json:"tenant,omitempty"`

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

type InterfaceType string

const (
	// VirtualInterfaces
	Virtual InterfaceType = "Virtual"
	Bridge  InterfaceType = "Bridge"
	LAG     InterfaceType = "LAG"
	// Ethernet (fixed)
	BASEFX  InterfaceType = "100BASE-FX (10/100ME FIBER)"
	BASELFX InterfaceType = "100BASE-LX (10/100ME FIBER)"
	BASETX  InterfaceType = "100BASE-SX (10/100ME FIBER)"
)

// const speed2interfaceType = map[int]InterfaceType{
// 	100: BASEFX,
// }

// Interface represents a physical data interface within a device.
type Interface struct {
	common.NetboxObject
	// Device is the device to which the interface belongs. This field is required.
	Device *Device `json:"device,omitempty"`
	// Name is the name of the interface. This field is required.
	Name string `json:"name,omitempty"`
	// InterfaceType is the type of interface. This field is required. Can only be one of the predetermined values.
	InterfaceType *InterfaceType `json:"type,omitempty"`

	// Related Intefaces
	// Parent is the parent interface, if any.
	ParentInterface *Interface `json:"parent,omitempty"`
	// BridgedInterface is the bridged interface, if any.
	BridgedInterface *Interface `json:"bridge,omitempty"`
	// LAG is the LAG to which the interface belongs, if any.
	LAG *Interface `json:"lag,omitempty"`
}
