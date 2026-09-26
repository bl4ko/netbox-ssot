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
		{name: "windows 10", osType: "win10", want: "Windows 10"},
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

func TestParseDiskSizeMiB(t *testing.T) {
	tests := []struct {
		name string
		item string
		want int
	}{
		{name: "256G disk", item: "size=256G", want: 262144},
		{name: "32G disk", item: "size=32G", want: 32768},
		{name: "1T disk", item: "size=1T", want: 1048576},
		{name: "not a size token", item: "backup=0", want: 0},
		{name: "missing equals sign", item: "sizeG", want: 0},
		{name: "unsupported suffix", item: "size=512M", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDiskSizeMiB(tt.item); got != tt.want {
				t.Errorf("parseDiskSizeMiB(%q) = %d, want %d", tt.item, got, tt.want)
			}
		})
	}
}
