package dnac

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
)

func newTestDnacSource() *DnacSource {
	return &DnacSource{
		Config: common.Config{
			Logger: &logger.Logger{Logger: log.New(os.Stdout, "", log.LstdFlags)},
			Ctx: context.WithValue(
				context.Background(),
				constants.CtxSourceKey,
				"dnac-test",
			),
			SourceConfig: &parser.SourceConfig{},
		},
	}
}

func TestGetInterfaceDuplex(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name   string
		duplex string
		want   *objects.InterfaceDuplex
	}{
		{"empty string", "", nil},
		{"FullDuplex", "FullDuplex", &objects.DuplexFull},
		{"AutoNegotiate", "AutoNegotiate", &objects.DuplexAuto},
		{"HalfDuplex", "HalfDuplex", &objects.DuplexHalf},
		{"unknown value", "SomeOther", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ds.getInterfaceDuplex(tt.duplex)
			if got != tt.want {
				t.Errorf("getInterfaceDuplex(%q) = %v, want %v", tt.duplex, got, tt.want)
			}
		})
	}
}

func TestGetInterfaceStatus(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name    string
		status  string
		want    bool
		wantErr bool
	}{
		{"down", "down", false, false},
		{"up", "up", true, false},
		{"unknown", "error", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ds.getInterfaceStatus(tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInterfaceStatus(%q) error = %v, wantErr %v", tt.status, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getInterfaceStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestGetInterfaceType(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name    string
		iType   string
		speed   int
		want    *objects.InterfaceType
		wantErr bool
	}{
		{
			name:  "Physical with known speed",
			iType: "Physical",
			speed: int(objects.GBPS1),
			want:  &objects.GE1FixedInterfaceType,
		},
		{
			name:  "Physical with unknown speed",
			iType: "Physical",
			speed: 999,
			want:  &objects.OtherInterfaceType,
		},
		{
			name:  "Virtual",
			iType: "Virtual",
			speed: 0,
			want:  &objects.VirtualInterfaceType,
		},
		{
			name:    "Unknown type",
			iType:   "Wireless",
			speed:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ds.getInterfaceType(tt.iType, tt.speed)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInterfaceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getInterfaceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateInterfaceName(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name    string
		iName   string
		iID     string
		filter  string
		wantErr bool
	}{
		{"empty name", "", "iface-1", "", true},
		{"valid name no filter", "eth0", "iface-1", "", false},
		{"filtered name", "eth0", "iface-1", "^eth.*", true},
		{"non-matching filter", "GigabitEthernet0/0", "iface-2", "^eth.*", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds.SourceConfig.InterfaceFilter = tt.filter
			err := ds.validateInterfaceName(tt.iName, tt.iID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInterfaceName(%q) error = %v, wantErr %v", tt.iName, err, tt.wantErr)
			}
		})
	}
}

func TestGetVlanModeAndAccessVlan(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name     string
		portMode string
		vlanID   string
		wantMode *objects.InterfaceMode
		wantErr  bool
	}{
		{
			name:     "access mode",
			portMode: "access",
			vlanID:   "100",
			wantMode: &objects.InterfaceModeAccess,
		},
		{
			name:     "trunk mode",
			portMode: "trunk",
			vlanID:   "1",
			wantMode: &objects.InterfaceModeTagged,
		},
		{
			name:     "dynamic_auto mode",
			portMode: "dynamic_auto",
			vlanID:   "1",
			wantMode: nil,
		},
		{
			name:     "routed mode",
			portMode: "routed",
			vlanID:   "1",
			wantMode: nil,
		},
		{
			name:     "unknown mode",
			portMode: "unknown",
			vlanID:   "1",
			wantErr:  true,
		},
		{
			name:     "invalid vlan ID",
			portMode: "access",
			vlanID:   "abc",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, _, err := ds.getVlanModeAndAccessVlan(tt.portMode, tt.vlanID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVlanModeAndAccessVlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && mode != tt.wantMode {
				t.Errorf("getVlanModeAndAccessVlan() mode = %v, want %v", mode, tt.wantMode)
			}
		})
	}
}
