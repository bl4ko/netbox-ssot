package dnac

import (
	"fmt"
	"net/http"

	dnac "github.com/cisco-en-programmability/dnacenter-go-sdk/v5/sdk"
)

func (ds *DnacSource) InitSites(c *dnac.Client) error {
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

func (ds *DnacSource) InitDevices(c *dnac.Client) error {
	offset := 0.
	limit := 100.
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
		if len(*devices.Response) < int(limit) {
			break
		}
		offset += limit
	}

	ds.Devices = make(map[string]dnac.ResponseDevicesGetDeviceListResponse, len(allDevices))
	ds.Vlans = make(map[int]dnac.ResponseDevicesGetDeviceInterfaceVLANsResponse)
	for _, device := range allDevices {
		ds.Devices[device.ID] = device
		ds.initVlansForDevice(c, device.ID)
	}
	return nil
}

// Function that gets all vlans for device id.
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

func (ds *DnacSource) InitInterfaces(c *dnac.Client) error {
	offset := 0
	limit := 100
	allInterfaces := make([]dnac.ResponseDevicesGetAllInterfacesResponse, 0)
	for {
		interfacesResponse, response, err := c.Devices.GetAllInterfaces(&dnac.GetAllInterfacesQueryParams{Offset: float64(offset), Limit: float64(limit)})
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
func (ds *DnacSource) InitMemberships(c *dnac.Client) error {
	offset := 0
	limit := 100
	ds.Site2Devices = make(map[string]map[string]bool)
	ds.Device2Site = make(map[string]string)
	for _, site := range ds.Sites {
		for {
			membershipResp, _, _ := c.Sites.GetMembership(site.ID, &dnac.GetMembershipQueryParams{Offset: offset, Limit: limit})
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
