package utils

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
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
			expected: "test-string",
		},
		{
			name:     "String with trailing spaces",
			input:    "    Te st    ",
			expected: "te-st",
		},
		{
			name:     "String with special characters",
			input:    "Test@#String$%^",
			expected: "teststring",
		},
		{
			name:     "String with mixed case letters",
			input:    "TeSt StRiNg",
			expected: "test-string",
		},
		{
			name:     "String with numbers",
			input:    "Test123 String456",
			expected: "test123-string456",
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
		{
			name: "Test valid regex relations",
			args: args{
				regexRelations: []string{"disney.* = disney", "p([a-z]+)ch = peach"},
			},
			wantErr: false,
		},
		{
			name: "Test missing equal sign in regex relations",
			args: args{
				regexRelations: []string{"disney.* = disney", "p([a-z]+)ch peach"},
			},
			wantErr: true,
		},
		{
			name: "Test wrong regex relations",
			args: args{
				regexRelations: []string{"a(b(c = disney", "p([a-z]+)ch = peach"},
			},
			wantErr: true,
		},
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
		{
			name: "TEst conversion of strings to regex pairs",
			args: args{
				input: []string{
					"regex1 = value1", "regex2 = value2", "regex3 = value3",
				},
			},
			want: map[string]string{"regex1": "value1", "regex2": "value2", "regex3": "value3"},
		},
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
		{
			name: "Match string to value error",
			args: args{
				input: "can't see me",
				patterns: map[string]string{
					"$$$wrongregex\\": "wrong",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Match string to value valid",
			args: args{
				input: "can't see me",
				patterns: map[string]string{
					"can't": "true",
				},
			},
			want:    "true",
			wantErr: false,
		},
		{
			name: "No error with no matches",
			args: args{
				input: "can't see me",
				patterns: map[string]string{
					"john cena": "true",
				},
			},
			want:    "",
			wantErr: false,
		},
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
		{
			name: "Test alphanumeric conversion",
			args: args{
				name: "Fix-me99 now",
			},
			want: "fixme99_now",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Alphanumeric(tt.args.name); got != tt.want {
				t.Errorf("Alphanumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsVMInterfaceNameValid(t *testing.T) {
	type args struct {
		vmIfaceName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test IsVMInterfaceNameValid works for all vm names",
			args: args{
				vmIfaceName: "\\$\\$\\",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterInterfaceName("docker", "docker")
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
		name        string
		args        args
		want        string
		shouldPanic bool
	}{
		{
			name: "Test valid extract function name",
			args: args{
				i: TestExtractFunctionName,
			},
			want:        "TestExtractFunctionName",
			shouldPanic: false,
		},
		{
			name: "Test panic for non-function",
			args: args{
				i: "not a function",
			},
			want:        "", // Irrelevant for a panic scenario
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.shouldPanic {
					t.Errorf("ExtractFunctionName panicked = %v, want %v", r != nil, tt.shouldPanic)
				}
			}()

			got := ExtractFunctionName(tt.args.i)
			if got != tt.want && !tt.shouldPanic {
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
		{
			name: "Test MatchNamesWithEmails match name",
			args: args{
				ctx:    context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				names:  []string{"John Doe", "Jane Doe"},
				emails: []string{"john.doe@example.com"},
				logger: &logger.Logger{Logger: log.Default()},
			},
			want: map[string]string{
				"John Doe": "john.doe@example.com",
			},
		},
		{
			name: "Test MatchNamesWithEmails no match name",
			args: args{
				ctx:    context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				names:  []string{"Jane Doe"},
				emails: []string{"john.doe@example.com"},
				logger: &logger.Logger{Logger: log.Default()},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchNamesWithEmails(tt.args.ctx, tt.args.names, tt.args.emails, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchNamesWithEmails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterInterfaceName(t *testing.T) {
	type args struct {
		ifaceName   string
		ifaceFilter string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test filter out interface",
			args: args{
				ifaceName:   "docker0",
				ifaceFilter: "^(docker|veth)\\w*",
			},
			want: true,
		},
		{
			name: "Test don't filter out interface",
			args: args{
				ifaceName:   "eth5",
				ifaceFilter: "^(docker|veth)\\w*",
			},
			want: false,
		},
		{
			name: "Test don't filter when empty filter",
			args: args{
				ifaceName:   "eth5",
				ifaceFilter: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterInterfaceName(tt.args.ifaceName, tt.args.ifaceFilter); got != tt.want {
				t.Errorf("FilterInterfaceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadExtraCert(t *testing.T) {
	tests := []struct {
		name     string
		certPath string
		wantErr  bool
	}{
		{
			name:     "No cert load (empty string)",
			certPath: "",
			wantErr:  false,
		},
		{
			name:     "Load valid cert",
			certPath: "../../testdata/certificate/cert.pem",
			wantErr:  false,
		},
		{
			name:     "Throw error on non-existent path",
			certPath: "non existent path",
			wantErr:  true,
		},
		{
			name:     "Throw error on invalid cert data",
			certPath: "invalid cert path",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadExtraCert(tt.certPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadExtraCert(%q) error = %v, wantErr %v", tt.certPath, err, tt.wantErr)
			}
			if err == nil && got == nil {
				t.Errorf("LoadExtraCert(%q) = nil, expected non-nil CertPool", tt.certPath)
			}
		})
	}
}

func TestLoadExtraCertInTransportConfig(t *testing.T) {
	_, err := LoadExtraCertInTransportConfig("wrong path")
	if err == nil {
		t.Errorf("should throw error on wrong path")
	}

	_, err = LoadExtraCertInTransportConfig("../../testdata/certificate/cert.pem")
	if err != nil {
		t.Errorf("error loading extra cert: %s", err)
	}
}

func TestGetManufacturerFromString(t *testing.T) {
	type args struct {
		manufacturer string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Cisco System Inc. should be mapped to Cisco",
			args: args{
				manufacturer: "Cisco Systems Inc.",
			},
			want: "Cisco",
		},
		{
			name: "Test GetManufacturerFromString with HP string",
			args: args{
				manufacturer: "HP",
			},
			want: "HPE",
		},
		{
			name: "Test HP Enterprise should be mapped to HPE",
			args: args{
				manufacturer: "HP Enterprise",
			},
			want: "HPE",
		},
		{
			name: "Test HP Inc. should be mapped to HPE",
			args: args{
				manufacturer: "HP Inc.",
			},
			want: "HPE",
		},
		{
			name: "Test Fortinet Inc. should be mapped to Fortinet",
			args: args{
				manufacturer: "Fortinet Inc.",
			},
			want: "Fortinet",
		},
		{
			name: "Test GetManufacturer which is not in map",
			args: args{
				manufacturer: "amd",
			},
			want: "amd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeManufacturerName(tt.args.manufacturer); got != tt.want {
				t.Errorf("GetManufacturerFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
