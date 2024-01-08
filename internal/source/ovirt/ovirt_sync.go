package ovirt

import (
	"fmt"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// Function that syncs all data from oVirt to Netbox
func (o *OVirtSource) Sync(nbi *inventory.NetBoxInventory) error {
	syncFunctions := []func(*inventory.NetBoxInventory) error{
		o.syncDatacenters,
		o.syncClusters,
		o.syncHosts,
		o.syncVms,
	}
	for _, syncFunc := range syncFunctions {
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OVirtSource) syncDatacenters(nbi *inventory.NetBoxInventory) error {
	// First sync oVirt DataCenters as NetBoxClusterGroups
	for _, datacenter := range o.DataCenters {
		name, exists := datacenter.Name()
		if !exists {
			return fmt.Errorf("failed to get name for oVirt datacenter %s", name)
		}
		description, _ := datacenter.Description()

		nbClusterGroup := &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{Description: description},
			Name:         name,
			Slug:         utils.Slugify(name),
		}
		_, err := nbi.AddClusterGroup(nbClusterGroup, o.SourceTags)
		if err != nil {
			return fmt.Errorf("failed to add oVirt data center %s as Netbox cluster group: %v", name, err)
		}
	}
	return nil
}

func (o *OVirtSource) syncClusters(nbi *inventory.NetBoxInventory) error {
	clusterType := &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: o.SourceTags,
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
		if _, ok := o.DataCenters[cluster.MustDataCenter().MustId()]; ok {

		} else {
			o.Logger.Warning("failed to get datacenter for oVirt cluster ", clusterName)
		}
		if dataCenter, ok := cluster.DataCenter(); ok {
			if dataCenterName, ok := dataCenter.Name(); ok {
				clusterGroup = nbi.ClusterGroupsIndexByName[dataCenterName]
			}
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
				Tags:        o.SourceTags,
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
// custom role Server
func (o *OVirtSource) syncHosts(nbi *inventory.NetBoxInventory) error {
	for hostId, host := range o.Hosts {
		hostName, exists := host.Name()
		if !exists {
			o.Logger.Warning("name for oVirt host with id ", hostId, " is empty.")
		}
		hostCluster := nbi.ClustersIndexByName[o.Clusters[host.MustCluster().MustId()].MustName()]

		var hostDescription string
		description, exists := host.Description()
		if exists {
			hostDescription = description
		}

		var hostSite *objects.Site
		if o.HostSiteRelations != nil {
			match, err := utils.MatchStringToValue(hostName, o.HostSiteRelations)
			if err != nil {
				return fmt.Errorf("error occurred when matching oVirt host %s to a Netbox site: %v", hostName, err)
			}
			if match != "" {
				if _, ok := nbi.SitesIndexByName[match]; !ok {
					return fmt.Errorf("failed to match oVirt host %s to a Netbox site: %v. Site with this name doesn't exist", hostName, match)
				}
				hostSite = nbi.SitesIndexByName[match]
			}
		}
		var hostTenant *objects.Tenant
		if o.HostTenantRelations != nil {
			match, err := utils.MatchStringToValue(hostName, o.HostTenantRelations)
			if err != nil {
				return fmt.Errorf("error occurred when matching oVirt host %s to a Netbox tenant: %v", hostName, err)
			}
			if match != "" {
				if _, ok := nbi.TenantsIndexByName[match]; !ok {
					return fmt.Errorf("failed to match oVirt host %s to a Netbox tenant: %v. Tenant with this name doesn't exist", hostName, match)
				}
				hostTenant = nbi.TenantsIndexByName[match]
			}
		}

		var hostSerialNumber, manufacturerName, hostAssetTag, hostModel string
		var err error
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
				hostModel = "Generic Model" // Model is also required for adding device type into netbox
			}
		} else {
			o.Logger.Warning("Hardware information for oVirt host ", hostName, " is empty, it can't be identified so it will be skipped.")
			continue
		}

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
		osType := "Generic OS"
		osVersion := "Generic Version"
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

		var hostComment string
		if comment, exists := host.Comment(); exists {
			hostComment = comment
		}

		var hostCpuCores string
		if cpu, exists := host.Cpu(); exists {
			hostCpuCores, exists = cpu.Name()
			if !exists {
				o.Logger.Warning("oVirt hostCpuCores of ", hostName, " is empty.")
			}
		}

		mem, _ := host.Memory()
		mem /= (1024 * 1024 * 1024) // Value is in Bytes, we convert to GB

		nbHost := &objects.Device{
			NetboxObject: objects.NetboxObject{Description: hostDescription, Tags: o.SourceTags},
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
			CustomFields: map[string]string{
				"host_cpu_cores": hostCpuCores,
				"host_memory":    fmt.Sprintf("%d GB", mem),
			},
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

func (o *OVirtSource) syncHostNics(nbi *inventory.NetBoxInventory, ovirtHost *ovirtsdk4.Host, nbHost *objects.Device) error {
	nics, exists := ovirtHost.Nics()
	master2slave := make(map[string][]string) // masterId: [slaveId1, slaveId2, ...]
	parent2child := make(map[string][]string) // parentId: [childId, ... ]
	processedNicsIds := make(map[string]bool)
	if exists {
		hostInterfaces := map[string]*objects.Interface{}

		// First loop through all nics
		for _, nic := range nics.Slice() {
			nicId, exists := nic.Id()
			if !exists {
				o.Logger.Warning("id for oVirt nic with id ", nicId, " is empty. This should not happen! Skipping...")
				continue
			}
			nicName, exists := nic.Name()
			if !exists {
				o.Logger.Warning("name for oVirt nic with id ", nicId, " is empty.")
			}
			// var nicType *objects.InterfaceType
			nicSpeedBips, exists := nic.Speed()
			if !exists {
				o.Logger.Warning("speed for oVirt nic with id ", nicId, " is empty.")
			}
			nicSpeedKbps := nicSpeedBips / 1000

			nicMtu, exists := nic.Mtu()
			if !exists {
				o.Logger.Warning("mtu for oVirt nic with id ", nicId, " is empty.")
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

			// bridged, exists := nic.Bridged()
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
				parent2child[nicBaseInterface] = append(parent2child[nicBaseInterface], nicId)
			}

			nicBonding, exists := nic.Bonding()
			if exists {
				// Bond interface, we give it a type of LAG
				nicType = &objects.LAGInterfaceType
				slaves, exists := nicBonding.Slaves()
				if exists {
					for _, slave := range slaves.Slice() {
						master2slave[nicId] = append(master2slave[nicId], slave.MustId())
					}
				}
			}

			if nicType == nil {
				// This is a physical interface.
				// TODO: depending on speed assign different nic type
				nicType = &objects.OtherInterfaceType
			}

			var nicVlan *objects.Vlan
			var err error
			vlan, exists := nic.Vlan()
			if exists {
				vlanId, exists := vlan.Id()
				if exists {
					var vlanStatus *objects.VlanStatus
					if nicEnabled {
						vlanStatus = &objects.VlanStatusActive
					} else {
						vlanStatus = &objects.VlanStatusReserved
					}
					nicVlan, err = nbi.AddVlan(&objects.Vlan{
						NetboxObject: objects.NetboxObject{
							Tags: o.SourceTags,
						},
						Name:   fmt.Sprintf("VLAN-%d", vlanId),
						Vid:    int(vlanId),
						Status: vlanStatus,
						Tenant: nbHost.Tenant,
					})
					if err != nil {
						return fmt.Errorf("failed to add oVirt vlan %s with error: %v", nicName, err)
					}
				}
			}

			var nicTaggedVlans []*objects.Vlan
			if nicVlan != nil {
				nicTaggedVlans = []*objects.Vlan{nicVlan}
			}

			newInterface := &objects.Interface{
				NetboxObject: objects.NetboxObject{
					Tags:        o.SourceTags,
					Description: nicComment,
				},
				Device: nbHost,
				Name:   nicName,
				Speed:  objects.InterfaceSpeed(nicSpeedKbps),
				Status: nicEnabled,
				MTU:    nicMtu,
				Type:   nicType,
				CustomFields: map[string]string{
					"source_id": nicId,
				},
				TaggedVlans: nicTaggedVlans,
			}

			processedNicsIds[nicId] = true
			hostInterfaces[nicId] = newInterface
		}

		// Second loop to add relations between interfaces (e.g. [eno1, eno2] -> bond1)
		for masterId, slavesIds := range master2slave {
			var err error
			masterInterface := hostInterfaces[masterId]
			if _, ok := processedNicsIds[masterId]; ok {
				masterInterface, err = nbi.AddInterface(masterInterface)
				if err != nil {
					return fmt.Errorf("failed to add oVirt master interface %s with error: %v", masterInterface.Name, err)
				}
				delete(processedNicsIds, masterId)
				hostInterfaces[masterId] = masterInterface
			}
			for _, slaveId := range slavesIds {
				slaveInterface := hostInterfaces[slaveId]
				slaveInterface.LAG = masterInterface
				slaveInterface, err := nbi.AddInterface(slaveInterface)
				if err != nil {
					return fmt.Errorf("failed to add oVirt slave interface %s with error: %v", slaveInterface.Name, err)
				}
				delete(processedNicsIds, slaveId)
				hostInterfaces[slaveId] = slaveInterface
			}
		}

		// Third loop we connect children with parents (e.g. [bond1.605, bond1.604, bond1.603] -> bond1)
		for parent, children := range parent2child {
			parentInterface := hostInterfaces[parent]
			if _, ok := processedNicsIds[parent]; ok {
				parentInterface, err := nbi.AddInterface(parentInterface)
				if err != nil {
					return fmt.Errorf("failed to add oVirt parent interface %s with error: %v", parentInterface.Name, err)
				}
				delete(processedNicsIds, parent)
			}
			for _, child := range children {
				childInterface := hostInterfaces[child]
				childInterface.ParentInterface = parentInterface
				childInterface, err := nbi.AddInterface(childInterface)
				if err != nil {
					return fmt.Errorf("failed to add oVirt child interface %s with error: %v", childInterface.Name, err)
				}
				hostInterfaces[child] = childInterface
				delete(processedNicsIds, child)
			}
		}
		// Now we check if there are any nics that were not processed
		for nicId := range processedNicsIds {
			_, err := nbi.AddInterface(hostInterfaces[nicId])
			if err != nil {
				return fmt.Errorf("failed to add oVirt interface %s with error: %v", hostInterfaces[nicId].Name, err)
			}
		}
	}
	return nil
}

func (o *OVirtSource) syncVms(nbi *inventory.NetBoxInventory) error {
	for vmId, vm := range o.Vms {
		// VM name, which is used as unique identifier for VMs in Netbox
		vmName, exists := vm.Name()
		if !exists {
			o.Logger.Warning("name for oVirt vm with id ", vmId, " is empty. VM has to have unique name to be synced to netbox. Skipping...")
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
		host, exists := vm.Host()
		if exists {
			if _, ok := o.Hosts[host.MustId()]; ok {
				vmHostDevice = nbi.DevicesIndexByUuid[o.Hosts[host.MustId()].MustHardwareInformation().MustUuid()]
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
		vmOsType := "Generic OS"
		vmOsVersion := "Generic Version"
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
			return fmt.Errorf("failed adding oVirt vm's Platform %v with error: %s", vmPlatform, err)
		}

		newVM, err := nbi.AddVM(&objects.VM{
			NetboxObject: objects.NetboxObject{
				Tags: o.SourceTags,
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
			Memory:      int(vmMemorySizeBytes / 1024 / 1024),      // MBs
			Disk:        int(vmDiskSizeBytes / 1024 / 1024 / 1024), // GBs
		})
		if err != nil {
			return fmt.Errorf("failed to sync oVirt vm: %v", err)
		}

		err = o.syncVmInterfaces(nbi, vm, newVM)
		if err != nil {
			return fmt.Errorf("failed to sync oVirt vm's interfaces: %v", err)
		}
	}

	return nil
}

// Syncs VM's interfaces to Netbox
func (o *OVirtSource) syncVmInterfaces(nbi *inventory.NetBoxInventory, ovirtVm *ovirtsdk4.Vm, netboxVm *objects.VM) error {
	var vmPrimaryIpv4 *objects.IPAddress
	var vmPrimaryIpv6 *objects.IPAddress
	if reportedDevices, exist := ovirtVm.ReportedDevices(); exist {
		for _, reportedDevice := range reportedDevices.Slice() {
			if reportedDeviceType, exist := reportedDevice.Type(); exist {
				if reportedDeviceType == "network" {
					// We add interface to the list
					var vmInterface *objects.VMInterface
					var err error
					if reportedDeviceName, exists := reportedDevice.Name(); exists {
						vmInterface, err = nbi.AddVMInterface(&objects.VMInterface{
							NetboxObject: objects.NetboxObject{
								Tags:        o.SourceTags,
								Description: reportedDevice.MustDescription(),
							},
							VM:         netboxVm,
							Name:       reportedDeviceName,
							MACAddress: strings.ToUpper(reportedDevice.MustMac().MustAddress()),
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

									ipAddress, err := nbi.AddIPAddress(&objects.IPAddress{
										NetboxObject: objects.NetboxObject{
											Tags: o.SourceTags,
										},
										Address:            ipAddress + ipMask,
										Tenant:             netboxVm.Tenant,
										Status:             &objects.IPAddressStatusActive,
										DNSName:            hostname,
										AssignedObjectType: objects.AssignedObjectTypeVMInterface,
										AssignedObjectId:   vmInterface.Id,
										AssignedObject:     vmInterface,
									})

									if err != nil {
										// TODO: return should be here, commented just for now
										// return fmt.Errorf("failed to sync oVirt vm's interface %s ip %s: %v", vmInterface, ip.MustAddress(), err)
										o.Logger.Error(fmt.Sprintf("failed to sync oVirt vm's interface %s ip %s: %v", vmInterface, ip.MustAddress(), err))

									}

									// TODO: criteria to determine if reported device is primary IP
									switch ipVersion {
									case "v4":
										if vmPrimaryIpv4 == nil {
											vmPrimaryIpv4 = ipAddress
										}
									case "v6":
										if vmPrimaryIpv6 == nil {
											vmPrimaryIpv6 = ipAddress
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
	// Update netboxVM with primary IPs
	// TODO: determine which ip is primary ipv4 and which is primary ipv6
	// TODO: then assign them to netboxVM
	// if vmPrimaryIpv4 != nil && (netboxVm.PrimaryIPv4 == nil || vmPrimaryIpv4.Address != netboxVm.PrimaryIPv4.Address) {
	// 	netboxVm.PrimaryIPv4 = vmPrimaryIpv4
	// 	if _, err := nbi.AddVM(netboxVm); err != nil {
	// 		return fmt.Errorf("failed to sync oVirt vm's primary ipv4: %v", err)
	// 	}
	// }
	// if vmPrimaryIpv6 != nil && (netboxVm.PrimaryIPv6 == nil || vmPrimaryIpv6.Address != netboxVm.PrimaryIPv6.Address) {
	// 	netboxVm.PrimaryIPv6 = vmPrimaryIpv6
	// 	if _, err := nbi.AddVM(netboxVm); err != nil {
	// 		return fmt.Errorf("failed to sync oVirt vm's primary ipv6: %v", err)
	// 	}
	// }

	return nil
}
