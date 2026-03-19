package source

import (
	"context"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/dnac"
	"github.com/bl4ko/netbox-ssot/internal/source/fmc"
	"github.com/bl4ko/netbox-ssot/internal/source/fortigate"
	iosxe "github.com/bl4ko/netbox-ssot/internal/source/ios-xe"
	"github.com/bl4ko/netbox-ssot/internal/source/ovirt"
	"github.com/bl4ko/netbox-ssot/internal/source/paloalto"
	"github.com/bl4ko/netbox-ssot/internal/source/proxmox"
	"github.com/bl4ko/netbox-ssot/internal/source/vmware"
)

func setupMockServer(t *testing.T) {
	t.Helper()
	mockServer := service.CreateMockServer()
	t.Cleanup(mockServer.Close)
	service.MockNetboxClient.BaseURL = mockServer.URL
}

func TestNewSource_AllTypes(t *testing.T) {
	setupMockServer(t)
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
	nbi := inventory.MockInventory

	tests := []struct {
		name       string
		sourceType constants.SourceType
	}{
		{name: "ovirt", sourceType: constants.Ovirt},
		{name: "vmware", sourceType: constants.Vmware},
		{name: "dnac", sourceType: constants.Dnac},
		{name: "proxmox", sourceType: constants.Proxmox},
		{name: "paloalto", sourceType: constants.PaloAlto},
		{name: "fortigate", sourceType: constants.Fortigate},
		{name: "fmc", sourceType: constants.FMC},
		{name: "ios-xe", sourceType: constants.IOSXE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &parser.SourceConfig{
				Name:     "test-" + tt.name,
				Type:     tt.sourceType,
				Tag:      "test-tag",
				TagColor: "00add8",
			}
			src, err := NewSource(ctx, config, nbi.Logger, nbi)
			if err != nil {
				t.Fatalf("NewSource(%s) returned error: %v", tt.sourceType, err)
			}
			if src == nil {
				t.Fatalf("NewSource(%s) returned nil", tt.sourceType)
			}

			switch tt.sourceType {
			case constants.Ovirt:
				if _, ok := src.(*ovirt.OVirtSource); !ok {
					t.Errorf("expected *ovirt.OVirtSource, got %T", src)
				}
			case constants.Vmware:
				if _, ok := src.(*vmware.VmwareSource); !ok {
					t.Errorf("expected *vmware.VmwareSource, got %T", src)
				}
			case constants.Dnac:
				if _, ok := src.(*dnac.DnacSource); !ok {
					t.Errorf("expected *dnac.DnacSource, got %T", src)
				}
			case constants.Proxmox:
				if _, ok := src.(*proxmox.ProxmoxSource); !ok {
					t.Errorf("expected *proxmox.ProxmoxSource, got %T", src)
				}
			case constants.PaloAlto:
				if _, ok := src.(*paloalto.PaloAltoSource); !ok {
					t.Errorf("expected *paloalto.PaloAltoSource, got %T", src)
				}
			case constants.Fortigate:
				if _, ok := src.(*fortigate.FortigateSource); !ok {
					t.Errorf("expected *fortigate.FortigateSource, got %T", src)
				}
			case constants.FMC:
				if _, ok := src.(*fmc.FMCSource); !ok {
					t.Errorf("expected *fmc.FMCSource, got %T", src)
				}
			case constants.IOSXE:
				if _, ok := src.(*iosxe.IOSXESource); !ok {
					t.Errorf("expected *iosxe.IOSXESource, got %T", src)
				}
			}
		})
	}
}

func TestNewSource_UnsupportedType(t *testing.T) {
	setupMockServer(t)
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
	nbi := inventory.MockInventory

	config := &parser.SourceConfig{
		Name:     "test-unknown",
		Type:     "nonexistent",
		Tag:      "test-tag",
		TagColor: "00add8",
	}
	_, err := NewSource(ctx, config, nbi.Logger, nbi)
	if err == nil {
		t.Fatal("expected error for unsupported source type, got nil")
	}
}
