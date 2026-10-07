package objects

import (
	"fmt"
	"testing"
)

func TestStringWithoutNestedObjects(t *testing.T) {
	tests := []struct {
		name   string
		object fmt.Stringer
		want   string
	}{
		{
			name:   "device type without manufacturer",
			object: DeviceType{Model: "PowerEdge R650"},
			want:   "DeviceType{Manufacturer: <nil>, Model: PowerEdge R650}",
		},
		{
			name:   "interface without device nor type",
			object: Interface{Name: "eth0"},
			want:   "Interface{Name: eth0, Device: <nil>, Type: <nil>}",
		},
		{
			name:   "vlan without group",
			object: Vlan{NetboxObject: NetboxObject{ID: 3}, Name: "mgmt", Vid: 10},
			want:   "Vlan{ID: 3, Name: mgmt, Vid: 10, Group: <nil>}",
		},
		{
			name:   "vm interface without vm",
			object: VMInterface{Name: "vNIC 1"},
			want:   "VMInterface{Name: vNIC 1, VM: <nil>}",
		},
		{
			name:   "virtual disk without vm",
			object: VirtualDisk{Name: "Hard disk 1"},
			want:   "VirtualDisk{Name: Hard disk 1, VM: <nil>}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprint(tt.object)
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
