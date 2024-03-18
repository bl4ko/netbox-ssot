package utils

import (
	"context"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/logger"
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
			result := SubnetContainsIPAddress(tt.ipAddress, tt.subnet)
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

func TestValidateRegexRelations(t *testing.T) {
	type args struct {
		regexRelations []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRegexRelations(tt.args.regexRelations); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegexRelations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertStringsToRegexPairs(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertStringsToRegexPairs(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertStringsToRegexPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchStringToValue(t *testing.T) {
	type args struct {
		input    string
		patterns map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MatchStringToValue(tt.args.input, tt.args.patterns)
			if (err != nil) != tt.wantErr {
				t.Errorf("MatchStringToValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MatchStringToValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphanumeric(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Alphanumeric(tt.args.name); got != tt.want {
				t.Errorf("Alphanumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePlatformName(t *testing.T) {
	type args struct {
		osType    string
		osVersion string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePlatformName(tt.args.osType, tt.args.osVersion); got != tt.want {
				t.Errorf("GeneratePlatformName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsVMInterfaceNameValid(t *testing.T) {
	type args struct {
		vmIfaceName string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsVMInterfaceNameValid(tt.args.vmIfaceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsVMInterfaceNameValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsVMInterfaceNameValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractFunctionName(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractFunctionName(tt.args.i); got != tt.want {
				t.Errorf("ExtractFunctionName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mnSet_Contains(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		m    mnSet
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Contains(tt.args.r); got != tt.want {
				t.Errorf("mnSet.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeDiacritics(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDiacritics(tt.args.s); got != tt.want {
				t.Errorf("removeDiacritics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchNamesWithEmails(t *testing.T) {
	type args struct {
		ctx    context.Context
		names  []string
		emails []string
		logger *logger.Logger
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchNamesWithEmails(tt.args.ctx, tt.args.names, tt.args.emails, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchNamesWithEmails() = %v, want %v", got, tt.want)
			}
		})
	}
}
