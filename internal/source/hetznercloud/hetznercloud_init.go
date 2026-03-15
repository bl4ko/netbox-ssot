package hetznercloud

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func (hcs *HetznerCloudSource) initLocations(ctx context.Context, client *hcloud.Client) error {
	locations, err := client.Location.All(ctx)
	if err != nil {
		return fmt.Errorf("error fetching locations: %s", err)
	}
	hcs.Locations = locations
	hcs.Logger.Debugf(hcs.Ctx, "Fetched %d locations", len(locations))
	return nil
}

func (hcs *HetznerCloudSource) initDatacenters(ctx context.Context, client *hcloud.Client) error {
	datacenters, err := client.Datacenter.All(ctx)
	if err != nil {
		return fmt.Errorf("error fetching datacenters: %s", err)
	}
	hcs.Datacenters = datacenters
	hcs.Logger.Debugf(hcs.Ctx, "Fetched %d datacenters", len(datacenters))
	return nil
}

func (hcs *HetznerCloudSource) initServers(ctx context.Context, client *hcloud.Client) error {
	servers, err := client.Server.All(ctx)
	if err != nil {
		return fmt.Errorf("error fetching servers: %s", err)
	}
	hcs.Servers = servers
	hcs.Logger.Debugf(hcs.Ctx, "Fetched %d servers", len(servers))
	return nil
}

func (hcs *HetznerCloudSource) initNetworks(ctx context.Context, client *hcloud.Client) error {
	networks, err := client.Network.All(ctx)
	if err != nil {
		return fmt.Errorf("error fetching networks: %s", err)
	}
	hcs.Networks = networks
	hcs.Logger.Debugf(hcs.Ctx, "Fetched %d networks", len(networks))
	return nil
}

func (hcs *HetznerCloudSource) initFloatingIPs(ctx context.Context, client *hcloud.Client) error {
	floatingIPs, err := client.FloatingIP.All(ctx)
	if err != nil {
		return fmt.Errorf("error fetching floating IPs: %s", err)
	}
	hcs.FloatingIPs = floatingIPs
	hcs.Logger.Debugf(hcs.Ctx, "Fetched %d floating IPs", len(floatingIPs))
	return nil
}

func (hcs *HetznerCloudSource) initPrimaryIPs(ctx context.Context, client *hcloud.Client) error {
	primaryIPs, err := client.PrimaryIP.All(ctx)
	if err != nil {
		return fmt.Errorf("error fetching primary IPs: %s", err)
	}
	hcs.PrimaryIPs = primaryIPs
	hcs.Logger.Debugf(hcs.Ctx, "Fetched %d primary IPs", len(primaryIPs))
	return nil
}
