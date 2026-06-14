package iosxe

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestIOSXEPortSpeedToLinkSpeed(t *testing.T) {
	tests := []struct {
		name      string
		portSpeed string
		want      objects.InterfaceSpeed
	}{
		{name: "10 Mbit", portSpeed: "SPEED_10MB", want: 10_000},
		{name: "100 Mbit", portSpeed: "SPEED_100MB", want: 100_000},
		{name: "1 Gbit", portSpeed: "SPEED_1GB", want: 1_000_000},
		{name: "10 Gbit", portSpeed: "SPEED_10GB", want: 10_000_000},
		{name: "25 Gbit", portSpeed: "SPEED_25GB", want: 25_000_000},
		{name: "40 Gbit", portSpeed: "SPEED_40GB", want: 40_000_000},
		{name: "100 Gbit", portSpeed: "SPEED_100GB", want: 100_000_000},
		{name: "unknown speed", portSpeed: "SPEED_400GB", want: 0},
		{name: "empty speed", portSpeed: "", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iosxePortSpeedToLinkSpeed(tt.portSpeed); got != tt.want {
				t.Errorf("iosxePortSpeedToLinkSpeed(%q) = %d, want %d", tt.portSpeed, got, tt.want)
			}
		})
	}
}
