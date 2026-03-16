package constants

import (
	"testing"
)

func TestSourceTagColorMapCompleteness(t *testing.T) {
	allSources := []SourceType{Ovirt, Vmware, Dnac, PaloAlto, Fortigate, FMC, IOSXE}
	for _, src := range allSources {
		if _, ok := SourceTagColorMap[src]; !ok {
			t.Errorf("SourceTagColorMap missing entry for source %q", src)
		}
	}
}

func TestSourceTypeTagColorMapCompleteness(t *testing.T) {
	allSources := []SourceType{Ovirt, Vmware, Dnac, PaloAlto, Fortigate, FMC, IOSXE}
	for _, src := range allSources {
		if _, ok := SourceTypeTagColorMap[src]; !ok {
			t.Errorf("SourceTypeTagColorMap missing entry for source %q", src)
		}
	}
}

func TestSourceTagColorMapAndSourceTypeTagColorMapSameKeys(t *testing.T) {
	for src := range SourceTagColorMap {
		if _, ok := SourceTypeTagColorMap[src]; !ok {
			t.Errorf("SourceTypeTagColorMap missing key %q present in SourceTagColorMap", src)
		}
	}
	for src := range SourceTypeTagColorMap {
		if _, ok := SourceTagColorMap[src]; !ok {
			t.Errorf("SourceTagColorMap missing key %q present in SourceTypeTagColorMap", src)
		}
	}
}

func TestColorConstantsAreValidHex(t *testing.T) {
	colors := map[string]string{
		"ColorDarkRed":   ColorDarkRed,
		"ColorRed":       ColorRed,
		"ColorGreen":     ColorGreen,
		"ColorBlue":      ColorBlue,
		"ColorGrey":      ColorGrey,
		"ColorBlack":     ColorBlack,
		"ColorWhite":     ColorWhite,
		"SsotTagColor":   SsotTagColor,
		"OrphanTagColor": OrphanTagColor,
	}
	hexChars := "0123456789abcdef"
	for name, color := range colors {
		if len(color) != 6 {
			t.Errorf("%s has length %d, want 6", name, len(color))
		}
		for _, c := range color {
			isHex := false
			for _, h := range hexChars {
				if c == h {
					isHex = true
					break
				}
			}
			if !isHex {
				t.Errorf("%s contains invalid hex character %c", name, c)
			}
		}
	}
}

func TestArch2BitCompleteness(t *testing.T) {
	expectedArchs := []string{"x86_64", "aarch64", "arm64", "arm", "unknown"}
	for _, arch := range expectedArchs {
		if _, ok := Arch2Bit[arch]; !ok {
			t.Errorf("Arch2Bit missing entry for architecture %q", arch)
		}
	}
}
