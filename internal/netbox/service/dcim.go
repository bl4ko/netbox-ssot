package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// PATCH /api/dcim/devices/{id}/
func (api *NetboxAPI) PatchDevice(diffMap map[string]interface{}, deviceId int) (*objects.Device, error) {
	api.Logger.Debug("Patching device ", deviceId, " with data: ", diffMap, " in Netbox")

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

	var deviceResponse objects.Device
	err = json.Unmarshal(response.Body, &deviceResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched device: ", deviceResponse)
	return &deviceResponse, nil
}

// POST /api/dcim/devices/
func (api *NetboxAPI) CreateDevice(device *objects.Device) (*objects.Device, error) {
	api.Logger.Debug("Creating device in Netbox with data: ", device)

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

	var deviceResponse objects.Device
	err = json.Unmarshal(response.Body, &deviceResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created device: ", deviceResponse)

	return &deviceResponse, nil
}

// POST /api/dcim/device-roles/
func (api *NetboxAPI) CreateDeviceRole(deviceRole *objects.DeviceRole) (*objects.DeviceRole, error) {
	api.Logger.Debug("Creating device role with data", deviceRole, " in Netbox")

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

	var deviceRoleResponse objects.DeviceRole
	err = json.Unmarshal(response.Body, &deviceRoleResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created device role: ", deviceRoleResponse)

	return &deviceRoleResponse, nil
}

// PATCH /api/dcim/device-roles/{id}/
func (api *NetboxAPI) PatchDeviceRole(diffMap map[string]interface{}, id int) (*objects.DeviceRole, error) {
	api.Logger.Debug("Patching device role ", id, " with data: ", diffMap, " in Netbox")

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

	var deviceRoleResponse objects.DeviceRole
	err = json.Unmarshal(response.Body, &deviceRoleResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched device role: ", deviceRoleResponse)

	return &deviceRoleResponse, nil
}

// PATCH /api/dcim/manufacturers/{id}/
func (api *NetboxAPI) PatchManufacturer(diffMap map[string]interface{}, manufacturerId int) (*objects.Manufacturer, error) {
	api.Logger.Debug("Patching manufacturer ", manufacturerId, " with data: ", diffMap, " in Netbox")

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

	var manufacturerResponse objects.Manufacturer
	err = json.Unmarshal(response.Body, &manufacturerResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched manufacturer: ", &manufacturerResponse)
	return &manufacturerResponse, nil
}

// POST /api/dcim/manufacturers/
func (api *NetboxAPI) CreateManufacturer(manufacturer *objects.Manufacturer) (*objects.Manufacturer, error) {
	api.Logger.Debug("Creating manufacturer with data: ", manufacturer, " in Netbox")

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

	var manufacturerResponse objects.Manufacturer
	err = json.Unmarshal(response.Body, &manufacturerResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created manufacturer: ", manufacturerResponse)

	return &manufacturerResponse, nil
}

// PATCH /api/dcim/platforms/{id}/
func (api *NetboxAPI) PatchPlatform(diffMap map[string]interface{}, platformId int) (*objects.Platform, error) {
	api.Logger.Debug("Patching platform ", platformId, " with data: ", diffMap, " in Netbox")

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

	var platformResponse objects.Platform
	err = json.Unmarshal(response.Body, &platformResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched platform: ", platformResponse)
	return &platformResponse, nil
}

// POST /api/dcim/platforms/
func (api *NetboxAPI) CreatePlatform(platform *objects.Platform) (*objects.Platform, error) {
	api.Logger.Debug("Creating platform in Netbox with data: ", platform)

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

	var platformResponse objects.Platform
	err = json.Unmarshal(response.Body, &platformResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created Platform: ", platformResponse)

	return &platformResponse, nil
}

// PATCH /api/dcim/device-types/{id}/
func (api *NetboxAPI) PatchDeviceType(diffMap map[string]interface{}, deviceTypeId int) (*objects.DeviceType, error) {
	api.Logger.Debug("Patching device type ", deviceTypeId, " with data: ", diffMap, " in Netbox")

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

	var deviceTypeResponse objects.DeviceType
	err = json.Unmarshal(response.Body, &deviceTypeResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched device type: ", deviceTypeResponse)
	return &deviceTypeResponse, nil
}

// POST /api/dcim/device-types/
func (api *NetboxAPI) CreateDeviceType(deviceType *objects.DeviceType) (*objects.DeviceType, error) {
	api.Logger.Debug("Creating device type in Netbox with data: ", deviceType)

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

	var deviceTypeResponse objects.DeviceType
	err = json.Unmarshal(response.Body, &deviceTypeResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created device type: ", deviceTypeResponse)

	return &deviceTypeResponse, nil
}

// PATCH /api/dcim/interfaces/{id}/
func (api *NetboxAPI) PatchInterface(diffMap map[string]interface{}, interfaceId int) (*objects.Interface, error) {
	api.Logger.Debug("Patching interface with id ", interfaceId, " with data: ", diffMap, " in Netbox")

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

	var interfaceResponse objects.Interface
	err = json.Unmarshal(response.Body, &interfaceResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched interface: ", interfaceResponse)
	return &interfaceResponse, nil
}

// POST /api/dcim/interfaces/
func (api *NetboxAPI) CreateInterface(iface *objects.Interface) (*objects.Interface, error) {
	api.Logger.Debug("Creating interface in Netbox with data: ", iface)

	requestBody, err := utils.NetboxJsonMarshal(iface)
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

	var interfaceResponse objects.Interface
	err = json.Unmarshal(response.Body, &interfaceResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created interface: ", interfaceResponse)

	return &interfaceResponse, nil
}
