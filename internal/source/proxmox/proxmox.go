package proxmox

import (
	"context"
	"fmt"
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
	Cluster         *proxmox.Cluster
	Nodes           []*proxmox.Node
	NodeIfaces      map[string][]*proxmox.NodeNetwork        // NodeName -> NodeNetworks (interfaces)
	Vms             map[string][]*proxmox.VirtualMachine     // NodeName -> VirtualMachines
	VMIfaces        map[string][]*proxmox.AgentNetworkIface  // VMName -> NetworkDevices
	Containers      map[string][]*proxmox.Container          // NodeName -> Contatiners
	ContainerIfaces map[string][]*proxmox.ContainerInterface // ContainerName -> ContainerInterfaces

	// Netbox related data for easier access. Initialized in sync functions.
	NetboxCluster *objects.Cluster
	NetboxNodes   map[string]*objects.Device // NodeName -> netbox device
}

// Function that collects all data from Proxmox API and stores it in ProxmoxSource struct.
func (ps *ProxmoxSource) Init() error {
	// Setup credentials for proxmox
	credentials := proxmox.Credentials{
		Username: ps.SourceConfig.Username,
		Password: ps.SourceConfig.Password,
	}

	// Create http client depending on ssl configuration
	HTTPClient, err := utils.NewHTTPClient(ps.SourceConfig.ValidateCert, ps.SourceConfig.CAFile)
	if err != nil {
		return fmt.Errorf("error creating new HTTP client: %s", err)
	}

	// Initialize proxmox client
	client := proxmox.NewClient(fmt.Sprintf("%s://%s:%d/api2/json",
		ps.SourceConfig.HTTPScheme, ps.SourceConfig.Hostname, ps.SourceConfig.Port),
		proxmox.WithCredentials(&credentials),
		proxmox.WithHTTPClient(HTTPClient),
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
		ps.Logger.Infof(ps.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"), duration.Seconds())
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
		ps.Logger.Infof(ps.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"), duration.Seconds())
	}
	return nil
}
