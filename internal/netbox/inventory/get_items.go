package inventory

import (
	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

// GetVlan returns the VLAN for the given groupID and vlanID.
// It returns nil if the VLAN is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVlan(groupID, vlanID int) (*objects.Vlan, bool) {
	nbi.VlansLock.Lock()
	defer nbi.VlansLock.Unlock()

	vlanGroup, groupExists := nbi.VlansIndexByVlanGroupIDAndVID[groupID]
	if !groupExists {
		return nil, false
	}

	vlan, vlanExists := vlanGroup[vlanID]
	if !vlanExists {
		return nil, false
	}

	// Remove the VLAN from the OrphanManager if found.
	nbi.OrphanManager.RemoveItem(constants.VlansAPIPath, &vlan.NetboxObject)
	return vlan, true
}

func (nbi *NetboxInventory) GetTenant(tenantName string) (*objects.Tenant, bool) {
	nbi.TenantsLock.Lock()
	defer nbi.TenantsLock.Unlock()
	tenant, tenantExists := nbi.TenantsIndexByName[tenantName]
	if !tenantExists {
		return nil, false
	}
	// Remove the Tenmant from the OrphanManager if found
	nbi.OrphanManager.RemoveItem(constants.TenantsAPIPath, &tenant.NetboxObject)
	return tenant, true
}

func (nbi *NetboxInventory) GetSite(siteName string) (*objects.Site, bool) {
	nbi.SitesLock.Lock()
	defer nbi.SitesLock.Unlock()
	site, siteExists := nbi.SitesIndexByName[siteName]
	if !siteExists {
		return nil, false
	}
	// Remove the Site from the OrphanManager if found
	nbi.OrphanManager.RemoveItem(constants.SitesAPIPath, &site.NetboxObject)
	return site, true
}

func (nbi *NetboxInventory) GetVlanGroup(vlanGroupName string) (*objects.VlanGroup, bool) {
	nbi.VlanGroupsLock.Lock()
	defer nbi.VlanGroupsLock.Unlock()
	vlanGroup, vlanGroupExists := nbi.VlanGroupsIndexByName[vlanGroupName]
	if !vlanGroupExists {
		return nil, false
	}
	// Remove the VlanGroup from the OrphanManager if found
	nbi.OrphanManager.RemoveItem(constants.VlanGroupsAPIPath, &vlanGroup.NetboxObject)
	return vlanGroup, true
}

func (nbi *NetboxInventory) GetClusterGroup(clusterGroupName string) (*objects.ClusterGroup, bool) {
	nbi.ClusterGroupsLock.Lock()
	defer nbi.ClusterGroupsLock.Unlock()
	clusterGroup, clusterGroupExists := nbi.ClusterGroupsIndexByName[clusterGroupName]
	if !clusterGroupExists {
		return nil, false
	}
	// Remove the clusterGroup from the OrphanManager if found
	nbi.OrphanManager.RemoveItem(constants.ClusterGroupsAPIPath, &clusterGroup.NetboxObject)
	return clusterGroup, true
}

func (nbi *NetboxInventory) GetCluster(clusterName string) (*objects.Cluster, bool) {
	nbi.ClustersLock.Lock()
	defer nbi.ClustersLock.Unlock()
	cluster, clusterExists := nbi.ClustersIndexByName[clusterName]
	if !clusterExists {
		return nil, false
	}
	// Remove the cluster from the OrphanManager if found
	nbi.OrphanManager.RemoveItem(constants.ClustersAPIPath, &cluster.NetboxObject)
	return cluster, true
}

func (nbi *NetboxInventory) GetDevice(deviceName string, siteID int) (*objects.Device, bool) {
	nbi.DevicesLock.Lock()
	defer nbi.DevicesLock.Unlock()
	device, deviceExists := nbi.DevicesIndexByNameAndSiteID[deviceName][siteID]
	if !deviceExists {
		return nil, false
	}
	// Remove the device from the OrphanManager if found
	nbi.OrphanManager.RemoveItem(constants.DevicesAPIPath, &device.NetboxObject)
	return device, true
}

func (nbi *NetboxInventory) GetDeviceRole(deviceRoleName string) (*objects.DeviceRole, bool) {
	nbi.DeviceRolesLock.Lock()
	defer nbi.DeviceRolesLock.Unlock()
	deviceRole, deviceRoleExists := nbi.DeviceRolesIndexByName[deviceRoleName]
	if !deviceRoleExists {
		return nil, false
	}
	// Remove the deviceRole from the OrphanManager if found
	nbi.OrphanManager.RemoveItem(constants.DeviceRolesAPIPath, &deviceRole.NetboxObject)
	return deviceRole, true
}
