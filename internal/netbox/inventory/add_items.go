package inventory

import (
	"context"
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// AddTag adds the newTag from source sourceName to the local inventory.
func (nbi *NetboxInventory) AddTag(ctx context.Context, newTag *objects.Tag) (*objects.Tag, error) {
	nbi.TagsLock.Lock()
	defer nbi.TagsLock.Unlock()
	if _, ok := nbi.TagsIndexByName[newTag.Name]; ok {
		oldTag := nbi.TagsIndexByName[newTag.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newTag, oldTag, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Tag ", newTag.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedTag, err := service.Patch[objects.Tag](ctx, nbi.NetboxAPI, oldTag.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.TagsIndexByName[newTag.Name] = patchedTag
		} else {
			nbi.Logger.Debug(ctx, "Tag ", newTag.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Tag ", newTag.Name, " does not exist in Netbox. Creating it...")
		createdTag, err := service.Create[objects.Tag](ctx, nbi.NetboxAPI, newTag)
		if err != nil {
			return nil, err
		}
		nbi.TagsIndexByName[newTag.Name] = createdTag
	}
	return nbi.TagsIndexByName[newTag.Name], nil
}

// AddTenants adds a new tenant to the local netbox inventory.
func (nbi *NetboxInventory) AddTenant(ctx context.Context, newTenant *objects.Tenant) (*objects.Tenant, error) {
	newTenant.Tags = append(newTenant.Tags, nbi.SsotTag)
	nbi.TenantsLock.Lock()
	defer nbi.TenantsLock.Unlock()
	if _, ok := nbi.TenantsIndexByName[newTenant.Name]; ok {
		oldTenant := nbi.TenantsIndexByName[newTenant.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newTenant, oldTenant, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Tenant ", newTenant.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedTenant, err := service.Patch[objects.Tenant](ctx, nbi.NetboxAPI, oldTenant.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.TenantsIndexByName[newTenant.Name] = patchedTenant
		} else {
			nbi.Logger.Debug(ctx, "Tenant ", newTenant.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Tenant ", newTenant.Name, " does not exist in Netbox. Creating it...")
		createdTag, err := service.Create[objects.Tenant](ctx, nbi.NetboxAPI, newTenant)
		if err != nil {
			return nil, err
		}
		nbi.TenantsIndexByName[newTenant.Name] = createdTag
	}
	return nbi.TenantsIndexByName[newTenant.Name], nil
}

// AddContact adds a contact to the local netbox inventory.
func (nbi *NetboxInventory) AddSite(ctx context.Context, newSite *objects.Site) (*objects.Site, error) {
	newSite.Tags = append(newSite.Tags, nbi.SsotTag)

	nbi.SitesLock.Lock()
	defer nbi.SitesLock.Unlock()
	if _, ok := nbi.SitesIndexByName[newSite.Name]; ok {
		oldSite := nbi.SitesIndexByName[newSite.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newSite, oldSite, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Site ", newSite.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedSite, err := service.Patch[objects.Site](ctx, nbi.NetboxAPI, oldSite.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.SitesIndexByName[newSite.Name] = patchedSite
		} else {
			nbi.Logger.Debug(ctx, "Site ", newSite.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Site ", newSite.Name, " does not exist in Netbox. Creating it...")
		createdContact, err := service.Create[objects.Site](ctx, nbi.NetboxAPI, newSite)
		if err != nil {
			return nil, err
		}
		nbi.SitesIndexByName[newSite.Name] = createdContact
	}
	return nbi.SitesIndexByName[newSite.Name], nil
}

// AddContactRole adds the newContactRole to the local netbox inventory.
func (nbi *NetboxInventory) AddContactRole(ctx context.Context, newContactRole *objects.ContactRole) (*objects.ContactRole, error) {
	newContactRole.NetboxObject.Tags = []*objects.Tag{nbi.SsotTag}

	nbi.ContactRolesLock.Lock()
	defer nbi.ContactRolesLock.Unlock()
	addSourceNameCustomField(ctx, &newContactRole.NetboxObject)
	if _, ok := nbi.ContactRolesIndexByName[newContactRole.Name]; ok {
		oldContactRole := nbi.ContactRolesIndexByName[newContactRole.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newContactRole, oldContactRole, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Contact role ", newContactRole.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContactRole, err := service.Patch[objects.ContactRole](ctx, nbi.NetboxAPI, oldContactRole.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ContactRolesIndexByName[newContactRole.Name] = patchedContactRole
		} else {
			nbi.Logger.Debug(ctx, "Contact role ", newContactRole.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Contact role ", newContactRole.Name, " does not exist in Netbox. Creating it...")
		newContactRole, err := service.Create[objects.ContactRole](ctx, nbi.NetboxAPI, newContactRole)
		if err != nil {
			return nil, err
		}
		nbi.ContactRolesIndexByName[newContactRole.Name] = newContactRole
	}
	return nbi.ContactRolesIndexByName[newContactRole.Name], nil
}

// AddContactGroup adds contact group to the local netbox inventory.
func (nbi *NetboxInventory) AddContactGroup(ctx context.Context, newContactGroup *objects.ContactGroup) (*objects.ContactGroup, error) {
	nbi.ContactGroupsLock.Lock()
	defer nbi.ContactGroupsLock.Unlock()
	if _, ok := nbi.ContactGroupsIndexByName[newContactGroup.Name]; ok {
		oldContactGroup := nbi.ContactGroupsIndexByName[newContactGroup.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newContactGroup, oldContactGroup, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Contact group ", newContactGroup.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContactGroup, err := service.Patch[objects.ContactGroup](ctx, nbi.NetboxAPI, oldContactGroup.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ContactGroupsIndexByName[newContactGroup.Name] = patchedContactGroup
		} else {
			nbi.Logger.Debug(ctx, "Contact group ", newContactGroup.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Contact group ", newContactGroup.Name, " does not exist in Netbox. Creating it...")
		newContactGroup, err := service.Create[objects.ContactGroup](ctx, nbi.NetboxAPI, newContactGroup)
		if err != nil {
			return nil, err
		}
		nbi.ContactGroupsIndexByName[newContactGroup.Name] = newContactGroup
	}
	return nbi.ContactGroupsIndexByName[newContactGroup.Name], nil
}

// AddContact adds a contact to the local netbox inventory.
func (nbi *NetboxInventory) AddContact(ctx context.Context, newContact *objects.Contact) (*objects.Contact, error) {
	newContact.Tags = append(newContact.Tags, nbi.SsotTag)

	nbi.ContactsLock.Lock()
	defer nbi.ContactsLock.Unlock()
	if _, ok := nbi.ContactsIndexByName[newContact.Name]; ok {
		oldContact := nbi.ContactsIndexByName[newContact.Name]
		delete(nbi.OrphanManager[constants.ContactsAPIPath], oldContact.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newContact, oldContact, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Contact ", newContact.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContact, err := service.Patch[objects.Contact](ctx, nbi.NetboxAPI, oldContact.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ContactsIndexByName[newContact.Name] = patchedContact
		} else {
			nbi.Logger.Debug(ctx, "Contact ", newContact.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Contact ", newContact.Name, " does not exist in Netbox. Creating it...")
		createdContact, err := service.Create[objects.Contact](ctx, nbi.NetboxAPI, newContact)
		if err != nil {
			return nil, err
		}
		nbi.ContactsIndexByName[newContact.Name] = createdContact
	}
	return nbi.ContactsIndexByName[newContact.Name], nil
}

// AddContact assignment adds a contact assignment to the local netbox inventory.
// TODO: Make index check less code and more universal, checking each level is ugly.
func (nbi *NetboxInventory) AddContactAssignment(ctx context.Context, newCA *objects.ContactAssignment) (*objects.ContactAssignment, error) {
	nbi.ContactAssignmentsLock.Lock()
	defer nbi.ContactAssignmentsLock.Unlock()
	if nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType] == nil {
		nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType] = make(map[int]map[int]map[int]*objects.ContactAssignment)
	}
	if nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID] == nil {
		nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID] = make(map[int]map[int]*objects.ContactAssignment)
	}
	if nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID] == nil {
		nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID] = make(map[int]*objects.ContactAssignment)
	}
	newCA.Tags = append(newCA.Tags, nbi.SsotTag)
	if _, ok := nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID]; ok {
		oldCA := nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID]
		delete(nbi.OrphanManager[constants.ContactAssignmentsAPIPath], oldCA.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newCA, oldCA, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "ContactAssignment ", newCA.ID, " already exists in Netbox but is out of date. Patching it... ")
			patchedCA, err := service.Patch[objects.ContactAssignment](ctx, nbi.NetboxAPI, oldCA.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID] = patchedCA
		} else {
			nbi.Logger.Debug(ctx, "ContactAssignment ", newCA.ID, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debugf(ctx, "ContactAssignment %s does not exist in Netbox. Creating it...", newCA)
		newCA, err := service.Create[objects.ContactAssignment](ctx, nbi.NetboxAPI, newCA)
		if err != nil {
			return nil, err
		}
		nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID] = newCA
	}
	return nbi.ContactAssignmentsIndexByObjectTypeAndObjectIDAndContactIDAndRoleID[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID], nil
}

// AddCustomField adds a custom field to the Netbox inventory.
// It takes a context and a newCf object as input and returns the created or patched custom field along with any error encountered.
// If the custom field already exists in Netbox but is out of date, it will be patched with the new values.
// If the custom field does not exist, it will be created.
func (nbi *NetboxInventory) AddCustomField(ctx context.Context, newCf *objects.CustomField) (*objects.CustomField, error) {
	nbi.CustomFieldsLock.Lock()
	defer nbi.CustomFieldsLock.Unlock()
	if _, ok := nbi.CustomFieldsIndexByName[newCf.Name]; ok {
		oldCustomField := nbi.CustomFieldsIndexByName[newCf.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newCf, oldCustomField, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Custom field ", newCf.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedCf, err := service.Patch[objects.CustomField](ctx, nbi.NetboxAPI, oldCustomField.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.CustomFieldsIndexByName[newCf.Name] = patchedCf
		} else {
			nbi.Logger.Debug(ctx, "Custom field ", newCf.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Custom field ", newCf.Name, " does not exist in Netbox. Creating it...")
		createdCf, err := service.Create[objects.CustomField](ctx, nbi.NetboxAPI, newCf)
		if err != nil {
			return nil, err
		}
		nbi.CustomFieldsIndexByName[createdCf.Name] = createdCf
	}
	return nbi.CustomFieldsIndexByName[newCf.Name], nil
}

// AddClusterGroup adds a new cluster group to the Netbox inventory.
// It takes a context and a newCg object as input and returns the newly created cluster group and an error (if any).
// If the cluster group already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the cluster group does not exist, it creates a new one.
// The function also updates the cluster group index by name and removes the ID from the orphan manager.
func (nbi *NetboxInventory) AddClusterGroup(ctx context.Context, newCg *objects.ClusterGroup) (*objects.ClusterGroup, error) {
	nbi.ClusterGroupsLock.Lock()
	defer nbi.ClusterGroupsLock.Unlock()
	newCg.Tags = append(newCg.Tags, nbi.SsotTag)
	addSourceNameCustomField(ctx, &newCg.NetboxObject)
	if _, ok := nbi.ClusterGroupsIndexByName[newCg.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldCg := nbi.ClusterGroupsIndexByName[newCg.Name]
		delete(nbi.OrphanManager[constants.ClusterGroupsAPIPath], oldCg.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newCg, oldCg, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Cluster group ", newCg.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedCg, err := service.Patch[objects.ClusterGroup](ctx, nbi.NetboxAPI, oldCg.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ClusterGroupsIndexByName[newCg.Name] = patchedCg
		} else {
			nbi.Logger.Debug(ctx, "Cluster group ", newCg.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Cluster group ", newCg.Name, " does not exist in Netbox. Creating it...")
		newCg, err := service.Create[objects.ClusterGroup](ctx, nbi.NetboxAPI, newCg)
		if err != nil {
			return nil, err
		}
		nbi.ClusterGroupsIndexByName[newCg.Name] = newCg
	}
	// Delete id from orphan manager
	return nbi.ClusterGroupsIndexByName[newCg.Name], nil
}

// AddClusterType adds a new cluster type to the Netbox inventory.
// It takes a context and a newClusterType object as input and returns the created or updated cluster type object and an error, if any.
// If the cluster type already exists in Netbox, it checks if it is up to date. If not, it patches the existing cluster type.
// If the cluster type does not exist, it creates a new one.
func (nbi *NetboxInventory) AddClusterType(ctx context.Context, newClusterType *objects.ClusterType) (*objects.ClusterType, error) {
	nbi.ClusterTypesLock.Lock()
	defer nbi.ClusterTypesLock.Unlock()
	newClusterType.Tags = append(newClusterType.Tags, nbi.SsotTag)
	if _, ok := nbi.ClusterTypesIndexByName[newClusterType.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldClusterType := nbi.ClusterTypesIndexByName[newClusterType.Name]
		delete(nbi.OrphanManager[constants.ClusterTypesAPIPath], oldClusterType.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newClusterType, oldClusterType, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Cluster type ", newClusterType.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedClusterType, err := service.Patch[objects.ClusterType](ctx, nbi.NetboxAPI, oldClusterType.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ClusterTypesIndexByName[newClusterType.Name] = patchedClusterType
			return patchedClusterType, nil
		}
		nbi.Logger.Debug(ctx, "Cluster type ", newClusterType.Name, " already exists in Netbox and is up to date...")
		existingClusterType := nbi.ClusterTypesIndexByName[newClusterType.Name]
		return existingClusterType, nil
	}
	nbi.Logger.Debug(ctx, "Cluster type ", newClusterType.Name, " does not exist in Netbox. Creating it...")
	newClusterType, err := service.Create[objects.ClusterType](ctx, nbi.NetboxAPI, newClusterType)
	if err != nil {
		return nil, err
	}
	nbi.ClusterTypesIndexByName[newClusterType.Name] = newClusterType
	return newClusterType, nil
}

// AddCluster adds a new cluster to the Netbox inventory.
// It takes a context and a pointer to a Cluster object as input.
// It returns the newly created cluster object and an error, if any.
// If the cluster already exists in Netbox, it checks if the existing cluster is up to date.
// If it is not up to date, it patches the existing cluster with the changes from the new cluster.
// If the cluster does not exist in Netbox, it creates a new cluster.
func (nbi *NetboxInventory) AddCluster(ctx context.Context, newCluster *objects.Cluster) (*objects.Cluster, error) {
	nbi.ClustersLock.Lock()
	defer nbi.ClustersLock.Unlock()
	newCluster.Tags = append(newCluster.Tags, nbi.SsotTag)
	addSourceNameCustomField(ctx, &newCluster.NetboxObject)
	if _, ok := nbi.ClustersIndexByName[newCluster.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldCluster := nbi.ClustersIndexByName[newCluster.Name]
		delete(nbi.OrphanManager[constants.ClustersAPIPath], oldCluster.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newCluster, oldCluster, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Cluster ", newCluster.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedCluster, err := service.Patch[objects.Cluster](ctx, nbi.NetboxAPI, oldCluster.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ClustersIndexByName[newCluster.Name] = patchedCluster
		} else {
			nbi.Logger.Debug(ctx, "Cluster ", newCluster.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Cluster ", newCluster.Name, " does not exist in Netbox. Creating it...")
		createdCluster, err := service.Create[objects.Cluster](ctx, nbi.NetboxAPI, newCluster)
		if err != nil {
			return nil, err
		}
		nbi.ClustersIndexByName[createdCluster.Name] = createdCluster
	}
	return nbi.ClustersIndexByName[newCluster.Name], nil
}

// AddDeviceRole adds a new device role to the Netbox inventory.
// It takes a context and a newDeviceRole object as input and returns the created device role object and an error, if any.
// If the device role already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the device role does not exist, it creates a new one.
func (nbi *NetboxInventory) AddDeviceRole(ctx context.Context, newDeviceRole *objects.DeviceRole) (*objects.DeviceRole, error) {
	nbi.DeviceRolesLock.Lock()
	defer nbi.DeviceRolesLock.Unlock()
	newDeviceRole.Tags = append(newDeviceRole.Tags, nbi.SsotTag)
	if _, ok := nbi.DeviceRolesIndexByName[newDeviceRole.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldDeviceRole := nbi.DeviceRolesIndexByName[newDeviceRole.Name]
		delete(nbi.OrphanManager[constants.DeviceRolesAPIPath], nbi.DeviceRolesIndexByName[newDeviceRole.Name].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newDeviceRole, oldDeviceRole, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Device role ", newDeviceRole.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedDeviceRole, err := service.Patch[objects.DeviceRole](ctx, nbi.NetboxAPI, oldDeviceRole.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.DeviceRolesIndexByName[newDeviceRole.Name] = patchedDeviceRole
		} else {
			nbi.Logger.Debug(ctx, "Device role ", newDeviceRole.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Device role ", newDeviceRole.Name, " does not exist in Netbox. Creating it...")
		newDeviceRole, err := service.Create[objects.DeviceRole](ctx, nbi.NetboxAPI, newDeviceRole)
		if err != nil {
			return nil, err
		}
		nbi.DeviceRolesIndexByName[newDeviceRole.Name] = newDeviceRole
	}
	return nbi.DeviceRolesIndexByName[newDeviceRole.Name], nil
}

// AddManufacturer adds a new manufacturer to the Netbox inventory.
// It takes a context, `ctx`, and a pointer to a `newManufacturer` object as input.
// The function returns a pointer to the newly created manufacturer and an error, if any.
// If the manufacturer already exists in Netbox, the function checks if it is up to date.
// If it is not up to date, the function patches the existing manufacturer with the updated information.
// If the manufacturer does not exist in Netbox, the function creates a new one.
func (nbi *NetboxInventory) AddManufacturer(ctx context.Context, newManufacturer *objects.Manufacturer) (*objects.Manufacturer, error) {
	nbi.ManufacturersLock.Lock()
	defer nbi.ManufacturersLock.Unlock()
	newManufacturer.Tags = append(newManufacturer.Tags, nbi.SsotTag)
	if _, ok := nbi.ManufacturersIndexByName[newManufacturer.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldManufacturer := nbi.ManufacturersIndexByName[newManufacturer.Name]
		delete(nbi.OrphanManager[constants.ManufacturersAPIPath], oldManufacturer.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newManufacturer, oldManufacturer, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Manufacturer ", newManufacturer.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedManufacturer, err := service.Patch[objects.Manufacturer](ctx, nbi.NetboxAPI, oldManufacturer.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ManufacturersIndexByName[newManufacturer.Name] = patchedManufacturer
		} else {
			nbi.Logger.Debug(ctx, "Manufacturer ", newManufacturer.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Manufacturer ", newManufacturer.Name, " does not exist in Netbox. Creating it...")
		newManufacturer, err := service.Create[objects.Manufacturer](ctx, nbi.NetboxAPI, newManufacturer)
		if err != nil {
			return nil, err
		}
		nbi.ManufacturersIndexByName[newManufacturer.Name] = newManufacturer
	}
	return nbi.ManufacturersIndexByName[newManufacturer.Name], nil
}

// AddDeviceType adds a new device type to the Netbox inventory.
// It takes a context and a newDeviceType object as input and returns the created or updated device type object and an error, if any.
// If the device type already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the device type does not exist, it creates a new one.
func (nbi *NetboxInventory) AddDeviceType(ctx context.Context, newDeviceType *objects.DeviceType) (*objects.DeviceType, error) {
	nbi.DeviceTypesLock.Lock()
	defer nbi.DeviceTypesLock.Unlock()
	newDeviceType.Tags = append(newDeviceType.Tags, nbi.SsotTag)
	if _, ok := nbi.DeviceTypesIndexByModel[newDeviceType.Model]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldDeviceType := nbi.DeviceTypesIndexByModel[newDeviceType.Model]
		delete(nbi.OrphanManager[constants.DeviceTypesAPIPath], oldDeviceType.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newDeviceType, oldDeviceType, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Device type ", newDeviceType.Model, " already exists in Netbox but is out of date. Patching it...")
			patchedDeviceType, err := service.Patch[objects.DeviceType](ctx, nbi.NetboxAPI, oldDeviceType.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.DeviceTypesIndexByModel[newDeviceType.Model] = patchedDeviceType
		} else {
			nbi.Logger.Debug(ctx, "Device type ", newDeviceType.Model, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Device type ", newDeviceType.Model, " does not exist in Netbox. Creating it...")
		newDeviceType, err := service.Create[objects.DeviceType](ctx, nbi.NetboxAPI, newDeviceType)
		if err != nil {
			return nil, err
		}
		nbi.DeviceTypesIndexByModel[newDeviceType.Model] = newDeviceType
	}
	return nbi.DeviceTypesIndexByModel[newDeviceType.Model], nil
}

// AddPlatform adds a new platform to the Netbox inventory.
// It takes a context and a newPlatform object as input and returns the created or updated platform object and an error, if any.
// If the platform already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the platform does not exist, it creates a new one.
func (nbi *NetboxInventory) AddPlatform(ctx context.Context, newPlatform *objects.Platform) (*objects.Platform, error) {
	nbi.PlatformsLock.Lock()
	newPlatform.Tags = append(newPlatform.Tags, nbi.SsotTag)
	defer nbi.PlatformsLock.Unlock()
	if _, ok := nbi.PlatformsIndexByName[newPlatform.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldPlatform := nbi.PlatformsIndexByName[newPlatform.Name]
		delete(nbi.OrphanManager[constants.PlatformsAPIPath], oldPlatform.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newPlatform, oldPlatform, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Platform ", newPlatform.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedPlatform, err := service.Patch[objects.Platform](ctx, nbi.NetboxAPI, oldPlatform.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.PlatformsIndexByName[newPlatform.Name] = patchedPlatform
		} else {
			nbi.Logger.Debug(ctx, "Platform ", newPlatform.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Platform ", newPlatform.Name, " does not exist in Netbox. Creating it...")
		newPlatform, err := service.Create[objects.Platform](ctx, nbi.NetboxAPI, newPlatform)
		if err != nil {
			return nil, err
		}
		nbi.PlatformsIndexByName[newPlatform.Name] = newPlatform
	}
	return nbi.PlatformsIndexByName[newPlatform.Name], nil
}

// AddRackRole adds a new rack role to the Netbox inventory.
// It takes a context and a newRackRole object as input and returns the created or updated rack role object and an error, if any.
// If the rack role already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the rack role does not exist, it creates a new one.
func (nbi *NetboxInventory) AddDevice(ctx context.Context, newDevice *objects.Device) (*objects.Device, error) {
	nbi.DevicesLock.Lock()
	defer nbi.DevicesLock.Unlock()
	newDevice.Tags = append(newDevice.Tags, nbi.SsotTag)
	addSourceNameCustomField(ctx, &newDevice.NetboxObject)
	if newDevice.Site == nil {
		return nil, fmt.Errorf("device %s is not assigned to a site, but it should be", newDevice)
	}
	if _, ok := nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID]; ok {
		oldDevice := nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID]
		delete(nbi.OrphanManager[constants.DevicesAPIPath], oldDevice.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newDevice, oldDevice, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Device ", newDevice.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedDevice, err := service.Patch[objects.Device](ctx, nbi.NetboxAPI, oldDevice.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID] = patchedDevice
		} else {
			nbi.Logger.Debug(ctx, "Device ", newDevice.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Device ", newDevice.Name, " does not exist in Netbox. Creating it...")
		newDevice, err := service.Create[objects.Device](ctx, nbi.NetboxAPI, newDevice)
		if err != nil {
			return nil, err
		}
		if nbi.DevicesIndexByNameAndSiteID[newDevice.Name] == nil {
			nbi.DevicesIndexByNameAndSiteID[newDevice.Name] = make(map[int]*objects.Device)
		}
		nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID] = newDevice
	}
	return nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID], nil
}

// AddVirtualDeviceContext adds new virtual device context to the local inventory.
// It takes a context and a newVDC object as input and returns the created or updated virtual device context object and an error, if any.
// If the virtual device context already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the virtual device context does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVirtualDeviceContext(ctx context.Context, newVDC *objects.VirtualDeviceContext) (*objects.VirtualDeviceContext, error) {
	nbi.DevicesLock.Lock()
	defer nbi.DevicesLock.Unlock()
	newVDC.Tags = append(newVDC.Tags, nbi.SsotTag)
	if newVDC.Device == nil {
		return nil, fmt.Errorf("VirtualDeviceContext %s is not assigned to a device, but it should be", newVDC)
	}
	addSourceNameCustomField(ctx, &newVDC.NetboxObject)
	if _, ok := nbi.VirtualDeviceContextsIndexByNameAndDeviceID[newVDC.Name][newVDC.Device.ID]; ok {
		oldVDC := nbi.VirtualDeviceContextsIndexByNameAndDeviceID[newVDC.Name][newVDC.Device.ID]
		delete(nbi.OrphanManager[constants.VirtualDeviceContextsAPIPath], oldVDC.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVDC, oldVDC, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "VirtualDeviceContext ", newVDC.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVDC, err := service.Patch[objects.VirtualDeviceContext](ctx, nbi.NetboxAPI, oldVDC.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VirtualDeviceContextsIndexByNameAndDeviceID[newVDC.Name][newVDC.Device.ID] = patchedVDC
		} else {
			nbi.Logger.Debug(ctx, "VirtualDeviceContext ", newVDC.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "VirtualDeviceContext ", newVDC.Name, " does not exist in Netbox. Creating it...")
		newDevice, err := service.Create[objects.VirtualDeviceContext](ctx, nbi.NetboxAPI, newVDC)
		if err != nil {
			return nil, err
		}
		if nbi.VirtualDeviceContextsIndexByNameAndDeviceID[newDevice.Name] == nil {
			nbi.VirtualDeviceContextsIndexByNameAndDeviceID[newDevice.Name] = make(map[int]*objects.VirtualDeviceContext)
		}
		nbi.VirtualDeviceContextsIndexByNameAndDeviceID[newDevice.Name][newDevice.Device.ID] = newDevice
	}
	return nbi.VirtualDeviceContextsIndexByNameAndDeviceID[newVDC.Name][newVDC.Device.ID], nil
}

// AddVlanGroup adds a new vlan group to the Netbox inventory.
// It takes a context and a newVlanGroup object as input and returns the created or updated vlan group object and an error, if any.
// If the vlan group already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the vlan group does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVlanGroup(ctx context.Context, newVlanGroup *objects.VlanGroup) (*objects.VlanGroup, error) {
	nbi.VlanGroupsLock.Lock()
	defer nbi.VlanGroupsLock.Unlock()
	newVlanGroup.Tags = append(newVlanGroup.Tags, nbi.SsotTag)
	if _, ok := nbi.VlanGroupsIndexByName[newVlanGroup.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldVlanGroup := nbi.VlanGroupsIndexByName[newVlanGroup.Name]
		delete(nbi.OrphanManager[constants.VlanGroupsAPIPath], oldVlanGroup.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVlanGroup, oldVlanGroup, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "VlanGroup ", newVlanGroup.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVlanGroup, err := service.Patch[objects.VlanGroup](ctx, nbi.NetboxAPI, oldVlanGroup.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VlanGroupsIndexByName[newVlanGroup.Name] = patchedVlanGroup
		} else {
			nbi.Logger.Debug(ctx, "Vlan ", newVlanGroup.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Vlan ", newVlanGroup.Name, " does not exist in Netbox. Creating it...")
		newVlan, err := service.Create[objects.VlanGroup](ctx, nbi.NetboxAPI, newVlanGroup)
		if err != nil {
			return nil, err
		}
		nbi.VlanGroupsIndexByName[newVlan.Name] = newVlan
	}
	return nbi.VlanGroupsIndexByName[newVlanGroup.Name], nil
}

// AddVlan adds a new vlan to the Netbox inventory.
// It takes a context and a newVlan object as input and returns the created or updated vlan object and an error, if any.
// If the vlan already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the vlan does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVlan(ctx context.Context, newVlan *objects.Vlan) (*objects.Vlan, error) {
	nbi.VlansLock.Lock()
	defer nbi.VlansLock.Unlock()
	newVlan.Tags = append(newVlan.Tags, nbi.SsotTag)
	addSourceNameCustomField(ctx, &newVlan.NetboxObject)
	if _, ok := nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldVlan := nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid]
		delete(nbi.OrphanManager[constants.VlansAPIPath], oldVlan.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVlan, oldVlan, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Vlan ", newVlan.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVlan, err := service.Patch[objects.Vlan](ctx, nbi.NetboxAPI, oldVlan.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid] = patchedVlan
		} else {
			nbi.Logger.Debug(ctx, "Vlan ", newVlan.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Vlan ", newVlan.Name, " does not exist in Netbox. Creating it...")
		newVlan, err := service.Create[objects.Vlan](ctx, nbi.NetboxAPI, newVlan)
		if err != nil {
			return nil, err
		}
		if nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID] == nil {
			nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID] = make(map[int]*objects.Vlan)
		}
		nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid] = newVlan
	}
	return nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid], nil
}

// AddInterface adds a new interface to the Netbox inventory.
// It takes a context and a newInterface object as input and returns the created or updated interface object and an error, if any.
// If the interface already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the interface does not exist, it creates a new one.
func (nbi *NetboxInventory) AddInterface(ctx context.Context, newInterface *objects.Interface) (*objects.Interface, error) {
	nbi.InterfacesLock.Lock()
	defer nbi.InterfacesLock.Unlock()
	newInterface.Tags = append(newInterface.Tags, nbi.SsotTag)
	addSourceNameCustomField(ctx, &newInterface.NetboxObject)
	if len(newInterface.Name) > constants.MaxInterfaceNameLength {
		newInterface.Name = newInterface.Name[:constants.MaxInterfaceNameLength]
	}
	if _, ok := nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[constants.InterfacesAPIPath], nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newInterface, nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name], false, nbi.SourcePriority)
		oldIntf := nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Interface ", newInterface.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedInterface, err := service.Patch[objects.Interface](ctx, nbi.NetboxAPI, oldIntf.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name] = patchedInterface
		} else {
			nbi.Logger.Debug(ctx, "Interface ", newInterface.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "Interface ", newInterface.Name, " does not exist in Netbox. Creating it...")
		newInterface, err := service.Create[objects.Interface](ctx, nbi.NetboxAPI, newInterface)
		if err != nil {
			return nil, err
		}
		if nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID] == nil {
			nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID] = make(map[string]*objects.Interface)
		}
		nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name] = newInterface
	}
	return nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name], nil
}

// AddVM adds a new virtual machine to the Netbox inventory.
// It takes a context and a newVM object as input and returns the created or updated virtual machine object and an error, if any.
// If the virtual machine already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the virtual machine does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVM(ctx context.Context, newVM *objects.VM) (*objects.VM, error) {
	nbi.VMsLock.Lock()
	defer nbi.VMsLock.Unlock()
	newVM.Tags = append(newVM.Tags, nbi.SsotTag)
	addSourceNameCustomField(ctx, &newVM.NetboxObject)
	newVMClusterID := -1
	if newVM.Cluster != nil {
		newVMClusterID = newVM.Cluster.ID
	}
	if len(newVM.Name) > constants.MaxVMNameLength {
		newVM.Name = newVM.Name[:constants.MaxVMNameLength]
	}
	if oldVM, ok := nbi.VMsIndexByNameAndClusterID[newVM.Name][newVMClusterID]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[constants.VirtualMachinesAPIPath], oldVM.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVM, oldVM, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(ctx, "%s already exists in Netbox but is out of date. Patching it...", newVM)
			patchedVM, err := service.Patch[objects.VM](ctx, nbi.NetboxAPI, oldVM.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VMsIndexByNameAndClusterID[newVM.Name][newVMClusterID] = patchedVM
		} else {
			nbi.Logger.Debugf(ctx, "%s already exists in Netbox and is up to date...", newVM)
		}
	} else {
		nbi.Logger.Debugf(ctx, "%s does not exist in Netbox. Creating it...", newVM)
		newVM, err := service.Create[objects.VM](ctx, nbi.NetboxAPI, newVM)
		if err != nil {
			return nil, err
		}
		if nbi.VMsIndexByNameAndClusterID[newVM.Name] == nil {
			nbi.VMsIndexByNameAndClusterID[newVM.Name] = make(map[int]*objects.VM)
		}
		nbi.VMsIndexByNameAndClusterID[newVM.Name][newVMClusterID] = newVM
		return newVM, nil
	}
	return nbi.VMsIndexByNameAndClusterID[newVM.Name][newVMClusterID], nil
}

// AddVMInterface adds a new virtual machine interface to the Netbox inventory.
// It takes a context and a newVMInterface object as input and returns the created or updated virtual machine interface object and an error, if any.
// If the virtual machine interface already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the virtual machine interface does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVMInterface(ctx context.Context, newVMInterface *objects.VMInterface) (*objects.VMInterface, error) {
	newVMInterface.Tags = append(newVMInterface.Tags, nbi.SsotTag)
	nbi.VMInterfacesLock.Lock()
	defer nbi.VMInterfacesLock.Unlock()
	if len(newVMInterface.Name) > constants.MaxVMInterfaceNameLength {
		newVMInterface.Name = newVMInterface.Name[:constants.MaxVMInterfaceNameLength]
	}
	if _, ok := nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[constants.VMInterfacesAPIPath], nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVMInterface, nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name], false, nbi.SourcePriority)
		oldVMIface := nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "VM interface ", newVMInterface.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVMInterface, err := service.Patch[objects.VMInterface](ctx, nbi.NetboxAPI, oldVMIface.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name] = patchedVMInterface
		} else {
			nbi.Logger.Debug(ctx, "VM interface ", newVMInterface.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "VM interface ", newVMInterface.Name, " does not exist in Netbox. Creating it...")
		newVMInterface, err := service.Create[objects.VMInterface](ctx, nbi.NetboxAPI, newVMInterface)
		if err != nil {
			return nil, err
		}
		if nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID] == nil {
			nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID] = make(map[string]*objects.VMInterface)
		}
		nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name] = newVMInterface
	}
	return nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name], nil
}

// AddIPAddress adds a new IP address to the Netbox inventory.
// It takes a context and a newIPAddress object as input and returns the created or updated IP address object and an error, if any.
// If the IP address already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the IP address does not exist, it creates a new one.
func (nbi *NetboxInventory) AddIPAddress(ctx context.Context, newIPAddress *objects.IPAddress) (*objects.IPAddress, error) {
	newIPAddress.Tags = append(newIPAddress.Tags, nbi.SsotTag)
	nbi.IPAddressesLock.Lock()
	defer nbi.IPAddressesLock.Unlock()
	addSourceNameCustomField(ctx, &newIPAddress.NetboxObject)
	if _, ok := nbi.IPAdressesIndexByAddress[newIPAddress.Address]; ok {
		// Delete id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[constants.IPAddressesAPIPath], nbi.IPAdressesIndexByAddress[newIPAddress.Address].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newIPAddress, nbi.IPAdressesIndexByAddress[newIPAddress.Address], false, nbi.SourcePriority)
		oldIPAddress := nbi.IPAdressesIndexByAddress[newIPAddress.Address]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "IP address ", newIPAddress.Address, " already exists in Netbox but is out of date. Patching it...")
			patchedIPAddress, err := service.Patch[objects.IPAddress](ctx, nbi.NetboxAPI, oldIPAddress.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.IPAdressesIndexByAddress[newIPAddress.Address] = patchedIPAddress
			return patchedIPAddress, nil
		}
		nbi.Logger.Debug(ctx, "IP address ", newIPAddress.Address, " already exists in Netbox and is up to date...")
	} else {
		nbi.Logger.Debug(ctx, "IP address ", newIPAddress.Address, " does not exist in Netbox. Creating it...")
		newIPAddress, err := service.Create[objects.IPAddress](ctx, nbi.NetboxAPI, newIPAddress)
		if err != nil {
			return nil, err
		}
		nbi.IPAdressesIndexByAddress[newIPAddress.Address] = newIPAddress
		return newIPAddress, nil
	}
	return nbi.IPAdressesIndexByAddress[newIPAddress.Address], nil
}

// AddPrefix adds a new prefix to the Netbox inventory.
// It takes a context and a newPrefix object as input and returns the created or updated prefix object and an error, if any.
// If the prefix already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the prefix does not exist, it creates a new one.
func (nbi *NetboxInventory) AddPrefix(ctx context.Context, newPrefix *objects.Prefix) (*objects.Prefix, error) {
	newPrefix.Tags = append(newPrefix.Tags, nbi.SsotTag)
	nbi.PrefixesLock.Lock()
	if newPrefix.NetboxObject.CustomFields == nil {
		newPrefix.NetboxObject.CustomFields = make(map[string]interface{})
	}
	newPrefix.NetboxObject.CustomFields[constants.CustomFieldSourceName] = ctx.Value(constants.CtxSourceKey).(string) //nolint:forcetypeassert
	defer nbi.PrefixesLock.Unlock()
	if _, ok := nbi.PrefixesIndexByPrefix[newPrefix.Prefix]; ok {
		// Delete id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[constants.PrefixesAPIPath], nbi.PrefixesIndexByPrefix[newPrefix.Prefix].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newPrefix, nbi.PrefixesIndexByPrefix[newPrefix.Prefix], false, nbi.SourcePriority)
		oldPrefix := nbi.PrefixesIndexByPrefix[newPrefix.Prefix]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "Prefix ", newPrefix.Prefix, " already exists in Netbox but is out of date. Patching it...")
			patchedPrefix, err := service.Patch[objects.Prefix](ctx, nbi.NetboxAPI, oldPrefix.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.PrefixesIndexByPrefix[newPrefix.Prefix] = patchedPrefix
		} else {
			nbi.Logger.Debug(ctx, "IP address ", newPrefix.Prefix, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "IP address ", newPrefix.Prefix, " does not exist in Netbox. Creating it...")
		newPrefix, err := service.Create[objects.Prefix](ctx, nbi.NetboxAPI, newPrefix)
		if err != nil {
			return nil, err
		}
		nbi.PrefixesIndexByPrefix[newPrefix.Prefix] = newPrefix
		return newPrefix, nil
	}
	return nbi.PrefixesIndexByPrefix[newPrefix.Prefix], nil
}

// AddWirelessLAN adds a new wireless LAN to the Netbox inventory.
// It takes a context and a newWirelessLan object as input and returns the created or updated wireless LAN object and an error, if any.
// If the wireless LAN already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the wireless LAN does not exist, it creates a new one.
func (nbi *NetboxInventory) AddWirelessLAN(ctx context.Context, newWirelessLan *objects.WirelessLAN) (*objects.WirelessLAN, error) {
	newWirelessLan.Tags = append(newWirelessLan.Tags, nbi.SsotTag)
	if _, ok := nbi.WirelessLANsIndexBySSID[newWirelessLan.SSID]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldWirelessLan := nbi.WirelessLANsIndexBySSID[newWirelessLan.SSID]
		delete(nbi.OrphanManager[constants.WirelessLANsAPIPath], oldWirelessLan.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newWirelessLan, oldWirelessLan, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "WirelessLAN ", newWirelessLan.SSID, " already exists in Netbox but is out of date. Patching it...")
			patchedWirelessLan, err := service.Patch[objects.WirelessLAN](ctx, nbi.NetboxAPI, oldWirelessLan.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.WirelessLANsIndexBySSID[newWirelessLan.SSID] = patchedWirelessLan
		} else {
			nbi.Logger.Debug(ctx, "WirelessLAN ", newWirelessLan.SSID, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "WirelessLAN ", newWirelessLan.SSID, " does not exist in Netbox. Creating it...")
		newWirelessLan, err := service.Create[objects.WirelessLAN](ctx, nbi.NetboxAPI, newWirelessLan)
		if err != nil {
			return nil, err
		}
		nbi.WirelessLANsIndexBySSID[newWirelessLan.SSID] = newWirelessLan
	}
	return nbi.WirelessLANsIndexBySSID[newWirelessLan.SSID], nil
}

// AddWirelessLANGroup adds a new wireless LAN group to the Netbox inventory.
// It takes a context and a newWirelessLANGroup object as input and returns the created or updated wireless LAN group object and an error, if any.
// If the wireless LAN group already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the wireless LAN group does not exist, it creates a new one.
func (nbi *NetboxInventory) AddWirelessLANGroup(ctx context.Context, newWirelessLANGroup *objects.WirelessLANGroup) (*objects.WirelessLANGroup, error) {
	newWirelessLANGroup.Tags = append(newWirelessLANGroup.Tags, nbi.SsotTag)
	if _, ok := nbi.WirelessLANGroupsIndexByName[newWirelessLANGroup.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldWirelessLANGroup := nbi.WirelessLANGroupsIndexByName[newWirelessLANGroup.Name]
		delete(nbi.OrphanManager[constants.WirelessLANGroupsAPIPath], oldWirelessLANGroup.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newWirelessLANGroup, oldWirelessLANGroup, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug(ctx, "WirelessLANGroup ", newWirelessLANGroup.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedWirelessLANGroup, err := service.Patch[objects.WirelessLANGroup](ctx, nbi.NetboxAPI, oldWirelessLANGroup.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.WirelessLANGroupsIndexByName[newWirelessLANGroup.Name] = patchedWirelessLANGroup
		} else {
			nbi.Logger.Debug(ctx, "WirelessLANGroup ", newWirelessLANGroup.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug(ctx, "WirelessLANGroup ", newWirelessLANGroup.Name, " does not exist in Netbox. Creating it...")
		newWirelessLANGroup, err := service.Create[objects.WirelessLANGroup](ctx, nbi.NetboxAPI, newWirelessLANGroup)
		if err != nil {
			return nil, err
		}
		nbi.WirelessLANGroupsIndexByName[newWirelessLANGroup.Name] = newWirelessLANGroup
	}
	return nbi.WirelessLANGroupsIndexByName[newWirelessLANGroup.Name], nil
}

// Helper function that adds source name to custom field of the netbox object.
func addSourceNameCustomField(ctx context.Context, netboxObject *objects.NetboxObject) {
	if netboxObject.CustomFields == nil {
		netboxObject.CustomFields = make(map[string]interface{})
	}
	netboxObject.CustomFields[constants.CustomFieldSourceName] = ctx.Value(constants.CtxSourceKey).(string) //nolint
}
