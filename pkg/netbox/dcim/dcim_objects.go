package dcim

import (
	"github.com/bl4ko/netbox-ssot/pkg/netbox/extras"
)

// Netbox's predetermined statuses, that we choose for some of
// our objects
type Status struct {
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
}

var (
	Active = Status{
		Value: "active",
		Label: "Active",
	}
	Offline = Status{
		Value: "offline",
		Label: "Offline",
	}
)

type Manafacturer struct {
	// ID is the unique numeric ID of the manufacturer.
	ID int `json:"id,omitempty"`
	// Name of the manufacturer (e.g. Cisco). This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`

	// Description of the manufacturer.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the manufacturer.
	Tags []*extras.Tag `json:"tags,omitempty"`
}

type Platform struct {
	// Netbox's ID of the platform.
	ID int `json:"id,omitempty"`
	// Name of the platform. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Manufacturer is the manufacturer of the platform.
	Manafacturer *Manafacturer `json:"manufacturer,omitempty"`
	// Description is a description of the platform.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the platform.
	Tags []*extras.Tag `json:"tags,omitempty"`
}

type Region struct {
	// Netbox's ID of the region.
	ID int `json:"id,omitempty"`
	// Name is the name of the region. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is a URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
}

// Site ares used for functional groupings.
// A site usually represents a building within a region.
type Site struct {
	ID int `json:"id,omitempty"`
	// Full name for the site. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Site status. This field is required.
	Status Status `json:"status,omitempty"`

	// Site description.
	Description string `json:"description,omitempty"`
	// Tags
	Tags []*extras.Tag `json:"tags,omitempty"`
}

// Location represents a physical location, such as a floor or room in a building.
type Location struct {
	// Netbox's ID of the location.
	ID int `json:"id,omitempty"`
	// Site is the site to which the location belongs. This field is required.
	Site *Site
	// Name is the name of the location. This field is required.
	Name string
	// URL-friendly unique shorthand. This field is required.
	Slug string
	// Status is the status of the location. This field is required.
	Status *Status

	// Location description.
	Description string
	// Tags
	Tags []*extras.Tag
}

type DeviceColor string

const (
	DeviceColorRed    DeviceColor = "Amber"
	DeviceColorOrange DeviceColor = "Orange"
	DeviceColorYellow DeviceColor = "Dark Orange"
	DeviceColorGreen  DeviceColor = "Brown"
	DeviceColorBlue   DeviceColor = "blue"
	DeviceColorPurple DeviceColor = "purple"
	DeviceColorBlack  DeviceColor = "black"
)

type DeviceType struct {
	// ID is the unique numeric ID of the device type.
	ID int `json:"id,omitempty"`
	// Manufacturer is the manufacturer of the device type. This field is required.
	Manafacturer *Manafacturer `json:"manufacturer,omitempty"`
	// Model is the model of the device type. This field is required.
	Model string `json:"model,omitempty"`
	// Slug is a URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Description is a description of the device type.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the device type.
	Tags []*extras.Tag `json:"tags,omitempty"`
	// UHeight is the height of the device type in rack units. This field is required.
	UHeight float32 `json:"u_height,omitempty"`
}

// DeviceRole represents the functional role of a device.
// For example, a device may play the role of a router, a switch, a firewall, etc.
type DeviceRole struct {
	// ID is the unique numeric ID of the device role.
	ID int `json:"id,omitempty"`
	// Name is the name of the device role. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Color of the device role. This field is required.
	Color *DeviceColor `json:"color,omitempty"`

	// Device role description.
	Description string `json:"description,omitempty"`
	// Tags
	Tags []*extras.Tag `json:"tags,omitempty"`
}

// Device can be any piece of physical hardware, such as a server, router, or switch.
type Device struct {
	// Netbox's ID of the device.
	ID int `json:"id,omitempty"`
	// Name is the name of the device.
	Name string `json:"name,omitempty"`
	// DeviceRole is the functional role of the device. This field is required.
	DeviceRole *DeviceRole `json:"role,omitempty"`

	// DeviceType is the type of device. This field is required.
	DeviceType *DeviceType `json:"device_type,omitempty"`
	// Site is the site to which the device belongs. This field is required.
	Site *Site `json:"site,omitempty"`

	// Description is a description of the device.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the device.
	Tags []*extras.Tag `json:"tags,omitempty"`
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
	// Netbox's ID of the interface.
	ID int `json:"id,omitempty"`
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

	// Tags is a list of tags for the interface.
	Tags []*extras.Tag `json:"tags,omitempty"`
	// Description is a description of the interface.
	Description string `json:"description,omitempty"`
}
