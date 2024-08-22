package inventory

import "github.com/bl4ko/netbox-ssot/internal/netbox/objects"

// GetVlan returns vlan for the given vlanGroupID and vlanID.
// Returns nil if vlan is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVlan(vlanGroupID int, vlanID int) *objects.Vlan {
	nbi.VlansLock.Lock()
	defer nbi.VlansLock.Unlock()
	if _, ok := nbi.VlansIndexByVlanGroupIDAndVID[vlanGroupID]; !ok {
		return nil
	}
	return nbi.VlansIndexByVlanGroupIDAndVID[vlanGroupID][vlanID]
}

func (nbi *NetboxInventory) GetTenant(tenantName string) (*objects.Tenant, bool) {
	nbi.TenantsLock.Lock()
	defer nbi.TenantsLock.Unlock()
	if _, ok := nbi.TenantsIndexByName[tenantName]; !ok {
		return nil, false
	}
	return nbi.TenantsIndexByName[tenantName], true
}

func (nbi *NetboxInventory) GetSite(siteName string) (*objects.Site, bool) {
	nbi.SitesLock.Lock()
	defer nbi.SitesLock.Unlock()
	if _, ok := nbi.SitesIndexByName[siteName]; !ok {
		return nil, false
	}
	return nbi.SitesIndexByName[siteName], true
}

func (nbi *NetboxInventory) GetVlanGroup(vlanGroupName string) (*objects.VlanGroup, bool) {
	nbi.VlanGroupsLock.Lock()
	defer nbi.VlanGroupsLock.Unlock()
	if _, ok := nbi.VlanGroupsIndexByName[vlanGroupName]; !ok {
		return nil, false
	}
	return nbi.VlanGroupsIndexByName[vlanGroupName], true
}

func (nbi *NetboxInventory) GetClusterGroup(clusterGroupName string) (*objects.ClusterGroup, bool) {
	nbi.ClusterGroupsLock.Lock()
	defer nbi.ClusterGroupsLock.Unlock()
	if _, ok := nbi.ClusterGroupsIndexByName[clusterGroupName]; !ok {
		return nil, false
	}
	return nbi.ClusterGroupsIndexByName[clusterGroupName], true
}

func (nbi *NetboxInventory) GetCluster(clusterName string) (*objects.Cluster, bool) {
	nbi.ClustersLock.Lock()
	defer nbi.ClustersLock.Unlock()
	if _, ok := nbi.ClustersIndexByName[clusterName]; !ok {
		return nil, false
	}
	return nbi.ClustersIndexByName[clusterName], true
}

func (nbi *NetboxInventory) GetDevice(deviceName string, siteID int) (*objects.Device, bool) {
	nbi.DevicesLock.Lock()
	defer nbi.DevicesLock.Unlock()
	if device, ok := nbi.DevicesIndexByNameAndSiteID[deviceName][siteID]; ok {
		return device, true
	}
	return nil, false
}
