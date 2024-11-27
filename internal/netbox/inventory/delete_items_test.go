package inventory

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestNetboxInventory_DeleteOrphans(t *testing.T) {
	type args struct {
		hard bool
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
			if err := tt.nbi.DeleteOrphans(tt.args.hard); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.DeleteOrphans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_hardDelete(t *testing.T) {
	type args struct {
		apiPath    string
		orphanItem objects.OrphanItem
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
			if err := tt.nbi.hardDelete(tt.args.apiPath, tt.args.orphanItem); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.hardDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetboxInventory_softDelete(t *testing.T) {
	type args struct {
		apiPath    string
		orphanItem objects.OrphanItem
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
			if err := tt.nbi.softDelete(tt.args.apiPath, tt.args.orphanItem); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.softDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
