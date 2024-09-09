package utils

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func TestExtractCPUArch(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Extract cpu arch from BOOT_IMAGE",
			args: args{
				input: "BOOT_IMAGE=(hd0,msdos1)/vmlinuz-5.15.0-207.156.6.el8uek.x86_64 root=/dev/mapper/ol-root ro crashkernel=auto resume=/dev/mapper/ol-swap00 rd.lvm.lv=ol/root rd.lvm.lv=ol/swap00 rhgb quiet intel_iommu=on",
			},
			want: "x86_64",
		},
		{
			name: "Extract cpu arch from kernel version",
			args: args{
				input: "Linux version 4.15.0-1051-oem (buildd@lgw01-amd64-016) (gcc version 7.4.0 (Ubuntu 7.4.0-1ubuntu1~18.04.1)) #60-Ubuntu SMP Fri Sep 13 13:51:54 UTC 2019 (x86_64)",
			},
			want: "x86_64",
		},
		{
			name: "Extract cpu arch from arch",
			args: args{
				input: "arch=aarch64 kernel_version=5.4.0-42-generic",
			},
			want: "aarch64",
		},
		{
			name: "No arch in string",
			args: args{
				input: "This is a test string with missing architecture",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractCPUArch(tt.args.input); got != tt.want {
				t.Errorf("ExtractCPUArch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePlatformName(t *testing.T) {
	type args struct {
		osDistribution string
		osMajorVersion string
		cpuArch        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test generation of platform name with empty parameters",
			args: args{
				osDistribution: "",
				osMajorVersion: "",
				cpuArch:        "",
			},
			want: constants.DefaultOSName,
		},
		{
			name: "Test generation of platform name with all parameters",
			args: args{
				osDistribution: "Linux",
				osMajorVersion: "8",
				cpuArch:        "x86_64",
			},
			want: "Linux 8 (64-bit)",
		},
		{
			name: "Test generation of platform name with missing osMajorVersion",
			args: args{
				osDistribution: "Linux",
				osMajorVersion: "",
				cpuArch:        "x86_64",
			},
			want: "Linux (64-bit)",
		},
		{
			name: "Test generation of platform name with missing cpuArch",
			args: args{
				osDistribution: "Linux",
				osMajorVersion: "8",
				cpuArch:        "",
			},
			want: "Linux 8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePlatformName(tt.args.osDistribution, tt.args.osMajorVersion, tt.args.cpuArch); got != tt.want {
				t.Errorf("GeneratePlatformName() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGenerateDeviceTypeSlug(t *testing.T) {
	type args struct {
		manufacturerName string
		modelName        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test generation of device type slug",
			args: args{
				manufacturerName: "Dell",
				modelName:        "PowerEdge R640",
			},
			want: "dell-poweredge-r640",
		},
		{
			name: "Test generation of device type slug with special characters",
			args: args{
				manufacturerName: "HPE",
				modelName:        "ProLiant DL380 Gen10",
			},
			want: "hpe-proliant-dl380-gen10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateDeviceTypeSlug(tt.args.manufacturerName, tt.args.modelName); got != tt.want {
				t.Errorf("GenerateDeviceTypeSlug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeManufacturerName(t *testing.T) {
	type args struct {
		manufacturer string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Serialize Cisco",
			args: args{
				manufacturer: "Cisco Systems",
			},
			want: "Cisco",
		},
		{
			name: "Serialize Dell",
			args: args{
				manufacturer: "Dell Inc.",
			},
			want: "Dell",
		},
		{
			name: "Serialize Fujitsu",
			args: args{
				manufacturer: "FTS Corp",
			},
			want: "Fujitsu",
		},
		{
			name: "No serialization needed",
			args: args{
				manufacturer: "Unknown Manufacturer",
			},
			want: "Unknown Manufacturer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeManufacturerName(tt.args.manufacturer); got != tt.want {
				t.Errorf("SerializeManufacturerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeOSName(t *testing.T) {
	type args struct {
		os string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "windows_2022",
			args: args{
				os: "windows_2022",
			},
			want: "Microsoft Windows 2022",
		},
		{
			name: "Redhat CoreOS",
			args: args{
				os: "Red Hat Enterprise Linux CoreOS",
			},
			want: "RHCOS",
		},
		{
			name: "Serialize Microsoft Windows Server",
			args: args{
				os: "Microsoft Windows Server",
			},
			want: "Microsoft Windows Server",
		},
		{
			name: "Serialize Windows",
			args: args{
				os: "Windows 10 Pro",
			},
			want: "Microsoft Windows",
		},
		{
			name: "Serialize Oracle Linux Server",
			args: args{
				os: "Oracle Linux Server 7.9",
			},
			want: "Oracle Linux",
		},
		{
			name: "Serialize Centos",
			args: args{
				os: "Centos 7",
			},
			want: "Centos Linux",
		},
		{
			name: "No serialization needed",
			args: args{
				os: "Custom OS",
			},
			want: "Custom OS",
		},
		{
			name: "Don't serialize default OS name",
			args: args{
				os: constants.DefaultOSName,
			},
			want: constants.DefaultOSName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeOSName(tt.args.os); got != tt.want {
				t.Errorf("SerializeOSName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCPUArchToBits(t *testing.T) {
	type args struct {
		arch string
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
			if got := CPUArchToBits(tt.args.arch); got != tt.want {
				t.Errorf("CPUArchToBits() = %v, want %v", got, tt.want)
			}
		})
	}
}
