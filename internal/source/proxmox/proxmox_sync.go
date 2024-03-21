package proxmox

import (
	"fmt"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/luthermonson/go-proxmox"
)

func (ps *ProxmoxSource) syncCluster(nbi *inventory.NetboxInventory) error {
	clusterSite, err := common.MatchClusterToSite(ps.Ctx, nbi, ps.Cluster.Name, ps.ClusterSiteRelations)
	if err != nil {
		return err
	}
	clusterTenant, err := common.MatchClusterToTenant(ps.Ctx, nbi, ps.Cluster.Name, ps.ClusterTenantRelations)
	if err != nil {
		return err
	}
	clusterType, err := nbi.AddClusterType(ps.Ctx, &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: ps.Config.SourceTags,
			CustomFields: map[string]string{
				constants.CustomFieldSourceName: ps.SourceConfig.Name,
			},
		},
		Name: "Proxmox",
		Slug: utils.Slugify("Proxmox"),
	})
	if err != nil {
		return fmt.Errorf("proxmox cluster type: %s", err)
	}

	// Check if proxmox is running standalon node (in that case cluster name is empty)
	if ps.Cluster.Name == "" {
		ps.Cluster.Name = "ProxmoxStandalone"
	}

	nbCluster, err := nbi.AddCluster(ps.Ctx, &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			Tags: ps.SourceTags,
			CustomFields: map[string]string{
				constants.CustomFieldSourceName:   ps.SourceConfig.Name,
				constants.CustomFieldSourceIDName: ps.Cluster.ID,
			},
		},
		Name:   ps.Cluster.Name,
		Type:   clusterType,
		Site:   clusterSite,
		Tenant: clusterTenant,
	})
	if err != nil {
		return fmt.Errorf("sync cluster: %s", err)
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
			hostSite, err = common.MatchHostToSite(ps.Ctx, nbi, node.Name, ps.HostSiteRelations)
			if err != nil {
				return fmt.Errorf("match host to site: %s", err)
			}
		}
		hostTenant, err := common.MatchHostToTenant(ps.Ctx, nbi, node.Name, ps.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("match host to tenant: %s", err)
		}
		// TODO: find a way to get device type info from proxmox
		hostManufacturer, err := nbi.AddManufacturer(ps.Ctx, &objects.Manufacturer{
			Name: constants.DefaultManufacturer,
			Slug: utils.Slugify(constants.DefaultManufacturer),
		})
		if err != nil {
			return fmt.Errorf("adding host manufacturer: %s", err)
		}
		hostDeviceType, err := nbi.AddDeviceType(ps.Ctx, &objects.DeviceType{
			Manufacturer: hostManufacturer,
			Model:        constants.DefaultModel,
			Slug:         utils.Slugify(hostManufacturer.Name + constants.DefaultModel),
		})
		if err != nil {
			return fmt.Errorf("adding host device type: %s", err)
		}

		nbHost, err := nbi.AddDevice(ps.Ctx, &objects.Device{
			NetboxObject: objects.NetboxObject{
				Tags: ps.Config.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName:       ps.SourceConfig.Name,
					constants.CustomFieldHostCPUCoresName: fmt.Sprintf("%d", node.CPUInfo.CPUs),
					constants.CustomFieldHostMemoryName:   fmt.Sprintf("%d GB", node.Memory.Total/constants.GiB),
				},
			},
			Name:       node.Name,
			DeviceRole: nbi.DeviceRolesIndexByName["Server"],
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
	for _, nodeNetwork := range ps.NodeNetworks[node.Name] {
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
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: ps.SourceConfig.Name,
				},
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
	for nodeName, vms := range ps.Vms {
		nbHost := ps.NetboxNodes[nodeName]
		for _, vm := range vms {
			// Determine VM status
			vmStatus := &objects.VMStatusActive
			if vm.Status == "stopped" {
				vmStatus = &objects.VMStatusOffline
			}
			// Determine VM tenant
			vmTenant, err := common.MatchVMToTenant(ps.Ctx, nbi, vm.Name, ps.VMTenantRelations)
			if err != nil {
				return fmt.Errorf("match vm to tenant: %s", err)
			}
			nbVM, err := nbi.AddVM(ps.Ctx, &objects.VM{
				NetboxObject: objects.NetboxObject{
					Tags: ps.SourceTags,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName:   ps.SourceConfig.Name,
						constants.CustomFieldSourceIDName: fmt.Sprintf("%d", vm.VMID),
					},
				},
				Host:    nbHost,
				Cluster: ps.NetboxCluster, // Default single proxmox cluster
				Tenant:  vmTenant,
				VCPUs:   float32(vm.CPUs),
				Memory:  int(vm.MaxMem / constants.MiB),  // Memory is in MB
				Disk:    int(vm.MaxDisk / constants.GiB), // Disk is in GB
				Site:    nbHost.Site,
				Name:    vm.Name,
				Status:  vmStatus,
			})
			if err != nil {
				return fmt.Errorf("new vm: %s", err)
			}

			err = ps.syncVMNetworks(nbi, nbVM)
			if err != nil {
				return fmt.Errorf("sync vm networks: %s", err)
			}
		}
	}
	return nil
}

func (ps *ProxmoxSource) syncVMNetworks(nbi *inventory.NetboxInventory, nbVM *objects.VM) error {
	vmIPv4Addresses := make([]*objects.IPAddress, 0)
	vmIPv6Addresses := make([]*objects.IPAddress, 0)
	for _, vmNetwork := range ps.VMNetworks[nbVM.Name] {
		if utils.FilterInterfaceName(vmNetwork.Name, ps.SourceConfig.InterfaceFilter) {
			ps.Logger.Debugf(ps.Ctx, "interface %s is filtered out with interface filter %s", vmNetwork.Name, ps.SourceConfig.InterfaceFilter)
			continue
		}
		nbVMIface, err := nbi.AddVMInterface(ps.Ctx, &objects.VMInterface{
			NetboxObject: objects.NetboxObject{
				Tags: ps.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: ps.SourceConfig.Name,
				},
			},
			Name:       vmNetwork.Name,
			MACAddress: strings.ToUpper(vmNetwork.HardwareAddress),
			VM:         nbVM,
		})
		if err != nil {
			return fmt.Errorf("add vm interface: %s", err)
		}

		for _, ipAddress := range vmNetwork.IPAddresses {
			if !utils.SubnetsContainIPAddress(ipAddress.IPAddress, ps.SourceConfig.IgnoredSubnets) {
				nbIPAddress, err := nbi.AddIPAddress(ps.Ctx, &objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: ps.SourceTags,
						CustomFields: map[string]string{
							constants.CustomFieldSourceName: ps.SourceConfig.Name,
						},
					},
					Address:            fmt.Sprintf("%s/%d", ipAddress.IPAddress, ipAddress.Prefix),
					DNSName:            utils.ReverseLookup(ipAddress.IPAddress),
					Tenant:             nbVM.Tenant,
					AssignedObjectType: objects.AssignedObjectTypeVMInterface,
					AssignedObjectID:   nbVMIface.ID,
					Status:             &objects.IPAddressStatusActive, //TODO: this is hardcoded
				})
				if err != nil {
					return fmt.Errorf("add ip address: %s", err)
				}
				switch ipAddress.IPAddressType {
				case "ipv4":
					vmIPv4Addresses = append(vmIPv4Addresses, nbIPAddress)
				case "ipv6":
					vmIPv6Addresses = append(vmIPv6Addresses, nbIPAddress)
				default:
					ps.Logger.Warningf(ps.Ctx, "wrong IP type: %s", ipAddress.IPAddressType)
				}
				prefix, err := utils.ExtractPrefixFromIPAddress(nbIPAddress.Address)
				if err != nil {
					ps.Logger.Warningf(ps.Ctx, "extract prefix from ip address: %s", err)
					continue
				}
				_, err = nbi.AddPrefix(ps.Ctx, &objects.Prefix{
					Prefix: prefix,
				})
				if err != nil {
					ps.Logger.Errorf(ps.Ctx, "adding prefix: %s", err)
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
