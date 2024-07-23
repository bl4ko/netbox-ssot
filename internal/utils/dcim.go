package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

// extractCPUArch extracts the CPU architecture from an input string.
// If no CPU architecture was found, empty string is returned.
func ExtractCPUArch(input string) string {
	// Define a regular expression pattern to match common CPU architectures
	re := regexp.MustCompile(`(?:x86_64|i[3-6]86|aarch64|arm64|ppc64le|s390x|mips64|riscv64)`)

	// Find the first match
	match := re.FindString(input)

	return match
}

// CPUArchToBits maps cpu architecture to corresponding bits of the architecture:
// x86_64 -> 64-bit, arm -> 32-bit....
// CPUArchToBits maps cpu architecture to corresponding bits of the architecture.
func CPUArchToBits(arch string) string {
	archMap := map[string]string{
		"x86_64":  "64-bit",
		"amd64":   "64-bit",
		"i386":    "32-bit",
		"i486":    "32-bit",
		"i586":    "32-bit",
		"i686":    "32-bit",
		"aarch64": "64-bit",
		"arm64":   "64-bit",
		"arm":     "32-bit",
		"arm32":   "32-bit",
		"ppc64le": "64-bit",
		"s390x":   "64-bit",
		"mips64":  "64-bit",
		"riscv64": "64-bit",
	}

	if bits, exists := archMap[arch]; exists {
		return bits
	}
	return arch
}

// Function that takes osDistribution (Linux, Windows, ...), osMajorVersion (8, 9, 10, ...)
// and cpuArch (x86_64, or 64bit, ....) and generates universal platform name in the format of
// "osDistrbution osMajorVersion (cpuArch)".
func GeneratePlatformName(osDistribution string, osMajorVersion string, cpuArch string) string {
	if osDistribution != "" {
		osDistribution = SerializeOSName(osDistribution)
	} else {
		return constants.DefaultOSName
	}
	if osMajorVersion != "" {
		if !strings.Contains(osDistribution, "(") {
			osMajorVersion = fmt.Sprintf(" %s", osMajorVersion)
		} else {
			osMajorVersion = ""
		}
	}
	if cpuArch != "" {
		// Check if cpuArch was extreacted from osDistribution
		if !strings.Contains(osDistribution, "(") {
			cpuArch = fmt.Sprintf(" (%s)", CPUArchToBits(cpuArch))
		} else {
			cpuArch = ""
		}
	}
	return fmt.Sprintf("%s%s%s", osDistribution, osMajorVersion, cpuArch)
}

// GenerateDeviceTypeSlug generates a device type slug from the given manufacturer and model.
func GenerateDeviceTypeSlug(manufacturerName string, modelName string) string {
	manufacturerSlug := Slugify(manufacturerName)
	modelSlug := Slugify(modelName)
	return fmt.Sprintf("%s-%s", manufacturerSlug, modelSlug)
}

// ManufacturerMap maps regex of manufacturer names to manufacturer name.
// Manufacturer names are compatible with device type library. See
// internal/devices/combined_data.go for more info.
var ManufacturerMap = map[string]string{
	".*Cisco.*":    "Cisco",
	".*Fortinet.*": "Fortinet",
	".*Dell.*":     "Dell",
	"FTS Corp":     "Fujitsu",
	".*Fujitsu.*":  "Fujitsu",
	"^HP$":         "HPE",
	"^HP .*":       "HPE",
	".*Huawei.*":   "Huawei",
	".*Inspur.*":   "Inspur",
	".*Intel.*":    "Intel",
	"LEN":          "Lenovo",
	".*Nvidea.*":   "Nvidia",
	".*Samsung.*":  "Samsung",
}

// GetManufactuerFromString returns manufacturer name from the given string.
func SerializeManufacturerName(manufacturer string) string {
	for regex, name := range ManufacturerMap {
		matched, _ := regexp.MatchString(regex, manufacturer)
		if matched {
			return name
		}
	}
	return manufacturer
}

// SpecificOSMap maps regex of OS names to serialized OS names.
var SpecificOSMap = map[string]string{
	"rhcos_x64":                           "RHCOS (64bit)",
	".*Red Hat Enterprise Linux CoreOS.*": "RHCOS",

	".*windows_2022.*": "Microsoft Windows 2022",

	".*ol_5x64.*":              "Oracle Linux 5 (64-bit)",
	".*ol_6x64.*":              "Oracle Linux 6 (64-bit)",
	".*ol_7x64.*":              "Oracle Linux 7 (64-bit)",
	".*ol_8x64.*":              "Oracle Linux 8 (64-bit)",
	".*ol_9x64.*":              "Oracle Linux 9 (64-bit)",
	"Microsoft Windows Server": "Microsoft Windows Server",
}

// Universal OSMap maps regex of OS names to serialized OS names.
var UniversalOSMap = map[string]string{
	".*Red Hat Enterprise Linux.*": "RHEL",

	".*Windows.*": "Microsoft Windows",

	".*ol_.*":                 "Oracle Linux",
	"^Oracle$":                "Oracle Linux",
	".*Oracle Linux Server.*": "Oracle Linux",

	".*Centos.*": "Centos Linux",

	".*Rocky.*": "Rocky Linux",

	".*Alma.*": "Alma Linux",

	".*Ubuntu.*": "Ubuntu Linux",
}

// SerializeOSName returns serialized OS name from the given string.
func SerializeOSName(os string) string {
	for regex, name := range SpecificOSMap {
		matched, _ := regexp.MatchString(regex, os)
		if matched {
			return name
		}
	}
	for regex, name := range UniversalOSMap {
		matched, _ := regexp.MatchString(regex, os)
		if matched {
			return name
		}
	}
	return os
}
