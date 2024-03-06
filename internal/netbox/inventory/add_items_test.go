package inventory

import (
	"context"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

var MockInventory = &NetboxInventory{
	Logger:          &logger.Logger{Logger: log.New(os.Stdout, "", log.LstdFlags)},
	TagsIndexByName: make(map[string]*objects.Tag),
	TagsLock:        sync.Mutex{},
	NetboxAPI:       nil, // TODO
}

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
		// {
		// 	name: "Test new tag add",
		// 	nbi:  MockInventory,
		// 	args: args{ctx: context.WithValue(context.Background(), constants.CtxSourceKey, "test"), newTag: &objects.Tag{Name: "new tag", Description: "New Tag", Color: objects.ColorBlack, Slug: "new_tag"}},
		// 	want: &objects.Tag{Name: "new tag", Description: "New Tag", Color: objects.ColorBlack, Slug: "new_tag"},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddTag(tt.args.ctx, tt.args.newTag)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddTag() error = %v, wantErr %v", err, tt.wantErr)
				return
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
		// TODO: Add test cases.
	}
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
		// TODO: Add test cases.
	}
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
	type args struct {
		ctx            context.Context
		newContactRole *objects.ContactRole
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.ContactRole
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddContactRole(tt.args.ctx, tt.args.newContactRole)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddContactRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddContactRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddContactGroup(t *testing.T) {
	type args struct {
		ctx             context.Context
		newContactGroup *objects.ContactGroup
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.ContactGroup
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddContactGroup(tt.args.ctx, tt.args.newContactGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddContactGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddContactGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddContact(t *testing.T) {
	type args struct {
		ctx        context.Context
		newContact *objects.Contact
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Contact
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddContact(tt.args.ctx, tt.args.newContact)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddContact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddContact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddContactAssignment(t *testing.T) {
	type args struct {
		ctx   context.Context
		newCA *objects.ContactAssignment
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.ContactAssignment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddContactAssignment(tt.args.ctx, tt.args.newCA)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddContactAssignment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddContactAssignment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddCustomField(t *testing.T) {
	type args struct {
		ctx   context.Context
		newCf *objects.CustomField
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.CustomField
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddCustomField(tt.args.ctx, tt.args.newCf)
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
	type args struct {
		ctx   context.Context
		newCg *objects.ClusterGroup
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.ClusterGroup
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddClusterGroup(tt.args.ctx, tt.args.newCg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddClusterGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddClusterGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddClusterType(t *testing.T) {
	type args struct {
		ctx            context.Context
		newClusterType *objects.ClusterType
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.ClusterType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddClusterType(tt.args.ctx, tt.args.newClusterType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddClusterType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddClusterType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddCluster(t *testing.T) {
	type args struct {
		ctx        context.Context
		newCluster *objects.Cluster
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Cluster
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddCluster(tt.args.ctx, tt.args.newCluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddDeviceRole(t *testing.T) {
	type args struct {
		ctx           context.Context
		newDeviceRole *objects.DeviceRole
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.DeviceRole
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddDeviceRole(tt.args.ctx, tt.args.newDeviceRole)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddDeviceRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddManufacturer(t *testing.T) {
	type args struct {
		ctx             context.Context
		newManufacturer *objects.Manufacturer
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Manufacturer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddManufacturer(tt.args.ctx, tt.args.newManufacturer)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddManufacturer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddManufacturer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddDeviceType(t *testing.T) {
	type args struct {
		ctx           context.Context
		newDeviceType *objects.DeviceType
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.DeviceType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddDeviceType(tt.args.ctx, tt.args.newDeviceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddDeviceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddDeviceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddPlatform(t *testing.T) {
	type args struct {
		ctx         context.Context
		newPlatform *objects.Platform
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Platform
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddPlatform(tt.args.ctx, tt.args.newPlatform)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddPlatform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddDevice(t *testing.T) {
	type args struct {
		ctx       context.Context
		newDevice *objects.Device
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Device
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddDevice(tt.args.ctx, tt.args.newDevice)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddVlanGroup(t *testing.T) {
	type args struct {
		ctx          context.Context
		newVlanGroup *objects.VlanGroup
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.VlanGroup
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddVlanGroup(tt.args.ctx, tt.args.newVlanGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVlanGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddVlanGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddVlan(t *testing.T) {
	type args struct {
		ctx     context.Context
		newVlan *objects.Vlan
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Vlan
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddVlan(tt.args.ctx, tt.args.newVlan)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddVlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddInterface(t *testing.T) {
	type args struct {
		ctx          context.Context
		newInterface *objects.Interface
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Interface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddInterface(tt.args.ctx, tt.args.newInterface)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddVM(t *testing.T) {
	type args struct {
		ctx   context.Context
		newVM *objects.VM
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.VM
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddVM(tt.args.ctx, tt.args.newVM)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddVMInterface(t *testing.T) {
	type args struct {
		ctx            context.Context
		newVMInterface *objects.VMInterface
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.VMInterface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddVMInterface(tt.args.ctx, tt.args.newVMInterface)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVMInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddVMInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddIPAddress(t *testing.T) {
	type args struct {
		ctx          context.Context
		newIPAddress *objects.IPAddress
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.IPAddress
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddIPAddress(tt.args.ctx, tt.args.newIPAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddIPAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddPrefix(t *testing.T) {
	type args struct {
		ctx       context.Context
		newPrefix *objects.Prefix
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		want    *objects.Prefix
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nbi.AddPrefix(tt.args.ctx, tt.args.newPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
