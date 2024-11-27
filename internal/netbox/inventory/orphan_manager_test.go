package inventory

import (
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestNewOrphanManager(t *testing.T) {
	type args struct {
		logger *logger.Logger
	}
	tests := []struct {
		name string
		args args
		want *OrphanManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrphanManager(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrphanManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrphanManager_AddItem(t *testing.T) {
	type args struct {
		itemAPIPath string
		orphanItem  objects.OrphanItem
	}
	tests := []struct {
		name          string
		orphanManager *OrphanManager
		args          args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tt.orphanManager.AddItem(tt.args.itemAPIPath, tt.args.orphanItem)
		})
	}
}

func TestOrphanManager_RemoveItem(t *testing.T) {
	type args struct {
		itemAPIPath string
		obj         objects.OrphanItem
	}
	tests := []struct {
		name          string
		orphanManager *OrphanManager
		args          args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tt.orphanManager.RemoveItem(tt.args.itemAPIPath, tt.args.obj)
		})
	}
}
