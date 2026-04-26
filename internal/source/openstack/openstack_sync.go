package openstack

import (
	"fmt"
	"regexp"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
)

const (
	// GBToMBConversion is the multiplier to convert GB to MB.
	GBToMBConversion = 1024
	// RegexMatchGroups is the expected number of groups in platform name regex.
	RegexMatchGroups = 3
	// IPv6Version is the version number for IPv6.
	IPv6Version = 6
)

func (oss *Source) syncServers(nbi *inventory.NetboxInventory) error {
	clusterTypeName := oss.SourceConfig.ClusterType
	if clusterTypeName == "" {
		clusterTypeName = "OpenStack"
	}
	clusterType, err := nbi.AddClusterType(oss.Ctx, &objects.ClusterType{
		NetboxObject: objects.NetboxObject{
			Tags: oss.GetSourceTags(),
		},
		Name: clusterTypeName,
		Slug: utils.Slugify(clusterTypeName),
	})
	if err != nil {
		return fmt.Errorf("error adding cluster type: %s", err)
	}

	cgName := fmt.Sprintf("OpenStack - %s", oss.SourceConfig.Name)
	clusterGroup, err := nbi.AddClusterGroup(oss.Ctx, &objects.ClusterGroup{
		NetboxObject: objects.NetboxObject{
			Tags: oss.GetSourceTags(),
		},
		Name: cgName,
		Slug: utils.Slugify(cgName),
	})
	if err != nil {
		return fmt.Errorf("error adding cluster group: %s", err)
	}

	clusterName := oss.SourceConfig.ClusterName
	if clusterName == "" {
		clusterName = "OpenStack Cloud"
	}
	cluster, err := nbi.AddCluster(oss.Ctx, &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			Tags:        oss.GetSourceTags(),
			Description: fmt.Sprintf("OpenStack Cluster for %s", oss.SourceConfig.Name),
		},
		Name:   clusterName,
		Type:   clusterType,
		Group:  clusterGroup,
		Status: objects.ClusterStatusActive,
	})
	if err != nil {
		return fmt.Errorf("error adding cluster: %s", err)
	}

	// Determine VM Role
	vmRole, err := nbi.AddVMDeviceRole(oss.Ctx)
	if err != nil {
		return fmt.Errorf("error adding vm device role: %s", err)
	}

	// 2. Iterate through servers and sync them as VirtualMachines
	for _, server := range oss.Servers {
		// Find flavor for resources
		var vcpus float32
		var memory int
		var disk int
		var flavorName string
		for _, flavor := range oss.Flavors {
			if server.Flavor != nil {
				if flavor.ID == server.Flavor["id"] {
					vcpus = float32(flavor.VCPUs)
					memory = flavor.RAM
					disk = flavor.Disk // GB
					flavorName = flavor.Name
					break
				}
			}
		}

		// Determine Platform
		platformName := oss.getPlatformName(&server)
		platform, err := nbi.AddPlatform(oss.Ctx, &objects.Platform{
			Name: platformName,
			Slug: utils.Slugify(platformName),
		})
		if err != nil {
			return fmt.Errorf("error adding platform: %s", err)
		}

		// Determine VM Status
		vmStatus := &objects.VMStatusActive
		if server.Status != "ACTIVE" && server.VmState != "active" {
			vmStatus = &objects.VMStatusOffline
		}

		vm := &objects.VM{
			NetboxObject: objects.NetboxObject{
				Tags:        oss.GetSourceTags(),
				Description: flavorName,
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceName:   oss.SourceConfig.Name,
					constants.CustomFieldSourceIDName: server.ID,
				},
			},
			Name:        server.Name,
			Cluster:     cluster,
			Status:      vmStatus,
			VCPUs:       vcpus,
			Memory:      memory,
			Disk:        disk,
			Role:        vmRole,
			Platform:    platform,
		}

		nbVM, err := nbi.AddVM(oss.Ctx, vm)
		if err != nil {
			return fmt.Errorf("error adding vm %s: %s", server.Name, err)
		}

		// Sync Volume/Disks
		err = oss.syncVMVolumes(nbi, nbVM, &server)
		if err != nil {
			oss.Logger.Errorf(oss.Ctx, "Error syncing volumes for VM %s: %v", nbVM.Name, err)
		}

		// Sync Interfaces and IPs
		err = oss.syncVMInterfaces(nbi, nbVM, &server)
		if err != nil {
			oss.Logger.Errorf(oss.Ctx, "Error syncing interfaces for VM %s: %v", nbVM.Name, err)
		}
	}

	return nil
}

func (oss *Source) findImageNameByID(imageID string) string {
	for _, img := range oss.Images {
		if img.ID == imageID {
			return img.Name
		}
	}
	return ""
}

func (oss *Source) getPlatformName(server *servers.Server) string {
	// Define a list of functions to try for getting platform name
	platformGetters := []func(*servers.Server) string{
		oss.getPlatformFromImageMap,
		oss.getPlatformFromImageMetadataNested,
		oss.getPlatformFromServerMetadata,
		oss.getPlatformFromVolumeMetadata,
		oss.getPlatformFromOSDistro,
	}

	for _, getter := range platformGetters {
		if name := getter(server); name != "" {
			return name
		}
	}

	return "Unknown"
}

func (oss *Source) getPlatformFromImageMap(server *servers.Server) string {
	if imgMap := server.Image; imgMap != nil {
		if imageID, ok := imgMap["id"].(string); ok && imageID != "" {
			return oss.findImageNameByID(imageID)
		}
	}
	return ""
}

func (oss *Source) getPlatformFromImageMetadataNested(server *servers.Server) string {
	if imgMap := server.Image; imgMap != nil {
		if imgMetadata, ok := imgMap["metadata"].(map[string]interface{}); ok {
			if val, ok := imgMetadata["base_image_ref"].(string); ok && val != "" {
				return oss.findImageNameByID(val)
			}
		}
	}
	return ""
}

func (oss *Source) getPlatformFromServerMetadata(server *servers.Server) string {
	if sMeta := server.Metadata; sMeta != nil {
		if val, ok := sMeta["image_name"]; ok && val != "" {
			return val
		}
		if val, ok := sMeta["image_id"]; ok && val != "" {
			return oss.findImageNameByID(val)
		}
		if val, ok := sMeta["base_image_ref"]; ok && val != "" {
			return oss.findImageNameByID(val)
		}
		if val, ok := sMeta["image_metadata.base_image_ref"]; ok && val != "" {
			return oss.findImageNameByID(val)
		}
	}
	return ""
}

func (oss *Source) getPlatformFromVolumeMetadata(server *servers.Server) string {
	for _, attachment := range server.AttachedVolumes {
		for _, vol := range oss.Volumes {
			if vol.ID == attachment.ID {
				if val, ok := vol.Metadata["image_name"]; ok && val != "" {
					return val
				}
				if val, ok := vol.Metadata["image_id"]; ok && val != "" {
					return oss.findImageNameByID(val)
				}
				if val, ok := vol.Metadata["base_image_ref"]; ok && val != "" {
					return oss.findImageNameByID(val)
				}
			}
		}
	}
	return ""
}

func (oss *Source) getPlatformFromOSDistro(server *servers.Server) string {
	if sMeta := server.Metadata; sMeta != nil {
		if distro, ok := sMeta["os_distro"]; ok && distro != "" {
			return oss.cleanPlatformName(distro)
		}
	}
	return ""
}

func (oss *Source) cleanPlatformName(name string) string {
	// e.g. almalinux9 -> Almalinux 9
	re := regexp.MustCompile(`([a-zA-Z]+)(\d+)`)
	matches := re.FindStringSubmatch(name)
	if len(matches) == RegexMatchGroups {
		return cases.Title(language.English).String(matches[1]) + " " + matches[2]
	}
	return cases.Title(language.English).String(name)
}

func (oss *Source) syncVMVolumes(
	nbi *inventory.NetboxInventory,
	nbVM *objects.VM,
	server *servers.Server,
) error {
	for _, attached := range server.AttachedVolumes {
		for _, vol := range oss.Volumes {
			if vol.ID == attached.ID {
				_, err := nbi.AddVirtualDisk(oss.Ctx, &objects.VirtualDisk{
					NetboxObject: objects.NetboxObject{
						Description: fmt.Sprintf("Volume ID: %s", vol.ID),
					},
					VM:   nbVM,
					Name: vol.Name,
					Size: vol.Size * GBToMBConversion, // GB to MB
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (oss *Source) syncVMInterfaces(
	nbi *inventory.NetboxInventory,
	nbVM *objects.VM,
	server *servers.Server,
) error {
	var primaryIPv4 *objects.IPAddress
	var primaryIPv6 *objects.IPAddress

	addrMap := server.Addresses
	if addrMap == nil {
		return nil
	}
	for netName, addrs := range addrMap {
		vmi, err := nbi.AddVMInterface(oss.Ctx, &objects.VMInterface{
			NetboxObject: objects.NetboxObject{
				Tags: oss.GetSourceTags(),
			},
			VM:      nbVM,
			Name:    netName,
			Enabled: true,
		})
		if err != nil {
			return err
		}

		// Handle list of addresses
		addrList, ok := addrs.([]interface{})
		if !ok {
			continue
		}

		for _, a := range addrList {
			addrMap, ok := a.(map[string]interface{})
			if !ok {
				continue
			}

			ipStr, _ := addrMap["addr"].(string)
			version, _ := addrMap["version"].(float64)

			if ipStr == "" {
				continue
			}

			// TODO: OpenStack API exposes subnet and prefix info available on the Network endpoint.
			// Currently, this is a known limitation where we hardcode /32 for IPv4 and /64 for IPv6.
			prefix := "32"
			if int(version) == IPv6Version {
				prefix = "64"
			}

			nbIP, err := nbi.AddIPAddress(oss.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: oss.GetSourceTags(),
				},
				Address:            fmt.Sprintf("%s/%s", ipStr, prefix),
				AssignedObjectType: constants.ContentTypeVirtualizationVMInterface,
				AssignedObjectID:   vmi.ID,
				Status:             &objects.IPAddressStatusActive,
			})
			if err != nil {
				oss.Logger.Errorf(oss.Ctx, "Error adding IP %s to interface %s: %v", ipStr, netName, err)
				continue
			}

			// Set primary if not already set
			if int(version) == 4 && primaryIPv4 == nil {
				primaryIPv4 = nbIP
			} else if int(version) == 6 && primaryIPv6 == nil {
				primaryIPv6 = nbIP
			}
		}
	}

	// Update VM with primary IPs if found
	if primaryIPv4 != nil || primaryIPv6 != nil {
		err := common.SetPrimaryIPAddressForObject(oss.Ctx, nbi, nbVM, primaryIPv4, primaryIPv6)
		if err != nil {
			return fmt.Errorf("error updating vm primary ip: %s", err)
		}
	}

	return nil
}
