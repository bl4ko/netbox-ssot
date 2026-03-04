package objects

import (
	"fmt"

	"github.com/src-doo/netbox-ssot/internal/constants"
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
	// VRF
    VRF *VRF `json:"vrf,omitempty"`
	// AssignedObjectType is either a DeviceInterface or a VMInterface.
	AssignedObjectType constants.ContentType `json:"assigned_object_type,omitempty"`
	// ID of the assigned object (either an ID of DeviceInterface or an ID of VMInterface).
	AssignedObjectID int `json:"assigned_object_id,omitempty"`
}

func (ip IPAddress) String() string {
	return fmt.Sprintf(
		"IPAddress{ID: %d, Address: %s, Status: %s, DNSName: %s}",
		ip.ID,
		ip.Address,
		ip.Status,
		ip.DNSName,
	)
}

// IPAddress implements IDItem interface.
func (ip *IPAddress) GetID() int {
	return ip.ID
}
func (ip *IPAddress) GetObjectType() constants.ContentType {
	return constants.ContentTypeIpamIPAddress
}
func (ip *IPAddress) GetAPIPath() constants.APIPath {
	return constants.IPAddressesAPIPath
}

// IPAddress implements OrphanItem interface.
func (ip *IPAddress) GetNetboxObject() *NetboxObject {
	return &ip.NetboxObject
}

type VRF struct {
	NetboxObject
	// Name of the VRF. This field is required.
	Name string `json:"name,omitempty"`
	// Route distinguisher
	RD string `json:"rd,omitempty"`
}

func (v VRF) String() string {
	return fmt.Sprintf("VRF{ID: %d, Name: %s, RD: %s}", v.ID, v.Name, v.RD)
}

func (v *VRF) GetID() int                           { return v.ID }
func (v *VRF) GetObjectType() constants.ContentType { return constants.ContentTypeIpamVRF }
func (v *VRF) GetAPIPath() constants.APIPath        { return constants.VRFsAPIPath }
func (v *VRF) GetNetboxObject() *NetboxObject       { return &v.NetboxObject }


type VidRange [2]int

type VlanGroup struct {
	NetboxObject
	// Name of the VlanGroup. This field is required.
	Name string `json:"name,omitempty"`
	// Slug of the VlanGroup. This field is required.
	Slug string `json:"slug,omitempty"`
	// VidRanges is a list of VID ranges that this VlanGroup can use.
	VidRanges []VidRange `json:"vid_ranges,omitempty"`
	// ScopeType is the scope of the VlanGroup.
	ScopeType constants.ContentType `json:"scope_type,omitempty"`
	// ScopeID is the ID of the scope object.
	ScopeID int `json:"scope_id,omitempty"`
}

func (vg VlanGroup) String() string {
	return fmt.Sprintf("VlanGroup{Name: %s, VidRanges: %v}", vg.Name, vg.VidRanges)
}

// VlanGroup implements IDItem interface.
func (vg *VlanGroup) GetID() int {
	return vg.ID
}
func (vg *VlanGroup) GetObjectType() constants.ContentType {
	return constants.ContentTypeIpamVlanGroup
}
func (vg *VlanGroup) GetAPIPath() constants.APIPath {
	return constants.VlanGroupsAPIPath
}

// VlanGroup implements OrphanItem interface.
func (vg *VlanGroup) GetNetboxObject() *NetboxObject {
	return &vg.NetboxObject
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
	return fmt.Sprintf(
		"Vlan{ID: %d, Name: %s, Vid: %d, Group: %s}",
		v.ID, v.Name, v.Vid, v.Group.Name)
}

// Vlan implements IDItem interface.
func (v *Vlan) GetID() int {
	return v.ID
}
func (v *Vlan) GetObjectType() constants.ContentType {
	return constants.ContentTypeIpamVlan
}
func (v *Vlan) GetAPIPath() constants.APIPath {
	return constants.VlansAPIPath
}

// Vlan implements OrphanItem interface.
func (v *Vlan) GetNetboxObject() *NetboxObject {
	return &v.NetboxObject
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
	ScopeID   int                   `json:"scope_id,omitempty"`
	ScopeType constants.ContentType `json:"scope_type,omitempty"`

	// Vlan that this prefix belongs to.
	Vlan *Vlan `json:"vlan,omitempty"`

	// Tenant that this prefix belongs to.
	Tenant *Tenant `json:"tenant,omitempty"`
	
	// VRF
	VRF *VRF `json:"vrf,omitempty"`

	Comments string `json:"comments,omitempty"`
}

func (p Prefix) String() string {
	return fmt.Sprintf("Prefix{Prefix: %s}", p.Prefix)
}

// Prefix implements IDItem interface.
func (p *Prefix) GetID() int {
	return p.ID
}
func (p *Prefix) GetObjectType() constants.ContentType {
	return constants.ContentTypeIpamPrefix
}
func (p *Prefix) GetAPIPath() constants.APIPath {
	return constants.PrefixesAPIPath
}

// Prefix implements OrphanItem interface.
func (p *Prefix) GetNetboxObject() *NetboxObject {
	return &p.NetboxObject
}

