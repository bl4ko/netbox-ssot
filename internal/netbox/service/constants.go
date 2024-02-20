package service

// Here all mappings are defined so we don't hardcode api paths of objects
// in our code.
const (
	ContactGroupsAPIPath      = "/api/tenancy/contact-groups/"
	ContactRolesAPIPath       = "/api/tenancy/contact-roles/"
	ContactsAPIPath           = "/api/tenancy/contacts/"
	TenantsAPIPath            = "/api/tenancy/tenants/"
	ContactAssignmentsAPIPath = "/api/tenancy/contact-assignments/"

	PrefixesAPIPath    = "/api/ipam/prefixes/"
	VlanGroupsAPIPath  = "/api/ipam/vlan-groups/"
	VlansAPIPath       = "/api/ipam/vlans/"
	IPAddressesAPIPath = "/api/ipam/ip-addresses/"

	ClusterTypesAPIPath    = "/api/virtualization/cluster-types/"
	ClusterGroupsAPIPath   = "/api/virtualization/cluster-groups/"
	ClustersAPIPath        = "/api/virtualization/clusters/"
	VirtualMachinesAPIPath = "/api/virtualization/virtual-machines/"
	VMInterfacesAPIPath    = "/api/virtualization/interfaces/"

	DevicesAPIPath       = "/api/dcim/devices/"
	DeviceRolesAPIPath   = "/api/dcim/device-roles/"
	DeviceTypesAPIPath   = "/api/dcim/device-types/"
	InterfacesAPIPath    = "/api/dcim/interfaces/"
	SitesAPIPath         = "/api/dcim/sites/"
	ManufacturersAPIPath = "/api/dcim/manufacturers/"
	PlatformsAPIPath     = "/api/dcim/platforms/"

	CustomFieldsAPIPath = "/api/extras/custom-fields/"
	TagsAPIPath         = "/api/extras/tags/"
)
