package proxmox

import "testing"

func TestProxmoxOSTypeToPlatformName(t *testing.T) {
	tests := []struct {
		name   string
		osType string
		want   string
	}{
		{name: "linux 2.6 kernel", osType: "l26", want: "Other 2.6.x Linux (64-bit)"},
		{name: "windows 11", osType: "win11", want: "Windows 11"},
		{name: "unknown type", osType: "win10", want: ""},
		{name: "empty type", osType: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := proxmoxOSTypeToPlatformName(tt.osType); got != tt.want {
				t.Errorf("proxmoxOSTypeToPlatformName(%q) = %q, want %q", tt.osType, got, tt.want)
			}
		})
	}
}
