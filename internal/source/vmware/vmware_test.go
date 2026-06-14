package vmware

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
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
