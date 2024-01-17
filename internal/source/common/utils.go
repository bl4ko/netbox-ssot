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
