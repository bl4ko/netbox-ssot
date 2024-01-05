package objects

import "fmt"

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

	// AssignedInterface
	// AssignedObjectType is either a DeviceInterface or a VMInterface
	AssignedObjectType AssignedObjectType `json:"assigned_object_type,omitempty"`
	// ID of the assigned object (either a DeviceInterface or a VMInterface)
	AssignedObjectId int `json:"assigned_object_id,omitempty"`
	// AssignedObject can be either a DeviceInterface or a VMInterface
	AssignedObject interface{} `json:"assigned_object,omitempty"`
}

func (ip IPAddress) String() string {
	return fmt.Sprintf("IPAddress{Id: %d, Address: %s, Status: %s, DNSName: %s}", ip.Id, ip.Address, ip.Status, ip.DNSName)
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
	// Status of the VLAN. This field is required. Default is "active".
	Status *VlanStatus `json:"status,omitempty"`
	// Tenant that this VLAN belongs to.
	Tenant *Tenant `json:"tenant,omitempty"`
	// Comments about this Vlan.
	Comments string `json:"comments,omitempty"`
}

func (v Vlan) String() string {
	return fmt.Sprintf("Vlan{Id: %d, Name: %s, Vid: %d, Status: %s}", v.Id, v.Name, v.Vid, v.Status)
}

type IpRange struct {
	NetboxObject
}
