package ovirt

import (
	"fmt"
	"strings"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// OVirtSource represents an oVirt source.
//
//nolint:revive
type OVirtSource struct {
	common.Config
	Disks       map[string]*ovirtsdk4.Disk
	DataCenters map[string]*ovirtsdk4.DataCenter
	Clusters    map[string]*ovirtsdk4.Cluster
	Hosts       map[string]*ovirtsdk4.Host
	Vms         map[string]*ovirtsdk4.Vm
	Networks    *NetworkData
}

type NetworkData struct {
	OVirtNetworks       map[string]*ovirtsdk4.Network
	VnicProfile2Network map[string]string // vnicProfileId -> networkId
	Vid2Name            map[int]string
}

// Function that initializes state from ovirt api to local storage.
func (o *OVirtSource) Init() error {
	// Build the connection
	o.Logger.Debug(o.Ctx, "Initializing oVirt source ", o.SourceConfig.Name)
	connBuilder := ovirtsdk4.NewConnectionBuilder().
		URL(fmt.Sprintf(
			"%s://%s:%d/ovirt-engine/api",
			o.SourceConfig.HTTPScheme,
			o.SourceConfig.Hostname,
			o.SourceConfig.Port,
		),
		).
		Username(o.SourceConfig.Username).
		Password(o.SourceConfig.Password).
		Insecure(!o.SourceConfig.ValidateCert).
		Compress(true).
		Timeout(time.Second * constants.DefaultAPITimeout).
		CAFile(o.Config.CAFile)

	if o.Config.CAFile != "" {
		connBuilder.CAFile(o.Config.CAFile)
	}

	// Initialize the connection
	conn, err := connBuilder.Build()
	if err != nil {
		return fmt.Errorf("failed to create oVirt connection: %v", err)
	}
	defer conn.Close()

	// Initialize items to local storage
	initFunctions := []func(*ovirtsdk4.Connection) error{
		o.initNetworks,
		o.initDisks,
		o.initDataCenters,
		o.initClusters,
		o.initHosts,
		o.initVms,
	}

	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(conn); err != nil {
			return fmt.Errorf(
				"failed to initialize oVirt %s: %v",
				strings.TrimPrefix(fmt.Sprintf("%T", initFunc), "*source.OVirtSource.Init"),
				err,
			)
		}
		duration := time.Since(startTime)
		o.Logger.Infof(
			o.Ctx,
			"Successfully initialized %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"),
			duration.Seconds(),
		)
	}
	return nil
}

// Function that syncs all data from oVirt to Netbox.
func (o *OVirtSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		o.syncNetworks,
		o.syncDatacenters,
		o.syncClusters,
		o.syncHosts,
		o.syncVMs,
	}
	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		o.Logger.Infof(
			o.Ctx,
			"Successfully synced %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"),
			duration.Seconds(),
		)
	}
	return nil
}
