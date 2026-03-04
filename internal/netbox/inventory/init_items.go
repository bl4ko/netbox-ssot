package inventory

import (
	"context"
	"fmt"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
	"github.com/src-doo/netbox-ssot/internal/netbox/service"
	"github.com/src-doo/netbox-ssot/internal/utils"
)

// Collect all tags from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initTags(ctx context.Context) error {
	extraArgs := fmt.Sprintf("&fields=%s", utils.ExtractJSONTagsFromStructIntoString(objects.Tag{}))
	nbTags, err := service.GetAll[objects.Tag](
		ctx,
		nbi.NetboxAPI,
		extraArgs,
	)
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
	ssotTag, err := nbi.AddTag(
		ctx,
		&objects.Tag{
			Name:        constants.SsotTagName,
			Slug:        constants.SsotTagName,
			Description: constants.SsotTagDescription,
			Color:       constants.SsotTagColor,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating default ssot  tag: %s", err)
	}

	nbi.SsotTag = ssotTag

	// Create default tag for orphaned objects
	orphanTag, err := nbi.AddTag(
		ctx,
		&objects.Tag{
			Name:        constants.OrphanTagName,
			Slug:        constants.OrphanTagName,
			Description: constants.OrphanTagDescription,
			Color:       constants.OrphanTagColor,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating default orphan tag: %s", err)
	}
	nbi.OrphanManager.Tag = orphanTag

	// Create default tag for device type overriding
	ignoreDeviceTypeTag, err := nbi.AddTag(
		ctx,
		&objects.Tag{
			Name:        constants.IgnoreDeviceTypeTagName,
			Slug:        constants.IgnoreDeviceTypeTagName,
			Description: constants.IgnoreDeviceTypeTagDescription,
			Color:       constants.IgnoreDeviceTypeTagColor,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating default ignore device type tag: %s", err)
	}

	nbi.IgnoreDeviceTypeTag = ignoreDeviceTypeTag
	return nil
}

// Collects all tenants from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initTenants(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Tenant{}),
	)
	nbTenants, err := service.GetAll[objects.Tenant](ctx, nbi.NetboxAPI, extraArgs)
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
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Contact{}),
	)
	nbContacts, err := service.GetAll[objects.Contact](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// We also create an index of contacts by name for easier access
	nbi.contactsIndexByName = make(map[string]*objects.Contact)
	for i := range nbContacts {
		contact := &nbContacts[i]
		nbi.contactsIndexByName[contact.Name] = contact
		nbi.OrphanManager.AddItem(contact)
	}
	nbi.Logger.Debug(ctx, "Successfully collected contacts from Netbox: ", nbi.contactsIndexByName)
	return nil
}

// Collects all contact roles from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initContactRoles(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.ContactRole{}),
	)
	nbContactRoles, err := service.GetAll[objects.ContactRole](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// We also create an index of contact roles by name for easier access
	nbi.contactRolesIndexByName = make(map[string]*objects.ContactRole)
	for i := range nbContactRoles {
		contactRole := &nbContactRoles[i]
		nbi.contactRolesIndexByName[contactRole.Name] = contactRole
	}
	nbi.Logger.Debug(
		ctx,
		"Successfully collected ContactRoles from Netbox: ",
		nbi.contactRolesIndexByName,
	)
	return nil
}

func (nbi *NetboxInventory) initContactAssignments(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.ContactAssignment{}),
	)
	nbCAs, err := service.GetAll[objects.ContactAssignment](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// We also create an index of contacts by name for easier access
	nbi.contactAssignmentsIndex = make(
		map[constants.ContentType]map[int]map[int]map[int]*objects.ContactAssignment,
	)
	debugIDs := map[int]bool{} // Netbox pagination bug duplicates
	for i := range nbCAs {
		cA := &nbCAs[i]
		if _, ok := debugIDs[cA.ID]; ok {
			fmt.Printf("Already been here: %d", cA.ID)
		}
		debugIDs[cA.ID] = true
		if nbi.contactAssignmentsIndex[cA.ModelType] == nil {
			nbi.contactAssignmentsIndex[cA.ModelType] = make(
				map[int]map[int]map[int]*objects.ContactAssignment,
			)
		}
		if nbi.contactAssignmentsIndex[cA.ModelType][cA.ObjectID] == nil {
			nbi.contactAssignmentsIndex[cA.ModelType][cA.ObjectID] = make(
				map[int]map[int]*objects.ContactAssignment,
			)
		}
		if nbi.contactAssignmentsIndex[cA.ModelType][cA.ObjectID][cA.Contact.ID] == nil {
			nbi.contactAssignmentsIndex[cA.ModelType][cA.ObjectID][cA.Contact.ID] = make(
				map[int]*objects.ContactAssignment,
			)
		}
		nbi.contactAssignmentsIndex[cA.ModelType][cA.ObjectID][cA.Contact.ID][cA.Role.ID] = cA
		nbi.OrphanManager.AddItem(cA)
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
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.ContactGroup{}),
	)
	nbContactGroups, err := service.GetAll[objects.ContactGroup](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// We also create an index of contact groups by name for easier access
	nbi.contactGroupsIndexByName = make(map[string]*objects.ContactGroup)
	for i := range nbContactGroups {
		contactGroup := &nbContactGroups[i]
		nbi.contactGroupsIndexByName[contactGroup.Name] = contactGroup
	}
	nbi.Logger.Debug(
		ctx,
		"Successfully collected ContactGroups from Netbox: ",
		nbi.contactGroupsIndexByName,
	)
	return nil
}

// Collects all sites from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initSites(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Site{}),
	)
	nbSites, err := service.GetAll[objects.Site](ctx, nbi.NetboxAPI, extraArgs)
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

// Collects all sites from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initSiteGroups(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.SiteGroup{}),
	)
	nbSiteGroups, err := service.GetAll[objects.SiteGroup](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// We also create an index of sites by name for easier access
	nbi.siteGroupsIndexByName = make(map[string]*objects.SiteGroup)
	for i := range nbSiteGroups {
		siteGroup := &nbSiteGroups[i]
		nbi.siteGroupsIndexByName[siteGroup.Name] = siteGroup
	}
	nbi.Logger.Debug(
		ctx,
		"Successfully collected SiteGroups from Netbox: ",
		nbi.siteGroupsIndexByName,
	)
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
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Manufacturer{}),
	)
	nbManufacturers, err := service.GetAll[objects.Manufacturer](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// Initialize internal index of manufacturers by name
	nbi.manufacturersIndexByName = make(map[string]*objects.Manufacturer)
	for i := range nbManufacturers {
		manufacturer := &nbManufacturers[i]
		nbi.manufacturersIndexByName[manufacturer.Name] = manufacturer
		nbi.OrphanManager.AddItem(manufacturer)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected manufacturers from Netbox: ",
		nbi.manufacturersIndexByName,
	)
	return nil
}

// Collects all platforms from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initPlatforms(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Platform{}),
	)
	nbPlatforms, err := service.GetAll[objects.Platform](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// Initialize internal index of platforms by name
	nbi.platformsIndexByName = make(map[string]*objects.Platform)

	for i, platform := range nbPlatforms {
		nbPlatform := &nbPlatforms[i]
		nbi.platformsIndexByName[platform.Name] = nbPlatform
		nbi.OrphanManager.AddItem(nbPlatform)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected platforms from Netbox: ",
		nbi.platformsIndexByName,
	)
	return nil
}

// Collect all devices from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initDevices(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Device{}),
	)
	nbDevices, err := service.GetAll[objects.Device](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// Initialize main index of devices by Name and SiteId
	nbi.devicesIndexByNameAndSiteID = make(map[string]map[int]*objects.Device)
	// Initialize helper index of devices by ID
	nbi.devicesIndexByID = make(map[int]*objects.Device)

	for i, device := range nbDevices {
		nbDevice := &nbDevices[i]
		nbi.devicesIndexByID[device.ID] = nbDevice
		if nbi.devicesIndexByNameAndSiteID[device.Name] == nil {
			nbi.devicesIndexByNameAndSiteID[device.Name] = make(map[int]*objects.Device)
		}
		nbi.devicesIndexByNameAndSiteID[device.Name][device.Site.ID] = nbDevice
		nbi.OrphanManager.AddItem(nbDevice)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected devices from Netbox: ",
		nbi.devicesIndexByNameAndSiteID,
	)
	return nil
}

// Collect all devices from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) initVirtualDeviceContexts(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.VirtualDeviceContext{}),
	)
	nbVirtualDeviceContexts, err := service.GetAll[objects.VirtualDeviceContext](
		ctx,
		nbi.NetboxAPI,
		extraArgs,
	)
	if err != nil {
		return err
	}
	// Initialize internal index of devices by Name and SiteId
	nbi.virtualDeviceContextsIndex = make(
		map[string]map[int]*objects.VirtualDeviceContext,
	)
	for i, virtualDeviceContext := range nbVirtualDeviceContexts {
		nbVirtualDeviceContext := &nbVirtualDeviceContexts[i]
		if nbi.virtualDeviceContextsIndex[virtualDeviceContext.Name] == nil {
			nbi.virtualDeviceContextsIndex[virtualDeviceContext.Name] = make(
				map[int]*objects.VirtualDeviceContext,
			)
		}
		nbi.virtualDeviceContextsIndex[virtualDeviceContext.Name][virtualDeviceContext.Device.ID] = nbVirtualDeviceContext
		nbi.OrphanManager.AddItem(nbVirtualDeviceContext)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected VirtualDeviceContexts from Netbox: ",
		nbi.virtualDeviceContextsIndex,
	)
	return nil
}

// Collects all deviceRoles from Netbox API and store them in the
// NetBoxInventory.
func (nbi *NetboxInventory) initDeviceRoles(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.DeviceRole{}),
	)
	nbDeviceRoles, err := service.GetAll[objects.DeviceRole](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// We also create an index of device roles by name for easier access
	nbi.deviceRolesIndexByName = make(map[string]*objects.DeviceRole)

	for i := range nbDeviceRoles {
		deviceRole := &nbDeviceRoles[i]
		nbi.deviceRolesIndexByName[deviceRole.Name] = deviceRole
		nbi.OrphanManager.AddItem(deviceRole)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected device roles from Netbox: ",
		nbi.deviceRolesIndexByName,
	)
	return nil
}

func (nbi *NetboxInventory) initCustomFields(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.CustomField{}),
	)
	customFields, err := service.GetAll[objects.CustomField](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// Initialize internal index of custom fields by name
	nbi.customFieldsIndexByName = make(map[string]*objects.CustomField, len(customFields))
	for i := range customFields {
		customField := &customFields[i]
		nbi.customFieldsIndexByName[customField.Name] = customField
	}
	nbi.Logger.Debug(
		ctx,
		"Successfully collected custom fields from Netbox: ",
		nbi.customFieldsIndexByName,
	)
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
		ObjectTypes: []constants.ContentType{
			constants.ContentTypeDcimDevice,
			constants.ContentTypeDcimDeviceRole,
			constants.ContentTypeDcimDeviceType,
			constants.ContentTypeDcimInterface,
			constants.ContentTypeDcimLocation,
			constants.ContentTypeDcimManufacturer,
			constants.ContentTypeDcimPlatform,
			constants.ContentTypeDcimRegion,
			constants.ContentTypeDcimSite,
			constants.ContentTypeDcimVirtualDeviceContext,
			constants.ContentTypeIpamIPAddress,
			constants.ContentTypeIpamVlanGroup,
			constants.ContentTypeIpamVlan,
			constants.ContentTypeIpamPrefix,
			constants.ContentTypeIpamVRF,
			constants.ContentTypeTenancyTenantGroup,
			constants.ContentTypeTenancyTenant,
			constants.ContentTypeTenancyContact,
			constants.ContentTypeTenancyContactAssignment,
			constants.ContentTypeTenancyContactGroup,
			constants.ContentTypeTenancyContactRole,
			constants.ContentTypeVirtualizationCluster,
			constants.ContentTypeVirtualizationClusterGroup,
			constants.ContentTypeVirtualizationClusterType,
			constants.ContentTypeVirtualizationVirtualMachine,
			constants.ContentTypeVirtualizationVMInterface,
			constants.ContentTypeWirelessLAN,
			constants.ContentTypeWirelessLANGroup,
			constants.ContentTypeVirtualizationVirtualDisk,
		},
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
		ObjectTypes: []constants.ContentType{
			constants.ContentTypeDcimDevice,
			constants.ContentTypeDcimDeviceRole,
			constants.ContentTypeDcimDeviceType,
			constants.ContentTypeDcimInterface,
			constants.ContentTypeDcimLocation,
			constants.ContentTypeDcimManufacturer,
			constants.ContentTypeDcimPlatform,
			constants.ContentTypeDcimRegion,
			constants.ContentTypeDcimSite,
			constants.ContentTypeDcimVirtualDeviceContext,
			constants.ContentTypeIpamIPAddress,
			constants.ContentTypeIpamVlanGroup,
			constants.ContentTypeIpamVlan,
			constants.ContentTypeIpamPrefix,
			constants.ContentTypeIpamVRF,
			constants.ContentTypeTenancyTenantGroup,
			constants.ContentTypeTenancyTenant,
			constants.ContentTypeTenancyContact,
			constants.ContentTypeTenancyContactAssignment,
			constants.ContentTypeTenancyContactGroup,
			constants.ContentTypeTenancyContactRole,
			constants.ContentTypeVirtualizationCluster,
			constants.ContentTypeVirtualizationClusterGroup,
			constants.ContentTypeVirtualizationClusterType,
			constants.ContentTypeVirtualizationVirtualMachine,
			constants.ContentTypeVirtualizationVMInterface,
			constants.ContentTypeWirelessLAN,
			constants.ContentTypeWirelessLANGroup,
			constants.ContentTypeDcimMACAddress,
			constants.ContentTypeVirtualizationVirtualDisk,
		},
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
		ObjectTypes: []constants.ContentType{
			constants.ContentTypeDcimDevice,
			constants.ContentTypeDcimDeviceRole,
			constants.ContentTypeDcimDeviceType,
			constants.ContentTypeDcimInterface,
			constants.ContentTypeDcimLocation,
			constants.ContentTypeDcimManufacturer,
			constants.ContentTypeDcimPlatform,
			constants.ContentTypeDcimRegion,
			constants.ContentTypeDcimSite,
			constants.ContentTypeDcimVirtualDeviceContext,
			constants.ContentTypeIpamIPAddress,
			constants.ContentTypeIpamVlanGroup,
			constants.ContentTypeIpamVlan,
			constants.ContentTypeIpamPrefix,
			constants.ContentTypeIpamVRF,
			constants.ContentTypeTenancyTenantGroup,
			constants.ContentTypeTenancyTenant,
			constants.ContentTypeTenancyContact,
			constants.ContentTypeTenancyContactAssignment,
			constants.ContentTypeTenancyContactGroup,
			constants.ContentTypeTenancyContactRole,
			constants.ContentTypeVirtualizationCluster,
			constants.ContentTypeVirtualizationClusterGroup,
			constants.ContentTypeVirtualizationClusterType,
			constants.ContentTypeVirtualizationVirtualMachine,
			constants.ContentTypeVirtualizationVMInterface,
			constants.ContentTypeVirtualizationVirtualDisk,
		},
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
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.ClusterGroup{}),
	)
	nbClusterGroups, err := service.GetAll[objects.ClusterGroup](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// Initialize internal index of cluster groups by name
	nbi.clusterGroupsIndexByName = make(map[string]*objects.ClusterGroup)

	for i := range nbClusterGroups {
		clusterGroup := &nbClusterGroups[i]
		nbi.clusterGroupsIndexByName[clusterGroup.Name] = clusterGroup
		nbi.OrphanManager.AddItem(clusterGroup)
	}
	nbi.Logger.Debug(
		ctx,
		"Successfully collected cluster groups from Netbox: ",
		nbi.clusterGroupsIndexByName,
	)
	return nil
}

// Collects all ClusterTypes from Netbox API and stores them in the NetBoxInventory.
func (nbi *NetboxInventory) initClusterTypes(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.ClusterType{}),
	)
	nbClusterTypes, err := service.GetAll[objects.ClusterType](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	// Initialize internal index of cluster types by name
	nbi.clusterTypesIndexByName = make(map[string]*objects.ClusterType)
	for i := range nbClusterTypes {
		clusterType := &nbClusterTypes[i]
		nbi.clusterTypesIndexByName[clusterType.Name] = clusterType
		nbi.OrphanManager.AddItem(clusterType)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected cluster types from Netbox: ",
		nbi.clusterTypesIndexByName,
	)
	return nil
}

// Collects all clusters from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initClusters(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Cluster{}),
	)
	nbClusters, err := service.GetAll[objects.Cluster](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	// Initialize internal index of clusters by name
	nbi.clustersIndexByName = make(map[string]*objects.Cluster)

	for i := range nbClusters {
		cluster := &nbClusters[i]
		nbi.clustersIndexByName[cluster.Name] = cluster
		nbi.OrphanManager.AddItem(cluster)
	}

	nbi.Logger.Debug(ctx, "Successfully collected clusters from Netbox: ", nbi.clustersIndexByName)
	return nil
}

func (nbi *NetboxInventory) initDeviceTypes(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.DeviceType{}),
	)
	nbDeviceTypes, err := service.GetAll[objects.DeviceType](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	// Initialize internal index of device types by model
	nbi.deviceTypesIndexByModel = make(map[string]*objects.DeviceType)
	for i := range nbDeviceTypes {
		deviceType := &nbDeviceTypes[i]
		nbi.deviceTypesIndexByModel[deviceType.Model] = deviceType
		nbi.OrphanManager.AddItem(deviceType)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected device types from Netbox: ",
		nbi.deviceTypesIndexByModel,
	)
	return nil
}

// Collects all interfaces from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initInterfaces(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Interface{}),
	)
	nbInterfaces, err := service.GetAll[objects.Interface](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	// Initialize main index of interfaces by device id and name
	nbi.interfacesIndexByDeviceIDAndName = make(map[int]map[string]*objects.Interface)
	// Initialize helper index for interfaces by ID
	nbi.interfacesIndexByID = make(map[int]*objects.Interface)

	for i := range nbInterfaces {
		intf := &nbInterfaces[i]
		nbi.interfacesIndexByID[intf.ID] = intf
		if nbi.interfacesIndexByDeviceIDAndName[intf.Device.ID] == nil {
			nbi.interfacesIndexByDeviceIDAndName[intf.Device.ID] = make(
				map[string]*objects.Interface,
			)
		}
		nbi.interfacesIndexByDeviceIDAndName[intf.Device.ID][intf.Name] = intf
		nbi.OrphanManager.AddItem(intf)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected interfaces from Netbox: ",
		nbi.interfacesIndexByDeviceIDAndName,
	)
	return nil
}

// Collects all vlans from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initVlanGroups(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.VlanGroup{}),
	)
	nbVlanGroups, err := service.GetAll[objects.VlanGroup](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	// Initialize internal index of vlans by name
	nbi.vlanGroupsIndexByName = make(map[string]*objects.VlanGroup)

	for i := range nbVlanGroups {
		vlanGroup := &nbVlanGroups[i]
		nbi.vlanGroupsIndexByName[vlanGroup.Name] = vlanGroup
		nbi.OrphanManager.AddItem(vlanGroup)
	}

	nbi.Logger.Debug(ctx, "Successfully collected vlans from Netbox: ", nbi.vlanGroupsIndexByName)
	return nil
}

// Collects all vlans from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initVlans(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Vlan{}),
	)
	nbVlans, err := service.GetAll[objects.Vlan](ctx, nbi.NetboxAPI, extraArgs)
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
			defaultVlanGroup, err := nbi.CreateDefaultVlanGroupForVlan(nbi.Ctx, vlan.Site)
			if err != nil {
				return fmt.Errorf("create default vlan group for vlan: %s", err)
			}
			vlan.Group = defaultVlanGroup
		}
		if nbi.vlansIndexByVlanGroupIDAndVID[vlan.Group.ID] == nil {
			nbi.vlansIndexByVlanGroupIDAndVID[vlan.Group.ID] = make(map[int]*objects.Vlan)
		}
		nbi.vlansIndexByVlanGroupIDAndVID[vlan.Group.ID][vlan.Vid] = vlan
		nbi.OrphanManager.AddItem(vlan)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected vlans from Netbox: ",
		nbi.vlansIndexByVlanGroupIDAndVID,
	)
	return nil
}

// Collects all vms from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initVMs(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.VM{}),
	)
	nbVMs, err := service.GetAll[objects.VM](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	// Initialize internal index of VMs by name and cluster id
	nbi.vmsIndexByNameAndClusterID = make(map[string]map[int]*objects.VM)
	nbi.vmsIndexByID = make(map[int]*objects.VM)

	for i := range nbVMs {
		vm := &nbVMs[i]
		nbi.vmsIndexByID[vm.ID] = vm
		if nbi.vmsIndexByNameAndClusterID[vm.Name] == nil {
			nbi.vmsIndexByNameAndClusterID[vm.Name] = make(map[int]*objects.VM)
		}
		if vm.Cluster == nil {
			nbi.vmsIndexByNameAndClusterID[vm.Name][-1] = vm
		} else {
			nbi.vmsIndexByNameAndClusterID[vm.Name][vm.Cluster.ID] = vm
		}
		nbi.OrphanManager.AddItem(vm)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected VMs from Netbox: ",
		nbi.vmsIndexByNameAndClusterID,
	)
	return nil
}

// Collects all VMInterfaces from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initVMInterfaces(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.VMInterface{}),
	)
	nbVMInterfaces, err := service.GetAll[objects.VMInterface](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return fmt.Errorf("Init vm interfaces: %s", err)
	}

	// Initialize internal index of VM interfaces by VM id and name
	nbi.vmInterfacesIndexByVMIdAndName = make(map[int]map[string]*objects.VMInterface)
	nbi.vmInterfacesIndexByID = make(map[int]*objects.VMInterface)
	for i := range nbVMInterfaces {
		vmIntf := &nbVMInterfaces[i]
		nbi.vmInterfacesIndexByID[vmIntf.ID] = vmIntf
		if nbi.vmInterfacesIndexByVMIdAndName[vmIntf.VM.ID] == nil {
			nbi.vmInterfacesIndexByVMIdAndName[vmIntf.VM.ID] = make(map[string]*objects.VMInterface)
		}
		nbi.vmInterfacesIndexByVMIdAndName[vmIntf.VM.ID][vmIntf.Name] = vmIntf
		nbi.OrphanManager.AddItem(vmIntf)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected VM interfaces from Netbox: ",
		nbi.vmInterfacesIndexByVMIdAndName,
	)
	return nil
}

// Collects all IP addresses from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initIPAddresses(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.IPAddress{}),
	)
	ipAddresses, err := service.GetAll[objects.IPAddress](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	nbi.ipAddressesIndex = make(
		map[constants.ContentType]map[string]map[string]map[string]*objects.IPAddress,
	)
	for i := range ipAddresses {
		ipAddr := &ipAddresses[i]
		ifaceType, ifaceName, ifaceParentName, err := nbi.getIndexValuesForIPAddress(ipAddr)
		if err != nil {
			return fmt.Errorf("get index values for ip address: %s", err)
		}
		// Skip IP addresses whose assigned interface is not in inventory
		if ipAddr.AssignedObjectType != "" && ifaceType == "" {
			continue
		}
		nbi.verifyIPAddressIndexExists(ifaceType, ifaceName, ifaceParentName)
		indexKey := ipAddressIndexKey(ipAddr)
		nbi.ipAddressesIndex[ifaceType][ifaceName][ifaceParentName][indexKey] = ipAddr
		nbi.OrphanManager.AddItem(ipAddr)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected IP addresses from Netbox: ",
		nbi.ipAddressesIndex,
	)
	return nil
}

func (nbi *NetboxInventory) initMACAddresses(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.MACAddress{}),
	)
	nbMACAddresses, err := service.GetAll[objects.MACAddress](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}
	// Initializes internal indexx
	nbi.macAddressesIndex = make(
		map[constants.ContentType]map[string]map[string]map[string]*objects.MACAddress,
	)
	for i := range nbMACAddresses {
		macAddress := &nbMACAddresses[i]
		ifaceType, ifaceName, ifaceParentName, err := nbi.getIndexValuesForMACAddress(
			macAddress,
		)
		if err != nil {
			return fmt.Errorf("get index values for mac address: %s", err)
		}
		// Skip MAC addresses whose assigned interface is not in inventory
		if macAddress.AssignedObjectType != "" && ifaceType == "" {
			continue
		}
		nbi.verifyMACAddressIndexExists(ifaceType, ifaceName, ifaceParentName)
		nbi.macAddressesIndex[ifaceType][ifaceName][ifaceParentName][macAddress.MAC] = macAddress
		nbi.OrphanManager.AddItem(macAddress)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected MAC addresses from Netbox: ",
		nbi.macAddressesIndex,
	)
	return nil
}

// Collects all Prefixes from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initPrefixes(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.Prefix{}),
	)
	prefixes, err := service.GetAll[objects.Prefix](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	nbi.prefixesIndexByPrefix = make(map[string]map[int]*objects.Prefix)

	for i := range prefixes {
		prefix := &prefixes[i]
		vrfID := 0
		if prefix.VRF != nil {
			vrfID = prefix.VRF.ID
		}
		if nbi.prefixesIndexByPrefix[prefix.Prefix] == nil {
			nbi.prefixesIndexByPrefix[prefix.Prefix] = make(map[int]*objects.Prefix)
		}
		nbi.prefixesIndexByPrefix[prefix.Prefix][vrfID] = prefix
		nbi.OrphanManager.AddItem(prefix)
	}

	nbi.Logger.Debug(
		ctx,
		"Successfully collected prefixes from Netbox: ",
		nbi.prefixesIndexByPrefix,
	)
	return nil
}

// Collects all WirelessLANs from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initWirelessLANs(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.WirelessLAN{}),
	)
	nbWirelessLans, err := service.GetAll[objects.WirelessLAN](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	// Initialize internal index of WirelessLANs by SSID
	nbi.wirelessLANsIndexBySSID = make(map[string]*objects.WirelessLAN)

	for i := range nbWirelessLans {
		wirelessLan := &nbWirelessLans[i]
		nbi.wirelessLANsIndexBySSID[wirelessLan.SSID] = wirelessLan
		nbi.OrphanManager.AddItem(wirelessLan)
	}
	nbi.Logger.Debug(
		ctx,
		"Successfully collected wireless-lans from Netbox: ",
		nbi.wirelessLANsIndexBySSID,
	)
	return nil
}

// Collects all WirelessLANGroups from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) initWirelessLANGroups(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.WirelessLANGroup{}),
	)
	nbWirelessLanGroups, err := service.GetAll[objects.WirelessLANGroup](
		ctx,
		nbi.NetboxAPI,
		extraArgs,
	)
	if err != nil {
		return err
	}

	// Initialize internal index of WirelessLanGroups by SSID
	nbi.wirelessLANGroupsIndexByName = make(map[string]*objects.WirelessLANGroup)

	for i := range nbWirelessLanGroups {
		wirelessLanGroup := &nbWirelessLanGroups[i]
		nbi.wirelessLANGroupsIndexByName[wirelessLanGroup.Name] = wirelessLanGroup
		nbi.OrphanManager.AddItem(wirelessLanGroup)
	}
	nbi.Logger.Debug(
		ctx,
		"Successfully collected wireless-lan-groups from Netbox: ",
		nbi.wirelessLANGroupsIndexByName,
	)
	return nil
}

// initVirtualDisks collects all virtual disks from Netbox API
// and stores them to local inventory.
func (nbi *NetboxInventory) initVirtualDisks(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.VirtualDisk{}),
	)
	nbVirtualDisks, err := service.GetAll[objects.VirtualDisk](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return err
	}

	// Initialize internal index of virtual disks by name and VM id
	nbi.virtualDisksIndexByVMIDAndName = make(map[int]map[string]*objects.VirtualDisk)

	for i := range nbVirtualDisks {
		virtualDisk := &nbVirtualDisks[i]
		if nbi.virtualDisksIndexByVMIDAndName[virtualDisk.VM.ID] == nil {
			nbi.virtualDisksIndexByVMIDAndName[virtualDisk.VM.ID] = make(
				map[string]*objects.VirtualDisk,
			)
		}
		nbi.virtualDisksIndexByVMIDAndName[virtualDisk.VM.ID][virtualDisk.Name] = virtualDisk
		nbi.OrphanManager.AddItem(virtualDisk)
	}
	return nil
}

// initVRFs collects all VRF from Netbox API
// and stores them to local inventory.
func (nbi *NetboxInventory) initVRFs(ctx context.Context) error {
	extraArgs := fmt.Sprintf(
		"&fields=%s",
		utils.ExtractJSONTagsFromStructIntoString(objects.VRF{}),
	)
	nbVRFs, err := service.GetAll[objects.VRF](ctx, nbi.NetboxAPI, extraArgs)
	if err != nil {
		return fmt.Errorf("get all vrfs: %s", err)
	}
	nbi.vrfsIndexByName = make(map[string]*objects.VRF, len(nbVRFs))
	for i := range nbVRFs {
		vrf := &nbVRFs[i]
		nbi.vrfsIndexByName[vrf.Name] = vrf
		nbi.OrphanManager.AddItem(vrf)
	}
	nbi.Logger.Debug(
	ctx,
	"Successfully collected VRF from Netbox: ",
	nbi.vrfsIndexByName,
	)
	return nil
}
