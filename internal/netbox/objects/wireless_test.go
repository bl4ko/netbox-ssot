package objects

import (
	"reflect"
	"testing"
)

func TestWirelessLAN_String(t *testing.T) {
	tests := []struct {
		name string
		wl   WirelessLAN
		want string
	}{
		{
			name: "Test string of a siimple WirelessLan",
			wl: WirelessLAN{
				NetboxObject: NetboxObject{
					ID:           1,
					CustomFields: map[string]interface{}{},
				},
				SSID: "Test",
			},
			want: "WirelessLAN{SSID: Test}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wl.String(); got != tt.want {
				t.Errorf("WirelessLan.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWirelessLANGroup_String(t *testing.T) {
	tests := []struct {
		name string
		wlg  WirelessLANGroup
		want string
	}{
		{
			name: "Test string of a simple WirelessLANGroup",
			wlg: WirelessLANGroup{
				NetboxObject: NetboxObject{
					ID:           1,
					CustomFields: map[string]interface{}{},
				},
				Name: "Test",
				Slug: "test",
			},
			want: "WirelessLANGroup{Name: Test, Slug: test}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wlg.String(); got != tt.want {
				t.Errorf("WirelessLANGroup.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWirelessLANGroup_GetID(t *testing.T) {
	tests := []struct {
		name string
		wlg  *WirelessLANGroup
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wlg.GetID(); got != tt.want {
				t.Errorf("WirelessLANGroup.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWirelessLANGroup_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		wlg  *WirelessLANGroup
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wlg.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WirelessLANGroup.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWirelessLAN_GetID(t *testing.T) {
	tests := []struct {
		name string
		wl   *WirelessLAN
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wl.GetID(); got != tt.want {
				t.Errorf("WirelessLAN.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWirelessLAN_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		wl   *WirelessLAN
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wl.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WirelessLAN.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
