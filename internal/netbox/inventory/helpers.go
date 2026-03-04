package inventory

import (
	"context"
	"fmt"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
	"github.com/src-doo/netbox-ssot/internal/utils"
)

// Inits default VlanGroup, which is required to group all Vlans that are not part of other
// vlangroups into it. Each vlan is indexed by their (vlanGroup, vid).
func (nbi *NetboxInventory) CreateDefaultVlanGroupForVlan(
	ctx context.Context,
	vlanSite *objects.Site,
) (*objects.VlanGroup, error) {
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
		defaultVlanGroup.Name = fmt.Sprintf("%s DefaultVlanGroup", vlanSite.Name)
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

// getIndexValuesForIPAddress returns the index values for the given IPAddress.
// Index values are the type, name and owner name of the assigned object of the IPAddress.
func (nbi *NetboxInventory) getIndexValuesForIPAddress(
	ipAddr *objects.IPAddress,
) (constants.ContentType, string, string, error) {
	var ipIfaceType constants.ContentType
	var ipIfaceName, ipIfaceParentName string
	if ipAddr.AssignedObjectType != "" {
		switch ipAddr.AssignedObjectType {
		case constants.ContentTypeDcimInterface:
			ipIfaceType = constants.ContentTypeDcimDevice
			ipIface := nbi.GetInterfaceByID(ipAddr.AssignedObjectID)
			if ipIface == nil {
				return "", "", "", nil // Skip — interface not in inventory
			}
			ipIfaceName = ipIface.Name
			if ipIface.Device != nil {
				ipIfaceParentName = ipIface.Device.Name
			}
		case constants.ContentTypeVirtualizationVMInterface:
			ipIfaceType = constants.ContentTypeVirtualizationVirtualMachine
			ipIface := nbi.GetVMInterfaceByID(ipAddr.AssignedObjectID)
			if ipIface == nil {
				return "", "", "", nil // Skip — interface not in inventory
			}
			ipIfaceName = ipIface.Name
			if ipIface.VM != nil {
				ipIfaceParentName = ipIface.VM.Name
			}
		default:
			return "", "", "", fmt.Errorf(
				"unsupported assigned object type for ip address %+v: %s",
				ipAddr,
				ipIfaceType,
			)
		}
	}
	return ipIfaceType, ipIfaceName, ipIfaceParentName, nil
}

func (nbi *NetboxInventory) getIndexValuesForMACAddress(
	macAddr *objects.MACAddress,
) (constants.ContentType, string, string, error) {
	var macIfaceType constants.ContentType
	var macIfaceName, macIfaceParentName string
	if macAddr.AssignedObjectType != "" {
		switch macAddr.AssignedObjectType {
		case constants.ContentTypeDcimInterface:
			macIfaceType = constants.ContentTypeDcimDevice
			macIface := nbi.GetInterfaceByID(macAddr.AssignedObjectID)
			if macIface == nil {
				return "", "", "", nil // Skip — interface not in inventory
			}
			macIfaceName = macIface.Name
			if macIface.Device != nil {
				macIfaceParentName = macIface.Device.Name
			}
		case constants.ContentTypeVirtualizationVMInterface:
			macIfaceType = constants.ContentTypeVirtualizationVirtualMachine
			macIface := nbi.GetVMInterfaceByID(macAddr.AssignedObjectID)
			if macIface == nil {
				return "", "", "", nil // Skip — interface not in inventory
			}
			macIfaceName = macIface.Name
			if macIface.VM != nil {
				macIfaceParentName = macIface.VM.Name
			}
		default:
			return "", "", "", fmt.Errorf(
				"unsupported assigned object type for mac address %+v: %s",
				macAddr,
				macIfaceType,
			)
		}
	}
	return macIfaceType, macIfaceName, macIfaceParentName, nil
}

func (nbi *NetboxInventory) verifyIPAddressIndexExists(
	ifaceType constants.ContentType,
	ifaceName string,
	ifaceParentName string,
) {
	nbi.ipAddressesLock.Lock()
	defer nbi.ipAddressesLock.Unlock()
	if nbi.ipAddressesIndex[ifaceType] == nil {
		nbi.ipAddressesIndex[ifaceType] = make(
			map[string]map[string]map[string]*objects.IPAddress,
		)
	}

	if nbi.ipAddressesIndex[ifaceType][ifaceName] == nil {
		nbi.ipAddressesIndex[ifaceType][ifaceName] = make(
			map[string]map[string]*objects.IPAddress,
		)
	}

	if nbi.ipAddressesIndex[ifaceType][ifaceName][ifaceParentName] == nil {
		nbi.ipAddressesIndex[ifaceType][ifaceName][ifaceParentName] = make(
			map[string]*objects.IPAddress,
		)
	}
}

func (nbi *NetboxInventory) verifyMACAddressIndexExists(
	ifaceType constants.ContentType,
	ifaceName string,
	ifaceParentName string,
) {
	nbi.macAddressesLock.Lock()
	defer nbi.macAddressesLock.Unlock()
	if nbi.macAddressesIndex[ifaceType] == nil {
		nbi.macAddressesIndex[ifaceType] = make(
			map[string]map[string]map[string]*objects.MACAddress,
		)
	}

	if nbi.macAddressesIndex[ifaceType][ifaceName] == nil {
		nbi.macAddressesIndex[ifaceType][ifaceName] = make(
			map[string]map[string]*objects.MACAddress,
		)
	}

	if nbi.macAddressesIndex[ifaceType][ifaceName][ifaceParentName] == nil {
		nbi.macAddressesIndex[ifaceType][ifaceName][ifaceParentName] = make(
			map[string]*objects.MACAddress,
		)
	}
}

// ipAddressIndexKey returns the index key for an IPAddress,
// incorporating the VRF ID to avoid collisions across VRFs.
func ipAddressIndexKey(ipAddress *objects.IPAddress) string {
	if ipAddress.VRF != nil {
		return fmt.Sprintf("vrf%d/%s", ipAddress.VRF.ID, ipAddress.Address)
	}
	return ipAddress.Address
}