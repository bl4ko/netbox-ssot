package inventory

import (
	"github.com/bl4ko/netbox-ssot/pkg/netbox/dcim"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/extras"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/tenancy"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/virtualization"
)

// Collect all tags from NetBox API and store them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitTags() error {
	nbTags, err := netboxInventory.NetboxApi.GetAllTags()
	if err != nil {
		return err
	}
	netboxInventory.Tags = nbTags
	netboxInventory.Logger.Debug("Successfully collected tags from NetBox: ", netboxInventory.Tags)

	// Custom tag for all netbox objects
	ssotTag, err := netboxInventory.NetboxApi.GetTagByName("netbox-ssot")
	if err != nil {
		return err
	}
	if ssotTag == nil {
		netboxInventory.Logger.Info("Tag netbox-ssot not found in NetBox. Creating it now...")
		newTag := extras.Tag{Name: "netbox-ssot", Slug: "netbox-ssot", Description: "Tag used by netbox-ssot to mark devices that are managed by it", Color: "00add8"}
		ssotTag, err = netboxInventory.NetboxApi.CreateTag(&newTag)
		if err != nil {
			return err
		}
	}
	netboxInventory.SsotTag = ssotTag
	return nil
}

// Collects all tenants from NetBox API and store them in the NetBoxInventory
func (NetBoxInventory *NetBoxInventory) InitTenants() error {
	nbTenants, err := NetBoxInventory.NetboxApi.GetAllTenants()
	if err != nil {
		return err
	}
	// We also create an index of tenants by name for easier access
	NetBoxInventory.TenantsIndexByName = make(map[string]*tenancy.Tenant)
	for _, tenant := range nbTenants {
		NetBoxInventory.TenantsIndexByName[tenant.Name] = tenant
	}
	NetBoxInventory.Logger.Debug("Successfully collected tenants from NetBox: ", NetBoxInventory.TenantsIndexByName)
	return nil
}

// Collects all sites from NetBox API and store them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitSites() error {
	nbSites, err := netboxInventory.NetboxApi.GetAllSites()
	if err != nil {
		return err
	}
	// We also create an index of sites by name for easier access
	netboxInventory.SitesIndexByName = make(map[string]*dcim.Site)
	for _, site := range nbSites {
		netboxInventory.SitesIndexByName[site.Name] = site
	}
	netboxInventory.Logger.Debug("Successfully collected sites from NetBox: ", netboxInventory.SitesIndexByName)
	return nil
}

// Collect all devices from NetBox API and store them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitDevices() error {
	nbDevices, err := netboxInventory.NetboxApi.GetAllDevices()
	if err != nil {
		return err
	}
	// We also create an index of devices by name for easier access
	netboxInventory.DevicesIndexByName = make(map[string]*dcim.Device)
	for _, device := range nbDevices {
		netboxInventory.DevicesIndexByName[device.Name] = device
	}
	netboxInventory.Logger.Debug("Successfully collected devices from NetBox: ", netboxInventory.DevicesIndexByName)
	return nil
}

// Collects all nbClusters from NetBox API and stores them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitClusterGroups() error {
	nbClusters, err := netboxInventory.NetboxApi.GetAllClusterGroups()
	if err != nil {
		return err
	}
	// We also create an index of cluster groups by name for easier access
	netboxInventory.ClusterGroupsIndexByName = make(map[string]*virtualization.ClusterGroup)
	for _, clusterGroup := range nbClusters {
		netboxInventory.ClusterGroupsIndexByName[clusterGroup.Name] = clusterGroup
	}
	netboxInventory.Logger.Debug("Successfully collected cluster groups from NetBox: ", netboxInventory.ClusterGroupsIndexByName)
	return nil
}

// Collects all ClusterTypes from NetBox API and stores them in the NetBoxInventory
func (netboxInventory *NetBoxInventory) InitClusterTypes() error {
	nbClusterTypes, err := netboxInventory.NetboxApi.GetAllClusterTypes()
	if err != nil {
		return err
	}
	netboxInventory.ClusterTypesIndexByName = make(map[string]*virtualization.ClusterType)
	for _, clusterType := range nbClusterTypes {
		netboxInventory.ClusterTypesIndexByName[clusterType.Name] = clusterType
	}
	netboxInventory.Logger.Debug("Successfully collected cluster types from NetBox: ", netboxInventory.ClusterTypesIndexByName)
	return nil
}

// Collects all clusters from NetBox API and stores them to local inventory
func (netboxInventory *NetBoxInventory) InitClusters() error {
	nbClusters, err := netboxInventory.NetboxApi.GetAllClusters()
	if err != nil {
		return err
	}
	netboxInventory.ClustersIndexByName = make(map[string]*virtualization.Cluster)
	for _, cluster := range nbClusters {
		netboxInventory.ClustersIndexByName[cluster.Name] = cluster
	}
	netboxInventory.Logger.Debug("Successfully collected clusters from NetBox: ", netboxInventory.ClustersIndexByName)
	return nil
}
