package proxmox

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/luthermonson/go-proxmox"
)

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

func TestSyncVMsReturnsErrorWhenFailedVMsOutnumberNodes(t *testing.T) {
	testLogger, err := logger.New("", logger.ERROR)
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
	nbi := inventory.NewNetboxInventory(ctx, testLogger, &parser.NetboxConfig{}, false)
	nbi.NetboxAPI = service.FailingMockNetboxClient

	vms := make([]*proxmox.VirtualMachine, 0, 3)
	for i := range 3 {
		vms = append(vms, &proxmox.VirtualMachine{
			Name:                 fmt.Sprintf("vm%d", i),
			Status:               "stopped",
			VirtualMachineConfig: &proxmox.VirtualMachineConfig{},
		})
	}
	ps := &ProxmoxSource{
		Config: common.Config{
			Logger:       testLogger,
			Ctx:          ctx,
			SourceConfig: &parser.SourceConfig{Name: "test"},
		},
		Vms: map[string][]*proxmox.VirtualMachine{"node1": vms},
	}

	done := make(chan error, 1)
	go func() { done <- ps.syncVMs(nbi) }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("syncVMs() error = nil, want platform creation error")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("syncVMs() did not return: error fan-out deadlocked")
	}
}
