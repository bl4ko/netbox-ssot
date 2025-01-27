package fortigate

import (
	"fmt"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// SyncDevice creates default device in netbox representing Fortigate firewall.
func (fs *FortigateSource) syncDevice(nbi *inventory.NetboxInventory) error {
	deviceName := fs.SystemInfo.Hostname
	if deviceName == "" {
		return fmt.Errorf("can't extract hostname from system info")
	}
	var deviceSerialNumber string
	if !fs.SourceConfig.IgnoreSerialNumbers {
		deviceSerialNumber = fs.SystemInfo.Serial
	}

	deviceModel := fs.SystemInfo.Hostname
	if deviceModel == "" {
		fs.Logger.Warningf(fs.Ctx, "model field in system info is empty. Using fallback mechanism.")
		deviceModel = constants.DefaultModel
	}
	deviceManufacturer, err := nbi.AddManufacturer(fs.Ctx, &objects.Manufacturer{
		Name: "Fortinet",
		Slug: utils.Slugify("Fortinet"),
	})
	if err != nil {
		return fmt.Errorf("failed adding manufacturer: %s", err)
	}
	deviceType, err := nbi.AddDeviceType(fs.Ctx, &objects.DeviceType{
		Manufacturer: deviceManufacturer,
		Model:        deviceModel,
		Slug:         utils.Slugify(deviceManufacturer.Name + deviceModel),
	})
	if err != nil {
		return fmt.Errorf("add device type: %s", err)
	}

	deviceTenant, err := common.MatchHostToTenant(fs.Ctx, nbi, deviceName, fs.SourceConfig.HostTenantRelations)
	if err != nil {
		return fmt.Errorf("match host to tenant: %s", err)
	}

	var deviceRole *objects.DeviceRole
	if len(fs.SourceConfig.HostRoleRelations) > 0 {
		deviceRole, err = common.MatchHostToRole(fs.Ctx, nbi, deviceName, fs.SourceConfig.HostRoleRelations)
		if err != nil {
			return fmt.Errorf("match host to role: %s", err)
		}
	}
	if deviceRole == nil {
		deviceRole, err = nbi.AddFirewallDeviceRole(fs.Ctx)
		if err != nil {
			return fmt.Errorf("add DeviceRole firewall: %s", err)
		}
	}
	deviceSite, err := common.MatchHostToSite(fs.Ctx, nbi, deviceName, fs.SourceConfig.HostSiteRelations)
	if err != nil {
		return fmt.Errorf("match host to site: %s", err)
	}
	devicePlatformName := fmt.Sprintf("FortiOS %s", fs.SystemInfo.Version)
	devicePlatform, err := nbi.AddPlatform(fs.Ctx, &objects.Platform{
		Name:         devicePlatformName,
		Slug:         utils.Slugify(devicePlatformName),
		Manufacturer: deviceManufacturer,
	})
	if err != nil {
		return fmt.Errorf("add platform: %s", err)
	}
	NBDevice, err := nbi.AddDevice(fs.Ctx, &objects.Device{
		NetboxObject: objects.NetboxObject{
			Tags: fs.SourceTags,
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

	fs.NBFirewall = NBDevice
	return nil
}

// syncInterfaces syncs all interfaces for firewall.
func (fs *FortigateSource) syncInterfaces(nbi *inventory.NetboxInventory) error {
	for _, iface := range fs.Ifaces {
		switch iface.Type {
		case "loopback":
		case "tunnel":
		case "vlan":
		case "aggregate":
		case "physical":
		case "hard-switch":
		}

		ifaceName := iface.Name
		if ifaceName == "" {
			fs.Logger.Warningf(fs.Ctx, "empty interface name - skipping")
			continue
		}

		if utils.FilterInterfaceName(ifaceName, fs.SourceConfig.InterfaceFilter) {
			fs.Logger.Debugf(fs.Ctx, "interface %s is filtered out with interfaceFilter %s", ifaceName, fs.SourceConfig.InterfaceFilter)
			continue
		}

		var interfaceStatus bool
		if iface.Status == "up" {
			interfaceStatus = true
		}
		interfaceMTU := iface.MTU
		interfaceMAC := iface.MAC

		var vdcs []*objects.VirtualDeviceContext
		if iface.Vdom != "" {
			vdom, err := nbi.AddVirtualDeviceContext(fs.Ctx, &objects.VirtualDeviceContext{
				NetboxObject: objects.NetboxObject{
					Tags: fs.SourceTags,
				},
				Name:   iface.Vdom,
				Device: fs.NBFirewall,
				Status: &objects.VDCStatusActive,
			})
			if err != nil {
				return fmt.Errorf("add VirtualDeviceContext: %s", err)
			}
			vdcs = append(vdcs, vdom)
		}
		NBIface, err := nbi.AddInterface(fs.Ctx, &objects.Interface{
			NetboxObject: objects.NetboxObject{
				Tags:        fs.SourceTags,
				Description: iface.Description,
			},
			Device: fs.NBFirewall,
			Type:   &objects.OtherInterfaceType,
			Name:   ifaceName,
			MTU:    interfaceMTU,
			MAC:    strings.ToUpper(interfaceMAC),
			Status: interfaceStatus,

			Vdcs: vdcs,
		})
		if err != nil {
			return fmt.Errorf("add interface: %s", err)
		}

		NBIPAddress, err := syncInterfaceIPs(fs, nbi, iface, NBIface)
		if err != nil {
			return fmt.Errorf("sync interface ips: %s", err)
		}

		if iface.Type == "vlan" {
			// Add Vlan for interface
			vlanID := iface.VlanID
			vlanName := fmt.Sprintf("Vlan%d", vlanID)
			vlanSite, err := common.MatchVlanToSite(fs.Ctx, nbi, vlanName, fs.SourceConfig.VlanSiteRelations)
			if err != nil {
				return fmt.Errorf("match vlan to site: %s", err)
			}
			vlanGroup, err := common.MatchVlanToGroup(fs.Ctx, nbi, vlanName, vlanSite, fs.SourceConfig.VlanGroupRelations, fs.SourceConfig.VlanGroupSiteRelations)
			if err != nil {
				return fmt.Errorf("match vlan to group: %s", err)
			}
			vlanTenant, err := common.MatchVlanToTenant(fs.Ctx, nbi, vlanName, fs.SourceConfig.VlanTenantRelations)
			if err != nil {
				return fmt.Errorf("match vlan to tenant: %s", err)
			}
			NBVlan, err := nbi.AddVlan(fs.Ctx, &objects.Vlan{
				NetboxObject: objects.NetboxObject{
					Tags: fs.SourceTags,
				},
				Status: &objects.VlanStatusActive,
				Name:   vlanName,
				Vid:    vlanID,
				Site:   vlanSite,
				Tenant: vlanTenant,
				Group:  vlanGroup,
			})
			if err != nil {
				return fmt.Errorf("add vlan: %s", err)
			}

			// Connect prefix with vlan
			if NBIPAddress != nil {
				prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(NBIPAddress.Address)
				if err != nil {
					fs.Logger.Warningf(fs.Ctx, "extract prefix from ip address: %s", err)
				} else if mask != constants.MaxIPv4MaskBits {
					_, err = nbi.AddPrefix(fs.Ctx, &objects.Prefix{
						Prefix: prefix,
						Tenant: NBVlan.Tenant,
						Vlan:   NBVlan,
					})
					if err != nil {
						return fmt.Errorf("add prefix: %s", err)
					}
				}
			}
		}
	}
	return nil
}

// syncInterfaceIPs is a helper function for syncInterfaces.
// it synces IPs for an interface.
func syncInterfaceIPs(fs *FortigateSource, nbi *inventory.NetboxInventory, iface InterfaceResponse, nbIface *objects.Interface) (*objects.IPAddress, error) {
	var NBIPAddress *objects.IPAddress
	ipAndMask := strings.Split(iface.IP, " ")
	if len(ipAndMask) == 2 && ipAndMask[0] != constants.WildcardIP {
		if utils.IsPermittedIPAddress(ipAndMask[0], fs.SourceConfig.PermittedSubnets, fs.SourceConfig.IgnoredSubnets) {
			maskBits, err := utils.MaskToBits(ipAndMask[1])
			if err != nil {
				return nil, fmt.Errorf("mask to bits: %s", err)
			}
			NBIPAddress, err = nbi.AddIPAddress(fs.Ctx, &objects.IPAddress{
				NetboxObject: objects.NetboxObject{
					Tags: fs.SourceTags,
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: false,
					},
				},
				Address:            fmt.Sprintf("%s/%d", ipAndMask[0], maskBits),
				AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
				AssignedObjectID:   nbIface.ID,
			})
			if err != nil {
				return nil, fmt.Errorf("add ip address: %s", err)
			}
		}
	}
	if len(iface.SecondaryIP) > 0 {
		for _, secondaryIP := range iface.SecondaryIP {
			ipAndMask := strings.Split(secondaryIP.IP, " ")
			if len(ipAndMask) == 2 && ipAndMask[0] != constants.WildcardIP {
				if utils.IsPermittedIPAddress(ipAndMask[0], fs.SourceConfig.PermittedSubnets, fs.SourceConfig.IgnoredSubnets) {
					maskBits, err := utils.MaskToBits(ipAndMask[1])
					if err != nil {
						return nil, fmt.Errorf("mask to bits: %s", err)
					}
					_, err = nbi.AddIPAddress(fs.Ctx, &objects.IPAddress{
						NetboxObject: objects.NetboxObject{
							Tags: fs.SourceTags,
							CustomFields: map[string]interface{}{
								constants.CustomFieldArpEntryName: false,
							},
						},
						Address:            fmt.Sprintf("%s/%d", ipAndMask[0], maskBits),
						AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
						AssignedObjectID:   nbIface.ID,
					})
					if err != nil {
						fs.Logger.Warningf(fs.Ctx, "add secondary ip address: %s", err)
					}
				}
			}
		}
	}

	if len(iface.VRRPIP) > 0 {
		for _, vrrp := range iface.VRRPIP {
			ipAndMask := []string{vrrp.VRIP, "255.255.255.255"}
			if len(ipAndMask) == 2 && ipAndMask[0] != constants.WildcardIP {
				if utils.IsPermittedIPAddress(ipAndMask[0], fs.SourceConfig.PermittedSubnets, fs.SourceConfig.IgnoredSubnets) {
					maskBits, err := utils.MaskToBits(ipAndMask[1])
					if err != nil {
						return nil, fmt.Errorf("mask to bits: %s", err)
					}
					_, err = nbi.AddIPAddress(fs.Ctx, &objects.IPAddress{
						NetboxObject: objects.NetboxObject{
							Tags: fs.SourceTags,
							CustomFields: map[string]interface{}{
								constants.CustomFieldArpEntryName: false,
							},
						},
						Address:            fmt.Sprintf("%s/%d", ipAndMask[0], maskBits),
						AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
						AssignedObjectID:   nbIface.ID,
						Role:               &objects.IPAddressRoleVRRP,
					})
					if err != nil {
						fs.Logger.Warningf(fs.Ctx, "add VRRP ip address: %s", err)
					}
				}
			}
		}
	}
	return NBIPAddress, nil
}
