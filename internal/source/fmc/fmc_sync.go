package fmc

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

func (fmcs *FMCSource) syncDevices(nbi *inventory.NetboxInventory) error {
	for deviceUUID, device := range fmcs.Devices {
		deviceName := device.Name
		if deviceName == "" {
			fmcs.Logger.Warningf(fmcs.Ctx, "device with empty name. Skipping...")
			continue
		}
		var deviceSerialNumber string
		if !fmcs.SourceConfig.IgnoreSerialNumbers {
			deviceSerialNumber = device.Metadata.SerialNumber
		}
		deviceModel := device.Model
		if deviceModel == "" {
			fmcs.Logger.Warning(fmcs.Ctx, "model field for device is emptpy. Using fallback model.")
			deviceModel = constants.DefaultModel
		}
		deviceManufacturer, err := nbi.AddManufacturer(fmcs.Ctx, &objects.Manufacturer{
			Name: "Cisco",
			Slug: utils.Slugify("Cisco"),
		})
		if err != nil {
			return fmt.Errorf("add manufacturer: %s", err)
		}
		deviceType, err := nbi.AddDeviceType(fmcs.Ctx, &objects.DeviceType{
			Manufacturer: deviceManufacturer,
			Model:        deviceModel,
			Slug:         utils.Slugify(deviceManufacturer.Name + deviceModel),
		})
		if err != nil {
			return fmt.Errorf("add device type: %s", err)
		}
		deviceTenant, err := common.MatchHostToTenant(fmcs.Ctx, nbi, deviceName, fmcs.SourceConfig.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("match host to tenant %s", err)
		}

		// Match host to a role. First test if user provided relations, if not
		// use default firewall role.
		var deviceRole *objects.DeviceRole
		if len(fmcs.SourceConfig.HostRoleRelations) > 0 {
			deviceRole, err = common.MatchHostToRole(fmcs.Ctx, nbi, deviceName, fmcs.SourceConfig.HostRoleRelations)
			if err != nil {
				return fmt.Errorf("match host to role: %s", err)
			}
		}
		if deviceRole == nil {
			deviceRole, err = nbi.AddFirewallDeviceRole(fmcs.Ctx)
			if err != nil {
				return fmt.Errorf("add DeviceRole firewall: %s", err)
			}
		}

		deviceSite, err := common.MatchHostToSite(fmcs.Ctx, nbi, deviceName, fmcs.SourceConfig.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("match host to site: %s", err)
		}
		devicePlatformName := fmt.Sprintf("FXOS %s", device.SWVersion)
		devicePlatform, err := nbi.AddPlatform(fmcs.Ctx, &objects.Platform{
			Name:         devicePlatformName,
			Slug:         utils.Slugify(devicePlatformName),
			Manufacturer: deviceManufacturer,
		})
		if err != nil {
			return fmt.Errorf("add platform: %s", err)
		}
		NBDevice, err := nbi.AddDevice(fmcs.Ctx, &objects.Device{
			NetboxObject: objects.NetboxObject{
				Description: device.Description,
				Tags:        fmcs.SourceTags,
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceIDName:     deviceUUID,
					constants.CustomFieldDeviceUUIDName:   deviceUUID,
					constants.CustomFieldHostCPUCoresName: device.Metadata.InventoryData.CPUCores,
					constants.CustomFieldHostMemoryName:   fmt.Sprintf("%sMB", device.Metadata.InventoryData.MemoryInMB),
				},
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
		err = fmcs.syncPhysicalInterfaces(nbi, NBDevice, deviceUUID)
		if err != nil {
			return fmt.Errorf("sync physical interfaces: %s", err)
		}
		err = fmcs.syncVlanInterfaces(nbi, NBDevice, deviceUUID)
		if err != nil {
			return fmt.Errorf("sync vlan interfaces: %s", err)
		}
	}
	return nil
}

func (fmcs *FMCSource) syncVlanInterfaces(nbi *inventory.NetboxInventory, nbDevice *objects.Device, deviceUUID string) error {
	if vlanIfaces, ok := fmcs.DeviceVlanIfaces[deviceUUID]; ok {
		for _, vlanIface := range vlanIfaces {
			// Add vlan
			ifaceTaggedVlans := []*objects.Vlan{}
			if vlanIface.VID != 0 {
				// Match vlan to site
				vlanSite, err := common.MatchVlanToSite(fmcs.Ctx, nbi, vlanIface.Name, fmcs.SourceConfig.VlanSiteRelations)
				if err != nil {
					return fmt.Errorf("match vlan to site: %s", err)
				}
				// Match vlan to group
				vlanGroup, err := common.MatchVlanToGroup(fmcs.Ctx, nbi, vlanIface.Name, vlanSite, fmcs.SourceConfig.VlanGroupRelations, fmcs.SourceConfig.VlanGroupSiteRelations)
				if err != nil {
					return fmt.Errorf("match vlan to group: %s", err)
				}
				vlanTenent, err := common.MatchVlanToTenant(fmcs.Ctx, nbi, vlanIface.Name, fmcs.SourceConfig.VlanTenantRelations)
				if err != nil {
					return fmt.Errorf("match vlan to tenant: %s", err)
				}
				vlan, err := nbi.AddVlan(fmcs.Ctx, &objects.Vlan{
					NetboxObject: objects.NetboxObject{
						Tags:        fmcs.SourceTags,
						Description: vlanIface.Description,
					},
					Status: &objects.VlanStatusActive,
					Name:   vlanIface.Name,
					Site:   vlanSite,
					Vid:    vlanIface.VID,
					Tenant: vlanTenent,
					Group:  vlanGroup,
				})
				if err != nil {
					return fmt.Errorf("add vlan: %s", err)
				}
				ifaceTaggedVlans = append(ifaceTaggedVlans, vlan)
			}

			NBIface, err := nbi.AddInterface(fmcs.Ctx, &objects.Interface{
				NetboxObject: objects.NetboxObject{
					Description: vlanIface.Description,
					Tags:        fmcs.SourceTags,
					CustomFields: map[string]interface{}{
						constants.CustomFieldSourceIDName: vlanIface.ID,
					},
				},
				Name:        vlanIface.Name,
				Device:      nbDevice,
				Status:      vlanIface.Enabled,
				MTU:         vlanIface.MTU,
				TaggedVlans: ifaceTaggedVlans,
				Type:        &objects.VirtualInterfaceType,
			})
			if err != nil {
				return fmt.Errorf("add vlan interface: %s", err)
			}

			if vlanIface.IPv4 != nil && vlanIface.IPv4.Static != nil {
				if utils.IsPermittedIPAddress(vlanIface.IPv4.Static.Address, fmcs.SourceConfig.PermittedSubnets, fmcs.SourceConfig.IgnoredSubnets) {
					address := fmt.Sprintf("%s/%s", vlanIface.IPv4.Static.Address, vlanIface.IPv4.Static.Netmask)
					dnsName := utils.ReverseLookup(vlanIface.IPv4.Static.Address)
					_, err := nbi.AddIPAddress(fmcs.Ctx, &objects.IPAddress{
						NetboxObject: objects.NetboxObject{
							Tags: fmcs.SourceTags,
							CustomFields: map[string]interface{}{
								constants.CustomFieldArpEntryName: false,
							},
						},
						Address:            address,
						DNSName:            dnsName,
						AssignedObjectID:   NBIface.ID,
						AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
					})
					if err != nil {
						return fmt.Errorf("add ip address")
					}
					// Also add prefix
					prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(address)
					if err != nil {
						fmcs.Logger.Debugf(fmcs.Ctx, "extract prefix from address: %s", err)
					} else if mask != constants.MaxIPv4MaskBits {
						var prefixTenant *objects.Tenant
						var prefixVlan *objects.Vlan
						if len(ifaceTaggedVlans) > 0 {
							prefixVlan = ifaceTaggedVlans[0]
							prefixTenant = prefixVlan.Tenant
						}
						_, err = nbi.AddPrefix(fmcs.Ctx, &objects.Prefix{
							Prefix: prefix,
							Tenant: prefixTenant,
							Vlan:   prefixVlan,
						})
						if err != nil {
							return fmt.Errorf("add prefix: %s", err)
						}
					}
				}
			}
		}
	}
	return nil
}

func (fmcs *FMCSource) syncPhysicalInterfaces(nbi *inventory.NetboxInventory, nbDevice *objects.Device, deviceUUID string) error {
	if physicalIfaces, ok := fmcs.DevicePhysicalIfaces[deviceUUID]; ok {
		for _, pIface := range physicalIfaces {
			NBIface, err := nbi.AddInterface(fmcs.Ctx, &objects.Interface{
				NetboxObject: objects.NetboxObject{
					Description: pIface.Description,
					Tags:        fmcs.SourceTags,
					CustomFields: map[string]interface{}{
						constants.CustomFieldSourceIDName: pIface.ID,
					},
				},
				Name:   pIface.Name,
				Device: nbDevice,
				Status: pIface.Enabled,
				MTU:    pIface.MTU,
				Type:   &objects.OtherInterfaceType, // TODO
			})
			if err != nil {
				return fmt.Errorf("add vlan interface: %s", err)
			}

			if pIface.IPv4 != nil && pIface.IPv4.Static != nil {
				if utils.IsPermittedIPAddress(pIface.IPv4.Static.Address, fmcs.SourceConfig.PermittedSubnets, fmcs.SourceConfig.IgnoredSubnets) {
					address := fmt.Sprintf("%s/%s", pIface.IPv4.Static.Address, pIface.IPv4.Static.Netmask)
					dnsName := utils.ReverseLookup(pIface.IPv4.Static.Address)
					_, err := nbi.AddIPAddress(fmcs.Ctx, &objects.IPAddress{
						NetboxObject: objects.NetboxObject{
							Tags: fmcs.SourceTags,
							CustomFields: map[string]interface{}{
								constants.CustomFieldArpEntryName: false,
							},
						},
						Address:            address,
						DNSName:            dnsName,
						AssignedObjectID:   NBIface.ID,
						AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
					})
					if err != nil {
						return fmt.Errorf("add ip address")
					}
				}
			}
		}
	}
	return nil
}
