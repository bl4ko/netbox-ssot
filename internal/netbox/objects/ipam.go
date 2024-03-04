package objects

import (
	"fmt"
)

type IPAddressStatus struct {
	Choice
}

var (
	IPAddressStatusActive   = IPAddressStatus{Choice{Value: "active", Label: "Active"}}
	IPAddressStatusReserved = IPAddressStatus{Choice{Value: "reserved", Label: "Reserved"}}
	IPAddressStatusDHCP     = IPAddressStatus{Choice{Value: "dhcp", Label: "DHCP"}}
	IPAddressStatusSLAAC    = IPAddressStatus{Choice{Value: "slaac", Label: "SLAAC"}}
)

type IPAddressRole struct {
	Choice
}

var (
	IPAddressRoleLoopback  = IPAddressRole{Choice{Value: "loopback", Label: "Loopback"}}
	IPAddressRoleSecondary = IPAddressRole{Choice{Value: "secondary", Label: "Secondary"}}
	IPAddressRoleAnycast   = IPAddressRole{Choice{Value: "anycast", Label: "Anycast"}}
	IPAddressRoleVIP       = IPAddressRole{Choice{Value: "vip", Label: "VIP"}}
	IPAddressRoleVRRP      = IPAddressRole{Choice{Value: "vrrp", Label: "VRRP"}}
	IPAddressRoleHSRP      = IPAddressRole{Choice{Value: "hsrp", Label: "HSRP"}}
	IPAddressRoleGLBP      = IPAddressRole{Choice{Value: "glbp", Label: "GLBP"}}
	IPAddressRoleCARP      = IPAddressRole{Choice{Value: "carp", Label: "CARP"}}
)

type AssignedObjectType string

const (
	AssignedObjectTypeVMInterface     = "virtualization.vminterface"
	AssignedObjectTypeDeviceInterface = "dcim.interface"
)

type IPAddress struct {
	NetboxObject
	// IPv4 or IPv6 address (with mask). This field is required.
	Address string `json:"address,omitempty"`
	// The status of this IP address.
	Status *IPAddressStatus `json:"status,omitempty"`
	// Role of the IP address.
	Role *IPAddressRole `json:"role,omitempty"`
	// Hostname or FQDN (not case-sensitive)
	DNSName string `json:"dns_name,omitempty"`
	// Tenancy
	Tenant *Tenant `json:"tenant,omitempty"`

	// AssignedObjectType is either a DeviceInterface or a VMInterface.
	AssignedObjectType AssignedObjectType `json:"assigned_object_type,omitempty"`
	// ID of the assigned object (either an ID of DeviceInterface or an ID of VMInterface).
	AssignedObjectID int `json:"assigned_object_id,omitempty"`
}

func (ip IPAddress) String() string {
	return fmt.Sprintf("IPAddress{Id: %d, Address: %s, Status: %s, DNSName: %s}", ip.ID, ip.Address, ip.Status, ip.DNSName)
}

const (
	// Default vlan group for all objects, that are not party of any other vlan group.
	DefaultVlanGroupName = "Default netbox-ssot vlan group"
)

type VlanGroup struct {
	NetboxObject
	// Name of the VlanGroup. This field is required.
	Name string `json:"name,omitempty"`
	// Slug of the VlanGroup. This field is required.
	Slug string `json:"slug,omitempty"`
	// MinVid is the minimal VID that can be assigned in this group. This field is required (default 1).
	MinVid int `json:"min_vid,omitempty"`
	// MaxVid is the maximal VID that can be assigned in this group. This field is required (default 4094).
	MaxVid int `json:"max_vid,omitempty"`
}

func (vg VlanGroup) String() string {
	return fmt.Sprintf("VlanGroup{Name: %s, MinVid: %d, MaxVid: %d}", vg.Name, vg.MinVid, vg.MaxVid)
}

type VlanStatus struct {
	Choice
}

var (
	VlanStatusActive     = VlanStatus{Choice{Value: "active", Label: "Active"}}
	VlanStatusReserved   = VlanStatus{Choice{Value: "reserved", Label: "Reserved"}}
	VlanStatusDeprecated = VlanStatus{Choice{Value: "deprecated", Label: "Deprecated"}}
)

type Vlan struct {
	NetboxObject
	// Name of the VLAN. This field is required.
	Name string `json:"name,omitempty"`
	// VID of the VLAN. This field is required.
	Vid int `json:"vid,omitempty"`
	// VlanGroup that this vlan belongs to.
	Group *VlanGroup `json:"group,omitempty"`
	// Status of the VLAN. This field is required. Default is "active".
	Status *VlanStatus `json:"status,omitempty"`
	// Tenant that this VLAN belongs to.
	Tenant *Tenant `json:"tenant,omitempty"`
	// Site that this VLAN belongs to.
	Site *Site `json:"site,omitempty"`
	// Comments about this Vlan.
	Comments string `json:"comments,omitempty"`
}

func (v Vlan) String() string {
	return fmt.Sprintf("Vlan{Id: %d, Name: %s, Vid: %d, Status: %s}", v.ID, v.Name, v.Vid, v.Status)
}

type IPRange struct {
	NetboxObject
}

type PrefixStatus struct {
	Choice
}

// https://github.com/netbox-community/netbox/blob/b408beaed52cb9fc5f7e197a7e00479af3714564/netbox/ipam/choices.py#L21
var (
	PrefixStatusContainer  = PrefixStatus{Choice{Value: "container", Label: "Container"}}
	PrefixStatusActive     = PrefixStatus{Choice{Value: "active", Label: "Active"}}
	PrefixStatusReserved   = PrefixStatus{Choice{Value: "reserved", Label: "Reserved"}}
	PrefixStatusDeprecated = PrefixStatus{Choice{Value: "deprecated", Label: "Deprecated"}}
)

type Prefix struct {
	NetboxObject
	// Prefix is a IPv4 or IPv6 network address (with mask). This field is required.
	Prefix string `json:"prefix,omitempty"`
	// Status of the prefix (default "active").
	Status *PrefixStatus `json:"status,omitempty"`

	// Site that this prefix belongs to.
	Site *Site `json:"site,omitempty"`
	// Vlan that this prefix belongs to.
	Vlan *Vlan `json:"vlan,omitempty"`

	// Tenant that this prefix belongs to.
	Tenant *Tenant `json:"tenant,omitempty"`

	Comments string `json:"comments,omitempty"`
}

func (p Prefix) String() string {
	return fmt.Sprintf("Prefix{Prefix: %s}", p.Prefix)
}
