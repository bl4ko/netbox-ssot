package inventory

import (
	"context"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestNetboxInventory_AddContainerDeviceRole(t *testing.T) {
	type args struct {
		ctx context.Context
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
			got, err := tt.nbi.AddContainerDeviceRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddContainerDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddContainerDeviceRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddFirewallDeviceRole(t *testing.T) {
	type args struct {
		ctx context.Context
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
			got, err := tt.nbi.AddFirewallDeviceRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddFirewallDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddFirewallDeviceRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddSwitchDeviceRole(t *testing.T) {
	type args struct {
		ctx context.Context
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
			got, err := tt.nbi.AddSwitchDeviceRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddSwitchDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddSwitchDeviceRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddServerDeviceRole(t *testing.T) {
	type args struct {
		ctx context.Context
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
			got, err := tt.nbi.AddServerDeviceRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddServerDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddServerDeviceRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddVMDeviceRole(t *testing.T) {
	type args struct {
		ctx context.Context
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
			got, err := tt.nbi.AddVMDeviceRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVMDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddVMDeviceRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_AddVMTemplateDeviceRole(t *testing.T) {
	type args struct {
		ctx context.Context
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
			got, err := tt.nbi.AddVMTemplateDeviceRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.AddVMTemplateDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxInventory.AddVMTemplateDeviceRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
