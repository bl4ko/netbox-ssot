package objects

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

type ClusterGroup struct {
	NetboxObject
	// Name is the name of the cluster group. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the cluster group name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Description is a description of the cluster group.
}

func (cg ClusterGroup) String() string {
	return fmt.Sprintf("ClusterGroup{Name: %s}", cg.Name)
}

// ClusterGroup implements IDItem interface.
func (cg *ClusterGroup) GetID() int {
	return cg.ID
}
func (cg *ClusterGroup) GetObjectType() constants.ContentType {
	return constants.ContentTypeVirtualizationClusterGroup
}
func (cg *ClusterGroup) GetAPIPath() constants.APIPath {
	return constants.ClusterGroupsAPIPath
}

// ClusterGroup implements OrphanItem interface.
func (cg *ClusterGroup) GetNetboxObject() *NetboxObject {
	return &cg.NetboxObject
}

type ClusterType struct {
	NetboxObject
	// Name is the name of the cluster type. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the cluster type name. This field is read-only.
	Slug string `json:"slug,omitempty"`
}

func (ct ClusterType) String() string {
	return fmt.Sprintf("ClusterType{Name: %s}", ct.Name)
}

// ClusterType implements IDItem interface.
func (ct *ClusterType) GetID() int {
	return ct.ID
}
func (ct *ClusterType) GetObjectType() constants.ContentType {
	return constants.ContentTypeVirtualizationClusterType
}
func (ct *ClusterType) GetAPIPath() constants.APIPath {
	return constants.ClusterTypesAPIPath
}

// ClusterType implements OrphanItem interface.
func (ct *ClusterType) GetNetboxObject() *NetboxObject {
	return &ct.NetboxObject
}

type ClusterStatus struct {
	Choice
}

var (
	ClusterStatusActive  = ClusterStatus{Choice{Value: "active", Label: "Active"}}
	ClusterStatusOffline = ClusterStatus{Choice{Value: "offline", Label: "Offline"}}
)

type Cluster struct {
	NetboxObject
	// Name is the name of the cluster. This field is required.
	Name string `json:"name,omitempty"`
	// Type is the type of the cluster. This field is required.
	// e.g. oVirt,VMware...
	Type *ClusterType `json:"type,omitempty"`
	// ClusterGroup is the cluster group to which this cluster belongs.
	Group *ClusterGroup `json:"group,omitempty"`
	// ScopeType is the scope of the cluster.
	ScopeType constants.ContentType `json:"scope_type,omitempty"`
	// ScopeID is the ID of the scope object.
	ScopeID int `json:"scope_id,omitempty"`
	// Status is the operational status of the cluster. This field is required.
	Status ClusterStatus `json:"status,omitempty"`
	// TenantGroup is the tenant group to which this cluster belongs.
	TenantGroup *TenantGroup `json:"tenant_group,omitempty"`
	// Tenant is the tenant to which this cluster belongs.
	Tenant *Tenant `json:"tenant,omitempty"`
}

func (c Cluster) String() string {
	return fmt.Sprintf("Cluster{Name: %s, Type: %s}", c.Name, c.Type)
}

// Cluster implements IDItem interface.
func (c *Cluster) GetID() int {
	return c.ID
}
func (c *Cluster) GetObjectType() constants.ContentType {
	return constants.ContentTypeVirtualizationCluster
}
func (c *Cluster) GetAPIPath() constants.APIPath {
	return constants.ClustersAPIPath
}

// Cluster implements OrphanItem interface.
func (c *Cluster) GetNetboxObject() *NetboxObject {
	return &c.NetboxObject
}

type VMStatus struct {
	Choice
}

var (
	VMStatusActive  = VMStatus{Choice{Value: "active", Label: "Active"}}
	VMStatusOffline = VMStatus{Choice{Value: "offline", Label: "Offline"}}
)

// VM represents a netbox's virtual machine.
type VM struct {
	NetboxObject
	// Name is the name of the virtual machine. This field is required.
	Name string `json:"name,omitempty"`
	// Status is the status of the virtual machine. This field is required.
	Status *VMStatus `json:"status,omitempty"`
	// Site is the site to which this virtual machine belongs.
	Site *Site `json:"site,omitempty"`
	// Cluster is the cluster to which this virtual machine belongs.
	Cluster *Cluster `json:"cluster,omitempty"`
	// Host is a specific host that this virtual machine is hosted on.
	Host *Device `json:"device,omitempty"`

	// TenantGroup is the datacenter that this virtual machine belongs to.
	TenantGroup *TenantGroup `json:"tenant_group,omitempty"`
	// Tenant is the tenant to which this virtual machine belongs.
	Tenant *Tenant `json:"tenant,omitempty"`

	// Platform is the platform of the virtual machine.
	Platform *Platform `json:"platform,omitempty"`
	// PrimaryIPv4 is the primary IPv4 address assigned to the virtual machine.
	PrimaryIPv4 *IPAddress `json:"primary_ip4,omitempty"`
	// PrimaryIPv6 is the primary IPv6 address assigned to the virtual machine.
	PrimaryIPv6 *IPAddress `json:"primary_ip6,omitempty"`

	// VCPUs is the number of virtual CPUs allocated to the virtual machine.
	VCPUs float32 `json:"vcpus,omitempty"`
	// Memory is the amount of memory allocated to the virtual machine in MB.
	Memory int `json:"memory,omitempty"`
	// Disk is the amount of disk space allocated to the virtual machine in GB.
	Disk int `json:"disk,omitempty"`
	// Role of the virtual machine.
	Role *DeviceRole `json:"role,omitempty"`

	// Additional Comments
	Comments string `json:"comments,omitempty"`
}

func (vm VM) String() string {
	return fmt.Sprintf("VM{Name: %s, Cluster: %s}", vm.Name, vm.Cluster)
}

// VM implements IDItem interface.
func (vm *VM) GetID() int {
	return vm.ID
}
func (vm *VM) GetObjectType() constants.ContentType {
	return constants.ContentTypeVirtualizationVirtualMachine
}
func (vm *VM) GetAPIPath() constants.APIPath {
	return constants.VirtualMachinesAPIPath
}

// VM implements IPAddressOwner interface.
func (vm *VM) GetPrimaryIPv4Address() *IPAddress {
	return vm.PrimaryIPv4
}
func (vm *VM) GetPrimaryIPv6Address() *IPAddress {
	return vm.PrimaryIPv6
}
func (vm *VM) SetPrimaryIPAddress(ip *IPAddress) {
	vm.PrimaryIPv4 = ip
}
func (vm *VM) SetPrimaryIPv6Address(ip *IPAddress) {
	vm.PrimaryIPv6 = ip
}

// VM implements OrphanItem interface.
func (vm *VM) GetNetboxObject() *NetboxObject {
	return &vm.NetboxObject
}

// 802.1Q VLAN Tagging Mode (Access, Tagged, Tagged All).
type VMInterfaceMode struct {
	Choice
}

var (
	VMInterfaceModeAccess    = VMInterfaceMode{Choice{Value: "access", Label: "Access"}}
	VMInterfaceModeTagged    = VMInterfaceMode{Choice{Value: "tagged", Label: "Tagged"}}
	VMInterfaceModeTaggedAll = VMInterfaceMode{Choice{Value: "tagged-all", Label: "Tagged All"}}
)

type VMInterface struct {
	NetboxObject
	// VM that this interface belongs to. This field is required.
	VM *VM `json:"virtual_machine,omitempty"`
	// Name is the name of the interface. This field is required.
	Name string `json:"name,omitempty"`
	// PrimaryMACAddress is the primary MAC address of the interface.
	PrimaryMACAddress *MACAddress `json:"primary_mac_address,omitempty"`
	// MTU of the interface.
	MTU int `json:"mtu,omitempty"`
	// Enabled is true if interface is enabled, false otherwise.
	Enabled bool `json:"enabled,omitempty"`
	// Related parent interface of this interface.
	ParentInterface *VMInterface `json:"parent,omitempty"`
	// Related bridged interface
	BridgedInterface *VMInterface `json:"bridge,omitempty"`
	// 802.1Q VLAN Tagging Mode
	Mode *VMInterfaceMode `json:"mode,omitempty"`
	// When Mode=VMInterfaceModeTagged: TaggedVlans is a list of all the VLANs that are tagged on the interface.
	TaggedVlans []*Vlan `json:"tagged_vlans,omitempty"`
	// When mode=VMInterfaceModeAccess: UntaggedVlan is the VLAN that is untagged on the interface.
	UntaggedVlan *Vlan `json:"untagged_vlan,omitempty"`
}

func (vmi VMInterface) String() string {
	return fmt.Sprintf("VMInterface{Name: %s, VM: %s}", vmi.Name, vmi.VM.Name)
}

// VMInterface implements IDItem interface.
func (vmi *VMInterface) GetID() int {
	return vmi.ID
}
func (vmi *VMInterface) GetObjectType() constants.ContentType {
	return constants.ContentTypeVirtualizationVMInterface
}
func (vmi *VMInterface) GetAPIPath() constants.APIPath {
	return constants.VMInterfacesAPIPath
}

// VMInterface also implements MACAddressOwner interface.
func (vmi *VMInterface) GetPrimaryMACAddress() *MACAddress {
	return vmi.PrimaryMACAddress
}
func (vmi *VMInterface) SetPrimaryMACAddress(mac *MACAddress) {
	vmi.PrimaryMACAddress = mac
}

// VMInterface implements OrphanItem interface.
func (vmi *VMInterface) GetNetboxObject() *NetboxObject {
	return &vmi.NetboxObject
}
