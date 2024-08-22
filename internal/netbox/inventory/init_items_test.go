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
			if err := tt.nbi.initTags(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initTenants(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initContacts(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initContactRoles(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initContactAssignments(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initAdminContactRole(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initContactGroups(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initSites(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initManufacturers(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initPlatforms(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initDevices(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initDeviceRoles(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitDeviceRoles() error = %v, wantErr %v", err, tt.wantErr)
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
			if err := tt.nbi.initCustomFields(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initSsotCustomFields(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initClusterGroups(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initClusterTypes(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initClusters(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initDeviceTypes(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initInterfaces(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initDefaultVlanGroup(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initVlanGroups(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initVlans(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initVMs(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initVMInterfaces(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initIPAddresses(tt.args.ctx); (err != nil) != tt.wantErr {
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
			if err := tt.nbi.initPrefixes(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.InitPrefixes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
