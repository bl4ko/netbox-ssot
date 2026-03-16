package inventory

import (
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/logger"
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
