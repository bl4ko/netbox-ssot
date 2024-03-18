package dnac

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/constants"
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
				Tags: ds.Config.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: ds.SourceConfig.Name,
				},
			},
			Name: site.Name,
			Slug: utils.Slugify(site.Name),
		}
		for _, additionalInfo := range site.AdditionalInfo {
			if additionalInfo.Namespace == "Location" {
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
		nbSite, err := nbi.AddSite(ds.Ctx, dnacSite)
		if err != nil {
			return fmt.Errorf("adding site: %s", err)
		}
		ds.SiteID2nbSite[site.ID] = nbSite
	}
	return nil
}

func (ds *DnacSource) SyncVlans(nbi *inventory.NetboxInventory) error {
	for vid, vlan := range ds.Vlans {
		vlanGroup, err := common.MatchVlanToGroup(ds.Ctx, nbi, vlan.InterfaceName, ds.VlanGroupRelations)
		if err != nil {
			return fmt.Errorf("vlanGroup: %s", err)
		}
		vlanTenant, err := common.MatchVlanToTenant(ds.Ctx, nbi, vlan.InterfaceName, ds.VlanTenantRelations)
		if err != nil {
			return fmt.Errorf("vlanTenant: %s", err)
		}
		newVlan, err := nbi.AddVlan(ds.Ctx, &objects.Vlan{
			NetboxObject: objects.NetboxObject{
				Tags:        ds.Config.SourceTags,
				Description: vlan.VLANType,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: ds.SourceConfig.Name,
				},
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
			_, err = nbi.AddPrefix(ds.Ctx, &objects.Prefix{
				NetboxObject: objects.NetboxObject{
					Tags: ds.Config.SourceTags,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: ds.SourceConfig.Name,
					},
				},
				Prefix: prefix,
				Tenant: vlanTenant,
				Vlan:   newVlan,
			})
			if err != nil {
				return fmt.Errorf("adding prefix: %s", err)
			}
		}

		ds.VID2nbVlan[vid] = newVlan
	}
	return nil
}

func (ds *DnacSource) SyncDevices(nbi *inventory.NetboxInventory) error {
	for _, device := range ds.Devices {
		var description, comments string
		if device.Description != "" {
			description = device.Description
		}
		if len(description) > objects.MaxDescriptionLength {
			comments = description
			description = "See comments"
		}

		ciscoManufacturer, err := nbi.AddManufacturer(ds.Ctx, &objects.Manufacturer{
			Name: "Cisco",
			Slug: utils.Slugify("Cisco"),
		})
		if err != nil {
			return fmt.Errorf("failed creating device: %s", err)
		}

		deviceRole, err := nbi.AddDeviceRole(ds.Ctx, &objects.DeviceRole{
			Name:   device.Family,
			Slug:   utils.Slugify(device.Family),
			Color:  constants.ColorAqua,
			VMRole: false,
		})
		if err != nil {
			return fmt.Errorf("adding dnac device role: %s", err)
		}

		platformName := device.SoftwareType
		if platformName == "" {
			platformName = device.PlatformID // Fallback name
		}

		platform, err := nbi.AddPlatform(ds.Ctx, &objects.Platform{
			Name:         platformName,
			Slug:         utils.Slugify(platformName),
			Manufacturer: ciscoManufacturer,
		})
		if err != nil {
			return fmt.Errorf("dnac platform: %s", err)
		}

		var deviceSite *objects.Site
		var ok bool
		if deviceSite, ok = ds.SiteID2nbSite[ds.Device2Site[device.ID]]; !ok {
			ds.Logger.Errorf(ds.Ctx, "DeviceSite is not existing for device %s, this should not happen. This device will be skipped", device.ID)
			continue
		}

		if device.Type == "" {
			ds.Logger.Errorf(ds.Ctx, "Device type for device %s is empty, this should not happen. This device will be skipped", device.ID)
		}

		deviceType, err := nbi.AddDeviceType(ds.Ctx, &objects.DeviceType{
			Manufacturer: ciscoManufacturer,
			Model:        device.Type,
			Slug:         utils.Slugify(device.Type),
		})
		if err != nil {
			return fmt.Errorf("add device type: %s", err)
		}

		deviceTenant, err := common.MatchHostToTenant(ds.Ctx, nbi, device.Hostname, ds.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("hostTenant: %s", err)
		}

		deviceStatus := &objects.DeviceStatusActive
		if device.ReachabilityStatus == "Unreachable" {
			deviceStatus = &objects.DeviceStatusOffline
		}

		nbDevice, err := nbi.AddDevice(ds.Ctx, &objects.Device{
			NetboxObject: objects.NetboxObject{
				Tags:        ds.Config.SourceTags,
				Description: description,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: ds.SourceConfig.Name,
				},
			},
			Name:         device.Hostname,
			Status:       deviceStatus,
			Tenant:       deviceTenant,
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

		ds.DeviceID2nbDevice[device.ID] = nbDevice
	}
	return nil
}

func (ds *DnacSource) SyncDeviceInterfaces(nbi *inventory.NetboxInventory) error {
	for ifaceID, iface := range ds.Interfaces {
		ifaceDescription := iface.Description
		ifaceDevice := ds.DeviceID2nbDevice[iface.DeviceID]
		var ifaceDuplex *objects.InterfaceDuplex
		switch iface.Duplex {
		case "FullDuplex":
			ifaceDuplex = &objects.DuplexFull
		case "AutoNegotiate":
			ifaceDuplex = &objects.DuplexAuto
		case "HalfDuplex":
			ifaceDuplex = &objects.DuplexHalf
		case "":
		default:
			ds.Logger.Warningf(ds.Ctx, "Unknown duplex value: %s", iface.Duplex)
		}

		var ifaceStatus bool
		switch iface.Status {
		case "down":
			ifaceStatus = false
		case "up":
			ifaceStatus = true
		default:
			ds.Logger.Errorf(ds.Ctx, "wrong interface status: %s", iface.Status)
		}

		ifaceSpeed, err := strconv.Atoi(iface.Speed)
		if err != nil {
			ds.Logger.Errorf(ds.Ctx, "wrong speed for iface %s", iface.Speed)
		}

		var ifaceType *objects.InterfaceType
		switch iface.InterfaceType {
		case "Physical":
			ifaceType = objects.IfaceSpeed2IfaceType[objects.InterfaceSpeed(ifaceSpeed)]
			if ifaceType == nil {
				ifaceType = &objects.OtherInterfaceType
			}
		case "Virtual":
			ifaceType = &objects.VirtualInterfaceType
		default:
			ds.Logger.Errorf(ds.Ctx, "Unknown interface type: %s. Skipping this device...", iface.InterfaceType)
			continue
		}

		ifaceName := iface.PortName
		if ifaceName == "" {
			ds.Logger.Errorf(ds.Ctx, "Unknown interface name for iface: %s", ifaceID)
			continue
		}

		var ifaceMode *objects.InterfaceMode
		var ifaceAccessVlan *objects.Vlan
		var ifaceTrunkVlans []*objects.Vlan
		vid, err := strconv.Atoi(iface.VLANID)
		if err != nil {
			ds.Logger.Errorf(ds.Ctx, "Can't parse vid for iface %s", iface.VLANID)
			continue
		}
		switch iface.PortMode {
		case "access":
			ifaceMode = &objects.InterfaceModeAccess
			ifaceAccessVlan = ds.VID2nbVlan[vid]
		case "trunk":
			ifaceMode = &objects.InterfaceModeTagged
			// TODO: ifaceTrunkVlans = append(ifaceTrunkVlans, ds.Vid2nbVlan[vid])
		case "dynamic_auto":
			// TODO: how to handle this mode in netbox
			ds.Logger.Debugf(ds.Ctx, "vlan mode 'dynamic_auto' is not implemented yet")
		case "routed":
			ds.Logger.Debugf(ds.Ctx, "vlan mode 'routed' is not implemented yet")
		default:
			ds.Logger.Errorf(ds.Ctx, "Unknown interface mode: '%s'", iface.PortMode)
		}

		nbIface, err := nbi.AddInterface(ds.Ctx, &objects.Interface{
			NetboxObject: objects.NetboxObject{
				Description: ifaceDescription,
				Tags:        ds.Config.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: ds.SourceConfig.Name,
				},
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

		// Add IP address to the interface
		if iface.IPv4Address != "" && !utils.SubnetsContainIPAddress(iface.IPv4Address, ds.SourceConfig.IgnoredSubnets) {
			defaultMask := 32
			if iface.IPv4Mask != "" {
				maskBits, err := utils.MaskToBits(iface.IPv4Mask)
				if err != nil {
					return fmt.Errorf("wrong mask: %s", err)
				}
				defaultMask = maskBits
			}
			nbIPAddress, err := nbi.AddIPAddress(ds.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: ds.Config.SourceTags,
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: ds.SourceConfig.Name,
					},
				},
				Address:            fmt.Sprintf("%s/%d", iface.IPv4Address, defaultMask),
				Status:             &objects.IPAddressStatusActive,
				DNSName:            utils.ReverseLookup(iface.IPv4Address),
				AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
				AssignedObjectID:   nbIface.ID,
			})
			if err != nil {
				return fmt.Errorf("adding ip address: %s", err)
			}

			// To determine if this interface, has the same IP address as the device's management IP
			// we need to check if management IP is in the same subnet as this interface
			deviceManagementIP := ds.Devices[iface.DeviceID].ManagementIPAddress
			if deviceManagementIP == iface.IPv4Address {
				deviceCopy := *ifaceDevice
				deviceCopy.PrimaryIPv4 = nbIPAddress
				_, err = nbi.AddDevice(ds.Ctx, &deviceCopy)
				if err != nil {
					return fmt.Errorf("adding primary ipv4 address: %s", err)
				}
			}
		}

		ds.InterfaceID2nbInterface[ifaceID] = nbIface
	}
	return nil
}
