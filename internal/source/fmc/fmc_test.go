package fmc

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/source/fmc/client"
)

func TestGetIPAddressForIface(t *testing.T) {
	tests := []struct {
		name string
		ipv4 *client.InterfaceIPv4
		want string
	}{
		{
			name: "nil ipv4",
			ipv4: nil,
			want: "",
		},
		{
			name: "static address",
			ipv4: &client.InterfaceIPv4{
				Static: &struct {
					Address string `json:"address"`
					Netmask string `json:"netmask"`
				}{
					Address: "10.0.0.1",
					Netmask: "255.255.255.0",
				},
			},
			want: "10.0.0.1/24",
		},
		{
			name: "dhcp address",
			ipv4: &client.InterfaceIPv4{
				Dhcp: &struct {
					Address string `json:"address"`
					Netmask string `json:"netmask"`
				}{
					Address: "192.168.1.100",
					Netmask: "255.255.255.0",
				},
			},
			want: "192.168.1.100/24",
		},
		{
			name: "empty static address",
			ipv4: &client.InterfaceIPv4{
				Static: &struct {
					Address string `json:"address"`
					Netmask string `json:"netmask"`
				}{
					Address: "",
					Netmask: "255.255.255.0",
				},
			},
			want: "",
		},
		{
			name: "static takes priority over dhcp",
			ipv4: &client.InterfaceIPv4{
				Static: &struct {
					Address string `json:"address"`
					Netmask string `json:"netmask"`
				}{
					Address: "10.0.0.1",
					Netmask: "255.255.0.0",
				},
				Dhcp: &struct {
					Address string `json:"address"`
					Netmask string `json:"netmask"`
				}{
					Address: "192.168.1.100",
					Netmask: "255.255.255.0",
				},
			},
			want: "10.0.0.1/16",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getIPAddressForIface(tt.ipv4)
			if got != tt.want {
				t.Errorf("getIPAddressForIface() = %q, want %q", got, tt.want)
			}
		})
	}
}
