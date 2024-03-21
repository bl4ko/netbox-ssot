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
			name: "Valid reverse lookup",
			args: args{
				ipAddress: "8.8.8.8",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubnetsContainIPAddress(tt.args.ipAddress, tt.args.subnets); got != tt.want {
				t.Errorf("SubnetsContainIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractPrefixFromIPAddress(t *testing.T) {
	type args struct {
		ipAddress string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Extract valid prefix",
			args: args{
				ipAddress: "172.16.1.2/16",
			},
			want:    "172.16.0.0/16",
			wantErr: false,
		},
		{
			name: "Extract valid prefix",
			args: args{
				ipAddress: "192.168.1.2/24",
			},
			want:    "192.168.1.0/24",
			wantErr: false,
		},
		{
			name: "Extract invalid prefix",
			args: args{
				ipAddress: "192.168.1.2",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Valid ipv6 prefix",
			args: args{
				ipAddress: "2001:db8::1/32",
			},
			want:    "2001:db8::/32",
			wantErr: false,
		},
		{
			name: "Invalid ipv6 prefix",
			args: args{
				ipAddress: "2001:db8::1",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Ignore 128 mask for ipv6",
			args: args{
				ipAddress: "2001:db8::1/128",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Ignore 32 mask for ipv4",
			args: args{
				ipAddress: "192:168:0:64/32",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPrefixFromIPAddress(tt.args.ipAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPrefixFromIPAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractPrefixFromIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
