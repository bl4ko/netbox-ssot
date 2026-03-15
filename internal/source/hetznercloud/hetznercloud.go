package hetznercloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/src-doo/netbox-ssot/internal/netbox/inventory"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
	"github.com/src-doo/netbox-ssot/internal/source/common"
	"github.com/src-doo/netbox-ssot/internal/utils"
)

type HetznerCloudSource struct {
	common.Config

	// Hetzner Cloud API data initialized in init functions
	Locations   []*hcloud.Location
	Datacenters []*hcloud.Datacenter
	Servers     []*hcloud.Server
	Networks    []*hcloud.Network
	FloatingIPs []*hcloud.FloatingIP
	PrimaryIPs  []*hcloud.PrimaryIP

	// Netbox related data for easier access
	NetboxSites     map[string]*objects.Site     // Location City -> Netbox Site
	NetboxLocations map[string]*objects.Location // Datacenter Name -> Netbox Location
}

func (hcs *HetznerCloudSource) Init() error {
	opts := []hcloud.ClientOption{
		hcloud.WithToken(hcs.SourceConfig.APIToken),
	}
	
	// Create http client depending on ssl configuration
	HTTPClient, err := utils.NewHTTPClient(hcs.SourceConfig.ValidateCert, hcs.SourceConfig.CAFile)
	if err != nil {
		return fmt.Errorf("error creating new HTTP client: %s", err)
	}
	opts = append(opts, hcloud.WithHTTPClient(HTTPClient))

	client := hcloud.NewClient(opts...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initFuncs := []func(context.Context, *hcloud.Client) error{
		hcs.initLocations,
		hcs.initDatacenters,
		hcs.initServers,
		hcs.initNetworks,
		hcs.initFloatingIPs,
		hcs.initPrimaryIPs,
	}

	for _, initFunc := range initFuncs {
		startTime := time.Now()
		if err := initFunc(ctx, client); err != nil {
			return fmt.Errorf("hetznercloud initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		hcs.Logger.Infof(
			hcs.Ctx,
			"Successfully initialized %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"),
			duration.Seconds(),
		)
	}

	return nil
}

func (hcs *HetznerCloudSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		hcs.syncLocationsAndDatacenters,
		hcs.syncServers,
		hcs.syncNetworks,
		hcs.syncFloatingIPs,
	}
	var encounteredErrors []error
	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		funcName := utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync")
		err := syncFunc(nbi)
		if err != nil {
			if hcs.SourceConfig.ContinueOnError {
				hcs.Logger.Errorf(
					hcs.Ctx,
					"Error syncing %s: %s (continuing due to continueOnError flag)",
					funcName,
					err,
				)
				encounteredErrors = append(encounteredErrors, fmt.Errorf("%s: %w", funcName, err))
			} else {
				return err
			}
		} else {
			duration := time.Since(startTime)
			hcs.Logger.Infof(
				hcs.Ctx,
				"Successfully synced %s in %f seconds",
				funcName,
				duration.Seconds(),
			)
		}
	}
	if len(encounteredErrors) > 0 {
		return fmt.Errorf("encountered %d errors during sync: %v", len(encounteredErrors), encounteredErrors)
	}
	return nil
}
