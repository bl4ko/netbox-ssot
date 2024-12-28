package proxmox

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/luthermonson/go-proxmox"
)

func (ps *ProxmoxSource) syncCluster(nbi *inventory.NetboxInventory) error {
	clusterSite, err := common.MatchClusterToSite(ps.Ctx, nbi, ps.Cluster.Name, ps.SourceConfig.ClusterSiteRelations)
	if err != nil {
		return err
	}
	clusterTenant, err := common.MatchClusterToTenant(ps.Ctx, nbi, ps.Cluster.Name, ps.SourceConfig.ClusterTenantRelations)
	if err != nil {
		return err
	}
	clusterTypeStruct := &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: ps.Config.SourceTags,
		},
		Name: "Proxmox",
		Slug: utils.Slugify("Proxmox"),
	}
	clusterType, err := nbi.AddClusterType(ps.Ctx, clusterTypeStruct)
	if err != nil {
		return fmt.Errorf("add cluster type %+v: %s", clusterTypeStruct, err)
	}

	// Check if proxmox is running standalon node.
	// in that case cluster name is empty and should set SourceConfig.Name for Cluster.Name
	if ps.Cluster.Name == "" {
		ps.Cluster.Name = ps.SourceConfig.Name
	}

	clusterStruct := &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			Tags: ps.SourceTags,
		},
		Name:   ps.Cluster.Name,
		Type:   clusterType,
		Site:   clusterSite,
		Tenant: clusterTenant,
	}
	nbCluster, err := nbi.AddCluster(ps.Ctx, clusterStruct)
	if err != nil {
		return fmt.Errorf("add cluster %+v: %s", clusterStruct, err)
	}
	ps.NetboxCluster = nbCluster
	return nil
}

func (ps *ProxmoxSource) syncNodes(nbi *inventory.NetboxInventory) error {
	ps.NetboxNodes = make(map[string]*objects.Device, len(ps.Nodes))
	for _, node := range ps.Nodes {
		hostSite := ps.NetboxCluster.Site
		var err error
		if hostSite == nil {
			hostSite, err = common.MatchHostToSite(ps.Ctx, nbi, node.Name, ps.SourceConfig.HostSiteRelations)
			if err != nil {
				return fmt.Errorf("match host to site: %s", err)
			}
		}
		hostTenant, err := common.MatchHostToTenant(ps.Ctx, nbi, node.Name, ps.SourceConfig.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("match host to tenant: %s", err)
		}
		// TODO: find a way to get device type info from proxmox
		manufacturerStruct := &objects.Manufacturer{
			NetboxObject: objects.NetboxObject{
				Description: constants.DefaultManufacturerDescription,
			},
			Name: constants.DefaultManufacturer,
			Slug: utils.Slugify(constants.DefaultManufacturer),
		}
		hostManufacturer, err := nbi.AddManufacturer(ps.Ctx, manufacturerStruct)
		if err != nil {
			return fmt.Errorf("adding host manufacturer %+v: %s", manufacturerStruct, err)
		}
		deviceTypeStruct := &objects.DeviceType{
			NetboxObject: objects.NetboxObject{
				Description: constants.DefaultDeviceTypeDescription,
			},
			Manufacturer: hostManufacturer,
			Model:        constants.DefaultModel,
			Slug:         utils.Slugify(hostManufacturer.Name + constants.DefaultModel),
		}
		hostDeviceType, err := nbi.AddDeviceType(ps.Ctx, deviceTypeStruct)
		if err != nil {
			return fmt.Errorf("adding host device type %+v: %s", deviceTypeStruct, err)
		}

		// Match host to a role. First test if user provided relations, if not
		// use default server role.
		var hostRole *objects.DeviceRole
		if len(ps.SourceConfig.HostRoleRelations) > 0 {
			hostRole, err = common.MatchHostToRole(ps.Ctx, nbi, node.Name, ps.SourceConfig.HostRoleRelations)
			if err != nil {
				return fmt.Errorf("match host to role: %s", err)
			}
		}
		if hostRole == nil {
			hostRole, err = nbi.AddServerDeviceRole(ps.Ctx)
			if err != nil {
				return fmt.Errorf("add server device role %s", err)
			}
		}

		nbHost, err := nbi.AddDevice(ps.Ctx, &objects.Device{
			NetboxObject: objects.NetboxObject{
				Tags: ps.Config.SourceTags,
				CustomFields: map[string]interface{}{
					constants.CustomFieldHostCPUCoresName: fmt.Sprintf("%d", node.CPUInfo.CPUs),
					constants.CustomFieldHostMemoryName:   fmt.Sprintf("%d GB", node.Memory.Total/constants.GiB),
				},
			},
			Name:       node.Name,
			DeviceRole: hostRole,
			Site:       hostSite,
			Tenant:     hostTenant,
			Cluster:    ps.NetboxCluster,
			DeviceType: hostDeviceType,
		})
		if err != nil {
			return fmt.Errorf("add device: %s", err)
		}
		ps.NetboxNodes[node.Name] = nbHost

		err = ps.syncNodeNetworks(nbi, node)
		if err != nil {
			return fmt.Errorf("sync node networks: %s", err)
		}
	}
	return nil
}

func (ps *ProxmoxSource) syncNodeNetworks(nbi *inventory.NetboxInventory, node *proxmox.Node) error {
	// hostIPv4Addresses := []*objects.IPAddress TODO
	// hostIPv6Addresses := []*objects.IPAddress TODO
	for _, nodeNetwork := range ps.NodeIfaces[node.Name] {
		active := false
		if nodeNetwork.Active == 1 {
			active = true
		}
		nbHost := ps.NetboxNodes[node.Name]
		if utils.FilterInterfaceName(nodeNetwork.Iface, ps.SourceConfig.InterfaceFilter) {
			ps.Logger.Debugf(ps.Ctx, "interface %s is filtered out with interfaceFilter %s", nodeNetwork.Iface, ps.SourceConfig.InterfaceFilter)
			continue
		}
		_, err := nbi.AddInterface(ps.Ctx, &objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags: ps.Config.SourceTags,
			},
			Device: nbHost,
			Name:   nodeNetwork.Iface,
			Status: active,
			Type:   &objects.OtherInterfaceType, // TODO
			// Speed: TODO
			// Mode: TODO
			// TaggedVlans: TODO
		})
		if err != nil {
			return fmt.Errorf("add host interface: %s", err)
		}
	}
	return nil
}

// Function that synces proxmox vms to the netbox inventory.
func (ps *ProxmoxSource) syncVMs(nbi *inventory.NetboxInventory) error {
	const maxGoroutines = 50
	guard := make(chan struct{}, maxGoroutines)
	errChan := make(chan error, len(ps.Vms))
	var wg sync.WaitGroup

	for nodeName, vms := range ps.Vms {
		nbHost := ps.NetboxNodes[nodeName]
		for _, vm := range vms {
			guard <- struct{}{} // Block if maxGoroutines are running
			wg.Add(1)

			go func(vm *proxmox.VirtualMachine, nbHost *objects.Device) {
				defer wg.Done()
				defer func() { <-guard }() // Release one spot in the semaphore

				err := ps.syncVM(nbi, vm, nbHost)
				if err != nil {
					errChan <- err
				}
			}(vm, nbHost)
		}
	}

	wg.Wait()
	close(errChan)
	close(guard)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *ProxmoxSource) syncVM(nbi *inventory.NetboxInventory, vm *proxmox.VirtualMachine, nbHost *objects.Device) error {
	// Determine VM status
	vmStatus := &objects.VMStatusActive
	if vm.Status == "stopped" {
		vmStatus = &objects.VMStatusOffline
	}

	// Determine VM tenant
	vmTenant, err := common.MatchVMToTenant(ps.Ctx, nbi, vm.Name, ps.SourceConfig.VMTenantRelations)
	if err != nil {
		return fmt.Errorf("match vm to tenant: %s", err)
	}

	var vmRole *objects.DeviceRole
	if len(ps.SourceConfig.VMRoleRelations) > 0 {
		vmRole, err = common.MatchVMToRole(ps.Ctx, nbi, vm.Name, ps.SourceConfig.VMRoleRelations)
		if err != nil {
			return fmt.Errorf("match vm to role: %s", err)
		}
	}
	if vmRole == nil {
		vmRole, err = nbi.AddVMDeviceRole(ps.Ctx)
		if err != nil {
			return fmt.Errorf("add vm device role: %s", err)
		}
	}

	// Add VM to Netbox
	vmStruct := &objects.VM{
		NetboxObject: objects.NetboxObject{
			Tags: ps.SourceTags,
			CustomFields: map[string]interface{}{
				constants.CustomFieldSourceName:   ps.SourceConfig.Name,
				constants.CustomFieldSourceIDName: fmt.Sprintf("%d", vm.VMID),
			},
		},
		Host:    nbHost,
		Cluster: ps.NetboxCluster, // Default single proxmox cluster
		Tenant:  vmTenant,
		Role:    vmRole,
		VCPUs:   float32(vm.CPUs),
		Memory:  int(vm.MaxMem / constants.MiB),  //nolint:gosec
		Disk:    int(vm.MaxDisk / constants.MiB), //nolint:gosec
		Site:    nbHost.Site,
		Name:    vm.Name,
		Status:  vmStatus,
	}
	nbVM, err := nbi.AddVM(ps.Ctx, vmStruct)
	if err != nil {
		return fmt.Errorf("add vm: %s", err)
	}

	// Sync VM networks
	err = ps.syncVMNetworks(nbi, nbVM)
	if err != nil {
		return fmt.Errorf("sync vm networks: %s", err)
	}

	return nil
}

func (ps *ProxmoxSource) syncVMNetworks(nbi *inventory.NetboxInventory, nbVM *objects.VM) error {
	vmIPv4Addresses := make([]*objects.IPAddress, 0)
	vmIPv6Addresses := make([]*objects.IPAddress, 0)
	for _, vmNetwork := range ps.VMIfaces[nbVM.Name] {
		if utils.FilterInterfaceName(vmNetwork.Name, ps.SourceConfig.InterfaceFilter) {
			ps.Logger.Debugf(ps.Ctx, "interface %s is filtered out with interface filter %s", vmNetwork.Name, ps.SourceConfig.InterfaceFilter)
			continue
		}
		vmInterfaceStruct := &objects.VMInterface{
			NetboxObject: objects.NetboxObject{
				Tags: ps.SourceTags,
			},
			Name:       vmNetwork.Name,
			MACAddress: strings.ToUpper(vmNetwork.HardwareAddress),
			VM:         nbVM,
		}
		nbVMIface, err := nbi.AddVMInterface(ps.Ctx, vmInterfaceStruct)
		if err != nil {
			return fmt.Errorf("add vm interface %+v: %s", vmInterfaceStruct, err)
		}

		for _, ipAddress := range vmNetwork.IPAddresses {
			if utils.IsPermittedIPAddress(ipAddress.IPAddress, ps.SourceConfig.PermittedSubnets, ps.SourceConfig.IgnoredSubnets) {
				ipAddress.IPAddress = utils.RemoveZoneIndexFromIPAddress(ipAddress.IPAddress)
				ipAddressStruct := &objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: ps.SourceTags,
						CustomFields: map[string]interface{}{
							constants.CustomFieldArpEntryName: false,
						},
					},
					Address:            fmt.Sprintf("%s/%d", ipAddress.IPAddress, ipAddress.Prefix),
					DNSName:            utils.ReverseLookup(ipAddress.IPAddress),
					Tenant:             nbVM.Tenant,
					AssignedObjectType: objects.AssignedObjectTypeVMInterface,
					AssignedObjectID:   nbVMIface.ID,
					Status:             &objects.IPAddressStatusActive, //TODO: this is hardcoded
				}
				nbIPAddress, err := nbi.AddIPAddress(ps.Ctx, ipAddressStruct)
				if err != nil {
					ps.Logger.Warningf(ps.Ctx, "failed adding ip address %s with: %s", ipAddressStruct, err)
					continue
				}
				switch ipAddress.IPAddressType {
				case "ipv4":
					vmIPv4Addresses = append(vmIPv4Addresses, nbIPAddress)
				case "ipv6":
					vmIPv6Addresses = append(vmIPv6Addresses, nbIPAddress)
				default:
					ps.Logger.Warningf(ps.Ctx, "wrong IP type: %s for ip %s", ipAddress.IPAddressType, ipAddress.IPAddress)
				}
				prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(nbIPAddress.Address)
				if err != nil {
					ps.Logger.Warningf(ps.Ctx, "failed extracting prefix from ip address %s with: %s", nbIPAddress.Address, err)
				} else if (ipAddress.IPAddressType == "ipv4" && mask != constants.MaxIPv4MaskBits) || (ipAddress.IPAddressType == "ipv6" && mask != constants.MaxIPv6MaskBits) {
					_, err = nbi.AddPrefix(ps.Ctx, &objects.Prefix{
						Prefix: prefix,
					})
					if err != nil {
						ps.Logger.Errorf(ps.Ctx, "adding prefix: %s", err)
					}
				}
			}
		}
	}
	// From all IPv4 addresses and IPv6 addresses determine primary ips
	if len(vmIPv4Addresses) > 0 || len(vmIPv6Addresses) > 0 {
		nbVMCopy := *nbVM
		if len(vmIPv4Addresses) > 0 {
			// TODO: add criteria for primary IPv4
			nbVMCopy.PrimaryIPv4 = vmIPv4Addresses[0]
		}
		if len(vmIPv6Addresses) > 0 {
			// TODO add criteria for primary IPv6
			nbVMCopy.PrimaryIPv6 = vmIPv6Addresses[0]
		}
		_, err := nbi.AddVM(ps.Ctx, &nbVMCopy)
		if err != nil {
			return fmt.Errorf("updating vm primary ip: %s", err)
		}
	}
	return nil
}

// Function that synces proxmox containers to the netbox inventory.
func (ps *ProxmoxSource) syncContainers(nbi *inventory.NetboxInventory) error {
	if len(ps.Containers) > 0 {
		// Create container role
		containerRole, err := nbi.AddContainerDeviceRole(ps.Ctx)
		if err != nil {
			return fmt.Errorf("create container role: %s", err)
		}
		for nodeName, containers := range ps.Containers {
			nbHost := ps.NetboxNodes[nodeName]
			for _, container := range containers {
				// Determine Container status
				containerStatus := &objects.VMStatusActive
				if container.Status == "stopped" {
					containerStatus = &objects.VMStatusOffline
				}
				// Determine Container tenant
				vmTenant, err := common.MatchVMToTenant(ps.Ctx, nbi, container.Name, ps.SourceConfig.VMTenantRelations)
				if err != nil {
					return fmt.Errorf("match vm to tenant: %s", err)
				}
				nbContainer, err := nbi.AddVM(ps.Ctx, &objects.VM{
					NetboxObject: objects.NetboxObject{
						Tags: ps.SourceTags,
						CustomFields: map[string]interface{}{
							constants.CustomFieldSourceIDName: fmt.Sprintf("%d", container.VMID),
						},
					},
					Host:    nbHost,
					Role:    containerRole,
					Cluster: ps.NetboxCluster, // Default single proxmox cluster
					Tenant:  vmTenant,
					VCPUs:   float32(container.CPUs),
					Memory:  int(container.MaxMem / constants.MiB),  //nolint:gosec
					Disk:    int(container.MaxDisk / constants.GiB), //nolint:gosec
					Site:    nbHost.Site,
					Name:    container.Name,
					Status:  containerStatus,
				})
				if err != nil {
					return fmt.Errorf("new vm: %s", err)
				}

				err = ps.syncContainerNetworks(nbi, nbContainer)
				if err != nil {
					return fmt.Errorf("sync container networks: %s", err)
				}
			}
		}
	}
	return nil
}

func (ps *ProxmoxSource) syncContainerNetworks(nbi *inventory.NetboxInventory, nbContainer *objects.VM) error {
	vmIPv4Addresses := make([]*objects.IPAddress, 0)
	vmIPv6Addresses := make([]*objects.IPAddress, 0)
	for _, containerIface := range ps.ContainerIfaces[nbContainer.Name] {
		if utils.FilterInterfaceName(containerIface.Name, ps.SourceConfig.InterfaceFilter) {
			ps.Logger.Debugf(ps.Ctx, "interface %s is filtered out with interface filter %s", containerIface.Name, ps.SourceConfig.InterfaceFilter)
			continue
		}
		vmIfaceStruct := &objects.VMInterface{
			NetboxObject: objects.NetboxObject{
				Tags: ps.SourceTags,
			},
			Name:       containerIface.Name,
			MACAddress: strings.ToUpper(containerIface.HWAddr),
			VM:         nbContainer,
		}
		nbVMIface, err := nbi.AddVMInterface(ps.Ctx, vmIfaceStruct)
		if err != nil {
			return fmt.Errorf("add vm interface: %s", err)
		}

		// Check if IPv4 address is present
		if utils.IsPermittedIPAddress(containerIface.Inet, ps.SourceConfig.PermittedSubnets, ps.SourceConfig.IgnoredSubnets) {
			// Check if IPv4 address is present
			if containerIface.Inet != "" {
				nbIPAddress, err := nbi.AddIPAddress(ps.Ctx, &objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: ps.SourceTags,
						CustomFields: map[string]interface{}{
							constants.CustomFieldArpEntryName: false,
						},
					},
					Address:            containerIface.Inet,
					DNSName:            utils.ReverseLookup(containerIface.Inet),
					Tenant:             nbContainer.Tenant,
					AssignedObjectType: objects.AssignedObjectTypeVMInterface,
					AssignedObjectID:   nbVMIface.ID,
					Status:             &objects.IPAddressStatusActive, //TODO: this is hardcoded
				})
				if err != nil {
					ps.Logger.Warningf(ps.Ctx, "add ip address: %s", err)
				} else {
					vmIPv4Addresses = append(vmIPv4Addresses, nbIPAddress)
					prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(nbIPAddress.Address)
					if err != nil {
						ps.Logger.Warningf(ps.Ctx, "extract prefix from ip address: %s", err)
					} else if mask != constants.MaxIPv4MaskBits {
						_, err = nbi.AddPrefix(ps.Ctx, &objects.Prefix{
							Prefix: prefix,
						})
						if err != nil {
							ps.Logger.Errorf(ps.Ctx, "adding prefix: %s", err)
						}
					}
				}
			}
		}
		// Check if IPv6 address is present
		if utils.IsPermittedIPAddress(containerIface.Inet6, ps.SourceConfig.PermittedSubnets, ps.SourceConfig.IgnoredSubnets) {
			if containerIface.Inet6 != "" {
				containerIface.Inet6 = utils.RemoveZoneIndexFromIPAddress(containerIface.Inet6)
				nbIPAddress, err := nbi.AddIPAddress(ps.Ctx, &objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: ps.SourceTags,
						CustomFields: map[string]interface{}{
							constants.CustomFieldArpEntryName: false,
						},
					},
					Address:            containerIface.Inet6,
					DNSName:            utils.ReverseLookup(containerIface.Inet6),
					Tenant:             nbContainer.Tenant,
					AssignedObjectType: objects.AssignedObjectTypeVMInterface,
					AssignedObjectID:   nbVMIface.ID,
					Status:             &objects.IPAddressStatusActive, //TODO: this is hardcoded
				})
				if err != nil {
					ps.Logger.Warningf(ps.Ctx, "add ipv6 address: %s", err)
				} else {
					vmIPv6Addresses = append(vmIPv6Addresses, nbIPAddress)
				}
			}
		}
	}
	// From all IPv4 addresses and IPv6 addresses determine primary ips
	if len(vmIPv4Addresses) > 0 || len(vmIPv6Addresses) > 0 {
		nbContainerCopy := *nbContainer
		if len(vmIPv4Addresses) > 0 {
			// TODO: add criteria for primary IPv4
			nbContainerCopy.PrimaryIPv4 = vmIPv4Addresses[0]
		}
		if len(vmIPv6Addresses) > 0 {
			// TODO add criteria for primary IPv6
			nbContainerCopy.PrimaryIPv6 = vmIPv6Addresses[0]
		}
		_, err := nbi.AddVM(ps.Ctx, &nbContainerCopy)
		if err != nil {
			return fmt.Errorf("updating vm primary ip: %s", err)
		}
	}
	return nil
}
