package source

import (
	"fmt"
	"time"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/common"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/tenancy"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/virtualization"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// OVirtSource represents an oVirt source
type OVirtSource struct {
	CommonConfig
	Disks       map[string]*ovirtsdk4.Disk
	DataCenters map[string]*ovirtsdk4.DataCenter
	Clusters    map[string]*ovirtsdk4.Cluster
	Hosts       map[string]*ovirtsdk4.Host
	Vms         map[string]*ovirtsdk4.Vm

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

	err = o.InitDisks(conn)
	if err != nil {
		return fmt.Errorf("failed to initialize oVirt disks: %v", err)
	}

	err = o.InitDataCenters(conn)
	if err != nil {
		return fmt.Errorf("failed to initialize oVirt data centers: %v", err)
	}

	err = o.InitClusters(conn)
	if err != nil {
		return fmt.Errorf("failed to initialize oVirt clusters: %v", err)
	}

	err = o.InitHosts(conn)
	if err != nil {
		return fmt.Errorf("failed to initialize oVirt hosts: %v", err)
	}

	err = o.InitVms(conn)
	if err != nil {
		return fmt.Errorf("failed to initialize oVirt vms: %v", err)
	}
	return nil
}

func (o *OVirtSource) InitDisks(conn *ovirtsdk4.Connection) error {
	// Get the disks
	disksResponse, err := conn.SystemService().DisksService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt disks: %v", err)
	}
	o.Disks = make(map[string]*ovirtsdk4.Disk)
	if disks, ok := disksResponse.Disks(); ok {
		for _, disk := range disks.Slice() {
			o.Disks[disk.MustId()] = disk
		}
		o.Logger.Debug("Successfully initalized oVirt disks: ", o.Disks)
	} else {
		o.Logger.Warning("Error initialising oVirt disks")
	}
	return nil
}

func (o *OVirtSource) InitDataCenters(conn *ovirtsdk4.Connection) error {
	dataCentersResponse, err := conn.SystemService().DataCentersService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt data centers: %v", err)
	}
	o.DataCenters = make(map[string]*ovirtsdk4.DataCenter)
	if dataCenters, ok := dataCentersResponse.DataCenters(); ok {
		for _, dataCenter := range dataCenters.Slice() {
			o.DataCenters[dataCenter.MustId()] = dataCenter
		}
		o.Logger.Debug("Successfully initalized oVirt data centers: ", o.DataCenters)
	} else {
		o.Logger.Warning("Error initialising oVirt data centers")
	}
	return nil
}

// Function that queries ovirt api for clustrers and stores them locally
func (o *OVirtSource) InitClusters(conn *ovirtsdk4.Connection) error {
	clustersResponse, err := conn.SystemService().ClustersService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt clusters: %v", err)
	}
	o.Clusters = make(map[string]*ovirtsdk4.Cluster)
	if clusters, ok := clustersResponse.Clusters(); ok {
		for _, cluster := range clusters.Slice() {
			o.Clusters[cluster.MustId()] = cluster
		}
		o.Logger.Debug("Successfully initalized oVirt clusters: ", o.Clusters)
	} else {
		o.Logger.Warning("Error initialising oVirt clusters")
	}
	return nil
}

// Function that queries ovirt api for hosts and stores them locally
func (o *OVirtSource) InitHosts(conn *ovirtsdk4.Connection) error {
	hostsResponse, err := conn.SystemService().HostsService().List().Follow("nics").Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt hosts: %+v", err)
	}
	o.Hosts = make(map[string]*ovirtsdk4.Host)
	if hosts, ok := hostsResponse.Hosts(); ok {
		for _, host := range hosts.Slice() {
			o.Hosts[host.MustId()] = host
		}
		o.Logger.Debug("Successfully initalized oVirt hosts: ", hosts)
	} else {
		o.Logger.Warning("Error initialising oVirt hosts")
	}
	return nil
}

// Function that quries the ovirt api for vms and stores them locally
func (o *OVirtSource) InitVms(conn *ovirtsdk4.Connection) error {
	vmsResponse, err := conn.SystemService().VmsService().List().Follow("nics,diskattachments,reporteddevices").Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt vms: %+v", err)
	}
	o.Vms = make(map[string]*ovirtsdk4.Vm)
	if vms, ok := vmsResponse.Vms(); ok {
		for _, vm := range vms.Slice() {
			o.Vms[vm.MustId()] = vm
		}
		o.Logger.Debug("Successfully initalized oVirt vms: ", vms)
	} else {
		o.Logger.Warning("Error initialising oVirt vms")
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
	err = o.SyncHosts(nbi)
	if err != nil {
		return err
	}

	// TODO
	// err = o.SyncVms(nbi)

	return nil
}

func (o *OVirtSource) SyncDatacenters(nbi *inventory.NetBoxInventory) error {
	// First sync oVirt DataCenters as NetBoxClusterGroups
	for _, datacenter := range o.DataCenters {
		name, exists := datacenter.Name()
		if !exists {
			return fmt.Errorf("failed to get name for oVirt datacenter %s", name)
		}
		description, exists := datacenter.Description()
		if !exists {
			o.Logger.Warning("description for oVirt datacenter ", name, " is empty.")
		}
		nbClusterGroup := &virtualization.ClusterGroup{
			NetboxObject: common.NetboxObject{Description: description},
			Name:         name,
			Slug:         utils.Slugify(name),
		}
		err := nbi.AddClusterGroup(nbClusterGroup, o.SourceTags)
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
	clusterType, err := nbi.AddClusterType(clusterType, o.SourceTags)
	if err != nil {
		return fmt.Errorf("failed to add oVirt cluster type: %v", err)
	}
	// Then sync oVirt Clusters as NetBoxClusters
	for _, cluster := range o.Clusters {
		clusterName, exists := cluster.Name()
		if !exists {
			return fmt.Errorf("failed to get name for oVirt cluster %s", clusterName)
		}
		description, exists := cluster.Description()
		if !exists {
			o.Logger.Warning("description for oVirt cluster ", clusterName, " is empty.")
		}
		var clusterGroup *virtualization.ClusterGroup
		if _, ok := o.DataCenters[cluster.MustDataCenter().MustId()]; ok {

		} else {
			o.Logger.Warning("failed to get datacenter for oVirt cluster ", clusterName)
		}
		if dataCenter, ok := cluster.DataCenter(); ok {
			if dataCenterName, ok := dataCenter.Name(); ok {
				clusterGroup = nbi.ClusterGroupsIndexByName[dataCenterName]
			}
		}
		var clusterSite *common.Site
		if o.ClusterSiteRelations != nil {
			match, err := utils.MatchStringToValue(clusterName, o.ClusterSiteRelations)
			if err != nil {
				return fmt.Errorf("failed to match oVirt cluster %s to a NetBox site: %v", clusterName, err)
			}
			if match != "" {
				if _, ok := nbi.SitesIndexByName[match]; !ok {
					return fmt.Errorf("failed to match oVirt cluster %s to a NetBox site: %v. Site with this name doesn't exist", clusterName, match)
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
					return fmt.Errorf("failed to match oVirt cluster %s to a NetBox tenant: %v. Tenant with this name doesn't exist", clusterName, match)
				}
				clusterTenant = nbi.TenantsIndexByName[match]
			}
		}

		nbCluster := &virtualization.Cluster{
			NetboxObject: common.NetboxObject{Description: description},
			Name:         clusterName,
			Type:         clusterType,
			Status:       virtualization.ClusterStatusActive,
			Group:        clusterGroup,
			Site:         clusterSite,
			Tenant:       clusterTenant,
		}
		err := nbi.AddCluster(nbCluster, o.SourceTags)
		if err != nil {
			return fmt.Errorf("failed to add oVirt cluster %s as NetBox cluster: %v", clusterName, err)
		}
	}
	return nil
}

// Host in oVirt is a represented as device in netbox with a
// custom role Server
func (o *OVirtSource) SyncHosts(nbi *inventory.NetBoxInventory) error {
	// for hostId, host := range o.Hosts {
	// 	hostCluster := nbi.ClustersIndexByName[o.Clusters[host.MustCluster().MustId()].MustName()]

	// 	var hostSite *common.Site
	// 	if o.HostSiteRelations != nil {
	// 		match, err := utils.MatchStringToValue(host.MustName(), o.HostSiteRelations)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to match oVirt host %s to a NetBox site: %v", host.MustName(), err)
	// 		}
	// 		if match != "" {
	// 			if _, ok := nbi.SitesIndexByName[match]; !ok {
	// 				return fmt.Errorf("failed to match oVirt host %s to a NetBox site: %v. Site with this name doesn't exist", host.MustName(), match)
	// 			}
	// 			hostSite = nbi.SitesIndexByName[match]
	// 		}
	// 	}
	// 	nbHost := &dcim.Device{
	// 		Name: 			host.MustName(),
	// 		DeviceRole: 	nbi.DeviceRolesIndexByName["Server"],
	// 		// DeviceType: , # TODO
	// 		Cluster: 		hostCluster,

	// 	}
	// }

	return nil
}
