package service

// Here all mappings are defined so we don't hardcode api paths of objects
// in our code.
const (
	ContactGroupApiPath = "/api/tenancy/contact-groups/"
	ContactRoleApiPath  = "/api/tenancy/contact-roles/"
	ContactApiPath      = "/api/tenancy/contacts/"
	TenantApiPath       = "/api/tenancy/tenants/"

	VlanGroupApiPath = "/api/ipam/vlan-groups/"
	VlanApiPath      = "/api/ipam/vlans/"
	IpAddressApiPath = "/api/ipam/ip-addresses/"

	ClusterTypeApiPath    = "/api/virtualization/cluster-types/"
	ClusterGroupApiPath   = "/api/virtualization/cluster-groups/"
	ClusterApiPath        = "/api/virtualization/clusters/"
	VirtualMachineApiPath = "/api/virtualization/virtual-machines/"
	VMInterfaceApiPath    = "/api/virtualization/interfaces/"

	DeviceApiPath       = "/api/dcim/devices/"
	DeviceRoleApiPath   = "/api/dcim/device-roles/"
	DeviceTypeApiPath   = "/api/dcim/device-types/"
	InterfaceApiPath    = "/api/dcim/interfaces/"
	SiteApiPath         = "/api/dcim/sites/"
	ManufacturerApiPath = "/api/dcim/manufacturers/"
	PlatformApiPath     = "/api/dcim/platforms/"

	CustomFieldApiPath = "/api/extras/custom-fields/"
	TagApiPath         = "/api/extras/tags/"
)
