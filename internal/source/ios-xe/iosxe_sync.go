package iosxe

import (
	"fmt"
	"strings"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Syncs dnac sites to netbox inventory.
func (is *IOSXESource) syncDevice(nbi *inventory.NetboxInventory) error {
	deviceName := is.SystemInfo.Hostname
	if deviceName == "" {
		return fmt.Errorf("hostname for device is empty")
	}
	deviceModel := constants.DefaultModel
	deviceManufacturer, err := nbi.AddManufacturer(is.Ctx, &objects.Manufacturer{
		Name: "Cisco",
		Slug: utils.Slugify("Cisco"),
	})
	if err != nil {
		return fmt.Errorf("failed adding manufacturer: %s", err)
	}
	deviceType, err := nbi.AddDeviceType(is.Ctx, &objects.DeviceType{
		Manufacturer: deviceManufacturer,
		Model:        deviceModel,
		Slug:         utils.Slugify(deviceManufacturer.Name + deviceModel),
	})
	if err != nil {
		return fmt.Errorf("add device type: %s", err)
	}
	deviceTenant, err := common.MatchHostToTenant(is.Ctx, nbi, deviceName, is.HostTenantRelations)
	if err != nil {
		return fmt.Errorf("match host to tenant: %s", err)
	}

	deviceRole, err := nbi.AddSwitchDeviceRole(is.Ctx)
	if err != nil {
		return fmt.Errorf("add device role: %s", err)
	}

	deviceSite, err := common.MatchHostToSite(is.Ctx, nbi, deviceName, is.HostSiteRelations)
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
			Tags: is.SourceTags,
		},
		Name:       deviceName,
		Site:       deviceSite,
		DeviceRole: deviceRole,
		Status:     &objects.DeviceStatusActive,
		DeviceType: deviceType,
		Tenant:     deviceTenant,
		Platform:   devicePlatform,
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
		var ifaceLinkSpeed objects.InterfaceSpeed
		switch iface.Ethernet.PortSpeed {
		case "SPEED_10MB":
			ifaceLinkSpeed = (10 * constants.MB) / constants.KB //nolint:mnd
			if _, ok := objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]; ok {
				ifaceType = objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]
			}
		case "SPEED_100MB":
			ifaceLinkSpeed = (100 * constants.MB) / constants.KB //nolint:mnd
			if _, ok := objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]; ok {
				ifaceType = objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]
			}
		case "SPEED_1GB":
			ifaceLinkSpeed = (1 * constants.GB) / constants.KB
			if _, ok := objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]; ok {
				ifaceType = objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]
			}
		case "SPEED_10GB":
			ifaceLinkSpeed = (10 * constants.GB) / constants.KB //nolint:mnd
			if _, ok := objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]; ok {
				ifaceType = objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]
			}
		case "SPEED_25GB":
			ifaceLinkSpeed = (25 * constants.GB) / constants.KB //nolint:mnd
			if _, ok := objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]; ok {
				ifaceType = objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]
			}
		case "SPEED_40GB":
			ifaceLinkSpeed = (40 * constants.GB) / constants.KB //nolint:mnd
			if _, ok := objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]; ok {
				ifaceType = objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]
			}
		case "SPEED_100GB":
			ifaceLinkSpeed = (100 * constants.GB) / constants.KB //nolint:mnd
			if _, ok := objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]; ok {
				ifaceType = objects.IfaceSpeed2IfaceType[ifaceLinkSpeed]
			}
		default:
		}

		nbIface, err := nbi.AddInterface(is.Ctx, &objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags: is.SourceTags,
			},
			Name:   ifaceName,
			Type:   ifaceType,
			Device: is.NBDevice,
			MAC:    strings.ToUpper(ifaceMAC),
			Speed:  ifaceLinkSpeed,
			Status: ifaceEnabled,
		})
		if err != nil {
			return fmt.Errorf("add interface: %s", err)
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
		if !utils.SubnetsContainIPAddress(arpEntry.Address, is.SourceConfig.IgnoredSubnets) {
			newTags := is.SourceTags
			newTags = append(newTags, arpTag)
			currentTime := time.Now()
			dnsName := utils.ReverseLookup(arpEntry.Address)
			defaultMask := 32
			addressWithMask := fmt.Sprintf("%s/%d", arpEntry.Address, defaultMask)
			_, err := nbi.AddIPAddress(is.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags:        newTags,
					Description: fmt.Sprintf("IP collected from %s arp table", is.SourceConfig.Name),
					CustomFields: map[string]interface{}{
						constants.CustomFieldOrphanLastSeenName: currentTime.Format(constants.CustomFieldOrphanLastSeenFormat),
						constants.CustomFieldArpEntryName:       true,
					},
				},
				Address: addressWithMask,
				DNSName: dnsName,
				Status:  &objects.IPAddressStatusActive,
			})
			if err != nil {
				is.Logger.Warningf(is.Ctx, "error creating ip address: %s", err)
			}
		}
	}
	return nil
}
