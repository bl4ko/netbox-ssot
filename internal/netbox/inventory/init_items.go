package inventory

import (
	"context"
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Collect all tags from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initTags(ctx context.Context) error {
	nbTags, err := service.GetAll[objects.Tag](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	nbi.tagsIndexByName = make(map[string]*objects.Tag)
	for i := range nbTags {
		tag := nbTags[i]
		nbi.tagsIndexByName[tag.Name] = &tag
	}
	nbi.Logger.Debug(ctx, "Successfully collected tags from Netbox: ", nbi.tagsIndexByName)

	// Create default tag for netbox-ssot microservice
	ssotTag, err := nbi.AddTag(ctx, &objects.Tag{Name: constants.SsotTagName, Slug: constants.SsotTagName, Description: constants.SsotTagDescription, Color: constants.SsotTagColor})
	if err != nil {
		return fmt.Errorf("error creating default ssot  tag: %s", err)
	}

	nbi.SsotTag = ssotTag

	// Create default tag for orphaned objects
	orphanTag, err := nbi.AddTag(ctx, &objects.Tag{Name: constants.OrphanTagName, Slug: constants.OrphanTagName, Description: constants.OrphanTagDescription, Color: constants.OrphanTagColor})
	if err != nil {
		return fmt.Errorf("error creating default orphan tag: %s", err)
	}
	nbi.OrphanManager.Tag = orphanTag
	return nil
}

// Collects all tenants from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initTenants(ctx context.Context) error {
	nbTenants, err := service.GetAll[objects.Tenant](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of tenants by name for easier access
	nbi.tenantsIndexByName = make(map[string]*objects.Tenant)
	for i := range nbTenants {
		tenant := &nbTenants[i]
		nbi.tenantsIndexByName[tenant.Name] = tenant
	}
	nbi.Logger.Debug(ctx, "Successfully collected tenants from Netbox: ", nbi.tenantsIndexByName)
	return nil
}

// Collects all contacts from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initContacts(ctx context.Context) error {
	nbContacts, err := service.GetAll[objects.Contact](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of contacts by name for easier access
	nbi.contactsIndexByName = make(map[string]*objects.Contact)
	for i := range nbContacts {
		contact := &nbContacts[i]
		nbi.contactsIndexByName[contact.Name] = contact
		nbi.OrphanManager.AddItem(constants.ContactsAPIPath, contact)
	}
	nbi.Logger.Debug(ctx, "Successfully collected contacts from Netbox: ", nbi.contactsIndexByName)
	return nil
}

// Collects all contact roles from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initContactRoles(ctx context.Context) error {
	nbContactRoles, err := service.GetAll[objects.ContactRole](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of contact roles by name for easier access
	nbi.contactRolesIndexByName = make(map[string]*objects.ContactRole)
	for i := range nbContactRoles {
		contactRole := &nbContactRoles[i]
		nbi.contactRolesIndexByName[contactRole.Name] = contactRole
	}
	nbi.Logger.Debug(ctx, "Successfully collected ContactRoles from Netbox: ", nbi.contactRolesIndexByName)
	return nil
}

func (nbi *NetboxInventory) initContactAssignments(ctx context.Context) error {
	nbCAs, err := service.GetAll[objects.ContactAssignment](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of contacts by name for easier access
	nbi.contactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID = make(map[constants.ContentType]map[int]map[int]map[int]*objects.ContactAssignment)
	debugIDs := map[int]bool{} // Netbox pagination bug duplicates
	for i := range nbCAs {
		cA := &nbCAs[i]
		if _, ok := debugIDs[cA.ID]; ok {
			fmt.Printf("Already been here: %d", cA.ID)
		}
		debugIDs[cA.ID] = true
		if nbi.contactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[cA.ModelType] == nil {
			nbi.contactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[cA.ModelType] = make(map[int]map[int]map[int]*objects.ContactAssignment)
		}
		if nbi.contactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[cA.ModelType][cA.ObjectID] == nil {
			nbi.contactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[cA.ModelType][cA.ObjectID] = make(map[int]map[int]*objects.ContactAssignment)
		}
		if nbi.contactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[cA.ModelType][cA.ObjectID][cA.Contact.ID] == nil {
			nbi.contactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[cA.ModelType][cA.ObjectID][cA.Contact.ID] = make(map[int]*objects.ContactAssignment)
		}
		nbi.contactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[cA.ModelType][cA.ObjectID][cA.Contact.ID][cA.Role.ID] = cA
		nbi.OrphanManager.AddItem(constants.ContactAssignmentsAPIPath, cA)
	}
	nbi.Logger.Debug(ctx, "Successfully collected contacts from Netbox: ", nbi.contactsIndexByName)
	return nil
}

// Initializes default admin contact role used for adding admin contacts of vms.
func (nbi *NetboxInventory) initAdminContactRole(ctx context.Context) error {
	_, err := nbi.AddContactRole(ctx, &objects.ContactRole{
		NetboxObject: objects.NetboxObject{
			Description: "Auto generated contact role by netbox-ssot for admins of vms.",
		},
		Name: objects.AdminContactRoleName,
		Slug: utils.Slugify(objects.AdminContactRoleName),
	})
	if err != nil {
		return fmt.Errorf("admin contact role: %s", err)
	}
	return nil
}

// Collects all contact groups from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initContactGroups(ctx context.Context) error {
	nbContactGroups, err := service.GetAll[objects.ContactGroup](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of contact groups by name for easier access
	nbi.contactGroupsIndexByName = make(map[string]*objects.ContactGroup)
	for i := range nbContactGroups {
		contactGroup := &nbContactGroups[i]
		nbi.contactGroupsIndexByName[contactGroup.Name] = contactGroup
	}
	nbi.Logger.Debug(ctx, "Successfully collected ContactGroups from Netbox: ", nbi.contactGroupsIndexByName)
	return nil
}

// Collects all sites from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initSites(ctx context.Context) error {
	nbSites, err := service.GetAll[objects.Site](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of sites by name for easier access
	nbi.sitesIndexByName = make(map[string]*objects.Site)
	for i := range nbSites {
		site := &nbSites[i]
		nbi.sitesIndexByName[site.Name] = site
	}
	nbi.Logger.Debug(ctx, "Successfully collected sites from Netbox: ", nbi.sitesIndexByName)
	return nil
}

// initDefaultSite inits default site, which is used for hosts that have no corresponding site.
// This is because site is required for adding new hosts.
func (nbi *NetboxInventory) initDefaultSite(ctx context.Context) error {
	_, err := nbi.AddSite(ctx, &objects.Site{
		NetboxObject: objects.NetboxObject{
			Tags:        []*objects.Tag{nbi.SsotTag},
			Description: "Default netbox-ssot site used for all hosts, that have no site matched.",
			CustomFields: map[string]interface{}{
				constants.CustomFieldSourceName: nbi.SsotTag.Name,
			},
		},
		Name: constants.DefaultSite,
		Slug: utils.Slugify(constants.DefaultSite),
	})
	if err != nil {
		return fmt.Errorf("init default site: %s", err)
	}
	return nil
}

// Collects all manufacturers from Netbox API and store them in NetBoxInventory.
func (nbi *NetboxInventory) initManufacturers(ctx context.Context) error {
	nbManufacturers, err := service.GetAll[objects.Manufacturer](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of manufacturers by name
	nbi.manufacturersIndexByName = make(map[string]*objects.Manufacturer)
	for i := range nbManufacturers {
		manufacturer := &nbManufacturers[i]
		nbi.manufacturersIndexByName[manufacturer.Name] = manufacturer
		nbi.OrphanManager.AddItem(constants.ManufacturersAPIPath, manufacturer)
	}

	nbi.Logger.Debug(ctx, "Successfully collected manufacturers from Netbox: ", nbi.manufacturersIndexByName)
	return nil
}

// Collects all platforms from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initPlatforms(ctx context.Context) error {
	nbPlatforms, err := service.GetAll[objects.Platform](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of platforms by name
	nbi.platformsIndexByName = make(map[string]*objects.Platform)

	for i, platform := range nbPlatforms {
		nbPlatform := &nbPlatforms[i]
		nbi.platformsIndexByName[platform.Name] = nbPlatform
		nbi.OrphanManager.AddItem(constants.PlatformsAPIPath, nbPlatform)
	}

	nbi.Logger.Debug(ctx, "Successfully collected platforms from Netbox: ", nbi.platformsIndexByName)
	return nil
}

// Collect all devices from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initDevices(ctx context.Context) error {
	nbDevices, err := service.GetAll[objects.Device](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of devices by Name and SiteId
	nbi.devicesIndexByNameAndSiteID = make(map[string]map[int]*objects.Device)

	for i, device := range nbDevices {
		nbDevice := &nbDevices[i]
		if nbi.devicesIndexByNameAndSiteID[device.Name] == nil {
			nbi.devicesIndexByNameAndSiteID[device.Name] = make(map[int]*objects.Device)
		}
		nbi.devicesIndexByNameAndSiteID[device.Name][device.Site.ID] = nbDevice
		nbi.OrphanManager.AddItem(constants.DevicesAPIPath, nbDevice)
	}

	nbi.Logger.Debug(ctx, "Successfully collected devices from Netbox: ", nbi.devicesIndexByNameAndSiteID)
	return nil
}

// Collect all devices from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initVirtualDeviceContexts(ctx context.Context) error {
	nbVirtualDeviceContexts, err := service.GetAll[objects.VirtualDeviceContext](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of devices by Name and SiteId
	nbi.virtualDeviceContextsIndexByNameAndDeviceID = make(map[string]map[int]*objects.VirtualDeviceContext)
	for i, virtualDeviceContext := range nbVirtualDeviceContexts {
		nbVirtualDeviceContext := &nbVirtualDeviceContexts[i]
		if nbi.virtualDeviceContextsIndexByNameAndDeviceID[virtualDeviceContext.Name] == nil {
			nbi.virtualDeviceContextsIndexByNameAndDeviceID[virtualDeviceContext.Name] = make(map[int]*objects.VirtualDeviceContext)
		}
		nbi.virtualDeviceContextsIndexByNameAndDeviceID[virtualDeviceContext.Name][virtualDeviceContext.Device.ID] = nbVirtualDeviceContext
		nbi.OrphanManager.AddItem(constants.VirtualDeviceContextsAPIPath, nbVirtualDeviceContext)
	}

	nbi.Logger.Debug(ctx, "Successfully collected VirtualDeviceContexts from Netbox: ", nbi.virtualDeviceContextsIndexByNameAndDeviceID)
	return nil
}

// Collects all deviceRoles from Netbox API and store them in the
// NetBoxInventory.
func (nbi *NetboxInventory) initDeviceRoles(ctx context.Context) error {
	nbDeviceRoles, err := service.GetAll[objects.DeviceRole](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of device roles by name for easier access
	nbi.deviceRolesIndexByName = make(map[string]*objects.DeviceRole)

	for i := range nbDeviceRoles {
		deviceRole := &nbDeviceRoles[i]
		nbi.deviceRolesIndexByName[deviceRole.Name] = deviceRole
		nbi.OrphanManager.AddItem(constants.DeviceRolesAPIPath, deviceRole)
	}

	nbi.Logger.Debug(ctx, "Successfully collected device roles from Netbox: ", nbi.deviceRolesIndexByName)
	return nil
}

func (nbi *NetboxInventory) initCustomFields(ctx context.Context) error {
	customFields, err := service.GetAll[objects.CustomField](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of custom fields by name
	nbi.customFieldsIndexByName = make(map[string]*objects.CustomField, len(customFields))
	for i := range customFields {
		customField := &customFields[i]
		nbi.customFieldsIndexByName[customField.Name] = customField
	}
	nbi.Logger.Debug(ctx, "Successfully collected custom fields from Netbox: ", nbi.customFieldsIndexByName)
	return nil
}

// This function Initializes all custom fields required for servers and other objects
// Currently these are two:
// - host_cpu_cores
// - host_memory
// - sourceId - this is used to store the ID of the source object in Netbox (interfaces).
func (nbi *NetboxInventory) initSsotCustomFields(ctx context.Context) error {
	// Custom field for storing object's source name.
	_, err := nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldSourceName,
		Label:                 constants.CustomFieldSourceLabel,
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldSourceDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ObjectTypes:           []constants.ContentType{constants.ContentTypeDcimDevice, constants.ContentTypeDcimDeviceRole, constants.ContentTypeDcimDeviceType, constants.ContentTypeDcimInterface, constants.ContentTypeDcimLocation, constants.ContentTypeDcimManufacturer, constants.ContentTypeDcimPlatform, constants.ContentTypeDcimRegion, constants.ContentTypeDcimSite, constants.ContentTypeVirtualDeviceContext, constants.ContentTypeIpamIPAddress, constants.ContentTypeIpamVlanGroup, constants.ContentTypeIpamVlan, constants.ContentTypeIpamPrefix, constants.ContentTypeTenancyTenantGroup, constants.ContentTypeTenancyTenant, constants.ContentTypeTenancyContact, constants.ContentTypeTenancyContactAssignment, constants.ContentTypeTenancyContactGroup, constants.ContentTypeTenancyContactRole, constants.ContentTypeVirtualizationCluster, constants.ContentTypeVirtualizationClusterGroup, constants.ContentTypeVirtualizationClusterType, constants.ContentTypeVirtualizationVirtualMachine, constants.ContentTypeVirtualizationVMInterface, constants.ContentTypeWirelessLAN, constants.ContentTypeWirelessLANGroup},
	})
	if err != nil {
		return fmt.Errorf("add source custom field %s", err)
	}
	// Custom field for marking when the object was last seen.
	// This is useful for orphan manager so we can delete objects
	// that haven't been seen for a while.
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldOrphanLastSeenName,
		Label:                 constants.CustomFieldOrphanLastSeenLabel,
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldOrphanLastSeenDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ObjectTypes:           []constants.ContentType{constants.ContentTypeDcimDevice, constants.ContentTypeDcimDeviceRole, constants.ContentTypeDcimDeviceType, constants.ContentTypeDcimInterface, constants.ContentTypeDcimLocation, constants.ContentTypeDcimManufacturer, constants.ContentTypeDcimPlatform, constants.ContentTypeDcimRegion, constants.ContentTypeDcimSite, constants.ContentTypeVirtualDeviceContext, constants.ContentTypeIpamIPAddress, constants.ContentTypeIpamVlanGroup, constants.ContentTypeIpamVlan, constants.ContentTypeIpamPrefix, constants.ContentTypeTenancyTenantGroup, constants.ContentTypeTenancyTenant, constants.ContentTypeTenancyContact, constants.ContentTypeTenancyContactAssignment, constants.ContentTypeTenancyContactGroup, constants.ContentTypeTenancyContactRole, constants.ContentTypeVirtualizationCluster, constants.ContentTypeVirtualizationClusterGroup, constants.ContentTypeVirtualizationClusterType, constants.ContentTypeVirtualizationVirtualMachine, constants.ContentTypeVirtualizationVMInterface, constants.ContentTypeWirelessLAN, constants.ContentTypeWirelessLANGroup},
	})
	if err != nil {
		return fmt.Errorf("add last seen custom field: %s", err)
	}
	// Custom field for storing object's source id.
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldSourceIDName,
		Label:                 constants.CustomFieldSourceIDLabel,
		Type:                  objects.CustomFieldTypeText,
		Default:               nil,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldSourceIDDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ObjectTypes:           []constants.ContentType{constants.ContentTypeDcimDevice, constants.ContentTypeDcimDeviceRole, constants.ContentTypeDcimDeviceType, constants.ContentTypeDcimInterface, constants.ContentTypeDcimLocation, constants.ContentTypeDcimManufacturer, constants.ContentTypeDcimPlatform, constants.ContentTypeDcimRegion, constants.ContentTypeDcimSite, constants.ContentTypeVirtualDeviceContext, constants.ContentTypeIpamIPAddress, constants.ContentTypeIpamVlanGroup, constants.ContentTypeIpamVlan, constants.ContentTypeIpamPrefix, constants.ContentTypeTenancyTenantGroup, constants.ContentTypeTenancyTenant, constants.ContentTypeTenancyContact, constants.ContentTypeTenancyContactAssignment, constants.ContentTypeTenancyContactGroup, constants.ContentTypeTenancyContactRole, constants.ContentTypeVirtualizationCluster, constants.ContentTypeVirtualizationClusterGroup, constants.ContentTypeVirtualizationClusterType, constants.ContentTypeVirtualizationVirtualMachine, constants.ContentTypeVirtualizationVMInterface},
	})
	if err != nil {
		return fmt.Errorf("add source_id custom field %s", err)
	}
	// Custom field for storing number of CPU cores for device (server).
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldHostCPUCoresName,
		Label:                 constants.CustomFieldHostCPUCoresLabel,
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		Default:               nil,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldHostCPUCoresDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ObjectTypes:           []constants.ContentType{constants.ContentTypeDcimDevice},
	})
	if err != nil {
		return fmt.Errorf("add host cpu cores custom field: %s", err)
	}
	// Custom field for storing the amount of the RAM on the device (server).
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldHostMemoryName,
		Label:                 constants.CustomFieldHostMemoryLabel,
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		Default:               nil,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldHostMemoryDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ObjectTypes:           []constants.ContentType{constants.ContentTypeDcimDevice},
	})
	if err != nil {
		return fmt.Errorf("add host memory custom field: %s", err)
	}
	// custom field for storing uuid of the device.
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldDeviceUUIDName,
		Label:                 constants.CustomFieldDeviceUUIDLabel,
		Type:                  objects.CustomFieldTypeText,
		Default:               nil,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldDeviceUUIDDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ObjectTypes:           []constants.ContentType{constants.ContentTypeDcimDevice},
	})
	if err != nil {
		return fmt.Errorf("add device uuid custom field: %s", err)
	}
	// Custom field for determining if an IP address was obtained from the arp table.
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldArpEntryName,
		Label:                 constants.CustomFieldArpEntryLabel,
		Type:                  objects.CustomFieldTypeBoolean,
		Default:               false,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldArpEntryDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ObjectTypes:           []constants.ContentType{constants.ContentTypeIpamIPAddress},
	})
	if err != nil {
		return fmt.Errorf("add arp entry custom field: %s", err)
	}
	return nil
}

// Collects all nbClusters from Netbox API and stores them in the NetBoxInventory.
func (nbi *NetboxInventory) initClusterGroups(ctx context.Context) error {
	nbClusterGroups, err := service.GetAll[objects.ClusterGroup](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of cluster groups by name
	nbi.clusterGroupsIndexByName = make(map[string]*objects.ClusterGroup)

	for i := range nbClusterGroups {
		clusterGroup := &nbClusterGroups[i]
		nbi.clusterGroupsIndexByName[clusterGroup.Name] = clusterGroup
		nbi.OrphanManager.AddItem(constants.ClusterGroupsAPIPath, clusterGroup)
	}
	nbi.Logger.Debug(ctx, "Successfully collected cluster groups from Netbox: ", nbi.clusterGroupsIndexByName)
	return nil
}

// Collects all ClusterTypes from Netbox API and stores them in the NetBoxInventory.
func (nbi *NetboxInventory) initClusterTypes(ctx context.Context) error {
	nbClusterTypes, err := service.GetAll[objects.ClusterType](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of cluster types by name
	nbi.clusterTypesIndexByName = make(map[string]*objects.ClusterType)
	for i := range nbClusterTypes {
		clusterType := &nbClusterTypes[i]
		nbi.clusterTypesIndexByName[clusterType.Name] = clusterType
		nbi.OrphanManager.AddItem(constants.ClusterTypesAPIPath, clusterType)
	}

	nbi.Logger.Debug(ctx, "Successfully collected cluster types from Netbox: ", nbi.clusterTypesIndexByName)
	return nil
}

// Collects all clusters from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initClusters(ctx context.Context) error {
	nbClusters, err := service.GetAll[objects.Cluster](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of clusters by name
	nbi.clustersIndexByName = make(map[string]*objects.Cluster)

	for i := range nbClusters {
		cluster := &nbClusters[i]
		nbi.clustersIndexByName[cluster.Name] = cluster
		nbi.OrphanManager.AddItem(constants.ClustersAPIPath, cluster)
	}

	nbi.Logger.Debug(ctx, "Successfully collected clusters from Netbox: ", nbi.clustersIndexByName)
	return nil
}

func (nbi *NetboxInventory) initDeviceTypes(ctx context.Context) error {
	nbDeviceTypes, err := service.GetAll[objects.DeviceType](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of device types by model
	nbi.deviceTypesIndexByModel = make(map[string]*objects.DeviceType)
	for i := range nbDeviceTypes {
		deviceType := &nbDeviceTypes[i]
		nbi.deviceTypesIndexByModel[deviceType.Model] = deviceType
		nbi.OrphanManager.AddItem(constants.DeviceTypesAPIPath, deviceType)
	}

	nbi.Logger.Debug(ctx, "Successfully collected device types from Netbox: ", nbi.deviceTypesIndexByModel)
	return nil
}

// Collects all interfaces from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initInterfaces(ctx context.Context) error {
	nbInterfaces, err := service.GetAll[objects.Interface](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of interfaces by device id and name
	nbi.interfacesIndexByDeviceIDAndName = make(map[int]map[string]*objects.Interface)

	for i := range nbInterfaces {
		intf := &nbInterfaces[i]
		if nbi.interfacesIndexByDeviceIDAndName[intf.Device.ID] == nil {
			nbi.interfacesIndexByDeviceIDAndName[intf.Device.ID] = make(map[string]*objects.Interface)
		}
		nbi.interfacesIndexByDeviceIDAndName[intf.Device.ID][intf.Name] = intf
		nbi.OrphanManager.AddItem(constants.InterfacesAPIPath, intf)
	}

	nbi.Logger.Debug(ctx, "Successfully collected interfaces from Netbox: ", nbi.interfacesIndexByDeviceIDAndName)
	return nil
}

// Collects all vlans from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initVlanGroups(ctx context.Context) error {
	nbVlanGroups, err := service.GetAll[objects.VlanGroup](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of vlans by name
	nbi.vlanGroupsIndexByName = make(map[string]*objects.VlanGroup)

	for i := range nbVlanGroups {
		vlanGroup := &nbVlanGroups[i]
		nbi.vlanGroupsIndexByName[vlanGroup.Name] = vlanGroup
		nbi.OrphanManager.AddItem(constants.VlanGroupsAPIPath, vlanGroup)
	}

	nbi.Logger.Debug(ctx, "Successfully collected vlans from Netbox: ", nbi.vlanGroupsIndexByName)
	return nil
}

// Collects all vlans from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initVlans(ctx context.Context) error {
	nbVlans, err := service.GetAll[objects.Vlan](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of vlans by VlanGroupId and Vid
	nbi.vlansIndexByVlanGroupIDAndVID = make(map[int]map[int]*objects.Vlan)

	for i := range nbVlans {
		vlan := &nbVlans[i]
		if vlan.Group == nil {
			// Update all existing vlans with default vlanGroup. This only happens
			// when there are predefined vlans in netbox. This is required because
			// vlans are indexed by vlan group.
			vlan.Group, err = nbi.CreateDefaultVlanGroupForVlan(nbi.Ctx, vlan.Site)
			if err != nil {
				return fmt.Errorf("create default vlan group for vlan: %s", err)
			}
			vlan, err = nbi.AddVlan(ctx, vlan)
			if err != nil {
				return err
			}
		}
		if nbi.vlansIndexByVlanGroupIDAndVID[vlan.Group.ID] == nil {
			nbi.vlansIndexByVlanGroupIDAndVID[vlan.Group.ID] = make(map[int]*objects.Vlan)
		}
		nbi.vlansIndexByVlanGroupIDAndVID[vlan.Group.ID][vlan.Vid] = vlan
		nbi.OrphanManager.AddItem(constants.VlansAPIPath, vlan)
	}

	nbi.Logger.Debug(ctx, "Successfully collected vlans from Netbox: ", nbi.vlansIndexByVlanGroupIDAndVID)
	return nil
}

// Collects all vms from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initVMs(ctx context.Context) error {
	nbVMs, err := service.GetAll[objects.VM](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of VMs by name and cluster id
	nbi.vmsIndexByNameAndClusterID = make(map[string]map[int]*objects.VM)

	for i := range nbVMs {
		vm := &nbVMs[i]
		if nbi.vmsIndexByNameAndClusterID[vm.Name] == nil {
			nbi.vmsIndexByNameAndClusterID[vm.Name] = make(map[int]*objects.VM)
		}
		if vm.Cluster == nil {
			nbi.vmsIndexByNameAndClusterID[vm.Name][-1] = vm
		} else {
			nbi.vmsIndexByNameAndClusterID[vm.Name][vm.Cluster.ID] = vm
		}
		nbi.OrphanManager.AddItem(constants.VirtualMachinesAPIPath, vm)
	}

	nbi.Logger.Debug(ctx, "Successfully collected VMs from Netbox: ", nbi.vmsIndexByNameAndClusterID)
	return nil
}

// Collects all VMInterfaces from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initVMInterfaces(ctx context.Context) error {
	nbVMInterfaces, err := service.GetAll[objects.VMInterface](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return fmt.Errorf("Init vm interfaces: %s", err)
	}

	// Initialize internal index of VM interfaces by VM id and name
	nbi.vmInterfacesIndexByVMIdAndName = make(map[int]map[string]*objects.VMInterface)
	for i := range nbVMInterfaces {
		vmIntf := &nbVMInterfaces[i]
		if nbi.vmInterfacesIndexByVMIdAndName[vmIntf.VM.ID] == nil {
			nbi.vmInterfacesIndexByVMIdAndName[vmIntf.VM.ID] = make(map[string]*objects.VMInterface)
		}
		nbi.vmInterfacesIndexByVMIdAndName[vmIntf.VM.ID][vmIntf.Name] = vmIntf
		nbi.OrphanManager.AddItem(constants.VMInterfacesAPIPath, vmIntf)
	}

	nbi.Logger.Debug(ctx, "Successfully collected VM interfaces from Netbox: ", nbi.vmInterfacesIndexByVMIdAndName)
	return nil
}

// Collects all IP addresses from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initIPAddresses(ctx context.Context) error {
	ipAddresses, err := service.GetAll[objects.IPAddress](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initializes internal index of IP addresses by address
	nbi.ipAdressesIndexByAddress = make(map[string]*objects.IPAddress)

	for i := range ipAddresses {
		ipAddr := &ipAddresses[i]
		if ipAddr.HasTag(nbi.SsotTag) {
			nbi.ipAdressesIndexByAddress[ipAddr.Address] = ipAddr
			nbi.OrphanManager.AddItem(constants.IPAddressesAPIPath, ipAddr)
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected IP addresses from Netbox: ", nbi.ipAdressesIndexByAddress)
	return nil
}

// Collects all Prefixes from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initPrefixes(ctx context.Context) error {
	prefixes, err := service.GetAll[objects.Prefix](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initializes internal index of prefixes by prefix
	nbi.prefixesIndexByPrefix = make(map[string]*objects.Prefix)

	for i := range prefixes {
		prefix := &prefixes[i]
		nbi.prefixesIndexByPrefix[prefix.Prefix] = prefix
		nbi.OrphanManager.AddItem(constants.PrefixesAPIPath, prefix)
	}

	nbi.Logger.Debug(ctx, "Successfully collected prefixes from Netbox: ", nbi.prefixesIndexByPrefix)
	return nil
}

// Collects all WirelessLANs from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initWirelessLANs(ctx context.Context) error {
	nbWirelessLans, err := service.GetAll[objects.WirelessLAN](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of WirelessLANs by SSID
	nbi.wirelessLANsIndexBySSID = make(map[string]*objects.WirelessLAN)

	for i := range nbWirelessLans {
		wirelessLan := &nbWirelessLans[i]
		nbi.wirelessLANsIndexBySSID[wirelessLan.SSID] = wirelessLan
		nbi.OrphanManager.AddItem(constants.WirelessLANsAPIPath, wirelessLan)
	}
	nbi.Logger.Debug(ctx, "Successfully collected wireless-lans from Netbox: ", nbi.wirelessLANsIndexBySSID)
	return nil
}

// Collects all WirelessLANGroups from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initWirelessLANGroups(ctx context.Context) error {
	nbWirelessLanGroups, err := service.GetAll[objects.WirelessLANGroup](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of WirelessLanGroups by SSID
	nbi.wirelessLANGroupsIndexByName = make(map[string]*objects.WirelessLANGroup)

	for i := range nbWirelessLanGroups {
		wirelessLanGroup := &nbWirelessLanGroups[i]
		nbi.wirelessLANGroupsIndexByName[wirelessLanGroup.Name] = wirelessLanGroup
		nbi.OrphanManager.AddItem(constants.WirelessLANGroupsAPIPath, wirelessLanGroup)
	}
	nbi.Logger.Debug(ctx, "Successfully collected wireless-lan-groups from Netbox: ", nbi.wirelessLANGroupsIndexByName)
	return nil
}
