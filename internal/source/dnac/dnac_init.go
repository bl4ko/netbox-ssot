package dnac

import (
	"fmt"
	"net/http"

	dnac "github.com/cisco-en-programmability/dnacenter-go-sdk/v6/sdk"
)

// Collects all sites from DNAC API and stores them in the
// local source inventory.
func (ds *DnacSource) initSites(c *dnac.Client) error {
	offset := 0
	limit := 100
	allSites := make([]dnac.ResponseSitesGetSiteResponse, 0)
	for {
		sites, response, err := c.Sites.GetSite(&dnac.GetSiteQueryParams{Offset: offset, Limit: limit})
		if err != nil {
			return fmt.Errorf("init sites: %s", err)
		}
		if response.StatusCode() != http.StatusOK {
			return fmt.Errorf("init sites response code: %s", response.String())
		}
		allSites = append(allSites, *sites.Response...)
		if len(*sites.Response) < limit {
			break
		}
		offset += limit
	}
	ds.Sites = make(map[string]dnac.ResponseSitesGetSiteResponse, len(allSites))
	for _, site := range allSites {
		ds.Sites[site.ID] = site
	}
	return nil
}

// Collects all devices from DNAC API and stores them in the
// local source inventory.
func (ds *DnacSource) initDevices(c *dnac.Client) error {
	offset := 0
	limit := 100
	allDevices := make([]dnac.ResponseDevicesGetDeviceListResponse, 0)
	for {
		devices, response, err := c.Devices.GetDeviceList(&dnac.GetDeviceListQueryParams{Offset: offset, Limit: limit})
		if err != nil {
			return fmt.Errorf("init devices: %s", err)
		}
		if response.StatusCode() != http.StatusOK {
			return fmt.Errorf("init devices response code: %s", response.String())
		}
		allDevices = append(allDevices, *devices.Response...)
		if len(*devices.Response) < limit {
			break
		}
		offset += limit
	}

	ds.Devices = make(map[string]dnac.ResponseDevicesGetDeviceListResponse, len(allDevices))
	ds.Vlans = make(map[int]dnac.ResponseDevicesGetDeviceInterfaceVLANsResponse)
	for _, device := range allDevices {
		ds.Devices[device.ID] = device
		ds.DeviceID2isMissingPrimaryIP.Store(device.ID, true)
		ds.initVlansForDevice(c, device.ID)
	}
	return nil
}

// initVlansForDevice collects all VLANs for a device from DNAC API
// and stores them in the local source inventory.
func (ds *DnacSource) initVlansForDevice(c *dnac.Client, deviceID string) {
	vlans, _, _ := c.Devices.GetDeviceInterfaceVLANs(deviceID, nil)
	if vlans != nil {
		for _, vlan := range *vlans.Response {
			if vlan.VLANNumber != nil {
				ds.Vlans[*vlan.VLANNumber] = vlan
			}
		}
	}
}

// Collects all interfaces from DNAC API and stores them in the
// local source inventory.
func (ds *DnacSource) initInterfaces(c *dnac.Client) error {
	offset := 0
	limit := 100
	allInterfaces := make([]dnac.ResponseDevicesGetAllInterfacesResponse, 0)
	for {
		interfacesResponse, response, err := c.Devices.GetAllInterfaces(&dnac.GetAllInterfacesQueryParams{Offset: offset, Limit: limit})
		if err != nil {
			return fmt.Errorf("init interfaces: %s", err)
		}
		if response.StatusCode() != http.StatusOK {
			return fmt.Errorf("init interfaces response code: %s", response.String())
		}
		interfaces := *interfacesResponse.Response
		allInterfaces = append(allInterfaces, interfaces...)
		if len(interfaces) < limit {
			break
		}
		offset += limit
	}
	ds.Interfaces = make(map[string]dnac.ResponseDevicesGetAllInterfacesResponse, len(allInterfaces))
	ds.DeviceID2InterfaceIDs = make(map[string][]string)
	for _, intf := range allInterfaces {
		ds.Interfaces[intf.ID] = intf
		if ds.DeviceID2InterfaceIDs[intf.DeviceID] == nil {
			ds.DeviceID2InterfaceIDs[intf.DeviceID] = make([]string, 0)
		}
		ds.DeviceID2InterfaceIDs[intf.DeviceID] = append(ds.DeviceID2InterfaceIDs[intf.DeviceID], intf.ID)
	}
	return nil
}

// For each site id finds the corresponding device ids.
// This is necessary to find relations between devices and sites.
//
// This function has to run after InitSites.
func (ds *DnacSource) initMemberships(c *dnac.Client) error {
	offset := 0
	limit := 100
	ds.Site2Devices = make(map[string]map[string]bool)
	ds.Device2Site = make(map[string]string)
	for _, site := range ds.Sites {
		for {
			membershipResp, _, _ := c.Sites.GetMembership(site.ID, &dnac.GetMembershipQueryParams{Offset: float64(offset), Limit: float64(limit)})
			if len(*membershipResp.Device) > 0 {
				deviceResponses := *membershipResp.Device
				for _, deviceResponse := range deviceResponses {
					siteID := deviceResponse.SiteID
					devices := *deviceResponse.Response
					for _, device := range devices {
						if deviceMap, ok := device.(map[string]interface{}); ok {
							if deviceID, ok := deviceMap["instanceUuid"].(string); ok {
								if ds.Site2Devices[siteID] == nil {
									ds.Site2Devices[siteID] = make(map[string]bool)
								}
								ds.Site2Devices[siteID][deviceID] = true
								ds.Device2Site[deviceID] = siteID
							}
						}
					}
				}
			}
			if len(*membershipResp.Device) < limit {
				break
			}
			offset += limit
		}
	}
	return nil
}

// initWirelessLANs collects all wireless profiles, dynamic interfaces
// and enterprise SSIDs from DNAC API and stores them in the local source inventory.
// All this data is necessary to create WirelessLANs and WirelessLANGroups in netbox.
func (ds *DnacSource) initWirelessLANs(c *dnac.Client) error {
	// Get all WirelessProfiles
	wirelessProfiles, response, err := c.Wireless.GetWirelessProfile(nil)
	if err != nil {
		return fmt.Errorf("init wireless profiles: %s", err)
	}
	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("init wireless profiles response code: %s", response.String())
	}

	// Get wireless lan vlan relation.
	wirelessDynamicInterfaces, response, err := c.Wireless.GetDynamicInterface(nil)
	if err != nil {
		return fmt.Errorf("init wireless dynamic interfaces: %s", err)
	}
	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("init wireless dynamic interfaces: %s", response.String())
	}

	// Get wireless lan data.
	enterpriseSsids, response, err := c.Wireless.GetEnterpriseSSID(nil)
	if err != nil {
		return fmt.Errorf("init enterprise ssids: %s", err)
	}
	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("init enterprise ssids response code: %s", response.String())
	}

	// Create a map of IntefaceName -> VLAN
	ds.WirelessLANInterfaceName2VlanID = make(map[string]int)
	for _, dynamicInterface := range *wirelessDynamicInterfaces {
		ds.WirelessLANInterfaceName2VlanID[dynamicInterface.InterfaceName] = int(*dynamicInterface.VLANID)
	}

	// Create a map of SSID -> WirelessLANGroup
	ds.SSID2WlanGroupName = make(map[string]string)
	ds.SSID2WirelessProfileDetails = make(map[string]dnac.ResponseItemWirelessGetWirelessProfileProfileDetailsSSIDDetails)
	for _, wirelessProfile := range *wirelessProfiles {
		for _, ssid := range *wirelessProfile.ProfileDetails.SSIDDetails {
			ds.SSID2WirelessProfileDetails[ssid.Name] = ssid
			ds.SSID2WlanGroupName[ssid.Name] = wirelessProfile.ProfileDetails.Name
		}
	}

	// sync enterprise SSID
	ds.SSID2SecurityDetails = make(map[string]dnac.ResponseItemWirelessGetEnterpriseSSIDSSIDDetails)
	for _, enterpriseSSID := range *enterpriseSsids {
		for _, SSIDDetails := range *enterpriseSSID.SSIDDetails {
			ds.SSID2SecurityDetails[SSIDDetails.Name] = SSIDDetails
		}
	}

	return nil
}
