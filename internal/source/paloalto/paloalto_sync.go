package paloalto

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Sync device creates default device in netbox representing Paloalto firewall.
func (pas *PaloAltoSource) syncDevice(nbi *inventory.NetboxInventory) error {
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
		Name: "Palo Alto",
		Slug: utils.Slugify("Palo Alto"),
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

func (pas *PaloAltoSource) syncInterfaces(nbi *inventory.NetboxInventory) error {
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

		var ifaceVdcs []*objects.VirtualDeviceContext
		if vdc := pas.getVirtualDeviceContext(nbi, iface.Name); vdc != nil {
			ifaceVdcs = []*objects.VirtualDeviceContext{vdc}
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
			Vdcs:   ifaceVdcs,
		})
		if err != nil {
			return fmt.Errorf("add interface %s", err)
		}

		if len(iface.StaticIps) > 0 {
			pas.syncIPs(nbi, nbIface, iface.StaticIps, nil)
		}

		for _, subIface := range pas.Iface2SubIfaces[iface.Name] {
			subIfaceName := subIface.Name
			if subIfaceName == "" {
				continue
			}
			var subIfaceVlan *objects.Vlan
			subIfaceVlans := []*objects.Vlan{}
			var subifaceMode *objects.InterfaceMode
			if subIface.Tag != 0 {
				// Extract Vlan
				vlanName := fmt.Sprintf("Vlan%d", subIface.Tag)
				vlanGroup, err := common.MatchVlanToGroup(pas.Ctx, nbi, vlanName, pas.VlanGroupRelations)
				if err != nil {
					return fmt.Errorf("match vlan to group: %s", err)
				}
				vlanTenant, err := common.MatchVlanToTenant(pas.Ctx, nbi, vlanName, pas.VlanTenantRelations)
				if err != nil {
					return fmt.Errorf("match vlan to tenant: %s", err)
				}
				subIfaceVlan, err = nbi.AddVlan(pas.Ctx, &objects.Vlan{
					NetboxObject: objects.NetboxObject{
						Tags:        pas.SourceTags,
						Description: subIface.Comment,
					},
					Status: &objects.VlanStatusActive,
					Name:   fmt.Sprintf("Vlan%d", subIface.Tag),
					Vid:    subIface.Tag,
					Tenant: vlanTenant,
					Group:  vlanGroup,
				})
				if err != nil {
					return fmt.Errorf("add vlan: %s", err)
				}
				subIfaceVlans = append(subIfaceVlans, subIfaceVlan)
				subifaceMode = &objects.InterfaceModeTagged
			}
			var vdcs []*objects.VirtualDeviceContext
			if vdc := pas.getVirtualDeviceContext(nbi, subIfaceName); vdc != nil {
				vdcs = []*objects.VirtualDeviceContext{vdc}
			}
			nbSubIface, err := nbi.AddInterface(pas.Ctx, &objects.Interface{
				NetboxObject: objects.NetboxObject{
					Tags:        pas.SourceTags,
					Description: subIface.Comment,
				},
				Name:            subIface.Name,
				Type:            &objects.VirtualInterfaceType,
				Device:          pas.NBFirewall,
				Mode:            subifaceMode,
				TaggedVlans:     subIfaceVlans,
				ParentInterface: nbIface,
				MTU:             subIface.Mtu,
				Vdcs:            vdcs,
			})
			if err != nil {
				return fmt.Errorf("add subinterface: %s", err)
			}
			if len(subIface.StaticIps) > 0 {
				pas.syncIPs(nbi, nbSubIface, subIface.StaticIps, subIfaceVlan)
			}
		}
	}
	return nil
}

// syncIPs adds all of the given ips to the given nbIface. It also
// Extracts prefixes from ips and connect them with prefix vlan.
func (pas *PaloAltoSource) syncIPs(nbi *inventory.NetboxInventory, nbIface *objects.Interface, ips []string, prefixVlan *objects.Vlan) {
	for _, ipAddress := range ips {
		if !utils.SubnetsContainIPAddress(ipAddress, pas.SourceConfig.IgnoredSubnets) {
			dnsName := utils.ReverseLookup(ipAddress)
			_, err := nbi.AddIPAddress(pas.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: pas.SourceTags,
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: false,
					},
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
			prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(ipAddress)
			if err != nil {
				pas.Logger.Warningf(pas.Ctx, "extract prefix from address: %s", err)
			} else if mask != constants.MaxIPv4MaskBits {
				var prefixTenant *objects.Tenant
				if prefixVlan != nil {
					prefixTenant = prefixVlan.Tenant
				}
				_, err = nbi.AddPrefix(pas.Ctx, &objects.Prefix{
					Prefix: prefix,
					Tenant: prefixTenant,
					Vlan:   prefixVlan,
				})
				if err != nil {
					pas.Logger.Errorf(pas.Ctx, "adding prefix: %s", err)
				}
			}
		}
	}
}

// syncSecurityZones syncs all security zones from palo alto as virtual device context in netbox.
// They are all added as part of main paloalto firewall device.
func (pas *PaloAltoSource) syncSecurityZones(nbi *inventory.NetboxInventory) error {
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

// getVirtualDeviceContext retrieves the virtual device context associated with the given interface name.
func (pas *PaloAltoSource) getVirtualDeviceContext(nbi *inventory.NetboxInventory, ifaceName string) *objects.VirtualDeviceContext {
	var virtualDeviceContext *objects.VirtualDeviceContext
	zoneName := pas.Iface2SecurityZone[ifaceName]
	if vdc, ok := nbi.VirtualDeviceContextsIndexByNameAndDeviceID[zoneName][pas.NBFirewall.ID]; ok {
		virtualDeviceContext = vdc
	}
	return virtualDeviceContext
}

func (pas *PaloAltoSource) syncArpTable(nbi *inventory.NetboxInventory) error {
	if !pas.SourceConfig.CollectArpData {
		pas.Logger.Info(pas.Ctx, "skipping collecting of arp data")
		return nil
	}

	// We tag it with special tag for arp data.
	arpTag, err := nbi.AddTag(pas.Ctx, &objects.Tag{
		Name:        constants.DefaultArpTagName,
		Slug:        utils.Slugify(constants.DefaultArpTagName),
		Color:       constants.DefaultArpTagColor,
		Description: "tag created for ip's collected from arp table",
	})
	if err != nil {
		return fmt.Errorf("add tag: %s", err)
	}
	// We create custom field for tracking when was arp entry last seen
	_, err = nbi.AddCustomField(pas.Ctx, &objects.CustomField{
		Name:                  constants.CustomFieldArpIPLastSeenName,
		Label:                 constants.CustomFieldArpIPLastSeenLabel,
		Type:                  objects.CustomFieldTypeText,
		FilterLogic:           objects.FilterLogicLoose,
		CustomFieldUIVisible:  &objects.CustomFieldUIVisibleAlways,
		CustomFieldUIEditable: &objects.CustomFieldUIEditableYes,
		DisplayWeight:         objects.DisplayWeightDefault,
		Description:           constants.CustomFieldArpIPLastSeenDescription,
		SearchWeight:          objects.SearchWeightDefault,
		ObjectTypes:           []objects.ObjectType{objects.ObjectTypeIpamIPAddress},
	})
	if err != nil {
		return fmt.Errorf("add custom field: %s", err)
	}
	for _, entry := range pas.ArpData {
		if entry.MAC != "(incomplete)" {
			newTags := pas.SourceTags
			newTags = append(newTags, arpTag)
			currentTime := time.Now()
			dnsName := utils.ReverseLookup(entry.IP)
			defaultMask := 32
			addressWithMask := fmt.Sprintf("%s/%d", entry.IP, defaultMask)
			_, err = nbi.AddIPAddress(pas.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags:        newTags,
					Description: fmt.Sprintf("IP collected from %s arp table", pas.SourceConfig.Name),
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpIPLastSeenName: currentTime.Format(constants.ArpLastSeenFormat),
						constants.CustomFieldArpEntryName:      true,
					},
				},
				Address: addressWithMask,
				DNSName: dnsName,
				Status:  &objects.IPAddressStatusActive,
			})
			if err != nil {
				return fmt.Errorf("add arp ip address: %s", err)
			}
		}
	}
	return nil
}
