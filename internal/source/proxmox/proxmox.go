package proxmox

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/luthermonson/go-proxmox"
)

//nolint:revive
type ProxmoxSource struct {
	common.Config

	// Proxmox API data initialized in init functions
	Cluster      *proxmox.Cluster
	Nodes        []*proxmox.Node
	NodeNetworks map[string][]*proxmox.NodeNetwork       // NodeName -> NodeNetworks (interfaces)
	Vms          map[string][]*proxmox.VirtualMachine    // NodeName -> VirtualMachines
	VMNetworks   map[string][]*proxmox.AgentNetworkIface // VMName -> NetworkDevices
	Containers   map[string][]*proxmox.Container         // NodeName -> Contatiners

	// Netbox related data for easier access. Initialized in sync functions.
	NetboxCluster *objects.Cluster
	NetboxNodes   map[string]*objects.Device // NodeName -> netbox device

	// Regex relations for matching objets
	ClusterSiteRelations   map[string]string
	ClusterTenantRelations map[string]string
	HostTenantRelations    map[string]string
	HostSiteRelations      map[string]string
	VMTenantRelations      map[string]string
	VlanGroupRelations     map[string]string
	VlanTenantRelations    map[string]string
}

// Function that collects all data from Proxmox API and stores it in ProxmoxSource struct.
func (ps *ProxmoxSource) Init() error {
	// Initialize regex relations
	ps.Logger.Debug(ps.Ctx, "Initializing regex relations for oVirt source ", ps.SourceConfig.Name)
	ps.HostSiteRelations = utils.ConvertStringsToRegexPairs(ps.SourceConfig.HostSiteRelations)
	ps.Logger.Debug(ps.Ctx, "HostSiteRelations: ", ps.HostSiteRelations)
	ps.ClusterSiteRelations = utils.ConvertStringsToRegexPairs(ps.SourceConfig.ClusterSiteRelations)
	ps.Logger.Debug(ps.Ctx, "ClusterSiteRelations: ", ps.ClusterSiteRelations)
	ps.ClusterTenantRelations = utils.ConvertStringsToRegexPairs(ps.SourceConfig.ClusterTenantRelations)
	ps.Logger.Debug(ps.Ctx, "ClusterTenantRelations: ", ps.ClusterTenantRelations)
	ps.HostTenantRelations = utils.ConvertStringsToRegexPairs(ps.SourceConfig.HostTenantRelations)
	ps.Logger.Debug(ps.Ctx, "HostTenantRelations: ", ps.HostTenantRelations)
	ps.VMTenantRelations = utils.ConvertStringsToRegexPairs(ps.SourceConfig.VMTenantRelations)
	ps.Logger.Debug(ps.Ctx, "VmTenantRelations: ", ps.VMTenantRelations)
	ps.VlanGroupRelations = utils.ConvertStringsToRegexPairs(ps.SourceConfig.VlanGroupRelations)
	ps.Logger.Debug(ps.Ctx, "VlanGroupRelations: ", ps.VlanGroupRelations)
	ps.VlanTenantRelations = utils.ConvertStringsToRegexPairs(ps.SourceConfig.VlanTenantRelations)
	ps.Logger.Debug(ps.Ctx, "VlanTenantRelations: ", ps.VlanTenantRelations)

	// Initialize the connection
	credentials := proxmox.Credentials{
		Username: ps.SourceConfig.Username,
		Password: ps.SourceConfig.Password,
	}
	HTTPClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !ps.SourceConfig.ValidateCert,
			},
		},
	}
	client := proxmox.NewClient(fmt.Sprintf("%s://%s:%d/api2/json",
		ps.SourceConfig.HTTPScheme, ps.SourceConfig.Hostname, ps.SourceConfig.Port),
		proxmox.WithCredentials(&credentials),
		proxmox.WithHTTPClient(&HTTPClient),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initFuncs := []func(context.Context, *proxmox.Client) error{
		ps.initCluster,
		ps.initNodes,
	}

	for _, initFunc := range initFuncs {
		startTime := time.Now()
		if err := initFunc(ctx, client); err != nil {
			return fmt.Errorf("proxmox initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		ps.Logger.Infof(ps.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionName(initFunc), duration.Seconds())
	}

	return nil
}

// Function that syncs all collected data to Netbox inventory.
func (ps *ProxmoxSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		ps.syncCluster,
		ps.syncNodes,
		ps.syncVMs,
		ps.syncContainers,
	}
	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		ps.Logger.Infof(ps.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionName(syncFunc), duration.Seconds())
	}
	return nil
}
