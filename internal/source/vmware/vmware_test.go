package vmware

import (
	"context"
	"slices"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func TestVmwareConnectionStateToDeviceStatus(t *testing.T) {
	tests := []struct {
		name  string
		state types.HostSystemConnectionState
		want  *objects.DeviceStatus
	}{
		{name: "connected", state: "connected", want: &objects.DeviceStatusActive},
		{name: "disconnected", state: "disconnected", want: &objects.DeviceStatusOffline},
		{name: "not responding", state: "notResponding", want: &objects.DeviceStatusOffline},
		{name: "empty", state: "", want: &objects.DeviceStatusOffline},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vmwareConnectionStateToDeviceStatus(tt.state)
			if *got != *tt.want {
				t.Errorf("vmwareConnectionStateToDeviceStatus(%q) = %v, want %v", tt.state, *got, *tt.want)
			}
		})
	}
}

func TestCollectVMInterfaceDataMatchesGuestNicMacCaseInsensitively(t *testing.T) {
	testLogger, err := logger.New("", logger.ERROR)
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	vc := &VmwareSource{
		Config: common.Config{
			Logger:       testLogger,
			Ctx:          context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
			SourceConfig: &parser.SourceConfig{Name: "test"},
		},
	}
	netboxVM := &objects.VM{Name: "vm1", Host: &objects.Device{Name: "host1"}}

	tests := []struct {
		name     string
		hwMac    string
		guestMac string
		wantIPv4 []string
		wantIPv6 []string
		wantMac  string
		wantUp   bool
	}{
		{
			name:     "case-only mismatch lower guest",
			hwMac:    "00:50:56:AB:CD:EF",
			guestMac: "00:50:56:ab:cd:ef",
			wantIPv4: []string{"10.0.0.5/24"},
			wantIPv6: []string{"fd00::5/64"},
			wantMac:  "00:50:56:AB:CD:EF",
			wantUp:   true,
		},
		{
			name:     "case-only mismatch lower hardware",
			hwMac:    "00:50:56:ab:cd:ef",
			guestMac: "00:50:56:AB:CD:EF",
			wantIPv4: []string{"10.0.0.5/24"},
			wantIPv6: []string{"fd00::5/64"},
			wantMac:  "00:50:56:AB:CD:EF",
			wantUp:   true,
		},
		{
			name:     "exact match",
			hwMac:    "00:50:56:ab:cd:ef",
			guestMac: "00:50:56:ab:cd:ef",
			wantIPv4: []string{"10.0.0.5/24"},
			wantIPv6: []string{"fd00::5/64"},
			wantMac:  "00:50:56:AB:CD:EF",
			wantUp:   true,
		},
		{
			name:     "different mac",
			hwMac:    "00:50:56:ab:cd:ef",
			guestMac: "00:50:56:ab:cd:00",
			wantIPv4: []string{},
			wantIPv6: []string{},
			wantMac:  "00:50:56:AB:CD:EF",
			wantUp:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := &types.VirtualEthernetCard{
				VirtualDevice: types.VirtualDevice{
					Backing:     &types.VirtualEthernetCardNetworkBackingInfo{},
					Connectable: &types.VirtualDeviceConnectInfo{Connected: false},
					DeviceInfo:  &types.Description{Label: "Network adapter 1"},
				},
				MacAddress: tt.hwMac,
			}
			vmwareVM := mo.VirtualMachine{
				Guest: &types.GuestInfo{
					Net: []types.GuestNicInfo{{
						MacAddress: tt.guestMac,
						Connected:  true,
						IpConfig: &types.NetIpConfigInfo{
							IpAddress: []types.NetIpConfigInfoIpAddress{
								{IpAddress: "10.0.0.5", PrefixLength: 24},
								{IpAddress: "fd00::5", PrefixLength: 64},
							},
						},
					}},
				},
			}

			ipv4, ipv6, iface, mac, err := vc.collectVMInterfaceData(nil, netboxVM, vmwareVM, card)
			if err != nil {
				t.Fatalf("collectVMInterfaceData() error = %v", err)
			}
			if !slices.Equal(ipv4, tt.wantIPv4) {
				t.Errorf("ipv4 = %v, want %v", ipv4, tt.wantIPv4)
			}
			if !slices.Equal(ipv6, tt.wantIPv6) {
				t.Errorf("ipv6 = %v, want %v", ipv6, tt.wantIPv6)
			}
			if mac != tt.wantMac {
				t.Errorf("mac = %q, want %q", mac, tt.wantMac)
			}
			if iface.Enabled != tt.wantUp {
				t.Errorf("enabled = %v, want %v", iface.Enabled, tt.wantUp)
			}
		})
	}
}
