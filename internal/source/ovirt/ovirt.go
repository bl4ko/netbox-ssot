package ovirt

import (
	"fmt"
	"strings"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// OVirtSource represents an oVirt source
type OVirtSource struct {
	common.CommonConfig
	Disks       map[string]*ovirtsdk4.Disk
	DataCenters map[string]*ovirtsdk4.DataCenter
	Clusters    map[string]*ovirtsdk4.Cluster
	Hosts       map[string]*ovirtsdk4.Host
	Vms         map[string]*ovirtsdk4.Vm
	Networks    *NetworkData

	HostSiteRelations      map[string]string
	ClusterSiteRelations   map[string]string
	ClusterTenantRelations map[string]string
	HostTenantRelations    map[string]string
	VmTenantRelations      map[string]string
	VlanGroupRelations     map[string]string
	VlanTenantRelations    map[string]string
}

type NetworkData struct {
	OVirtNetworks map[string]*ovirtsdk4.Network
	Vid2Name      map[int]string
}

// Function that initializes state from ovirt api to local storage
func (o *OVirtSource) Init() error {
	// Initialize regex relations
	o.Logger.Debug("Initializing regex relations for oVirt source ", o.SourceConfig.Name)
	o.HostSiteRelations = utils.ConvertStringsToRegexPairs(o.SourceConfig.HostSiteRelations)
	o.Logger.Debug("HostSiteRelations: ", o.HostSiteRelations)
	o.ClusterSiteRelations = utils.ConvertStringsToRegexPairs(o.SourceConfig.ClusterSiteRelations)
	o.Logger.Debug("ClusterSiteRelations: ", o.ClusterSiteRelations)
	o.ClusterTenantRelations = utils.ConvertStringsToRegexPairs(o.SourceConfig.ClusterTenantRelations)
	o.Logger.Debug("ClusterTenantRelations: ", o.ClusterTenantRelations)
	o.HostTenantRelations = utils.ConvertStringsToRegexPairs(o.SourceConfig.HostTenantRelations)
	o.Logger.Debug("HostTenantRelations: ", o.HostTenantRelations)
	o.VmTenantRelations = utils.ConvertStringsToRegexPairs(o.SourceConfig.VmTenantRelations)
	o.Logger.Debug("VmTenantRelations: ", o.VmTenantRelations)
	o.VlanGroupRelations = utils.ConvertStringsToRegexPairs(o.SourceConfig.VlanGroupRelations)
	o.Logger.Debug("VlanGroupRelations: ", o.VlanGroupRelations)
	o.VlanTenantRelations = utils.ConvertStringsToRegexPairs(o.SourceConfig.VlanTenantRelations)
	o.Logger.Debug("VlanTenantRelations: ", o.VlanTenantRelations)

	// Initialize the connection
	o.Logger.Debug("Initializing oVirt source ", o.SourceConfig.Name)
	conn, err := ovirtsdk4.NewConnectionBuilder().
		URL(fmt.Sprintf("%s://%s:%d/ovirt-engine/api", o.SourceConfig.HTTPScheme, o.SourceConfig.Hostname, o.SourceConfig.Port)).
		Username(o.SourceConfig.Username).
		Password(o.SourceConfig.Password).
		Insecure(!o.SourceConfig.ValidateCert).
		Compress(true).
		Timeout(time.Second * 10).
		Build()
	if err != nil {
		return fmt.Errorf("failed to create oVirt connection: %v", err)
	}
	defer conn.Close()

	// Initialize items to local storage
	initFunctions := []func(*ovirtsdk4.Connection) error{
		o.InitNetworks,
		o.InitDisks,
		o.InitDataCenters,
		o.InitClusters,
		o.InitHosts,
		o.InitVms,
	}

	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(conn); err != nil {
			return fmt.Errorf("failed to initialize oVirt %s: %v", strings.TrimPrefix(fmt.Sprintf("%T", initFunc), "*source.OVirtSource.Init"), err)
		}
		duration := time.Since(startTime)
		o.Logger.Infof("Successfully initialized %s in %f seconds", utils.ExtractFunctionName(initFunc), duration.Seconds())
	}

	return nil
}

// Function that syncs all data from oVirt to Netbox
func (o *OVirtSource) Sync(nbi *inventory.NetBoxInventory) error {
	syncFunctions := []func(*inventory.NetBoxInventory) error{
		o.syncNetworks,
		o.syncDatacenters,
		o.syncClusters,
		o.syncHosts,
		o.syncVms,
	}
	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		o.Logger.Infof("Successfully synced %s in %f seconds", utils.ExtractFunctionName(syncFunc), duration.Seconds())
	}
	return nil
}
