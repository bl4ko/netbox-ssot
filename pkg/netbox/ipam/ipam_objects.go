package ipam

import (
	"github.com/bl4ko/netbox-ssot/pkg/netbox/common"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/tenancy"
)

type IPAddressStatus struct {
	common.Choice
}

var (
	IPAddressStatusActive   = IPAddressStatus{common.Choice{Value: "active", Label: "Active"}}
	IPAddressStatusReserved = IPAddressStatus{common.Choice{Value: "reserved", Label: "Reserved"}}
	IPAddressStatusDHCP     = IPAddressStatus{common.Choice{Value: "dhcp", Label: "DHCP"}}
	IPAddressStatusSLAAC    = IPAddressStatus{common.Choice{Value: "slaac", Label: "SLAAC"}}
)

type IPAddressRole struct {
	common.Choice
}

var (
	IPAddressRoleLoopback  = IPAddressRole{common.Choice{Value: "loopback", Label: "Loopback"}}
	IPAddressRoleSecondary = IPAddressRole{common.Choice{Value: "secondary", Label: "Secondary"}}
	IPAddressRoleAnycast   = IPAddressRole{common.Choice{Value: "anycast", Label: "Anycast"}}
	IPAddressRoleVIP       = IPAddressRole{common.Choice{Value: "vip", Label: "VIP"}}
	IPAddressRoleVRRP      = IPAddressRole{common.Choice{Value: "vrrp", Label: "VRRP"}}
	IPAddressRoleHSRP      = IPAddressRole{common.Choice{Value: "hsrp", Label: "HSRP"}}
	IPAddressRoleGLBP      = IPAddressRole{common.Choice{Value: "glbp", Label: "GLBP"}}
	IPAddressRoleCARP      = IPAddressRole{common.Choice{Value: "carp", Label: "CARP"}}
)

type AssignedObjectType string

const (
	AssignedObjectTypeVMInterface     = "virtualization.vminterface"
	AssignedObjectTypeDeviceInterface = "dcim.interface"
)

type IPAddress struct {
	common.NetboxObject
	// IPv4 or IPv6 address (with mask). This field is required.
	Address string `json:"address,omitempty"`
	// The status of this IP address.
	Status IPAddressStatus `json:"status,omitempty"`
	// Role of the IP address.
	Role IPAddressRole `json:"role,omitempty"`
	// Hostanme or FQDN (not case-sensitive)
	DNSName string `json:"dns_name,omitempty"`

	// Tenancy
	Tenant *tenancy.Tenant `json:"tenant,omitempty"`

	// AssignedInterface
	AssignedObjectType AssignedObjectType `json:"assigned_interface,omitempty"`
	// ID of the assigned object (either a DeviceInterface or a VMInterface)
	AssignedObjectID int `json:"assigned_object_id,omitempty"`
	// AssignedObject can be either a DeviceInterface or a VMInterface
	AssignedObject interface{} `json:"assigned_object,omitempty"`
}

type Vlan struct {
	common.NetboxObject
}

type IpRange struct {
	common.NetboxObject
}
