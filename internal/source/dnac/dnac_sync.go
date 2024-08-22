package dnac

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	dnac "github.com/cisco-en-programmability/dnacenter-go-sdk/v5/sdk"
)

// Syncs dnac sites to netbox inventory.
func (ds *DnacSource) syncSites(nbi *inventory.NetboxInventory) error {
	for _, site := range ds.Sites {
		dnacSite := &objects.Site{
			NetboxObject: objects.NetboxObject{
				Tags: ds.Config.SourceTags,
				CustomFields: map[string]interface{}{
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
		ds.SiteID2nbSite.Store(site.ID, nbSite)
	}
	return nil
}

// Syncs dnac vlans to netbox inventory.
func (ds *DnacSource) syncVlans(nbi *inventory.NetboxInventory) error {
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
				CustomFields: map[string]interface{}{
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
					CustomFields: map[string]interface{}{
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

		ds.VID2nbVlan.Store(vid, newVlan)
	}
	return nil
}
func (ds *DnacSource) syncDevices(nbi *inventory.NetboxInventory) error {
	const maxGoroutines = 50
	guard := make(chan struct{}, maxGoroutines)
	errChan := make(chan error, len(ds.Devices))
	var wg sync.WaitGroup

	for deviceID, device := range ds.Devices {
		guard <- struct{}{} // Block if maxGoroutines are running
		wg.Add(1)

		go func(deviceID string, device dnac.ResponseDevicesGetDeviceListResponse) {
			defer wg.Done()
			defer func() { <-guard }() // Release one spot in the semaphore

			err := ds.syncDevice(nbi, deviceID, device)
			if err != nil {
				errChan <- err
			}
		}(deviceID, device)
	}

	wg.Wait()
	close(errChan)
	close(guard)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *DnacSource) syncDevice(nbi *inventory.NetboxInventory, deviceID string, device dnac.ResponseDevicesGetDeviceListResponse) error {
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
	} else {
		platformName = strings.Trim(fmt.Sprintf("%s %s", device.SoftwareType, device.SoftwareVersion), " ")
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
	if site, ok := ds.SiteID2nbSite.Load(ds.Device2Site[device.ID]); ok {
		if deviceSite, ok = site.(*objects.Site); !ok {
			ds.Logger.Errorf(ds.Ctx, "Type assertion to *objects.Site failed for device %s, this should not happen. This device will be skipped", device.ID)
			return nil
		}
	} else {
		ds.Logger.Errorf(ds.Ctx, "DeviceSite is not existing for device %s, this should not happen. This device will be skipped", device.ID)
		return nil
	}

	if device.Type == "" {
		ds.Logger.Errorf(ds.Ctx, "Device type for device %s is empty, this should not happen. This device will be skipped", device.ID)
		return nil
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
			CustomFields: map[string]interface{}{
				constants.CustomFieldSourceName:     ds.SourceConfig.Name,
				constants.CustomFieldSourceIDName:   deviceID,
				constants.CustomFieldDeviceUUIDName: device.InstanceUUID,
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

	ds.DeviceID2nbDevice.Store(device.ID, nbDevice)
	return nil
}

func (ds *DnacSource) syncDeviceInterfaces(nbi *inventory.NetboxInventory) error {
	const maxGoroutines = 50
	guard := make(chan struct{}, maxGoroutines)
	errChan := make(chan error, len(ds.Interfaces))
	var wg sync.WaitGroup

	for ifaceID, iface := range ds.Interfaces {
		guard <- struct{}{}
		wg.Add(1)

		go func(ifaceID string, iface dnac.ResponseDevicesGetAllInterfacesResponse) {
			defer wg.Done()
			defer func() { <-guard }()

			err := ds.syncDeviceInterface(nbi, ifaceID, iface)
			if err != nil {
				errChan <- err
			}
		}(ifaceID, iface)
	}

	wg.Wait()
	close(errChan)
	close(guard)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
func (ds *DnacSource) syncDeviceInterface(nbi *inventory.NetboxInventory, ifaceID string, iface dnac.ResponseDevicesGetAllInterfacesResponse) error {
	ifaceDescription := iface.Description

	ifaceDevice, err := ds.getDevice(iface.DeviceID)
	if err != nil {
		ds.Logger.Errorf(ds.Ctx, "%s This interface will be skipped", err)
		return nil
	}

	ifaceDuplex := ds.getInterfaceDuplex(iface.Duplex)
	ifaceStatus, err := ds.getInterfaceStatus(iface.Status)
	if err != nil {
		ds.Logger.Errorf(ds.Ctx, "%s", err)
		return nil
	}

	var ifaceSpeed objects.InterfaceSpeed
	var ifaceType *objects.InterfaceType
	speed, err := strconv.Atoi(iface.Speed)
	if err != nil {
		ds.Logger.Errorf(ds.Ctx, "wrong speed for iface %s", iface.Speed)
	} else {
		ifaceSpeed = objects.InterfaceSpeed(speed)
		typeI, err := ds.getInterfaceType(iface.InterfaceType, speed)
		if err != nil {
			ds.Logger.Errorf(ds.Ctx, "%s. Skipping this device...", err)
			return nil
		}
		ifaceType = typeI
	}

	ifaceName := iface.PortName
	if err := ds.validateInterfaceName(ifaceName, ifaceID); err != nil {
		ds.Logger.Errorf(ds.Ctx, "%s", err)
		return nil
	}

	ifaceMode, ifaceAccessVlan, err := ds.getVlanModeAndAccessVlan(iface.PortMode, iface.VLANID)
	if err != nil {
		ds.Logger.Errorf(ds.Ctx, "%s", err)
		return nil
	}

	nbIface, err := nbi.AddInterface(ds.Ctx, &objects.Interface{
		NetboxObject: objects.NetboxObject{
			Description: ifaceDescription,
			Tags:        ds.Config.SourceTags,
			CustomFields: map[string]interface{}{
				constants.CustomFieldSourceName: ds.SourceConfig.Name,
			},
		},
		Name:         ifaceName,
		MAC:          strings.ToUpper(iface.MacAddress),
		Speed:        ifaceSpeed,
		Status:       ifaceStatus,
		Duplex:       ifaceDuplex,
		Device:       ifaceDevice,
		Type:         ifaceType,
		Mode:         ifaceMode,
		UntaggedVlan: ifaceAccessVlan,
		TaggedVlans:  nil, // placeholder for tagged VLANs
	})
	if err != nil {
		return fmt.Errorf("add device interface: %s", err)
	}

	err = ds.addIPAddressToInterface(nbi, nbIface, iface)
	if err != nil {
		ds.Logger.Errorf(ds.Ctx, "adding IP address: %s", err)
	}

	ds.InterfaceID2nbInterface.Store(ifaceID, nbIface)
	return nil
}

func (ds *DnacSource) getDevice(deviceID string) (*objects.Device, error) {
	if device, ok := ds.DeviceID2nbDevice.Load(deviceID); ok {
		if ifaceDevice, ok := device.(*objects.Device); ok {
			return ifaceDevice, nil
		}
		return nil, fmt.Errorf("type assertion to *objects.Device failed for device %s", deviceID)
	}
	return nil, fmt.Errorf("device %s not found", deviceID)
}

func (ds *DnacSource) getInterfaceDuplex(duplex string) *objects.InterfaceDuplex {
	switch duplex {
	case "FullDuplex":
		return &objects.DuplexFull
	case "AutoNegotiate":
		return &objects.DuplexAuto
	case "HalfDuplex":
		return &objects.DuplexHalf
	default:
		ds.Logger.Warningf(ds.Ctx, "Unknown duplex value: %s", duplex)
		return nil
	}
}

func (ds *DnacSource) getInterfaceStatus(status string) (bool, error) {
	switch status {
	case "down":
		return false, nil
	case "up":
		return true, nil
	default:
		return false, fmt.Errorf("wrong interface status: %s", status)
	}
}

func (ds *DnacSource) getInterfaceType(interfaceType string, speed int) (*objects.InterfaceType, error) {
	switch interfaceType {
	case "Physical":
		ifaceType := objects.IfaceSpeed2IfaceType[objects.InterfaceSpeed(speed)]
		if ifaceType == nil {
			return &objects.OtherInterfaceType, nil
		}
		return ifaceType, nil
	case "Virtual":
		return &objects.VirtualInterfaceType, nil
	default:
		return nil, fmt.Errorf("unknown interface type: %s", interfaceType)
	}
}

func (ds *DnacSource) validateInterfaceName(ifaceName, ifaceID string) error {
	if ifaceName == "" {
		return fmt.Errorf("unknown interface name for iface: %s", ifaceID)
	}
	if utils.FilterInterfaceName(ifaceName, ds.SourceConfig.InterfaceFilter) {
		return fmt.Errorf("interface %s is filtered out with interfaceFilter %s", ifaceName, ds.SourceConfig.InterfaceFilter)
	}
	return nil
}

func (ds *DnacSource) getVlanModeAndAccessVlan(portMode, vlanID string) (*objects.InterfaceMode, *objects.Vlan, error) {
	vid, err := strconv.Atoi(vlanID)
	if err != nil {
		return nil, nil, fmt.Errorf("can't parse vid for iface %s", vlanID)
	}

	switch portMode {
	case "access":
		ifaceMode := &objects.InterfaceModeAccess
		if accessVlan, ok := ds.VID2nbVlan.Load(vid); ok {
			if ifaceAccessVlan, ok := accessVlan.(*objects.Vlan); ok {
				return ifaceMode, ifaceAccessVlan, nil
			}
			return nil, nil, fmt.Errorf("type assertion to *objects.Vlan failed for vlan %d", vid)
		}
		return ifaceMode, nil, nil
	case "trunk":
		return &objects.InterfaceModeTagged, nil, nil
	case "dynamic_auto", "routed":
		ds.Logger.Debugf(ds.Ctx, "vlan mode '%s' is not implemented yet", portMode)
		return nil, nil, nil
	default:
		return nil, nil, fmt.Errorf("unknown interface mode: '%s'", portMode)
	}
}

func (ds *DnacSource) addIPAddressToInterface(nbi *inventory.NetboxInventory, iface *objects.Interface, ifaceDetails dnac.ResponseDevicesGetAllInterfacesResponse) error {
	if ifaceDetails.IPv4Address == "" || utils.SubnetsContainIPAddress(ifaceDetails.IPv4Address, ds.SourceConfig.IgnoredSubnets) {
		return nil
	}

	defaultMask := 32
	if ifaceDetails.IPv4Mask != "" {
		maskBits, err := utils.MaskToBits(ifaceDetails.IPv4Mask)
		if err != nil {
			return fmt.Errorf("wrong mask: %s", err)
		}
		defaultMask = maskBits
	}

	nbIPAddress, err := nbi.AddIPAddress(ds.Ctx, &objects.IPAddress{
		NetboxObject: objects.NetboxObject{
			Tags: ds.Config.SourceTags,
			CustomFields: map[string]interface{}{
				constants.CustomFieldSourceName:   ds.SourceConfig.Name,
				constants.CustomFieldArpEntryName: false,
			},
		},
		Address:            fmt.Sprintf("%s/%d", ifaceDetails.IPv4Address, defaultMask),
		Status:             &objects.IPAddressStatusActive,
		DNSName:            utils.ReverseLookup(ifaceDetails.IPv4Address),
		AssignedObjectType: objects.AssignedObjectTypeDeviceInterface,
		AssignedObjectID:   iface.ID,
		Tenant:             iface.Device.Tenant,
	})
	if err != nil {
		return fmt.Errorf("adding IP address: %s", err)
	}

	// Optionally, add the prefix to NetBox
	prefix, mask, err := utils.GetPrefixAndMaskFromIPAddress(nbIPAddress.Address)
	if err != nil {
		ds.Logger.Warningf(ds.Ctx, "failed extracting prefix from IPAddress: %s", err)
	} else if mask != constants.MaxIPv4MaskBits {
		_, err = nbi.AddPrefix(ds.Ctx, &objects.Prefix{
			Prefix: prefix,
			Tenant: iface.Device.Tenant,
		})
		if err != nil {
			ds.Logger.Errorf(ds.Ctx, "adding prefix: %s", err)
		}
	}

	// Set the interface as the primary IPv4 if it matches the device's management IP
	deviceManagementIP := ds.Devices[ifaceDetails.DeviceID].ManagementIPAddress
	if deviceManagementIP == ifaceDetails.IPv4Address {
		deviceCopy := *iface.Device
		deviceCopy.PrimaryIPv4 = nbIPAddress
		_, err := nbi.AddDevice(ds.Ctx, &deviceCopy)
		if err != nil {
			ds.Logger.Errorf(ds.Ctx, "adding primary IPv4 address: %s", err)
		}
	}

	return nil
}
