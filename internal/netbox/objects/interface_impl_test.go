package objects

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func TestGetObjectType(t *testing.T) {
	tests := []struct {
		name string
		item IDItem
		want constants.ContentType
	}{
		{"Site", &Site{}, constants.ContentTypeDcimSite},
		{"SiteGroup", &SiteGroup{}, constants.ContentTypeDcimSiteGroup},
		{"Platform", &Platform{}, constants.ContentTypeDcimPlatform},
		{"Region", &Region{}, constants.ContentTypeDcimRegion},
		{"Location", &Location{}, constants.ContentTypeDcimLocation},
		{"Manufacturer", &Manufacturer{}, constants.ContentTypeDcimManufacturer},
		{"DeviceType", &DeviceType{}, constants.ContentTypeDcimDeviceType},
		{"DeviceRole", &DeviceRole{}, constants.ContentTypeDcimDeviceRole},
		{"Device", &Device{}, constants.ContentTypeDcimDevice},
		{"Interface", &Interface{}, constants.ContentTypeDcimInterface},
		{"VirtualDeviceContext", &VirtualDeviceContext{}, constants.ContentTypeDcimVirtualDeviceContext},
		{"MACAddress", &MACAddress{}, constants.ContentTypeDcimMACAddress},
		{"IPAddress", &IPAddress{}, constants.ContentTypeIpamIPAddress},
		{"VlanGroup", &VlanGroup{}, constants.ContentTypeIpamVlanGroup},
		{"Vlan", &Vlan{}, constants.ContentTypeIpamVlan},
		{"Prefix", &Prefix{}, constants.ContentTypeIpamPrefix},
		{"VRF", &VRF{}, constants.ContentTypeIpamVRF},
		{"TenantGroup", &TenantGroup{}, constants.ContentTypeTenancyTenantGroup},
		{"Tenant", &Tenant{}, constants.ContentTypeTenancyTenant},
		{"Contact", &Contact{}, constants.ContentTypeTenancyContact},
		{"ContactAssignment", &ContactAssignment{}, constants.ContentTypeTenancyContactAssignment},
		{"ClusterGroup", &ClusterGroup{}, constants.ContentTypeVirtualizationClusterGroup},
		{"ClusterType", &ClusterType{}, constants.ContentTypeVirtualizationClusterType},
		{"Cluster", &Cluster{}, constants.ContentTypeVirtualizationCluster},
		{"VM", &VM{}, constants.ContentTypeVirtualizationVirtualMachine},
		{"VMInterface", &VMInterface{}, constants.ContentTypeVirtualizationVMInterface},
		{"VirtualDisk", &VirtualDisk{}, constants.ContentTypeVirtualizationVirtualDisk},
		{"WirelessLANGroup", &WirelessLANGroup{}, constants.ContentTypeWirelessLANGroup},
		{"WirelessLAN", &WirelessLAN{}, constants.ContentTypeWirelessLAN},
		{"Tag", &Tag{}, constants.ContentTypeExtrasTag},
		{"CustomField", &CustomField{}, constants.ContentTypeExtrasCustomField},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.GetObjectType(); got != tt.want {
				t.Errorf("%s.GetObjectType() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestGetAPIPath(t *testing.T) {
	tests := []struct {
		name string
		item IDItem
		want constants.APIPath
	}{
		{"Site", &Site{}, constants.SitesAPIPath},
		{"SiteGroup", &SiteGroup{}, constants.SiteGroupsAPIPath},
		{"Platform", &Platform{}, constants.PlatformsAPIPath},
		{"Region", &Region{}, constants.RegionsAPIPath},
		{"Location", &Location{}, constants.LocationsAPIPath},
		{"Manufacturer", &Manufacturer{}, constants.ManufacturersAPIPath},
		{"DeviceType", &DeviceType{}, constants.DeviceTypesAPIPath},
		{"DeviceRole", &DeviceRole{}, constants.DeviceRolesAPIPath},
		{"Device", &Device{}, constants.DevicesAPIPath},
		{"Interface", &Interface{}, constants.InterfacesAPIPath},
		{"VirtualDeviceContext", &VirtualDeviceContext{}, constants.VirtualDeviceContextsAPIPath},
		{"MACAddress", &MACAddress{}, constants.MACAddressesAPIPath},
		{"IPAddress", &IPAddress{}, constants.IPAddressesAPIPath},
		{"VlanGroup", &VlanGroup{}, constants.VlanGroupsAPIPath},
		{"Vlan", &Vlan{}, constants.VlansAPIPath},
		{"Prefix", &Prefix{}, constants.PrefixesAPIPath},
		{"VRF", &VRF{}, constants.VRFsAPIPath},
		{"TenantGroup", &TenantGroup{}, constants.TenantGroupsAPIPath},
		{"Tenant", &Tenant{}, constants.TenantsAPIPath},
		{"Contact", &Contact{}, constants.ContactsAPIPath},
		{"ContactAssignment", &ContactAssignment{}, constants.ContactAssignmentsAPIPath},
		{"ClusterGroup", &ClusterGroup{}, constants.ClusterGroupsAPIPath},
		{"ClusterType", &ClusterType{}, constants.ClusterTypesAPIPath},
		{"Cluster", &Cluster{}, constants.ClustersAPIPath},
		{"VM", &VM{}, constants.VirtualMachinesAPIPath},
		{"VMInterface", &VMInterface{}, constants.VMInterfacesAPIPath},
		{"VirtualDisk", &VirtualDisk{}, constants.VirtualDisksAPIPath},
		{"WirelessLANGroup", &WirelessLANGroup{}, constants.WirelessLANGroupsAPIPath},
		{"WirelessLAN", &WirelessLAN{}, constants.WirelessLANsAPIPath},
		{"Tag", &Tag{}, constants.TagsAPIPath},
		{"CustomField", &CustomField{}, constants.CustomFieldsAPIPath},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.GetAPIPath(); got != tt.want {
				t.Errorf("%s.GetAPIPath() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestDevice_IPAddressOwner(t *testing.T) {
	ipv4 := &IPAddress{NetboxObject: NetboxObject{ID: 1}, Address: "10.0.0.1/24"}
	ipv6 := &IPAddress{NetboxObject: NetboxObject{ID: 2}, Address: "::1/128"}
	d := &Device{}

	if d.GetPrimaryIPv4Address() != nil {
		t.Error("expected nil primary IPv4 initially")
	}
	if d.GetPrimaryIPv6Address() != nil {
		t.Error("expected nil primary IPv6 initially")
	}

	d.SetPrimaryIPAddress(ipv4)
	d.SetPrimaryIPv6Address(ipv6)

	if got := d.GetPrimaryIPv4Address(); got != ipv4 {
		t.Errorf("Device.GetPrimaryIPv4Address() = %v, want %v", got, ipv4)
	}
	if got := d.GetPrimaryIPv6Address(); got != ipv6 {
		t.Errorf("Device.GetPrimaryIPv6Address() = %v, want %v", got, ipv6)
	}
}

func TestVM_IPAddressOwner(t *testing.T) {
	ipv4 := &IPAddress{NetboxObject: NetboxObject{ID: 1}, Address: "10.0.0.1/24"}
	ipv6 := &IPAddress{NetboxObject: NetboxObject{ID: 2}, Address: "::1/128"}
	vm := &VM{}

	if vm.GetPrimaryIPv4Address() != nil {
		t.Error("expected nil primary IPv4 initially")
	}
	if vm.GetPrimaryIPv6Address() != nil {
		t.Error("expected nil primary IPv6 initially")
	}

	vm.SetPrimaryIPAddress(ipv4)
	vm.SetPrimaryIPv6Address(ipv6)

	if got := vm.GetPrimaryIPv4Address(); got != ipv4 {
		t.Errorf("VM.GetPrimaryIPv4Address() = %v, want %v", got, ipv4)
	}
	if got := vm.GetPrimaryIPv6Address(); got != ipv6 {
		t.Errorf("VM.GetPrimaryIPv6Address() = %v, want %v", got, ipv6)
	}
}

func TestInterface_MACAddressOwner(t *testing.T) {
	mac := &MACAddress{NetboxObject: NetboxObject{ID: 1}, MAC: "00:11:22:33:44:55"}
	iface := &Interface{}

	if iface.GetPrimaryMACAddress() != nil {
		t.Error("expected nil primary MAC initially")
	}

	iface.SetPrimaryMACAddress(mac)

	if got := iface.GetPrimaryMACAddress(); got != mac {
		t.Errorf("Interface.GetPrimaryMACAddress() = %v, want %v", got, mac)
	}
}

func TestVMInterface_MACAddressOwner(t *testing.T) {
	mac := &MACAddress{NetboxObject: NetboxObject{ID: 1}, MAC: "00:11:22:33:44:55"}
	vmi := &VMInterface{}

	if vmi.GetPrimaryMACAddress() != nil {
		t.Error("expected nil primary MAC initially")
	}

	vmi.SetPrimaryMACAddress(mac)

	if got := vmi.GetPrimaryMACAddress(); got != mac {
		t.Errorf("VMInterface.GetPrimaryMACAddress() = %v, want %v", got, mac)
	}
}

func TestMACAddress_String(t *testing.T) {
	mac := MACAddress{
		NetboxObject:       NetboxObject{ID: 1},
		MAC:                "00:11:22:33:44:55",
		AssignedObjectType: constants.ContentTypeDcimInterface,
		AssignedObjectID:   42,
	}
	want := "MACAddress{MAC: 00:11:22:33:44:55, AssignedObjectType: dcim.interface, AssignedObjectID: 42}"
	if got := mac.String(); got != want {
		t.Errorf("MACAddress.String() = %v, want %v", got, want)
	}
}

func TestMACAddress_GetID(t *testing.T) {
	mac := &MACAddress{NetboxObject: NetboxObject{ID: 5}}
	if got := mac.GetID(); got != 5 {
		t.Errorf("MACAddress.GetID() = %v, want 5", got)
	}
}

func TestMACAddress_GetNetboxObject(t *testing.T) {
	nbo := NetboxObject{ID: 3}
	mac := &MACAddress{NetboxObject: nbo}
	if got := mac.GetNetboxObject(); got.ID != 3 {
		t.Errorf("MACAddress.GetNetboxObject().ID = %v, want 3", got.ID)
	}
}

func TestVirtualDisk_String(t *testing.T) {
	vd := VirtualDisk{
		Name: "disk-0",
		VM:   &VM{Name: "test-vm"},
	}
	want := "VirtualDisk{Name: disk-0, VM: test-vm}"
	if got := vd.String(); got != want {
		t.Errorf("VirtualDisk.String() = %v, want %v", got, want)
	}
}

func TestVirtualDisk_GetID(t *testing.T) {
	vd := &VirtualDisk{NetboxObject: NetboxObject{ID: 7}}
	if got := vd.GetID(); got != 7 {
		t.Errorf("VirtualDisk.GetID() = %v, want 7", got)
	}
}

func TestVirtualDisk_GetNetboxObject(t *testing.T) {
	nbo := NetboxObject{ID: 4}
	vd := &VirtualDisk{NetboxObject: nbo}
	if got := vd.GetNetboxObject(); got.ID != 4 {
		t.Errorf("VirtualDisk.GetNetboxObject().ID = %v, want 4", got.ID)
	}
}

func TestContactGroup_GetObjectType(t *testing.T) {
	cg := &ContactGroup{}
	if got := cg.GetObjectType(); got != constants.ContentTypeTenancyContactGroup {
		t.Errorf("ContactGroup.GetObjectType() = %v, want %v", got, constants.ContentTypeTenancyContactGroup)
	}
}

func TestContactRole_GetObjectType(t *testing.T) {
	cr := &ContactRole{}
	if got := cr.GetObjectType(); got != constants.ContentTypeTenancyContactRole {
		t.Errorf("ContactRole.GetObjectType() = %v, want %v", got, constants.ContentTypeTenancyContactRole)
	}
}

func TestSiteGroup_String(t *testing.T) {
	sg := SiteGroup{Name: "test-group"}
	want := "SiteGroup{Name: test-group}"
	if got := sg.String(); got != want {
		t.Errorf("SiteGroup.String() = %v, want %v", got, want)
	}
}

func TestSiteGroup_GetID(t *testing.T) {
	sg := &SiteGroup{NetboxObject: NetboxObject{ID: 2}}
	if got := sg.GetID(); got != 2 {
		t.Errorf("SiteGroup.GetID() = %v, want 2", got)
	}
}

func TestSiteGroup_GetNetboxObject(t *testing.T) {
	nbo := NetboxObject{ID: 6}
	sg := &SiteGroup{NetboxObject: nbo}
	if got := sg.GetNetboxObject(); got.ID != 6 {
		t.Errorf("SiteGroup.GetNetboxObject().ID = %v, want 6", got.ID)
	}
}

func TestVRF_String(t *testing.T) {
	v := VRF{NetboxObject: NetboxObject{ID: 1}, Name: "main", RD: "65000:1"}
	want := "VRF{ID: 1, Name: main, RD: 65000:1}"
	if got := v.String(); got != want {
		t.Errorf("VRF.String() = %v, want %v", got, want)
	}
}

func TestVRF_GetID(t *testing.T) {
	v := &VRF{NetboxObject: NetboxObject{ID: 9}}
	if got := v.GetID(); got != 9 {
		t.Errorf("VRF.GetID() = %v, want 9", got)
	}
}

func TestVRF_GetNetboxObject(t *testing.T) {
	nbo := NetboxObject{ID: 8}
	v := &VRF{NetboxObject: nbo}
	if got := v.GetNetboxObject(); got.ID != 8 {
		t.Errorf("VRF.GetNetboxObject().ID = %v, want 8", got.ID)
	}
}
