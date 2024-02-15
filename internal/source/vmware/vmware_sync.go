package vmware

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func (vc *VmwareSource) syncNetworks(nbi *inventory.NetboxInventory) error {
	for _, dvpg := range vc.Networks.DistributedVirtualPortgroups {
		// TODO: currently we are syncing only vlans
		// Get vlanGroup from relations
		vlanGroup, err := common.MatchVlanToGroup(nbi, dvpg.Name, vc.VlanGroupRelations)
		if err != nil {
			return fmt.Errorf("vlanGroup: %s", err)
		}
		// Get tenant from relations
		vlanTenant, err := common.MatchVlanToTenant(nbi, dvpg.Name, vc.VlanTenantRelations)
		if err != nil {
			return fmt.Errorf("vlanTenant: %s", err)
		}
		if len(dvpg.VlanIds) == 1 && len(dvpg.VlanIdRanges) == 0 {
			_, err := nbi.AddVlan(&objects.Vlan{
				NetboxObject: objects.NetboxObject{
					Tags: vc.CommonConfig.SourceTags,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: vc.SourceConfig.Name,
					},
				},
				Name:   dvpg.Name,
				Group:  vlanGroup,
				Vid:    dvpg.VlanIds[0],
				Status: &objects.VlanStatusActive,
				Tenant: vlanTenant,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vc *VmwareSource) syncDatacenters(nbi *inventory.NetboxInventory) error {
	for _, dc := range vc.DataCenters {

		nbClusterGroup := &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{
				Description: fmt.Sprintf("Datacenter from source %s", vc.SourceConfig.Hostname),
				Tags:        vc.CommonConfig.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: vc.SourceConfig.Name,
				},
			},
			Name: dc.Name,
			Slug: utils.Slugify(dc.Name),
		}
		_, err := nbi.AddClusterGroup(nbClusterGroup)
		if err != nil {
			return fmt.Errorf("failed to add vmware datacenter %s as Netbox ClusterGroup: %v", dc.Name, err)
		}
	}
	return nil
}

func (vc *VmwareSource) syncClusters(nbi *inventory.NetboxInventory) error {
	clusterType := &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: vc.CommonConfig.SourceTags,
			CustomFields: map[string]string{
				constants.CustomFieldSourceName: vc.SourceConfig.Name,
			},
		},
		Name: "Vmware ESXi",
		Slug: utils.Slugify("Vmware ESXi"),
	}
	clusterType, err := nbi.AddClusterType(clusterType)
	if err != nil {
		return fmt.Errorf("failed to add vmware ClusterType: %v", err)
	}
	// Then sync vmware Clusters as NetBoxClusters
	for clusterId, cluster := range vc.Clusters {

		clusterName := cluster.Name

		var clusterGroup *objects.ClusterGroup
		datacenterId := vc.Cluster2Datacenter[clusterId]
		clusterGroup = nbi.ClusterGroupsIndexByName[vc.DataCenters[datacenterId].Name]

		var clusterSite *objects.Site
		if vc.ClusterSiteRelations != nil {
			match, err := utils.MatchStringToValue(clusterName, vc.ClusterSiteRelations)
			if err != nil {
				return fmt.Errorf("failed to match vmware cluster %s to a Netbox site: %v", clusterName, err)
			}
			if match != "" {
				if _, ok := nbi.SitesIndexByName[match]; !ok {
					return fmt.Errorf("failed to match vmware cluster %s to a Netbox site: %v. Site with this name doesn't exist", clusterName, match)
				}
				clusterSite = nbi.SitesIndexByName[match]
			}
		}

		var clusterTenant *objects.Tenant
		if vc.ClusterTenantRelations != nil {
			match, err := utils.MatchStringToValue(clusterName, vc.ClusterTenantRelations)
			if err != nil {
				return fmt.Errorf("error occurred when matching vmware cluster %s to a Netbox tenant: %v", clusterName, err)
			}
			if match != "" {
				if _, ok := nbi.TenantsIndexByName[match]; !ok {
					return fmt.Errorf("failed to match vmware cluster %s to a Netbox tenant: %v. Tenant with this name doesn't exist", clusterName, match)
				}
				clusterTenant = nbi.TenantsIndexByName[match]
			}
		}

		nbCluster := &objects.Cluster{
			NetboxObject: objects.NetboxObject{
				Tags: vc.CommonConfig.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: vc.SourceConfig.Name,
				},
			},
			Name:   clusterName,
			Type:   clusterType,
			Status: objects.ClusterStatusActive,
			Group:  clusterGroup,
			Site:   clusterSite,
			Tenant: clusterTenant,
		}
		err := nbi.AddCluster(nbCluster)
		if err != nil {
			return fmt.Errorf("failed to add vmware cluster %s as Netbox cluster: %v", clusterName, err)
		}
	}
	return nil
}

// Host in vmware is a represented as device in netbox with a
// custom role Server
func (vc *VmwareSource) syncHosts(nbi *inventory.NetboxInventory) error {
	for hostId, host := range vc.Hosts {
		var err error
		hostName := host.Name
		hostCluster := nbi.ClustersIndexByName[vc.Clusters[vc.Host2Cluster[hostId]].Name]

		hostSite, err := common.MatchHostToSite(nbi, hostName, vc.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("hostSite: %s", err)
		}
		hostTenant, err := common.MatchHostToTenant(nbi, hostName, vc.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("hostTenant: %s", err)
		}
		hostAssetTag := host.Summary.Hardware.Uuid
		hostModel := host.Summary.Hardware.Model

		var hostSerialNumber string
		// find serial number from  host summary.hardware.OtherIdentifyingInfo (vmware specific logic)
		serialInfoTypes := map[string]bool{
			"EnclosureSerialNumberTag": true,
			"ServiceTag":               true,
			"SerialNumberTag":          true,
		}
		for _, info := range host.Summary.Hardware.OtherIdentifyingInfo {
			infoType := info.IdentifierType.GetElementDescription().Key
			if serialInfoTypes[infoType] {
				if info.IdentifierValue != "" {
					hostSerialNumber = info.IdentifierValue
					break
				}
			}
		}

		manufacturerName := host.Summary.Hardware.Vendor
		var hostManufacturer *objects.Manufacturer
		if manufacturerName == "" {
			manufacturerName = "Generic Manufacturer"
		}
		hostManufacturer, err = nbi.AddManufacturer(&objects.Manufacturer{
			Name: manufacturerName,
			Slug: utils.Slugify(manufacturerName),
		})
		if err != nil {
			return fmt.Errorf("failed adding vmware Manufacturer %v with error: %s", hostManufacturer, err)
		}

		var hostDeviceType *objects.DeviceType
		hostDeviceType, err = nbi.AddDeviceType(&objects.DeviceType{
			Manufacturer: hostManufacturer,
			Model:        hostModel,
			Slug:         utils.Slugify(hostModel),
		})
		if err != nil {
			return fmt.Errorf("failed adding vmware DeviceType %v with error: %s", hostDeviceType, err)
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
		platformName := utils.GeneratePlatformName(osType, osVersion)
		hostPlatform, err = nbi.AddPlatform(&objects.Platform{
			Name: platformName,
			Slug: utils.Slugify(platformName),
		})
		if err != nil {
			return fmt.Errorf("failed adding vmware Platform %v with error: %s", hostPlatform, err)
		}

		hostCpuCores := host.Summary.Hardware.NumCpuCores
		hostMemGB := host.Summary.Hardware.MemorySize / 1024 / 1024 / 1024

		nbHost := &objects.Device{
			NetboxObject: objects.NetboxObject{Tags: vc.CommonConfig.SourceTags, CustomFields: map[string]string{
				constants.CustomFieldSourceName:       vc.SourceConfig.Name,
				constants.CustomFieldHostCpuCoresName: fmt.Sprintf("%d", hostCpuCores),
				constants.CustomFieldHostMemoryName:   fmt.Sprintf("%d GB", hostMemGB),
			}},
			Name:         hostName,
			Status:       hostStatus,
			Platform:     hostPlatform,
			DeviceRole:   nbi.DeviceRolesIndexByName["Server"],
			Site:         hostSite,
			Tenant:       hostTenant,
			Cluster:      hostCluster,
			SerialNumber: hostSerialNumber,
			AssetTag:     hostAssetTag,
			DeviceType:   hostDeviceType,
		}
		nbHost, err = nbi.AddDevice(nbHost)
		if err != nil {
			return fmt.Errorf("failed to add vmware host %s with error: %v", host.Name, err)
		}

		// We also need to sync nics separately, because nic is a separate object in netbox
		err = vc.syncHostNics(nbi, host, nbHost)
		if err != nil {
			return fmt.Errorf("failed to sync vmware host %s nics with error: %v", host.Name, err)
		}
	}

	return nil
}

func (vc *VmwareSource) syncHostNics(nbi *inventory.NetboxInventory, vcHost mo.HostSystem, nbHost *objects.Device) error {

	// Sync host's physical interfaces
	err := vc.syncHostPhysicalNics(nbi, vcHost, nbHost)
	if err != nil {
		return fmt.Errorf("physical interfaces sync: %s", err)
	}

	// Sync host's virtual interfaces
	err = vc.syncHostVirtualNics(nbi, vcHost, nbHost)
	if err != nil {
		return fmt.Errorf("virtual interfaces sync: %s", err)
	}

	return nil
}

func (vc *VmwareSource) syncHostPhysicalNics(nbi *inventory.NetboxInventory, vcHost mo.HostSystem, nbHost *objects.Device) error {

	// Collect data from physical interfaces
	for _, pnic := range vcHost.Config.Network.Pnic {
		pnicName := pnic.Device
		var pnicLinkSpeed int32
		if pnic.LinkSpeed != nil {
			pnicLinkSpeed = pnic.LinkSpeed.SpeedMb
		}
		if pnicLinkSpeed == 0 {
			if pnic.Spec.LinkSpeed != nil {
				pnicLinkSpeed = pnic.LinkSpeed.SpeedMb
			}
			if pnicLinkSpeed == 0 {
				if pnic.ValidLinkSpecification != nil {
					pnicLinkSpeed = pnic.ValidLinkSpecification[0].SpeedMb
				}
			}
		}
		var pnicDescription string
		if pnicLinkSpeed >= 1000 {
			pnicDescription = fmt.Sprintf("%dGB/s", pnicLinkSpeed/1000)
		} else {
			pnicDescription = fmt.Sprintf("%dMB/s", pnicLinkSpeed)
		}
		pnicDescription += " pNIC"
		// netbox stores pnicSpeed in kbps
		pnicLinkSpeed *= 1000

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
		vlanIdMap := map[int]*objects.Vlan{} // set of vlans
		for portgroupName, portgroupData := range vc.Networks.HostPortgroups[nbHost.Name] {
			if slices.Contains(portgroupData.nics, pnicName) {
				if portgroupData.vlanId == 0 || portgroupData.vlanId > 4094 {
					vlanIdMap[portgroupData.vlanId] = &objects.Vlan{Vid: portgroupData.vlanId}
					continue
				}
				// Check if vlan with this vid already exists, else create it
				if vlanName, ok := vc.Networks.Vid2Name[portgroupData.vlanId]; ok {
					vlanGroup, err := common.MatchVlanToGroup(nbi, vlanName, vc.VlanGroupRelations)
					if err != nil {
						return fmt.Errorf("vlanGroup: %s", err)
					}
					vlanIdMap[portgroupData.vlanId] = nbi.VlansIndexByVlanGroupIdAndVid[vlanGroup.Id][portgroupData.vlanId]
				} else {
					vlanGroup, err := common.MatchVlanToGroup(nbi, portgroupName, vc.VlanGroupRelations)
					if err != nil {
						return fmt.Errorf("vlanGroup: %s", err)
					}
					var newVlan *objects.Vlan
					var ok bool
					newVlan, ok = nbi.VlansIndexByVlanGroupIdAndVid[vlanGroup.Id][portgroupData.vlanId]
					if !ok {
						newVlan, err = nbi.AddVlan(&objects.Vlan{
							NetboxObject: objects.NetboxObject{
								Tags: vc.CommonConfig.SourceTags,
								CustomFields: map[string]string{
									constants.CustomFieldSourceName: vc.SourceConfig.Name,
								},
							},
							Status: &objects.VlanStatusActive,
							Name:   fmt.Sprintf("VLAN%d_%s", portgroupData.vlanId, portgroupName),
							Vid:    portgroupData.vlanId,
							Group:  vlanGroup,
						})
						if err != nil {
							return fmt.Errorf("new vlan: %s", err)
						}
					}
					vlanIdMap[portgroupData.vlanId] = newVlan
				}
			}
		}

		// Determine interface mode for non VM traffic NIC, from vlans data
		var taggedVlanList []*objects.Vlan // when mode="tagged"
		if len(vlanIdMap) > 0 {
			if len(vlanIdMap) == 1 && vlanIdMap[0] != nil {
				pnicMode = &objects.InterfaceModeAccess
			} else if vlanIdMap[4095] != nil {
				pnicMode = &objects.InterfaceModeTaggedAll
			} else {
				pnicMode = &objects.InterfaceModeTagged
			}
			taggedVlanList = []*objects.Vlan{}
			if pnicMode == &objects.InterfaceModeTagged {
				for vid, vlan := range vlanIdMap {
					if vid == 0 {
						continue
					}
					taggedVlanList = append(taggedVlanList, vlan)
				}
			}
		}

		pnicType := objects.IfaceSpeed2IfaceType[objects.InterfaceSpeed(pnicLinkSpeed)]
		if pnicType == nil {
			pnicType = &objects.OtherInterfaceType
		}

		// After collecting all of the data add interface to nbi
		_, err := nbi.AddInterface(&objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        vc.CommonConfig.SourceTags,
				Description: pnicDescription,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: vc.SourceConfig.Name,
				},
			},
			Device:      nbHost,
			Name:        pnicName,
			Status:      true,
			Type:        pnicType,
			Speed:       objects.InterfaceSpeed(pnicLinkSpeed),
			MTU:         pnicMtu,
			MAC:         strings.ToUpper(pnic.Mac),
			Mode:        pnicMode,
			TaggedVlans: taggedVlanList,
		})
		if err != nil {
			return fmt.Errorf("failed adding physical interface: %s", err)
		}
	}
	return nil
}

func (vc *VmwareSource) syncHostVirtualNics(nbi *inventory.NetboxInventory, vcHost mo.HostSystem, nbHost *objects.Device) error {
	// Collect data over all virtual interfaces
	for _, vnic := range vcHost.Config.Network.Vnic {
		vnicName := vnic.Device
		vnicPortgroupData, vnicPortgroupDataOk := vc.Networks.HostPortgroups[vcHost.Name][vnic.Portgroup]
		vnicDvPortgroupKey := ""
		if vnic.Spec.DistributedVirtualPort != nil {
			vnicDvPortgroupKey = vnic.Spec.DistributedVirtualPort.PortgroupKey
		}
		vnicDvPortgroupData, vnicDvPortgroupDataOk := vc.Networks.DistributedVirtualPortgroups[vnicDvPortgroupKey]
		vnicPortgroupVlanId := 0
		vnicDvPortgroupVlanIds := []int{}
		var vnicMode *objects.InterfaceMode
		var vlanDescription, vnicDescription string

		// Get data from local portgroup, or distributed portgroup
		if vnicPortgroupDataOk {
			vnicPortgroupVlanId = vnicPortgroupData.vlanId
			vnicSwitch := vnicPortgroupData.vswitch
			vnicDescription = fmt.Sprintf("%s (%s, vlan ID: %d)", vnic.Portgroup, vnicSwitch, vnicPortgroupVlanId)
		} else if vnicDvPortgroupDataOk {
			vnicDescription = vnicDvPortgroupData.Name
			vnicDvPortgroupVlanIds = vnicDvPortgroupData.VlanIds
			if len(vnicDvPortgroupVlanIds) == 1 && vnicDvPortgroupData.VlanIds[0] == 4095 {
				vnicDescription = "all vlans"
				vnicMode = &objects.InterfaceModeTaggedAll
			} else {
				if len(vnicDvPortgroupData.VlanIdRanges) > 0 {
					vlanDescription = fmt.Sprintf("vlan IDs: %s", strings.Join(vnicDvPortgroupData.VlanIdRanges, ","))
				} else {
					vlanDescription = fmt.Sprintf("vlan ID: %d", vnicDvPortgroupData.VlanIds[0])
				}
				if len(vnicDvPortgroupData.VlanIds) == 1 && vnicDvPortgroupData.VlanIds[0] == 0 {
					vnicMode = &objects.InterfaceModeAccess
				} else {
					vnicMode = &objects.InterfaceModeTagged
				}
			}
			vnicDvPortgroupDwSwitchUuid := vnic.Spec.DistributedVirtualPort.SwitchUuid
			vnicVswitch, vnicVswitchOk := vc.Networks.HostVirtualSwitches[vcHost.Name][vnicDvPortgroupDwSwitchUuid]
			if vnicVswitchOk {
				vnicDescription = fmt.Sprintf("%s (%v, %s)", vnicDescription, vnicVswitch, vlanDescription)
			}
		}

		var vnicUntaggedVlan *objects.Vlan
		var vnicTaggedVlans []*objects.Vlan
		if vnicPortgroupData != nil && vnicPortgroupVlanId != 0 {
			vnicUntaggedVlanGroup, err := common.MatchVlanToGroup(nbi, vc.Networks.Vid2Name[vnicPortgroupVlanId], vc.VlanGroupRelations)
			if err != nil {
				return fmt.Errorf("vlan group: %s", err)
			}
			vnicUntaggedVlan = nbi.VlansIndexByVlanGroupIdAndVid[vnicUntaggedVlanGroup.Id][vnicPortgroupVlanId]
			vnicMode = &objects.InterfaceModeAccess
			// vnicUntaggedVlan = &objects.Vlan{
			// 	Name:   fmt.Sprintf("ESXi %s (ID: %d) (%s)", vnic.Portgroup, vnicPortgroupVlanId, nbHost.Site.Name),
			// 	Vid:    vnicPortgroupVlanId,
			// 	Tenant: nbHost.Tenant,
			// }
		} else if vnicDvPortgroupData != nil {
			for _, vnicDvPortgroupDataVlanId := range vnicDvPortgroupVlanIds {
				if vnicMode != &objects.InterfaceModeTagged {
					break
				}
				if vnicDvPortgroupDataVlanId == 0 {
					continue
				}
				vnicTaggedVlanGroup, err := common.MatchVlanToGroup(nbi, vc.Networks.Vid2Name[vnicDvPortgroupDataVlanId], vc.VlanGroupRelations)
				if err != nil {
					return fmt.Errorf("vlan group: %s", err)
				}
				vnicTaggedVlans = append(vnicTaggedVlans, nbi.VlansIndexByVlanGroupIdAndVid[vnicTaggedVlanGroup.Id][vnicDvPortgroupDataVlanId])
				// vnicTaggedVlans = append(vnicTaggedVlans, &objects.Vlan{
				// 	Name:   fmt.Sprintf("%s-%d", vnicDvPortgroupData.Name, vnicDvPortgroupDataVlanId),
				// 	Vid:    vnicDvPortgroupDataVlanId,
				// 	Tenant: nbHost.Tenant,
				// })
			}
		}

		nbVnic, err := nbi.AddInterface(&objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        vc.CommonConfig.SourceTags,
				Description: vnicDescription,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: vc.SourceConfig.Name,
				},
			},
			Device:       nbHost,
			Name:         vnicName,
			Status:       true,
			Type:         &objects.VirtualInterfaceType,
			MTU:          int(vnic.Spec.Mtu),
			Mode:         vnicMode,
			TaggedVlans:  vnicTaggedVlans,
			UntaggedVlan: vnicUntaggedVlan,
		})
		if err != nil {
			return err
		}

		var ipv4Address *objects.IPAddress
		// Get IPv4 address for this vnic. TODO: filter
		ipv4_address := vnic.Spec.Ip.IpAddress
		ipv4_mask_bits, err := utils.MaskToBits(vnic.Spec.Ip.SubnetMask)
		if err != nil {
			return fmt.Errorf("mask to bits: %s", err)
		}
		ipv4_dns := utils.ReverseLookup(ipv4_address)
		ipv4Address, err = nbi.AddIPAddress(&objects.IPAddress{
			NetboxObject: objects.NetboxObject{
				Tags: vc.CommonConfig.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: vc.SourceConfig.Name,
				},
			},
			Address:            fmt.Sprintf("%s/%d", ipv4_address, ipv4_mask_bits),
			Status:             &objects.IPAddressStatusActive, // TODO
			DNSName:            ipv4_dns,
			Tenant:             nbHost.Tenant,
			AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
			AssignedObjectId:   nbVnic.Id,
		})
		if err != nil {
			return err
		}

		var ipv6Address *objects.IPAddress
		if vnic.Spec.Ip.IpV6Config != nil {
			for _, ipv6_entry := range vnic.Spec.Ip.IpV6Config.IpV6Address {
				ipv6_address := ipv6_entry.IpAddress
				ipv6_mask := ipv6_entry.PrefixLength
				// TODO: Filter out ipv6 addresses
				ipv6Address, err = nbi.AddIPAddress(&objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: vc.CommonConfig.SourceTags,
						CustomFields: map[string]string{
							constants.CustomFieldSourceName: vc.SourceConfig.Name,
						},
					},
					Address:            fmt.Sprintf("%s/%d", ipv6_address, ipv6_mask),
					Status:             &objects.IPAddressStatusActive, // TODO
					Tenant:             nbHost.Tenant,
					AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
					AssignedObjectId:   nbVnic.Id,
				})
				if err != nil {
					return err
				}
			}
		}

		// Update host's primary ipv4: TODO, determine if primary or not
		if nbHost.PrimaryIPv4 == nil && ipv4Address != nil {
			newNbHost := *nbHost
			newNbHost.PrimaryIPv4 = ipv4Address
			nbHost, err = nbi.AddDevice(&newNbHost)
			if err != nil {
				return fmt.Errorf("new Host's primaryIpv4: %s", err)
			}
		}

		// Update host's primary ipv4: TODO, determine if primary or not
		if nbHost.PrimaryIPv6 == nil && ipv6Address != nil {
			newNbHost := *nbHost
			newNbHost.PrimaryIPv6 = ipv6Address
			nbHost, err = nbi.AddDevice(&newNbHost)
			if err != nil {
				return fmt.Errorf("new Host's primaryIpv6: %s", err)
			}
		}
	}
	return nil
}

func (vc *VmwareSource) syncVms(nbi *inventory.NetboxInventory) error {
	for vmKey, vm := range vc.Vms {
		// Check if vm is a template, we don't add templates into netbox.
		if vm.Config != nil {
			if vm.Config.Template {
				continue
			}
		}

		vmName := vm.Name
		vmHostName := vc.Hosts[vc.Vm2Host[vmKey]].Name

		// Tenant is received from VmTenantRelations
		vmTenant, err := common.MatchVmToTenant(nbi, vmName, vc.VmTenantRelations)
		if err != nil {
			return fmt.Errorf("vm's Tenant: %s", err)
		}

		// Site is the same as the Host
		vmSite, err := common.MatchHostToSite(nbi, vmHostName, vc.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("vm's Site: %s", err)
		}
		vmHost := nbi.DevicesIndexByNameAndSiteId[vmHostName][vmSite.Id]

		// Cluster of the vm is same as the host
		vmCluster := vmHost.Cluster

		// VM status
		vmStatus := &objects.VMStatusOffline
		vmPowerState := vm.Runtime.PowerState
		if vmPowerState == types.VirtualMachinePowerStatePoweredOn {
			vmStatus = &objects.VMStatusActive
		}

		// vmVCPUs
		vmVCPUs := vm.Config.Hardware.NumCPU

		// vmMemory
		vmMemory := vm.Config.Hardware.MemoryMB

		// DisksSize
		vmDiskSizeB := int64(0)
		for _, hwDevice := range vm.Config.Hardware.Device {
			if disk, ok := hwDevice.(*types.VirtualDisk); ok {
				vmDiskSizeB += disk.CapacityInBytes
			}
		}

		// vmPlatform
		vmPlatformName := vm.Config.GuestFullName
		if vmPlatformName == "" {
			vmPlatformName = vm.Guest.GuestFullName
		}
		if vmPlatformName == "" {
			vmPlatformName = utils.GeneratePlatformName("Generic OS", "Generic Version")
		}
		vmPlatform, err := nbi.AddPlatform(&objects.Platform{
			Name: vmPlatformName,
			Slug: utils.Slugify(vmPlatformName),
		})
		if err != nil {
			return fmt.Errorf("failed adding vmware vm's Platform %v with error: %s", vmPlatform, err)
		}

		// Extract additional info from CustomFields
		var vmOwners []string
		var vmOwnerEmails []string
		var vmDescription string
		vmCustomFields := map[string]string{}
		if len(vm.Summary.CustomValue) > 0 {
			for _, field := range vm.Summary.CustomValue {
				if field, ok := field.(*types.CustomFieldStringValue); ok {
					fieldName := vc.CustomFieldId2Name[field.Key]

					if mappedField, ok := vc.CustomFieldMappings[fieldName]; ok {
						switch mappedField {
						case "owner":
							vmOwners = strings.Split(field.Value, ",")
						case "email":
							vmOwnerEmails = strings.Split(field.Value, ",")
						case "description":
							vmDescription = strings.TrimSpace(field.Value)
						}
					} else {
						fieldName = utils.Alphanumeric(fieldName)
						if _, ok := nbi.CustomFieldsIndexByName[fieldName]; !ok {
							err := nbi.AddCustomField(&objects.CustomField{
								Name:                  fieldName,
								Type:                  objects.CustomFieldTypeText,
								CustomFieldUIVisible:  &objects.CustomFieldUIVisibleIfSet,
								CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
								ContentTypes:          []string{"virtualization.virtualmachine"},
							})
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
		vmCustomFields[constants.CustomFieldSourceIdName] = vm.Self.Value

		// netbox description has constraint <= len(200 characters)
		// In this case we make a comment
		var vmComments string
		if len(vmDescription) >= 200 {
			vmDescription = "See comments."
			vmComments = vmDescription
		}

		newVM, err := nbi.AddVM(&objects.VM{
			NetboxObject: objects.NetboxObject{
				Tags:         vc.CommonConfig.SourceTags,
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
			Memory:   int(vmMemory),                         // MBs
			Disk:     int(vmDiskSizeB / 1024 / 1024 / 1024), // GBs
			Comments: vmComments,
		})

		if err != nil {
			return fmt.Errorf("failed to sync vmware vm: %v", err)
		}

		// If vm owner name was found we also add contact assignment to the vm
		var vmMailMapFallback bool
		if len(vmOwners) > 0 && len(vmOwnerEmails) > 0 && len(vmOwners) != len(vmOwnerEmails) {
			vc.Logger.Warningf("vm owner names and emails mismatch (len(vmOwnerEmails) != len(vmOwners), using fallback mechanism")
			vmMailMapFallback = true
		}
		vmOwner2Email := utils.MatchNamesWithEmails(vmOwners, vmOwnerEmails, vc.Logger)
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
				contact, err := nbi.AddContact(
					&objects.Contact{
						Name:  strings.TrimSpace(vmOwners[i]),
						Email: vmOwnerEmail,
					},
				)
				if err != nil {
					return fmt.Errorf("creating vm contact: %s", err)
				}
				_, err = nbi.AddContactAssignment(&objects.ContactAssignment{
					ContentType: "virtualization.virtualmachine",
					ObjectId:    newVM.Id,
					Contact:     contact,
					Role:        nbi.ContactRolesIndexByName[objects.AdminContactRoleName],
				})
				if err != nil {
					return fmt.Errorf("add contact assignment for vm: %s", err)
				}
			}
		}

		// Sync vm interfaces
		err = vc.syncVmInterfaces(nbi, vm, newVM)
		if err != nil {
			return fmt.Errorf("failed to sync vmware vm's interfaces: %v", err)
		}
	}
	return nil
}

// Syncs VM's interfaces to Netbox
func (vc *VmwareSource) syncVmInterfaces(nbi *inventory.NetboxInventory, vmwareVm mo.VirtualMachine, netboxVm *objects.VM) error {
	var vmPrimaryIpv4 *objects.IPAddress
	var vmPrimaryIpv6 *objects.IPAddress
	var vmDefaultGatewayIpv4 string
	var vmDefaultGatewayIpv6 string

	// From vm's routing determine the default interface
	if len(vmwareVm.Guest.IpStack) > 0 {
		for _, route := range vmwareVm.Guest.IpStack[0].IpRouteConfig.IpRoute {
			if route.PrefixLength == 0 {
				ipAddress := route.Network
				if ipAddress == "" {
					continue
				}
				gatewayIpAddress := route.Gateway.IpAddress
				if gatewayIpAddress == "" {
					continue
				}

				// Get version from ipAddress (v4 or v6)
				ipVersion := utils.GetIPVersion(ipAddress)
				if ipVersion == 4 {
					vmDefaultGatewayIpv4 = gatewayIpAddress
				} else if ipVersion == 6 {
					vmDefaultGatewayIpv6 = gatewayIpAddress
				}
			}
		}
	}

	nicIps := map[string][]string{}

	for _, vmDevice := range vmwareVm.Config.Hardware.Device {
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
			intMac := vmEthernetCard.MacAddress
			intConnected := vmEthernetCard.Connectable.Connected
			intDeviceBackingInfo := vmEthernetCard.Backing
			intDeviceInfo := vmEthernetCard.DeviceInfo
			var intMtu int
			var intNetworkName string
			var intNetworkPrivate bool
			var intMode *objects.VMInterfaceMode
			intNetworkVlanIds := []int{}
			intNetworkVlanIdRanges := []string{}

			// Get info from local vSwitches if possible, else from DistributedPortGroup
			if backingInfo, ok := intDeviceBackingInfo.(*types.VirtualEthernetCardNetworkBackingInfo); ok {
				intNetworkName = backingInfo.DeviceName
				intHostPgroup := vc.Networks.HostPortgroups[netboxVm.Host.Name][intNetworkName]

				if intHostPgroup != nil {
					intNetworkVlanIds = []int{intHostPgroup.vlanId}
					intNetworkVlanIdRanges = []string{strconv.Itoa(intHostPgroup.vlanId)}
					intVswitchName := intHostPgroup.vswitch
					intVswitchData := vc.Networks.HostVirtualSwitches[netboxVm.Host.Name][intVswitchName]
					if intVswitchData != nil {
						intMtu = intVswitchData.mtu
					}

				}
			} else if backingInfo, ok := intDeviceBackingInfo.(*types.VirtualEthernetCardDistributedVirtualPortBackingInfo); ok {
				dvsPortgroupKey := backingInfo.Port.PortgroupKey
				intPortgroupData := vc.Networks.DistributedVirtualPortgroups[dvsPortgroupKey]

				if intPortgroupData != nil {
					intNetworkName = intPortgroupData.Name
					intNetworkVlanIds = intPortgroupData.VlanIds
					intNetworkVlanIdRanges = intPortgroupData.VlanIdRanges
					if len(intNetworkVlanIdRanges) == 0 {
						intNetworkVlanIdRanges = []string{strconv.Itoa(intNetworkVlanIds[0])}
					}
					intNetworkPrivate = intPortgroupData.Private
				}

				intDvswitchUuid := backingInfo.Port.SwitchUuid
				intDvswitchData := vc.Networks.HostProxySwitches[netboxVm.Host.Name][intDvswitchUuid]

				if intDvswitchData != nil {
					intMtu = intDvswitchData.mtu
				}
			}

			var vlanDescription string
			intLabel := intDeviceInfo.GetDescription().Label
			splitStr := strings.Split(intLabel, " ")
			intName := fmt.Sprintf("vNic %s", splitStr[len(splitStr)-1])
			intFullName := intName
			if intNetworkName != "" {
				intFullName = fmt.Sprintf("%s (%s)", intFullName, intNetworkName)
			}
			intDescription := intLabel
			if len(intNetworkVlanIds) > 0 {
				if len(intNetworkVlanIds) == 1 && intNetworkVlanIds[0] == 4095 {
					vlanDescription = "all vlans"
					intMode = &objects.VMInterfaceModeTaggedAll
				} else {
					vlanDescription = fmt.Sprintf("vlan ID: %s", strings.Join(intNetworkVlanIdRanges, ", "))
					if len(intNetworkVlanIds) == 1 {
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
			for _, guestNic := range vmwareVm.Guest.Net {
				if intMac != guestNic.MacAddress {
					continue
				}
				intConnected = guestNic.Connected

				if _, ok := nicIps[intFullName]; !ok {
					nicIps[intFullName] = []string{}
				}

				if guestNic.IpConfig != nil {
					for _, intIp := range guestNic.IpConfig.IpAddress {
						intIpAddress := fmt.Sprintf("%s/%d", intIp.IpAddress, intIp.PrefixLength)
						nicIps[intFullName] = append(nicIps[intFullName], intIpAddress)

						// Check if primary gateways are in the subnet of this IP address
						// if it matches IP gets chosen as primary ip
						if vmDefaultGatewayIpv4 != "" && utils.SubnetContainsIpAddress(vmDefaultGatewayIpv4, intIpAddress) {
							ipDns := utils.ReverseLookup(intIp.IpAddress)
							vmPrimaryIpv4 = &objects.IPAddress{
								NetboxObject: objects.NetboxObject{
									Tags: vc.CommonConfig.SourceTags,
									CustomFields: map[string]string{
										constants.CustomFieldSourceName: vc.SourceConfig.Name,
									},
								},
								Address: intIpAddress,
								Status:  &objects.IPAddressStatusActive,
								DNSName: ipDns,
							}
						}
						if vmDefaultGatewayIpv6 != "" && utils.SubnetContainsIpAddress(vmDefaultGatewayIpv6, intIpAddress) {
							ipDns := utils.ReverseLookup(intIp.IpAddress)
							vmPrimaryIpv6 = &objects.IPAddress{
								NetboxObject: objects.NetboxObject{
									Tags: vc.CommonConfig.SourceTags,
									CustomFields: map[string]string{
										constants.CustomFieldSourceName: vc.SourceConfig.Name,
									},
								},
								Address: intIpAddress,
								Status:  &objects.IPAddressStatusActive,
								DNSName: ipDns,
							}
						}
					}
				}
			}
			var intUntaggedVlan *objects.Vlan
			var intTaggedVlanList []*objects.Vlan
			if len(intNetworkVlanIds) > 0 && intMode != &objects.VMInterfaceModeTaggedAll {
				if len(intNetworkVlanIds) == 1 && intNetworkVlanIds[0] != 0 {
					vidId := intNetworkVlanIds[0]
					nicUntaggedVlanGroup, err := common.MatchVlanToGroup(nbi, vc.Networks.Vid2Name[vidId], vc.VlanGroupRelations)
					if err != nil {
						return fmt.Errorf("vlan group: %s", err)
					}
					intUntaggedVlan = nbi.VlansIndexByVlanGroupIdAndVid[nicUntaggedVlanGroup.Id][vidId]
				} else {
					intTaggedVlanList = []*objects.Vlan{}
					for _, intNetworkVlanId := range intNetworkVlanIds {
						if intNetworkVlanId == 0 {
							continue
						}
						// nicTaggedVlanList = append(nicTaggedVlanList, nbi.get[intNetworkVlanId])
					}
				}
			}
			nbVmInterface, err := nbi.AddVMInterface(&objects.VMInterface{
				NetboxObject: objects.NetboxObject{
					Tags:        vc.CommonConfig.SourceTags,
					Description: intDescription,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: vc.SourceConfig.Name,
					},
				},
				VM:           netboxVm,
				Name:         intFullName,
				MACAddress:   strings.ToUpper(intMac),
				MTU:          intMtu,
				Mode:         intMode,
				Enabled:      intConnected,
				TaggedVlans:  intTaggedVlanList,
				UntaggedVlan: intUntaggedVlan,
			})
			if err != nil {
				return fmt.Errorf("adding VmInterface: %s", err)
			}

			// Add primary IPs to the netbox
			if vmPrimaryIpv4 != nil {
				vmPrimaryIpv4.AssignedObjectType = objects.AssignedObjectTypeVMInterface
				vmPrimaryIpv4.AssignedObjectId = nbVmInterface.Id
				vmPrimaryIpv4, err = nbi.AddIPAddress(vmPrimaryIpv4)
				if err != nil {
					vc.Logger.Warningf("adding vm's primary ipv4: %s", err)
				}
			}
			if vmPrimaryIpv6 != nil {
				vmPrimaryIpv6.AssignedObjectType = objects.AssignedObjectTypeVMInterface
				vmPrimaryIpv6.AssignedObjectId = nbVmInterface.Id
				vmPrimaryIpv6, err = nbi.AddIPAddress(vmPrimaryIpv6)
				if err != nil {
					vc.Logger.Warningf("adding vm's primary ipv6: %s", err)
				}
			}

			// Update the vms with primary addresses
			if vmPrimaryIpv4 != nil && (netboxVm.PrimaryIPv4 == nil || vmPrimaryIpv4.Address != netboxVm.PrimaryIPv4.Address) {
				// Shallow copy netboxVm to newNetboxVM
				newNetboxVm := *netboxVm
				newNetboxVm.PrimaryIPv4 = vmPrimaryIpv4
				_, err = nbi.AddVM(&newNetboxVm)
				if err != nil {
					vc.Logger.Warningf("adding vm: %s", err)
				}
			}
		}
	}

	return nil
}
