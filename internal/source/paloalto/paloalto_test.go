package paloalto

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestPaloAltoLinkDuplexToNetbox(t *testing.T) {
	tests := []struct {
		name   string
		duplex string
		want   *objects.InterfaceDuplex
	}{
		{name: "full", duplex: "full", want: &objects.DuplexFull},
		{name: "auto", duplex: "auto", want: &objects.DuplexAuto},
		{name: "half", duplex: "half", want: &objects.DuplexHalf},
		{name: "empty", duplex: "", want: nil},
		{name: "unknown", duplex: "quarter", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := paloAltoLinkDuplexToNetbox(tt.duplex)
			switch {
			case tt.want == nil && got != nil:
				t.Errorf("paloAltoLinkDuplexToNetbox(%q) = %v, want nil", tt.duplex, *got)
			case tt.want != nil && got == nil:
				t.Errorf("paloAltoLinkDuplexToNetbox(%q) = nil, want %v", tt.duplex, *tt.want)
			case tt.want != nil && got != nil && *got != *tt.want:
				t.Errorf("paloAltoLinkDuplexToNetbox(%q) = %v, want %v", tt.duplex, *got, *tt.want)
			}
		})
	}
}
