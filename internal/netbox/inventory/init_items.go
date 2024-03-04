package inventory

import (
	"context"
	"fmt"
	"slices"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Collect all tags from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) InitTags(ctx context.Context) error {
	nbTags, err := service.GetAll[objects.Tag](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	nbi.Tags = make([]*objects.Tag, len(nbTags))
	for i := range nbTags {
		nbi.Tags[i] = &nbTags[i]
	}
	nbi.Logger.Debug(ctx, "Successfully collected tags from Netbox: ", nbi.Tags)

	// Custom tag for all netbox objects
	ssotTags, err := service.GetAll[objects.Tag](ctx, nbi.NetboxAPI, "&name=netbox-ssot")
	if err != nil {
		return err
	}
	if len(ssotTags) == 0 {
		nbi.Logger.Info(ctx, "Tag netbox-ssot not found in Netbox. Creating it now...")
		newTag := objects.Tag{Name: "netbox-ssot", Slug: "netbox-ssot", Description: "Tag used by netbox-ssot to mark devices that are managed by it", Color: "00add8"}
		ssotTag, err := service.Create[objects.Tag](ctx, nbi.NetboxAPI, &newTag)
		if err != nil {
			return err
		}
		nbi.SsotTag = ssotTag
	} else {
		nbi.SsotTag = &ssotTags[0]
	}
	return nil
}

// Collects all tenants from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) InitTenants(ctx context.Context) error {
	nbTenants, err := service.GetAll[objects.Tenant](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of tenants by name for easier access
	nbi.TenantsIndexByName = make(map[string]*objects.Tenant)
	for i := range nbTenants {
		tenant := &nbTenants[i]
		nbi.TenantsIndexByName[tenant.Name] = tenant
	}
	nbi.Logger.Debug(ctx, "Successfully collected tenants from Netbox: ", nbi.TenantsIndexByName)
	return nil
}

// Collects all contacts from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) InitContacts(ctx context.Context) error {
	nbContacts, err := service.GetAll[objects.Contact](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of contacts by name for easier access
	nbi.ContactsIndexByName = make(map[string]*objects.Contact)
	nbi.OrphanManager[service.ContactsAPIPath] = make(map[int]bool, len(nbContacts))
	for i := range nbContacts {
		contact := &nbContacts[i]
		nbi.ContactsIndexByName[contact.Name] = contact
		if slices.IndexFunc(contact.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.ContactsAPIPath][contact.ID] = true
		}
	}
	nbi.Logger.Debug(ctx, "Successfully collected contacts from Netbox: ", nbi.ContactsIndexByName)
	return nil
}

// Collects all contact roles from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) InitContactRoles(ctx context.Context) error {
	nbContactRoles, err := service.GetAll[objects.ContactRole](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of contact roles by name for easier access
	nbi.ContactRolesIndexByName = make(map[string]*objects.ContactRole)
	for i := range nbContactRoles {
		contactRole := &nbContactRoles[i]
		nbi.ContactRolesIndexByName[contactRole.Name] = contactRole
	}
	nbi.Logger.Debug(ctx, "Successfully collected ContactRoles from Netbox: ", nbi.ContactRolesIndexByName)
	return nil
}

func (nbi *NetboxInventory) InitContactAssignments(ctx context.Context) error {
	nbCAs, err := service.GetAll[objects.ContactAssignment](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of contacts by name for easier access
	nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID = make(map[string]map[int]map[int]map[int]*objects.ContactAssignment)
	nbi.OrphanManager[service.ContactAssignmentsAPIPath] = make(map[int]bool, len(nbCAs))
	debugIDs := map[int]bool{} // Netbox pagination bug duplicates
	for i := range nbCAs {
		cA := &nbCAs[i]
		if _, ok := debugIDs[cA.ID]; ok {
			fmt.Printf("Already been here: %d", cA.ID)
		}
		debugIDs[cA.ID] = true
		if nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[cA.ContentType] == nil {
			nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[cA.ContentType] = make(map[int]map[int]map[int]*objects.ContactAssignment)
		}
		if nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[cA.ContentType][cA.ObjectID] == nil {
			nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[cA.ContentType][cA.ObjectID] = make(map[int]map[int]*objects.ContactAssignment)
		}
		if nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[cA.ContentType][cA.ObjectID][cA.Contact.ID] == nil {
			nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[cA.ContentType][cA.ObjectID][cA.Contact.ID] = make(map[int]*objects.ContactAssignment)
		}
		nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[cA.ContentType][cA.ObjectID][cA.Contact.ID][cA.Role.ID] = cA
		if slices.IndexFunc(cA.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.ContactAssignmentsAPIPath][cA.ID] = true
		}
	}
	nbi.Logger.Debug(ctx, "Successfully collected contacts from Netbox: ", nbi.ContactsIndexByName)
	return nil
}

// Initializes default admin contact role used for adding admin contacts of vms.
func (nbi *NetboxInventory) InitAdminContactRole(ctx context.Context) error {
	_, err := nbi.AddContactRole(ctx, &objects.ContactRole{
		NetboxObject: objects.NetboxObject{
			Description: "Auto generated contact role by netbox-ssot for admins of vms.",
			CustomFields: map[string]string{
				constants.CustomFieldSourceName: nbi.SsotTag.Name,
			},
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
func (nbi *NetboxInventory) InitContactGroups(ctx context.Context) error {
	nbContactGroups, err := service.GetAll[objects.ContactGroup](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of contact groups by name for easier access
	nbi.ContactGroupsIndexByName = make(map[string]*objects.ContactGroup)
	for i := range nbContactGroups {
		contactGroup := &nbContactGroups[i]
		nbi.ContactGroupsIndexByName[contactGroup.Name] = contactGroup
	}
	nbi.Logger.Debug(ctx, "Successfully collected ContactGroups from Netbox: ", nbi.ContactGroupsIndexByName)
	return nil
}

// Collects all sites from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) InitSites(ctx context.Context) error {
	nbSites, err := service.GetAll[objects.Site](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of sites by name for easier access
	nbi.SitesIndexByName = make(map[string]*objects.Site)
	for i := range nbSites {
		site := &nbSites[i]
		nbi.SitesIndexByName[site.Name] = site
	}
	nbi.Logger.Debug(ctx, "Successfully collected sites from Netbox: ", nbi.SitesIndexByName)
	return nil
}

// Collects all manufacturers from Netbox API and store them in NetBoxInventory.
func (nbi *NetboxInventory) InitManufacturers(ctx context.Context) error {
	nbManufacturers, err := service.GetAll[objects.Manufacturer](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of manufacturers by name
	nbi.ManufacturersIndexByName = make(map[string]*objects.Manufacturer)
	// OrphanManager takes care of all manufacturers created by netbox-ssot
	nbi.OrphanManager[service.ManufacturersAPIPath] = make(map[int]bool)

	for i := range nbManufacturers {
		manufacturer := &nbManufacturers[i]
		nbi.ManufacturersIndexByName[manufacturer.Name] = manufacturer
		if slices.IndexFunc(manufacturer.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.ManufacturersAPIPath][manufacturer.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected manufacturers from Netbox: ", nbi.ManufacturersIndexByName)
	return nil
}

// Collects all platforms from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) InitPlatforms(ctx context.Context) error {
	nbPlatforms, err := service.GetAll[objects.Platform](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of platforms by name
	nbi.PlatformsIndexByName = make(map[string]*objects.Platform)
	// OrphanManager takes care of all platforms created by netbox-ssot
	nbi.OrphanManager[service.PlatformsAPIPath] = make(map[int]bool, 0)

	for i, platform := range nbPlatforms {
		nbi.PlatformsIndexByName[platform.Name] = &nbPlatforms[i]
		if slices.IndexFunc(platform.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.PlatformsAPIPath][platform.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected platforms from Netbox: ", nbi.PlatformsIndexByName)
	return nil
}

// Collect all devices from Netbox API and store them in the NetBoxInventory.
func (nbi *NetboxInventory) InitDevices(ctx context.Context) error {
	nbDevices, err := service.GetAll[objects.Device](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of devices by Name and SiteId
	nbi.DevicesIndexByNameAndSiteID = make(map[string]map[int]*objects.Device)
	// OrphanManager takes care of all devices created by netbox-ssot
	nbi.OrphanManager[service.DevicesAPIPath] = make(map[int]bool)

	for i, device := range nbDevices {
		if nbi.DevicesIndexByNameAndSiteID[device.Name] == nil {
			nbi.DevicesIndexByNameAndSiteID[device.Name] = make(map[int]*objects.Device)
		}
		nbi.DevicesIndexByNameAndSiteID[device.Name][device.Site.ID] = &nbDevices[i]
		if slices.IndexFunc(device.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.DevicesAPIPath][device.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected devices from Netbox: ", nbi.DevicesIndexByNameAndSiteID)
	return nil
}

// Collects all deviceRoles from Netbox API and store them in the
// NetBoxInventory.
func (nbi *NetboxInventory) InitDeviceRoles(ctx context.Context) error {
	nbDeviceRoles, err := service.GetAll[objects.DeviceRole](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// We also create an index of device roles by name for easier access
	nbi.DeviceRolesIndexByName = make(map[string]*objects.DeviceRole)
	// OrphanManager takes care of all device roles created by netbox-ssot
	nbi.OrphanManager[service.DeviceRolesAPIPath] = make(map[int]bool, 0)

	for i := range nbDeviceRoles {
		deviceRole := &nbDeviceRoles[i]
		nbi.DeviceRolesIndexByName[deviceRole.Name] = deviceRole
		if slices.IndexFunc(deviceRole.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.DeviceRolesAPIPath][deviceRole.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected device roles from Netbox: ", nbi.DeviceRolesIndexByName)
	return nil
}

// Ensures that attribute ServerDeviceRole is proper initialized.
func (nbi *NetboxInventory) InitServerDeviceRole(ctx context.Context) error {
	_, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{Name: "Server", Slug: "server", Color: "00add8", VMRole: true})
	if err != nil {
		return err
	}
	return nil
}

func (nbi *NetboxInventory) InitCustomFields(ctx context.Context) error {
	customFields, err := service.GetAll[objects.CustomField](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of custom fields by name
	nbi.CustomFieldsIndexByName = make(map[string]*objects.CustomField, len(customFields))
	for i := range customFields {
		customField := &customFields[i]
		nbi.CustomFieldsIndexByName[customField.Name] = customField
	}
	nbi.Logger.Debug(ctx, "Successfully collected custom fields from Netbox: ", nbi.CustomFieldsIndexByName)
	return nil
}

// This function Initializes all custom fields required for servers and other objects
// Currently these are two:
// - host_cpu_cores
// - host_memory
// - sourceId - this is used to store the ID of the source object in Netbox (interfaces).
func (nbi *NetboxInventory) InitSsotCustomFields(ctx context.Context) error {
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
		ContentTypes:          []string{"dcim.device", "dcim.devicerole", "dcim.devicetype", "dcim.interface", "dcim.location", "dcim.manufacturer", "dcim.platform", "dcim.region", "dcim.site", "ipam.ipaddress", "ipam.vlangroup", "ipam.vlan", "ipam.prefix", "tenancy.tenantgroup", "tenancy.tenant", "tenancy.contact", "tenancy.contactassignment", "tenancy.contactgroup", "tenancy.contactrole", "virtualization.cluster", "virtualization.clustergroup", "virtualization.clustertype", "virtualization.virtualmachine", "virtualization.vminterface"},
	})
	if err != nil {
		return err
	}
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldSourceIDName,
		Label:                 constants.CustomFieldSourceIDLabel,
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldSourceIDDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ContentTypes:          []string{"dcim.device", "dcim.devicerole", "dcim.devicetype", "dcim.interface", "dcim.location", "dcim.manufacturer", "dcim.platform", "dcim.region", "dcim.site", "ipam.ipaddress", "ipam.vlangroup", "ipam.vlan", "ipam.prefix", "tenancy.tenantgroup", "tenancy.tenant", "tenancy.contact", "tenancy.contactassignment", "tenancy.contactgroup", "tenancy.contactrole", "virtualization.cluster", "virtualization.clustergroup", "virtualization.clustertype", "virtualization.virtualmachine", "virtualization.vminterface"},
	})
	if err != nil {
		return err
	}
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldHostCPUCoresName,
		Label:                 constants.CustomFieldHostCPUCoresLabel,
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldHostCPUCoresDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ContentTypes:          []string{"dcim.device"},
	})
	if err != nil {
		return err
	}
	_, err = nbi.AddCustomField(ctx, &objects.CustomField{
		Name:                  constants.CustomFieldHostMemoryName,
		Label:                 constants.CustomFieldHostMemoryLabel,
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldHostMemoryDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ContentTypes:          []string{"dcim.device"},
	})
	if err != nil {
		return err
	}
	return nil
}

// Collects all nbClusters from Netbox API and stores them in the NetBoxInventory.
func (nbi *NetboxInventory) InitClusterGroups(ctx context.Context) error {
	nbClusterGroups, err := service.GetAll[objects.ClusterGroup](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}
	// Initialize internal index of cluster groups by name
	nbi.ClusterGroupsIndexByName = make(map[string]*objects.ClusterGroup)
	// OrphanManager takes care of all cluster groups created by netbox-ssot
	nbi.OrphanManager[service.ClusterGroupsAPIPath] = make(map[int]bool, 0)

	for i := range nbClusterGroups {
		clusterGroup := &nbClusterGroups[i]
		nbi.ClusterGroupsIndexByName[clusterGroup.Name] = clusterGroup
		if slices.IndexFunc(clusterGroup.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.ClusterGroupsAPIPath][clusterGroup.ID] = true
		}
	}
	nbi.Logger.Debug(ctx, "Successfully collected cluster groups from Netbox: ", nbi.ClusterGroupsIndexByName)
	return nil
}

// Collects all ClusterTypes from Netbox API and stores them in the NetBoxInventory.
func (nbi *NetboxInventory) InitClusterTypes(ctx context.Context) error {
	nbClusterTypes, err := service.GetAll[objects.ClusterType](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of cluster types by name
	nbi.ClusterTypesIndexByName = make(map[string]*objects.ClusterType)
	// OrphanManager takes care of all cluster types created by netbox-ssot
	nbi.OrphanManager[service.ClusterTypesAPIPath] = make(map[int]bool, 0)

	for i := range nbClusterTypes {
		clusterType := &nbClusterTypes[i]
		nbi.ClusterTypesIndexByName[clusterType.Name] = clusterType
		if slices.IndexFunc(clusterType.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.ClusterTypesAPIPath][clusterType.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected cluster types from Netbox: ", nbi.ClusterTypesIndexByName)
	return nil
}

// Collects all clusters from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) InitClusters(ctx context.Context) error {
	nbClusters, err := service.GetAll[objects.Cluster](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of clusters by name
	nbi.ClustersIndexByName = make(map[string]*objects.Cluster)
	// OrphanManager takes care of all clusters created by netbox-ssot
	nbi.OrphanManager[service.ClustersAPIPath] = make(map[int]bool, 0)

	for i := range nbClusters {
		cluster := &nbClusters[i]
		nbi.ClustersIndexByName[cluster.Name] = cluster
		if slices.IndexFunc(cluster.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.ClustersAPIPath][cluster.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected clusters from Netbox: ", nbi.ClustersIndexByName)
	return nil
}

func (nbi *NetboxInventory) InitDeviceTypes(ctx context.Context) error {
	nbDeviceTypes, err := service.GetAll[objects.DeviceType](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of device types by model
	nbi.DeviceTypesIndexByModel = make(map[string]*objects.DeviceType)
	// OrphanManager takes care of all device types created by netbox-ssot
	nbi.OrphanManager[service.DeviceTypesAPIPath] = make(map[int]bool, 0)

	for i := range nbDeviceTypes {
		deviceType := &nbDeviceTypes[i]
		nbi.DeviceTypesIndexByModel[deviceType.Model] = deviceType
		if slices.IndexFunc(deviceType.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.DeviceTypesAPIPath][deviceType.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected device types from Netbox: ", nbi.DeviceTypesIndexByModel)
	return nil
}

// Collects all interfaces from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) InitInterfaces(ctx context.Context) error {
	nbInterfaces, err := service.GetAll[objects.Interface](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of interfaces by device id and name
	nbi.InterfacesIndexByDeviceIDAndName = make(map[int]map[string]*objects.Interface)
	// OrphanManager takes care of all interfaces created by netbox-ssot
	nbi.OrphanManager[service.InterfacesAPIPath] = make(map[int]bool, 0)

	for i := range nbInterfaces {
		intf := &nbInterfaces[i]
		if nbi.InterfacesIndexByDeviceIDAndName[intf.Device.ID] == nil {
			nbi.InterfacesIndexByDeviceIDAndName[intf.Device.ID] = make(map[string]*objects.Interface)
		}
		nbi.InterfacesIndexByDeviceIDAndName[intf.Device.ID][intf.Name] = intf
		if slices.IndexFunc(intf.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.InterfacesAPIPath][intf.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected interfaces from Netbox: ", nbi.InterfacesIndexByDeviceIDAndName)
	return nil
}

// Inits default VlanGroup, which is required to group all Vlans that are not part of other
// vlangroups into it. Each vlan is indexed by their (vlanGroup, vid).
func (nbi *NetboxInventory) InitDefaultVlanGroup(ctx context.Context) error {
	_, err := nbi.AddVlanGroup(ctx, &objects.VlanGroup{
		NetboxObject: objects.NetboxObject{
			Tags:        []*objects.Tag{nbi.SsotTag},
			Description: "Default netbox-ssot VlanGroup for all vlans that are not part of any other vlanGroup. This group is required for netbox-ssot vlan index to work.",
			CustomFields: map[string]string{
				constants.CustomFieldSourceName: nbi.SsotTag.Name,
			},
		},
		Name:   objects.DefaultVlanGroupName,
		Slug:   utils.Slugify(objects.DefaultVlanGroupName),
		MinVid: 1,
		MaxVid: constants.MaxVID,
	})
	if err != nil {
		return fmt.Errorf("init default vlan group: %s", err)
	}
	return nil
}

// Collects all vlans from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) InitVlanGroups(ctx context.Context) error {
	nbVlanGroups, err := service.GetAll[objects.VlanGroup](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of vlans by name
	nbi.VlanGroupsIndexByName = make(map[string]*objects.VlanGroup)
	// Add VlanGroups to orphan manager
	nbi.OrphanManager[service.VlanGroupsAPIPath] = make(map[int]bool, 0)

	for i := range nbVlanGroups {
		vlanGroup := &nbVlanGroups[i]
		nbi.VlanGroupsIndexByName[vlanGroup.Name] = vlanGroup
		if slices.IndexFunc(vlanGroup.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.VlanGroupsAPIPath][vlanGroup.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected vlans from Netbox: ", nbi.VlanGroupsIndexByName)
	return nil
}

// Collects all vlans from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) InitVlans(ctx context.Context) error {
	nbVlans, err := service.GetAll[objects.Vlan](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of vlans by VlanGroupId and Vid
	nbi.VlansIndexByVlanGroupIDAndVID = make(map[int]map[int]*objects.Vlan)
	// Add vlans to orphan manager
	nbi.OrphanManager[service.VlansAPIPath] = make(map[int]bool, 0)

	for i := range nbVlans {
		vlan := &nbVlans[i]
		if vlan.Group == nil {
			// Update all existing vlans with default vlanGroup. This only happens
			// when there are predefined vlans in netbox.
			vlan.Group = nbi.VlanGroupsIndexByName[objects.DefaultVlanGroupName] // This should not fail, because InitDefaultVlanGroup executes before InitVlans
			vlan, err = nbi.AddVlan(ctx, vlan)
			if err != nil {
				return err
			}
		}
		if nbi.VlansIndexByVlanGroupIDAndVID[vlan.Group.ID] == nil {
			nbi.VlansIndexByVlanGroupIDAndVID[vlan.Group.ID] = make(map[int]*objects.Vlan)
		}
		nbi.VlansIndexByVlanGroupIDAndVID[vlan.Group.ID][vlan.Vid] = vlan
		if slices.IndexFunc(vlan.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.VlansAPIPath][vlan.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected vlans from Netbox: ", nbi.VlansIndexByVlanGroupIDAndVID)
	return nil
}

// Collects all vms from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) InitVMs(ctx context.Context) error {
	nbVMs, err := service.GetAll[objects.VM](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initialize internal index of VMs by name
	nbi.VMsIndexByName = make(map[string]*objects.VM)
	// Add VMs to orphan manager
	nbi.OrphanManager[service.VirtualMachinesAPIPath] = make(map[int]bool, 0)

	for i := range nbVMs {
		vm := &nbVMs[i]
		nbi.VMsIndexByName[vm.Name] = vm
		if slices.IndexFunc(vm.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.VirtualMachinesAPIPath][vm.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected VMs from Netbox: ", nbi.VMsIndexByName)
	return nil
}

// Collects all VMInterfaces from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) InitVMInterfaces(ctx context.Context) error {
	nbVMInterfaces, err := service.GetAll[objects.VMInterface](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return fmt.Errorf("Init vm interfaces: %s", err)
	}

	// Initialize internal index of VM interfaces by VM id and name
	nbi.VMInterfacesIndexByVMIdAndName = make(map[int]map[string]*objects.VMInterface)
	// Add VMInterfaces to orphan manager
	nbi.OrphanManager[service.VMInterfacesAPIPath] = make(map[int]bool)

	for i := range nbVMInterfaces {
		vmIntf := &nbVMInterfaces[i]
		if nbi.VMInterfacesIndexByVMIdAndName[vmIntf.VM.ID] == nil {
			nbi.VMInterfacesIndexByVMIdAndName[vmIntf.VM.ID] = make(map[string]*objects.VMInterface)
		}
		nbi.VMInterfacesIndexByVMIdAndName[vmIntf.VM.ID][vmIntf.Name] = vmIntf
		if slices.IndexFunc(vmIntf.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.VMInterfacesAPIPath][vmIntf.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected VM interfaces from Netbox: ", nbi.VMInterfacesIndexByVMIdAndName)
	return nil
}

// Collects all IP addresses from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) InitIPAddresses(ctx context.Context) error {
	ipAddresses, err := service.GetAll[objects.IPAddress](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initializes internal index of IP addresses by address
	nbi.IPAdressesIndexByAddress = make(map[string]*objects.IPAddress)
	// Add IP addresses to orphan manager
	nbi.OrphanManager[service.IPAddressesAPIPath] = make(map[int]bool, 0)

	for i := range ipAddresses {
		ipAddr := &ipAddresses[i]
		nbi.IPAdressesIndexByAddress[ipAddr.Address] = ipAddr
		if slices.IndexFunc(ipAddr.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.IPAddressesAPIPath][ipAddr.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected IP addresses from Netbox: ", nbi.IPAdressesIndexByAddress)
	return nil
}

// Collects all Prefixes from Netbox API and stores them to local inventory.
func (nbi *NetboxInventory) InitPrefixes(ctx context.Context) error {
	prefixes, err := service.GetAll[objects.Prefix](ctx, nbi.NetboxAPI, "")
	if err != nil {
		return err
	}

	// Initializes internal index of prefixes by prefix
	nbi.PrefixesIndexByPrefix = make(map[string]*objects.Prefix)
	// Add prefixes to orphan manager
	nbi.OrphanManager[service.PrefixesAPIPath] = make(map[int]bool, 0)

	for i := range prefixes {
		prefix := &prefixes[i]
		nbi.PrefixesIndexByPrefix[prefix.Prefix] = prefix
		if slices.IndexFunc(prefix.Tags, func(t *objects.Tag) bool { return t.Slug == nbi.SsotTag.Slug }) >= 0 {
			nbi.OrphanManager[service.PrefixesAPIPath][prefix.ID] = true
		}
	}

	nbi.Logger.Debug(ctx, "Successfully collected prefixes from Netbox: ", nbi.PrefixesIndexByPrefix)
	return nil
}
