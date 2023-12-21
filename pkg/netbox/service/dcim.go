package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/common"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/dcim"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
)

type SiteResponse struct {
	Count    int           `json:"count"`
	Next     *string       `json:"next"`
	Previous *string       `json:"previous"`
	Results  []common.Site `json:"results"`
}

// GET /api/dcim/sites/?limit=0
func (api *NetboxAPI) GetAllSites() ([]*common.Site, error) {
	api.Logger.Debug("Getting all sites from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/sites/?limit=0", nil)
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

	sites := make([]*common.Site, len(siteResponse.Results))
	for i := range siteResponse.Results {
		sites[i] = &siteResponse.Results[i]
	}
	api.Logger.Debug("Successfully received all sites from netbox: ", siteResponse.Results)

	return sites, nil
}

// top-level JSON object returned by Netbox API
type DeviceResponse struct {
	Count    int           `json:"count"`
	Next     *string       `json:"next"`
	Previous *string       `json:"previous"`
	Results  []dcim.Device `json:"results"`
}

// GET /api/dcim/devices/?limit=0
func (api *NetboxAPI) GetAllDevices() ([]*dcim.Device, error) {
	api.Logger.Debug("Getting all devices from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/devices/?limit=0", nil)
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
	api.Logger.Debug("Successfully received all devices from netbox: ", deviceResponse.Results)

	return devices, nil
}

// PATCH /api/dcim/devices/{id}/
func (api *NetboxAPI) PatchDevice(diffMap map[string]interface{}, deviceId int) (*dcim.Device, error) {
	api.Logger.Debug("Patching device ", deviceId, " with data: ", diffMap, " in NetBox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/dcim/devices/%d/", deviceId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var deviceResponse dcim.Device
	err = json.Unmarshal(response.Body, &deviceResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched device: ", deviceResponse)
	return &deviceResponse, nil
}

// POST /api/dcim/devices/
func (api *NetboxAPI) CreateDevice(device *dcim.Device) (*dcim.Device, error) {
	api.Logger.Debug("Creating device in NetBox with data: ", device)

	requestBody, err := utils.NetboxJsonMarshal(device)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/dcim/devices/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var deviceResponse dcim.Device
	err = json.Unmarshal(response.Body, &deviceResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created device: ", deviceResponse)

	return &deviceResponse, nil
}

type DeviceRoleResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []dcim.DeviceRole `json:"results"`
}

// GET /api/dcim/device-roles/?limit=0
func (api *NetboxAPI) GetAllDeviceRoles() ([]*dcim.DeviceRole, error) {
	api.Logger.Debug("Getting all device roles from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/device-roles/?limit=0", nil)
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
	api.Logger.Debug("Successfully received device roles: ", deviceRoleResponse.Results)

	return deviceRoles, nil
}

// POST /api/dcim/device-roles/
func (api *NetboxAPI) CreateDeviceRole(deviceRole *dcim.DeviceRole) (*dcim.DeviceRole, error) {
	api.Logger.Debug("Creating device role with data", deviceRole, " in NetBox")

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

	api.Logger.Debug("Successfully created device role: ", deviceRoleResponse)

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

	api.Logger.Debug("Successfully patched device role: ", deviceRoleResponse)

	return &deviceRoleResponse, nil
}

type ManufacturerResponse struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"`
	Previous *string               `json:"previous"`
	Results  []common.Manufacturer `json:"results"`
}

// GET /api/dcim/manufacturers/?limit=0
func (api *NetboxAPI) GetAllManufacturers() ([]*common.Manufacturer, error) {
	api.Logger.Debug("Getting all manufacturers from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/manufacturers/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var manufacturerResponse ManufacturerResponse
	err = json.Unmarshal(response.Body, &manufacturerResponse)
	if err != nil {
		return nil, err
	}

	manufacturers := make([]*common.Manufacturer, len(manufacturerResponse.Results))
	for i := range manufacturerResponse.Results {
		manufacturers[i] = &manufacturerResponse.Results[i]
	}
	api.Logger.Debug("Successfully received manufacturers: ", manufacturerResponse.Results)

	return manufacturers, nil
}

// PATCH /api/dcim/manufacturers/{id}/
func (api *NetboxAPI) PatchManufacturer(diffMap map[string]interface{}, manufacturerId int) (*common.Manufacturer, error) {
	api.Logger.Debug("Patching manufacturer ", manufacturerId, " with data: ", diffMap, " in NetBox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/dcim/manufacturers/%d/", manufacturerId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var manufacturerResponse common.Manufacturer
	err = json.Unmarshal(response.Body, &manufacturerResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched manufacturer: ", &manufacturerResponse)
	return &manufacturerResponse, nil
}

// POST /api/dcim/manufacturers/
func (api *NetboxAPI) CreateManufacturer(manufacturer *common.Manufacturer) (*common.Manufacturer, error) {
	api.Logger.Debug("Creating manufacturer with data: ", manufacturer, " in NetBox")

	requestBody, err := utils.NetboxJsonMarshal(manufacturer)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/dcim/manufacturers/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var manufacturerResponse common.Manufacturer
	err = json.Unmarshal(response.Body, &manufacturerResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created manufacturer: ", manufacturerResponse)

	return &manufacturerResponse, nil
}

type PlatformResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []common.Platform `json:"results"`
}

// GET /api/dcim/platforms/?limit=0
func (api *NetboxAPI) GetAllPlatforms() ([]*common.Platform, error) {
	api.Logger.Debug("Getting all platforms from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/platforms/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var platformResponse PlatformResponse
	err = json.Unmarshal(response.Body, &platformResponse)
	if err != nil {
		return nil, err
	}

	platforms := make([]*common.Platform, len(platformResponse.Results))
	for i := range platformResponse.Results {
		platforms[i] = &platformResponse.Results[i]
	}
	api.Logger.Debug("Successfully received platforms from netbox: ", platformResponse.Results)

	return platforms, nil
}

// PATCH /api/dcim/platforms/{id}/
func (api *NetboxAPI) PatchPlatform(diffMap map[string]interface{}, platformId int) (*common.Platform, error) {
	api.Logger.Debug("Patching platform ", platformId, " with data: ", diffMap, " in NetBox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/dcim/platforms/%d/", platformId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var platformResponse common.Platform
	err = json.Unmarshal(response.Body, &platformResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched platform: ", platformResponse)
	return &platformResponse, nil
}

// POSST /api/dcim/platforms/
func (api *NetboxAPI) CreatePlatform(platform *common.Platform) (*common.Platform, error) {
	api.Logger.Debug("Creating platform in NetBox with data: ", platform)

	requestBody, err := utils.NetboxJsonMarshal(platform)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/dcim/platforms/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var platformResponse common.Platform
	err = json.Unmarshal(response.Body, &platformResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created Platform: ", platformResponse)

	return &platformResponse, nil
}

type DeviceTypeResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []dcim.DeviceType `json:"results"`
}

// GET /api/dcim/device-types/?limit=0
func (api *NetboxAPI) GetAllDeviceTypes() ([]*dcim.DeviceType, error) {
	api.Logger.Debug("Getting all device types from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/device-types/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var deviceTypeResponse DeviceTypeResponse
	err = json.Unmarshal(response.Body, &deviceTypeResponse)
	if err != nil {
		return nil, err
	}

	deviceTypes := make([]*dcim.DeviceType, len(deviceTypeResponse.Results))
	for i := range deviceTypeResponse.Results {
		deviceTypes[i] = &deviceTypeResponse.Results[i]
	}
	api.Logger.Debug("Successfully received device types: ", deviceTypeResponse.Results)

	return deviceTypes, nil
}

// PATCH /api/dcim/device-types/{id}/
func (api *NetboxAPI) PatchDeviceType(diffMap map[string]interface{}, deviceTypeId int) (*dcim.DeviceType, error) {
	api.Logger.Debug("Patching device type ", deviceTypeId, " with data: ", diffMap, " in NetBox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/dcim/device-types/%d/", deviceTypeId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var deviceTypeResponse dcim.DeviceType
	err = json.Unmarshal(response.Body, &deviceTypeResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched device type: ", deviceTypeResponse)
	return &deviceTypeResponse, nil
}

// POST /api/dcim/device-types/
func (api *NetboxAPI) CreateDeviceType(deviceType *dcim.DeviceType) (*dcim.DeviceType, error) {
	api.Logger.Debug("Creating device type in NetBox with data: ", deviceType)

	requestBody, err := utils.NetboxJsonMarshal(deviceType)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/dcim/device-types/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var deviceTypeResponse dcim.DeviceType
	err = json.Unmarshal(response.Body, &deviceTypeResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created device type: ", deviceTypeResponse)

	return &deviceTypeResponse, nil
}

type InterfaceResponse struct {
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []dcim.Interface `json:"results"`
}

// GET /api/dcim/interfaces/?limit=0
func (api *NetboxAPI) GetAllInterfaces() ([]*dcim.Interface, error) {
	api.Logger.Debug("Getting all interfaces from NetBox")

	response, err := api.doRequest(MethodGet, "/api/dcim/interfaces/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var interfaceResponse InterfaceResponse
	err = json.Unmarshal(response.Body, &interfaceResponse)
	if err != nil {
		return nil, err
	}

	interfaces := make([]*dcim.Interface, len(interfaceResponse.Results))
	for i := range interfaceResponse.Results {
		interfaces[i] = &interfaceResponse.Results[i]
	}
	api.Logger.Debug("Successfully received interfaces: ", interfaceResponse.Results)

	return interfaces, nil
}

// PATCH /api/dcim/interfaces/{id}/
func (api *NetboxAPI) PatchInterface(diffMap map[string]interface{}, interfaceId int) (*dcim.Interface, error) {
	api.Logger.Debug("Patching interface with id ", interfaceId, " with data: ", diffMap, " in NetBox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/dcim/interfaces/%d/", interfaceId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var interfaceResponse dcim.Interface
	err = json.Unmarshal(response.Body, &interfaceResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched interface: ", interfaceResponse)
	return &interfaceResponse, nil
}

// POST /api/dcim/interfaces/
func (api *NetboxAPI) CreateInterface(interf *dcim.Interface) (*dcim.Interface, error) {
	api.Logger.Debug("Creating interface in NetBox with data: ", interf)

	requestBody, err := utils.NetboxJsonMarshal(interf)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/dcim/interfaces/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var interfaceResponse dcim.Interface
	err = json.Unmarshal(response.Body, &interfaceResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created interface: ", interfaceResponse)

	return &interfaceResponse, nil
}
