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
