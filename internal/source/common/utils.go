package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/inventory"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
	"github.com/src-doo/netbox-ssot/internal/utils"
)

// Function that matches cluster to tenant using regexRelationsMap.
//
// In case there is no match or regexRelations is nil, it will return nil.
func MatchClusterToTenant(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	clusterName string,
	clusterTenantRelations map[string]string,
) (*objects.Tenant, error) {
	if clusterTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(clusterName, clusterTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching cluster to tenant: %s", err)
	}
	if tenantName != "" {
		tenant, ok := nbi.GetTenant(tenantName)
		if !ok {
			tenant, err := nbi.AddTenant(ctx, &objects.Tenant{
				Name: tenantName,
				Slug: utils.Slugify(tenantName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new tenant: %s", err)
			}
			return tenant, nil
		}
		return tenant, nil
	}
	return nil, nil
}

// Function that matches cluster to tenant using regexRelationsMap.
//
// In case there is no match or regexRelations is nil, it will return nil.
func MatchClusterToSite(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	clusterName string,
	clusterSiteRelations map[string]string,
) (*objects.Site, error) {
	if clusterSiteRelations == nil {
		return nil, nil
	}
	siteName, err := utils.MatchStringToValue(clusterName, clusterSiteRelations)
	if err != nil {
		return nil, fmt.Errorf("matching cluster to tenant: %s", err)
	}
	if siteName != "" {
		site, ok := nbi.GetSite(siteName)
		if !ok {
			newSite, err := nbi.AddSite(ctx, &objects.Site{
				Name: siteName,
				Slug: utils.Slugify(siteName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new site: %s", err)
			}
			return newSite, nil
		}
		return site, nil
	}
	return nil, nil
}

// Function that matches vlanName to vlanGroupName using regexRelationsMap.
//
// In case there is no match or regexRelations is nil, it will return default VlanGroup.
func MatchVlanToGroup(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	vlanName string,
	vlanSite *objects.Site,
	vlanGroupRelations map[string]string,
	vlanGroupSiteRelations map[string]string,
) (*objects.VlanGroup, error) {
	if vlanGroupRelations == nil {
		vlanGroup, err := nbi.CreateDefaultVlanGroupForVlan(ctx, vlanSite)
		if err != nil {
			return nil, fmt.Errorf("create default vlan group for vlan %s: %s", vlanName, err)
		}
		return vlanGroup, nil
	}
	vlanGroupName, err := utils.MatchStringToValue(vlanName, vlanGroupRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vlan to group: %s", err)
	}
	var vlanGroupSite *objects.Site
	if vlanGroupSiteRelations != nil {
		siteName, err := utils.MatchStringToValue(vlanName, vlanGroupSiteRelations)
		if err != nil {
			return nil, fmt.Errorf("matching vlan to site: %s", err)
		}
		if siteName != "" {
			vlanGroupSite, err = nbi.AddSite(ctx, &objects.Site{
				Name: siteName,
				Slug: utils.Slugify(siteName),
			})
			if err != nil {
				return nil, fmt.Errorf("add site: %s", err)
			}
		}
	}

	if vlanGroupName != "" {
		vlanGroup := &objects.VlanGroup{
			Name:      vlanGroupName,
			Slug:      utils.Slugify(vlanGroupName),
			VidRanges: []objects.VidRange{{constants.DefaultVID, constants.MaxVID}},
		}
		if vlanGroupSite != nil {
			vlanGroup.ScopeType = constants.ContentTypeDcimSite
			vlanGroup.ScopeID = vlanGroupSite.ID
		}
		vlanGroup, err := nbi.AddVlanGroup(ctx, vlanGroup)
		if err != nil {
			return nil, fmt.Errorf("add vlan group %+v: %s", vlanGroup, err)
		}
		return vlanGroup, nil
	}

	// No vlan group was matched create default one.
	vlanGroup, err := nbi.CreateDefaultVlanGroupForVlan(ctx, vlanSite)
	if err != nil {
		return nil, fmt.Errorf("create default vlan group for vlan %s: %s", vlanName, err)
	}
	return vlanGroup, nil
}

// Function that matches vlanName to tenant using vlanTenantRelations regex relations map.
//
// In case there is no match or vlanTenantRelations is nil, it will return nil.
func MatchVlanToTenant(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	vlanName string,
	vlanTenantRelations map[string]string,
) (*objects.Tenant, error) {
	if vlanTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(vlanName, vlanTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vlan to tenant: %s", err)
	}
	if tenantName != "" {
		tenant, ok := nbi.GetTenant(tenantName)
		if !ok {
			tenant, err := nbi.AddTenant(ctx, &objects.Tenant{
				Name: tenantName,
				Slug: utils.Slugify(tenantName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new tenant: %s", err)
			}
			return tenant, nil
		}
		return tenant, nil
	}

	return nil, nil
}

// MathcVlanToSite matches vlanName to Site using vlanSiteRelations.
//
// In case there is no match or vlanSiteRelations is nil, it returns nil.
func MatchVlanToSite(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	vlanName string,
	vlanSiteRelations map[string]string,
) (*objects.Site, error) {
	if vlanSiteRelations == nil {
		return nil, nil
	}
	siteName, err := utils.MatchStringToValue(vlanName, vlanSiteRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vlan to site: %s", err)
	}
	if siteName != "" {
		site, ok := nbi.GetSite(siteName)
		if !ok {
			newSite, err := nbi.AddSite(ctx, &objects.Site{
				Name: siteName,
				Slug: utils.Slugify(siteName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new site: %s", err)
			}
			return newSite, nil
		}
		return site, nil
	}
	return nil, nil
}

// Function that matches Host from hostName to Site using hostSiteRelations.
//
// In case that there is not match or hostSiteRelations is nil, it will return default site.
func MatchHostToSite(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	hostName string,
	hostSiteRelations map[string]string,
) (*objects.Site, error) {
	if hostSiteRelations == nil {
		return nil, nil
	}
	siteName, err := utils.MatchStringToValue(hostName, hostSiteRelations)
	if err != nil {
		return nil, fmt.Errorf("matching host to site: %s", err)
	}
	if siteName != "" {
		newSite, err := nbi.AddSite(ctx, &objects.Site{
			Name: siteName,
			Slug: utils.Slugify(siteName),
		})
		if err != nil {
			return nil, fmt.Errorf("add new site: %s", err)
		}
		return newSite, nil
	}
	site, _ := nbi.GetSite(constants.DefaultSite)
	return site, nil
}

// Function that matches Host from hostName to Tenant using hostTenantRelations.
//
// In case that there is not match or hostTenantRelations is nil, it will return nil.
func MatchHostToTenant(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	hostName string,
	hostTenantRelations map[string]string,
) (*objects.Tenant, error) {
	if hostTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(hostName, hostTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching host to tenant: %s", err)
	}
	if tenantName != "" {
		tenant, err := nbi.AddTenant(ctx, &objects.Tenant{
			Name: tenantName,
			Slug: utils.Slugify(tenantName),
		})
		if err != nil {
			return nil, fmt.Errorf("add new tenant: %s", err)
		}
		return tenant, nil
	}
	return nil, nil
}

// MatchHostToRole matches Host from hostName to DeviceRole using hostRoleRelations.
//
// In case that there is not match or hostRoleRelations is nil, it will return nil.
func MatchHostToRole(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	hostName string,
	hostRoleRelations map[string]string,
) (*objects.DeviceRole, error) {
	if hostRoleRelations == nil {
		return nil, nil
	}
	roleName, err := utils.MatchStringToValue(hostName, hostRoleRelations)
	if err != nil {
		return nil, fmt.Errorf("matching host to role: %s", err)
	}
	if roleName != "" {
		role, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
			Name: roleName,
			Slug: utils.Slugify(roleName),
		})
		if err != nil {
			return nil, fmt.Errorf("add new host role: %s", err)
		}
		return role, nil
	}
	return nil, nil
}

// Function that matches Vm from vmName to Tenant using vmTenantRelations.
//
// In case that there is not match or hostTenantRelations is nil, it will return nil.
func MatchVMToTenant(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	vmName string,
	vmTenantRelations map[string]string,
) (*objects.Tenant, error) {
	if vmTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(vmName, vmTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vm to tenant: %s", err)
	}
	if tenantName != "" {
		site, ok := nbi.GetTenant(tenantName)
		if !ok {
			tenant, err := nbi.AddTenant(ctx, &objects.Tenant{
				Name: tenantName,
				Slug: utils.Slugify(tenantName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new tenant: %s", err)
			}
			return tenant, nil
		}
		return site, nil
	}
	return nil, nil
}

// MatchVMToRole matches VM from vmName to DeviceRole using vmRoleRelations.
//
// In case that there is not match or hostRoleRelations is nil, it will return nil.
func MatchVMToRole(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	vmName string,
	vmRoleRelations map[string]string,
) (*objects.DeviceRole, error) {
	if vmRoleRelations == nil {
		return nil, nil
	}
	roleName, err := utils.MatchStringToValue(vmName, vmRoleRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vm to role: %s", err)
	}
	if roleName != "" {
		role, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
			Name: roleName,
			Slug: utils.Slugify(roleName),
		})
		if err != nil {
			return nil, fmt.Errorf("add new vm role: %s", err)
		}
		return role, nil
	}
	return nil, nil
}

// CreateMACAddressForObjectType creates MAC address for object type.
func CreateMACAddressForObjectType(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	mac string,
	targetInterface objects.MACAddressOwner,
) (*objects.MACAddress, error) {
	macAddress := &objects.MACAddress{
		MAC:                mac,
		AssignedObjectType: targetInterface.GetObjectType(),
		AssignedObjectID:   targetInterface.GetID(),
	}
	nbMACAddress, err := nbi.AddMACAddress(ctx, macAddress)
	if err != nil {
		return nil, fmt.Errorf(
			"add mac address %+v for interface %+v: %s",
			macAddress,
			targetInterface,
			err,
		)
	}
	return nbMACAddress, nil
}

func SetPrimaryIPAddressForObject(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	targetObject objects.IPAddressOwner,
	ipv4 *objects.IPAddress,
	ipv6 *objects.IPAddress,
) error {
	switch targetObject := targetObject.(type) {
	case *objects.Device:
		deviceCopy := *targetObject
		deviceCopy.PrimaryIPv4 = ipv4
		deviceCopy.PrimaryIPv6 = ipv6
		_, err := nbi.AddDevice(ctx, &deviceCopy)
		if err != nil {
			return fmt.Errorf("set primary ip for device %+v: %s", deviceCopy, err)
		}
	case *objects.VM:
		vmCopy := *targetObject
		vmCopy.PrimaryIPv4 = ipv4
		vmCopy.PrimaryIPv6 = ipv6
		_, err := nbi.AddVM(ctx, &vmCopy)
		if err != nil {
			return fmt.Errorf("set primary ip for vm %+v: %s", vmCopy, err)
		}
	}
	return nil
}

func SetPrimaryMACForInterface(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	targetInterface objects.MACAddressOwner,
	mac *objects.MACAddress,
) error {
	switch targetInterface := targetInterface.(type) {
	case *objects.Interface:
		interfaceCopy := *targetInterface
		interfaceCopy.PrimaryMACAddress = mac
		_, err := nbi.AddInterface(ctx, &interfaceCopy)
		if err != nil {
			return fmt.Errorf("set primary mac for interface %+v: %s", interfaceCopy, err)
		}
	case *objects.VMInterface:
		vmInterfaceCopy := *targetInterface
		vmInterfaceCopy.PrimaryMACAddress = mac
		_, err := nbi.AddVMInterface(ctx, &vmInterfaceCopy)
		if err != nil {
			return fmt.Errorf("set primary mac for interface %+v: %s", vmInterfaceCopy, err)
		}
	}
	return nil
}

// MatchIPToVRF matches an IP address to a VRF using ipVrfRelations regex map.
// The IP can be in CIDR format ("10.0.0.1/24") or plain ("10.0.0.1").
//
// In case there is no match or ipVrfRelations is nil, it returns nil (global routing table).
func MatchIPToVRF(
	ctx context.Context,
	nbi *inventory.NetboxInventory,
	ipAddress string,
	ipVrfRelations map[string]string,
) (*objects.VRF, error) {
	if ipVrfRelations == nil {
		return nil, nil
	}
	// Strip mask if present: "10.0.0.1/24" -> "10.0.0.1"
	ip := ipAddress
	if idx := strings.Index(ipAddress, "/"); idx != -1 {
		ip = ipAddress[:idx]
	}
	vrfName, err := utils.MatchStringToValue(ip, ipVrfRelations)
	if err != nil {
		return nil, fmt.Errorf("matching ip to vrf: %s", err)
	}
	if vrfName != "" {
		vrf, ok := nbi.GetVRF(vrfName)
		if !ok {
			return nil, fmt.Errorf(
				"VRF %q not found in NetBox: create it manually before syncing",
				vrfName,
			)
		}
		return vrf, nil
	}
	return nil, nil
}