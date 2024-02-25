package proxmox

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/luthermonson/go-proxmox"
)

func (ps *ProxmoxSource) syncCluster(nbi *inventory.NetboxInventory) error {
	clusterSite, err := common.MatchClusterToSite(nbi, ps.Cluster.Name, ps.ClusterSiteRelations)
	if err != nil {
		return err
	}
	clusterTenant, err := common.MatchClusterToTenant(nbi, ps.Cluster.Name, ps.ClusterTenantRelations)
	if err != nil {
		return err
	}
	clusterType, err := nbi.AddClusterType(
		&objects.ClusterType{
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

	nbCluster, err := nbi.AddCluster(&objects.Cluster{
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
			hostSite, err = common.MatchHostToSite(nbi, node.Name, ps.HostSiteRelations)
			if err != nil {
				return fmt.Errorf("match host to site: %s", err)
			}
		}
		hostTenant, err := common.MatchHostToTenant(nbi, node.Name, ps.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("match host to tenant: %s", err)
		}
		// TODO: find a way to get device type info from proxmox
		hostManufacturer, err := nbi.AddManufacturer(&objects.Manufacturer{
			Name: constants.DefaultManufacturerName,
			Slug: utils.Slugify(constants.DefaultManufacturerName),
		})
		if err != nil {
			return fmt.Errorf("adding host manufacturer: %s", err)
		}
		hostDeviceType, err := nbi.AddDeviceType(&objects.DeviceType{
			Manufacturer: hostManufacturer,
			Model:        constants.DefaultModelName,
			Slug:         utils.Slugify(hostManufacturer.Name + constants.DefaultModelName),
		})
		if err != nil {
			return fmt.Errorf("adding host device type: %s", err)
		}

		nbHost, err := nbi.AddDevice(
			&objects.Device{
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

		ps.syncNodeNetworks(nbi, node)
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
		nbInterface, err := nbi.AddInterface(&objects.Interface{
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
