package proxmox

import (
	"context"
	"testing"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/luthermonson/go-proxmox"
)

// newTestSource returns a ProxmoxSource wired to a site-scoped NetBox cluster,
// ready to be synced against inventory.MockInventory.
func newTestSource(t *testing.T, sourceConfig *parser.SourceConfig) *ProxmoxSource {
	t.Helper()
	lg, err := logger.New("", logger.ERROR)
	if err != nil {
		t.Fatal(err)
	}
	if sourceConfig.Name == "" {
		sourceConfig.Name = "test-proxmox"
	}
	return &ProxmoxSource{
		Config: common.Config{
			Logger:        lg,
			Ctx:           context.WithValue(context.Background(), constants.CtxSourceKey, sourceConfig.Name),
			SourceConfig:  sourceConfig,
			SourceNameTag: &objects.Tag{ID: 9001, Name: "Source: " + sourceConfig.Name, Slug: "source-test-proxmox"},
			SourceTypeTag: &objects.Tag{ID: 9002, Name: "proxmox", Slug: "proxmox"},
		},
		NetboxCluster: &objects.Cluster{
			NetboxObject: objects.NetboxObject{ID: 1},
			Name:         "test-cluster",
			ScopeType:    constants.ContentTypeDcimSite,
			ScopeID:      1,
		},
	}
}

func newTestVM(name string, vmid uint64) *proxmox.VirtualMachine {
	osType := "win11"
	return &proxmox.VirtualMachine{
		Name:                 name,
		VMID:                 proxmox.StringOrUint64(vmid),
		VirtualMachineConfig: &proxmox.VirtualMachineConfig{OSType: &osType},
	}
}

func TestSyncVMsReturnsWhenMoreVMsThanNodesFail(t *testing.T) {
	nbi := inventory.MockInventory
	saved := nbi.NetboxAPI
	nbi.NetboxAPI = service.FailingMockNetboxClient
	defer func() { nbi.NetboxAPI = saved }()

	ps := newTestSource(t, &parser.SourceConfig{})
	ps.NetboxNodes = map[string]*objects.Device{
		"n1": {
			NetboxObject: objects.NetboxObject{ID: 1},
			Name:         "n1",
			Site:         &objects.Site{NetboxObject: objects.NetboxObject{ID: 1}},
		},
	}
	// One node and two VMs whose platform creation fails against the failing client.
	ps.Vms = map[string][]*proxmox.VirtualMachine{
		"n1": {newTestVM("failing-vm1", 101), newTestVM("failing-vm2", 102)},
	}

	done := make(chan error, 1)
	go func() { done <- ps.syncVMs(nbi) }()
	select {
	case err := <-done:
		if err == nil {
			t.Errorf("syncVMs returned nil, want an error")
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("syncVMs did not return within 5s (1 node, 2 failing VMs)")
	}
}
