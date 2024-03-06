package objects

import "testing"

func TestIPAddress_String(t *testing.T) {
	tests := []struct {
		name string
		ip   IPAddress
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ip.String(); got != tt.want {
				t.Errorf("IPAddress.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVlanGroup_String(t *testing.T) {
	tests := []struct {
		name string
		vg   VlanGroup
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vg.String(); got != tt.want {
				t.Errorf("VlanGroup.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVlan_String(t *testing.T) {
	tests := []struct {
		name string
		v    Vlan
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("Vlan.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrefix_String(t *testing.T) {
	tests := []struct {
		name string
		p    Prefix
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Prefix.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
