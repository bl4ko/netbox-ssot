package inventory

import (
	"context"
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Inits default VlanGroup, which is required to group all Vlans that are not part of other
// vlangroups into it. Each vlan is indexed by their (vlanGroup, vid).
func (nbi *NetboxInventory) CreateDefaultVlanGroupForVlan(ctx context.Context, vlanSite *objects.Site) (*objects.VlanGroup, error) {
	defaultVlanGroup := &objects.VlanGroup{
		NetboxObject: objects.NetboxObject{
			Tags:        []*objects.Tag{nbi.SsotTag},
			Description: constants.DefaultVlanGroupDescription,
			CustomFields: map[string]interface{}{
				constants.CustomFieldSourceName: nbi.SsotTag.Name,
			},
		},
		VidRanges: []objects.VidRange{{constants.DefaultVID, constants.MaxVID}}}

	if vlanSite != nil {
		defaultVlanGroup.Name = fmt.Sprintf("%sDefaultVlanGroup", vlanSite.Name)
		defaultVlanGroup.ScopeType = constants.ContentTypeDcimSite
		defaultVlanGroup.ScopeID = vlanSite.ID
	} else {
		defaultVlanGroup.Name = constants.DefaultVlanGroupName
	}
	defaultVlanGroup.Slug = utils.Slugify(defaultVlanGroup.Name)

	nbVlanGroup, err := nbi.AddVlanGroup(ctx, defaultVlanGroup)

	if err != nil {
		return nil, fmt.Errorf("add vlan group %+v: %s", defaultVlanGroup, err)
	}
	return nbVlanGroup, nil
}
