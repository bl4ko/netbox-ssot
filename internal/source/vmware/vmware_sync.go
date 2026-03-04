package vmware

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"

	devices "github.com/src-doo/go-devicetype-library/pkg"
	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/inventory"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
	"github.com/src-doo/netbox-ssot/internal/source/common"
	"github.com/src-doo/netbox-ssot/internal/utils"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func (vc *VmwareSource) syncTags(nbi *inventory.NetboxInventory) error {
	objectNames2NBTags := make(map[string][]*objects.Tag)
	for objectName, tags := range vc.Object2Tags {
		for _, tag := range tags {
			var description string
			if tag.Description != "" {
				description = fmt.Sprintf("Tag synced from vmware:%s", tag.Description)
			} else {
				description = "Tag synced from vmware"
			}
			nbTag, err := nbi.AddTag(vc.Ctx, &objects.Tag{
				Name:        tag.Name,
				Slug:        utils.Slugify(tag.Name),
				Color:       constants.ColorGreen,
				Description: description,
			})
			if err != nil {
				return fmt.Errorf("add tag %+v: %s", tag, err)
			}
			objectNames2NBTags[objectName] = append(objectNames2NBTags[objectName], nbTag)
		}
	}
	vc.Object2NBTags = objectNames2NBTags
	return nil
}

func (vc *VmwareSource) syncNetworks(nbi *inventory.NetboxInventory) error {
	for dvpgID, dvpg := range vc.Networks.DistributedVirtualPortgroups {
		// TODO: currently we are syncing only vlans
		// Get vlanGroup from relations
		vlanSite, err := common.MatchVlanToSite(
			vc.Ctx,
			nbi,
			dvpg.Name,
			vc.SourceConfig.VlanSiteRelations,
		)
		if err != nil {
			return fmt.Errorf("match vlan to site: %s", err)
		}
		vlanGroup, err := common.MatchVlanToGroup(
			vc.Ctx,
			nbi,
			dvpg.Name,
			vlanSite,
			vc.SourceConfig.VlanGroupRelations,
			vc.SourceConfig.VlanGroupSiteRelations,
		)
		if err != nil {
			return fmt.Errorf("match vlan to group: %s", err)
		}
		// Get tenant from relations
		vlanTenant, err := common.MatchVlanToTenant(
			vc.Ctx,
			nbi,
			dvpg.Name,
			vc.SourceConfig.VlanTenantRelations,
		)
		if err != nil {
			return fmt.Errorf("vlanTenant: %s", err)
		}
		if len(dvpg.VlanIDs) == 1 && len(dvpg.VlanIDRanges) == 0 && dvpg.VlanIDs[0] != 0 {
			vlanName := dvpg.Name
			if !strings.HasPrefix(vlanName, vc.SourceConfig.VlanPrefix) {
				vlanName = fmt.Sprintf("%s%04d_%s", vc.SourceConfig.VlanPrefix, dvpg.VlanIDs[0], vlanName)
			}
			networkTags := vc.Object2NBTags[dvpgID]
			vlanStruct := &objects.Vlan{
				NetboxObject: objects.NetboxObject{
					Tags: append(vc.Config.GetSourceTags(), networkTags...),
					CustomFields: map[string]interface{}{
						constants.CustomFieldSourceIDName: dvpgID,
					},
				},
				Name:   vlanName,
				Group:  vlanGroup,
				Site:   vlanSite,
				Vid:    dvpg.VlanIDs[0],
				Status: &objects.VlanStatusActive,
				Tenant: vlanTenant,
			}
			_, err := nbi.AddVlan(vc.Ctx, vlanStruct)
			if err != nil {
				return fmt.Errorf("add vlan %+v: %s", vlanStruct, err)
			}
		}
	}
	return nil
}

func (vc *VmwareSource) syncDatacenters(nbi *inventory.NetboxInventory) error {
	for dcID, dc := range vc.DataCenters {
		netboxClusterGroupName := dc.Name
		if mappedClusterGroupName, ok := vc.SourceConfig.DatacenterClusterGroupRelations[netboxClusterGroupName]; ok {
			netboxClusterGroupName = mappedClusterGroupName
			vc.Logger.Debugf(
				vc.Ctx,
				"mapping datacenter name %s to cluster group name %s",
				dc.Name,
				mappedClusterGroupName,
			)
		}
		clusterGroupStruct := &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{
				Description: fmt.Sprintf("Datacenter from source %s", vc.SourceConfig.Hostname),
				Tags:        vc.Config.GetSourceTags(),
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceIDName: dcID,
				},
			},
			Name: netboxClusterGroupName,
			Slug: utils.Slugify(netboxClusterGroupName),
		}
		_, err := nbi.AddClusterGroup(vc.Ctx, clusterGroupStruct)
		if err != nil {
			return fmt.Errorf(
				"failed to add vmware datacenter %+v as Netbox ClusterGroup: %v",
				clusterGroupStruct,
				err,
			)
		}
	}
	return nil
}

func (vc *VmwareSource) syncClusters(nbi *inventory.NetboxInventory) error {
	clusterType, err := vc.createVmwareClusterType(nbi)
	if err != nil {
		return fmt.Errorf("failed to add vmware ClusterType: %v", err)
	}
	// Then sync vmware Clusters as NetBoxClusters
	for clusterID, cluster := range vc.Clusters {
		clusterName := cluster.Name
		clusterTags := vc.Object2NBTags[clusterID]

		var clusterGroup *objects.ClusterGroup
		datacenterID := vc.Cluster2Datacenter[clusterID]
		clusterGroupName := vc.DataCenters[datacenterID].Name
		if mappedName, ok := vc.SourceConfig.DatacenterClusterGroupRelations[clusterGroupName]; ok {
			clusterGroupName = mappedName
		}
		clusterGroup, _ = nbi.GetClusterGroup(clusterGroupName)

		var clusterScopeType constants.ContentType
		var clusterScopeID int
		clusterSite, err := common.MatchClusterToSite(
			vc.Ctx,
			nbi,
			clusterName,
			vc.SourceConfig.ClusterSiteRelations,
		)
		if err != nil {
			return fmt.Errorf("match cluster to site: %s", err)
		}
		if clusterSite != nil {
			clusterScopeType = constants.ContentTypeDcimSite
			clusterScopeID = clusterSite.ID
		}

		clusterTenant, err := common.MatchClusterToTenant(
			vc.Ctx,
			nbi,
			clusterName,
			vc.SourceConfig.ClusterTenantRelations,
		)
		if err != nil {
			return fmt.Errorf("match cluster to tenant: %s", err)
		}

		clusterStruct := &objects.Cluster{
			NetboxObject: objects.NetboxObject{
				Tags: append(vc.Config.GetSourceTags(), clusterTags...),
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceIDName: clusterID,
				},
			},
			Name:      clusterName,
			Type:      clusterType,
			Status:    objects.ClusterStatusActive,
			Group:     clusterGroup,
			ScopeType: clusterScopeType,
			ScopeID:   clusterScopeID,
			Tenant:    clusterTenant,
		}
		_, err = nbi.AddCluster(vc.Ctx, clusterStruct)
		if err != nil {
			return fmt.Errorf(
				"failed to add vmware cluster %+v as Netbox cluster: %v",
				clusterStruct,
				err,
			)
		}
	}
	return nil
}

// Host in vmware is a represented as device in netbox with a
// custom role Server.
func (vc *VmwareSource) syncHosts(nbi *inventory.NetboxInventory) error {
	for hostID, host := range vc.Hosts {
		var err error
		hostName := host.Name

		hostSite, err := common.MatchHostToSite(
			vc.Ctx,
			nbi,
			hostName,
			vc.SourceConfig.HostSiteRelations,
		)
		if err != nil {
			return fmt.Errorf("hostSite: %s", err)
		}

		hostTenant, err := common.MatchHostToTenant(
			vc.Ctx,
			nbi,
			hostName,
			vc.SourceConfig.HostTenantRelations,
		)
		if err != nil {
			return fmt.Errorf("hostTenant: %s", err)
		}

		hostCluster, _ := nbi.GetCluster(vc.Clusters[vc.Host2Cluster[hostID]].Name)
		if hostCluster == nil {
			// Create a hypothetical cluster https://github.com/src-doo/netbox-ssot/issues/141
			hostCluster, err = vc.createHypotheticalCluster(nbi, hostName, hostSite, hostTenant)
			if err != nil {
				return fmt.Errorf("add hypothetical cluster: %s", err)
			}
		}

		hostTags := vc.Object2NBTags[hostID]

		// Extract host hardware info
		var hostUUID, hostModel, hostManufacturerName string
		if host.Summary.Hardware != nil {
			hostUUID = host.Summary.Hardware.Uuid
			hostModel = host.Summary.Hardware.Model
			hostManufacturerName = host.Summary.Hardware.Vendor
			// Serialize manufacturer names so they match device type library
			if hostManufacturerName != "" {
				hostManufacturerName = utils.SerializeManufacturerName(hostManufacturerName)
			}
		}

		if hostModel == "" {
			hostModel = constants.DefaultModel
		}
		if hostManufacturerName == "" {
			hostManufacturerName = constants.DefaultManufacturer
		}

		// Enrich data from device type library if possible
		var deviceSlug string
		deviceData, hasDeviceData := devices.DeviceTypesMap[hostManufacturerName][hostModel]
		if hasDeviceData {
			deviceSlug = deviceData.Slug
		} else {
			deviceSlug = utils.GenerateDeviceTypeSlug(hostManufacturerName, hostModel)
		}

		manufacturerStruct := &objects.Manufacturer{
			Name: hostManufacturerName,
			Slug: utils.Slugify(hostManufacturerName),
		}
		hostManufacturer, err := nbi.AddManufacturer(vc.Ctx, manufacturerStruct)
		if err != nil {
			return fmt.Errorf(
				"failed adding vmware Manufacturer %v with error: %s",
				manufacturerStruct,
				err,
			)
		}

		// Create device type
		deviceTypeStruct := &objects.DeviceType{
			Manufacturer: hostManufacturer,
			Model:        hostModel,
			Slug:         deviceSlug,
		}
		hostDeviceType, err := nbi.AddDeviceType(vc.Ctx, deviceTypeStruct)
		if err != nil {
			return fmt.Errorf(
				"failed adding vmware DeviceType %+v with error: %s",
				deviceTypeStruct,
				err,
			)
		}

		// Find serial number from host summary.hardware.OtherIdentifyingInfo (vmware specific logic)
		var hostSerialNumber string
		serialInfoTypes := map[string]bool{
			"EnclosureSerialNumberTag": true,
			"ServiceTag":               true,
			"SerialNumberTag":          true,
		}
		var assetTag string
		for _, info := range host.Summary.Hardware.OtherIdentifyingInfo {
			infoType := info.IdentifierType.GetElementDescription().Key
			infoValue := strings.Trim(info.IdentifierValue, " ") // remove blank spaces from value
			if infoType == "AssetTag" {
				if infoValue == "No Asset Tag" {
					infoValue = ""
				}
				if !vc.SourceConfig.IgnoreAssetTags {
					assetTag = infoValue
				}
			} else if serialInfoTypes[infoType] {
				if info.IdentifierValue != "" {
					if !vc.SourceConfig.IgnoreSerialNumbers {
						hostSerialNumber = infoValue
					}
					break
				}
			}
		}

		var hostStatus *objects.DeviceStatus
		switch host.Summary.Runtime.ConnectionState {
		case "connected":
			hostStatus = &objects.DeviceStatusActive
		default:
			hostStatus = &objects.DeviceStatusOffline
		}

		var hostPlatform *objects.Platform
		osType := host.Summary.Config.Product.Name
		osVersion := host.Summary.Config.Product.Version
		platformName := utils.GeneratePlatformName(osType, osVersion, "")
		platformStruct := &objects.Platform{
			Name: platformName,
			Slug: utils.Slugify(platformName),
		}
		hostPlatform, err = nbi.AddPlatform(vc.Ctx, platformStruct)
		if err != nil {
			return fmt.Errorf(
				"failed adding vmware Platform %+v with error: %s",
				platformStruct,
				err,
			)
		}

		// Match host to a role. First test if user provided relations, if not
		// use default server role.
		var hostRole *objects.DeviceRole
		if len(vc.SourceConfig.HostRoleRelations) > 0 {
			hostRole, err = common.MatchHostToRole(
				vc.Ctx,
				nbi,
				hostName,
				vc.SourceConfig.HostRoleRelations,
			)
			if err != nil {
				return fmt.Errorf("match host to role: %s", err)
			}
		}
		if hostRole == nil {
			hostRole, err = nbi.AddServerDeviceRole(vc.Ctx)
			if err != nil {
				return fmt.Errorf("add server device role %s", err)
			}
		}

		hostCPUCores := host.Summary.Hardware.NumCpuCores
		hostMemGB := host.Summary.Hardware.MemorySize / constants.KiB / constants.KiB / constants.KiB

		hostStruct := &objects.Device{
			NetboxObject: objects.NetboxObject{
				Tags: append(vc.Config.GetSourceTags(), hostTags...),
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceIDName:     hostID,
					constants.CustomFieldDeviceUUIDName:   hostUUID,
					constants.CustomFieldHostCPUCoresName: fmt.Sprintf("%d", hostCPUCores),
					constants.CustomFieldHostMemoryName:   fmt.Sprintf("%d GB", hostMemGB),
				}},
			Name:         hostName,
			Status:       hostStatus,
			Platform:     hostPlatform,
			DeviceRole:   hostRole,
			Site:         hostSite,
			Tenant:       hostTenant,
			Cluster:      hostCluster,
			SerialNumber: hostSerialNumber,
			AssetTag:     assetTag,
			DeviceType:   hostDeviceType,
		}
		nbHost, err := nbi.AddDevice(vc.Ctx, hostStruct)
		if err != nil {
			return fmt.Errorf("failed to add vmware host %+v with error: %v", hostStruct, err)
		}

		// We also need to sync nics separately, because nic is a separate object in netbox
		err = vc.syncHostNics(nbi, host, nbHost, deviceData)
		if err != nil {
			return fmt.Errorf("failed to sync vmware host %s nics with error: %v", host.Name, err)
		}
	}
	return nil
}

func (vc *VmwareSource) syncHostNics(
	nbi *inventory.NetboxInventory,
	vcHost mo.HostSystem,
	nbHost *objects.Device,
	deviceData *devices.DeviceData,
) error {
	// Variable for storeing all ipAddresses from all host interfaces,
	// we use them to determine the primary ip of the host.
	hostIPv4Addresses := []*objects.IPAddress{}
	hostIPv6Addresses := []*objects.IPAddress{}

	// Sync host's physical interfaces
	err := vc.syncHostPhysicalNics(nbi, vcHost, nbHost, deviceData)
	if err != nil {
		return fmt.Errorf("physical interfaces sync: %s", err)
	}

	// Sync host's virtual interfaces
	err = vc.syncHostVirtualNics(nbi, vcHost, nbHost, hostIPv4Addresses, hostIPv6Addresses)
	if err != nil {
		return fmt.Errorf("virtual interfaces sync: %s", err)
	}

	// Set host's private ip address from collected ips
	err = vc.setHostPrimaryIPAddress(nbi, nbHost, hostIPv4Addresses, hostIPv6Addresses)
	if err != nil {
		return fmt.Errorf("adding host primary ip addresses: %s", err)
	}

	return nil
}

func (vc *VmwareSource) syncHostPhysicalNics(
	nbi *inventory.NetboxInventory,
	vcHost mo.HostSystem,
	nbHost *objects.Device,
	deviceData *devices.DeviceData,
) error {
	// Collect data over all physical interfaces
	if vcHost.Config != nil && vcHost.Config.Network != nil && vcHost.Config.Network.Pnic != nil {
		for _, pnic := range vcHost.Config.Network.Pnic {
			// Fetch host pnic data
			hostPnic, macAddress, err := vc.collectHostPhysicalNicData(nbi, nbHost, pnic, deviceData)
			if err != nil {
				return err
			}

			// Filter host pnic
			if utils.FilterInterfaceName(hostPnic.Name, vc.SourceConfig.InterfaceFilter) {
				vc.Logger.Debugf(
					vc.Ctx,
					"interface %s is filtered out with interfaceFilter %s",
					hostPnic.Name,
					vc.SourceConfig.InterfaceFilter,
				)
				continue
			}

			// After collecting all of the data add interface to nbi
			nbHostPnic, err := nbi.AddInterface(vc.Ctx, hostPnic)
			if err != nil {
				return fmt.Errorf("failed adding physical interface %+v: %s", hostPnic, err)
			}

			// Create MAC address
			if macAddress != "" {
				nbMACAddress, err := common.CreateMACAddressForObjectType(
					vc.Ctx,
					nbi,
					macAddress,
					nbHostPnic,
				)
				if err != nil {
					return fmt.Errorf("create mac address for object type: %s", err)
				}
				if err = common.SetPrimaryMACForInterface(vc.Ctx, nbi, nbHostPnic, nbMACAddress); err != nil {
					return fmt.Errorf("set primary mac for interface %+v: %s", nbHostPnic, err)
				}
			}
		}
	}
	return nil
}

//
//nolint:gocyclo
func (vc *VmwareSource) collectHostPhysicalNicData(
	nbi *inventory.NetboxInventory,
	nbHost *objects.Device,
	pnic types.PhysicalNic,
	_ *devices.DeviceData,
) (*objects.Interface, string, error) {
	pnicName := pnic.Device
	var pnicLinkSpeedMb int32
	if pnic.LinkSpeed != nil {
		pnicLinkSpeedMb = pnic.LinkSpeed.SpeedMb
	}
	if pnicLinkSpeedMb == 0 {
		if pnic.Spec.LinkSpeed != nil {
			pnicLinkSpeedMb = pnic.LinkSpeed.SpeedMb
		}
		if pnicLinkSpeedMb == 0 {
			if pnic.ValidLinkSpecification != nil {
				pnicLinkSpeedMb = pnic.ValidLinkSpecification[0].SpeedMb
			}
		}
	}
	var pnicDescription string
	if pnicLinkSpeedMb/constants.KB >= 1 {
		pnicDescription = fmt.Sprintf("%dGB/s", pnicLinkSpeedMb/constants.KB)
	} else {
		pnicDescription = fmt.Sprintf("%dMB/s", pnicLinkSpeedMb)
	}
	pnicDescription += " pNIC"
	var pnicMtu int
	var pnicMode *objects.InterfaceMode
	// Check virtual switches for data
	for vswitch, vswitchData := range vc.Networks.HostVirtualSwitches[nbHost.Name] {
		if slices.Contains(vswitchData.pnics, pnic.Key) {
			pnicDescription = fmt.Sprintf("%s (%s)", pnicDescription, vswitch)
			pnicMtu = vswitchData.mtu
		}
	}

	// Check proxy switches for data
	for _, pswitchData := range vc.Networks.HostProxySwitches[nbHost.Name] {
		if slices.Contains(pswitchData.pnics, pnic.Key) {
			pnicDescription = fmt.Sprintf("%s (%s)", pnicDescription, pswitchData.name)
			pnicMtu = pswitchData.mtu
			pnicMode = &objects.InterfaceModeTaggedAll
		}
	}

	// Check vlans on this pnic
	vlanIDMap := map[int]*objects.Vlan{} // set of vlans
	for portgroupName, portgroupData := range vc.Networks.HostPortgroups[nbHost.Name] {
		if slices.Contains(portgroupData.nics, pnicName) {
			if portgroupData.vlanID == 0 || portgroupData.vlanID > 4094 {
				vlanIDMap[portgroupData.vlanID] = &objects.Vlan{Vid: portgroupData.vlanID}
				continue
			}
			// Check if vlan with this vid already exists, else create it
			if vlanName, ok := vc.Networks.Vid2Name[portgroupData.vlanID]; ok {
				vlanSite, err := common.MatchVlanToSite(
					vc.Ctx,
					nbi,
					vlanName,
					vc.SourceConfig.VlanSiteRelations,
				)
				if err != nil {
					return nil, "", fmt.Errorf("match vlan to site: %s", err)
				}
				vlanGroup, err := common.MatchVlanToGroup(
					vc.Ctx,
					nbi,
					vlanName,
					vlanSite,
					vc.SourceConfig.VlanGroupRelations,
					vc.SourceConfig.VlanGroupSiteRelations,
				)
				if err != nil {
					return nil, "", fmt.Errorf("match vlan to group: %s", err)
				}
				vlan, vlanExists := nbi.GetVlan(vlanGroup.ID, portgroupData.vlanID)
				if vlanExists {
					vlanIDMap[portgroupData.vlanID] = vlan
				}
			} else {
				vlanName := portgroupName
				if !strings.HasPrefix(vlanName, vc.SourceConfig.VlanPrefix) {
					vlanName = fmt.Sprintf("%s%04d_%s", vc.SourceConfig.VlanPrefix, portgroupData.vlanID, vlanName)
				}
				vlanSite, err := common.MatchVlanToSite(vc.Ctx, nbi, vlanName, vc.SourceConfig.VlanSiteRelations)
				if err != nil {
					return nil, "", fmt.Errorf("match vlan to site: %s", err)
				}
				vlanGroup, err := common.MatchVlanToGroup(
					vc.Ctx,
					nbi,
					vlanName,
					vlanSite,
					vc.SourceConfig.VlanGroupRelations,
					vc.SourceConfig.VlanGroupSiteRelations,
				)
				if err != nil {
					return nil, "", fmt.Errorf("match vlan to group: %s", err)
				}
				vlanTenant, err := common.MatchVlanToTenant(vc.Ctx, nbi, vlanName, vc.SourceConfig.VlanTenantRelations)
				if err != nil {
					return nil, "", fmt.Errorf("match vlan to tenant: %s", err)
				}
				newVlan, newVlanExists := nbi.GetVlan(vlanGroup.ID, portgroupData.vlanID)
				if !newVlanExists {
					vlanStruct := &objects.Vlan{
						NetboxObject: objects.NetboxObject{
							Tags: vc.Config.GetSourceTags(),
						},
						Status: &objects.VlanStatusActive,
						Name:   vlanName,
						Site:   vlanSite,
						Vid:    portgroupData.vlanID,
						Tenant: vlanTenant,
						Group:  vlanGroup,
					}
					newVlan, err = nbi.AddVlan(vc.Ctx, vlanStruct)
					if err != nil {
						return nil, "", fmt.Errorf("add vlan %+v: %s", vlanStruct, err)
					}
				}
				vlanIDMap[portgroupData.vlanID] = newVlan
			}
		}
	}

	// Determine interface mode for non VM traffic NIC, from vlans data
	var taggedVlanList []*objects.Vlan // when mode="tagged"
	if len(vlanIDMap) > 0 {
		switch {
		case len(vlanIDMap) == 1 && vlanIDMap[0] != nil:
			pnicMode = &objects.InterfaceModeAccess
		case vlanIDMap[4095] != nil:
			pnicMode = &objects.InterfaceModeTaggedAll
		default:
			pnicMode = &objects.InterfaceModeTagged
			taggedVlanList = []*objects.Vlan{}
			if pnicMode == &objects.InterfaceModeTagged {
				for vid, vlan := range vlanIDMap {
					if vid == 0 {
						continue
					}
					taggedVlanList = append(taggedVlanList, vlan)
				}
			}
		}
	}

	pnicLinkSpeedKb := pnicLinkSpeedMb * constants.KB
	pnicType := objects.IfaceSpeed2IfaceType[objects.InterfaceSpeed(pnicLinkSpeedKb)]

	if pnicType == nil {
		pnicType = &objects.OtherInterfaceType
	}
	return &objects.Interface{
		NetboxObject: objects.NetboxObject{
			Tags:        vc.Config.GetSourceTags(),
			Description: pnicDescription,
			CustomFields: map[string]interface{}{
				constants.CustomFieldSourceName: vc.SourceConfig.Name,
			},
		},
		Device:      nbHost,
		Name:        pnicName,
		Status:      true,
		Type:        pnicType,
		Speed:       objects.InterfaceSpeed(pnicLinkSpeedKb),
		MTU:         pnicMtu,
		Mode:        pnicMode,
		TaggedVlans: taggedVlanList,
	}, strings.ToUpper(pnic.Mac), nil
}

func (vc *VmwareSource) syncHostVirtualNics(
	nbi *inventory.NetboxInventory,
	vcHost mo.HostSystem,
	nbHost *objects.Device,
	hostIPv4Addresses []*objects.IPAddress,
	hostIPv6Addresses []*objects.IPAddress,
) error {
	// Collect data over all virtual interfaces
	if vcHost.Config != nil && vcHost.Config.Network != nil && vcHost.Config.Network.Vnic != nil {
		for _, vnic := range vcHost.Config.Network.Vnic {
			// Fetch host vnic data
			hostVnic, macAddress, err := vc.collectHostVirtualNicData(nbi, nbHost, vcHost, vnic)
			if err != nil {
				return err
			}

			// Filter host vnic
			if utils.FilterInterfaceName(hostVnic.Name, vc.SourceConfig.InterfaceFilter) {
				vc.Logger.Debugf(
					vc.Ctx,
					"interface %s is filtered out with interfaceFilter %s",
					hostVnic.Name,
					vc.SourceConfig.InterfaceFilter,
				)
				continue
			}

			// After collecting all of the data add interface to nbi
			nbHostVnic, err := nbi.AddInterface(vc.Ctx, hostVnic)
			if err != nil {
				return fmt.Errorf("failed adding virtual interface %+v: %s", hostVnic, err)
			}

			// Create MAC address
			if macAddress != "" {
				nbMACAddress, err := common.CreateMACAddressForObjectType(
					vc.Ctx,
					nbi,
					macAddress,
					nbHostVnic,
				)
				if err != nil {
					return fmt.Errorf("create mac address for object type: %s", err)
				}
				if err = common.SetPrimaryMACForInterface(vc.Ctx, nbi, nbHostVnic, nbMACAddress); err != nil {
					return fmt.Errorf("set primary mac for interface %+v: %s", nbHostVnic, err)
				}
			}
			// Get IPv4 address for this vnic
			ipv4Address := vnic.Spec.Ip.IpAddress
			if ipv4Address == "" {
				continue
			}
			// VRF
			ipVRF, err := common.MatchIPToVRF(vc.Ctx, nbi, ipv4Address, vc.SourceConfig.IPVrfRelations)
			if err != nil {
				vc.Logger.Warningf(vc.Ctx, "match ip to vrf for %s: %s", ipv4Address, err)
			}
			if utils.IsPermittedIPAddress(
				ipv4Address,
				vc.SourceConfig.PermittedSubnets,
				vc.SourceConfig.IgnoredSubnets,
			) {
				ipv4MaskBits, err := utils.MaskToBits(vnic.Spec.Ip.SubnetMask)
				if err != nil {
					return fmt.Errorf("mask to bits: %s", err)
				}
				ipv4DNS := utils.ReverseLookup(ipv4Address)
				nbIPv4Address, err := nbi.AddIPAddress(vc.Ctx, &objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: vc.Config.GetSourceTags(),
						CustomFields: map[string]interface{}{
							constants.CustomFieldArpEntryName: false,
						},
					},
					Address:            fmt.Sprintf("%s/%d", ipv4Address, ipv4MaskBits),
					Status:             &objects.IPAddressStatusActive, // TODO
					DNSName:            ipv4DNS,
					Tenant:             nbHost.Tenant,
					AssignedObjectType: constants.ContentTypeDcimInterface,
					AssignedObjectID:   nbHostVnic.ID,
					VRF: ipVRF,
				})
				if err != nil {
					vc.Logger.Errorf(vc.Ctx, "add ipv4 address: %s", err)
					continue
				}
				hostIPv4Addresses = append(hostIPv4Addresses, nbIPv4Address)

				prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(nbIPv4Address.Address)
				if err != nil {
					vc.Logger.Warningf(vc.Ctx, "extract prefix from ip address: %s", err)
				} else if mask != constants.MaxIPv4MaskBits {
					_, err = nbi.AddPrefix(vc.Ctx, &objects.Prefix{
						Prefix: prefix,
						VRF:    ipVRF,
					})
					if err != nil {
						vc.Logger.Errorf(vc.Ctx, "add prefix: %s", err)
					}
				}
			}

			// Get IPv6 address for this vnic
			if vnic.Spec.Ip.IpV6Config != nil {
				for _, ipv6Entry := range vnic.Spec.Ip.IpV6Config.IpV6Address {
					ipv6Address := ipv6Entry.IpAddress
					ipv6Mask := ipv6Entry.PrefixLength
					if utils.IsPermittedIPAddress(
						ipv6Address,
						vc.SourceConfig.PermittedSubnets,
						vc.SourceConfig.IgnoredSubnets,
					) {
						nbIPv6Address, err := nbi.AddIPAddress(vc.Ctx, &objects.IPAddress{
							NetboxObject: objects.NetboxObject{
								Tags: vc.Config.GetSourceTags(),
								CustomFields: map[string]interface{}{
									constants.CustomFieldArpEntryName: false,
								},
							},
							Address:            fmt.Sprintf("%s/%d", ipv6Address, ipv6Mask),
							Status:             &objects.IPAddressStatusActive, // TODO
							Tenant:             nbHost.Tenant,
							AssignedObjectType: constants.ContentTypeDcimInterface,
							AssignedObjectID:   nbHostVnic.ID,
						})
						if err != nil {
							vc.Logger.Errorf(vc.Ctx, "add ipv6 address: %s", err)
							continue
						}
						hostIPv6Addresses = append(hostIPv6Addresses, nbIPv6Address)
					}
				}
			}
		}
	}
	return nil
}

func (vc *VmwareSource) setHostPrimaryIPAddress(
	nbi *inventory.NetboxInventory,
	nbHost *objects.Device,
	hostIPv4Addresses []*objects.IPAddress,
	hostIPv6Addresses []*objects.IPAddress,
) error {
	if len(hostIPv4Addresses) > 0 || len(hostIPv6Addresses) > 0 {
		var hostPrimaryIPv4 *objects.IPAddress
		for _, addr := range hostIPv4Addresses {
			if hostPrimaryIPv4 == nil || utils.Lookup(nbHost.Name) == addr.Address {
				hostPrimaryIPv4 = addr
			}
		}
		var hostPrimaryIPv6 *objects.IPAddress
		for _, addr := range hostIPv6Addresses {
			if hostPrimaryIPv6 == nil || utils.Lookup(nbHost.Name) == addr.Address {
				hostPrimaryIPv6 = addr
			}
		}
		newHost := *nbHost
		newHost.PrimaryIPv4 = hostPrimaryIPv4
		newHost.PrimaryIPv6 = hostPrimaryIPv6
		_, err := nbi.AddDevice(vc.Ctx, &newHost)
		if err != nil {
			return fmt.Errorf("updating host's primary ip: %s", err)
		}
	}

	return nil
}

func (vc *VmwareSource) collectHostVirtualNicData(
	nbi *inventory.NetboxInventory,
	nbHost *objects.Device,
	vcHost mo.HostSystem,
	vnic types.HostVirtualNic,
) (*objects.Interface, string, error) {
	vnicName := vnic.Device
	vnicPortgroupData, vnicPortgroupDataOk := vc.Networks.HostPortgroups[vcHost.Name][vnic.Portgroup]
	vnicDvPortgroupKey := ""
	if vnic.Spec.DistributedVirtualPort != nil {
		vnicDvPortgroupKey = vnic.Spec.DistributedVirtualPort.PortgroupKey
	}
	vnicDvPortgroupData, vnicDvPortgroupDataOk := vc.Networks.DistributedVirtualPortgroups[vnicDvPortgroupKey]
	vnicPortgroupVlanID := 0
	vnicDvPortgroupVlanIDs := []int{}
	var vnicMode *objects.InterfaceMode
	var vlanDescription, vnicDescription string

	// Get data from local portgroup, or distributed portgroup
	if vnicPortgroupDataOk {
		vnicPortgroupVlanID = vnicPortgroupData.vlanID
		vnicSwitch := vnicPortgroupData.vswitch
		vnicDescription = fmt.Sprintf(
			"%s (%s, vlan ID: %d)",
			vnic.Portgroup,
			vnicSwitch,
			vnicPortgroupVlanID,
		)
	} else if vnicDvPortgroupDataOk {
		vnicDescription = vnicDvPortgroupData.Name
		vnicDvPortgroupVlanIDs = vnicDvPortgroupData.VlanIDs
		if len(vnicDvPortgroupVlanIDs) == 1 && vnicDvPortgroupData.VlanIDs[0] == 4095 {
			vnicDescription = "all vlans"
			vnicMode = &objects.InterfaceModeTaggedAll
		} else {
			if len(vnicDvPortgroupData.VlanIDRanges) > 0 {
				vlanDescription = fmt.Sprintf("vlan IDs: %s", strings.Join(vnicDvPortgroupData.VlanIDRanges, ","))
			} else {
				vlanDescription = fmt.Sprintf("vlan ID: %d", vnicDvPortgroupData.VlanIDs[0])
			}
			if len(vnicDvPortgroupData.VlanIDs) == 1 && vnicDvPortgroupData.VlanIDs[0] == 0 {
				vnicMode = &objects.InterfaceModeAccess
			} else {
				vnicMode = &objects.InterfaceModeTagged
			}
		}
		vnicDvPortgroupDwSwitchUUID := vnic.Spec.DistributedVirtualPort.SwitchUuid
		vnicVswitch, vnicVswitchOk := vc.Networks.HostVirtualSwitches[vcHost.Name][vnicDvPortgroupDwSwitchUUID]
		if vnicVswitchOk {
			vnicDescription = fmt.Sprintf("%s (%v, %s)", vnicDescription, vnicVswitch, vlanDescription)
		}
	}

	var vnicUntaggedVlan *objects.Vlan
	var vnicTaggedVlans []*objects.Vlan
	if vnicPortgroupData != nil && vnicPortgroupVlanID != 0 {
		vnicUntaggedVlanSite, err := common.MatchVlanToSite(
			vc.Ctx,
			nbi,
			vc.Networks.Vid2Name[vnicPortgroupVlanID],
			vc.SourceConfig.VlanSiteRelations,
		)
		if err != nil {
			return nil, "", fmt.Errorf("vlan site: %s", err)
		}
		vnicUntaggedVlanGroup, err := common.MatchVlanToGroup(
			vc.Ctx,
			nbi,
			vc.Networks.Vid2Name[vnicPortgroupVlanID],
			vnicUntaggedVlanSite,
			vc.SourceConfig.VlanGroupRelations,
			vc.SourceConfig.VlanGroupSiteRelations,
		)
		if err != nil {
			return nil, "", fmt.Errorf("vlan group: %s", err)
		}
		vnicUntaggedVlan, _ = nbi.GetVlan(vnicUntaggedVlanGroup.ID, vnicPortgroupVlanID)
		vnicMode = &objects.InterfaceModeAccess
		// vnicUntaggedVlan = &objects.Vlan{
		// 	Name:   fmt.Sprintf("ESXi %s (ID: %d) (%s)", vnic.Portgroup, vnicPortgroupVlanId, nbHost.Site.Name),
		// 	Vid:    vnicPortgroupVlanId,
		// 	Tenant: nbHost.Tenant,
		// }
	} else if vnicDvPortgroupData != nil {
		for _, vnicDvPortgroupDataVlanID := range vnicDvPortgroupVlanIDs {
			if vnicMode != &objects.InterfaceModeTagged {
				break
			}
			if vnicDvPortgroupDataVlanID == 0 {
				continue
			}
			vnicTaggedVlanSite, err := common.MatchVlanToSite(
				vc.Ctx,
				nbi,
				vc.Networks.Vid2Name[vnicDvPortgroupDataVlanID],
				vc.SourceConfig.VlanSiteRelations,
			)
			if err != nil {
				return nil, "", fmt.Errorf("match vlan to site: %s", err)
			}
			vnicTaggedVlanGroup, err := common.MatchVlanToGroup(
				vc.Ctx,
				nbi,
				vc.Networks.Vid2Name[vnicDvPortgroupDataVlanID],
				vnicTaggedVlanSite,
				vc.SourceConfig.VlanGroupRelations,
				vc.SourceConfig.VlanGroupSiteRelations,
			)
			if err != nil {
				return nil, "", fmt.Errorf("match vlan to vlan group: %s", err)
			}
			taggedVlan, taggedVlanExists := nbi.GetVlan(vnicTaggedVlanGroup.ID, vnicDvPortgroupDataVlanID)
			if taggedVlanExists {
				vnicTaggedVlans = append(vnicTaggedVlans, taggedVlan)
			}
			// vnicTaggedVlans = append(vnicTaggedVlans, &objects.Vlan{
			// 	Name:   fmt.Sprintf("%s-%d", vnicDvPortgroupData.Name, vnicDvPortgroupDataVlanId),
			// 	Vid:    vnicDvPortgroupDataVlanId,
			// 	Tenant: nbHost.Tenant,
			// })
		}
	}
	return &objects.Interface{
		NetboxObject: objects.NetboxObject{
			Tags:        vc.Config.GetSourceTags(),
			Description: vnicDescription,
		},
		Device:       nbHost,
		Name:         vnicName,
		Status:       true,
		Type:         &objects.VirtualInterfaceType,
		MTU:          int(vnic.Spec.Mtu),
		Mode:         vnicMode,
		TaggedVlans:  vnicTaggedVlans,
		UntaggedVlan: vnicUntaggedVlan,
	}, strings.ToUpper(vnic.Spec.Mac), nil
}

// syncVMs syncs VMs from the source to Netbox.
func (vc *VmwareSource) syncVMs(nbi *inventory.NetboxInventory) error {
	const maxGoroutines = 50 // Maximum number of goroutines to run concurrently
	// Use a guard channel as semaphore to limit the number of goroutines
	guard := make(chan struct{}, maxGoroutines)
	// Use errChan to collect errors from goroutines
	errChan := make(chan error, len(vc.Vms))
	// Use a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Iterate over each VM and start a goroutine to sync it
	for vmKey, vm := range vc.Vms {
		guard <- struct{}{} // Block if max goroutines are running
		wg.Add(1)

		go func(vmKey string, vm mo.VirtualMachine) {
			defer wg.Done()
			defer func() { <-guard }() // Release one spot in the semaphore

			err := vc.syncVM(nbi, vmKey, vm)
			if err != nil {
				errChan <- err
			}
		}(vmKey, vm)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Collect any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// syncVM synces VM from the source to Netbox.
//
//nolint:gocyclo
func (vc *VmwareSource) syncVM(
	nbi *inventory.NetboxInventory,
	vmKey string,
	vm mo.VirtualMachine,
) error {
	isTemplate := false
	if vm.Config != nil && vm.Config.Template {
		isTemplate = true
	}

	if vc.SourceConfig.IgnoreVMTemplates && isTemplate {
		return nil
	}

	vmName := vm.Name
	vmHostName := vc.Hosts[vc.VM2Host[vmKey]].Name

	// Map to a vm role
	var vmRole *objects.DeviceRole
	var err error
	if len(vc.SourceConfig.VMRoleRelations) > 0 {
		vmRole, err = common.MatchVMToRole(vc.Ctx, nbi, vmHostName, vc.SourceConfig.VMRoleRelations)
		if err != nil {
			return fmt.Errorf("match vm to role: %s", err)
		}
	}
	if vmRole == nil {
		if isTemplate {
			vmRole, err = nbi.AddVMTemplateDeviceRole(vc.Ctx)
			if err != nil {
				return fmt.Errorf("add template device role: %s", err)
			}
		} else {
			vmRole, err = nbi.AddVMDeviceRole(vc.Ctx)
			if err != nil {
				return fmt.Errorf("get vm device role: %s", err)
			}
		}
	}

	// Tenant is received from VmTenantRelations
	vmTenant, err := common.MatchVMToTenant(vc.Ctx, nbi, vmName, vc.SourceConfig.VMTenantRelations)
	if err != nil {
		return fmt.Errorf("vm's Tenant: %s", err)
	}

	// Site is the same as the Host
	vmSite, err := common.MatchHostToSite(
		vc.Ctx,
		nbi,
		vmHostName,
		vc.SourceConfig.HostSiteRelations,
	)
	if err != nil {
		return fmt.Errorf("vm's Site: %s", err)
	}
	vmHost, _ := nbi.GetDevice(vmHostName, vmSite.ID)
	if vmHost == nil {
		return fmt.Errorf("host device %q not found in site %d, skipping VM", vmHostName, vmSite.ID)
	}

	// Cluster of the vm is same as the host
	vmCluster := vmHost.Cluster

	// VM status
	vmStatus := &objects.VMStatusOffline
	if vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
		vmStatus = &objects.VMStatusActive
	}

	// vmVCPUs and vmMemory
	var vmVCPUs, vmMemoryMB int32
	if vm.Config != nil {
		vmVCPUs = vm.Config.Hardware.NumCPU
		vmMemoryMB = vm.Config.Hardware.MemoryMB
	}

	// DisksSize
	// vmTotalDiskSizeMiB := int64(0)
	vmDisks := make([]*objects.VirtualDisk, 0)
	for _, hwDevice := range vm.Config.Hardware.Device {
		if disk, ok := hwDevice.(*types.VirtualDisk); ok {
			vmDiskSizeMiB := disk.CapacityInBytes / constants.MiB
			vmDiskName := disk.DeviceInfo.GetDescription().Label
			vmDiskSummary := disk.DeviceInfo.GetDescription().Summary
			vmDisks = append(vmDisks, &objects.VirtualDisk{
				NetboxObject: objects.NetboxObject{
					Description: vmDiskSummary,
				},
				Name: vmDiskName,
				Size: int(vmDiskSizeMiB),
			})
			// vmTotalDiskSizeMiB += vmDiskSizeMiB
		}
	}

	// Determine guest OS using fallback mechanisms
	var platformName string
	switch {
	case vm.Summary.Guest != nil && vm.Summary.Guest.GuestFullName != "":
		platformName = vm.Summary.Guest.GuestFullName
	case vm.Config.GuestFullName != "":
		platformName = vm.Config.GuestFullName
	case vm.Guest.GuestFullName != "":
		platformName = vm.Guest.GuestFullName
	}

	platformStruct := &objects.Platform{
		Name: platformName,
		Slug: utils.Slugify(platformName),
	}
	vmPlatform, err := nbi.AddPlatform(vc.Ctx, platformStruct)
	if err != nil {
		return fmt.Errorf(
			"failed adding vmware vm's Platform %+v with error: %s",
			platformStruct,
			err,
		)
	}

	// Extract additional info from CustomFields
	var vmOwners []string
	var vmOwnerEmails []string
	var vmDescription string
	vmCustomFields := map[string]interface{}{}
	if len(vm.Summary.CustomValue) > 0 {
		for _, field := range vm.Summary.CustomValue {
			if field, ok := field.(*types.CustomFieldStringValue); ok {
				fieldName := vc.CustomFieldID2Name[field.Key]

				if mappedField, ok := vc.SourceConfig.CustomFieldMappings[fieldName]; ok {
					switch mappedField {
					case "owner":
						vmOwners = utils.SerializeOwners(strings.Split(field.Value, ","))
					case "email":
						vmOwnerEmails = utils.SerializeEmails(strings.Split(field.Value, ","))
					case "description":
						vmDescription = strings.TrimSpace(field.Value)
					}
				} else {
					fieldName = utils.Alphanumeric(fieldName)
					if _, ok := nbi.GetCustomField(fieldName); !ok {
						customFieldStruct := &objects.CustomField{
							Name:                  fieldName,
							Type:                  objects.CustomFieldTypeText,
							CustomFieldUIVisible:  &objects.CustomFieldUIVisibleIfSet,
							CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
							ObjectTypes:           []constants.ContentType{constants.ContentTypeVirtualizationVirtualMachine},
						}
						_, err := nbi.AddCustomField(vc.Ctx, customFieldStruct)
						if err != nil {
							return fmt.Errorf("vm's custom field %s: %s", fieldName, err)
						}
					}
					vmCustomFields[fieldName] = field.Value
				}
			}
		}
	}
	vmCustomFields[constants.CustomFieldSourceName] = vc.SourceConfig.Name

	// netbox description has constraint <= len(200 characters)
	// In this case we make a comment
	var vmComments string
	if len(vmDescription) >= objects.MaxDescriptionLength {
		vmDescription = "See comments."
		vmComments = vmDescription
	}

	vmStruct := &objects.VM{
		NetboxObject: objects.NetboxObject{
			Tags:         append(vc.Config.GetSourceTags(), vc.Object2NBTags[vmKey]...),
			Description:  vmDescription,
			CustomFields: vmCustomFields,
		},
		Name:     vmName,
		Cluster:  vmCluster,
		Site:     vmSite,
		Tenant:   vmTenant,
		Status:   vmStatus,
		Host:     vmHost,
		Platform: vmPlatform,
		VCPUs:    float32(vmVCPUs),
		Memory:   int(vmMemoryMB),
		// Disk:     int(vmTotalDiskSizeMiB),
		Comments: vmComments,
		Role:     vmRole,
	}
	newVM, err := nbi.AddVM(vc.Ctx, vmStruct)
	if err != nil {
		return fmt.Errorf("failed to sync vmware VM %s: %v", vmName, err)
	}

	// For non template VMS also sync contacts and their network interfaces
	if !isTemplate {
		err = vc.addVMContact(nbi, newVM, vmOwners, vmOwnerEmails)
		if err != nil {
			return fmt.Errorf("adding %s's contact: %s", newVM, err)
		}

		// Sync vm interfaces
		err = vc.syncVMInterfaces(nbi, vm, newVM)
		if err != nil {
			return fmt.Errorf("failed to sync vmware %s's interfaces: %v", newVM, err)
		}

		// Sync vm disks
		err := vc.syncVMDisks(nbi, newVM, vmDisks)
		if err != nil {
			return fmt.Errorf("failed to sync vm's %+v disks: %s", newVM, err)
		}
	}
	return nil
}

// syncVMDisks syncs VM's disks to Netbox.
func (vc *VmwareSource) syncVMDisks(
	nbi *inventory.NetboxInventory,
	vm *objects.VM,
	vmDisks []*objects.VirtualDisk,
) error {
	for _, disk := range vmDisks {
		disk.VM = vm
		_, err := nbi.AddVirtualDisk(vc.Ctx, disk)
		if err != nil {
			return fmt.Errorf("adding VirtualDisk %+v: %s", disk, err)
		}
	}
	return nil
}

// Syncs VM's interfaces to Netbox.
func (vc *VmwareSource) syncVMInterfaces(
	nbi *inventory.NetboxInventory,
	vmwareVM mo.VirtualMachine,
	netboxVM *objects.VM,
) error {
	// Data to determine the primary IP address of the vm
	var vmDefaultGatewayIpv4 string
	var vmDefaultGatewayIpv6 string
	vmIPv4Addresses := make([]*objects.IPAddress, 0)
	vmIPv6Addresses := make([]*objects.IPAddress, 0)

	// From vm's routing determine the default interface
	if len(vmwareVM.Guest.IpStack) > 0 {
		for _, route := range vmwareVM.Guest.IpStack[0].IpRouteConfig.IpRoute {
			if route.PrefixLength == 0 {
				ipAddress := route.Network
				if ipAddress == "" {
					continue
				}
				gatewayIPAddress := route.Gateway.IpAddress
				if gatewayIPAddress == "" {
					continue
				}

				// Get version from ipAddress (v4 or v6)
				ipVersion := utils.GetIPVersion(ipAddress)
				if ipVersion == constants.IPv4 {
					vmDefaultGatewayIpv4 = gatewayIPAddress
				} else if ipVersion == constants.IPv6 {
					vmDefaultGatewayIpv6 = gatewayIPAddress
				}
			}
		}
	}

	for _, vmDevice := range vmwareVM.Config.Hardware.Device {
		// TODO: Refactor this to avoid hardcoded typecasting. Ensure all types
		// that compose VirtualEthernetCard are properly handled.
		var vmEthernetCard *types.VirtualEthernetCard
		switch v := vmDevice.(type) {
		case *types.VirtualPCNet32:
			vmEthernetCard = &v.VirtualEthernetCard
		case *types.VirtualVmxnet3:
			vmEthernetCard = &v.VirtualEthernetCard
		case *types.VirtualVmxnet2:
			vmEthernetCard = &v.VirtualEthernetCard
		case *types.VirtualVmxnet:
			vmEthernetCard = &v.VirtualEthernetCard
		case *types.VirtualE1000e:
			vmEthernetCard = &v.VirtualEthernetCard
		case *types.VirtualE1000:
			vmEthernetCard = &v.VirtualEthernetCard
		case *types.VirtualSriovEthernetCard:
			vmEthernetCard = &v.VirtualEthernetCard
		case *types.VirtualEthernetCard:
			vmEthernetCard = v
		default:
			continue
		}

		if vmEthernetCard != nil {
			nicIPv4Addresses, nicIPv6Addresses, collectedVMIface, macAddress, err := vc.collectVMInterfaceData(
				nbi,
				netboxVM,
				vmwareVM,
				vmEthernetCard,
			)
			if err != nil {
				return err
			}

			// Apply filter to VMIface name
			if utils.FilterInterfaceName(collectedVMIface.Name, vc.SourceConfig.InterfaceFilter) {
				vc.Logger.Debugf(
					vc.Ctx,
					"interface %s is filtered out with interfaceFilter %s",
					collectedVMIface.Name,
					vc.SourceConfig.InterfaceFilter,
				)
				continue
			}

			nbVMInterface, err := nbi.AddVMInterface(vc.Ctx, collectedVMIface)
			if err != nil {
				return fmt.Errorf("adding VmInterface %+v: %s", collectedVMIface, err)
			}
			if macAddress != "" {
				nbMACAddress, err := common.CreateMACAddressForObjectType(
					vc.Ctx,
					nbi,
					macAddress,
					nbVMInterface,
				)
				if err != nil {
					return fmt.Errorf("creating MAC address for %+v: %s", collectedVMIface, err)
				}
				if err = common.SetPrimaryMACForInterface(vc.Ctx, nbi, nbVMInterface, nbMACAddress); err != nil {
					return fmt.Errorf("setting primary MAC for %+v: %s", collectedVMIface, err)
				}
			}

			vmIPv4Addresses, vmIPv6Addresses = vc.addVMInterfaceIPs(
				nbi,
				netboxVM,
				nbVMInterface,
				nicIPv4Addresses,
				nicIPv6Addresses,
				vmIPv4Addresses,
				vmIPv6Addresses,
			)
		}
	}
	vc.setVMPrimaryIPAddress(
		nbi,
		netboxVM,
		vmDefaultGatewayIpv4,
		vmDefaultGatewayIpv6,
		vmIPv4Addresses,
		vmIPv6Addresses,
	)
	return nil
}

func (vc *VmwareSource) collectVMInterfaceData(
	nbi *inventory.NetboxInventory,
	netboxVM *objects.VM,
	vmwareVM mo.VirtualMachine,
	vmEthernetCard *types.VirtualEthernetCard,
) ([]string, []string, *objects.VMInterface, string, error) {
	intMac := vmEthernetCard.MacAddress
	intConnected := vmEthernetCard.Connectable.Connected
	intDeviceBackingInfo := vmEthernetCard.Backing
	intDeviceInfo := vmEthernetCard.DeviceInfo
	nicIPv4Addresses := []string{}
	nicIPv6Addresses := []string{}
	var intMtu int
	var intNetworkName string
	var intNetworkPrivate bool
	var intMode *objects.VMInterfaceMode
	intNetworkVlanIDs := []int{}
	intNetworkVlanIDRanges := []string{}

	// Get info from local vSwitches if possible, else from DistributedPortGroup
	if backingInfo, ok := intDeviceBackingInfo.(*types.VirtualEthernetCardNetworkBackingInfo); ok {
		intNetworkName = backingInfo.DeviceName
		intHostPgroup := vc.Networks.HostPortgroups[netboxVM.Host.Name][intNetworkName]

		if intHostPgroup != nil {
			intNetworkVlanIDs = []int{intHostPgroup.vlanID}
			intNetworkVlanIDRanges = []string{strconv.Itoa(intHostPgroup.vlanID)}
			intVswitchName := intHostPgroup.vswitch
			intVswitchData := vc.Networks.HostVirtualSwitches[netboxVM.Host.Name][intVswitchName]
			if intVswitchData != nil {
				intMtu = intVswitchData.mtu
			}
		}
	} else if backingInfo, ok := intDeviceBackingInfo.(*types.VirtualEthernetCardDistributedVirtualPortBackingInfo); ok {
		dvsPortgroupKey := backingInfo.Port.PortgroupKey
		intPortgroupData := vc.Networks.DistributedVirtualPortgroups[dvsPortgroupKey]

		if intPortgroupData != nil {
			intNetworkName = intPortgroupData.Name
			intNetworkVlanIDs = intPortgroupData.VlanIDs
			intNetworkVlanIDRanges = intPortgroupData.VlanIDRanges
			if len(intNetworkVlanIDRanges) == 0 {
				intNetworkVlanIDRanges = []string{strconv.Itoa(intNetworkVlanIDs[0])}
			}
			intNetworkPrivate = intPortgroupData.Private
		}

		intDvswitchUUID := backingInfo.Port.SwitchUuid
		intDvswitchData := vc.Networks.HostProxySwitches[netboxVM.Host.Name][intDvswitchUUID]

		if intDvswitchData != nil {
			intMtu = intDvswitchData.mtu
		}
	}

	var vlanDescription string
	intLabel := intDeviceInfo.GetDescription().Label
	splitStr := strings.Split(intLabel, " ")
	intName := fmt.Sprintf("vNIC %s", splitStr[len(splitStr)-1])
	intFullName := intName
	if intNetworkName != "" {
		intFullName = fmt.Sprintf("%s (%s)", intFullName, intNetworkName)
	}
	intDescription := intLabel
	if len(intNetworkVlanIDs) > 0 {
		if len(intNetworkVlanIDs) == 1 && intNetworkVlanIDs[0] == 4095 {
			vlanDescription = "all vlans"
			intMode = &objects.VMInterfaceModeTaggedAll
		} else {
			vlanDescription = fmt.Sprintf("vlan ID: %s", strings.Join(intNetworkVlanIDRanges, ", "))
			if len(intNetworkVlanIDs) == 1 {
				intMode = &objects.VMInterfaceModeAccess
			} else {
				intMode = &objects.VMInterfaceModeTagged
			}
		}

		if intNetworkPrivate {
			vlanDescription += "(private)"
		}
		intDescription = fmt.Sprintf("%s (%s)", intDescription, vlanDescription)
	}
	// Find corresponding guest NIC and get IP addresses and connected status
	for _, guestNic := range vmwareVM.Guest.Net {
		if intMac != guestNic.MacAddress {
			continue
		}
		intConnected = guestNic.Connected

		if guestNic.IpConfig != nil {
			for _, intIP := range guestNic.IpConfig.IpAddress {
				intIPAddress := fmt.Sprintf("%s/%d", intIP.IpAddress, intIP.PrefixLength)
				ipVersion := utils.GetIPVersion(intIP.IpAddress)
				switch ipVersion {
				case constants.IPv4:
					nicIPv4Addresses = append(nicIPv4Addresses, intIPAddress)
				case constants.IPv6:
					nicIPv6Addresses = append(nicIPv6Addresses, intIPAddress)
				default:
					return nicIPv4Addresses, nicIPv6Addresses, nil, "", fmt.Errorf(
						"unknown ip version: %s",
						intIPAddress,
					)
				}
			}
		}
	}
	var intUntaggedVlan *objects.Vlan
	var intTaggedVlanList []*objects.Vlan
	if len(intNetworkVlanIDs) > 0 && intMode != &objects.VMInterfaceModeTaggedAll {
		if len(intNetworkVlanIDs) == 1 && intNetworkVlanIDs[0] != 0 {
			vidID := intNetworkVlanIDs[0]
			nicUntaggedVlanSite, err := common.MatchVlanToSite(
				vc.Ctx,
				nbi,
				vc.Networks.Vid2Name[vidID],
				vc.SourceConfig.VlanSiteRelations,
			)
			if err != nil {
				return nicIPv4Addresses, nicIPv6Addresses, nil, "", fmt.Errorf(
					"match vlan to site: %s",
					err,
				)
			}
			nicUntaggedVlanGroup, err := common.MatchVlanToGroup(
				vc.Ctx,
				nbi,
				vc.Networks.Vid2Name[vidID],
				nicUntaggedVlanSite,
				vc.SourceConfig.VlanGroupRelations,
				vc.SourceConfig.VlanGroupSiteRelations,
			)
			if err != nil {
				return nicIPv4Addresses, nicIPv6Addresses, nil, "", fmt.Errorf(
					"mathc vlan to vlan group: %s",
					err,
				)
			}
			intUntaggedVlan, _ = nbi.GetVlan(nicUntaggedVlanGroup.ID, vidID)
		} else {
			intTaggedVlanList = []*objects.Vlan{}
			for _, intNetworkVlanID := range intNetworkVlanIDs {
				if intNetworkVlanID == 0 {
					continue
				}
				// nicTaggedVlanList = append(nicTaggedVlanList, nbi.get[intNetworkVlanId])
			}
		}
	}
	return nicIPv4Addresses, nicIPv6Addresses, &objects.VMInterface{
		NetboxObject: objects.NetboxObject{
			Tags:        vc.Config.GetSourceTags(),
			Description: intDescription,
		},
		VM:           netboxVM,
		Name:         intFullName,
		MTU:          intMtu,
		Mode:         intMode,
		Enabled:      intConnected,
		TaggedVlans:  intTaggedVlanList,
		UntaggedVlan: intUntaggedVlan,
	}, strings.ToUpper(intMac), nil
}

// Function that adds all collected IPs for the vm's interface to netbox.
func (vc *VmwareSource) addVMInterfaceIPs(
	nbi *inventory.NetboxInventory,
	netboxVM *objects.VM,
	nbVMInterface *objects.VMInterface,
	nicIPv4Addresses []string,
	nicIPv6Addresses []string,
	vmIPv4Addresses []*objects.IPAddress,
	vmIPv6Addresses []*objects.IPAddress,
) ([]*objects.IPAddress, []*objects.IPAddress) {
	// Add all collected ipv4 addresses for the interface to netbox
	for _, ipv4Address := range nicIPv4Addresses {
		if utils.IsPermittedIPAddress(
			ipv4Address,
			vc.SourceConfig.PermittedSubnets,
			vc.SourceConfig.IgnoredSubnets,
		) {
			// VRF
			ipVRF, err := common.MatchIPToVRF(vc.Ctx, nbi, ipv4Address, vc.SourceConfig.IPVrfRelations)
            if err != nil {
                vc.Logger.Warningf(vc.Ctx, "match ip to vrf for %s: %s", ipv4Address, err)
            }

			ipAddressStruct := &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: vc.Config.GetSourceTags(),
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: false,
					},
				},
				Address:            ipv4Address,
				DNSName:            utils.ReverseLookup(ipv4Address),
				AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
				AssignedObjectID:   nbVMInterface.ID,
				Tenant:             netboxVM.Tenant,
				VRF:                ipVRF,
			}
			nbIPv4Address, err := nbi.AddIPAddress(vc.Ctx, ipAddressStruct)
			if err != nil {
				vc.Logger.Warningf(vc.Ctx, "adding ipv4 address %s: %s", ipAddressStruct, err)
				continue
			}
			vmIPv4Addresses = append(vmIPv4Addresses, nbIPv4Address)
			prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(nbIPv4Address.Address)
			if err != nil {
				vc.Logger.Warningf(vc.Ctx, "extract prefix from ip address: %s", err)
			} else if mask != constants.MaxIPv4MaskBits {
				prefixStruct := &objects.Prefix{
					Prefix: prefix,
					VRF:    ipVRF,
				}
				_, err = nbi.AddPrefix(vc.Ctx, prefixStruct)
				if err != nil {
					vc.Logger.Errorf(vc.Ctx, "add prefix %+v: %s", prefixStruct, err)
				}
			}
		}
	}

	// Add all collected ipv6 addresses for the interface to netbox
	for _, ipv6Address := range nicIPv6Addresses {
		if utils.IsPermittedIPAddress(
			ipv6Address,
			vc.SourceConfig.PermittedSubnets,
			vc.SourceConfig.IgnoredSubnets,
		) {
			// VRF
            ipVRF, err := common.MatchIPToVRF(vc.Ctx, nbi, ipv6Address, vc.SourceConfig.IPVrfRelations)
            if err != nil {
                vc.Logger.Warningf(vc.Ctx, "match ip to vrf for %s: %s", ipv6Address, err)
            }

			nbIPv6Address, err := nbi.AddIPAddress(vc.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: vc.Config.GetSourceTags(),
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: false,
					},
				},
				Address:            ipv6Address,
				DNSName:            utils.ReverseLookup(ipv6Address),
				AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
				AssignedObjectID:   nbVMInterface.ID,
				VRF:                ipVRF,
			})
			if err != nil {
				vc.Logger.Warningf(vc.Ctx, "adding ipv6 address: %s", err)
				continue
			}
			vmIPv6Addresses = append(vmIPv6Addresses, nbIPv6Address)
			prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(nbIPv6Address.Address)
			if err != nil {
				vc.Logger.Warningf(vc.Ctx, "extract prefix from ip address: %s", err)
			} else if mask != constants.MaxIPv6MaskBits {
				prefixStruct := &objects.Prefix{
					Prefix: prefix,
					VRF:    ipVRF,
				}
				_, err = nbi.AddPrefix(vc.Ctx, prefixStruct)
				if err != nil {
					vc.Logger.Errorf(vc.Ctx, "add prefix: %s", err)
				}
			}
		}
	}
	return vmIPv4Addresses, vmIPv6Addresses
}

// setVMPrimaryIPAddress updates the vm's primary IP in the following way:
// we loop through all of the collected IPv4 and IPv6 addresses for the vm.
// If any of the ips is in the same subnet as the default gateway, we choose it.
// If there is no ip in the subnet of the default gateway, we choose the first one.
func (vc *VmwareSource) setVMPrimaryIPAddress(
	nbi *inventory.NetboxInventory,
	netboxVM *objects.VM,
	vmDefaultGatewayIpv4 string,
	vmDefaultGatewayIpv6 string,
	vmIPv4Addresses []*objects.IPAddress,
	vmIPv6Addresses []*objects.IPAddress,
) {
	if len(vmIPv4Addresses) > 0 || len(vmIPv6Addresses) > 0 {
		var vmIPv4PrimaryAddress *objects.IPAddress
		for _, addr := range vmIPv4Addresses {
			if vmIPv4PrimaryAddress == nil ||
				utils.SubnetContainsIPAddress(vmDefaultGatewayIpv4, addr.Address) {
				vmIPv4PrimaryAddress = addr
			}
		}
		var vmIPv6PrimaryAddress *objects.IPAddress
		for _, addr := range vmIPv6Addresses {
			if vmIPv6PrimaryAddress == nil ||
				utils.SubnetContainsIPAddress(vmDefaultGatewayIpv6, addr.Address) {
				vmIPv6PrimaryAddress = addr
			}
		}
		newNetboxVM := *netboxVM
		newNetboxVM.PrimaryIPv4 = vmIPv4PrimaryAddress
		newNetboxVM.PrimaryIPv6 = vmIPv6PrimaryAddress
		_, err := nbi.AddVM(vc.Ctx, &newNetboxVM)
		if err != nil {
			vc.Logger.Warningf(vc.Ctx, "updating vm's primary ip: %s", err)
		}
	}
}

func (vc *VmwareSource) addVMContact(
	nbi *inventory.NetboxInventory,
	nbVM *objects.VM,
	vmOwners []string,
	vmOwnerEmails []string,
) error {
	// If vm owner name was found we also add contact assignment to the vm
	var vmMailMapFallback bool
	if len(vmOwners) > 0 && len(vmOwnerEmails) > 0 && len(vmOwners) != len(vmOwnerEmails) {
		vc.Logger.Debugf(
			vc.Ctx,
			"vm owner names and emails mismatch len(vmOwnerEmails) != len(vmOwners), using fallback mechanism",
		)
		vmMailMapFallback = true
	}
	vmOwner2Email := utils.MatchNamesWithEmails(vc.Ctx, vmOwners, vmOwnerEmails, vc.Logger)
	for i, vmOwnerName := range vmOwners {
		if vmOwnerName != "" {
			var vmOwnerEmail string
			if len(vmOwnerEmails) > 0 {
				if vmMailMapFallback {
					if match, ok := vmOwner2Email[vmOwnerName]; ok {
						vmOwnerEmail = match
					}
				} else {
					vmOwnerEmail = vmOwnerEmails[i]
				}
			}
			contactStruct := &objects.Contact{
				Name:  strings.TrimSpace(vmOwners[i]),
				Email: vmOwnerEmail,
			}
			contact, err := nbi.AddContact(
				vc.Ctx,
				contactStruct,
			)
			if err != nil {
				return fmt.Errorf("creating vm contact %+v: %s", contactStruct, err)
			}
			contactRole, _ := nbi.GetContactRole(objects.AdminContactRoleName)

			contactAssignmentSturct := &objects.ContactAssignment{
				ModelType: constants.ContentTypeVirtualizationVirtualMachine,
				ObjectID:  nbVM.ID,
				Contact:   contact,
				Role:      contactRole,
			}
			_, err = nbi.AddContactAssignment(vc.Ctx, contactAssignmentSturct)
			if err != nil {
				return fmt.Errorf("add contact assignment for vm: %s", err)
			}
		}
	}
	return nil
}

// createVmwareClusterType creates a new VMware cluster type in Netbox.
// It takes a NetboxInventory object as input and returns the created
// ClusterType object and an error, if any.
func (vc *VmwareSource) createVmwareClusterType(
	nbi *inventory.NetboxInventory,
) (*objects.ClusterType, error) {
	clusterType := &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: []*objects.Tag{vc.Config.SourceTypeTag},
		},
		Name: "VMware ESXi",
		Slug: utils.Slugify("VMware ESXi"),
	}
	clusterType, err := nbi.AddClusterType(vc.Ctx, clusterType)
	if err != nil {
		return nil, fmt.Errorf("failed to add vmware ClusterType %+v: %v", clusterType, err)
	}
	return clusterType, nil
}

// createHypotheticalCluster creates a cluster with name clusterName. This function is needed
// for all hosts that are not assigned to cluster so we can assign them to hypotheticalCluster.
// for more see: https://github.com/src-doo/netbox-ssot/issues/141
func (vc *VmwareSource) createHypotheticalCluster(
	nbi *inventory.NetboxInventory,
	hostName string,
	hostSite *objects.Site,
	hostTenant *objects.Tenant,
) (*objects.Cluster, error) {
	clusterType, err := vc.createVmwareClusterType(nbi)
	if err != nil {
		return nil, fmt.Errorf("failed to add vmware ClusterType: %v", err)
	}
	var clusterScopeType constants.ContentType
	var clusterScopeID int
	if hostSite != nil {
		clusterScopeType = constants.ContentTypeDcimSite
		clusterScopeID = hostSite.ID
	}
	clusterStruct := &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			Tags: vc.Config.GetSourceTags(),
		},
		Name:      hostName,
		Type:      clusterType,
		Status:    objects.ClusterStatusActive,
		ScopeType: clusterScopeType,
		ScopeID:   clusterScopeID,
		Tenant:    hostTenant,
	}
	nbCluster, err := nbi.AddCluster(vc.Ctx, clusterStruct)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to add vmware hypothetical cluster %+v: %v",
			clusterStruct,
			err,
		)
	}

	return nbCluster, nil
}
