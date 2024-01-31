package utils

import (
	"reflect"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string",
			input:    "Test",
			expected: "test",
		},
		{
			name:     "String with spaces",
			input:    "Test String",
			expected: "test_string",
		},
		{
			name:     "String with trailing spaces",
			input:    "    Te st    ",
			expected: "te_st",
		},
		{
			name:     "String with special characters",
			input:    "Test@#String$%^",
			expected: "teststring",
		},
		{
			name:     "String with mixed case letters",
			input:    "TeSt StRiNg",
			expected: "test_string",
		},
		{
			name:     "String with numbers",
			input:    "Test123 String456",
			expected: "test123_string456",
		},
		{
			name:     "String with underscores",
			input:    "Test_String",
			expected: "test_string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := Slugify(tt.input)
			if slug != tt.expected {
				t.Errorf("Slugify() = %v, want %v", slug, tt.expected)
			}
		})
	}
}

func TestFilterVMInterfaceNames(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "No interfaces",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "No VM interfaces",
			input:    []string{"eth0", "eth1", "eth2"},
			expected: []string{"eth0", "eth1", "eth2"},
		},
		{
			name:     "One VM interface",
			input:    []string{"eth0", "docker0", "eth1", "cali7839a755dc1"},
			expected: []string{"eth0", "eth1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filteredSlice := make([]string, 0)
			for _, iface := range tt.input {
				filtered, err := IsVMInterfaceNameValid(iface)
				if err != nil {
					t.Errorf("FilterVMInterfaceNames() error = %v", err)
				}
				if filtered == true {
					filteredSlice = append(filteredSlice, iface)
				}
			}

			if !reflect.DeepEqual(filteredSlice, tt.expected) {
				t.Errorf("FilterVMInterfaceNames() = %v, want %v", filteredSlice, tt.expected)
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

func TestSubnetContainsIpAddress(t *testing.T) {
	tests := []struct {
		name      string
		ipAddress string
		subnet    string
		expected  bool
	}{
		{
			name:      "IP in Subnet",
			ipAddress: "172.31.4.129",
			subnet:    "172.31.4.145/25",
			expected:  true,
		},
		{
			name:      "IP not in Subnet",
			ipAddress: "192.168.1.1",
			subnet:    "172.31.4.145/25",
			expected:  false,
		},
		{
			name:      "Invalid IP Address",
			ipAddress: "invalid",
			subnet:    "172.31.4.145/25",
			expected:  false,
		},
		{
			name:      "Invalid Subnet",
			ipAddress: "172.31.4.129",
			subnet:    "invalid",
			expected:  false,
		},
		{
			name:      "Empty IP Address",
			ipAddress: "",
			subnet:    "172.31.4.145/25",
			expected:  false,
		},
		{
			name:      "Empty Subnet",
			ipAddress: "172.31.4.129",
			subnet:    "",
			expected:  false,
		},
		{
			name:      "IPv6 IP in Subnet",
			ipAddress: "2001:db8::1",
			subnet:    "2001:db8::/32",
			expected:  true,
		},
		{
			name:      "IPv6 IP not in Subnet",
			ipAddress: "2001:db8::1",
			subnet:    "2001:db7::/32",
			expected:  false,
		},
		{
			name:      "Invalid IPv6 Address",
			ipAddress: "2001:db8::zzz",
			subnet:    "2001:db8::/32",
			expected:  false,
		},
		{
			name:      "Invalid IPv6 CIDR",
			ipAddress: "2001:db8::1",
			subnet:    "2001:db8::zzz/32",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SubnetContainsIpAddress(tt.ipAddress, tt.subnet)
			if result != tt.expected {
				t.Errorf("SubnetContainsIpAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertStringsToPairs(t *testing.T) {
	input := []string{"key1=value1", "key2=value2", "key3=value3"}
	output := ConvertStringsToPairs(input)
	desiredOutput := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	if !reflect.DeepEqual(output, desiredOutput) {
		t.Errorf("ConvertStringsToPairs() = %v, want %v", output, desiredOutput)
	}
}
