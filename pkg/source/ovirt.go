package source

import (
	"fmt"
	"strings"
	"time"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/dcim"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/tenancy"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/virtualization"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// OVirtSource represents an oVirt source
type OVirtSource struct {
	CommonConfig
	Disks                  *ovirtsdk4.DiskSlice
	DataCenters            *ovirtsdk4.DataCenterSlice
	Clusters               *ovirtsdk4.ClusterSlice
	Hosts                  *ovirtsdk4.HostSlice
	Vms                    *ovirtsdk4.VmSlice
	HostSiteRelations      map[string]string
	ClusterSiteRelations   map[string]string
	ClusterTenantRelations map[string]string
	HostTenantRelations    map[string]string
	VmTenantRelations      map[string]string
}

func (o *OVirtSource) Init() error {
	// Initialise regex relations
	o.Logger.Info("Initializing regex relations for oVirt source ", o.SourceConfig.Name)
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

	// Initialise the connection
	o.Logger.Info("Initializing oVirt source ", o.SourceConfig.Name)
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

	// Get the disks
	disksResponse, err := conn.SystemService().DisksService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt disks: %v", err)
	}
	if disks, ok := disksResponse.Disks(); ok {
		o.Disks = disks
		o.Logger.Debug("Successfully initalized oVirt disks: ", disks)
	}

	// Get the DataCenters
	dataCentersResponse, err := conn.SystemService().DataCentersService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt data centers: %v", err)
	}
	if dataCenters, ok := dataCentersResponse.DataCenters(); ok {
		o.DataCenters = dataCenters
		o.Logger.Debug("Successfully initalized oVirt data centers: ", o.DataCenters)
	}

	// Get the clusters
	clustersResponse, err := conn.SystemService().ClustersService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt clusters: %v", err)
	}
	if clusters, ok := clustersResponse.Clusters(); ok {
		// Extract extra data for each cluster
		for _, cluster := range clusters.Slice() {
			datacenter, err := conn.FollowLink(cluster.MustDataCenter())
			if err != nil {
				return fmt.Errorf("failed to get datacenter for cluster %s: %v", cluster.MustName(), err)
			}
			cluster.SetDataCenter(datacenter.(*ovirtsdk4.DataCenter))
		}
		o.Clusters = clusters
		o.Logger.Debug("Successfully initalized oVirt clusters: ", o.Clusters)
	}

	//Get the hosts
	hostsResponse, err := conn.SystemService().HostsService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt hosts: %+v", err)
	}
	if hosts, ok := hostsResponse.Hosts(); ok {
		o.Hosts = hosts
		o.Logger.Debug("Successfully initalized oVirt hosts: ", hosts)
	}
	return nil
}

func (o *OVirtSource) Sync(nbi *inventory.NetBoxInventory) error {
	o.Logger.Info("Syncing oVirt source ", o.SourceConfig.Name, " with netbox")
	err := o.SyncDatacenters(nbi)
	if err != nil {
		return err
	}
	err = o.SyncClusters(nbi)
	if err != nil {
		return err
	}

	return nil
}

func (o *OVirtSource) SyncDatacenters(nbi *inventory.NetBoxInventory) error {
	// First sync oVirt DataCenters as NetBoxClusterGroups
	for _, datacenter := range o.DataCenters.Slice() {
		name, exists := datacenter.Name()
		if !exists {
			return fmt.Errorf("failed to get name for oVirt datacenter %s", name)
		}
		description, exists := datacenter.Description()
		if !exists {
			o.Logger.Warning("description for oVirt datacenter ", name, " is empty.")
		}
		nbClusterGroup := &virtualization.ClusterGroup{
			Name:        name,
			Slug:        strings.ToLower(name),
			Description: description,
		}
		err := nbi.AddClusterGroup(nbClusterGroup, o.SourceTag)
		if err != nil {
			return fmt.Errorf("failed to add oVirt data center %s as NetBox cluster group: %v", name, err)
		}
	}
	return nil
}

func (o *OVirtSource) SyncClusters(nbi *inventory.NetBoxInventory) error {
	clusterType := &virtualization.ClusterType{
		Name: "oVirt",
		Slug: "ovirt",
	}
	clusterType, err := nbi.AddClusterType(clusterType, o.SourceTag)
	if err != nil {
		return fmt.Errorf("failed to add oVirt cluster type: %v", err)
	}
	// Then sync oVirt Clusters as NetBoxClusters
	for _, cluster := range o.Clusters.Slice() {
		clusterName, exists := cluster.Name()
		if !exists {
			return fmt.Errorf("failed to get name for oVirt cluster %s", clusterName)
		}
		description, exists := cluster.Description()
		if !exists {
			o.Logger.Warning("description for oVirt cluster ", clusterName, " is empty.")
		}
		var clusterGroup *virtualization.ClusterGroup
		if dataCenter, ok := cluster.DataCenter(); ok {
			if dataCenterName, ok := dataCenter.Name(); ok {
				clusterGroup = nbi.ClusterGroupsIndexByName[dataCenterName]
			}
		}
		var clusterSite *dcim.Site
		if o.ClusterSiteRelations != nil {
			match, err := utils.MatchStringToValue(clusterName, o.ClusterSiteRelations)
			if err != nil {
				return fmt.Errorf("failed to match oVirt cluster %s to a NetBox site: %v", clusterName, err)
			}
			if match != "" {
				if _, ok := nbi.SitesIndexByName[match]; !ok {
					return fmt.Errorf("failed to match oVirt cluster %s to a NetBox site: %v. Site with this name doesn't exist!", clusterName, match)
				}
				clusterSite = nbi.SitesIndexByName[match]
			}
		}
		var clusterTenant *tenancy.Tenant
		if o.ClusterTenantRelations != nil {
			match, err := utils.MatchStringToValue(clusterName, o.ClusterTenantRelations)
			if err != nil {
				return fmt.Errorf("failed to match oVirt cluster %s to a NetBox tenant: %v", clusterName, err)
			}
			if match != "" {
				if _, ok := nbi.TenantsIndexByName[match]; !ok {
					return fmt.Errorf("failed to match oVirt cluster %s to a NetBox tenant: %v. Tenant with this name doesn't exist!", clusterName, match)
				}
				clusterTenant = nbi.TenantsIndexByName[match]
			}
		}

		nbCluster := &virtualization.Cluster{
			Name:        clusterName,
			Type:        clusterType,
			Status:      &dcim.Active,
			Group:       clusterGroup,
			Description: description,
			Site:        clusterSite,
			Tenant:      clusterTenant,
		}
		err := nbi.AddCluster(nbCluster, o.SourceTag)
		if err != nil {
			return fmt.Errorf("failed to add oVirt cluster %s as NetBox cluster: %v", clusterName, err)
		}
	}
	return nil
}
