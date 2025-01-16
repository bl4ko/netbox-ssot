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
	var vlanGroupName string
	if vlanSite != nil {
		vlanGroupName = fmt.Sprintf("%sDefaultVlanGroup", vlanSite.Name)
	} else {
		vlanGroupName = constants.DefaultVlanGroupName
	}
	vlanGroup, err := nbi.AddVlanGroup(ctx, &objects.VlanGroup{
		NetboxObject: objects.NetboxObject{
			Tags:        []*objects.Tag{nbi.SsotTag},
			Description: constants.DefaultVlanGroupDescription,
			CustomFields: map[string]interface{}{
				constants.CustomFieldSourceName: nbi.SsotTag.Name,
			},
		},
		Name:      vlanGroupName,
		Slug:      utils.Slugify(constants.DefaultVlanGroupName),
		VidRanges: []objects.VidRange{{constants.DefaultVID, constants.MaxVID}},
	})
	if err != nil {
		return nil, fmt.Errorf("add vlan group: %s", err)
	}
	return vlanGroup, nil
}
