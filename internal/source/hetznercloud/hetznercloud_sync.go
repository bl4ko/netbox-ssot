package hetznercloud

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/inventory"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
	"github.com/src-doo/netbox-ssot/internal/source/common"
	"github.com/src-doo/netbox-ssot/internal/utils"
)

func (hcs *HetznerCloudSource) syncLocationsAndDatacenters(nbi *inventory.NetboxInventory) error {
	hcs.NetboxSites = make(map[string]*objects.Site)
	
	for _, loc := range hcs.Locations {
		// Only create a Site if we haven't created one for this city yet
		if _, exists := hcs.NetboxSites[loc.City]; exists {
			continue
		}
		
		site := &objects.Site{
			NetboxObject: objects.NetboxObject{
				Tags: hcs.GetSourceTags(),
				Description: fmt.Sprintf("Hetzner Location in %s", loc.City),
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceName: hcs.SourceConfig.Name,
					constants.CustomFieldSourceIDName: fmt.Sprintf("loc-%d", loc.ID),
				},
			},
			Name: loc.City,
			Slug: utils.Slugify(loc.City),
			Status: &objects.SiteStatusActive,
		}
		
		netboxSite, err := nbi.AddSite(hcs.Ctx, site)
		if err != nil {
			return fmt.Errorf("syncing site %s: %s", loc.City, err)
		}
		
		hcs.NetboxSites[loc.City] = netboxSite
	}
	
	hcs.NetboxLocations = make(map[string]*objects.Location)
	
	for _, dc := range hcs.Datacenters {
		var site *objects.Site
		if dc.Location != nil {
			site = hcs.NetboxSites[dc.Location.City]
		}
		
		location := &objects.Location{
			NetboxObject: objects.NetboxObject{
				Tags: hcs.GetSourceTags(),
				Description: dc.Description,
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceName: hcs.SourceConfig.Name,
					constants.CustomFieldSourceIDName: fmt.Sprintf("dc-%d", dc.ID),
				},
			},
			Name: dc.Name,
			Slug: utils.Slugify(dc.Name),
			Status: &objects.SiteStatusActive,
			Site: site,
		}
		
		netboxLocation, err := nbi.AddLocation(hcs.Ctx, location)
		if err != nil {
			return fmt.Errorf("syncing datacenter location %s: %s", dc.Name, err)
		}
		
		hcs.NetboxLocations[dc.Name] = netboxLocation
	}

	return nil
}

func (hcs *HetznerCloudSource) syncServers(nbi *inventory.NetboxInventory) error {
	// Create a default cluster for Hetzner Cloud Servers
	clusterName := fmt.Sprintf("HetznerCloud-%s", hcs.SourceConfig.Name)
	clusterType := &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: hcs.GetSourceTags(),
		},
		Name: "Hetzner Cloud",
		Slug: utils.Slugify("Hetzner Cloud"),
	}
	
	netboxClusterType, err := nbi.AddClusterType(hcs.Ctx, clusterType)
	if err != nil {
		return fmt.Errorf("error creating cluster type: %s", err)
	}

	cluster := &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			Tags: hcs.GetSourceTags(),
			Description: fmt.Sprintf("Hetzner Cloud Cluster for %s", hcs.SourceConfig.Name),
		},
		Name: clusterName,
		Type: netboxClusterType,
		Status: objects.ClusterStatusActive,
	}

	netboxCluster, err := nbi.AddCluster(hcs.Ctx, cluster)
	if err != nil {
		return fmt.Errorf("error creating cluster: %s", err)
	}

	// Create a default role for VMs
	role := &objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			Tags:        hcs.GetSourceTags(),
			Description: "Virtual Machine",
		},
		Name:   "VM",
		Slug:   utils.Slugify("VM"),
		Color:  constants.ColorBlue,
		VMRole: true,
	}
	netboxRole, err := nbi.AddDeviceRole(hcs.Ctx, role)
	if err != nil {
		return fmt.Errorf("error creating role: %s", err)
	}

	for _, server := range hcs.Servers {
		var site *objects.Site
		if server.Datacenter != nil && server.Datacenter.Location != nil {
			site = hcs.NetboxSites[server.Datacenter.Location.City]
		}
		
		var netboxPlatform *objects.Platform
		if server.Image != nil && server.Image.Description != "" {
			platform := &objects.Platform{
				NetboxObject: objects.NetboxObject{
					Tags: hcs.GetSourceTags(),
					Description: fmt.Sprintf("Platform %s", server.Image.Name),
				},
				Name: server.Image.Description,
				Slug: utils.Slugify(server.Image.Description),
			}
			netboxPlatform, _ = nbi.AddPlatform(hcs.Ctx, platform)
		}

		status := &objects.VMStatusActive
		if server.Status == hcloud.ServerStatusOff {
			status = &objects.VMStatusOffline
		}

		vm := &objects.VM{
			NetboxObject: objects.NetboxObject{
				Tags: hcs.GetSourceTags(),
				Description: server.ServerType.Name,
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceName: hcs.SourceConfig.Name,
					constants.CustomFieldSourceIDName: fmt.Sprintf("%d", server.ID),
					constants.CustomFieldServerCPUTypeName: string(server.ServerType.CPUType),
					constants.CustomFieldServerCategoryName: server.ServerType.Category,
					constants.CustomFieldServerDeprecatedName: server.ServerType.IsDeprecated(),
				},
			},
			Name: server.Name,
			Cluster: netboxCluster,
			Site: site,
			Status: status,
			Role: netboxRole,
			Platform: netboxPlatform,
			VCPUs: float32(server.ServerType.Cores),
			Memory: int(server.ServerType.Memory * 1024), // GB to MB
			Disk: server.ServerType.Disk * 1000, // GB to MB
		}
		
		netboxVM, err := nbi.AddVM(hcs.Ctx, vm)
		if err != nil {
			return fmt.Errorf("syncing server %s: %s", server.Name, err)
		}

		// Sync Public Network ("eth0")
		var nbPrimaryIPv4, nbPrimaryIPv6 *objects.IPAddress
		eth0Name := "eth0" // standard naming for Hetzner Cloud public interface
		eth0Interface := &objects.VMInterface{
			NetboxObject: objects.NetboxObject{
				Tags: hcs.GetSourceTags(),
				Description: "Public Interface",
			},
			VM: netboxVM,
			Name: eth0Name,
			Enabled: true,
		}
		
		netboxEth0, err := nbi.AddVMInterface(hcs.Ctx, eth0Interface)
		if err != nil {
			return fmt.Errorf("syncing eth0 interface for server %s: %s", server.Name, err)
		}

		// Handle Public IPv4
		if server.PublicNet.IPv4.IP != nil {
			ipv4Addr := fmt.Sprintf("%s/32", server.PublicNet.IPv4.IP.String())
			ipAddr := &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: hcs.GetSourceTags(),
				},
				Address: ipv4Addr,
				Status: &objects.IPAddressStatusActive,
				DNSName: server.PublicNet.IPv4.DNSPtr,
				AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
				AssignedObjectID: netboxEth0.ID,
			}
			
			netboxIP, err := nbi.AddIPAddress(hcs.Ctx, ipAddr)
			if err != nil {
				return fmt.Errorf("syncing ipv4 for server %s: %s", server.Name, err)
			}
			
			nbPrimaryIPv4 = netboxIP
		}

		// Handle Public IPv6
		if server.PublicNet.IPv6.IP != nil && server.PublicNet.IPv6.Network != nil {
			// Extract prefix length from the network
			ones, _ := server.PublicNet.IPv6.Network.Mask.Size()
			
			ipStr := server.PublicNet.IPv6.IP.String()
			// Hetzner returns the IPv6 network (e.g., 2a01:4f8:x:y::) for the server. 
			// We append the "1" to represent the typical host IP inside this /64 interface.
			if strings.HasSuffix(ipStr, "::") {
				ipStr = ipStr + "1"
			}
			
			ipv6Addr := fmt.Sprintf("%s/%d", ipStr, ones)
			
			ipAddr := &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: hcs.GetSourceTags(),
				},
				Address: ipv6Addr,
				Status: &objects.IPAddressStatusActive,
				AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
				AssignedObjectID: netboxEth0.ID,
			}
			
			netboxIP, err := nbi.AddIPAddress(hcs.Ctx, ipAddr)
			if err != nil {
				return fmt.Errorf("syncing ipv6 for server %s: %s", server.Name, err)
			}

			nbPrimaryIPv6 = netboxIP
		}

		// Set primary IPs on VM
		err = common.SetPrimaryIPAddressForObject(hcs.Ctx, nbi, netboxVM, nbPrimaryIPv4, nbPrimaryIPv6)
		if err != nil {
			return fmt.Errorf("setting primary ips for server %s: %s", server.Name, err)
		}

		// Handle Private Networks (eth1, eth2, etc.)
		for i, privateNet := range server.PrivateNet {
			// Hetzner attaches private networks starting at eth1 and going sequentially
			ifaceName := fmt.Sprintf("eth%d", i+1)
			
			macAddr := &objects.MACAddress{
				MAC: privateNet.MACAddress,
			}
			
			privInterface := &objects.VMInterface{
				NetboxObject: objects.NetboxObject{
					Tags: hcs.GetSourceTags(),
					Description: fmt.Sprintf("Private Network %d", privateNet.Network.ID),
				},
				VM: netboxVM,
				Name: ifaceName,
				Enabled: true,
				PrimaryMACAddress: macAddr,
			}
			
			netboxPrivIface, err := nbi.AddVMInterface(hcs.Ctx, privInterface)
			if err != nil {
				return fmt.Errorf("syncing private interface %s for server %s: %s", ifaceName, server.Name, err)
			}

			if privateNet.IP != nil {
				// Typically Hetzner private networks use /32 on the server side in their routing structure,
				// or /24 depending on perspective. Assuming /32 as the specific host IP binding in Netbox
				privIPAddr := fmt.Sprintf("%s/32", privateNet.IP.String())
				
				ipAddr := &objects.IPAddress{
					NetboxObject: objects.NetboxObject{
						Tags: hcs.GetSourceTags(),
					},
					Address: privIPAddr,
					Status: &objects.IPAddressStatusActive,
					AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
					AssignedObjectID: netboxPrivIface.ID,
				}
				
				_, err = nbi.AddIPAddress(hcs.Ctx, ipAddr)
				if err != nil {
					return fmt.Errorf("syncing private ip %s for server %s: %s", privIPAddr, server.Name, err)
				}
			}
		}
	}
	
	return nil
}

func (hcs *HetznerCloudSource) syncNetworks(nbi *inventory.NetboxInventory) error {
	for _, network := range hcs.Networks {
		prefix := &objects.Prefix{
			NetboxObject: objects.NetboxObject{
				Tags: hcs.GetSourceTags(),
				Description: network.Name,
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceName: hcs.SourceConfig.Name,
					constants.CustomFieldSourceIDName: fmt.Sprintf("%d", network.ID),
				},
			},
			Prefix: network.IPRange.String(),
			Status: &objects.PrefixStatusActive,
		}
		
		_, err := nbi.AddPrefix(hcs.Ctx, prefix)
		if err != nil {
			hcs.Logger.Errorf(hcs.Ctx, "Error syncing network %s: %s", network.Name, err)
		}
	}
	return nil
}

func (hcs *HetznerCloudSource) syncFloatingIPs(nbi *inventory.NetboxInventory) error {
	for _, fip := range hcs.FloatingIPs {
		ipAddr := &objects.IPAddress{
			NetboxObject: objects.NetboxObject{
				Tags: hcs.GetSourceTags(),
				Description: fip.Description,
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceName: hcs.SourceConfig.Name,
					constants.CustomFieldSourceIDName: fmt.Sprintf("%d", fip.ID),
				},
			},
			Address: fmt.Sprintf("%s/%d", fip.IP.String(), 32),
			Status: &objects.IPAddressStatusActive,
		}
		
		_, err := nbi.AddIPAddress(hcs.Ctx, ipAddr)
		if err != nil {
			hcs.Logger.Errorf(hcs.Ctx, "Error syncing floating IP %s: %s", fip.IP.String(), err)
		}
	}
	return nil
}
