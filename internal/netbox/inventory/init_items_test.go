package inventory

import (
	"context"
	"testing"
)

func TestNetboxInventory_InitTags(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitTags(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitTenants(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitTenants(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitTenants() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitContacts(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitContacts(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitContacts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitContactRoles(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitContactRoles(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitContactRoles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitContactAssignments(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitContactAssignments(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitContactAssignments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitAdminContactRole(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitAdminContactRole(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitAdminContactRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitContactGroups(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitContactGroups(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitContactGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitSites(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitSites(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitSites() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitManufacturers(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitManufacturers(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitManufacturers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitPlatforms(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitPlatforms(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitPlatforms() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitDevices(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitDevices(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitDevices() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitDeviceRoles(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitDeviceRoles(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitDeviceRoles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitServerDeviceRole(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitServerDeviceRole(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitServerDeviceRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitCustomFields(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitCustomFields(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitCustomFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitSsotCustomFields(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitSsotCustomFields(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitSsotCustomFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitClusterGroups(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitClusterGroups(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitClusterGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitClusterTypes(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitClusterTypes(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitClusterTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitClusters(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitClusters(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitClusters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitDeviceTypes(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitDeviceTypes(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitDeviceTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitInterfaces(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitInterfaces(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitInterfaces() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitDefaultVlanGroup(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitDefaultVlanGroup(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitDefaultVlanGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitVlanGroups(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitVlanGroups(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitVlanGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitVlans(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitVlans(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitVlans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitVMs(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitVMs(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitVMs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitVMInterfaces(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitVMInterfaces(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitVMInterfaces() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitIPAddresses(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitIPAddresses(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitIPAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_InitPrefixes(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.InitPrefixes(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitPrefixes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
