package utils

import (
	"slices"
	"testing"
)

func TestReverseLookup(t *testing.T) {
	type args struct {
		ipAddress string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Invalid reverse lookup",
			args: args{
				ipAddress: "890.1.324.234",
			},
			want: "",
		},
		{
			name: "Valid rerse lookup",
			args: args{
				ipAddress: "1.1.1.1",
			},
			want: "one.one.one.one",
		},
		{
			name: "Valid reverse lookup",
			args: args{
				ipAddress: "8.8.8.8",
			},
			want: "dns.google",
		},
		{
			name: "Lookup IP with mask",
			args: args{
				ipAddress: "8.8.8.8/24",
			},
			want: "dns.google",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseLookup(tt.args.ipAddress); got != tt.want {
				t.Errorf("ReverseLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	type args struct {
		hostname string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test valid address lookup",
			args: args{
				hostname: "dns.google",
			},
			want: []string{"8.8.4.4", "8.8.8.8"},
		},
		{
			name: "Test invalid address lookup",
			args: args{
				hostname: "example.invalid",
			},
			want: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Lookup(tt.args.hostname); !slices.Contains(tt.want, got) {
				t.Errorf("Lookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskToBits(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		err      bool
	}{
		{
			name:     "Valid mask 255.255.255.128",
			input:    "255.255.255.128",
			expected: 25,
			err:      false,
		},
		{
			name:     "Valid mask 255.255.255.0",
			input:    "255.255.255.0",
			expected: 24,
			err:      false,
		},
		{
			name:     "Invalid mask",
			input:    "255.255.255.256",
			expected: 0,
			err:      true,
		},
		{
			name:     "Empty mask",
			input:    "",
			expected: 0,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bits, err := MaskToBits(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("maskToBits() error = %v, wantErr %v", err, tt.err)
				return
			}
			if bits != tt.expected {
				t.Errorf("maskToBits() = %v, want %v", bits, tt.expected)
			}
		})
	}
}

func TestGetIPVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Valid IPv4",
			input:    "192.168.1.1",
			expected: 4,
		},
		{
			name:     "Valid IPv6",
			input:    "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: 6,
		},
		{
			name:     "Invalid IP",
			input:    "invalid",
			expected: 0,
		},
		{
			name:     "Empty IP",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := GetIPVersion(tt.input)
			if version != tt.expected {
				t.Errorf("GetIPVersion() = %v, want %v", version, tt.expected)
			}
		})
	}
}

func TestSubnetContainsIPAddress(t *testing.T) {
	type args struct {
		ipAddress string
		subnet    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Subnet contains ip address",
			args: args{
				ipAddress: "192.168.1.1",
				subnet:    "192.168.1.0/24",
			},
			want: true,
		},
		{
			name: "Subnet doesn't contain ip address",
			args: args{
				ipAddress: "192.168.1.1",
				subnet:    "192.168.0.0/24",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubnetContainsIPAddress(tt.args.ipAddress, tt.args.subnet); got != tt.want {
				t.Errorf("SubnetContainsIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifySubnet(t *testing.T) {
	type args struct {
		subnet string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test valid subnet",
			args: args{
				subnet: "172.16.0.0/24",
			},
			want: true,
		},
		{
			name: "Test invalid subnet",
			args: args{
				subnet: "172.16.0.257/24",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifySubnet(tt.args.subnet); got != tt.want {
				t.Errorf("VerifySubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubnetsContainIPAddress(t *testing.T) {
	type args struct {
		ipAddress string
		subnets   []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test subnets contain ip address",
			args: args{
				ipAddress: "192.168.1.2",
				subnets:   []string{"172.16.0.0/16", "10.0.0.0/8", "192.168.1.0/24"},
			},
			want: true,
		},
		{
			name: "Test subnets contain ip address",
			args: args{
				ipAddress: "192.168.1.2",
				subnets:   []string{"172.16.0.0/16", "10.0.0.0/8"},
			},
			want: false,
		},
		{
			name: "Test ipv6 should be filtered out",
			args: args{
				ipAddress: "fe80::2744:6e22:ce45:c4b3",
				subnets:   []string{"::/0"},
			},
			want: true,
		},
		{
			name: "Test ipv6 with zone should be filtered out",
			args: args{
				ipAddress: "fe80::2744:6e22:ce45:c4b3%27/64",
				subnets:   []string{"::/0"},
			},
			want: true,
		},
		{
			name: "Doesn't contain ipv6 address",
			args: args{
				ipAddress: "fe80::2744:6e22:ce45:c4b3",
				subnets:   []string{"10.0.0.0/8"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubnetsContainIPAddress(tt.args.ipAddress, tt.args.subnets); got != tt.want {
				t.Errorf("SubnetsContainIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetMaskAndPrefixFromIPAddress(t *testing.T) {
	type args struct {
		ipAddress string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		{
			name: "Test valid ip address with valid mask",
			args: args{
				ipAddress: "192.168.1.1/24",
			},
			want:  "192.168.1.0/24",
			want1: 24,
		},
		{
			name: "Test valid ip address with valid mask",
			args: args{
				ipAddress: "192.168.1.1/32",
			},
			want:  "192.168.1.1/32",
			want1: 32,
		},
		{
			name: "Error wrong ipv4 address",
			args: args{
				ipAddress: "300.168.1.1/24",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetPrefixAndMaskFromIPAddress(tt.args.ipAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMaskAndPrefixFromIPAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMaskAndPrefixFromIPAddress() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetMaskAndPrefixFromIPAddress() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetPrefixAndMaskFromIPAddress(t *testing.T) {
	type args struct {
		ipAddress string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetPrefixAndMaskFromIPAddress(tt.args.ipAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPrefixAndMaskFromIPAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPrefixAndMaskFromIPAddress() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetPrefixAndMaskFromIPAddress() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRemoveZoneIndexFromIPAddress(t *testing.T) {
	type args struct {
		ipAddress string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test valid ip address",
			args: args{
				ipAddress: "192.168.1.0",
			},
			want: "192.168.1.0",
		},
		{
			name: "Test zone index ipv6 address",
			args: args{
				ipAddress: "fe80::2744:6e22:ce45:c4b3%27/64",
			},
			want: "fe80::2744:6e22:ce45:c4b3/64",
		},
		{
			name: "Test zone index ipv6 address",
			args: args{
				ipAddress: "fe80::2744:6e22:ce45:c4b3%eth1",
			},
			want: "fe80::2744:6e22:ce45:c4b3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveZoneIndexFromIPAddress(tt.args.ipAddress); got != tt.want {
				t.Errorf("RemoveZoneIndexFromIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
