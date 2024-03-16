package objects

import "testing"

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
