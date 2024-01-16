package vmware

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/vmware/govmomi/vim25/mo"
)

func (vc *VmwareSource) syncDatacenters(nbi *inventory.NetBoxInventory) error {
	for _, dc := range vc.DataCenters {

		nbClusterGroup := &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{
				Description: fmt.Sprintf("Datacenter from source %s", vc.SourceConfig.Hostname),
				Tags:        vc.CommonConfig.SourceTags,
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

func (vc *VmwareSource) syncClusters(nbi *inventory.NetBoxInventory) error {
	clusterType := &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: vc.SourceTags,
		},
		Name: "Vmware ESXi",
		Slug: utils.Slugify("Vmware ESXi"),
	}
	clusterType, err := nbi.AddClusterType(clusterType)
	if err != nil {
		return fmt.Errorf("failed to add vmware ClusterType: %v", err)
	}
	// Then sync oVirt Clusters as NetBoxClusters
	for clusterId := range vc.Clusters {

		clusterName := clusterId

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
				Tags: vc.SourceTags,
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

// Host in oVirt is a represented as device in netbox with a
// custom role Server
func (vc *VmwareSource) syncHosts(nbi *inventory.NetBoxInventory) error {
	for hostId, host := range vc.Hosts {
		hostName := host.Name
		hostCluster := nbi.ClustersIndexByName[vc.Clusters[vc.Host2Cluster[hostId]].Name]

		var hostSite *objects.Site
		if vc.HostSiteRelations != nil {
			match, err := utils.MatchStringToValue(hostName, vc.HostSiteRelations)
			if err != nil {
				return fmt.Errorf("error occurred when matching vmware host %s to a Netbox site: %v", hostName, err)
			}
			if match != "" {
				if _, ok := nbi.SitesIndexByName[match]; !ok {
					return fmt.Errorf("failed to match vmware host %s to a Netbox site: %v. Site with this name doesn't exist", hostName, match)
				}
				hostSite = nbi.SitesIndexByName[match]
			}
		}
		var hostTenant *objects.Tenant
		if vc.HostTenantRelations != nil {
			match, err := utils.MatchStringToValue(hostName, vc.HostTenantRelations)
			if err != nil {
				return fmt.Errorf("error occurred when matching vmware host %s to a Netbox tenant: %v", hostName, err)
			}
			if match != "" {
				if _, ok := nbi.TenantsIndexByName[match]; !ok {
					return fmt.Errorf("failed to match vmware host %s to a Netbox tenant: %v. Tenant with this name doesn't exist", hostName, match)
				}
				hostTenant = nbi.TenantsIndexByName[match]
			}
		}

		var err error
		hostAssetTag := host.Summary.Hardware.Uuid
		hostModel := host.Summary.Hardware.Model

		var hostSerialNumber string
		// find serial number from  host summary.hardware.OtherIdentifyingInfo
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
			return fmt.Errorf("failed adding oVirt Manufacturer %v with error: %s", hostManufacturer, err)
		}

		var hostDeviceType *objects.DeviceType
		hostDeviceType, err = nbi.AddDeviceType(&objects.DeviceType{
			NetboxObject: objects.NetboxObject{
				Tags: vc.SourceTags,
			},
			Manufacturer: hostManufacturer,
			Model:        hostModel,
			Slug:         utils.Slugify(hostModel),
		})
		if err != nil {
			return fmt.Errorf("failed adding oVirt DeviceType %v with error: %s", hostDeviceType, err)
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
			return fmt.Errorf("failed adding oVirt Platform %v with error: %s", hostPlatform, err)
		}

		hostCpuCores := host.Summary.Hardware.NumCpuCores
		hostMemGB := host.Summary.Hardware.MemorySize / 1024 / 1024 / 1024

		nbHost := &objects.Device{
			NetboxObject: objects.NetboxObject{Tags: vc.SourceTags},
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
			CustomFields: map[string]string{
				"host_cpu_cores": fmt.Sprintf("%d", hostCpuCores),
				"host_memory":    fmt.Sprintf("%d GB", hostMemGB),
			},
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

func (vc *VmwareSource) syncHostNics(nbi *inventory.NetBoxInventory, vcHost *mo.HostSystem, nbHost *objects.Device) error {

	// Sync host's physical interfaces
	err := vc.syncHostPhysicalNics(nbi, vcHost, nbHost)
	if err != nil {
		return fmt.Errorf("physicial interfaces sync: %s", err)
	}

	// Sync host's virtual interfaces
	err = vc.syncHostVirtualNics(nbi, vcHost, nbHost)
	if err != nil {
		return fmt.Errorf("virtual interfaces sync: %s", err)
	}

	return nil
}

func (vc *VmwareSource) syncHostPhysicalNics(nbi *inventory.NetBoxInventory, vcHost *mo.HostSystem, nbHost *objects.Device) error {

	// Collect data from physical interfaces
	for _, pnic := range vcHost.Config.Network.Pnic {
		pnicName := pnic.Device
		pnicLinkSpeed := pnic.LinkSpeed.SpeedMb
		var pnicMode *objects.InterfaceMode
		var pnicMtu int
		var taggedVlanList []*objects.Vlan // when mode="tagged"
		if pnicLinkSpeed == 0 {
			pnicLinkSpeed = pnic.Spec.LinkSpeed.SpeedMb
			if pnicLinkSpeed == 0 {
				pnicLinkSpeed = pnic.ValidLinkSpecification[0].SpeedMb
			}
		}
		var pnicDescription string
		if pnicLinkSpeed > 1000 {
			pnicDescription = fmt.Sprintf("%dGB/s", pnicLinkSpeed/1000)
		} else {
			pnicDescription = fmt.Sprintf("%dMB/s", pnicLinkSpeed)
		}
		pnicDescription += " pNIC"

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
		pnicVlans := []*objects.Vlan{}
		for portgroupName, portgroupData := range vc.Networks.HostPortgroups[nbHost.Name] {
			if slices.Contains(portgroupData.nics, pnicName) {
				pnicVlans = append(pnicVlans, &objects.Vlan{
					Name: portgroupName,
					Vid:  portgroupData.vlanId,
				})
			}
		}

		// Determine interface mode for non VM traffic NIC, from vlans data
		if len(pnicVlans) > 0 {
			vlanIdSet := map[int]bool{} // set of vlans
			for _, pnicVlan := range pnicVlans {
				vlanIdSet[pnicVlan.Vid] = true
			}
			if len(vlanIdSet) == 1 && vlanIdSet[0] {
				pnicMode = &objects.InterfaceModeAccess
			} else if vlanIdSet[4095] {
				pnicMode = &objects.InterfaceModeTaggedAll
			} else {
				pnicMode = &objects.InterfaceModeTagged
			}
			taggedVlanList = []*objects.Vlan{}
			if pnicMode == &objects.InterfaceModeTagged {
				for _, pnicVlan := range pnicVlans {
					if pnicVlan.Vid == 0 {
						continue
					}
					taggedVlanList = append(taggedVlanList, pnicVlan)
				}
			}
		}

		// After collecting all of the data add interface to nbi
		_, err := nbi.AddInterface(&objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        vc.SourceTags,
				Description: pnicDescription,
			},
			Device:      nbHost,
			Name:        pnicName,
			Status:      true,
			Type:        &objects.OtherInterfaceType, //  TODO: Get type from link speed
			Speed:       objects.InterfaceSpeed(pnicLinkSpeed),
			MTU:         pnicMtu,
			Mode:        pnicMode,
			TaggedVlans: taggedVlanList,
		})
		if err != nil {
			return fmt.Errorf("failed adding physical interface: %s", err)
		}
	}
	return nil
}

func (vc *VmwareSource) syncHostVirtualNics(nbi *inventory.NetBoxInventory, vcHost *mo.HostSystem, nbHost *objects.Device) error {
	// Collect data over all virtual itnerfaces
	for _, vnic := range vcHost.Config.Network.Vnic {
		vnicName := vnic.Device
		vnicPortgroupData, vnicPortgroupDataOk := vc.Networks.HostPortgroups[vcHost.Name][vnic.Portgroup]
		vnicDvPortgroupKey := vnic.Spec.DistributedVirtualPort.PortgroupKey
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
			vnicUntaggedVlan = &objects.Vlan{
				Name:   fmt.Sprintf("ESXi %s (ID: %d) (%s)", vnic.Portgroup, vnicPortgroupVlanId, nbHost.Site.Name),
				Vid:    vnicPortgroupVlanId,
				Tenant: nbHost.Tenant,
			}
		} else if vnicDvPortgroupData != nil {
			for _, vnicDvPortgroupDataVlanId := range vnicDvPortgroupVlanIds {
				if vnicMode != &objects.InterfaceModeTagged {
					break
				}
				if vnicDvPortgroupDataVlanId == 0 {
					continue
				}
				vnicTaggedVlans = append(vnicTaggedVlans, &objects.Vlan{
					Name:   fmt.Sprintf("%s-%d", vnicDvPortgroupData.Name, vnicDvPortgroupDataVlanId),
					Vid:    vnicDvPortgroupDataVlanId,
					Tenant: nbHost.Tenant,
				})
			}
		}

		nbVnic, err := nbi.AddInterface(&objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        vc.SourceTags,
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
		})
		if err != nil {
			return err
		}

		// Get IPv4 address for this vnic. TODO: filter
		ipv4_address := vnic.Spec.Ip.IpAddress
		ipv4_mask := vnic.Spec.Ip.SubnetMask
		ipv4_dns := utils.ReverseLookup(ipv4_address)
		_, err = nbi.AddIPAddress(&objects.IPAddress{
			NetboxObject: objects.NetboxObject{
				Tags: vc.SourceTags,
			},
			Address:            fmt.Sprintf("%s/%s", ipv4_address, ipv4_mask),
			Status:             &objects.IPAddressStatusActive, // TODO
			DNSName:            ipv4_dns,
			Tenant:             nbHost.Tenant,
			AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
			AssignedObjectId:   nbVnic.Id,
		})
		if err != nil {
			return err
		}

		for _, ipv6_entry := range vnic.Spec.Ip.IpV6Config.IpV6Address {
			ipv6_address := ipv6_entry.IpAddress
			ipv6_mask := ipv6_entry.PrefixLength
			// TODO: Filter out ipv6 addresses
			_, err = nbi.AddIPAddress(&objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: vc.SourceTags,
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
	return nil
}

func (vc *VmwareSource) syncVms(nbi *inventory.NetBoxInventory) error {
	// for vmId, vm := range vc.Vms {
	// 	// VM name, which is used as unique identifier for VMs in Netbox
	// 	vmName, exists := vm.Name()
	// 	if !exists {
	// 		vc.Logger.Warning("name for oVirt vm with id ", vmId, " is empty. VM has to have unique name to be synced to netbox. Skipping...")
	// 	}

	// 	// VM's Cluster
	// 	var vmCluster *objects.Cluster
	// 	cluster, exists := vm.Cluster()
	// 	if exists {
	// 		if _, ok := vc.Clusters[cluster.MustId()]; ok {
	// 			vmCluster = nbi.ClustersIndexByName[vc.Clusters[cluster.MustId()].MustName()]
	// 		}
	// 	}

	// 	// Get VM's site,tenant and platform from cluster
	// 	var vmTenantGroup *objects.TenantGroup
	// 	var vmTenant *objects.Tenant
	// 	var vmSite *objects.Site
	// 	if vmCluster != nil {
	// 		vmTenantGroup = vmCluster.TenantGroup
	// 		vmTenant = vmCluster.Tenant
	// 		vmSite = vmCluster.Site
	// 	}

	// 	// VM's Status
	// 	var vmStatus *objects.VMStatus
	// 	status, exists := vm.Status()
	// 	if exists {
	// 		switch status {
	// 		case ovirtsdk4.VMSTATUS_UP:
	// 			vmStatus = &objects.VMStatusActive
	// 		default:
	// 			vmStatus = &objects.VMStatusOffline
	// 		}
	// 	}

	// 	// VM's Host Device (server)
	// 	var vmHostDevice *objects.Device
	// 	host, exists := vm.Host()
	// 	if exists {
	// 		if _, ok := vc.Hosts[host.MustId()]; ok {
	// 			vmHostDevice = nbi.DevicesIndexByUuid[vc.Hosts[host.MustId()].MustHardwareInformation().MustUuid()]
	// 		}
	// 	}

	// 	// vmVCPUs
	// 	var vmVCPUs float32
	// 	if cpuData, exists := vm.Cpu(); exists {
	// 		if cpuTopology, exists := cpuData.Topology(); exists {
	// 			if cores, exists := cpuTopology.Cores(); exists {
	// 				vmVCPUs = float32(cores)
	// 			}
	// 		}
	// 	}

	// 	// Memory
	// 	var vmMemorySizeBytes int64
	// 	if memory, exists := vm.Memory(); exists {
	// 		vmMemorySizeBytes = memory
	// 	}

	// 	// Disks
	// 	var vmDiskSizeBytes int64
	// 	if diskAttachment, exists := vm.DiskAttachments(); exists {
	// 		for _, diskAttachment := range diskAttachment.Slice() {
	// 			if ovirtDisk, exists := diskAttachment.Disk(); exists {
	// 				disk := vc.Disks[ovirtDisk.MustId()]
	// 				if provisionedDiskSize, exists := disk.ProvisionedSize(); exists {
	// 					vmDiskSizeBytes += provisionedDiskSize
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// VM's comments
	// 	var vmComments string
	// 	if comments, exists := vm.Comment(); exists {
	// 		vmComments = comments
	// 	}

	// 	// VM's Platform
	// 	var vmPlatform *objects.Platform
	// 	vmOsType := "Generic OS"
	// 	vmOsVersion := "Generic Version"
	// 	if guestOs, exists := vm.GuestOperatingSystem(); exists {
	// 		if guestOsType, exists := guestOs.Distribution(); exists {
	// 			vmOsType = guestOsType
	// 		}
	// 		if guestOsKernel, exists := guestOs.Kernel(); exists {
	// 			if guestOsVersion, exists := guestOsKernel.Version(); exists {
	// 				if osFullVersion, exists := guestOsVersion.FullVersion(); exists {
	// 					vmOsVersion = osFullVersion
	// 				}
	// 			}
	// 		}
	// 	} else {
	// 		if os, exists := vm.Os(); exists {
	// 			if ovirtOsType, exists := os.Type(); exists {
	// 				vmOsType = ovirtOsType
	// 			}
	// 			if ovirtOsVersion, exists := os.Version(); exists {
	// 				if osFullVersion, exists := ovirtOsVersion.FullVersion(); exists {
	// 					vmOsVersion = osFullVersion
	// 				}
	// 			}
	// 		}
	// 	}
	// 	platformName := utils.GeneratePlatformName(vmOsType, vmOsVersion)
	// 	vmPlatform, err := nbi.AddPlatform(&objects.Platform{
	// 		Name: platformName,
	// 		Slug: utils.Slugify(platformName),
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("failed adding oVirt vm's Platform %v with error: %s", vmPlatform, err)
	// 	}

	// 	newVM, err := nbi.AddVM(&objects.VM{
	// 		NetboxObject: objects.NetboxObject{
	// 			Tags: vc.SourceTags,
	// 		},
	// 		Name:        vmName,
	// 		Cluster:     vmCluster,
	// 		Site:        vmSite,
	// 		Tenant:      vmTenant,
	// 		TenantGroup: vmTenantGroup,
	// 		Status:      vmStatus,
	// 		Host:        vmHostDevice,
	// 		Platform:    vmPlatform,
	// 		Comments:    vmComments,
	// 		VCPUs:       vmVCPUs,
	// 		Memory:      int(vmMemorySizeBytes / 1024 / 1024),      // MBs
	// 		Disk:        int(vmDiskSizeBytes / 1024 / 1024 / 1024), // GBs
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("failed to sync oVirt vm: %v", err)
	// 	}

	// 	err = vc.syncVmInterfaces(nbi, vm, newVM)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to sync oVirt vm's interfaces: %v", err)
	// 	}
	// }

	return nil
}

// Syncs VM's interfaces to Netbox
// func (vc *VmwareSource) syncVmInterfaces(nbi *inventory.NetBoxInventory, ovirtVm *ovirtsdk4.Vm, netboxVm *objects.VM) error {
// 	// 	var vmPrimaryIpv4 *objects.IPAddress
// 	// 	var vmPrimaryIpv6 *objects.IPAddress
// 	// 	if reportedDevices, exist := ovirtVm.ReportedDevices(); exist {
// 	// 		for _, reportedDevice := range reportedDevices.Slice() {
// 	// 			if reportedDeviceType, exist := reportedDevice.Type(); exist {
// 	// 				if reportedDeviceType == "network" {
// 	// 					// We add interface to the list
// 	// 					var vmInterface *objects.VMInterface
// 	// 					var err error
// 	// 					if reportedDeviceName, exists := reportedDevice.Name(); exists {
// 	// 						vmInterface, err = nbi.AddVMInterface(&objects.VMInterface{
// 	// 							NetboxObject: objects.NetboxObject{
// 	// 								Tags:        vc.SourceTags,
// 	// 								Description: reportedDevice.MustDescription(),
// 	// 							},
// 	// 							VM:         netboxVm,
// 	// 							Name:       reportedDeviceName,
// 	// 							MACAddress: strings.ToUpper(reportedDevice.MustMac().MustAddress()),
// 	// 						})
// 	// 						if err != nil {
// 	// 							return fmt.Errorf("failed to sync oVirt vm's interface %s: %v", reportedDeviceName, err)
// 	// 						}
// 	// 					} else {
// 	// 						vc.Logger.Warning("name for oVirt vm's reported device is empty. Skipping...")
// 	// 						continue
// 	// 					}

// 	// 					if reportedDeviceIps, exist := reportedDevice.Ips(); exist {
// 	// 						for _, ip := range reportedDeviceIps.Slice() {
// 	// 							if ipAddress, exists := ip.Address(); exists {
// 	// 								if ipVersion, exists := ip.Version(); exists {

// 	// 									// Filter IPs, we won't sync IPs from specific interfaces
// 	// 									// like docker, flannel, calico, etc. interfaces
// 	// 									valid, err := utils.IsVMInterfaceNameValid(vmInterface.Name)
// 	// 									if err != nil {
// 	// 										return fmt.Errorf("failed to match oVirt vm's interface %s to a Netbox interface filter: %v", vmInterface.Name, err)
// 	// 									}
// 	// 									if !valid {
// 	// 										continue
// 	// 									}

// 	// 									// Try to do reverse lookup of IP to get DNS name
// 	// 									hostname := utils.ReverseLookup(ipAddress)

// 	// 									// Set default mask
// 	// 									var ipMask string
// 	// 									if netMask, exists := ip.Netmask(); exists {
// 	// 										ipMask = fmt.Sprintf("/%s", netMask)
// 	// 									} else {
// 	// 										switch ipVersion {
// 	// 										case "v4":
// 	// 											ipMask = "/32"
// 	// 										case "v6":
// 	// 											ipMask = "/128"
// 	// 										}
// 	// 									}

// 	// 									ipAddress, err := nbi.AddIPAddress(&objects.IPAddress{
// 	// 										NetboxObject: objects.NetboxObject{
// 	// 											Tags: vc.SourceTags,
// 	// 										},
// 	// 										Address:            ipAddress + ipMask,
// 	// 										Tenant:             netboxVm.Tenant,
// 	// 										Status:             &objects.IPAddressStatusActive,
// 	// 										DNSName:            hostname,
// 	// 										AssignedObjectType: objects.AssignedObjectTypeVMInterface,
// 	// 										AssignedObjectId:   vmInterface.Id,
// 	// 									})

// 	// 									if err != nil {
// 	// 										// TODO: return should be here, commented just for now
// 	// 										// return fmt.Errorf("failed to sync oVirt vm's interface %s ip %s: %v", vmInterface, ip.MustAddress(), err)
// 	// 										vc.Logger.Error(fmt.Sprintf("failed to sync oVirt vm's interface %s ip %s: %v", vmInterface, ip.MustAddress(), err))

// 	// 									}

// 	// 									// TODO: criteria to determine if reported device is primary IP
// 	// 									switch ipVersion {
// 	// 									case "v4":
// 	// 										if vmPrimaryIpv4 == nil {
// 	// 											vmPrimaryIpv4 = ipAddress
// 	// 										}
// 	// 									case "v6":
// 	// 										if vmPrimaryIpv6 == nil {
// 	// 											vmPrimaryIpv6 = ipAddress
// 	// 										}
// 	// 									}
// 	// 								}
// 	// 							}
// 	// 						}
// 	// 					}
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 	}
// 	// 	// Update netboxVM with primary IPs
// 	// 	// TODO: determine which ip is primary ipv4 and which is primary ipv6
// 	// 	// TODO: then assign them to netboxVM
// 	// 	// if vmPrimaryIpv4 != nil && (netboxVm.PrimaryIPv4 == nil || vmPrimaryIpv4.Address != netboxVm.PrimaryIPv4.Address) {
// 	// 	// 	netboxVm.PrimaryIPv4 = vmPrimaryIpv4
// 	// 	// 	if _, err := nbi.AddVM(netboxVm); err != nil {
// 	// 	// 		return fmt.Errorf("failed to sync oVirt vm's primary ipv4: %v", err)
// 	// 	// 	}
// 	// 	// }
// 	// 	// if vmPrimaryIpv6 != nil && (netboxVm.PrimaryIPv6 == nil || vmPrimaryIpv6.Address != netboxVm.PrimaryIPv6.Address) {
// 	// 	// 	netboxVm.PrimaryIPv6 = vmPrimaryIpv6
// 	// 	// 	if _, err := nbi.AddVM(netboxVm); err != nil {
// 	// 	// 		return fmt.Errorf("failed to sync oVirt vm's primary ipv6: %v", err)
// 	// 	// 	}
// 	// 	// }

// 	return nil
// }

func (vc *VmwareSource) syncNetworks(nbi *inventory.NetBoxInventory) error {
	vc.Logger.Info("Syncing networks...")
	for _, dvpg := range vc.Networks.DistributedVirtualPortgroups {
		// TODO: currently we are syncing only vlans
		if len(dvpg.VlanIds) == 1 && len(dvpg.VlanIdRanges) == 0 {
			_, err := nbi.AddVlan(&objects.Vlan{
				NetboxObject: objects.NetboxObject{
					Tags: vc.SourceTags,
				},
				Name:   dvpg.Name,
				Vid:    dvpg.VlanIds[0],
				Status: &objects.VlanStatusActive,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
