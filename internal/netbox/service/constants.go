package service

// Here all mappings are defined so we don't hardcode api paths of objects
// in our code.
const (
	ContactGroupsApiPath      = "/api/tenancy/contact-groups/"
	ContactRolesApiPath       = "/api/tenancy/contact-roles/"
	ContactsApiPath           = "/api/tenancy/contacts/"
	TenantsApiPath            = "/api/tenancy/tenants/"
	ContactAssignmentsApiPath = "/api/tenancy/contact-assignments/"

	PrefixesApiPath    = "/api/ipam/prefixes/"
	VlanGroupsApiPath  = "/api/ipam/vlan-groups/"
	VlansApiPath       = "/api/ipam/vlans/"
	IpAddressesApiPath = "/api/ipam/ip-addresses/"

	ClusterTypesApiPath    = "/api/virtualization/cluster-types/"
	ClusterGroupsApiPath   = "/api/virtualization/cluster-groups/"
	ClustersApiPath        = "/api/virtualization/clusters/"
	VirtualMachinesApiPath = "/api/virtualization/virtual-machines/"
	VMInterfacesApiPath    = "/api/virtualization/interfaces/"

	DevicesApiPath       = "/api/dcim/devices/"
	DeviceRolesApiPath   = "/api/dcim/device-roles/"
	DeviceTypesApiPath   = "/api/dcim/device-types/"
	InterfacesApiPath    = "/api/dcim/interfaces/"
	SitesApiPath         = "/api/dcim/sites/"
	ManufacturersApiPath = "/api/dcim/manufacturers/"
	PlatformsApiPath     = "/api/dcim/platforms/"

	CustomFieldsApiPath = "/api/extras/custom-fields/"
	TagsApiPath         = "/api/extras/tags/"
)
