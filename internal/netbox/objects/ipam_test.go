package objects

import (
	"reflect"
	"testing"
)

func TestIPAddress_String(t *testing.T) {
	tests := []struct {
		name string
		ip   IPAddress
		want string
	}{
		{
			name: "Test ip address correct string",
			ip: IPAddress{
				NetboxObject: NetboxObject{
					ID: 1,
				},
				Address: "172.18.19.20",
				Status:  &IPAddressStatusActive,
				DNSName: "test.example.com",
			},
			want: "IPAddress{ID: 1, Address: 172.18.19.20, Status: active, DNSName: test.example.com}",
		},
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
		{
			name: "Test vlan group correct string",
			vg: VlanGroup{
				Name:   "Test vlan group",
				MinVid: 1,
				MaxVid: 4094,
			},
			want: "VlanGroup{Name: Test vlan group, MinVid: 1, MaxVid: 4094}",
		},
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
		{
			name: "Test vlan correct string",
			v: Vlan{
				NetboxObject: NetboxObject{
					ID: 1,
				},
				Name:   "Test vlan",
				Vid:    10,
				Status: &VlanStatusActive,
			},
			want: "Vlan{ID: 1, Name: Test vlan, Vid: 10, Status: active}",
		},
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
		{
			name: "Test prefix correct string",
			p: Prefix{
				Prefix: "172.18.19.20/24",
			},
			want: "Prefix{Prefix: 172.18.19.20/24}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Prefix.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPAddress_GetID(t *testing.T) {
	tests := []struct {
		name string
		ip   *IPAddress
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ip.GetID(); got != tt.want {
				t.Errorf("IPAddress.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPAddress_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		ip   *IPAddress
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ip.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPAddress.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVlanGroup_GetID(t *testing.T) {
	tests := []struct {
		name string
		vg   *VlanGroup
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vg.GetID(); got != tt.want {
				t.Errorf("VlanGroup.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVlanGroup_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		vg   *VlanGroup
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vg.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VlanGroup.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVlan_GetID(t *testing.T) {
	tests := []struct {
		name string
		v    *Vlan
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.GetID(); got != tt.want {
				t.Errorf("Vlan.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVlan_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		v    *Vlan
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vlan.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrefix_GetID(t *testing.T) {
	tests := []struct {
		name string
		p    *Prefix
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetID(); got != tt.want {
				t.Errorf("Prefix.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrefix_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		p    *Prefix
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prefix.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
