package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
)

func TestCleanPlatformName(t *testing.T) {
	oss := &Source{}
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Should insert space between letters and numbers",
			input:    "almalinux9",
			expected: "Almalinux 9",
		},
		{
			name:     "Should capitalize without numbers",
			input:    "ubuntu",
			expected: "Ubuntu",
		},
		{
			name:     "Should format correctly with capitals in input",
			input:    "CentOS7",
			expected: "Centos 7",
		},
		{
			name:     "Should handle empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := oss.cleanPlatformName(tc.input)
			if result != tc.expected {
				t.Errorf("cleanPlatformName(%q) = %q; expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestFindImageNameByID(t *testing.T) {
	oss := &Source{
		Images: []images.Image{
			{ID: "image1-id", Name: "Ubuntu 22.04"},
			{ID: "image2-id", Name: "Debian 11"},
		},
	}

	tests := []struct {
		name     string
		imageID  string
		expected string
	}{
		{
			name:     "Valid existing image ID",
			imageID:  "image1-id",
			expected: "Ubuntu 22.04",
		},
		{
			name:     "Another valid existing image",
			imageID:  "image2-id",
			expected: "Debian 11",
		},
		{
			name:     "Image ID not found",
			imageID:  "nonexistent-id",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := oss.findImageNameByID(tc.imageID)
			if result != tc.expected {
				t.Errorf("findImageNameByID(%q) = %q; expected %q", tc.imageID, result, tc.expected)
			}
		})
	}
}
