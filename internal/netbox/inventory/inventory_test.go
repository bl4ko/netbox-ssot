package inventory

import (
	"context"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/parser"
)

func TestNetboxInventory_String(t *testing.T) {
	tests := []struct {
		name string
		nbi  *NetboxInventory
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.nbi.String(); got != tt.want {
				t.Errorf("NetboxInventory.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNetboxInventory(t *testing.T) {
	type args struct {
		ctx      context.Context
		logger   *logger.Logger
		nbConfig *parser.NetboxConfig
	}
	tests := []struct {
		name string
		args args
		want *NetboxInventory
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNetboxInventory(tt.args.ctx, tt.args.logger, tt.args.nbConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNetboxInventory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxInventory_Init(t *testing.T) {
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.Init(); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
