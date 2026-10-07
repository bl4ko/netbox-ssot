package vmware

import (
	"context"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// newTestSource returns a VmwareSource with empty vSphere data, ready to be synced
// against inventory.MockInventory in dry-run mode.
func newTestSource(t *testing.T, sourceConfig *parser.SourceConfig) *VmwareSource {
	t.Helper()
	lg, err := logger.New("", logger.ERROR)
	if err != nil {
		t.Fatal(err)
	}
	if sourceConfig.Name == "" {
		sourceConfig.Name = "test-vmware"
	}
	return &VmwareSource{
		Config: common.Config{
			Logger:        lg,
			Ctx:           context.WithValue(context.Background(), constants.CtxSourceKey, sourceConfig.Name),
			SourceConfig:  sourceConfig,
			SourceNameTag: &objects.Tag{ID: 9101, Name: "Source: " + sourceConfig.Name, Slug: "source-test-vmware"},
			SourceTypeTag: &objects.Tag{ID: 9102, Name: "vmware", Slug: "vmware"},
		},
		Networks: NetworkData{
			DistributedVirtualPortgroups: map[string]*DistributedPortgroupData{},
			Vid2Name:                     map[int]string{},
			HostVirtualSwitches:          map[string]map[string]*HostVirtualSwitchData{},
			HostProxySwitches:            map[string]map[string]*HostProxySwitchData{},
			HostPortgroups:               map[string]map[string]*HostPortgroupData{},
		},
	}
}

// newTestInventory returns the dry-run mock inventory, with the default site that
// NetboxInventory.Init creates in a real run.
func newTestInventory(t *testing.T) (*inventory.NetboxInventory, *objects.Site) {
	t.Helper()
	service.MockNetboxClient.DryRun = true
	nbi := inventory.MockInventory
	defaultSite, err := nbi.AddSite(context.Background(), &objects.Site{
		Name: constants.DefaultSite,
		Slug: "defaultsite",
	})
	if err != nil {
		t.Fatal(err)
	}
	return nbi, defaultSite
}

// newTestHost returns an ESXi host with the summary fields vCenter reports for a
// connected host.
func newTestHost(id, name string) mo.HostSystem {
	host := mo.HostSystem{
		Summary: types.HostListSummary{
			Hardware: &types.HostHardwareSummary{Vendor: "Dell Inc.", Model: "PowerEdge R650", NumCpuCores: 32},
			Runtime:  &types.HostRuntimeInfo{ConnectionState: types.HostSystemConnectionStateConnected},
			Config: types.HostConfigSummary{
				Product: &types.AboutInfo{Name: "VMware ESXi", Version: "8.0.3"},
			},
		},
	}
	host.Self = types.ManagedObjectReference{Type: "HostSystem", Value: id}
	host.Name = name
	return host
}

func TestSyncHostsWithoutHostSiteRelationsUsesDefaultSite(t *testing.T) {
	nbi, defaultSite := newTestInventory(t)
	vc := newTestSource(t, &parser.SourceConfig{})
	vc.Hosts = map[string]mo.HostSystem{"host-1": newTestHost("host-1", "esx-default-site")}

	if err := vc.syncHosts(nbi); err != nil {
		t.Fatalf("syncHosts() error = %v, want the host synced to the default site", err)
	}
	if _, ok := nbi.GetDevice("esx-default-site", defaultSite.ID); !ok {
		t.Errorf("host esx-default-site not found in site %s", constants.DefaultSite)
	}
}
