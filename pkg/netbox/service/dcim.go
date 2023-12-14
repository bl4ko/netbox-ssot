package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/dcim"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
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

type DeviceRoleResponse struct {
	Count    int               `json:"count"`
	Next     int               `json:"next"`
	Previous int               `json:"previous"`
	Results  []dcim.DeviceRole `json:"results"`
}

// GET /api/dcim/device-roles/
func (api *NetboxAPI) GetAllDeviceRoles() ([]*dcim.DeviceRole, error) {
	api.Logger.Debug("Getting all device roles from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/device-roles/", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var deviceRoleResponse DeviceRoleResponse
	err = json.Unmarshal(response.Body, &deviceRoleResponse)
	if err != nil {
		return nil, err
	}

	deviceRoles := make([]*dcim.DeviceRole, len(deviceRoleResponse.Results))
	for i := range deviceRoleResponse.Results {
		deviceRoles[i] = &deviceRoleResponse.Results[i]
	}
	api.Logger.Debug("Device roles: ", deviceRoleResponse.Results)

	return deviceRoles, nil
}

// POST /api/dcim/device-roles/
func (api *NetboxAPI) CreateDeviceRole(deviceRole *dcim.DeviceRole) (*dcim.DeviceRole, error) {
	api.Logger.Debug("Creating device role in NetBox")

	requestBody, err := utils.NetboxJsonMarshal(deviceRole)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/dcim/device-roles/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var deviceRoleResponse dcim.DeviceRole
	err = json.Unmarshal(response.Body, &deviceRoleResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Device role: ", deviceRoleResponse)

	return &deviceRoleResponse, nil
}

// PATCH /api/dcim/device-roles/{id}/
func (api *NetboxAPI) PatchDeviceRole(diffMap map[string]interface{}, id int) (*dcim.DeviceRole, error) {
	api.Logger.Debug("Patching device role ", id, " with data: ", diffMap, " in NetBox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/dcim/device-roles/%d/", id), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var deviceRoleResponse dcim.DeviceRole
	err = json.Unmarshal(response.Body, &deviceRoleResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Device role: ", deviceRoleResponse)

	return &deviceRoleResponse, nil
}
