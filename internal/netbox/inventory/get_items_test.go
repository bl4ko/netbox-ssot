package inventory

import (
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestNetboxInventory_GetTag(t *testing.T) {
	type args struct {
		tagName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.Tag
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetTag(tt.args.tagName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetTag() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetTag() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetManufacturer(t *testing.T) {
	type args struct {
		manufacturerName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.Manufacturer
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetManufacturer(tt.args.manufacturerName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetManufacturer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetManufacturer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetCustomField(t *testing.T) {
	type args struct {
		customFieldName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.CustomField
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetCustomField(tt.args.customFieldName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetCustomField() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetCustomField() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetVlan(t *testing.T) {
	type args struct {
		groupID int
		vlanID  int
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.Vlan
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetVlan(tt.args.groupID, tt.args.vlanID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetVlan() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetVlan() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetTenant(t *testing.T) {
	type args struct {
		tenantName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.Tenant
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetTenant(tt.args.tenantName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetTenant() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetTenant() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetSite(t *testing.T) {
	type args struct {
		siteName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.Site
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetSite(tt.args.siteName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetSite() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetSite() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetVlanGroup(t *testing.T) {
	type args struct {
		vlanGroupName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.VlanGroup
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetVlanGroup(tt.args.vlanGroupName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetVlanGroup() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetVlanGroup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetClusterGroup(t *testing.T) {
	type args struct {
		clusterGroupName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.ClusterGroup
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetClusterGroup(tt.args.clusterGroupName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetClusterGroup() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetClusterGroup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetCluster(t *testing.T) {
	type args struct {
		clusterName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.Cluster
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetCluster(tt.args.clusterName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetCluster() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetCluster() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetDevice(t *testing.T) {
	type args struct {
		deviceName string
		siteID     int
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.Device
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetDevice(tt.args.deviceName, tt.args.siteID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetDevice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetDevice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetDeviceRole(t *testing.T) {
	type args struct {
		deviceRoleName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.DeviceRole
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetDeviceRole(tt.args.deviceRoleName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetDeviceRole() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetDeviceRole() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetContactRole(t *testing.T) {
	type args struct {
		contactRoleName string
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.ContactRole
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetContactRole(tt.args.contactRoleName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetContactRole() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetContactRole() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetVirtualDeviceContext(t *testing.T) {
	type args struct {
		zoneName string
		deviceID int
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.VirtualDeviceContext
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetVirtualDeviceContext(tt.args.zoneName, tt.args.deviceID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"NetboxInventory.GetVirtualDeviceContext() got = %v, want %v",
					got,
					tt.want,
				)
			}
			if got1 != tt.want1 {
				t.Errorf(
					"NetboxInventory.GetVirtualDeviceContext() got1 = %v, want %v",
					got1,
					tt.want1,
				)
			}
		})
	}
}

func TestNetboxInventory_GetInterface(t *testing.T) {
	type args struct {
		interfaceName string
		deviceID      int
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.Interface
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetInterface(tt.args.interfaceName, tt.args.deviceID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetInterface() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NetboxInventory.GetInterface() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetboxInventory_GetContactAssignment(t *testing.T) {
	type args struct {
		contentType constants.ContentType
		objectID    int
		contactID   int
		roleID      int
	}
	tests := []struct {
		name  string
		nbi   *NetboxInventory
		args  args
		want  *objects.ContactAssignment
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.nbi.GetContactAssignment(
				tt.args.contentType,
				tt.args.objectID,
				tt.args.contactID,
				tt.args.roleID,
			)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.GetContactAssignment() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf(
					"NetboxInventory.GetContactAssignment() got1 = %v, want %v",
					got1,
					tt.want1,
				)
			}
		})
	}
}
