package ovirt

import (
	"fmt"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// Syncs networks received from oVirt API to the netbox.
func (o *OVirtSource) syncNetworks(nbi *inventory.NetboxInventory) error {
	for _, network := range o.Networks.OVirtNetworks {
		name, exists := network.Name()
		if !exists {
			return fmt.Errorf("network %v has no name", network)
		}
		description, _ := network.Description()
		// TODO: handle other networks
		if networkVlan, exists := network.Vlan(); exists {
			// Get vlanGroup from relation
			vlanGroup, err := common.MatchVlanToGroup(nbi, name, o.VlanGroupRelations)
			if err != nil {
				return err
			}
			// Get tenant from relation
			vlanTenant, err := common.MatchVlanToTenant(nbi, name, o.VlanTenantRelations)
			if err != nil {
				return err
			}
			if networkVlanID, exists := networkVlan.Id(); exists {
				_, err := nbi.AddVlan(&objects.Vlan{
					NetboxObject: objects.NetboxObject{
						Description: description,
						Tags:        o.Config.SourceTags,
						CustomFields: map[string]string{
							constants.CustomFieldSourceName: o.SourceConfig.Name,
						},
					},
					Name:     name,
					Group:    vlanGroup,
					Vid:      int(networkVlanID),
					Status:   &objects.VlanStatusActive,
					Tenant:   vlanTenant,
					Comments: network.MustComment(),
				})
				if err != nil {
					return fmt.Errorf("adding vlan: %v", err)
				}
			}
		}
	}
	return nil
}

func (o *OVirtSource) syncDatacenters(nbi *inventory.NetboxInventory) error {
	// First sync oVirt DataCenters as NetBoxClusterGroups
	for _, datacenter := range o.DataCenters {
		name, exists := datacenter.Name()
		if !exists {
			return fmt.Errorf("failed to get name for oVirt datacenter %s", name)
		}
		description, _ := datacenter.Description()

		nbClusterGroup := &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{
				Description: description,
				Tags:        o.Config.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: o.SourceConfig.Name,
				},
			},
			Name: name,
			Slug: utils.Slugify(name),
		}
		_, err := nbi.AddClusterGroup(nbClusterGroup)
		if err != nil {
			return fmt.Errorf("failed to add oVirt data center %s as Netbox cluster group: %v", name, err)
		}
	}
	return nil
}

func (o *OVirtSource) syncClusters(nbi *inventory.NetboxInventory) error {
	clusterType := &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: o.Config.SourceTags,
			CustomFields: map[string]string{
				constants.CustomFieldSourceName: o.SourceConfig.Name,
			},
		},
		Name: "oVirt",
		Slug: "ovirt",
	}
	clusterType, err := nbi.AddClusterType(clusterType)
	if err != nil {
		return fmt.Errorf("failed to add oVirt cluster type: %v", err)
	}
	// Then sync oVirt Clusters as NetBoxClusters
	for _, cluster := range o.Clusters {
		clusterName, exists := cluster.Name()
		if !exists {
			return fmt.Errorf("failed to get name for oVirt cluster %s", clusterName)
		}
		description, exists := cluster.Description()
		if !exists {
			o.Logger.Warning("description for oVirt cluster ", clusterName, " is empty.")
		}
		var clusterGroup *objects.ClusterGroup
		var clusterGroupName string
		if _, ok := o.DataCenters[cluster.MustDataCenter().MustId()]; ok {
			clusterGroupName = o.DataCenters[cluster.MustDataCenter().MustId()].MustName()
		} else {
			o.Logger.Warning("failed to get datacenter for oVirt cluster ", clusterName)
		}
		if clusterGroupName != "" {
			clusterGroup = nbi.ClusterGroupsIndexByName[clusterGroupName]
		}
		var clusterSite *objects.Site
		if o.ClusterSiteRelations != nil {
			match, err := utils.MatchStringToValue(clusterName, o.ClusterSiteRelations)
			if err != nil {
				return fmt.Errorf("failed to match oVirt cluster %s to a Netbox site: %v", clusterName, err)
			}
			if match != "" {
				if _, ok := nbi.SitesIndexByName[match]; !ok {
					return fmt.Errorf("failed to match oVirt cluster %s to a Netbox site: %v. Site with this name doesn't exist", clusterName, match)
				}
				clusterSite = nbi.SitesIndexByName[match]
			}
		}
		var clusterTenant *objects.Tenant
		if o.ClusterTenantRelations != nil {
			match, err := utils.MatchStringToValue(clusterName, o.ClusterTenantRelations)
			if err != nil {
				return fmt.Errorf("error occurred when matching oVirt cluster %s to a Netbox tenant: %v", clusterName, err)
			}
			if match != "" {
				if _, ok := nbi.TenantsIndexByName[match]; !ok {
					return fmt.Errorf("failed to match oVirt cluster %s to a Netbox tenant: %v. Tenant with this name doesn't exist", clusterName, match)
				}
				clusterTenant = nbi.TenantsIndexByName[match]
			}
		}

		nbCluster := &objects.Cluster{
			NetboxObject: objects.NetboxObject{
				Description: description,
				Tags:        o.Config.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: o.SourceConfig.Name,
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
			return fmt.Errorf("failed to add oVirt cluster %s as Netbox cluster: %v", clusterName, err)
		}
	}
	return nil
}

// Host in oVirt is a represented as device in netbox with a
// custom role Server.
func (o *OVirtSource) syncHosts(nbi *inventory.NetboxInventory) error {
	for hostID, host := range o.Hosts {
		hostName, exists := host.Name()
		if !exists {
			o.Logger.Warningf("name of host with id=%s is empty", hostID)
		}
		hostCluster := nbi.ClustersIndexByName[o.Clusters[host.MustCluster().MustId()].MustName()]

		hostSite, err := common.MatchHostToSite(nbi, hostName, o.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("hostSite: %s", err)
		}
		hostTenant, err := common.MatchHostToTenant(nbi, hostName, o.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("hostTenant: %s", err)
		}

		var hostSerialNumber, manufacturerName, hostAssetTag, hostModel string
		hwInfo, exists := host.HardwareInformation()
		if exists {
			hostAssetTag, exists = hwInfo.Uuid()
			if !exists {
				o.Logger.Warning("Uuid (asset tag) for oVirt host ", hostName, " is empty. Can't identify it, so it will be skipped...")
				continue
			}
			hostSerialNumber, exists = hwInfo.SerialNumber()
			if !exists {
				o.Logger.Warning("Serial number for oVirt host ", hostName, " is empty.")
			}
			manufacturerName, _ = hwInfo.Manufacturer()
			manufacturerName, err = utils.MatchStringToValue(manufacturerName, objects.ManufacturerMap)
			if err != nil {
				return fmt.Errorf("error occurred when matching oVirt host %s to a Netbox manufacturer: %v", hostName, err)
			}

			hostModel, exists = hwInfo.ProductName()
			if !exists {
				hostModel = constants.DefaultModel // Model is also required for adding device type into netbox
			}
		} else {
			o.Logger.Warning("Hardware information for oVirt host ", hostName, " is empty, it can't be identified so it will be skipped.")
			continue
		}

		var hostManufacturer *objects.Manufacturer
		if manufacturerName == "" {
			manufacturerName = constants.DefaultManufacturer
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
			Manufacturer: hostManufacturer,
			Model:        hostModel,
			Slug:         utils.Slugify(hostModel),
		})
		if err != nil {
			return fmt.Errorf("failed adding oVirt DeviceType %v with error: %s", hostDeviceType, err)
		}

		var hostStatus *objects.DeviceStatus
		ovirtStatus, exists := host.Status()
		if exists {
			switch ovirtStatus {
			case ovirtsdk4.HOSTSTATUS_UP:
				hostStatus = &objects.DeviceStatusActive
			default:
				hostStatus = &objects.DeviceStatusOffline
			}
		}

		var hostPlatform *objects.Platform
		var osType, osVersion string
		if os, exists := host.Os(); exists {
			if ovirtOsType, exists := os.Type(); exists {
				osType = ovirtOsType
			}
			if ovirtOsVersion, exists := os.Version(); exists {
				if osFullVersion, exists := ovirtOsVersion.FullVersion(); exists {
					osVersion = osFullVersion
				}
			}
		}
		platformName := utils.GeneratePlatformName(osType, osVersion)
		hostPlatform, err = nbi.AddPlatform(&objects.Platform{
			Name: platformName,
			Slug: utils.Slugify(platformName),
		})
		if err != nil {
			return fmt.Errorf("failed adding oVirt Platform %v with error: %s", hostPlatform, err)
		}

		var hostDescription string
		if description, exists := host.Description(); exists {
			hostDescription = description
		}

		var hostComment string
		if comment, exists := host.Comment(); exists {
			hostComment = comment
		}

		var hostCPUCores string
		if cpu, exists := host.Cpu(); exists {
			hostCPUCores, exists = cpu.Name()
			if !exists {
				o.Logger.Warning("oVirt hostCpuCores of ", hostName, " is empty.")
			}
		}

		mem, _ := host.Memory()
		mem /= (constants.KiB * constants.KiB * constants.KiB) // Value is in Bytes, we convert to GB

		nbHost := &objects.Device{
			NetboxObject: objects.NetboxObject{
				Description: hostDescription,
				Tags:        o.Config.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName:       o.SourceConfig.Name,
					constants.CustomFieldHostCPUCoresName: hostCPUCores,
					constants.CustomFieldHostMemoryName:   fmt.Sprintf("%d GB", mem),
				},
			},
			Name:         hostName,
			Status:       hostStatus,
			Platform:     hostPlatform,
			DeviceRole:   nbi.DeviceRolesIndexByName["Server"],
			Site:         hostSite,
			Tenant:       hostTenant,
			Cluster:      hostCluster,
			Comments:     hostComment,
			SerialNumber: hostSerialNumber,
			AssetTag:     hostAssetTag,
			DeviceType:   hostDeviceType,
		}
		nbHost, err = nbi.AddDevice(nbHost)
		if err != nil {
			return fmt.Errorf("failed to add oVirt host %s with error: %v", host.MustName(), err)
		}

		// We also need to sync nics separately, because nic is a separate object in netbox
		err = o.syncHostNics(nbi, host, nbHost)
		if err != nil {
			return fmt.Errorf("failed to sync oVirt host %s nics with error: %v", host.MustName(), err)
		}
	}

	return nil
}

func (o *OVirtSource) syncHostNics(nbi *inventory.NetboxInventory, ovirtHost *ovirtsdk4.Host, nbHost *objects.Device) error {
	if nics, exists := ovirtHost.Nics(); exists {
		master2slave := make(map[string][]string) // masterId: [slaveId1, slaveId2, ...]
		parent2child := make(map[string][]string) // parentId: [childId, ... ]
		processedNicsIDs := make(map[string]bool) // set of all nic ids that have already been processed

		nicID2nic := map[string]*objects.Interface{} // nicId: nic
		nicID2IPv4 := map[string]string{}            // nicId: ipAv4address/mask
		nicID2IPv6 := map[string]string{}            // nicId: ipv6Address/mask

		var hostIP string
		if hostAddress, exists := ovirtHost.Address(); exists {
			hostIP = utils.Lookup(hostAddress)
		}

		// First loop, we loop through all the nics and collect all the information
		err := o.collectHostNicsData(nbHost, nbi, nics, parent2child, master2slave, nicID2nic, processedNicsIDs, nicID2IPv4, nicID2IPv6)
		if err != nil {
			return fmt.Errorf("collect host nics data: %s", err)
		}

		// Second loop to add relations between interfaces (e.g. [eno1, eno2] -> bond1)
		for masterID, slavesIDs := range master2slave {
			var err error
			masterInterface := nicID2nic[masterID]
			if _, ok := processedNicsIDs[masterID]; ok {
				masterInterface, err = nbi.AddInterface(masterInterface)
				if err != nil {
					return fmt.Errorf("failed to add oVirt master interface %s with error: %v", masterInterface.Name, err)
				}
				delete(processedNicsIDs, masterID)
				nicID2nic[masterID] = masterInterface
			}
			for _, slaveID := range slavesIDs {
				slaveInterface := nicID2nic[slaveID]
				slaveInterface.LAG = masterInterface
				slaveInterface, err := nbi.AddInterface(slaveInterface)
				if err != nil {
					return fmt.Errorf("failed to add oVirt slave interface %s with error: %v", slaveInterface.Name, err)
				}
				delete(processedNicsIDs, slaveID)
				nicID2nic[slaveID] = slaveInterface
			}
		}

		// Third loop we connect children with parents (e.g. [bond1.605, bond1.604, bond1.603] -> bond1)
		for parent, children := range parent2child {
			parentInterface := nicID2nic[parent]
			if _, ok := processedNicsIDs[parent]; ok {
				parentInterface, err := nbi.AddInterface(parentInterface)
				if err != nil {
					return fmt.Errorf("failed to add oVirt parent interface %s with error: %v", parentInterface.Name, err)
				}
				nicID2nic[parent] = parentInterface
				delete(processedNicsIDs, parent)
			}
			for _, child := range children {
				childInterface := nicID2nic[child]
				childInterface.ParentInterface = parentInterface
				childInterface, err := nbi.AddInterface(childInterface)
				if err != nil {
					return fmt.Errorf("failed to add oVirt child interface %s with error: %v", childInterface.Name, err)
				}
				nicID2nic[child] = childInterface
				delete(processedNicsIDs, child)
			}
		}

		// Fourth loop we check if there are any nics that were not processed
		for nicID := range processedNicsIDs {
			nbNic, err := nbi.AddInterface(nicID2nic[nicID])
			if err != nil {
				return fmt.Errorf("failed to add oVirt interface %s with error: %v", nicID2nic[nicID].Name, err)
			}
			nicID2nic[nicID] = nbNic
		}

		// Fifth loop we add ip addresses to interfaces
		for nicID, ipv4 := range nicID2IPv4 {
			nbNic := nicID2nic[nicID]
			address := strings.Split(ipv4, "/")[0]
			nbIPAddress, err := nbi.AddIPAddress(&objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: o.Config.SourceTags,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: o.SourceConfig.Name,
					},
				},
				Address:            ipv4,
				Status:             &objects.IPAddressStatusActive, // TODO
				DNSName:            utils.ReverseLookup(address),
				AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
				AssignedObjectID:   nbNic.ID,
			})
			if err != nil {
				return fmt.Errorf("add ipv4 address: %s", err)
			}
			if address == hostIP {
				hostCopy := *nbHost
				hostCopy.PrimaryIPv4 = nbIPAddress
				_, err := nbi.AddDevice(&hostCopy)
				if err != nil {
					return fmt.Errorf("adding primary ipv4 address: %s", err)
				}
			}
		}
		for nicID, ipv6 := range nicID2IPv6 {
			nbNic := nicID2nic[nicID]
			address := strings.Split(ipv6, "/")[0]
			_, err := nbi.AddIPAddress(&objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: o.Config.SourceTags,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: o.SourceConfig.Name,
					},
				},
				Address:            ipv6,
				Status:             &objects.IPAddressStatusActive, // TODO
				DNSName:            utils.ReverseLookup(address),
				AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
				AssignedObjectID:   nbNic.ID,
			})
			if err != nil {
				return fmt.Errorf("add ipv6 address: %s", err)
			}
		}
	}
	return nil
}

func (o *OVirtSource) collectHostNicsData(nbHost *objects.Device, nbi *inventory.NetboxInventory, nics *ovirtsdk4.HostNicSlice, parent2child map[string][]string, master2slave map[string][]string, nicID2nic map[string]*objects.Interface, processedNicsIDs map[string]bool, nicID2IPv4 map[string]string, nicID2IPv6 map[string]string) error {
	for _, nic := range nics.Slice() {
		nicID, exists := nic.Id()
		if !exists {
			o.Logger.Warning("id for oVirt nic with id ", nicID, " is empty. This should not happen! Skipping...")
			continue
		}
		nicName, exists := nic.Name()
		if !exists {
			o.Logger.Warning("name for oVirt nic with id ", nicID, " is empty.")
		}
		// var nicType *objects.InterfaceType
		nicSpeedBips, exists := nic.Speed()
		if !exists {
			o.Logger.Warning("speed for oVirt nic with id ", nicID, " is empty.")
		}
		nicSpeedKbps := nicSpeedBips / constants.KB

		nicMtu, exists := nic.Mtu()
		if !exists {
			o.Logger.Warning("mtu for oVirt nic with id ", nicID, " is empty.")
		}

		nicComment, _ := nic.Comment()

		var nicEnabled bool
		ovirtNicStatus, exists := nic.Status()
		if exists {
			switch ovirtNicStatus {
			case ovirtsdk4.NICSTATUS_UP:
				nicEnabled = true
			default:
				nicEnabled = false
			}
		}

		// bridged, exists := nic.Bridged() // TODO: bridged interface
		// if exists {
		// 	if bridged {
		// 		// This interface is bridged
		// 		fmt.Printf("nic[%s] is bridged\n", nicName)
		// 	}
		// }

		// Determine nic type (virtual, physical, bond...)
		var nicType *objects.InterfaceType
		nicBaseInterface, exists := nic.BaseInterface()
		if exists {
			// This interface is a vlan bond. We treat is as a virtual interface
			nicType = &objects.VirtualInterfaceType
			parent2child[nicBaseInterface] = append(parent2child[nicBaseInterface], nicID)
		}

		nicBonding, exists := nic.Bonding()
		if exists {
			// Bond interface, we give it a type of LAG
			nicType = &objects.LAGInterfaceType
			slaves, exists := nicBonding.Slaves()
			if exists {
				for _, slave := range slaves.Slice() {
					master2slave[nicID] = append(master2slave[nicID], slave.MustId())
				}
			}
		}

		if nicType == nil {
			// This is a physical interface.
			nicType = objects.IfaceSpeed2IfaceType[objects.InterfaceSpeed(nicSpeedKbps)]
			if nicType == nil {
				nicType = &objects.OtherInterfaceType
			}
		}

		var nicVlan *objects.Vlan
		vlan, exists := nic.Vlan()
		if exists {
			vlanID, exists := vlan.Id()
			if exists {
				vlanName := o.Networks.Vid2Name[int(vlanID)]
				// Get vlanGroup from relation
				vlanGroup, err := common.MatchVlanToGroup(nbi, vlanName, o.VlanGroupRelations)
				if err != nil {
					return err
				}
				// Get vlan from inventory
				nicVlan = nbi.VlansIndexByVlanGroupIDAndVID[vlanGroup.ID][int(vlanID)]
			}
		}

		var nicTaggedVlans []*objects.Vlan
		if nicVlan != nil {
			nicTaggedVlans = []*objects.Vlan{nicVlan}
		}

		newInterface := &objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        o.Config.SourceTags,
				Description: nicComment,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName:   o.SourceConfig.Name,
					constants.CustomFieldSourceIDName: nicID,
				},
			},
			Device:      nbHost,
			Name:        nicName,
			Speed:       objects.InterfaceSpeed(nicSpeedKbps),
			Status:      nicEnabled,
			MTU:         int(nicMtu),
			Type:        nicType,
			TaggedVlans: nicTaggedVlans,
		}

		var err error
		// Extract ip info
		if nicIPv4, exists := nic.Ip(); exists {
			if nicAddress, exists := nicIPv4.Address(); exists {
				mask := 32
				if nicMask, exists := nicIPv4.Netmask(); exists {
					mask, err = utils.MaskToBits(nicMask)
					if err != nil {
						return fmt.Errorf("mask to bits: %s", err)
					}
				}
				ipv4Address := fmt.Sprintf("%s/%d", nicAddress, mask)
				nicID2IPv4[nicID] = ipv4Address
			}
		}
		if nicIPv6, exists := nic.Ipv6(); exists {
			if nicAddress, exists := nicIPv6.Address(); exists {
				mask := 128
				if nicMask, exists := nicIPv6.Netmask(); exists {
					mask, err = utils.MaskToBits(nicMask)
					if err != nil {
						return fmt.Errorf("mask to bits: %s", err)
					}
				}
				ipv6Address := fmt.Sprintf("%s/%d", nicAddress, mask)
				nicID2IPv6[nicID] = ipv6Address
			}
		}

		processedNicsIDs[nicID] = true
		nicID2nic[nicID] = newInterface
	}
	return nil
}

func (o *OVirtSource) syncVms(nbi *inventory.NetboxInventory) error {
	for vmID, ovirtVM := range o.Vms {
		collectedVM, err := o.extractVMData(nbi, vmID, ovirtVM)
		if err != nil {
			return err
		}

		nbVM, err := nbi.AddVM(collectedVM)
		if err != nil {
			return fmt.Errorf("failed to sync oVirt vm: %v", err)
		}

		err = o.syncVMInterfaces(nbi, ovirtVM, nbVM)
		if err != nil {
			return fmt.Errorf("failed to sync oVirt vm's interfaces: %v", err)
		}
	}

	return nil
}

func (o *OVirtSource) extractVMData(nbi *inventory.NetboxInventory, vmID string, vm *ovirtsdk4.Vm) (*objects.VM, error) {
	// VM name, which is used as unique identifier for VMs in Netbox
	vmName, exists := vm.Name()
	if !exists {
		o.Logger.Warning("name for oVirt vm with id ", vmID, " is empty. VM has to have unique name to be synced to netbox. Skipping...")
	}

	// VM's Cluster
	var vmCluster *objects.Cluster
	cluster, exists := vm.Cluster()
	if exists {
		if _, ok := o.Clusters[cluster.MustId()]; ok {
			vmCluster = nbi.ClustersIndexByName[o.Clusters[cluster.MustId()].MustName()]
		}
	}

	// Get VM's site,tenant and platform from cluster
	var vmTenantGroup *objects.TenantGroup
	var vmTenant *objects.Tenant
	var vmSite *objects.Site
	if vmCluster != nil {
		vmTenantGroup = vmCluster.TenantGroup
		vmTenant = vmCluster.Tenant
		vmSite = vmCluster.Site
	}

	// VM's Status
	var vmStatus *objects.VMStatus
	status, exists := vm.Status()
	if exists {
		switch status {
		case ovirtsdk4.VMSTATUS_UP:
			vmStatus = &objects.VMStatusActive
		default:
			vmStatus = &objects.VMStatusOffline
		}
	}

	// VM's Host Device (server)
	var vmHostDevice *objects.Device
	if host, exists := vm.Host(); exists {
		if oHost, ok := o.Hosts[host.MustId()]; ok {
			if oHostName, ok := oHost.Name(); ok {
				vmHostDevice = nbi.DevicesIndexByNameAndSiteID[oHostName][vmSite.ID]
			}
		}
	}

	// vmVCPUs
	var vmVCPUs float32
	if cpuData, exists := vm.Cpu(); exists {
		if cpuTopology, exists := cpuData.Topology(); exists {
			if cores, exists := cpuTopology.Cores(); exists {
				vmVCPUs = float32(cores)
			}
		}
	}

	// Memory
	var vmMemorySizeBytes int64
	if memory, exists := vm.Memory(); exists {
		vmMemorySizeBytes = memory
	}

	// Disks
	var vmDiskSizeBytes int64
	if diskAttachment, exists := vm.DiskAttachments(); exists {
		for _, diskAttachment := range diskAttachment.Slice() {
			if ovirtDisk, exists := diskAttachment.Disk(); exists {
				disk := o.Disks[ovirtDisk.MustId()]
				if provisionedDiskSize, exists := disk.ProvisionedSize(); exists {
					vmDiskSizeBytes += provisionedDiskSize
				}
			}
		}
	}

	// VM's comments
	var vmComments string
	if comments, exists := vm.Comment(); exists {
		vmComments = comments
	}

	// VM's Platform
	var vmPlatform *objects.Platform
	vmOsType := constants.DefaultOSName
	vmOsVersion := constants.DefaultOSVersion
	if guestOs, exists := vm.GuestOperatingSystem(); exists {
		if guestOsType, exists := guestOs.Distribution(); exists {
			vmOsType = guestOsType
		}
		if guestOsKernel, exists := guestOs.Kernel(); exists {
			if guestOsVersion, exists := guestOsKernel.Version(); exists {
				if osFullVersion, exists := guestOsVersion.FullVersion(); exists {
					vmOsVersion = osFullVersion
				}
			}
		}
	} else {
		if os, exists := vm.Os(); exists {
			if ovirtOsType, exists := os.Type(); exists {
				vmOsType = ovirtOsType
			}
			if ovirtOsVersion, exists := os.Version(); exists {
				if osFullVersion, exists := ovirtOsVersion.FullVersion(); exists {
					vmOsVersion = osFullVersion
				}
			}
		}
	}
	platformName := utils.GeneratePlatformName(vmOsType, vmOsVersion)
	vmPlatform, err := nbi.AddPlatform(&objects.Platform{
		Name: platformName,
		Slug: utils.Slugify(platformName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed adding oVirt vm's Platform %v with error: %s", vmPlatform, err)
	}

	return &objects.VM{
		NetboxObject: objects.NetboxObject{
			Tags: o.Config.SourceTags,
			CustomFields: map[string]string{
				constants.CustomFieldSourceName: o.SourceConfig.Name,
			},
		},
		Name:        vmName,
		Cluster:     vmCluster,
		Site:        vmSite,
		Tenant:      vmTenant,
		TenantGroup: vmTenantGroup,
		Status:      vmStatus,
		Host:        vmHostDevice,
		Platform:    vmPlatform,
		Comments:    vmComments,
		VCPUs:       vmVCPUs,
		Memory:      int(vmMemorySizeBytes / constants.KiB / constants.KiB),               // MBs
		Disk:        int(vmDiskSizeBytes / constants.KiB / constants.KiB / constants.KiB), // GBs
	}, nil
}

// Syncs VM's interfaces to Netbox.
func (o *OVirtSource) syncVMInterfaces(nbi *inventory.NetboxInventory, ovirtVM *ovirtsdk4.Vm, netboxVM *objects.VM) error {
	if reportedDevices, exist := ovirtVM.ReportedDevices(); exist {
		for _, reportedDevice := range reportedDevices.Slice() {
			if reportedDeviceType, exist := reportedDevice.Type(); exist {
				if reportedDeviceType == "network" {
					// We add interface to the list
					var vmInterface *objects.VMInterface
					var err error
					if reportedDeviceName, exists := reportedDevice.Name(); exists {
						vmInterface, err = nbi.AddVMInterface(&objects.VMInterface{
							NetboxObject: objects.NetboxObject{
								Tags:        o.Config.SourceTags,
								Description: reportedDevice.MustDescription(),
								CustomFields: map[string]string{
									constants.CustomFieldSourceName: o.SourceConfig.Name,
								},
							},
							VM:         netboxVM,
							Name:       reportedDeviceName,
							MACAddress: strings.ToUpper(reportedDevice.MustMac().MustAddress()),
							Enabled:    true, // TODO
						})
						if err != nil {
							return fmt.Errorf("failed to sync oVirt vm's interface %s: %v", reportedDeviceName, err)
						}
					} else {
						o.Logger.Warning("name for oVirt vm's reported device is empty. Skipping...")
						continue
					}

					if reportedDeviceIps, exist := reportedDevice.Ips(); exist {
						for _, ip := range reportedDeviceIps.Slice() {
							if ipAddress, exists := ip.Address(); exists {
								if ipVersion, exists := ip.Version(); exists {
									// Filter IPs, we won't sync IPs from specific interfaces
									// like docker, flannel, calico, etc. interfaces
									valid, err := utils.IsVMInterfaceNameValid(vmInterface.Name)
									if err != nil {
										return fmt.Errorf("failed to match oVirt vm's interface %s to a Netbox interface filter: %v", vmInterface.Name, err)
									}
									if !valid {
										continue
									}

									// Try to do reverse lookup of IP to get DNS name
									hostname := utils.ReverseLookup(ipAddress)

									// Set default mask
									var ipMask string
									if netMask, exists := ip.Netmask(); exists {
										ipMask = fmt.Sprintf("/%s", netMask)
									} else {
										switch ipVersion {
										case "v4":
											ipMask = "/32"
										case "v6":
											ipMask = "/128"
										}
									}

									newIPAddress, err := nbi.AddIPAddress(&objects.IPAddress{
										NetboxObject: objects.NetboxObject{
											Tags: o.Config.SourceTags,
											CustomFields: map[string]string{
												constants.CustomFieldSourceName: o.SourceConfig.Name,
											},
										},
										Address:            ipAddress + ipMask,
										Tenant:             netboxVM.Tenant,
										Status:             &objects.IPAddressStatusActive,
										DNSName:            hostname,
										AssignedObjectType: objects.AssignedObjectTypeVMInterface,
										AssignedObjectID:   vmInterface.ID,
									})

									if err != nil {
										// TODO: return should be here, commented just for now
										// return fmt.Errorf("failed to sync oVirt vm's interface %s ip %s: %v", vmInterface, ip.MustAddress(), err)
										o.Logger.Error(fmt.Sprintf("failed to sync oVirt vm's interface %s ip %s: %v", vmInterface, ip.MustAddress(), err))
									}

									// Check if ip is primary
									if ipVersion == "v4" {
										vmIP := utils.Lookup(netboxVM.Name)
										if vmIP != "" && vmIP == ipAddress || netboxVM.PrimaryIPv4 == nil {
											vmCopy := *netboxVM
											vmCopy.PrimaryIPv4 = newIPAddress
											_, err := nbi.AddVM(&vmCopy)
											if err != nil {
												return fmt.Errorf("adding primary ipv4 address: %s", err)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
