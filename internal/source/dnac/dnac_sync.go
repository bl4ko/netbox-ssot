package dnac

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
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
		ds.SiteId2nbSite[site.ID] = nbSite
	}
	return nil
}

func (ds *DnacSource) SyncVlans(nbi *inventory.NetboxInventory) error {
	for vid, vlan := range ds.Vlans {
		vlanGroup, err := common.MatchVlanToGroup(nbi, vlan.InterfaceName, ds.VlanGroupRelations)
		if err != nil {
			return fmt.Errorf("vlanGroup: %s", err)
		}
		vlanTenant, err := common.MatchVlanToTenant(nbi, vlan.InterfaceName, ds.VlanTenantRelations)
		if err != nil {
			return fmt.Errorf("vlanTenant: %s", err)
		}
		newVlan, err := nbi.AddVlan(&objects.Vlan{
			NetboxObject: objects.NetboxObject{
				Tags: ds.SourceTags,
			},
			Name:   vlan.InterfaceName,
			Group:  vlanGroup,
			Vid:    vid,
			Tenant: vlanTenant,
		})
		if err != nil {
			return fmt.Errorf("adding vlan: %s", err)
		}

		if vlan.Prefix != "" && vlan.NetworkAddress != "" {
			// Create prefix for this vlan
			prefix := fmt.Sprintf("%s/%s", vlan.NetworkAddress, vlan.Prefix)
			_, err = nbi.AddPrefix(&objects.Prefix{
				NetboxObject: objects.NetboxObject{
					Tags: ds.SourceTags,
				},
				Prefix: prefix,
				Tenant: vlanTenant,
				Vlan:   newVlan,
			})
			if err != nil {
				return fmt.Errorf("adding prefix: %s", err)
			}
		}
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

		ds.DeviceId2nbDevice[device.ID] = nbDevice
	}

	return nil
}

func (ds *DnacSource) SyncDeviceInterfaces(nbi *inventory.NetboxInventory) error {
	for ifaceId, iface := range ds.Interfaces {
		ifaceDescription := iface.Description
		ifaceDevice := ds.DeviceId2nbDevice[iface.DeviceID]
		var ifaceDuplex *objects.InterfaceDuplex
		if iface.Duplex != "" {
			switch iface.Duplex {
			case "FullDuplex":
				ifaceDuplex = &objects.DuplexFull
			case "AutoNegotiate":
				ifaceDuplex = &objects.DuplexAuto
			case "HalfDuplex":
				ifaceDuplex = &objects.DuplexHalf
			default:
				ds.Logger.Errorf("Wrong duplex value: %s", iface.Duplex)
			}

			var ifaceStatus bool
			switch iface.Status {
			case "down":
				ifaceStatus = false
			case "up":
				ifaceStatus = true
			default:
				ds.Logger.Errorf("wrong interface status: %s", iface.Status)
			}

			ifaceSpeed, err := strconv.Atoi(iface.Speed)
			if err != nil {
				ds.Logger.Errorf("wrong speed for iface %s", iface.Speed)
			}

			var ifaceType *objects.InterfaceType
			switch iface.InterfaceType {
			case "Physical":
				ifaceType = &objects.OtherInterfaceType // TODO: get from speed
			case "Virtual":
				ifaceType = &objects.VirtualInterfaceType
			default:
				ds.Logger.Errorf("Unknown interface type: %s. Skipping this device...", iface.InterfaceType)
				continue
			}

			ifaceName := iface.PortName
			if ifaceName == "" {
				ds.Logger.Errorf("Unknown interface name for iface: %s", ifaceId)
				continue
			}

			var ifaceMode *objects.InterfaceMode
			var ifaceAccessVlan *objects.Vlan
			var ifaceTrunkVlans []*objects.Vlan
			switch iface.PortMode {
			case "access":
				ifaceMode = &objects.InterfaceModeAccess
			case "trunk":
				ifaceMode = &objects.InterfaceModeTagged
			case "dynamic_auto":
				// TODO: how to handle this mode in netbox
			default:
				ds.Logger.Errorf("Unknown interface mode: %s. Skipping this device...", iface.PortMode)
			}

			nbIface, err := nbi.AddInterface(&objects.Interface{
				NetboxObject: objects.NetboxObject{
					Description: ifaceDescription,
					Tags:        ds.SourceTags,
				},
				Name:         iface.PortName,
				MAC:          strings.ToUpper(iface.MacAddress),
				Speed:        objects.InterfaceSpeed(ifaceSpeed),
				Status:       ifaceStatus,
				Duplex:       ifaceDuplex,
				Device:       ifaceDevice,
				Type:         ifaceType,
				Mode:         ifaceMode,
				UntaggedVlan: ifaceAccessVlan,
				TaggedVlans:  ifaceTrunkVlans,
			})
			if err != nil {
				return fmt.Errorf("add device interface: %s", err)
			}
			ds.InterfaceId2nbInterface[ifaceId] = nbIface
		}
	}
	return nil
}
