package dnac

import (
	"fmt"
	"strconv"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Syncs dnac sites to netbox inventory.
func (ds *DnacSource) SyncSites(nbi *inventory.NetboxInventory) error {
	for _, site := range ds.Sites {
		dnacSite := &objects.Site{
			NetboxObject: objects.NetboxObject{
				Tags: ds.SourceTags,
			},
			Name: site.Name,
			Slug: utils.Slugify(site.Name),
		}
		for _, additionalInfo := range site.AdditionalInfo {
			switch additionalInfo.Namespace {
			case "Location":
				dnacSite.PhysicalAddress = additionalInfo.Attributes.Address
				longitude, err := strconv.ParseFloat(additionalInfo.Attributes.Longitude, 64)
				if err != nil {
					dnacSite.Longitude = longitude
				}
				latitude, err := strconv.ParseFloat(additionalInfo.Attributes.Latitude, 64)
				if err != nil {
					dnacSite.Latitude = latitude
				}
			}
		}
		nbSite, err := nbi.AddSite(dnacSite)
		if err != nil {
			return fmt.Errorf("adding site: %s", err)
		}
		if ds.SiteId2nbSite == nil {
			ds.SiteId2nbSite = make(map[string]*objects.Site)
		}
		ds.SiteId2nbSite[site.ID] = nbSite
	}
	return nil
}

func (ds *DnacSource) SyncDevices(nbi *inventory.NetboxInventory) error {
	for _, device := range ds.Devices {
		var description, comments string
		if device.Description != "" {
			description = device.Description
		}
		if len(description) > 200 {
			comments = description
			description = "See comments"
		}

		ciscoManufacturer, err := nbi.AddManufacturer(&objects.Manufacturer{
			Name: "Cisco",
			Slug: utils.Slugify("Cisco"),
		})
		if err != nil {
			return fmt.Errorf("failed creating device: %s", err)
		}

		deviceRole, err := nbi.AddDeviceRole(&objects.DeviceRole{
			Name:   device.Family,
			Slug:   utils.Slugify(device.Family),
			Color:  objects.COLOR_AQUA, // TODO
			VMRole: false,
		})
		if err != nil {
			return fmt.Errorf("adding dnac device role: %s", err)
		}

		platformName := device.SoftwareType
		if platformName == "" {
			platformName = device.PlatformID // Fallback name
		}

		platform, err := nbi.AddPlatform(&objects.Platform{
			Name:         platformName,
			Slug:         utils.Slugify(platformName),
			Manufacturer: ciscoManufacturer,
		})
		if err != nil {
			return fmt.Errorf("dnac platform: %s", err)
		}

		var deviceSite *objects.Site
		var ok bool
		if deviceSite, ok = ds.SiteId2nbSite[ds.Device2Site[device.ID]]; !ok {
			ds.Logger.Errorf("DeviceSite is not existing for device %s, this should not happen. This device will be skipped", device.ID)
			continue
		}

		if device.Type == "" {
			ds.Logger.Errorf("Device type for device %s is empty, this should not happen. This device will be skipped", device.ID)
		}

		deviceType, err := nbi.AddDeviceType(&objects.DeviceType{
			Manufacturer: ciscoManufacturer,
			Model:        device.Type,
			Slug:         utils.Slugify(device.Type),
		})
		if err != nil {
			return fmt.Errorf("add device type: %s", err)
		}

		nbDevice, err := nbi.AddDevice(&objects.Device{
			NetboxObject: objects.NetboxObject{
				Tags:        ds.SourceTags,
				Description: description,
			},
			Name:         device.Hostname,
			DeviceRole:   deviceRole,
			SerialNumber: device.SerialNumber,
			Platform:     platform,
			Comments:     comments,
			Site:         deviceSite,
			DeviceType:   deviceType,
		})

		if err != nil {
			return fmt.Errorf("adding dnac device: %s", err)
		}

		if ds.DeviceId2nbDevice == nil {
			ds.DeviceId2nbDevice = make(map[string]*objects.Device)
		}
		ds.DeviceId2nbDevice[device.ID] = nbDevice
	}

	return nil
}

func (ds *DnacSource) SyncDeviceInterfaces(nbi *inventory.NetboxInventory) error {
	return nil
}

func (ds *DnacSource) SyncVlans(nbi *inventory.NetboxInventory) error {
	return nil
}
