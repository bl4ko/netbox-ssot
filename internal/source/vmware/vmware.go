package vmware

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

// VmwareSource represents an vsphere source
type VmwareSource struct {
	common.CommonConfig
	Disks       map[string]mo.Datastore
	DataCenters map[string]mo.Datacenter
	Clusters    map[string]mo.ClusterComputeResource
	Hosts       map[string]mo.HostSystem
	Vms         map[string]mo.VirtualMachine
	Networks    NetworkData

	// Relations between objects "object_id": "object_id"
	Cluster2Datacenter map[string]string // ClusterKey -> DatacenterKey
	Host2Cluster       map[string]string // HostKey -> ClusterKey
	Vm2Host            map[string]string // VmKey ->  HostKey

	// Netbox relations
	ClusterSiteRelations   map[string]string
	ClusterTenantRelations map[string]string
	HostTenantRelations    map[string]string
	HostSiteRelations      map[string]string
	VmTenantRelations      map[string]string
	VlanGroupRelations     map[string]string
	VlanTenantRelations    map[string]string
}

type NetworkData struct {
	DistributedVirtualPortgroups map[string]*DistributedPortgroupData         // Portgroup.key -> PortgroupData
	Vid2Name                     map[int]string                               // Helper map, for quickly obtaining name of the vid
	HostVirtualSwitches          map[string]map[string]*HostVirtualSwitchData // hostName -> VSwitchName-> VSwitchData
	HostProxySwitches            map[string]map[string]*HostProxySwitchData   // hostName -> PSwitchName ->
	HostPortgroups               map[string]map[string]*HostPortgroupData     // hostname -> Portgroup.Spec.Name -> HostPortgroupData
}

type DistributedPortgroupData struct {
	Name         string
	VlanIds      []int
	VlanIdRanges []string
	Private      bool
	Tenant       *objects.Tenant
}

type HostVirtualSwitchData struct {
	mtu   int
	pnics []string
}

type HostProxySwitchData struct {
	name  string
	mtu   int
	pnics []string
}

type HostPortgroupData struct {
	vlanId  int
	vswitch string
	nics    []string
}

func (vc *VmwareSource) Init() error {
	// Initialize regex relations
	vc.Logger.Debug("Initializing regex relations for oVirt source ", vc.SourceConfig.Name)
	vc.HostSiteRelations = utils.ConvertStringsToRegexPairs(vc.SourceConfig.HostSiteRelations)
	vc.Logger.Debug("HostSiteRelations: ", vc.HostSiteRelations)
	vc.ClusterSiteRelations = utils.ConvertStringsToRegexPairs(vc.SourceConfig.ClusterSiteRelations)
	vc.Logger.Debug("ClusterSiteRelations: ", vc.ClusterSiteRelations)
	vc.ClusterTenantRelations = utils.ConvertStringsToRegexPairs(vc.SourceConfig.ClusterTenantRelations)
	vc.Logger.Debug("ClusterTenantRelations: ", vc.ClusterTenantRelations)
	vc.HostTenantRelations = utils.ConvertStringsToRegexPairs(vc.SourceConfig.HostTenantRelations)
	vc.Logger.Debug("HostTenantRelations: ", vc.HostTenantRelations)
	vc.VmTenantRelations = utils.ConvertStringsToRegexPairs(vc.SourceConfig.VmTenantRelations)
	vc.Logger.Debug("VmTenantRelations: ", vc.VmTenantRelations)
	vc.VlanGroupRelations = utils.ConvertStringsToRegexPairs(vc.SourceConfig.VlanGroupRelations)
	vc.Logger.Debug("VlanGroupRelations: ", vc.VlanGroupRelations)
	vc.VlanTenantRelations = utils.ConvertStringsToRegexPairs(vc.SourceConfig.VlanTenantRelations)
	vc.Logger.Debug("VlanTenantRelations: ", vc.VlanTenantRelations)

	// Initialize the connection
	vc.Logger.Debug("Initializing oVirt source ", vc.SourceConfig.Name)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Correctly handle backslashes in username and password
	escapedUsername := url.PathEscape(vc.SourceConfig.Username)
	escapedPassword := url.PathEscape(vc.SourceConfig.Password)

	vcUrl := fmt.Sprintf("%s://%s:%s@%s:%d/sdk", vc.SourceConfig.HTTPScheme, escapedUsername, escapedPassword, vc.SourceConfig.Hostname, vc.SourceConfig.Port)

	url, err := url.Parse(vcUrl)
	if err != nil {
		return fmt.Errorf("failed parsing url for %s with error %s", vc.SourceConfig.Hostname, err)
	}

	conn, err := govmomi.NewClient(ctx, url, !vc.SourceConfig.ValidateCert)
	if err != nil {
		return fmt.Errorf("failed creating a govmomi client with an error: %s", err)
	}

	// View manager is used to create and manage views. Views are a mechanism in vSphere
	// to group and manage objects in the inventory.
	viewManager := view.NewManager(conn.Client)

	// viewType specifies the types of objects to be included in our container view.
	// Each string in this slice represents a different vSphere Managed Object type.
	viewType := []string{
		"Datastore", "Datacenter", "ClusterComputeResource", "HostSystem", "VirtualMachine", "Network",
	}

	// A container view is a subset of the vSphere inventory, focusing on the specified
	// object types, making it easier to manage and retrieve data for these objects.
	containerView, err := viewManager.CreateContainerView(ctx, conn.Client.ServiceContent.RootFolder, viewType, true)
	if err != nil {
		return fmt.Errorf("failed creating containerView: %s", err)
	}
	defer containerView.Destroy(ctx)

	vc.Logger.Debug("Connection to vmware source ", vc.SourceConfig.Hostname, " established successfully")

	// Find relation between datacenters and clusters. Currently we have to manually traverse
	// the tree to get this relation.
	vc.CreateClusterDataCenterRelation(ctx, conn.Client)

	// Initialize items to local storage
	initFunctions := []func(context.Context, *view.ContainerView) error{
		vc.InitNetworks,
		vc.InitDisks,
		vc.InitDataCenters,
		vc.InitClusters,
		vc.InitHosts,
		vc.InitVms,
	}

	for _, initFunc := range initFunctions {
		if err := initFunc(ctx, containerView); err != nil {
			return fmt.Errorf("vmware initialization failure: %v", err)
		}
	}

	vc.Logger.Debug("Successfully initialized objects from source ", vc.CommonConfig.SourceConfig.Name, ".")

	err = conn.Logout(ctx)
	if err != nil {
		return fmt.Errorf("error occurred when ending vmware connection to host %s: %s", vc.SourceConfig.Hostname, err)
	}

	vc.Logger.Debug("Successfully closed connection to vmware host: ", vc.SourceConfig.Hostname)

	return nil
}

// Function that syncs all data from oVirt to Netbox
func (vc *VmwareSource) Sync(nbi *inventory.NetBoxInventory) error {
	syncFunctions := []func(*inventory.NetBoxInventory) error{
		vc.syncNetworks,
		vc.syncDatacenters,
		vc.syncClusters,
		vc.syncHosts,
		vc.syncVms,
	}
	for _, syncFunc := range syncFunctions {
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
	}
	return nil
}

// Currently we have to traverse the vsphere tree to get datacenter to cluster relation
// For other objects relations are available in with containerView.
func (vc *VmwareSource) CreateClusterDataCenterRelation(ctx context.Context, client *vim25.Client) error {
	finder := find.NewFinder(client, true)
	datacenters, err := finder.DatacenterList(ctx, "*")
	if err != nil {
		return fmt.Errorf("finder failed creating datacenter list: %s", err)
	}
	vc.Cluster2Datacenter = make(map[string]string)
	for _, dc := range datacenters {
		finder.SetDatacenter(dc)
		clusters, err := finder.ClusterComputeResourceList(ctx, "*")
		if err != nil {
			return fmt.Errorf("finder failed finding clusters for datacenter: %s", err)
		}
		for _, cluster := range clusters {
			vc.Cluster2Datacenter[cluster.Reference().Value] = dc.Reference().Value
		}
	}
	return nil
}
