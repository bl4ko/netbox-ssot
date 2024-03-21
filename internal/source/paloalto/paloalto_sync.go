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

// Sync device creates default device in netbox representing Paloalto firewall.
func (pas *PaloAltoSource) SyncDevice(nbi *inventory.NetboxInventory) error {
	deviceName := pas.SystemInfo["devicename"]
	if deviceName == "" {
		return fmt.Errorf("can't extract device name from system info")
	}
	deviceSerialNumber := pas.SystemInfo["serial"]
	deviceModel := pas.SystemInfo["model"]
	if deviceModel == "" {
		pas.Logger.Warningf(pas.Ctx, "model field in system info is empty. Using fallback mechanism.")
		deviceModel = constants.DefaultModel
	}
	deviceManufacturer, err := nbi.AddManufacturer(pas.Ctx, &objects.Manufacturer{
		Name: "Palo Alto Networks, Inc.",
		Slug: utils.Slugify("Palo Alto Networks, Inc."),
	})
	if err != nil {
		return fmt.Errorf("failed adding manufacturer: %s", err)
	}
	deviceType, err := nbi.AddDeviceType(pas.Ctx, &objects.DeviceType{
		Manufacturer: deviceManufacturer,
		Model:        deviceModel,
		Slug:         utils.Slugify(deviceManufacturer.Name + deviceModel),
	})
	if err != nil {
		return fmt.Errorf("add device type: %s", err)
	}

	deviceTenant, err := common.MatchHostToTenant(pas.Ctx, nbi, deviceName, pas.HostTenantRelations)
	if err != nil {
		return fmt.Errorf("match host to tenant: %s", err)
	}

	deviceRole, err := nbi.AddDeviceRole(pas.Ctx, &objects.DeviceRole{
		Name:   constants.DeviceRoleFirewall,
		Slug:   utils.Slugify(constants.DeviceRoleFirewall),
		Color:  constants.DeviceRoleFirewallColor,
		VMRole: false,
	})
	if err != nil {
		return fmt.Errorf("add DeviceRole: %s", err)
	}
	deviceSite, err := common.MatchHostToSite(pas.Ctx, nbi, deviceName, pas.HostSiteRelations)
	if err != nil {
		return fmt.Errorf("match host to site: %s", err)
	}
	devicePlatformName := fmt.Sprintf("PAN-OS %s", pas.SystemInfo["sw-version"])
	devicePlatform, err := nbi.AddPlatform(pas.Ctx, &objects.Platform{
		Name:         devicePlatformName,
		Slug:         utils.Slugify(devicePlatformName),
		Manufacturer: deviceManufacturer,
	})
	if err != nil {
		return fmt.Errorf("add platform: %s", err)
	}
	NBDevice, err := nbi.AddDevice(pas.Ctx, &objects.Device{
		NetboxObject: objects.NetboxObject{
			Tags: pas.SourceTags,
		},
		Name:         deviceName,
		Site:         deviceSite,
		DeviceRole:   deviceRole,
		Status:       &objects.DeviceStatusActive,
		DeviceType:   deviceType,
		Tenant:       deviceTenant,
		Platform:     devicePlatform,
		SerialNumber: deviceSerialNumber,
	})
	if err != nil {
		return fmt.Errorf("add device: %s", err)
	}

	pas.NBFirewall = NBDevice
	return nil
}

func (pas *PaloAltoSource) SyncInterfaces(nbi *inventory.NetboxInventory) error {
	for _, iface := range pas.Ifaces {
		if iface.Name == "" {
			pas.Logger.Debugf(pas.Ctx, "empty interface name. Skipping...")
			continue
		}
		if utils.FilterInterfaceName(iface.Name, pas.SourceConfig.InterfaceFilter) {
			pas.Logger.Debugf(pas.Ctx, "interface %s is filtered out with interface filter %s", iface.Name, pas.SourceConfig.InterfaceFilter)
			continue
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

		nbIface, err := nbi.AddInterface(pas.Ctx, &objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        pas.SourceTags,
				Description: iface.Comment,
			},
			Name:   iface.Name,
			Type:   ifaceType,
			Duplex: ifaceDuplex,
			Device: pas.NBFirewall,
			MTU:    iface.Mtu,
			Speed:  ifaceLinkSpeed,
			Vdcs:   []*objects.VirtualDeviceContext{pas.getVirtualDeviceContext(nbi, iface.Name)},
		})
		if err != nil {
			return fmt.Errorf("add interface %s", err)
		}

		if len(iface.StaticIps) > 0 {
			pas.syncIPs(nbi, nbIface, iface.StaticIps)
		}

		ifaceVlans := []*objects.Vlan{}
		for _, subIface := range pas.Iface2SubIfaces[iface.Name] {
			subIfaceName := subIface.Name
			if subIfaceName == "" {
				continue
			}
			var subIfaceVlan *objects.Vlan
			if subIface.Tag != 0 {
				// Extract Vlan
				vlanGroup, err := common.MatchVlanToGroup(pas.Ctx, nbi, fmt.Sprintf("Vlan%d", subIface.Tag), pas.VlanGroupRelations)
				if err != nil {
					return fmt.Errorf("match vlan to group: %s", err)
				}
				subIfaceVlan, err = nbi.AddVlan(pas.Ctx, &objects.Vlan{
					NetboxObject: objects.NetboxObject{
						Tags:        pas.SourceTags,
						Description: subIface.Comment,
					},
					Status: &objects.VlanStatusActive,
					Name:   fmt.Sprintf("Vlan%d", subIface.Tag),
					Vid:    subIface.Tag,
					Group:  vlanGroup,
				})
				if err != nil {
					return fmt.Errorf("add vlan: %s", err)
				}
				ifaceVlans = append(ifaceVlans, subIfaceVlan)
			}
			nbSubIface, err := nbi.AddInterface(pas.Ctx, &objects.Interface{
				NetboxObject: objects.NetboxObject{
					Tags:        pas.SourceTags,
					Description: subIface.Comment,
				},
				Name:            subIface.Name,
				Type:            &objects.VirtualInterfaceType,
				Device:          pas.NBFirewall,
				Mode:            &objects.InterfaceModeTagged,
				TaggedVlans:     []*objects.Vlan{subIfaceVlan},
				ParentInterface: nbIface,
				MTU:             subIface.Mtu,
				Vdcs:            []*objects.VirtualDeviceContext{pas.getVirtualDeviceContext(nbi, subIfaceName)},
			})
			if err != nil {
				return fmt.Errorf("add subinterface: %s", err)
			}
			if len(subIface.StaticIps) > 0 {
				pas.syncIPs(nbi, nbSubIface, subIface.StaticIps)
			}
		}

		if len(ifaceVlans) > 0 {
			nbIfaceUpdate := *nbIface
			nbIfaceUpdate.Mode = &objects.InterfaceModeTagged
			nbIfaceUpdate.TaggedVlans = ifaceVlans
			_, err = nbi.AddInterface(pas.Ctx, &nbIfaceUpdate)
			if err != nil {
				pas.Logger.Errorf(pas.Ctx, "updating ifaceVlans: %s", err)
			}
		}
	}
	return nil
}

func (pas *PaloAltoSource) syncIPs(nbi *inventory.NetboxInventory, nbIface *objects.Interface, ips []string) {
	for _, ipAddress := range ips {
		if !utils.SubnetsContainIPAddress(ipAddress, pas.SourceConfig.IgnoredSubnets) {
			dnsName := utils.ReverseLookup(ipAddress)
			_, err := nbi.AddIPAddress(pas.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: pas.SourceTags,
				},
				Address:            ipAddress,
				AssignedObjectID:   nbIface.ID,
				DNSName:            dnsName,
				AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
			})
			if err != nil {
				pas.Logger.Errorf(pas.Ctx, "adding ip address %s failed with error: %s", ipAddress, err)
				continue
			}
			prefix, err := utils.ExtractPrefixFromIPAddress(ipAddress)
			if err != nil {
				pas.Logger.Warningf(pas.Ctx, "extract prefix from address: %s", err)
			} else {
				_, err = nbi.AddPrefix(pas.Ctx, &objects.Prefix{
					Prefix: prefix,
				})
				if err != nil {
					pas.Logger.Errorf(pas.Ctx, "adding prefix: %s", err)
				}
			}
		}
	}
}

func (pas *PaloAltoSource) SyncSecurityZones(nbi *inventory.NetboxInventory) error {
	for _, securityZone := range pas.SecurityZones {
		_, err := nbi.AddVirtualDeviceContext(pas.Ctx, &objects.VirtualDeviceContext{
			NetboxObject: objects.NetboxObject{
				Tags: pas.SourceTags,
			},
			Name:   securityZone.Name,
			Device: pas.NBFirewall,
			Status: &objects.VDCStatusActive,
		})
		if err != nil {
			return fmt.Errorf("add VirtualDeviceContext: %s", err)
		}
	}
	return nil
}

func (pas *PaloAltoSource) getVirtualDeviceContext(nbi *inventory.NetboxInventory, ifaceName string) *objects.VirtualDeviceContext {
	var virtualDeviceContext *objects.VirtualDeviceContext
	zoneName := pas.Iface2SecurityZone[ifaceName]
	if vdc, ok := nbi.VirtualDeviceContextsIndexByNameAndDeviceID[zoneName][pas.NBFirewall.ID]; ok {
		virtualDeviceContext = vdc
	}
	return virtualDeviceContext
}
