package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
)

func (oss *Source) initServers(ctx context.Context) error {
	allPages, err := servers.List(oss.ComputeClient, servers.ListOpts{}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error listing servers: %s", err)
	}

	err = servers.ExtractServersInto(allPages, &oss.Servers)
	if err != nil {
		return fmt.Errorf("error extracting servers: %s", err)
	}

	return nil
}

func (oss *Source) initFlavors(ctx context.Context) error {
	allPages, err := flavors.ListDetail(oss.ComputeClient, flavors.ListOpts{}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error listing flavors: %s", err)
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		return fmt.Errorf("error extracting flavors: %s", err)
	}

	oss.Flavors = allFlavors
	return nil
}

func (oss *Source) initNetworks(ctx context.Context) error {
	allPages, err := networks.List(oss.NetworkClient, networks.ListOpts{}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error listing networks: %s", err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return fmt.Errorf("error extracting networks: %s", err)
	}

	oss.Networks = allNetworks
	return nil
}

func (oss *Source) initVolumes(ctx context.Context) error {
	allPages, err := volumes.List(oss.BlockStorageClient, volumes.ListOpts{}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error listing volumes: %s", err)
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return fmt.Errorf("error extracting volumes: %s", err)
	}

	oss.Volumes = allVolumes
	return nil
}

func (oss *Source) initImages(ctx context.Context) error {
	allPages, err := images.List(oss.ImageClient, images.ListOpts{}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error listing images: %s", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return fmt.Errorf("error extracting images: %s", err)
	}

	oss.Images = allImages
	return nil
}
