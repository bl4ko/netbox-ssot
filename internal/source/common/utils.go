package common

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Function that matches vlanName to vlanGroupName using regexRelationsMap.
//
// In case there is no match or regexRelations is nil, it will return default VlanGroup.
func MatchVlanToGroup(nbi *inventory.NetBoxInventory, vlanName string, regexRelations map[string]string) (*objects.VlanGroup, error) {
	if regexRelations == nil {
		return nbi.VlanGroupsIndexByName[objects.DefaultVlanGroupName], nil
	}
	vlanGroupName, err := utils.MatchStringToValue(vlanName, regexRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vlan to group: %s", err)
	}
	if vlanGroupName != "" {
		vlanGroup, ok := nbi.VlanGroupsIndexByName[vlanGroupName]
		if !ok {
			return nil, fmt.Errorf("no vlan group exists with name: %s", vlanGroupName)
		}
		return vlanGroup, nil
	}

	return nbi.VlanGroupsIndexByName[objects.DefaultVlanGroupName], nil
}

// Function that matches Host from hostName to Site using hostSiteRelations.
//
// In case that there is not match or hostSiteRelations is nil, it will return nil.
func MatchHostToSite(nbi *inventory.NetBoxInventory, hostName string, hostSiteRelations map[string]string) (*objects.Site, error) {
	if hostSiteRelations == nil {
		return nil, nil
	}
	siteName, err := utils.MatchStringToValue(hostName, hostSiteRelations)
	if err != nil {
		return nil, fmt.Errorf("matching host to site: %s", err)
	}
	if siteName != "" {
		site, ok := nbi.SitesIndexByName[siteName]
		if !ok {
			return nil, fmt.Errorf("site with name %s doesn't exist", siteName)
		}
		return site, nil
	}
	return nil, nil
}

// Function that matches Host from hostName to Tenant using hostTenantRelations.
//
// In case that there is not match or hostTenantRelations is nil, it will return nil.
func MatchHostToTenant(nbi *inventory.NetBoxInventory, hostName string, hostTenantRelations map[string]string) (*objects.Tenant, error) {
	if hostTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(hostName, hostTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching host to tenant: %s", err)
	}
	if tenantName != "" {
		site, ok := nbi.TenantsIndexByName[tenantName]
		if !ok {
			return nil, fmt.Errorf("tenant with name %s doesn't exist", tenantName)
		}
		return site, nil
	}
	return nil, nil
}
