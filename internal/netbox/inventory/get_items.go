package inventory

import (
	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
)

// GetTag returns the Tag for the given tagName.
// It returns nil if the Tag is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetTag(tagName string) (*objects.Tag, bool) {
	nbi.tagsLock.Lock()
	defer nbi.tagsLock.Unlock()
	tag, tagExists := nbi.tagsIndexByName[tagName]
	if !tagExists {
		return nil, false
	}
	return tag, true
}

// GetManufacturer returns the Manufacturer for the given manufacturerName.
// It returns nil if the Manufacturer is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetManufacturer(manufacturerName string) (*objects.Manufacturer, bool) {
	nbi.manufacturersLock.Lock()
	defer nbi.manufacturersLock.Unlock()
	manufacturer, manufacturerExists := nbi.manufacturersIndexByName[manufacturerName]
	if !manufacturerExists {
		return nil, false
	}
	return manufacturer, true
}

// GetCustomField returns the CustomField for the given customFieldName.
// It returns nil if the CustomField is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetCustomField(customFieldName string) (*objects.CustomField, bool) {
	nbi.customFieldsLock.Lock()
	defer nbi.customFieldsLock.Unlock()

	customField, customFieldExists := nbi.customFieldsIndexByName[customFieldName]
	if !customFieldExists {
		return nil, false
	}
	return customField, true
}

// GetVlan returns the VLAN for the given groupID and vlanID.
// It returns nil if the VLAN is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVlan(groupID, vlanID int) (*objects.Vlan, bool) {
	nbi.vlansLock.Lock()
	defer nbi.vlansLock.Unlock()

	vlanGroup, groupExists := nbi.vlansIndexByVlanGroupIDAndVID[groupID]
	if !groupExists {
		return nil, false
	}

	vlan, vlanExists := vlanGroup[vlanID]
	if !vlanExists {
		return nil, false
	}
	return vlan, true
}

// GetTenant returns the Tenant for the given tenantName.
// It returns nil if the Tenant is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetTenant(tenantName string) (*objects.Tenant, bool) {
	nbi.tenantsLock.Lock()
	defer nbi.tenantsLock.Unlock()
	tenant, tenantExists := nbi.tenantsIndexByName[tenantName]
	if !tenantExists {
		return nil, false
	}
	return tenant, true
}

// GetSite returns the Site for the given siteName.
// It returns nil if the Site is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetSite(siteName string) (*objects.Site, bool) {
	nbi.sitesLock.Lock()
	defer nbi.sitesLock.Unlock()
	site, siteExists := nbi.sitesIndexByName[siteName]
	if !siteExists {
		return nil, false
	}
	return site, true
}

// GetSiteByID returns the Site for the given siteID.
// It returns nil if the Site is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetSiteByID(siteID int) *objects.Site {
	nbi.sitesLock.Lock()
	defer nbi.sitesLock.Unlock()
	for _, site := range nbi.sitesIndexByName {
		if site.ID == siteID {
			return site
		}
	}
	return nil
}

// GetVlanGroup returns the VlanGroup for the given vlanGroupName.
// It returns nil if the VlanGroup is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVlanGroup(vlanGroupName string) (*objects.VlanGroup, bool) {
	nbi.vlanGroupsLock.Lock()
	defer nbi.vlanGroupsLock.Unlock()
	vlanGroup, vlanGroupExists := nbi.vlanGroupsIndexByName[vlanGroupName]
	if !vlanGroupExists {
		return nil, false
	}
	return vlanGroup, true
}

// GetClusterGroup returns the ClusterGroup for the given clusterGroupName.
// It returns nil if the ClusterGroup is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetClusterGroup(clusterGroupName string) (*objects.ClusterGroup, bool) {
	nbi.clusterGroupsLock.Lock()
	defer nbi.clusterGroupsLock.Unlock()
	clusterGroup, clusterGroupExists := nbi.clusterGroupsIndexByName[clusterGroupName]
	if !clusterGroupExists {
		return nil, false
	}
	return clusterGroup, true
}

// GetCluster returns the Cluster for the given clusterName.
// It returns nil if the Cluster is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetCluster(clusterName string) (*objects.Cluster, bool) {
	nbi.clustersLock.Lock()
	defer nbi.clustersLock.Unlock()
	cluster, clusterExists := nbi.clustersIndexByName[clusterName]
	if !clusterExists {
		return nil, false
	}
	// Remove the cluster from the OrphanManager if found
	return cluster, true
}

func (nbi *NetboxInventory) GetDevice(deviceName string, siteID int) (*objects.Device, bool) {
	nbi.devicesLock.Lock()
	defer nbi.devicesLock.Unlock()
	device, deviceExists := nbi.devicesIndexByNameAndSiteID[deviceName][siteID]
	if !deviceExists {
		return nil, false
	}
	return device, true
}

func (nbi *NetboxInventory) GetDeviceRole(deviceRoleName string) (*objects.DeviceRole, bool) {
	nbi.deviceRolesLock.Lock()
	defer nbi.deviceRolesLock.Unlock()
	deviceRole, deviceRoleExists := nbi.deviceRolesIndexByName[deviceRoleName]
	if !deviceRoleExists {
		return nil, false
	}
	return deviceRole, true
}

// GetContactRole returns the ContactRole for the given contactRoleName.
// It returns nil if the ContactRole is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetContactRole(contactRoleName string) (*objects.ContactRole, bool) {
	nbi.contactRolesLock.Lock()
	defer nbi.contactRolesLock.Unlock()

	contactRole, contactRoleExists := nbi.contactRolesIndexByName[contactRoleName]
	if !contactRoleExists {
		return nil, false
	}
	return contactRole, true
}

// GetVirtualDeviceContext returns the VirtualDeviceContext for the given zoneName and deviceID.
// It returns nil if the VirtualDeviceContext is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVirtualDeviceContext(
	zoneName string,
	deviceID int,
) (*objects.VirtualDeviceContext, bool) {
	nbi.virtualDeviceContextsLock.Lock()
	defer nbi.virtualDeviceContextsLock.Unlock()
	vdc, vdcExists := nbi.virtualDeviceContextsIndex[zoneName][deviceID]
	if !vdcExists {
		return nil, false
	}
	return vdc, true
}

// GetInterface returns the Interface for the given interfaceName and deviceID.
// It returns nil if the Interface is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetInterface(
	interfaceName string,
	deviceID int,
) (*objects.Interface, bool) {
	nbi.interfacesLock.Lock()
	defer nbi.interfacesLock.Unlock()

	iface, ifaceExists := nbi.interfacesIndexByDeviceIDAndName[deviceID][interfaceName]
	if !ifaceExists {
		return nil, false
	}
	return iface, true
}

// GetContactAssignment returns the ContactAssignment for the given contentType, objectID, contactID and roleID.
// It returns nil if the ContactAssignment is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetContactAssignment(
	contentType constants.ContentType,
	objectID int,
	contactID int,
	roleID int,
) (*objects.ContactAssignment, bool) {
	nbi.contactAssignmentsLock.Lock()
	defer nbi.contactAssignmentsLock.Unlock()
	contactAssignment, contactAssignmentExists := nbi.contactAssignmentsIndex[contentType][objectID][contactID][roleID]
	if !contactAssignmentExists {
		return nil, false
	}
	return contactAssignment, true
}

// GetInterfaceByID returns the Interface for the given interfaceID.
// It returns nil if the Interface is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetInterfaceByID(interfaceID int) *objects.Interface {
	nbi.interfacesLock.Lock()
	defer nbi.interfacesLock.Unlock()
	return nbi.interfacesIndexByID[interfaceID]
}

// GetVMInterfaceByID returns the VMInterface for the given vmInterfaceID.
// It returns nil if the VMInterface is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVMInterfaceByID(vmInterfaceID int) *objects.VMInterface {
	nbi.vmInterfacesLock.Lock()
	defer nbi.vmInterfacesLock.Unlock()
	return nbi.vmInterfacesIndexByID[vmInterfaceID]
}

// GetDeviceByID returns the Device for the given deviceID.
// It returns nil if the Device is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetDeviceByID(deviceID int) *objects.Device {
	nbi.devicesLock.Lock()
	defer nbi.devicesLock.Unlock()
	return nbi.devicesIndexByID[deviceID]
}

// GetVMByID returns the VirtualMachine for the given vmID.
// It returns nil if the VirtualMachine is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVMByID(vmID int) *objects.VM {
	nbi.vmsLock.Lock()
	defer nbi.vmsLock.Unlock()
	return nbi.vmsIndexByID[vmID]
}

// GetVRF returns the VRF for the given vrfName.
// It returns nil if the VRF is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVRF(vrfName string) (*objects.VRF, bool) {
	nbi.vrfsLock.Lock()
	defer nbi.vrfsLock.Unlock()
	vrf, vrfExists := nbi.vrfsIndexByName[vrfName]
	if !vrfExists {
		return nil, false
	}
	return vrf, true
}