package mapper

import (
	"reflect"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
)

func reverseMap(m map[reflect.Type]constants.APIPath) map[constants.APIPath]reflect.Type {
	reversed := make(map[constants.APIPath]reflect.Type, len(m))
	for k, v := range m {
		reversed[v] = k
	}
	return reversed
}

var Type2Path = map[reflect.Type]constants.APIPath{
	reflect.TypeOf((*objects.VlanGroup)(nil)).Elem():            constants.VlanGroupsAPIPath,
	reflect.TypeOf((*objects.Vlan)(nil)).Elem():                 constants.VlansAPIPath,
	reflect.TypeOf((*objects.IPAddress)(nil)).Elem():            constants.IPAddressesAPIPath,
	reflect.TypeOf((*objects.ClusterType)(nil)).Elem():          constants.ClusterTypesAPIPath,
	reflect.TypeOf((*objects.ClusterGroup)(nil)).Elem():         constants.ClusterGroupsAPIPath,
	reflect.TypeOf((*objects.Cluster)(nil)).Elem():              constants.ClustersAPIPath,
	reflect.TypeOf((*objects.VM)(nil)).Elem():                   constants.VirtualMachinesAPIPath,
	reflect.TypeOf((*objects.VMInterface)(nil)).Elem():          constants.VMInterfacesAPIPath,
	reflect.TypeOf((*objects.Device)(nil)).Elem():               constants.DevicesAPIPath,
	reflect.TypeOf((*objects.MACAddress)(nil)).Elem():           constants.MACAddressesAPIPath,
	reflect.TypeOf((*objects.VirtualDeviceContext)(nil)).Elem(): constants.VirtualDeviceContextsAPIPath,
	reflect.TypeOf((*objects.DeviceRole)(nil)).Elem():           constants.DeviceRolesAPIPath,
	reflect.TypeOf((*objects.DeviceType)(nil)).Elem():           constants.DeviceTypesAPIPath,
	reflect.TypeOf((*objects.Interface)(nil)).Elem():            constants.InterfacesAPIPath,
	reflect.TypeOf((*objects.Site)(nil)).Elem():                 constants.SitesAPIPath,
	reflect.TypeOf((*objects.SiteGroup)(nil)).Elem():            constants.SiteGroupsAPIPath,
	reflect.TypeOf((*objects.Manufacturer)(nil)).Elem():         constants.ManufacturersAPIPath,
	reflect.TypeOf((*objects.Platform)(nil)).Elem():             constants.PlatformsAPIPath,
	reflect.TypeOf((*objects.Tenant)(nil)).Elem():               constants.TenantsAPIPath,
	reflect.TypeOf((*objects.ContactGroup)(nil)).Elem():         constants.ContactGroupsAPIPath,
	reflect.TypeOf((*objects.ContactRole)(nil)).Elem():          constants.ContactRolesAPIPath,
	reflect.TypeOf((*objects.Contact)(nil)).Elem():              constants.ContactsAPIPath,
	reflect.TypeOf((*objects.CustomField)(nil)).Elem():          constants.CustomFieldsAPIPath,
	reflect.TypeOf((*objects.Tag)(nil)).Elem():                  constants.TagsAPIPath,
	reflect.TypeOf((*objects.ContactAssignment)(nil)).Elem():    constants.ContactAssignmentsAPIPath,
	reflect.TypeOf((*objects.Prefix)(nil)).Elem():               constants.PrefixesAPIPath,
	reflect.TypeOf((*objects.WirelessLAN)(nil)).Elem():          constants.WirelessLANsAPIPath,
	reflect.TypeOf((*objects.WirelessLANGroup)(nil)).Elem():     constants.WirelessLANGroupsAPIPath,
	reflect.TypeOf((*objects.VirtualDisk)(nil)).Elem():          constants.VirtualDisksAPIPath,
	reflect.TypeOf((*objects.VRF)(nil)).Elem():                  constants.VRFsAPIPath,
}

var Path2Type = reverseMap(Type2Path)
