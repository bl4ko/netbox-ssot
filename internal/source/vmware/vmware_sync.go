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

func (vc *Source) syncNetworks(nbi *inventory.NetboxInventory) error {
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
		if len(dvpg.VlanIDs) == 1 && len(dvpg.VlanIDRanges) == 0 {
			_, err := nbi.AddVlan(&objects.Vlan{
				NetboxObject: objects.NetboxObject{
					Tags: vc.Config.SourceTags,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: vc.SourceConfig.Name,
					},
				},
				Name:   dvpg.Name,
				Group:  vlanGroup,
				Vid:    dvpg.VlanIDs[0],
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

func (vc *Source) syncDatacenters(nbi *inventory.NetboxInventory) error {
	for _, dc := range vc.DataCenters {
		nbClusterGroup := &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{
				Description: fmt.Sprintf("Datacenter from source %s", vc.SourceConfig.Hostname),
				Tags:        vc.Config.SourceTags,
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

func (vc *Source) syncClusters(nbi *inventory.NetboxInventory) error {
	clusterType := &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: vc.Config.SourceTags,
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
	for clusterID, cluster := range vc.Clusters {
		clusterName := cluster.Name

		var clusterGroup *objects.ClusterGroup
		datacenterID := vc.Cluster2Datacenter[clusterID]
		clusterGroup = nbi.ClusterGroupsIndexByName[vc.DataCenters[datacenterID].Name]

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
				Tags: vc.Config.SourceTags,
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
// custom role Server.
func (vc *Source) syncHosts(nbi *inventory.NetboxInventory) error {
	for hostID, host := range vc.Hosts {
		var err error
		hostName := host.Name
		hostCluster := nbi.ClustersIndexByName[vc.Clusters[vc.Host2Cluster[hostID]].Name]

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

		hostCPUCores := host.Summary.Hardware.NumCpuCores
		hostMemGB := host.Summary.Hardware.MemorySize / constants.KiB / constants.KiB / constants.KiB

		nbHost := &objects.Device{
			NetboxObject: objects.NetboxObject{Tags: vc.Config.SourceTags, CustomFields: map[string]string{
				constants.CustomFieldSourceName:       vc.SourceConfig.Name,
				constants.CustomFieldHostCPUCoresName: fmt.Sprintf("%d", hostCPUCores),
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

func (vc *Source) syncHostNics(nbi *inventory.NetboxInventory, vcHost mo.HostSystem, nbHost *objects.Device) error {
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

func (vc *Source) syncHostPhysicalNics(nbi *inventory.NetboxInventory, vcHost mo.HostSystem, nbHost *objects.Device) error {
	// Collect data from physical interfaces
	for _, pnic := range vcHost.Config.Network.Pnic {
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
		if pnicLinkSpeedMb*constants.MB >= constants.GB {
			pnicDescription = fmt.Sprintf("%dGB/s", pnicLinkSpeedMb/constants.KB)
		} else {
			pnicDescription = fmt.Sprintf("%dMB/s", pnicLinkSpeedMb)
		}
		pnicDescription += " pNIC"
		// netbox stores pnicSpeed in kbps

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
					vlanGroup, err := common.MatchVlanToGroup(nbi, vlanName, vc.VlanGroupRelations)
					if err != nil {
						return fmt.Errorf("vlanGroup: %s", err)
					}
					vlanIDMap[portgroupData.vlanID] = nbi.VlansIndexByVlanGroupIDAndVID[vlanGroup.ID][portgroupData.vlanID]
				} else {
					vlanGroup, err := common.MatchVlanToGroup(nbi, portgroupName, vc.VlanGroupRelations)
					if err != nil {
						return fmt.Errorf("vlanGroup: %s", err)
					}
					var newVlan *objects.Vlan
					var ok bool
					newVlan, ok = nbi.VlansIndexByVlanGroupIDAndVID[vlanGroup.ID][portgroupData.vlanID]
					if !ok {
						newVlan, err = nbi.AddVlan(&objects.Vlan{
							NetboxObject: objects.NetboxObject{
								Tags: vc.Config.SourceTags,
								CustomFields: map[string]string{
									constants.CustomFieldSourceName: vc.SourceConfig.Name,
								},
							},
							Status: &objects.VlanStatusActive,
							Name:   fmt.Sprintf("VLAN%d_%s", portgroupData.vlanID, portgroupName),
							Vid:    portgroupData.vlanID,
							Group:  vlanGroup,
						})
						if err != nil {
							return fmt.Errorf("new vlan: %s", err)
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

		pnicType := objects.IfaceSpeed2IfaceType[objects.InterfaceSpeed(pnicLinkSpeedMb)]
		if pnicType == nil {
			pnicType = &objects.OtherInterfaceType
		}

		// After collecting all of the data add interface to nbi
		_, err := nbi.AddInterface(&objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        vc.Config.SourceTags,
				Description: pnicDescription,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: vc.SourceConfig.Name,
				},
			},
			Device:      nbHost,
			Name:        pnicName,
			Status:      true,
			Type:        pnicType,
			Speed:       objects.InterfaceSpeed(pnicLinkSpeedMb / constants.KB),
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

func (vc *Source) syncHostVirtualNics(nbi *inventory.NetboxInventory, vcHost mo.HostSystem, nbHost *objects.Device) error {
	// Collect data over all virtual interfaces
	for _, vnic := range vcHost.Config.Network.Vnic {
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
			vnicDescription = fmt.Sprintf("%s (%s, vlan ID: %d)", vnic.Portgroup, vnicSwitch, vnicPortgroupVlanID)
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
			vnicUntaggedVlanGroup, err := common.MatchVlanToGroup(nbi, vc.Networks.Vid2Name[vnicPortgroupVlanID], vc.VlanGroupRelations)
			if err != nil {
				return fmt.Errorf("vlan group: %s", err)
			}
			vnicUntaggedVlan = nbi.VlansIndexByVlanGroupIDAndVID[vnicUntaggedVlanGroup.ID][vnicPortgroupVlanID]
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
				vnicTaggedVlanGroup, err := common.MatchVlanToGroup(nbi, vc.Networks.Vid2Name[vnicDvPortgroupDataVlanID], vc.VlanGroupRelations)
				if err != nil {
					return fmt.Errorf("vlan group: %s", err)
				}
				vnicTaggedVlans = append(vnicTaggedVlans, nbi.VlansIndexByVlanGroupIDAndVID[vnicTaggedVlanGroup.ID][vnicDvPortgroupDataVlanID])
				// vnicTaggedVlans = append(vnicTaggedVlans, &objects.Vlan{
				// 	Name:   fmt.Sprintf("%s-%d", vnicDvPortgroupData.Name, vnicDvPortgroupDataVlanId),
				// 	Vid:    vnicDvPortgroupDataVlanId,
				// 	Tenant: nbHost.Tenant,
				// })
			}
		}

		nbVnic, err := nbi.AddInterface(&objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        vc.Config.SourceTags,
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

		var nbIPAddress *objects.IPAddress
		// Get IPv4 address for this vnic. TODO: filter
		ipv4Address := vnic.Spec.Ip.IpAddress
		ipv4MaskBits, err := utils.MaskToBits(vnic.Spec.Ip.SubnetMask)
		if err != nil {
			return fmt.Errorf("mask to bits: %s", err)
		}
		ipv4DNS := utils.ReverseLookup(ipv4Address)
		nbIPAddress, err = nbi.AddIPAddress(&objects.IPAddress{
			NetboxObject: objects.NetboxObject{
				Tags: vc.Config.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: vc.SourceConfig.Name,
				},
			},
			Address:            fmt.Sprintf("%s/%d", ipv4Address, ipv4MaskBits),
			Status:             &objects.IPAddressStatusActive, // TODO
			DNSName:            ipv4DNS,
			Tenant:             nbHost.Tenant,
			AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
			AssignedObjectID:   nbVnic.ID,
		})
		if err != nil {
			return err
		}

		var nbIPv6Address *objects.IPAddress
		if vnic.Spec.Ip.IpV6Config != nil {
			for _, ipv6Entry := range vnic.Spec.Ip.IpV6Config.IpV6Address {
				ipv6Address := ipv6Entry.IpAddress
				ipv6Mask := ipv6Entry.PrefixLength
				// TODO: Filter out ipv6 addresses
				nbIPv6Address, err = nbi.AddIPAddress(&objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: vc.Config.SourceTags,
						CustomFields: map[string]string{
							constants.CustomFieldSourceName: vc.SourceConfig.Name,
						},
					},
					Address:            fmt.Sprintf("%s/%d", ipv6Address, ipv6Mask),
					Status:             &objects.IPAddressStatusActive, // TODO
					Tenant:             nbHost.Tenant,
					AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
					AssignedObjectID:   nbVnic.ID,
				})
				if err != nil {
					return err
				}
			}
		}

		// Update host's primary ipv4: TODO, determine if primary or not
		if nbHost.PrimaryIPv4 == nil && nbIPAddress != nil {
			newNbHost := *nbHost
			newNbHost.PrimaryIPv4 = nbIPAddress
			nbHost, err = nbi.AddDevice(&newNbHost)
			if err != nil {
				return fmt.Errorf("new Host's primaryIpv4: %s", err)
			}
		}

		// Update host's primary ipv4: TODO, determine if primary or not
		if nbHost.PrimaryIPv6 == nil && nbIPv6Address != nil {
			newNbHost := *nbHost
			newNbHost.PrimaryIPv6 = nbIPv6Address
			nbHost, err = nbi.AddDevice(&newNbHost)
			if err != nil {
				return fmt.Errorf("new Host's primaryIpv6: %s", err)
			}
		}
	}
	return nil
}

func (vc *Source) syncVms(nbi *inventory.NetboxInventory) error {
	for vmKey, vm := range vc.Vms {
		// Check if vm is a template, we don't add templates into netbox.
		if vm.Config != nil {
			if vm.Config.Template {
				continue
			}
		}

		vmName := vm.Name
		vmHostName := vc.Hosts[vc.VM2Host[vmKey]].Name

		// Tenant is received from VmTenantRelations
		vmTenant, err := common.MatchVMToTenant(nbi, vmName, vc.VMTenantRelations)
		if err != nil {
			return fmt.Errorf("vm's Tenant: %s", err)
		}

		// Site is the same as the Host
		vmSite, err := common.MatchHostToSite(nbi, vmHostName, vc.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("vm's Site: %s", err)
		}
		vmHost := nbi.DevicesIndexByNameAndSiteID[vmHostName][vmSite.ID]

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
					fieldName := vc.CustomFieldID2Name[field.Key]

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
		vmCustomFields[constants.CustomFieldSourceIDName] = vm.Self.Value

		// netbox description has constraint <= len(200 characters)
		// In this case we make a comment
		var vmComments string
		if len(vmDescription) >= objects.MaxDescriptionLength {
			vmDescription = "See comments."
			vmComments = vmDescription
		}

		newVM, err := nbi.AddVM(&objects.VM{
			NetboxObject: objects.NetboxObject{
				Tags:         vc.Config.SourceTags,
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
			Memory:   int(vmMemory),                                                    // MBs
			Disk:     int(vmDiskSizeB / constants.KiB / constants.KiB / constants.KiB), // GBs
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
					ObjectID:    newVM.ID,
					Contact:     contact,
					Role:        nbi.ContactRolesIndexByName[objects.AdminContactRoleName],
				})
				if err != nil {
					return fmt.Errorf("add contact assignment for vm: %s", err)
				}
			}
		}

		// Sync vm interfaces
		err = vc.syncVMInterfaces(nbi, vm, newVM)
		if err != nil {
			return fmt.Errorf("failed to sync vmware vm's interfaces: %v", err)
		}
	}
	return nil
}

// Syncs VM's interfaces to Netbox.
func (vc *Source) syncVMInterfaces(nbi *inventory.NetboxInventory, vmwareVM mo.VirtualMachine, netboxVM *objects.VM) error {
	var vmDefaultGatewayIpv4 string
	var vmDefaultGatewayIpv6 string

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
			intMac := vmEthernetCard.MacAddress
			intConnected := vmEthernetCard.Connectable.Connected
			intDeviceBackingInfo := vmEthernetCard.Backing
			intDeviceInfo := vmEthernetCard.DeviceInfo
			nicIPv4Addresses := []string{}
			primaryIPv4Address := ""
			nicIPv6Addresses := []string{}
			primaryIPv6Address := ""
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
			intName := fmt.Sprintf("vNic %s", splitStr[len(splitStr)-1])
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
							if vmDefaultGatewayIpv4 != "" && utils.SubnetContainsIPAddress(vmDefaultGatewayIpv4, intIPAddress) {
								primaryIPv4Address = intIPAddress
							}
						case constants.IPv6:
							nicIPv6Addresses = append(nicIPv6Addresses, intIPAddress)
							if vmDefaultGatewayIpv6 != "" && utils.SubnetContainsIPAddress(vmDefaultGatewayIpv6, intIPAddress) {
								primaryIPv6Address = intIPAddress
							}
						default:
							return fmt.Errorf("unknown ip version: %s", intIPAddress)
						}
					}
				}
			}
			var intUntaggedVlan *objects.Vlan
			var intTaggedVlanList []*objects.Vlan
			if len(intNetworkVlanIDs) > 0 && intMode != &objects.VMInterfaceModeTaggedAll {
				if len(intNetworkVlanIDs) == 1 && intNetworkVlanIDs[0] != 0 {
					vidID := intNetworkVlanIDs[0]
					nicUntaggedVlanGroup, err := common.MatchVlanToGroup(nbi, vc.Networks.Vid2Name[vidID], vc.VlanGroupRelations)
					if err != nil {
						return fmt.Errorf("vlan group: %s", err)
					}
					intUntaggedVlan = nbi.VlansIndexByVlanGroupIDAndVID[nicUntaggedVlanGroup.ID][vidID]
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
			nbVMInterface, err := nbi.AddVMInterface(&objects.VMInterface{
				NetboxObject: objects.NetboxObject{
					Tags:        vc.Config.SourceTags,
					Description: intDescription,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: vc.SourceConfig.Name,
					},
				},
				VM:           netboxVM,
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

			// Setup Primary ipv4 address
			var nbPrimaryIPv4 *objects.IPAddress
			if primaryIPv4Address == "" {
				// Fallback mechanism, we choose the first ipv4 address on the interface
				if len(nicIPv4Addresses) > 0 {
					primaryIPv4Address = nicIPv4Addresses[0]
				}
			}
			if primaryIPv4Address != "" {
				nbPrimaryIPv4, err = nbi.AddIPAddress(&objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: vc.Config.SourceTags,
						CustomFields: map[string]string{
							constants.CustomFieldSourceName: vc.SourceConfig.Name,
						},
					},
					Address:            primaryIPv4Address,
					DNSName:            utils.ReverseLookup(primaryIPv4Address),
					AssignedObjectType: objects.AssignedObjectTypeVMInterface,
					AssignedObjectID:   nbVMInterface.ID,
				})
				if err != nil {
					vc.Logger.Errorf("adding ipv4 address: %s", err)
				}
			}

			// Setup Primary ipv6 address
			var nbPrimaryIPv6 *objects.IPAddress
			if primaryIPv6Address != "" {
				// Fallback mechanism, we choose the first ipv6 address on the interface
				if len(nicIPv6Addresses) > 0 {
					primaryIPv6Address = nicIPv6Addresses[0]
				}
			}
			if primaryIPv6Address != "" {
				nbPrimaryIPv6, err = nbi.AddIPAddress(&objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: vc.Config.SourceTags,
						CustomFields: map[string]string{
							constants.CustomFieldSourceName: vc.SourceConfig.Name,
						},
					},
					Address:            primaryIPv6Address,
					DNSName:            utils.ReverseLookup(primaryIPv6Address),
					AssignedObjectType: objects.AssignedObjectTypeVMInterface,
					AssignedObjectID:   nbVMInterface.ID,
				})
				if err != nil {
					vc.Logger.Errorf("adding ipv6 address: %s", err)
				}
			}

			// Update the vms with primary addresses
			if nbPrimaryIPv4 != nil && (netboxVM.PrimaryIPv4 == nil || nbPrimaryIPv4.Address != netboxVM.PrimaryIPv4.Address) || nbPrimaryIPv6 != nil && (netboxVM.PrimaryIPv6 == nil || nbPrimaryIPv6.Address != netboxVM.PrimaryIPv6.Address) {
				// Shallow copy netboxVm to newNetboxVM
				newNetboxVM := *netboxVM
				newNetboxVM.PrimaryIPv4 = nbPrimaryIPv4
				newNetboxVM.PrimaryIPv6 = nbPrimaryIPv6
				_, err = nbi.AddVM(&newNetboxVM)
				if err != nil {
					vc.Logger.Warningf("adding vm: %s", err)
				}
			}
		}
	}

	return nil
}
