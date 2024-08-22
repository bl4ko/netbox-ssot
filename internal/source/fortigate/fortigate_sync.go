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
	deviceSerialNumber := fs.SystemInfo.Serial
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

	deviceTenant, err := common.MatchHostToTenant(fs.Ctx, nbi, deviceName, fs.HostTenantRelations)
	if err != nil {
		return fmt.Errorf("match host to tenant: %s", err)
	}

	deviceRole, err := nbi.AddDeviceRole(fs.Ctx, &objects.DeviceRole{
		Name:   constants.DeviceRoleFirewall,
		Slug:   utils.Slugify(constants.DeviceRoleFirewall),
		Color:  constants.DeviceRoleFirewallColor,
		VMRole: false,
	})
	if err != nil {
		return fmt.Errorf("add DeviceRole: %s", err)
	}
	deviceSite, err := common.MatchHostToSite(fs.Ctx, nbi, deviceName, fs.HostSiteRelations)
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
			MAC:    interfaceMAC,
			Status: interfaceStatus,

			Vdcs: vdcs,
		})
		if err != nil {
			return fmt.Errorf("add interface: %s", err)
		}

		var NBIPAddress *objects.IPAddress
		ipAndMask := strings.Split(iface.IP, " ")
		if len(ipAndMask) == 2 && ipAndMask[0] != "0.0.0.0" {
			maskBits, err := utils.MaskToBits(ipAndMask[1])
			if err != nil {
				return fmt.Errorf("mask to bits: %s", err)
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
				AssignedObjectID:   NBIface.ID,
			})
			if err != nil {
				return fmt.Errorf("add ip address: %s", err)
			}
		}

		if iface.Type == "vlan" {
			// Add Vlan for interface
			vlanID := iface.VlanID
			vlanName := fmt.Sprintf("Vlan%d", vlanID)
			vlanGroup, err := common.MatchVlanToGroup(fs.Ctx, nbi, vlanName, fs.VlanGroupRelations)
			if err != nil {
				return fmt.Errorf("match vlan to group: %s", err)
			}
			vlanTenant, err := common.MatchVlanToTenant(fs.Ctx, nbi, vlanName, fs.VlanTenantRelations)
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
