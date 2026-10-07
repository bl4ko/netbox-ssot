package inventory

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
)

func TestNetboxInventory_AddTag(t *testing.T) {
	type args struct {
		ctx    context.Context
		newTag *objects.Tag
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Tag
		wantErr bool
	}{
		{
			name: "Test add new tag",
			nbi:  MockInventory,
			args: args{
				ctx: context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newTag: &objects.Tag{
					Name:        "new tag",
					Description: "New Tag",
					Color:       constants.ColorBlack,
					Slug:        "new_tag",
				},
			},
			// Mock echoes request body with injected ID.
			want: &objects.Tag{
				ID:          1,
				Name:        "new tag",
				Description: "New Tag",
				Color:       constants.ColorBlack,
				Slug:        "new_tag",
			},
		},
		{
			name: "Test update existing tag",
			nbi:  MockInventory,
			args: args{
				ctx: context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newTag: &objects.Tag{
					Name:        "existing_tag1",
					Description: "New Tag",
					Color:       constants.ColorBlack,
					Slug:        "new_tag",
				},
			},
			want: &service.MockTagPatchResponse,
		},
		{
			name: "Test add the same tag",
			nbi:  MockInventory,
			args: args{
				ctx:    context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newTag: MockExistingTags["existing_tag2"],
			},
			want: MockExistingTags["existing_tag2"],
		},
	}

	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddTag(tt.args.ctx, tt.args.newTag)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Errorf("unmarshal test data: %s", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddTenant(t *testing.T) {
	type args struct {
		ctx       context.Context
		newTenant *objects.Tenant
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Tenant
		wantErr bool
	}{
		{
			name: "Test add new tenant",
			nbi:  MockInventory,
			args: args{
				ctx:       context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newTenant: &objects.Tenant{Name: "new tenant", Slug: "new_tenant"},
			},
			// Mock echoes request body with injected ID. AddTenant adds SsotTag before create.
			want: &objects.Tenant{
				NetboxObject: objects.NetboxObject{
					ID:   3,
					Tags: []*objects.Tag{MockInventory.SsotTag},
				},
				Name: "new tenant",
				Slug: "new_tenant",
			},
		},
		{
			name: "Test update existing tenant",
			nbi:  MockInventory,
			args: args{
				ctx:       context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newTenant: &objects.Tenant{Name: "existing_tenant1", Slug: "new_tenant"},
			},
			want: &service.MockTenantPatchResponse,
		},
		{
			name: "Test add the same tenant",
			nbi:  MockInventory,
			args: args{
				ctx:       context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newTenant: &objects.Tenant{Name: "existing_tenant2"},
			},
			want: MockExistingTenants["existing_tenant2"],
		},
	}
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddTenant(tt.args.ctx, tt.args.newTenant)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddSite(t *testing.T) {
	type args struct {
		ctx     context.Context
		newSite *objects.Site
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Site
		wantErr bool
	}{
		{
			name: "Test add new site",
			nbi:  MockInventory,
			args: args{
				ctx:     context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newSite: &objects.Site{Name: "new site", Slug: "new_site"},
			},
			// Mock echoes request body with injected ID. AddSite adds SsotTag before create.
			want: &objects.Site{
				NetboxObject: objects.NetboxObject{
					ID:   3,
					Tags: []*objects.Tag{MockInventory.SsotTag},
				},
				Name: "new site",
				Slug: "new_site",
			},
		},
		{
			name: "Test update existing site",
			nbi:  MockInventory,
			args: args{
				ctx:     context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newSite: &objects.Site{Name: "existing_site1", Slug: "new_site"},
			},
			want: &service.MockSitePatchResponse,
		},
		{
			name: "Test add the same site",
			nbi:  MockInventory,
			args: args{
				ctx:     context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newSite: &objects.Site{Name: "existing_site2"},
			},
			want: MockExistingSites["existing_site2"],
		},
	}
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddSite(tt.args.ctx, tt.args.newSite)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddSite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddSite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddContactRole(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.ContactRole
		wantErr bool
	}{
		{
			name:    "Create new contact role",
			args:    &objects.ContactRole{Name: "new_contact_role", Slug: "new_contact_role"},
			wantErr: false,
		},
		{
			name:    "Update existing contact role",
			args:    &objects.ContactRole{Name: "existing_contact_role1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddContactRole(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddContactRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddContactRole() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddContactGroup(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.ContactGroup
		wantErr bool
	}{
		{
			name:    "Create new contact group",
			args:    &objects.ContactGroup{Name: "new_contact_group", Slug: "new_contact_group"},
			wantErr: false,
		},
		{
			name:    "Update existing contact group",
			args:    &objects.ContactGroup{Name: "existing_contact_group1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddContactGroup(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddContactGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddContactGroup() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddContact(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.Contact
		wantErr bool
	}{
		{
			name:    "Create new contact",
			args:    &objects.Contact{Name: "new_contact"},
			wantErr: false,
		},
		{
			name:    "Update existing contact",
			args:    &objects.Contact{Name: "existing_contact1", Email: "new@email.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddContact(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddContact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddContact() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddContactAssignment(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.ContactAssignment
		wantErr bool
	}{
		{
			name: "Existing contact assignment triggers patch",
			args: &objects.ContactAssignment{
				ModelType: constants.ContentTypeDcimDevice,
				ObjectID:  1,
				Contact:   &objects.Contact{NetboxObject: objects.NetboxObject{ID: 1}, Name: "existing_contact1"},
				Role:      &objects.ContactRole{NetboxObject: objects.NetboxObject{ID: 1}, Name: "existing_contact_role1"},
			},
			wantErr: false,
		},
		{
			name: "New content type creates new assignment index path",
			args: &objects.ContactAssignment{
				ModelType: constants.ContentTypeVirtualizationVirtualMachine,
				ObjectID:  1,
				Contact:   &objects.Contact{NetboxObject: objects.NetboxObject{ID: 1}, Name: "existing_contact1"},
				Role:      &objects.ContactRole{NetboxObject: objects.NetboxObject{ID: 1}, Name: "existing_contact_role1"},
			},
			wantErr: true, // Echo-based mock can't unmarshal nested Contact back
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddContactAssignment(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NetboxInventory.AddContactAssignment() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NetboxInventory.AddContactAssignment() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddCustomField(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.CustomField
		want    *objects.CustomField
		wantErr bool
	}{
		{
			name: "Update existing custom field",
			args: &objects.CustomField{
				Name:        "existing_cf1",
				Type:        objects.CustomFieldTypeText,
				Description: "updated",
			},
			want: &service.MockCustomFieldPatchResponse,
		},
		{
			name: "Same custom field no-op",
			args: &objects.CustomField{
				Name: "existing_cf2",
				Type: objects.CustomFieldTypeText,
			},
			want: MockExistingCustomFields["existing_cf2"],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddCustomField(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddCustomField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddCustomField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddClusterGroup(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.ClusterGroup
		wantErr bool
	}{
		{
			name:    "Create new cluster group",
			args:    &objects.ClusterGroup{Name: "new_cluster_group", Slug: "new_cluster_group"},
			wantErr: false,
		},
		{
			name:    "Update existing cluster group",
			args:    &objects.ClusterGroup{Name: "existing_cluster_group1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddClusterGroup(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddClusterGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddClusterGroup() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddClusterType(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.ClusterType
		wantErr bool
	}{
		{
			name:    "Create new cluster type",
			args:    &objects.ClusterType{Name: "new_cluster_type", Slug: "new_cluster_type"},
			wantErr: false,
		},
		{
			name:    "Update existing cluster type",
			args:    &objects.ClusterType{Name: "existing_cluster_type1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddClusterType(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddClusterType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddClusterType() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddCluster(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.Cluster
		wantErr bool
	}{
		{
			name:    "Create new cluster",
			args:    &objects.Cluster{Name: "new_cluster"},
			wantErr: false,
		},
		{
			name:    "Update existing cluster",
			args:    &objects.Cluster{Name: "existing_cluster1", Status: objects.ClusterStatusOffline},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddCluster(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddCluster() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddDeviceRole(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.DeviceRole
		wantErr bool
	}{
		{
			name:    "Create new device role",
			args:    &objects.DeviceRole{Name: "new_device_role", Slug: "new_device_role", Color: "aa1409"},
			wantErr: false,
		},
		{
			name: "Update existing device role",
			args: &objects.DeviceRole{
				Name:  "existing_device_role1",
				Slug:  "updated_slug",
				Color: constants.Color(constants.DeviceRoleServerColor),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddDeviceRole(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddDeviceRole() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddManufacturer(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.Manufacturer
		wantErr bool
	}{
		{
			name:    "Create new manufacturer",
			args:    &objects.Manufacturer{Name: "new_manufacturer", Slug: "new_manufacturer"},
			wantErr: false,
		},
		{
			name:    "Update existing manufacturer",
			args:    &objects.Manufacturer{Name: "existing_manufacturer1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddManufacturer(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddManufacturer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddManufacturer() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddDeviceType(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.DeviceType
		wantErr bool
	}{
		{
			name:    "Create new device type",
			args:    &objects.DeviceType{Model: "new_device_type", Slug: "new_device_type"},
			wantErr: false,
		},
		{
			name:    "Update existing device type",
			args:    &objects.DeviceType{Model: "existing_device_type1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddDeviceType(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddDeviceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddDeviceType() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddPlatform(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.Platform
		wantErr bool
	}{
		{
			name:    "Create new platform",
			args:    &objects.Platform{Name: "new_platform", Slug: "new_platform"},
			wantErr: false,
		},
		{
			name:    "Update existing platform",
			args:    &objects.Platform{Name: "existing_platform1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddPlatform(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddPlatform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddPlatform() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddDevice(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.Device
		wantErr bool
	}{
		{
			name:    "Device without site returns error",
			args:    &objects.Device{Name: "no_site_device"},
			wantErr: true,
		},
		{
			name: "Existing device triggers diff",
			args: &objects.Device{
				Name:       "existing_device1",
				Site:       mockSite1,
				DeviceRole: mockDeviceRole1,
				DeviceType: mockDeviceType1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddDevice(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NetboxInventory.AddDevice() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddVlanGroup(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.VlanGroup
		wantErr bool
	}{
		{
			name:    "Create new vlan group",
			args:    &objects.VlanGroup{Name: "new_vlan_group", Slug: "new_vlan_group"},
			wantErr: false,
		},
		{
			name:    "Update existing vlan group",
			args:    &objects.VlanGroup{Name: "existing_vlan_group1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddVlanGroup(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVlanGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddVlanGroup() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddVlan(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.Vlan
		wantErr bool
	}{
		{
			name: "Existing vlan triggers diff",
			args: &objects.Vlan{
				Name:  "existing_vlan100",
				Vid:   100,
				Group: mockVlanGroup1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddVlan(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddVlan() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddInterface(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.Interface
		wantErr bool
	}{
		{
			name: "Existing interface triggers diff",
			args: &objects.Interface{
				Name:   "eth0",
				Device: mockDevice1,
				Type:   &objects.VirtualInterfaceType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddInterface(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddInterface() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddVM(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.VM
		wantErr bool
	}{
		{
			name:    "Create VM without cluster",
			args:    &objects.VM{Name: "new_vm_no_cluster"},
			wantErr: false,
		},
		{
			name:    "Existing VM triggers diff",
			args:    &objects.VM{Name: "existing_vm1", Cluster: mockCluster1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddVM(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddVM() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddVMInterface(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.VMInterface
		wantErr bool
	}{
		{
			name:    "Existing VM interface triggers diff",
			args:    &objects.VMInterface{Name: "vmeth0", VM: mockVM1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddVMInterface(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVMInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddVMInterface() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddIPAddress(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.IPAddress
		wantErr bool
	}{
		{
			name: "Create new IP address",
			args: &objects.IPAddress{
				Address:            "10.0.0.2/24",
				AssignedObjectType: constants.ContentTypeDcimInterface,
				AssignedObjectID:   1,
			},
			wantErr: false,
		},
		{
			name: "Existing IP address triggers diff",
			args: &objects.IPAddress{
				Address:            "10.0.0.1/24",
				AssignedObjectType: constants.ContentTypeDcimInterface,
				AssignedObjectID:   1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddIPAddress(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddIPAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddIPAddress() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddMACAddress(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.MACAddress
		wantErr bool
	}{
		{
			name: "Create new MAC address",
			args: &objects.MACAddress{
				MAC:                "11:22:33:44:55:66",
				AssignedObjectType: constants.ContentTypeDcimInterface,
				AssignedObjectID:   1,
			},
			wantErr: false,
		},
		{
			name: "Existing MAC address triggers diff",
			args: &objects.MACAddress{
				MAC:                "AA:BB:CC:DD:EE:FF",
				AssignedObjectType: constants.ContentTypeDcimInterface,
				AssignedObjectID:   1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddMACAddress(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddMACAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddMACAddress() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddPrefix(t *testing.T) {
	// Start mock NetBox server that validates custom_fields payloads
	// (rejects nested objects with "display" — mimics NetBox 4.2.x behavior)
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	type args struct {
		ctx       context.Context
		newPrefix *objects.Prefix
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		{
			name: "Create new prefix succeeds",
			nbi:  MockInventory,
			args: args{
				ctx: context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newPrefix: &objects.Prefix{
					Prefix: "192.168.1.0/24",
				},
			},
			wantErr: false,
		},
		{
			name: "Patch existing prefix with object-type custom fields succeeds (sanitized)",
			nbi:  MockInventory,
			args: args{
				ctx: context.WithValue(context.Background(), constants.CtxSourceKey, "new-source"),
				// This prefix already exists in MockExistingPrefixes with a site_ref
				// custom field that is a nested object (as returned by the NetBox API).
				// The source field differs ("new-source" vs "test") which triggers a PATCH.
				// Without sanitization, the PATCH payload would include the nested site_ref
				// object with "display", causing NetBox to return 400.
				newPrefix: &objects.Prefix{
					Prefix: "10.0.0.0/24",
				},
			},
			wantErr: false,
		},
		{
			name: "Same prefix no changes - no PATCH needed",
			nbi:  MockInventory,
			args: args{
				ctx: context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newPrefix: &objects.Prefix{
					Prefix: "10.0.0.0/24",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddPrefix(tt.args.ctx, tt.args.newPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddPrefix() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddVirtualDeviceContext(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.VirtualDeviceContext
		wantErr bool
	}{
		{
			name: "VDC without device returns error",
			args: &objects.VirtualDeviceContext{
				Name:   "no_device_vdc",
				Status: &objects.VDCStatusActive,
			},
			wantErr: true,
		},
		{
			name: "Existing VDC triggers diff",
			args: &objects.VirtualDeviceContext{
				Name:   "existing_vdc1",
				Device: mockDevice1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddVirtualDeviceContext(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NetboxInventory.AddVirtualDeviceContext() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NetboxInventory.AddVirtualDeviceContext() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddWirelessLAN(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.WirelessLAN
		wantErr bool
	}{
		{
			name:    "Create new wireless LAN",
			args:    &objects.WirelessLAN{SSID: "new_wlan"},
			wantErr: false,
		},
		{
			name:    "Update existing wireless LAN",
			args:    &objects.WirelessLAN{SSID: "existing_wlan1", AuthPsk: "new_psk"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddWirelessLAN(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddWirelessLAN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddWirelessLAN() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddWirelessLANGroup(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.WirelessLANGroup
		wantErr bool
	}{
		{
			name:    "Create new wireless LAN group",
			args:    &objects.WirelessLANGroup{Name: "new_wlan_group", Slug: "new_wlan_group"},
			wantErr: false,
		},
		{
			name:    "Update existing wireless LAN group",
			args:    &objects.WirelessLANGroup{Name: "existing_wlan_group1", Slug: "updated_slug"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddWirelessLANGroup(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NetboxInventory.AddWirelessLANGroup() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddWirelessLANGroup() returned nil")
			}
		})
	}
}

func TestNetboxInventory_AddVirtualDisk(t *testing.T) {
	mockServer := service.CreateMockServer()
	defer mockServer.Close()
	service.MockNetboxClient.BaseURL = mockServer.URL

	tests := []struct {
		name    string
		args    *objects.VirtualDisk
		wantErr bool
	}{
		{
			name: "Existing virtual disk triggers diff",
			args: &objects.VirtualDisk{
				Name: "existing_disk1",
				Size: 100,
				VM:   mockVM1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			got, err := MockInventory.AddVirtualDisk(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVirtualDisk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NetboxInventory.AddVirtualDisk() returned nil")
			}
		})
	}
}

func Test_addSourceNameCustomField(t *testing.T) {
	type args struct {
		ctx          context.Context
		netboxObject *objects.NetboxObject
	}
	tests := []struct {
		name string
		args args
		want *objects.NetboxObject
	}{
		{
			name: "Add source custom field to netbox object",
			args: args{
				ctx: context.WithValue(
					context.Background(),
					constants.CtxSourceKey,
					"testSource",
				),
				netboxObject: &objects.NetboxObject{},
			},
			want: &objects.NetboxObject{
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceName: "testSource",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addSourceNameCustomField(tt.args.ctx, tt.args.netboxObject)
			if !reflect.DeepEqual(tt.want, tt.args.netboxObject) {
				t.Errorf("%+v != %+v", tt.want, tt.args.netboxObject)
			}
		})
	}
}

func TestNetboxInventory_applyDeviceFieldLengthLimitations(t *testing.T) {
	type args struct {
		device *objects.Device
	}
	tests := []struct {
		name string
		nbi  *NetboxInventory
		args args
		want *objects.Device
	}{
		{
			name: "Test applyDeviceFieldLengthLimitations",
			nbi:  MockInventory,
			args: args{
				device: &objects.Device{
					Name:         "too_long_name_too_long_name_too_long_name_too_long_name_too_long_name",
					SerialNumber: "sjpqnnivlllbehccexqvlsxovizypvqdhyaaqptvaktnscbfjfownkdhwzckdhjzvpvkllaawxocwliaxhc",
					AssetTag:     "sjpqnnivlllbehccexqvlsxovizypvqdhyaaqptvaktnscbfjfownkdhwzckdhjzvpvkllaawxocwliaxhc",
				},
			},
			want: &objects.Device{
				Name:         "too_long_name_too_long_name_too_long_name_too_long_name_too_long",
				SerialNumber: "sjpqnnivlllbehccexqvlsxovizypvqdhyaaqptvaktnscbfjf",
				AssetTag:     "sjpqnnivlllbehccexqvlsxovizypvqdhyaaqptvaktnscbfjf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.nbi.applyDeviceFieldLengthLimitations(tt.args.device)
			if !reflect.DeepEqual(tt.want, tt.args.device) {
				t.Errorf("%+v != %+v", tt.want, tt.args)
			}
		})
	}
}

func TestAddPrefixDoesNotAdoptManualPrefix(t *testing.T) {
	manualPrefix := &objects.Prefix{
		NetboxObject: objects.NetboxObject{
			ID:          12,
			Description: "office LAN",
			Tags:        []*objects.Tag{{ID: 40, Name: "manual"}},
		},
		Prefix: "192.0.2.0/24",
	}
	nbi := &NetboxInventory{
		Logger:                mockLogger,
		OrphanManager:         NewOrphanManager(mockLogger),
		SsotTag:               &objects.Tag{ID: 1, Name: constants.SsotTagName},
		SourcePriority:        map[string]int{},
		NetboxAPI:             service.FailingMockNetboxClient,
		prefixesIndexByPrefix: map[string]map[int]*objects.Prefix{"192.0.2.0/24": {0: manualPrefix}},
	}
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "proxmox-a")

	got, err := nbi.AddPrefix(ctx, &objects.Prefix{Prefix: "192.0.2.0/24"})
	if err != nil {
		t.Fatalf("AddPrefix() error = %v, want the manual prefix returned without any API call", err)
	}
	if got != manualPrefix {
		t.Errorf("AddPrefix() = %v, want the existing manual prefix", got)
	}
	if manualPrefix.HasTagByName(constants.SsotTagName) {
		t.Errorf("manual prefix was tagged %s", constants.SsotTagName)
	}
}

// newDuplicateIPServer fakes a NetBox that refuses to create a duplicate IP address
// and lists existingIPsJSON for it. It records the IDs of the patched IP addresses.
func newDuplicateIPServer(t *testing.T, existingIPsJSON string) (*service.NetboxClient, func() []string) {
	t.Helper()
	var mu sync.Mutex
	patched := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"address":["Duplicate IP address found in global table: 10.0.0.5/24"]}`)
		case http.MethodGet:
			_, _ = io.WriteString(w, fmt.Sprintf(`{"count":1,"next":null,"previous":null,"results":%s}`, existingIPsJSON))
		case http.MethodPatch:
			mu.Lock()
			patched = append(patched, strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/ipam/ip-addresses/"), "/"))
			mu.Unlock()
			_, _ = io.WriteString(w, `{"id":1,"address":"10.0.0.5/24"}`)
		}
	}))
	t.Cleanup(server.Close)
	client := &service.NetboxClient{
		HTTPClient: &http.Client{},
		Logger:     mockLogger,
		BaseURL:    server.URL,
		APIToken:   "testtoken",
		Timeout:    constants.DefaultAPITimeout,
	}
	return client, func() []string {
		mu.Lock()
		defer mu.Unlock()
		out := make([]string, len(patched))
		copy(out, patched)
		return out
	}
}

func TestAddIPAddressDuplicateOnlyReassignsUnassignedIPOfSameVRF(t *testing.T) {
	tests := []struct {
		name        string
		existingIPs string
		wantErr     bool
		wantPatched []string
	}{
		{
			name:        "IP assigned to another object is not taken over",
			existingIPs: `[{"id":100,"address":"10.0.0.5/24","assigned_object_type":"dcim.interface","assigned_object_id":7}]`,
			wantErr:     true,
			wantPatched: []string{},
		},
		{
			name:        "unassigned IP of another VRF is not taken over",
			existingIPs: `[{"id":101,"address":"10.0.0.5/24","vrf":{"id":5,"name":"mgmt"}}]`,
			wantErr:     true,
			wantPatched: []string{},
		},
		{
			name: "unassigned IP of the same VRF is reassigned",
			existingIPs: `[{"id":102,"address":"10.0.0.5/24","assigned_object_type":"dcim.interface","assigned_object_id":7},` +
				`{"id":103,"address":"10.0.0.5/24"}]`,
			wantErr:     false,
			wantPatched: []string{"103"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, patched := newDuplicateIPServer(t, tt.existingIPs)
			nbi := &NetboxInventory{
				Logger:           mockLogger,
				OrphanManager:    NewOrphanManager(mockLogger),
				SsotTag:          &objects.Tag{ID: 1, Name: constants.SsotTagName},
				SourcePriority:   map[string]int{},
				NetboxAPI:        client,
				ipAddressesIndex: map[constants.ContentType]map[string]map[string]map[string]*objects.IPAddress{},
			}
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "proxmox-a")
			_, err := nbi.AddIPAddress(ctx, &objects.IPAddress{
				Address:            "10.0.0.5/24",
				AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
				AssignedObjectID:   60,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("AddIPAddress() error = %v, wantErr %t", err, tt.wantErr)
			}
			if got := patched(); !reflect.DeepEqual(got, tt.wantPatched) {
				t.Errorf("patched IP addresses = %v, want %v", got, tt.wantPatched)
			}
		})
	}
}

func TestAddTagIfMissingLeavesExistingTagUntouched(t *testing.T) {
	existingTag := &objects.Tag{ID: 7, Name: "prod", Slug: "prod", Color: "f44336"}
	nbi := &NetboxInventory{
		Logger:          mockLogger,
		SourcePriority:  map[string]int{},
		NetboxAPI:       service.FailingMockNetboxClient,
		tagsIndexByName: map[string]*objects.Tag{"prod": existingTag},
	}
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "proxmox-a")

	got, err := nbi.AddTagIfMissing(ctx, &objects.Tag{Name: "prod", Slug: "prod", Color: constants.ColorGreen})
	if err != nil {
		t.Fatalf("AddTagIfMissing() error = %v, want the existing tag without any API call", err)
	}
	if got != existingTag || existingTag.Color != "f44336" {
		t.Errorf("AddTagIfMissing() = %v, want the existing red tag untouched", got)
	}
}

func TestAddClusterIsScopedBySource(t *testing.T) {
	ssotTag := &objects.Tag{ID: 1, Name: constants.SsotTagName}
	unownedCluster := &objects.Cluster{NetboxObject: objects.NetboxObject{ID: 10}, Name: "legacy-cluster"}
	otherSourceCluster := &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			ID:           11,
			CustomFields: map[string]interface{}{constants.CustomFieldSourceName: "proxmox-b"},
		},
		Name: "pve-cluster0",
	}
	nbi := &NetboxInventory{
		Logger:         mockLogger,
		OrphanManager:  NewOrphanManager(mockLogger),
		SsotTag:        ssotTag,
		SourcePriority: map[string]int{},
		NetboxAPI:      &service.NetboxClient{Logger: mockLogger, DryRun: true},
		clustersIndexByNameAndSource: map[string]map[string]*objects.Cluster{
			"legacy-cluster": {"": unownedCluster},
			"pve-cluster0":   {"proxmox-b": otherSourceCluster},
		},
	}
	ctxA := context.WithValue(context.Background(), constants.CtxSourceKey, "proxmox-a")
	ctxB := context.WithValue(context.Background(), constants.CtxSourceKey, "proxmox-b")

	clusterA, err := nbi.AddCluster(ctxA, &objects.Cluster{Name: "pve-cluster0"})
	if err != nil {
		t.Fatalf("AddCluster(proxmox-a) error = %v", err)
	}
	if clusterA.ID == otherSourceCluster.ID {
		t.Errorf("proxmox-a got the cluster of proxmox-b (ID %d), want its own", clusterA.ID)
	}
	if got, ok := nbi.GetCluster(ctxB, "pve-cluster0"); !ok || got.ID != otherSourceCluster.ID {
		t.Errorf("GetCluster(proxmox-b) = %v, want cluster %d", got, otherSourceCluster.ID)
	}
	if got, ok := nbi.GetCluster(ctxA, "pve-cluster0"); !ok || got.ID != clusterA.ID {
		t.Errorf("GetCluster(proxmox-a) = %v, want cluster %d", got, clusterA.ID)
	}

	adopted, err := nbi.AddCluster(ctxA, &objects.Cluster{Name: "legacy-cluster"})
	if err != nil {
		t.Fatalf("AddCluster(legacy-cluster) error = %v", err)
	}
	if adopted.ID != unownedCluster.ID {
		t.Errorf("AddCluster(legacy-cluster) ID = %d, want the existing unowned cluster %d", adopted.ID, unownedCluster.ID)
	}
}

func TestAddIPAddressDuplicateReturnsLookupAndPatchErrors(t *testing.T) {
	permissionDenied := `{"detail":"You do not have permission to perform this action."}`
	tests := []struct {
		name      string
		failGet   bool
		failPatch bool
	}{
		{name: "lookup of the existing IP fails", failGet: true},
		{name: "reassignment of the existing IP fails", failPatch: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddIPAddressDuplicateError(t, permissionDenied, tt.failGet, tt.failPatch)
		})
	}
}

func testAddIPAddressDuplicateError(t *testing.T, failure string, failGet, failPatch bool) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"address":["Duplicate IP address found in global table: 10.0.0.5/24"]}`)
		case r.Method == http.MethodGet && failGet, r.Method == http.MethodPatch && failPatch:
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, failure)
		case r.Method == http.MethodGet:
			_, _ = io.WriteString(w, `{"count":1,"next":null,"previous":null,"results":[{"id":103,"address":"10.0.0.5/24"}]}`)
		case r.Method == http.MethodPatch:
			_, _ = io.WriteString(w, `{"id":103,"address":"10.0.0.5/24"}`)
		}
	}))
	defer server.Close()
	nbi := &NetboxInventory{
		Logger:         mockLogger,
		OrphanManager:  NewOrphanManager(mockLogger),
		SsotTag:        &objects.Tag{ID: 1, Name: constants.SsotTagName},
		SourcePriority: map[string]int{},
		NetboxAPI: &service.NetboxClient{
			HTTPClient: &http.Client{},
			Logger:     mockLogger,
			BaseURL:    server.URL,
			APIToken:   "testtoken",
			Timeout:    constants.DefaultAPITimeout,
		},
		ipAddressesIndex: map[constants.ContentType]map[string]map[string]map[string]*objects.IPAddress{},
	}
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "proxmox-a")

	_, err := nbi.AddIPAddress(ctx, &objects.IPAddress{Address: "10.0.0.5/24"})
	if err == nil {
		t.Fatalf("AddIPAddress() = nil error, want the NetBox failure")
	}
	if !strings.Contains(err.Error(), "permission") {
		t.Errorf("AddIPAddress() error = %q, want the NetBox failure (permission denied)", err)
	}
}
