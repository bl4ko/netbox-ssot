package inventory

import "github.com/bl4ko/netbox-ssot/internal/netbox/objects"

// GetVlan returns vlan for the given vlanGroupID and vlanID.
// Returns nil if vlan is not found.
// This function is thread-safe.
func (nbi *NetboxInventory) GetVlan(vlanGroupID int, vlanID int) *objects.Vlan {
	nbi.VlansLock.Lock()
	defer nbi.VlansLock.Unlock()
	if _, ok := nbi.VlansIndexByVlanGroupIDAndVID[vlanGroupID]; !ok {
		return nil
	}
	return nbi.VlansIndexByVlanGroupIDAndVID[vlanGroupID][vlanID]
}
