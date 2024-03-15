package paloalto

import (
	"fmt"
	"strconv"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

func (pas *PaloAltoSource) SyncDevices(nbi *inventory.NetboxInventory) error {
	pas.DeviceName2NbDevice = make(map[string]*objects.Device)
	for _, device := range pas.Devices {
		deviceRole, err := nbi.AddDeviceRole(pas.Ctx, &objects.DeviceRole{
			Name:   "VirtualRouter",
			Slug:   utils.Slugify("VirtualRouter"),
			Color:  constants.ColorOrange,
			VMRole: false,
		})
		if err != nil {
			return fmt.Errorf("create device role: %s", err)
		}
		deviceSite, err := common.MatchHostToSite(pas.Ctx, nbi, device.Name, pas.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("match host to site: %s", err)
		}
		deviceTenant, err := common.MatchHostToTenant(pas.Ctx, nbi, device.Name, pas.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("match host to tenant: %s", err)
		}
		genericManufacturer, err := nbi.AddManufacturer(pas.Ctx, &objects.Manufacturer{
			Name: constants.DefaultManufacturer,
			Slug: utils.Slugify(constants.DefaultManufacturer),
		})
		if err != nil {
			return fmt.Errorf("new manufacturer: %s", err)
		}
		platformName := utils.GeneratePlatformName(constants.DefaultOSName, constants.DefaultOSVersion)
		platform, err := nbi.AddPlatform(pas.Ctx, &objects.Platform{
			Name:         platformName,
			Slug:         utils.Slugify(platformName),
			Manufacturer: genericManufacturer,
		})
		if err != nil {
			return fmt.Errorf("adding platform: %s", err)
		}
		deviceType, err := nbi.AddDeviceType(pas.Ctx, &objects.DeviceType{
			Manufacturer: genericManufacturer,
			Model:        constants.DefaultModel,
			Slug:         utils.Slugify(genericManufacturer.Name + constants.DefaultModel),
		})
		if err != nil {
			return fmt.Errorf("add device type: %s", err)
		}
		nbDevice, err := nbi.AddDevice(pas.Ctx, &objects.Device{
			NetboxObject: objects.NetboxObject{
				Tags: pas.Config.SourceTags,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: pas.SourceConfig.Name,
				},
			},
			Name:       device.Name,
			Site:       deviceSite,
			Tenant:     deviceTenant,
			Status:     &objects.DeviceStatusActive,
			Platform:   platform,
			DeviceRole: deviceRole,
			DeviceType: deviceType,
		})
		if err != nil {
			return fmt.Errorf("add device: %s", err)
		}
		pas.DeviceName2NbDevice[device.Name] = nbDevice
	}
	return nil
}

func (pas *PaloAltoSource) SyncInterfaces(nbi *inventory.NetboxInventory) error {
	for _, iface := range pas.Interfaces {
		if iface.Name == "" {
			pas.Logger.Debugf(pas.Ctx, "empty interface name. Skipping...")
			continue
		}
		if utils.FilterInterfaceName(iface.Name, pas.SourceConfig.InterfaceFilter) {
			pas.Logger.Debugf(pas.Ctx, "interface %s is filtered out with interface filter %s", iface.Name, pas.SourceConfig.InterfaceFilter)
			continue
		}
		ifaceDeviceName, ok := pas.Interface2Router[iface.Name]
		if !ok {
			pas.Logger.Warningf(pas.Ctx, "no matched device for iface %s", iface.Name)
		}
		var ifaceLinkSpeed objects.InterfaceSpeed
		ifaceType := &objects.OtherInterfaceType
		if iface.LinkSpeed != "" {
			speed, _ := strconv.Atoi(iface.LinkSpeed)
			ifaceLinkSpeed = objects.InterfaceSpeed(speed)
			if _, ok := objects.IfaceSpeed2IfaceType[objects.InterfaceSpeed(speed)]; ok {
				ifaceType = objects.IfaceSpeed2IfaceType[objects.InterfaceSpeed(speed)]
			}
		}
		var ifaceDuplex *objects.InterfaceDuplex
		if iface.LinkDuplex != "" {
			switch iface.LinkDuplex {
			case "full":
				ifaceDuplex = &objects.DuplexFull
			case "auto":
				ifaceDuplex = &objects.DuplexAuto
			case "half":
				ifaceDuplex = &objects.DuplexHalf
			case "":
			default:
				pas.Logger.Debugf(pas.Ctx, "not implemented duplex value %s", iface.LinkDuplex)
			}
		}

		ifaceDevice := pas.DeviceName2NbDevice[ifaceDeviceName]
		nbIface, err := nbi.AddInterface(pas.Ctx, &objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        pas.SourceTags,
				Description: iface.Comment,
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: pas.SourceConfig.Name,
				},
			},
			Name:   iface.Name,
			Type:   ifaceType,
			Duplex: ifaceDuplex,
			Device: ifaceDevice,
			MTU:    iface.Mtu,
			Speed:  ifaceLinkSpeed,
		})
		if err != nil {
			return fmt.Errorf("add interface %s", err)
		}

		if len(iface.StaticIps) > 0 {
			err := pas.syncIPs(nbi, nbIface, iface.StaticIps)
			if err != nil {
				return fmt.Errorf("interface ip sync: %s", err)
			}
		}
	}
	return nil
}

func (pas *PaloAltoSource) syncIPs(nbi *inventory.NetboxInventory, nbIface *objects.Interface, ips []string) error {
	for _, ipAddress := range ips {
		_, err := nbi.AddIPAddress(pas.Ctx, &objects.IPAddress{
			Address:            ipAddress,
			AssignedObjectID:   nbIface.ID,
			AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
		})
		if err != nil {
			return fmt.Errorf("add ip address: %s", err)
		}
	}
	return nil
}
