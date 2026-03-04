package inventory

import (
	"context"
	"fmt"
	"strings"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
	"github.com/src-doo/netbox-ssot/internal/netbox/service"
	"github.com/src-doo/netbox-ssot/internal/utils"
)

// AddTag adds the newTag from source sourceName to the local inventory.
func (nbi *NetboxInventory) AddTag(ctx context.Context, newTag *objects.Tag) (*objects.Tag, error) {
	nbi.tagsLock.Lock()
	defer nbi.tagsLock.Unlock()
	if _, ok := nbi.tagsIndexByName[newTag.Name]; ok {
		oldTag := nbi.tagsIndexByName[newTag.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newTag, oldTag, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Tag %s already exists in Netbox but is out of date. Patching it... ",
				newTag.Name,
			)
			patchedTag, err := service.Patch[objects.Tag](ctx, nbi.NetboxAPI, oldTag.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.tagsIndexByName[newTag.Name] = patchedTag
		} else {
			nbi.Logger.Debugf(ctx, "Tag %s already exists in Netbox and is up to date...", newTag.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Tag %s does not exist in Netbox. Creating it...", newTag.Name)
		createdTag, err := service.Create(ctx, nbi.NetboxAPI, newTag)
		if err != nil {
			return nil, err
		}
		nbi.tagsIndexByName[newTag.Name] = createdTag
	}
	return nbi.tagsIndexByName[newTag.Name], nil
}

// AddTenants adds a new tenant to the local netbox inventory.
func (nbi *NetboxInventory) AddTenant(
	ctx context.Context,
	newTenant *objects.Tenant,
) (*objects.Tenant, error) {
	newTenant.NetboxObject.AddTag(nbi.SsotTag)
	nbi.tenantsLock.Lock()
	defer nbi.tenantsLock.Unlock()
	if _, ok := nbi.tenantsIndexByName[newTenant.Name]; ok {
		oldTenant := nbi.tenantsIndexByName[newTenant.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newTenant, oldTenant, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Tenant %s already exists in Netbox but is out of date. Patching it...",
				newTenant.Name,
			)
			patchedTenant, err := service.Patch[objects.Tenant](
				ctx,
				nbi.NetboxAPI,
				oldTenant.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.tenantsIndexByName[newTenant.Name] = patchedTenant
		} else {
			nbi.Logger.Debugf(ctx, "Tenant %s already exists in Netbox and is up to date...", newTenant.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Tenant %s does not exist in Netbox. Creating it...", newTenant.Name)
		createdTag, err := service.Create(ctx, nbi.NetboxAPI, newTenant)
		if err != nil {
			return nil, err
		}
		nbi.tenantsIndexByName[newTenant.Name] = createdTag
	}
	return nbi.tenantsIndexByName[newTenant.Name], nil
}

// AddSite adds a site to the local netbox inventory.
func (nbi *NetboxInventory) AddSite(
	ctx context.Context,
	newSite *objects.Site,
) (*objects.Site, error) {
	newSite.NetboxObject.AddTag(nbi.SsotTag)
	nbi.sitesLock.Lock()
	defer nbi.sitesLock.Unlock()
	if _, ok := nbi.sitesIndexByName[newSite.Name]; ok {
		oldSite := nbi.sitesIndexByName[newSite.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newSite, oldSite, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Site %s already exists in Netbox but is out of date. Patching it... ",
				newSite.Name,
			)
			patchedSite, err := service.Patch[objects.Site](ctx, nbi.NetboxAPI, oldSite.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.sitesIndexByName[newSite.Name] = patchedSite
		} else {
			nbi.Logger.Debugf(ctx, "Site %s already exists in Netbox and is up to date...", newSite.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Site %s does not exist in Netbox. Creating it...", newSite.Name)
		createdContact, err := service.Create(ctx, nbi.NetboxAPI, newSite)
		if err != nil {
			return nil, err
		}
		nbi.sitesIndexByName[newSite.Name] = createdContact
	}
	return nbi.sitesIndexByName[newSite.Name], nil
}

// AddSiteGroup adds a SiteGroup to the local netbox inventory.
func (nbi *NetboxInventory) AddSiteGroup(
	ctx context.Context,
	newSiteGroup *objects.SiteGroup,
) (*objects.SiteGroup, error) {
	newSiteGroup.NetboxObject.AddTag(nbi.SsotTag)
	nbi.siteGroupsLock.Lock()
	defer nbi.sitesLock.Unlock()
	if _, ok := nbi.siteGroupsIndexByName[newSiteGroup.Name]; ok {
		oldSiteGroup := nbi.siteGroupsIndexByName[newSiteGroup.Name]
		diffMap, err := utils.JSONDiffMapExceptID(
			newSiteGroup,
			oldSiteGroup,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"SiteGroup %s already exists in Netbox but is out of date. Patching it...",
				newSiteGroup.Name,
			)
			patchedSiteGroup, err := service.Patch[objects.SiteGroup](
				ctx,
				nbi.NetboxAPI,
				oldSiteGroup.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.siteGroupsIndexByName[newSiteGroup.Name] = patchedSiteGroup
		} else {
			nbi.Logger.Debugf(ctx, "SiteGroup %s already exists in Netbox and is up to date...", newSiteGroup.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "SiteGroup %s does not exist in Netbox. Creating it...", newSiteGroup.Name)
		createdSiteGroup, err := service.Create(ctx, nbi.NetboxAPI, newSiteGroup)
		if err != nil {
			return nil, err
		}
		nbi.siteGroupsIndexByName[newSiteGroup.Name] = createdSiteGroup
	}
	return nbi.siteGroupsIndexByName[newSiteGroup.Name], nil
}

// AddContactRole adds the newContactRole to the local netbox inventory.
func (nbi *NetboxInventory) AddContactRole(
	ctx context.Context,
	newContactRole *objects.ContactRole,
) (*objects.ContactRole, error) {
	newContactRole.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newContactRole.NetboxObject)
	nbi.contactRolesLock.Lock()
	defer nbi.contactRolesLock.Unlock()
	if _, ok := nbi.contactRolesIndexByName[newContactRole.Name]; ok {
		oldContactRole := nbi.contactRolesIndexByName[newContactRole.Name]
		diffMap, err := utils.JSONDiffMapExceptID(
			newContactRole,
			oldContactRole,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Contact role %s already exists in Netbox but is out of date. Patching it...",
				newContactRole.Name,
			)
			patchedContactRole, err := service.Patch[objects.ContactRole](
				ctx,
				nbi.NetboxAPI,
				oldContactRole.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.contactRolesIndexByName[newContactRole.Name] = patchedContactRole
		} else {
			nbi.Logger.Debugf(ctx, "Contact role %s already exists in Netbox and is up to date...", newContactRole.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Contact role %s does not exist in Netbox. Creating it...", newContactRole.Name)
		newContactRole, err := service.Create(ctx, nbi.NetboxAPI, newContactRole)
		if err != nil {
			return nil, err
		}
		nbi.contactRolesIndexByName[newContactRole.Name] = newContactRole
	}
	return nbi.contactRolesIndexByName[newContactRole.Name], nil
}

// AddContactGroup adds contact group to the local netbox inventory.
func (nbi *NetboxInventory) AddContactGroup(
	ctx context.Context,
	newContactGroup *objects.ContactGroup,
) (*objects.ContactGroup, error) {
	newContactGroup.NetboxObject.AddTag(nbi.SsotTag)
	nbi.contactGroupsLock.Lock()
	defer nbi.contactGroupsLock.Unlock()
	if _, ok := nbi.contactGroupsIndexByName[newContactGroup.Name]; ok {
		oldContactGroup := nbi.contactGroupsIndexByName[newContactGroup.Name]
		diffMap, err := utils.JSONDiffMapExceptID(
			newContactGroup,
			oldContactGroup,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Contact group %s already exists in Netbox but is out of date. Patching it...",
				newContactGroup.Name,
			)
			patchedContactGroup, err := service.Patch[objects.ContactGroup](
				ctx,
				nbi.NetboxAPI,
				oldContactGroup.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.contactGroupsIndexByName[newContactGroup.Name] = patchedContactGroup
		} else {
			nbi.Logger.Debugf(ctx, "Contact group %s already exists in Netbox and is up to date...", newContactGroup.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Contact group %s does not exist in Netbox. Creating it...", newContactGroup.Name)
		newContactGroup, err := service.Create(ctx, nbi.NetboxAPI, newContactGroup)
		if err != nil {
			return nil, err
		}
		nbi.contactGroupsIndexByName[newContactGroup.Name] = newContactGroup
	}
	return nbi.contactGroupsIndexByName[newContactGroup.Name], nil
}

// AddContact adds a contact to the local netbox inventory.
func (nbi *NetboxInventory) AddContact(
	ctx context.Context,
	newContact *objects.Contact,
) (*objects.Contact, error) {
	newContact.NetboxObject.AddTag(nbi.SsotTag)
	newContact.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.contactsLock.Lock()
	defer nbi.contactsLock.Unlock()
	if _, ok := nbi.contactsIndexByName[newContact.Name]; ok {
		oldContact := nbi.contactsIndexByName[newContact.Name]
		nbi.OrphanManager.RemoveItem(oldContact)
		diffMap, err := utils.JSONDiffMapExceptID(newContact, oldContact, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Contact %s already exists in Netbox but is out of date. Patching it...",
				newContact.Name,
			)
			patchedContact, err := service.Patch[objects.Contact](
				ctx,
				nbi.NetboxAPI,
				oldContact.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.contactsIndexByName[newContact.Name] = patchedContact
		} else {
			nbi.Logger.Debugf(ctx, "Contact %s already exists in Netbox and is up to date...", newContact.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Contact %s does not exist in Netbox. Creating it...", newContact.Name)
		createdContact, err := service.Create(ctx, nbi.NetboxAPI, newContact)
		if err != nil {
			return nil, err
		}
		nbi.contactsIndexByName[newContact.Name] = createdContact
	}
	return nbi.contactsIndexByName[newContact.Name], nil
}

// AddContact assignment adds a contact assignment to the local netbox inventory.
// TODO: Make index check less code and more universal, checking each level is ugly.
func (nbi *NetboxInventory) AddContactAssignment(
	ctx context.Context,
	newCA *objects.ContactAssignment,
) (*objects.ContactAssignment, error) {
	newCA.NetboxObject.AddTag(nbi.SsotTag)
	newCA.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.contactAssignmentsLock.Lock()
	defer nbi.contactAssignmentsLock.Unlock()
	if nbi.contactAssignmentsIndex[newCA.ModelType] == nil {
		nbi.contactAssignmentsIndex[newCA.ModelType] = make(
			map[int]map[int]map[int]*objects.ContactAssignment,
		)
	}
	if nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID] == nil {
		nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID] = make(
			map[int]map[int]*objects.ContactAssignment,
		)
	}
	if nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID] == nil {
		nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID] = make(
			map[int]*objects.ContactAssignment,
		)
	}
	newCA.Tags = append(newCA.Tags, nbi.SsotTag)
	if _, ok := nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID]; ok {
		oldCA := nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID]
		nbi.OrphanManager.RemoveItem(oldCA)
		diffMap, err := utils.JSONDiffMapExceptID(newCA, oldCA, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"ContactAssignment %d already exists in Netbox but is out of date. Patching it...",
				newCA.ID,
			)
			patchedCA, err := service.Patch[objects.ContactAssignment](
				ctx,
				nbi.NetboxAPI,
				oldCA.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID] = patchedCA
		} else {
			nbi.Logger.Debugf(ctx, "ContactAssignment %d already exists in Netbox and is up to date...", newCA.ID)
		}
	} else {
		nbi.Logger.Debugf(ctx, "ContactAssignment %s does not exist in Netbox. Creating it...", newCA)
		newCA, err := service.Create(ctx, nbi.NetboxAPI, newCA)
		if err != nil {
			return nil, err
		}
		nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID] = newCA
	}
	return nbi.contactAssignmentsIndex[newCA.ModelType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID], nil
}

// AddCustomField adds a custom field to the Netbox inventory.
// It takes a context and a newCf object as input and
// returns the created or patched custom field along with any error encountered.
// If the custom field already exists in Netbox but is out of date, it will be patched with the new values.
// If the custom field does not exist, it will be created.
func (nbi *NetboxInventory) AddCustomField(
	ctx context.Context,
	newCf *objects.CustomField,
) (*objects.CustomField, error) {
	nbi.customFieldsLock.Lock()
	defer nbi.customFieldsLock.Unlock()
	if _, ok := nbi.customFieldsIndexByName[newCf.Name]; ok {
		oldCustomField := nbi.customFieldsIndexByName[newCf.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newCf, oldCustomField, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Custom field %s already exists in Netbox but is out of date. Patching it...",
				newCf.Name,
			)
			patchedCf, err := service.Patch[objects.CustomField](
				ctx,
				nbi.NetboxAPI,
				oldCustomField.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.customFieldsIndexByName[newCf.Name] = patchedCf
		} else {
			nbi.Logger.Debugf(ctx, "Custom field %s already exists in Netbox and is up to date...", newCf.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Custom field %s does not exist in Netbox. Creating it...", newCf.Name)
		createdCf, err := service.Create(ctx, nbi.NetboxAPI, newCf)
		if err != nil {
			return nil, err
		}
		nbi.customFieldsIndexByName[createdCf.Name] = createdCf
	}
	return nbi.customFieldsIndexByName[newCf.Name], nil
}

// AddClusterGroup adds a new cluster group to the Netbox inventory.
// It takes a context and a newCg object as input and returns the newly created cluster group and an error (if any).
// If the cluster group already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the cluster group does not exist, it creates a new one.
// The function also updates the cluster group index by name and removes the ID from the orphan manager.
func (nbi *NetboxInventory) AddClusterGroup(
	ctx context.Context,
	newCg *objects.ClusterGroup,
) (*objects.ClusterGroup, error) {
	newCg.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newCg.NetboxObject)
	newCg.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.clusterGroupsLock.Lock()
	defer nbi.clusterGroupsLock.Unlock()
	if _, ok := nbi.clusterGroupsIndexByName[newCg.Name]; ok {
		oldCg := nbi.clusterGroupsIndexByName[newCg.Name]
		nbi.OrphanManager.RemoveItem(oldCg)
		diffMap, err := utils.JSONDiffMapExceptID(newCg, oldCg, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Cluster group %s already exists in Netbox but is out of date. Patching it...",
				newCg.Name,
			)
			patchedCg, err := service.Patch[objects.ClusterGroup](
				ctx,
				nbi.NetboxAPI,
				oldCg.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.clusterGroupsIndexByName[newCg.Name] = patchedCg
		} else {
			nbi.Logger.Debugf(ctx, "Cluster group %s already exists in Netbox and is up to date...", newCg.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Cluster group %s does not exist in Netbox. Creating it...", newCg.Name)
		newCg, err := service.Create(ctx, nbi.NetboxAPI, newCg)
		if err != nil {
			return nil, err
		}
		nbi.clusterGroupsIndexByName[newCg.Name] = newCg
	}
	// Delete id from orphan manager
	return nbi.clusterGroupsIndexByName[newCg.Name], nil
}

// AddClusterType adds a new cluster type to the Netbox inventory.
// It takes a context and a newClusterType object as input and
// returns the created or updated cluster type object and an error, if any.
// If the cluster type already exists in Netbox, it checks if it is up to date.
// If not, it patches the existing cluster type.
// If the cluster type does not exist, it creates a new one.
func (nbi *NetboxInventory) AddClusterType(
	ctx context.Context,
	newClusterType *objects.ClusterType,
) (*objects.ClusterType, error) {
	newClusterType.NetboxObject.AddTag(nbi.SsotTag)
	newClusterType.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.clusterTypesLock.Lock()
	defer nbi.clusterTypesLock.Unlock()
	if _, ok := nbi.clusterTypesIndexByName[newClusterType.Name]; ok {
		oldClusterType := nbi.clusterTypesIndexByName[newClusterType.Name]
		nbi.OrphanManager.RemoveItem(oldClusterType)
		diffMap, err := utils.JSONDiffMapExceptID(
			newClusterType,
			oldClusterType,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Cluster type %s already exists in Netbox but is out of date. Patching it...",
				newClusterType.Name,
			)
			patchedClusterType, err := service.Patch[objects.ClusterType](
				ctx,
				nbi.NetboxAPI,
				oldClusterType.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.clusterTypesIndexByName[newClusterType.Name] = patchedClusterType
			return patchedClusterType, nil
		}
		nbi.Logger.Debugf(
			ctx,
			"Cluster type %s already exists in Netbox and is up to date...",
			newClusterType.Name,
		)
		existingClusterType := nbi.clusterTypesIndexByName[newClusterType.Name]
		return existingClusterType, nil
	}
	nbi.Logger.Debugf(
		ctx,
		"Cluster type %s does not exist in Netbox. Creating it...",
		newClusterType.Name,
	)
	newClusterType, err := service.Create(ctx, nbi.NetboxAPI, newClusterType)
	if err != nil {
		return nil, err
	}
	nbi.clusterTypesIndexByName[newClusterType.Name] = newClusterType
	return newClusterType, nil
}

// AddCluster adds a new cluster to the Netbox inventory.
// It takes a context and a pointer to a Cluster object as input.
// It returns the newly created cluster object and an error, if any.
// If the cluster already exists in Netbox, it checks if the existing cluster is up to date.
// If it is not up to date, it patches the existing cluster with the changes from the new cluster.
// If the cluster does not exist in Netbox, it creates a new cluster.
func (nbi *NetboxInventory) AddCluster(
	ctx context.Context,
	newCluster *objects.Cluster,
) (*objects.Cluster, error) {
	newCluster.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newCluster.NetboxObject)
	newCluster.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.clustersLock.Lock()
	defer nbi.clustersLock.Unlock()
	if _, ok := nbi.clustersIndexByName[newCluster.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldCluster := nbi.clustersIndexByName[newCluster.Name]
		nbi.OrphanManager.RemoveItem(oldCluster)
		diffMap, err := utils.JSONDiffMapExceptID(newCluster, oldCluster, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Cluster %s already exists in Netbox but is out of date. Patching it...",
				newCluster.Name,
			)
			patchedCluster, err := service.Patch[objects.Cluster](
				ctx,
				nbi.NetboxAPI,
				oldCluster.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.clustersIndexByName[newCluster.Name] = patchedCluster
		} else {
			nbi.Logger.Debugf(ctx, "Cluster %s already exists in Netbox and is up to date...", newCluster.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Cluster %s does not exist in Netbox. Creating it...", newCluster.Name)
		createdCluster, err := service.Create(ctx, nbi.NetboxAPI, newCluster)
		if err != nil {
			return nil, err
		}
		nbi.clustersIndexByName[createdCluster.Name] = createdCluster
	}
	return nbi.clustersIndexByName[newCluster.Name], nil
}

// AddDeviceRole adds a new device role to the Netbox inventory.
// It takes a context and a newDeviceRole object as input and
// returns the created device role object and an error, if any.
// If the device role already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the device role does not exist, it creates a new one.
func (nbi *NetboxInventory) AddDeviceRole(
	ctx context.Context,
	newDeviceRole *objects.DeviceRole,
) (*objects.DeviceRole, error) {
	newDeviceRole.NetboxObject.AddTag(nbi.SsotTag)
	newDeviceRole.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.deviceRolesLock.Lock()
	defer nbi.deviceRolesLock.Unlock()
	if _, ok := nbi.deviceRolesIndexByName[newDeviceRole.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldDeviceRole := nbi.deviceRolesIndexByName[newDeviceRole.Name]
		nbi.OrphanManager.RemoveItem(oldDeviceRole)
		diffMap, err := utils.JSONDiffMapExceptID(
			newDeviceRole,
			oldDeviceRole,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Device role %s already exists in Netbox but is out of date. Patching it...",
				newDeviceRole.Name,
			)
			patchedDeviceRole, err := service.Patch[objects.DeviceRole](
				ctx,
				nbi.NetboxAPI,
				oldDeviceRole.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.deviceRolesIndexByName[newDeviceRole.Name] = patchedDeviceRole
		} else {
			nbi.Logger.Debugf(ctx, "Device role %s already exists in Netbox and is up to date...", newDeviceRole.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Device role %s does not exist in Netbox. Creating it...", newDeviceRole.Name)
		newDeviceRole, err := service.Create(ctx, nbi.NetboxAPI, newDeviceRole)
		if err != nil {
			return nil, err
		}
		nbi.deviceRolesIndexByName[newDeviceRole.Name] = newDeviceRole
	}
	return nbi.deviceRolesIndexByName[newDeviceRole.Name], nil
}

// AddManufacturer adds a new manufacturer to the Netbox inventory.
// It takes a context, `ctx`, and a pointer to a `newManufacturer` object as input.
// The function returns a pointer to the newly created manufacturer and an error, if any.
// If the manufacturer already exists in Netbox, the function checks if it is up to date.
// If it is not up to date, the function patches the existing manufacturer with the updated information.
// If the manufacturer does not exist in Netbox, the function creates a new one.
func (nbi *NetboxInventory) AddManufacturer(
	ctx context.Context,
	newManufacturer *objects.Manufacturer,
) (*objects.Manufacturer, error) {
	newManufacturer.NetboxObject.AddTag(nbi.SsotTag)
	newManufacturer.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.manufacturersLock.Lock()
	defer nbi.manufacturersLock.Unlock()
	if _, ok := nbi.manufacturersIndexByName[newManufacturer.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldManufacturer := nbi.manufacturersIndexByName[newManufacturer.Name]
		nbi.OrphanManager.RemoveItem(oldManufacturer)
		diffMap, err := utils.JSONDiffMapExceptID(
			newManufacturer,
			oldManufacturer,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Manufacturer %s already exists in Netbox but is out of date. Patching it...",
				newManufacturer.Name,
			)
			patchedManufacturer, err := service.Patch[objects.Manufacturer](
				ctx,
				nbi.NetboxAPI,
				oldManufacturer.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.manufacturersIndexByName[newManufacturer.Name] = patchedManufacturer
		} else {
			nbi.Logger.Debugf(ctx, "Manufacturer %s already exists in Netbox and is up to date...", newManufacturer.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Manufacturer %s does not exist in Netbox. Creating it...", newManufacturer.Name)
		newManufacturer, err := service.Create(ctx, nbi.NetboxAPI, newManufacturer)
		if err != nil {
			return nil, err
		}
		nbi.manufacturersIndexByName[newManufacturer.Name] = newManufacturer
	}
	return nbi.manufacturersIndexByName[newManufacturer.Name], nil
}

// AddDeviceType adds a new device type to the Netbox inventory.
// It takes a context and a newDeviceType object as input and
// returns the created or updated device type object and an error, if any.
// If the device type already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the device type does not exist, it creates a new one.
func (nbi *NetboxInventory) AddDeviceType(
	ctx context.Context,
	newDeviceType *objects.DeviceType,
) (*objects.DeviceType, error) {
	newDeviceType.NetboxObject.AddTag(nbi.SsotTag)
	newDeviceType.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.deviceTypesLock.Lock()
	defer nbi.deviceTypesLock.Unlock()
	if _, ok := nbi.deviceTypesIndexByModel[newDeviceType.Model]; ok {
		oldDeviceType := nbi.deviceTypesIndexByModel[newDeviceType.Model]
		nbi.OrphanManager.RemoveItem(oldDeviceType)
		diffMap, err := utils.JSONDiffMapExceptID(
			newDeviceType,
			oldDeviceType,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Device type %s already exists in Netbox but is out of date. Patching it...",
				newDeviceType.Model,
			)
			patchedDeviceType, err := service.Patch[objects.DeviceType](
				ctx,
				nbi.NetboxAPI,
				oldDeviceType.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.deviceTypesIndexByModel[newDeviceType.Model] = patchedDeviceType
		} else {
			nbi.Logger.Debugf(ctx, "Device type %s already exists in Netbox and is up to date...", newDeviceType.Model)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Device type %s does not exist in Netbox. Creating it...", newDeviceType.Model)
		newDeviceType, err := service.Create(ctx, nbi.NetboxAPI, newDeviceType)
		if err != nil {
			return nil, err
		}
		nbi.deviceTypesIndexByModel[newDeviceType.Model] = newDeviceType
	}
	return nbi.deviceTypesIndexByModel[newDeviceType.Model], nil
}

// AddPlatform adds a new platform to the Netbox inventory.
// It takes a context and a newPlatform object as input and
// returns the created or updated platform object and an error, if any.
// If the platform already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the platform does not exist, it creates a new one.
func (nbi *NetboxInventory) AddPlatform(
	ctx context.Context,
	newPlatform *objects.Platform,
) (*objects.Platform, error) {
	newPlatform.NetboxObject.AddTag(nbi.SsotTag)
	newPlatform.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.platformsLock.Lock()
	defer nbi.platformsLock.Unlock()
	if _, ok := nbi.platformsIndexByName[newPlatform.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldPlatform := nbi.platformsIndexByName[newPlatform.Name]
		nbi.OrphanManager.RemoveItem(oldPlatform)
		diffMap, err := utils.JSONDiffMapExceptID(
			newPlatform,
			oldPlatform,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Platform %s already exists in Netbox but is out of date. Patching it...",
				newPlatform.Name,
			)
			patchedPlatform, err := service.Patch[objects.Platform](
				ctx,
				nbi.NetboxAPI,
				oldPlatform.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.platformsIndexByName[newPlatform.Name] = patchedPlatform
		} else {
			nbi.Logger.Debugf(ctx, "Platform %s already exists in Netbox and is up to date...", newPlatform.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Platform %s does not exist in Netbox. Creating it...", newPlatform.Name)
		newPlatform, err := service.Create(ctx, nbi.NetboxAPI, newPlatform)
		if err != nil {
			return nil, err
		}
		nbi.platformsIndexByName[newPlatform.Name] = newPlatform
	}
	return nbi.platformsIndexByName[newPlatform.Name], nil
}

// AddRackRole adds a new rack role to the Netbox inventory.
// It takes a context and a newRackRole object as input and
// returns the created or updated rack role object and an error, if any.
// If the rack role already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the rack role does not exist, it creates a new one.
func (nbi *NetboxInventory) AddDevice(
	ctx context.Context,
	newDevice *objects.Device,
) (*objects.Device, error) {
	newDevice.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newDevice.NetboxObject)
	nbi.applyDeviceFieldLengthLimitations(newDevice)
	newDevice.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.devicesLock.Lock()
	defer nbi.devicesLock.Unlock()
	if newDevice.Site == nil {
		return nil, fmt.Errorf("device %s is not assigned to a site, but it should be", newDevice)
	}
	if _, ok := nbi.devicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID]; ok {
		oldDevice := nbi.devicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID]
		nbi.OrphanManager.RemoveItem(oldDevice)

		// Allow manual override device type
		if newDevice.DeviceType != nil && oldDevice.DeviceType != nil &&
			newDevice.DeviceType.ID != oldDevice.DeviceType.ID &&
			oldDevice.NetboxObject.HasTag(nbi.IgnoreDeviceTypeTag) {
			// Preserve manually set device type from NetBox and keep ignore tag
			newDevice.DeviceType = oldDevice.DeviceType
			newDevice.NetboxObject.AddTag(nbi.IgnoreDeviceTypeTag)
		}

		diffMap, err := utils.JSONDiffMapExceptID(newDevice, oldDevice, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Device %s already exists in Netbox but is out of date. Patching it...",
				newDevice.Name,
			)
			patchedDevice, err := service.Patch[objects.Device](
				ctx,
				nbi.NetboxAPI,
				oldDevice.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.devicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID] = patchedDevice
			nbi.devicesIndexByID[patchedDevice.ID] = patchedDevice
		} else {
			nbi.Logger.Debugf(ctx, "Device %s already exists in Netbox and is up to date...", newDevice.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Device %s does not exist in Netbox. Creating it...", newDevice.Name)
		newDevice, err := service.Create(ctx, nbi.NetboxAPI, newDevice)
		if err != nil {
			return nil, err
		}
		if nbi.devicesIndexByNameAndSiteID[newDevice.Name] == nil {
			nbi.devicesIndexByNameAndSiteID[newDevice.Name] = make(map[int]*objects.Device)
		}
		nbi.devicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID] = newDevice
		nbi.devicesIndexByID[newDevice.ID] = newDevice
	}
	return nbi.devicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID], nil
}

// AddVirtualDeviceContext adds new virtual device context to the local inventory.
// It takes a context and a newVDC object as input and
// returns the created or updated virtual device context object and an error, if any.
// If the virtual device context already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the virtual device context does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVirtualDeviceContext(
	ctx context.Context,
	newVDC *objects.VirtualDeviceContext,
) (*objects.VirtualDeviceContext, error) {
	newVDC.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newVDC.NetboxObject)
	newVDC.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.virtualDeviceContextsLock.Lock()
	defer nbi.virtualDeviceContextsLock.Unlock()
	if newVDC.Device == nil {
		return nil, fmt.Errorf(
			"VirtualDeviceContext %s is not assigned to a device, but it should be",
			newVDC,
		)
	}
	if _, ok := nbi.virtualDeviceContextsIndex[newVDC.Name][newVDC.Device.ID]; ok {
		oldVDC := nbi.virtualDeviceContextsIndex[newVDC.Name][newVDC.Device.ID]
		nbi.OrphanManager.RemoveItem(oldVDC)
		diffMap, err := utils.JSONDiffMapExceptID(newVDC, oldVDC, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"VirtualDeviceContext %s already exists in Netbox but is out of date. Patching it...",
				newVDC.Name,
			)
			patchedVDC, err := service.Patch[objects.VirtualDeviceContext](
				ctx,
				nbi.NetboxAPI,
				oldVDC.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.virtualDeviceContextsIndex[newVDC.Name][newVDC.Device.ID] = patchedVDC
		} else {
			nbi.Logger.Debugf(ctx, "VirtualDeviceContext %s already exists in Netbox and is up to date...", newVDC.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "VirtualDeviceContext %s does not exist in Netbox. Creating it...", newVDC.Name)
		newDevice, err := service.Create(ctx, nbi.NetboxAPI, newVDC)
		if err != nil {
			return nil, err
		}
		if nbi.virtualDeviceContextsIndex[newDevice.Name] == nil {
			nbi.virtualDeviceContextsIndex[newDevice.Name] = make(map[int]*objects.VirtualDeviceContext)
		}
		nbi.virtualDeviceContextsIndex[newDevice.Name][newDevice.Device.ID] = newDevice
	}
	return nbi.virtualDeviceContextsIndex[newVDC.Name][newVDC.Device.ID], nil
}

// AddVlanGroup adds a new vlan group to the Netbox inventory.
// It takes a context and a newVlanGroup object as input and
// returns the created or updated vlan group object and an error, if any.
// If the vlan group already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the vlan group does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVlanGroup(
	ctx context.Context,
	newVlanGroup *objects.VlanGroup,
) (*objects.VlanGroup, error) {
	newVlanGroup.NetboxObject.AddTag(nbi.SsotTag)
	newVlanGroup.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.vlanGroupsLock.Lock()
	defer nbi.vlanGroupsLock.Unlock()
	if _, ok := nbi.vlanGroupsIndexByName[newVlanGroup.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldVlanGroup := nbi.vlanGroupsIndexByName[newVlanGroup.Name]
		nbi.OrphanManager.RemoveItem(oldVlanGroup)
		diffMap, err := utils.JSONDiffMapExceptID(
			newVlanGroup,
			oldVlanGroup,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"VlanGroup %s already exists in Netbox but is out of date. Patching it...",
				newVlanGroup.Name,
			)
			patchedVlanGroup, err := service.Patch[objects.VlanGroup](
				ctx,
				nbi.NetboxAPI,
				oldVlanGroup.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.vlanGroupsIndexByName[newVlanGroup.Name] = patchedVlanGroup
		} else {
			nbi.Logger.Debugf(ctx, "VlanGroup %s already exists in Netbox and is up to date...", newVlanGroup.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "VlanGroup %s does not exist in Netbox. Creating it...", newVlanGroup.Name)
		newVlan, err := service.Create(ctx, nbi.NetboxAPI, newVlanGroup)
		if err != nil {
			return nil, err
		}
		nbi.vlanGroupsIndexByName[newVlan.Name] = newVlan
	}
	return nbi.vlanGroupsIndexByName[newVlanGroup.Name], nil
}

// AddVlan adds a new vlan to the Netbox inventory.
// It takes a context and a newVlan object as input and returns the created or updated vlan object and an error, if any.
// If the vlan already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the vlan does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVlan(
	ctx context.Context,
	newVlan *objects.Vlan,
) (*objects.Vlan, error) {
	newVlan.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newVlan.NetboxObject)
	newVlan.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.vlansLock.Lock()
	defer nbi.vlansLock.Unlock()
	if _, ok := nbi.vlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldVlan := nbi.vlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid]
		nbi.OrphanManager.RemoveItem(oldVlan)
		diffMap, err := utils.JSONDiffMapExceptID(newVlan, oldVlan, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Vlan %s already exists in Netbox but is out of date. Patching it...",
				newVlan.Name,
			)
			patchedVlan, err := service.Patch[objects.Vlan](ctx, nbi.NetboxAPI, oldVlan.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.vlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid] = patchedVlan
		} else {
			nbi.Logger.Debugf(ctx, "Vlan %s already exists in Netbox and is up to date...", newVlan.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Vlan %s does not exist in Netbox. Creating it...", newVlan.Name)
		newVlan, err := service.Create(ctx, nbi.NetboxAPI, newVlan)
		if err != nil {
			return nil, err
		}
		if nbi.vlansIndexByVlanGroupIDAndVID[newVlan.Group.ID] == nil {
			nbi.vlansIndexByVlanGroupIDAndVID[newVlan.Group.ID] = make(map[int]*objects.Vlan)
		}
		nbi.vlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid] = newVlan
	}
	return nbi.vlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid], nil
}

// AddInterface adds a new interface to the Netbox inventory.
// It takes a context and a newInterface object as input and
// returns the created or updated interface object and an error, if any.
// If the interface already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the interface does not exist, it creates a new one.
func (nbi *NetboxInventory) AddInterface(
	ctx context.Context,
	newInterface *objects.Interface,
) (*objects.Interface, error) {
	newInterface.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newInterface.NetboxObject)
	newInterface.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	if len(newInterface.Name) > constants.MaxInterfaceNameLength {
		newInterface.Name = newInterface.Name[:constants.MaxInterfaceNameLength]
	}
	nbi.interfacesLock.Lock()
	defer nbi.interfacesLock.Unlock()
	if _, ok := nbi.interfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name]; ok {
		oldInterface := nbi.interfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name]
		nbi.OrphanManager.RemoveItem(oldInterface)
		diffMap, err := utils.JSONDiffMapExceptID(
			newInterface,
			oldInterface,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Interface %s/%s already exists in Netbox but is out of date. Patching it...",
				newInterface.Device.Name,
				newInterface.Name,
			)
			patchedInterface, err := service.Patch[objects.Interface](
				ctx,
				nbi.NetboxAPI,
				oldInterface.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.interfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name] = patchedInterface
			nbi.interfacesIndexByID[patchedInterface.ID] = patchedInterface
		} else {
			nbi.Logger.Debugf(
				ctx,
				"Interface %s/%s already exists in Netbox and is up to date...",
				newInterface.Device.Name, newInterface.Name,
			)
		}
	} else {
		nbi.Logger.Debugf(
			ctx,
			"Interface %s/%s does not exist in Netbox. Creating it...",
			newInterface.Device.Name, newInterface.Name,
		)
		newInterface, err := service.Create(ctx, nbi.NetboxAPI, newInterface)
		if err != nil {
			return nil, err
		}
		if nbi.interfacesIndexByDeviceIDAndName[newInterface.Device.ID] == nil {
			nbi.interfacesIndexByDeviceIDAndName[newInterface.Device.ID] = make(map[string]*objects.Interface)
		}
		nbi.interfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name] = newInterface
		nbi.interfacesIndexByID[newInterface.ID] = newInterface
		return newInterface, nil
	}
	return nbi.interfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name], nil
}

// AddVM adds a new virtual machine to the Netbox inventory.
// It takes a context and a newVM object as input and
// returns the created or updated virtual machine object and an error, if any.
// If the virtual machine already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the virtual machine does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVM(ctx context.Context, newVM *objects.VM) (*objects.VM, error) {
	newVM.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newVM.NetboxObject)
	newVM.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.vmsLock.Lock()
	defer nbi.vmsLock.Unlock()
	newVMClusterID := -1
	if newVM.Cluster != nil {
		newVMClusterID = newVM.Cluster.ID
	}
	if len(newVM.Name) > constants.MaxVMNameLength {
		newVM.Name = newVM.Name[:constants.MaxVMNameLength]
	}
	if oldVM, ok := nbi.vmsIndexByNameAndClusterID[newVM.Name][newVMClusterID]; ok {
		nbi.OrphanManager.RemoveItem(oldVM)
		diffMap, err := utils.JSONDiffMapExceptID(newVM, oldVM, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"VM %s already exists in Netbox but is out of date. Patching it...",
				newVM,
			)
			patchedVM, err := service.Patch[objects.VM](ctx, nbi.NetboxAPI, oldVM.ID, diffMap)
			if err != nil {
				nbi.Logger.Errorf(ctx, "Error while patching %s : %s", newVM.Name, err)
				return nil, err
			}
			nbi.vmsIndexByNameAndClusterID[newVM.Name][newVMClusterID] = patchedVM
			nbi.vmsIndexByID[patchedVM.ID] = patchedVM
		} else {
			nbi.Logger.Debugf(ctx, "VM %s already exists in Netbox and is up to date...", newVM)
		}
	} else {
		nbi.Logger.Debugf(ctx, "VM %s does not exist in Netbox. Creating it...", newVM)
		newVM, err := service.Create(ctx, nbi.NetboxAPI, newVM)
		if err != nil {
			return nil, err
		}
		if nbi.vmsIndexByNameAndClusterID[newVM.Name] == nil {
			nbi.vmsIndexByNameAndClusterID[newVM.Name] = make(map[int]*objects.VM)
		}
		nbi.vmsIndexByNameAndClusterID[newVM.Name][newVMClusterID] = newVM
		nbi.vmsIndexByID[newVM.ID] = newVM
		return newVM, nil
	}
	return nbi.vmsIndexByNameAndClusterID[newVM.Name][newVMClusterID], nil
}

// AddVMInterface adds a new virtual machine interface to the Netbox inventory.
// It takes a context and a newVMInterface object as input and
// returns the created or updated virtual machine interface object and an error, if any.
// If the virtual machine interface already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the virtual machine interface does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVMInterface(
	ctx context.Context,
	newVMInterface *objects.VMInterface,
) (*objects.VMInterface, error) {
	newVMInterface.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newVMInterface.NetboxObject)
	newVMInterface.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.vmInterfacesLock.Lock()
	defer nbi.vmInterfacesLock.Unlock()
	if len(newVMInterface.Name) > constants.MaxVMInterfaceNameLength {
		newVMInterface.Name = newVMInterface.Name[:constants.MaxVMInterfaceNameLength]
	}
	if _, ok := nbi.vmInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name]; ok {
		oldVMIface := nbi.vmInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name]
		nbi.OrphanManager.RemoveItem(oldVMIface)
		diffMap, err := utils.JSONDiffMapExceptID(
			newVMInterface,
			oldVMIface,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"VM interface %s already exists in Netbox but is out of date. Patching it...",
				newVMInterface.Name,
			)
			patchedVMInterface, err := service.Patch[objects.VMInterface](
				ctx,
				nbi.NetboxAPI,
				oldVMIface.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.vmInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name] = patchedVMInterface
			nbi.vmInterfacesIndexByID[patchedVMInterface.ID] = patchedVMInterface
		} else {
			nbi.Logger.Debugf(ctx, "VM interface %s already exists in Netbox and is up to date...", newVMInterface.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "VM interface %s does not exist in Netbox. Creating it...", newVMInterface.Name)
		newVMInterface, err := service.Create(ctx, nbi.NetboxAPI, newVMInterface)
		if err != nil {
			return nil, err
		}
		if nbi.vmInterfacesIndexByVMIdAndName[newVMInterface.VM.ID] == nil {
			nbi.vmInterfacesIndexByVMIdAndName[newVMInterface.VM.ID] = make(map[string]*objects.VMInterface)
		}
		nbi.vmInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name] = newVMInterface
		nbi.vmInterfacesIndexByID[newVMInterface.ID] = newVMInterface
	}
	return nbi.vmInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name], nil
}

// AddIPAddress adds a new IP address to the Netbox inventory.
// It takes a context and a newIPAddress object as input and
// returns the created or updated IP address object and an error, if any.
// If the IP address already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the IP address does not exist, it creates a new one.
func (nbi *NetboxInventory) AddIPAddress(
	ctx context.Context,
	newIPAddress *objects.IPAddress,
) (*objects.IPAddress, error) {
	newIPAddress.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newIPAddress.NetboxObject)
	newIPAddress.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)

	objType, objName, ifaceName, err := nbi.getIndexValuesForIPAddress(newIPAddress)
	if err != nil {
		return nil, fmt.Errorf("get index values for ip address %+v: %s", newIPAddress, err)
	}
	nbi.verifyIPAddressIndexExists(objType, objName, ifaceName)

	indexKey := ipAddressIndexKey(newIPAddress) // ← clé composite

	nbi.ipAddressesLock.Lock()
	defer nbi.ipAddressesLock.Unlock()
	if _, ok := nbi.ipAddressesIndex[objType][objName][ifaceName][indexKey]; ok {
		oldIPAddress := nbi.ipAddressesIndex[objType][objName][ifaceName][indexKey]
		nbi.OrphanManager.RemoveItem(oldIPAddress)
		diffMap, err := utils.JSONDiffMapExceptID(
			newIPAddress,
			oldIPAddress,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"IP address %s already exists in Netbox but is out of date. Patching it...",
				newIPAddress.Address,
			)
			patchedIPAddress, err := service.Patch[objects.IPAddress](
				ctx,
				nbi.NetboxAPI,
				oldIPAddress.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.ipAddressesIndex[objType][objName][ifaceName][indexKey] = patchedIPAddress
			return patchedIPAddress, nil
		}
		nbi.Logger.Debugf(
			ctx,
			"IP address %s already exists in Netbox and is up to date...",
			newIPAddress.Address,
		)
	} else {
		nbi.Logger.Debugf(ctx, "IP address %s does not exist in Netbox. Creating it...", newIPAddress.Address)
		newIPAddress, err := service.Create(ctx, nbi.NetboxAPI, newIPAddress)
		if err != nil {
			return nil, err
		}
		nbi.ipAddressesIndex[objType][objName][ifaceName][indexKey] = newIPAddress
		return newIPAddress, nil
	}
	return nbi.ipAddressesIndex[objType][objName][ifaceName][indexKey], nil
}

// AddMACAddress adds a new MAC address to the Netbox inventory.
// It takes a context and a newMACAddress object as input and
// returns the created or updated MAC address object and an error, if any.
// If the MAC address already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the MAC address does not exist, it creates a new one.
func (nbi *NetboxInventory) AddMACAddress(
	ctx context.Context,
	newMACAddress *objects.MACAddress,
) (*objects.MACAddress, error) {
	newMACAddress.NetboxObject.AddTag(nbi.SsotTag)
	newMACAddress.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)

	// Get index values with helper function.
	objType, objName, ifaceName, err := nbi.getIndexValuesForMACAddress(newMACAddress)
	if err != nil {
		return nil, fmt.Errorf("get index values for mac address %+v: %s", newMACAddress, err)
	}

	// ensure index is not nil
	nbi.verifyMACAddressIndexExists(objType, objName, ifaceName)

	// ensure MAC address is uppercase
	newMACAddress.MAC = strings.ToUpper(newMACAddress.MAC)

	nbi.macAddressesLock.Lock()
	defer nbi.macAddressesLock.Unlock()
	if _, ok := nbi.macAddressesIndex[objType][objName][ifaceName][newMACAddress.MAC]; ok {
		oldMACAddress := nbi.macAddressesIndex[objType][objName][ifaceName][newMACAddress.MAC]
		nbi.OrphanManager.RemoveItem(oldMACAddress)

		diffMap, err := utils.JSONDiffMapExceptID(
			newMACAddress,
			oldMACAddress,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}

		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"MAC address %s already exists in Netbox but is out of date. Patching it...",
				newMACAddress.MAC,
			)
			patchedMACAddress, err := service.Patch[objects.MACAddress](
				ctx,
				nbi.NetboxAPI,
				oldMACAddress.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.macAddressesIndex[objType][objName][ifaceName][newMACAddress.MAC] = patchedMACAddress
			return patchedMACAddress, nil
		}
		nbi.Logger.Debugf(
			ctx,
			"MAC address %s already exists in Netbox and is up to date...",
			newMACAddress.MAC,
		)
	} else {
		nbi.Logger.Debugf(ctx, "MAC address %s does not exist in Netbox. Creating it...", newMACAddress.MAC)
		newMACAddress, err := service.Create(ctx, nbi.NetboxAPI, newMACAddress)
		if err != nil {
			return nil, err
		}
		nbi.macAddressesIndex[objType][objName][ifaceName][newMACAddress.MAC] = newMACAddress
		return newMACAddress, nil
	}
	return nbi.macAddressesIndex[objType][objName][ifaceName][newMACAddress.MAC], nil
}

// AddPrefix adds a new prefix to the Netbox inventory.
// It takes a context and a newPrefix object as input and
// returns the created or updated prefix object and an error, if any.
// If the prefix already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the prefix does not exist, it creates a new one.
func (nbi *NetboxInventory) AddPrefix(
	ctx context.Context,
	newPrefix *objects.Prefix,
) (*objects.Prefix, error) {
	newPrefix.NetboxObject.AddTag(nbi.SsotTag)
	newPrefix.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	if newPrefix.NetboxObject.CustomFields == nil {
		newPrefix.NetboxObject.CustomFields = make(map[string]interface{})
	}
	//nolint:forcetypeassert
	newPrefix.NetboxObject.CustomFields[constants.CustomFieldSourceName] = ctx.Value(constants.CtxSourceKey).(string)

	// Determine VRF ID for index key (0 = global table)
	vrfID := 0
	if newPrefix.VRF != nil {
		vrfID = newPrefix.VRF.ID
	}

	nbi.prefixesLock.Lock()
	defer nbi.prefixesLock.Unlock()

	if nbi.prefixesIndexByPrefix[newPrefix.Prefix] == nil {
		nbi.prefixesIndexByPrefix[newPrefix.Prefix] = make(map[int]*objects.Prefix)
	}

	if _, ok := nbi.prefixesIndexByPrefix[newPrefix.Prefix][vrfID]; ok {
		oldPrefix := nbi.prefixesIndexByPrefix[newPrefix.Prefix][vrfID]
		nbi.OrphanManager.RemoveItem(oldPrefix)
		diffMap, err := utils.JSONDiffMapExceptID(newPrefix, oldPrefix, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"Prefix %s already exists in Netbox but is out of date. Patching it...",
				newPrefix.Prefix,
			)
			patchedPrefix, err := service.Patch[objects.Prefix](
				ctx,
				nbi.NetboxAPI,
				oldPrefix.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.prefixesIndexByPrefix[newPrefix.Prefix][vrfID] = patchedPrefix
		} else {
			nbi.Logger.Debugf(ctx, "Prefix %s already exists in Netbox and is up to date...", newPrefix.Prefix)
		}
	} else {
		nbi.Logger.Debugf(ctx, "Prefix %s does not exist in Netbox. Creating it...", newPrefix.Prefix)
		newPrefix, err := service.Create(ctx, nbi.NetboxAPI, newPrefix)
		if err != nil {
			return nil, err
		}
		nbi.prefixesIndexByPrefix[newPrefix.Prefix][vrfID] = newPrefix
		return newPrefix, nil
	}
	return nbi.prefixesIndexByPrefix[newPrefix.Prefix][vrfID], nil
}

// AddWirelessLAN adds a new wireless LAN to the Netbox inventory.
// It takes a context and a newWirelessLan object as input and
// returns the created or updated wireless LAN object and an error, if any.
// If the wireless LAN already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the wireless LAN does not exist, it creates a new one.
func (nbi *NetboxInventory) AddWirelessLAN(
	ctx context.Context,
	newWirelessLan *objects.WirelessLAN,
) (*objects.WirelessLAN, error) {
	newWirelessLan.NetboxObject.AddTag(nbi.SsotTag)
	newWirelessLan.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.wirelessLANsLock.Lock()
	defer nbi.wirelessLANsLock.Unlock()
	if _, ok := nbi.wirelessLANsIndexBySSID[newWirelessLan.SSID]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldWirelessLan := nbi.wirelessLANsIndexBySSID[newWirelessLan.SSID]
		nbi.OrphanManager.RemoveItem(oldWirelessLan)
		diffMap, err := utils.JSONDiffMapExceptID(
			newWirelessLan,
			oldWirelessLan,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"WirelessLAN %s already exists in Netbox but is out of date. Patching it...",
				newWirelessLan.SSID,
			)
			patchedWirelessLan, err := service.Patch[objects.WirelessLAN](
				ctx,
				nbi.NetboxAPI,
				oldWirelessLan.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.wirelessLANsIndexBySSID[newWirelessLan.SSID] = patchedWirelessLan
		} else {
			nbi.Logger.Debugf(ctx, "WirelessLAN %s already exists in Netbox and is up to date...", newWirelessLan.SSID)
		}
	} else {
		nbi.Logger.Debugf(ctx, "WirelessLAN %s does not exist in Netbox. Creating it...", newWirelessLan.SSID)
		newWirelessLan, err := service.Create(ctx, nbi.NetboxAPI, newWirelessLan)
		if err != nil {
			return nil, err
		}
		nbi.wirelessLANsIndexBySSID[newWirelessLan.SSID] = newWirelessLan
	}
	return nbi.wirelessLANsIndexBySSID[newWirelessLan.SSID], nil
}

// AddWirelessLANGroup adds a new wireless LAN group to the Netbox inventory.
// It takes a context and a newWirelessLANGroup object as input and
// returns the created or updated wireless LAN group object and an error, if any.
// If the wireless LAN group already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the wireless LAN group does not exist, it creates a new one.
func (nbi *NetboxInventory) AddWirelessLANGroup(
	ctx context.Context,
	newWirelessLANGroup *objects.WirelessLANGroup,
) (*objects.WirelessLANGroup, error) {
	newWirelessLANGroup.NetboxObject.AddTag(nbi.SsotTag)
	newWirelessLANGroup.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	nbi.wirelessLANGroupsLock.Lock()
	defer nbi.wirelessLANGroupsLock.Unlock()
	if _, ok := nbi.wirelessLANGroupsIndexByName[newWirelessLANGroup.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldWirelessLANGroup := nbi.wirelessLANGroupsIndexByName[newWirelessLANGroup.Name]
		nbi.OrphanManager.RemoveItem(oldWirelessLANGroup)
		diffMap, err := utils.JSONDiffMapExceptID(
			newWirelessLANGroup,
			oldWirelessLANGroup,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"WirelessLANGroup %s already exists in Netbox but is out of date. Patching it...",
				newWirelessLANGroup.Name,
			)
			patchedWirelessLANGroup, err := service.Patch[objects.WirelessLANGroup](
				ctx,
				nbi.NetboxAPI,
				oldWirelessLANGroup.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.wirelessLANGroupsIndexByName[newWirelessLANGroup.Name] = patchedWirelessLANGroup
		} else {
			nbi.Logger.Debugf(ctx, "WirelessLANGroup %s already exists in Netbox and is up to date...", newWirelessLANGroup.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "WirelessLANGroup %s does not exist in Netbox. Creating it...", newWirelessLANGroup.Name)
		newWirelessLANGroup, err := service.Create(ctx, nbi.NetboxAPI, newWirelessLANGroup)
		if err != nil {
			return nil, err
		}
		nbi.wirelessLANGroupsIndexByName[newWirelessLANGroup.Name] = newWirelessLANGroup
	}
	return nbi.wirelessLANGroupsIndexByName[newWirelessLANGroup.Name], nil
}

// AddVirtualDisk adds a new virtual disk to the Netbox inventory.
// It takes a context and a newVirtualDisk object as input and
// returns the created or updated virtual disk object and an error, if any.
// If the virtual disk already exists in Netbox, it checks if it is up to date and patches it if necessary.
// If the virtual disk does not exist, it creates a new one.
func (nbi *NetboxInventory) AddVirtualDisk(
	ctx context.Context,
	newVirtualDisk *objects.VirtualDisk,
) (*objects.VirtualDisk, error) {
	newVirtualDisk.NetboxObject.AddTag(nbi.SsotTag)
	addSourceNameCustomField(ctx, &newVirtualDisk.NetboxObject)
	newVirtualDisk.SetCustomField(constants.CustomFieldOrphanLastSeenName, nil)
	if len(newVirtualDisk.Name) > constants.MaxVirtualDiskNameLength {
		nbi.Logger.Debugf(
			nbi.Ctx,
			"VirtualDisk name %s is too long, truncating to %d characters",
			newVirtualDisk.Name,
			constants.MaxVirtualDiskNameLength,
		)
		newVirtualDisk.Name = newVirtualDisk.Name[:constants.MaxVirtualDiskNameLength]
	}
	nbi.virtualDisksLock.Lock()
	defer nbi.virtualDisksLock.Unlock()
	if _, ok := nbi.virtualDisksIndexByVMIDAndName[newVirtualDisk.VM.ID][newVirtualDisk.Name]; ok {
		oldVirtualDisk := nbi.virtualDisksIndexByVMIDAndName[newVirtualDisk.VM.ID][newVirtualDisk.Name]
		nbi.OrphanManager.RemoveItem(oldVirtualDisk)
		diffMap, err := utils.JSONDiffMapExceptID(
			newVirtualDisk,
			oldVirtualDisk,
			false,
			nbi.SourcePriority,
		)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debugf(
				ctx,
				"VirtualDisk %s already exists in Netbox but is out of date. Patching it...",
				newVirtualDisk.Name,
			)
			patchedVirtualDisk, err := service.Patch[objects.VirtualDisk](
				ctx,
				nbi.NetboxAPI,
				oldVirtualDisk.ID,
				diffMap,
			)
			if err != nil {
				return nil, err
			}
			nbi.virtualDisksIndexByVMIDAndName[newVirtualDisk.VM.ID][newVirtualDisk.Name] = patchedVirtualDisk
		} else {
			nbi.Logger.Debugf(ctx, "VirtualDisk %s already exists in Netbox and is up to date...", newVirtualDisk.Name)
		}
	} else {
		nbi.Logger.Debugf(ctx, "VirtualDisk %s does not exist in Netbox. Creating it...", newVirtualDisk.Name)
		newVirtualDisk, err := service.Create(ctx, nbi.NetboxAPI, newVirtualDisk)
		if err != nil {
			return nil, err
		}
		if nbi.virtualDisksIndexByVMIDAndName[newVirtualDisk.VM.ID] == nil {
			nbi.virtualDisksIndexByVMIDAndName[newVirtualDisk.VM.ID] = make(map[string]*objects.VirtualDisk)
		}
		nbi.virtualDisksIndexByVMIDAndName[newVirtualDisk.VM.ID][newVirtualDisk.Name] = newVirtualDisk
	}
	return nbi.virtualDisksIndexByVMIDAndName[newVirtualDisk.VM.ID][newVirtualDisk.Name], nil
}

// Helper function that adds source name to custom field of the netbox object.
func addSourceNameCustomField(ctx context.Context, netboxObject *objects.NetboxObject) {
	if netboxObject.CustomFields == nil {
		netboxObject.CustomFields = make(map[string]interface{})
	}
	netboxObject.CustomFields[constants.CustomFieldSourceName] = ctx.Value(constants.CtxSourceKey).(string) //nolint
}

// applyDeviceFieldLengthLimitations applies field length limitations
// to the device object.
func (nbi *NetboxInventory) applyDeviceFieldLengthLimitations(device *objects.Device) {
	if len(device.Name) > constants.MaxDeviceNameLength {
		nbi.Logger.Warningf(
			nbi.Ctx,
			"Device name %s is too long, truncating to %d characters",
			device.Name,
			constants.MaxDeviceNameLength,
		)
		device.Name = device.Name[:constants.MaxDeviceNameLength]
	}
	if len(device.SerialNumber) > constants.MaxSerialNumberLength {
		nbi.Logger.Warningf(
			nbi.Ctx,
			"Device serial %s is too long, truncating to %d characters",
			device.SerialNumber,
			constants.MaxSerialNumberLength,
		)
		device.SerialNumber = device.SerialNumber[:constants.MaxSerialNumberLength]
	}
	if len(device.AssetTag) > constants.MaxAssetTagLength {
		nbi.Logger.Warningf(
			nbi.Ctx,
			"Device asset tag %s is too long, truncating to %d characters",
			device.AssetTag,
			constants.MaxAssetTagLength,
		)
		device.AssetTag = device.AssetTag[:constants.MaxAssetTagLength]
	}
}
