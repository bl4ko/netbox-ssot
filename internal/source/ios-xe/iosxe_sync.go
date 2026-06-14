package iosxe

import (
	"fmt"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	devices "github.com/src-doo/go-devicetype-library/pkg"
)

// Syncs dnac sites to netbox inventory.
func (is *IOSXESource) syncDevice(nbi *inventory.NetboxInventory) error {
	var err error
	deviceName := is.SystemInfo.Hostname
	if deviceName == "" {
		return fmt.Errorf("hostname for device is empty")
	}

	var deviceModel, serialNumber, description string
	if len(is.HardwareInfo.Inventory) > 0 {
		for _, inv := range is.HardwareInfo.Inventory {
			if inv.Type == "hw-type-chassis" {
				deviceModel = inv.PartNumber
				serialNumber = inv.SerialNumber
				description = inv.Description
			}
		}
	}
	if deviceModel == "" {
		deviceModel = constants.DefaultModel
	}
	deviceManufacturer, err := nbi.AddManufacturer(is.Ctx, &objects.Manufacturer{
		Name: "Cisco",
		Slug: utils.Slugify("Cisco"),
	})
	var deviceType *objects.DeviceType
	if deviceData, ok := devices.DeviceTypesMap[deviceManufacturer.Name][deviceModel]; ok {
		deviceType, err = nbi.AddDeviceType(is.Ctx, &objects.DeviceType{
			Manufacturer: deviceManufacturer,
			Model:        deviceModel,
			Slug:         deviceData.Slug,
		})
		if err != nil {
			return fmt.Errorf("add device type: %s", err)
		}
	} else {
		if err != nil {
			return fmt.Errorf("failed adding manufacturer: %s", err)
		}
		deviceType, err = nbi.AddDeviceType(is.Ctx, &objects.DeviceType{
			Manufacturer: deviceManufacturer,
			Model:        deviceModel,
			Slug:         utils.GenerateDeviceTypeSlug(deviceManufacturer.Name, deviceModel),
		})
		if err != nil {
			return fmt.Errorf("add device type: %s", err)
		}
	}

	deviceTenant, err := common.MatchHostToTenant(
		is.Ctx,
		nbi,
		deviceName,
		is.SourceConfig.HostTenantRelations,
	)
	if err != nil {
		return fmt.Errorf("match host to tenant: %s", err)
	}

	// Match host to a role. First test if user provided relations, if
	// not use default switch role.
	var deviceRole *objects.DeviceRole
	if len(is.SourceConfig.HostRoleRelations) > 0 {
		deviceRole, err = common.MatchHostToRole(
			is.Ctx,
			nbi,
			deviceName,
			is.SourceConfig.HostRoleRelations,
		)
		if err != nil {
			return fmt.Errorf("match host to role: %s", err)
		}
	}
	if deviceRole == nil {
		deviceRole, err = nbi.AddSwitchDeviceRole(is.Ctx)
		if err != nil {
			return fmt.Errorf("add device role: %s", err)
		}
	}

	deviceSite, err := common.MatchHostToSite(
		is.Ctx,
		nbi,
		deviceName,
		is.SourceConfig.HostSiteRelations,
	)
	if err != nil {
		return fmt.Errorf("match host to site: %s", err)
	}

	devicePlatformName := "IOS-XE" // TODO
	devicePlatform, err := nbi.AddPlatform(is.Ctx, &objects.Platform{
		Name:         devicePlatformName,
		Slug:         utils.Slugify(devicePlatformName),
		Manufacturer: deviceManufacturer,
	})
	if err != nil {
		return fmt.Errorf("add platform: %s", err)
	}
	NBDevice, err := nbi.AddDevice(is.Ctx, &objects.Device{
		NetboxObject: objects.NetboxObject{
			Tags:        is.GetSourceTags(),
			Description: description,
		},
		Name:         deviceName,
		SerialNumber: serialNumber,
		Site:         deviceSite,
		DeviceRole:   deviceRole,
		Status:       &objects.DeviceStatusActive,
		DeviceType:   deviceType,
		Tenant:       deviceTenant,
		Platform:     devicePlatform,
	})
	if err != nil {
		return fmt.Errorf("add device: %s", err)
	}
	is.NBDevice = NBDevice

	return nil
}

func (is *IOSXESource) syncInterfaces(nbi *inventory.NetboxInventory) error {
	is.NBInterfaces = make(map[string]*objects.Interface)
	for ifaceName, iface := range is.Interfaces {
		ifaceEnabled := iface.State.Enabled
		ifaceMAC := iface.Ethernet.MACAddress
		ifaceType := &objects.OtherInterfaceType
		ifaceLinkSpeed := iosxePortSpeedToLinkSpeed(iface.Ethernet.PortSpeed)
		if t, ok := objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]; ok {
			ifaceType = t
		}

		nbIface, err := nbi.AddInterface(is.Ctx, &objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags: is.GetSourceTags(),
			},
			Name:   ifaceName,
			Type:   ifaceType,
			Device: is.NBDevice,
			Speed:  ifaceLinkSpeed,
			Status: ifaceEnabled,
		})
		if err != nil {
			return fmt.Errorf("add interface: %s", err)
		}
		if ifaceMAC != "" {
			nbMACAddress, err := common.CreateMACAddressForObjectType(
				is.Ctx,
				nbi,
				ifaceMAC,
				nbIface,
			)
			if err != nil {
				return fmt.Errorf("create mac address for object type: %s", err)
			}
			if err = common.SetPrimaryMACForInterface(is.Ctx, nbi, nbIface, nbMACAddress); err != nil {
				return fmt.Errorf("set primary mac for interface: %s", err)
			}
		}
		is.NBInterfaces[ifaceName] = nbIface
	}
	return nil
}

func (is *IOSXESource) syncArpTable(nbi *inventory.NetboxInventory) error {
	if !is.SourceConfig.CollectArpData {
		is.Logger.Info(is.Ctx, "skipping collecting of arp data")
		return nil
	}

	// We tag it with special tag for arp data.
	arpTag, err := nbi.AddTag(is.Ctx, &objects.Tag{
		Name:        constants.DefaultArpTagName,
		Slug:        utils.Slugify(constants.DefaultArpTagName),
		Color:       constants.DefaultArpTagColor,
		Description: "tag created for ip's collected from arp table",
	})
	if err != nil {
		return fmt.Errorf("add tag: %s", err)
	}

	for _, arpEntry := range is.ArpEntries {
		if utils.IsPermittedIPAddress(
			arpEntry.Address,
			is.SourceConfig.PermittedSubnets,
			is.SourceConfig.IgnoredSubnets,
		) {
			newTags := is.GetSourceTags()
			newTags = append(newTags, arpTag)
			currentTime := time.Now()
			dnsName := utils.ReverseLookup(arpEntry.Address)
			defaultMask := 32
			addressWithMask := fmt.Sprintf("%s/%d", arpEntry.Address, defaultMask)
			// VRF
			ipVRF, err := common.MatchIPToVRF(is.Ctx, nbi, arpEntry.Address, is.SourceConfig.IPVrfRelations)
			if err != nil {
				is.Logger.Warningf(is.Ctx, "match ip to vrf for %s: %s", arpEntry.Address, err)
			}
			_, err = nbi.AddIPAddress(is.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: newTags,
					Description: fmt.Sprintf(
						"IP collected from %s arp table",
						is.SourceConfig.Name,
					),
					CustomFields: map[string]interface{}{
						constants.CustomFieldOrphanLastSeenName: currentTime.Format(
							constants.CustomFieldOrphanLastSeenFormat,
						),
						constants.CustomFieldArpEntryName: true,
					},
				},
				Address: addressWithMask,
				DNSName: dnsName,
				Status:  &objects.IPAddressStatusActive,
				VRF:     ipVRF,
			})
			if err != nil {
				is.Logger.Warningf(is.Ctx, "error creating ip address: %s", err)
			}
		}
	}
	return nil
}

// iosxePortSpeedToLinkSpeed maps a Cisco IOS-XE Ethernet PortSpeed enum to the
// netbox interface speed in kbps. It returns 0 for unknown speeds.
func iosxePortSpeedToLinkSpeed(portSpeed string) objects.InterfaceSpeed {
	switch portSpeed {
	case "SPEED_10MB":
		return (10 * constants.MB) / constants.KB //nolint:mnd
	case "SPEED_100MB":
		return (100 * constants.MB) / constants.KB //nolint:mnd
	case "SPEED_1GB":
		return (1 * constants.GB) / constants.KB
	case "SPEED_10GB":
		return (10 * constants.GB) / constants.KB //nolint:mnd
	case "SPEED_25GB":
		return (25 * constants.GB) / constants.KB //nolint:mnd
	case "SPEED_40GB":
		return (40 * constants.GB) / constants.KB //nolint:mnd
	case "SPEED_100GB":
		return (100 * constants.GB) / constants.KB //nolint:mnd
	default:
		return 0
	}
}
