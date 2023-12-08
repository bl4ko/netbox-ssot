package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/dcim"
)

type SiteResponse struct {
	Count    int         `json:"count"`
	Next     int         `json:"next"`
	Previous int         `json:"previous"`
	Results  []dcim.Site `json:"results"`
}

// GET /api/dcim/sites/
func (api *NetboxAPI) GetAllSites() ([]*dcim.Site, error) {
	api.Logger.Debug("Getting all sites from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/sites/", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var siteResponse SiteResponse
	err = json.Unmarshal(response.Body, &siteResponse)
	if err != nil {
		return nil, err
	}

	sites := make([]*dcim.Site, len(siteResponse.Results))
	for i := range siteResponse.Results {
		sites[i] = &siteResponse.Results[i]
	}
	api.Logger.Debug("Sites: ", siteResponse.Results)

	return sites, nil
}

// top-level JSON object returned by Netbox API
type DeviceResponse struct {
	Count    int           `json:"count"`
	Next     *string       `json:"next"`
	Previous *string       `json:"previous"`
	Results  []dcim.Device `json:"results"`
}

// GET /api/dcim/devices/
func (api *NetboxAPI) GetAllDevices() ([]*dcim.Device, error) {
	api.Logger.Debug("Getting all devices from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/devices/", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var deviceResponse DeviceResponse
	err = json.Unmarshal(response.Body, &deviceResponse)
	if err != nil {
		return nil, err
	}

	devices := make([]*dcim.Device, len(deviceResponse.Results))
	for i := range deviceResponse.Results {
		devices[i] = &deviceResponse.Results[i]
	}
	api.Logger.Debug("Devices: ", deviceResponse.Results)

	return devices, nil
}
