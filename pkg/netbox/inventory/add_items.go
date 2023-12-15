package inventory

import (
	"slices"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/common"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/dcim"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/extras"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/virtualization"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
)

// Function for adding a new tag to NetBoxInventory
func (ni *NetBoxInventory) AddTag(newTag *common.Tag) (*common.Tag, error) {
	existingTagIndex := slices.IndexFunc(ni.Tags, func(t *common.Tag) bool {
		return t.Name == newTag.Name
	})
	if existingTagIndex == -1 {
		ni.Logger.Debug("Tag ", newTag.Name, " does not exist in NetBox. Creating it...")
		createdTag, err := ni.NetboxApi.CreateTag(newTag)
		if err != nil {
			return nil, err
		}
		ni.Tags = append(ni.Tags, createdTag)
		return createdTag, nil
	} else {
		ni.Logger.Debug("Tag ", newTag.Name, " already exists in NetBox...")
		existingTag := ni.Tags[existingTagIndex]
		diffMap, err := utils.JsonDiffMapExceptId(newTag, existingTag)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			patchedTag, err := ni.NetboxApi.PatchTag(diffMap, existingTag.ID)
			if err != nil {
				return nil, err
			}
			ni.Tags[existingTagIndex] = patchedTag
			return patchedTag, nil
		} else {
			return existingTag, nil
		}
	}
}

// Adding new custom-field to inventory
func (ni *NetBoxInventory) AddCustomField(newCf *extras.CustomField) error {
	if _, ok := ni.CustomFieldsIndexByName[newCf.Name]; ok {
		existingCf := ni.CustomFieldsIndexByName[newCf.Name]
		diffMap, err := utils.JsonDiffMapExceptId(newCf, existingCf)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Custom field ", newCf.Name, " already exists in NetBox but is out of date. Patching it... ")
			patchedCf, err := ni.NetboxApi.PatchCustomField(diffMap, existingCf.ID)
			if err != nil {
				return err
			}
			ni.CustomFieldsIndexByName[newCf.Name] = patchedCf
		} else {
			ni.Logger.Debug("Custom field ", newCf.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Custom field ", newCf.Name, " does not exist in NetBox. Creating it...")
		newCf, err := ni.NetboxApi.CreateCustomField(newCf)
		if err != nil {
			return err
		}
		ni.CustomFieldsIndexByName[newCf.Name] = newCf
	}
	return nil
}

// Add Cluster to NetBoxInventory
func (ni *NetBoxInventory) AddClusterGroup(newCg *virtualization.ClusterGroup, newTags []*common.Tag) error {
	if _, ok := ni.ClusterGroupsIndexByName[newCg.Name]; ok {
		newCg.Tags = append([]*common.Tag{ni.SsotTag}, newTags...)
		diffMap, err := utils.JsonDiffMapExceptId(newCg, ni.ClusterGroupsIndexByName[newCg.Name])
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster group ", newCg.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedCg, err := ni.NetboxApi.PatchClusterGroup(diffMap, ni.ClusterGroupsIndexByName[newCg.Name].ID)
			if err != nil {
				return err
			}
			ni.ClusterGroupsIndexByName[newCg.Name] = patchedCg
		} else {
			ni.Logger.Debug("Cluster group ", newCg.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Cluster group ", newCg.Name, " does not exist in NetBox. Creating it...")
		newCg, err := ni.NetboxApi.CreateClusterGroup(newCg)
		if err != nil {
			return err
		}
		ni.ClusterGroupsIndexByName[newCg.Name] = newCg
	}
	return nil
}

// Add ClusterType to NetBoxInventory
func (ni *NetBoxInventory) AddClusterType(newClusterType *virtualization.ClusterType, newTags []*common.Tag) (*virtualization.ClusterType, error) {
	if _, ok := ni.ClusterTypesIndexByName[newClusterType.Name]; ok {
		newClusterType.Tags = append([]*common.Tag{ni.SsotTag}, newTags...)
		diffMap, err := utils.JsonDiffMapExceptId(newClusterType, ni.ClusterTypesIndexByName[newClusterType.Name])
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedClusterType, err := ni.NetboxApi.PatchClusterType(diffMap, ni.ClusterTypesIndexByName[newClusterType.Name].ID)
			if err != nil {
				return nil, err
			}
			ni.ClusterTypesIndexByName[newClusterType.Name] = patchedClusterType
			return patchedClusterType, nil
		} else {
			ni.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in NetBox and is up to date...")
			existingClusterType := ni.ClusterTypesIndexByName[newClusterType.Name]
			return existingClusterType, nil
		}
	} else {
		ni.Logger.Debug("Cluster type ", newClusterType.Name, " does not exist in NetBox. Creating it...")
		newClusterType, err := ni.NetboxApi.CreateClusterType(newClusterType)
		if err != nil {
			return nil, err
		}
		ni.ClusterTypesIndexByName[newClusterType.Name] = newClusterType
		return newClusterType, nil
	}
}

func (ni *NetBoxInventory) AddCluster(newCluster *virtualization.Cluster, newTags []*common.Tag) error {
	if _, ok := ni.ClustersIndexByName[newCluster.Name]; ok {
		newCluster.Tags = append([]*common.Tag{ni.SsotTag}, newTags...)
		diffMap, err := utils.JsonDiffMapExceptId(newCluster, ni.ClustersIndexByName[newCluster.Name])
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Cluster ", newCluster.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedCluster, err := ni.NetboxApi.PatchCluster(diffMap, ni.ClustersIndexByName[newCluster.Name].ID)
			if err != nil {
				return err
			}
			ni.ClustersIndexByName[newCluster.Name] = patchedCluster
		} else {
			ni.Logger.Debug("Cluster ", newCluster.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Cluster ", newCluster.Name, " does not exist in NetBox. Creating it...")
		newCluster, err := ni.NetboxApi.CreateCluster(newCluster)
		if err != nil {
			return err
		}
		ni.ClustersIndexByName[newCluster.Name] = newCluster
	}
	return nil
}

func (ni *NetBoxInventory) AddDeviceRole(newDeviceRole *dcim.DeviceRole, newTags []*common.Tag) error {
	if _, ok := ni.DeviceRolesIndexByName[newDeviceRole.Name]; ok {
		newDeviceRole.Tags = append([]*common.Tag{ni.SsotTag}, newTags...)
		diffMap, err := utils.JsonDiffMapExceptId(newDeviceRole, ni.DeviceRolesIndexByName[newDeviceRole.Name])
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			ni.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in NetBox but is out of date. Patching it...")
			patchedDeviceRole, err := ni.NetboxApi.PatchDeviceRole(diffMap, ni.DeviceRolesIndexByName[newDeviceRole.Name].ID)
			if err != nil {
				return err
			}
			ni.DeviceRolesIndexByName[newDeviceRole.Name] = patchedDeviceRole
		} else {
			ni.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in NetBox and is up to date...")
		}
	} else {
		ni.Logger.Debug("Device role ", newDeviceRole.Name, " does not exist in NetBox. Creating it...")
		newDeviceRole, err := ni.NetboxApi.CreateDeviceRole(newDeviceRole)
		if err != nil {
			return err
		}
		ni.DeviceRolesIndexByName[newDeviceRole.Name] = newDeviceRole
	}
	return nil
}
