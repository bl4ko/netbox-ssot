package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Standard response format from Netbox's API
type Response[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

var type2path = map[reflect.Type]string{
	reflect.TypeOf((*objects.VlanGroup)(nil)).Elem():    VlanGroupApiPath,
	reflect.TypeOf((*objects.Vlan)(nil)).Elem():         VlanApiPath,
	reflect.TypeOf((*objects.IPAddress)(nil)).Elem():    IpAddressApiPath,
	reflect.TypeOf((*objects.ClusterType)(nil)).Elem():  ClusterTypeApiPath,
	reflect.TypeOf((*objects.ClusterGroup)(nil)).Elem(): ClusterGroupApiPath,
	reflect.TypeOf((*objects.Cluster)(nil)).Elem():      ClusterApiPath,
	reflect.TypeOf((*objects.VM)(nil)).Elem():           VirtualMachineApiPath,
	reflect.TypeOf((*objects.VMInterface)(nil)).Elem():  VMInterfaceApiPath,
	reflect.TypeOf((*objects.Device)(nil)).Elem():       DeviceApiPath,
	reflect.TypeOf((*objects.DeviceRole)(nil)).Elem():   DeviceRoleApiPath,
	reflect.TypeOf((*objects.DeviceType)(nil)).Elem():   DeviceTypeApiPath,
	reflect.TypeOf((*objects.Interface)(nil)).Elem():    InterfaceApiPath,
	reflect.TypeOf((*objects.Site)(nil)).Elem():         SiteApiPath,
	reflect.TypeOf((*objects.Manufacturer)(nil)).Elem(): ManufacturerApiPath,
	reflect.TypeOf((*objects.Platform)(nil)).Elem():     PlatformApiPath,
	reflect.TypeOf((*objects.Tenant)(nil)).Elem():       TenantApiPath,
	reflect.TypeOf((*objects.ContactGroup)(nil)).Elem(): ContactGroupApiPath,
	reflect.TypeOf((*objects.ContactRole)(nil)).Elem():  ContactRoleApiPath,
	reflect.TypeOf((*objects.Contact)(nil)).Elem():      ContactApiPath,
	reflect.TypeOf((*objects.CustomField)(nil)).Elem():  CustomFieldApiPath,
	reflect.TypeOf((*objects.Tag)(nil)).Elem():          TagApiPath,
}

// GetAll queries all objects of type T from Netbox's API.
// It is querying objects via pagination of limit=100.
//
// extraParams in a string format of: &extraParam1=...&extraParam2=...
func GetAll[T any](api *NetboxAPI, extraParams string) ([]T, error) {
	var allResults []T
	var dummy T // Dummy variable for extracting type of generic
	path := type2path[reflect.TypeOf(dummy)]
	limit := 100
	offset := 0

	api.Logger.Debugf("Getting all %T from Netbox", dummy)

	for {
		queryPath := fmt.Sprintf("%s?limit=%d&offset=%d%s", path, limit, offset, extraParams)
		response, err := api.doRequest(MethodGet, queryPath, nil)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, response.Body)
		}

		var responseObj Response[T]
		err = json.Unmarshal(response.Body, &responseObj)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, responseObj.Results...)

		if responseObj.Next == nil {
			break
		}
		offset += limit
	}

	api.Logger.Debugf("Successfully received all %T: %v", dummy, allResults)

	return allResults, nil
}

// Patch func patches the object of type T, with the given api path and body.
// Path of the object (must contain the id), for example /api/dcim/devices/1/
func Patch[T any](api *NetboxAPI, objectId int, body map[string]interface{}) (*T, error) {
	var dummy T // dummy variable for printf
	path := type2path[reflect.TypeOf(dummy)]
	path = fmt.Sprintf("%s%d/", path, objectId)
	api.Logger.Debugf("Patching %T with path %s with data: %v", dummy, path, body)

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, path, requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var objectResponse T
	err = json.Unmarshal(response.Body, &objectResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debugf("Successfully patched %T: %v", dummy, objectResponse)
	return &objectResponse, nil
}

// Create func creates the new NetboxObject of type T, with the given api path and body.
func Create[T any](api *NetboxAPI, object *T) (*T, error) {
	var dummy T // dummy variable for printf
	path := type2path[reflect.TypeOf(dummy)]
	api.Logger.Debugf("Creating %T with path %s with data: %v", dummy, path, object)

	requestBody, err := utils.NetboxJsonMarshal(object)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPost, path, requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var objectResponse T
	err = json.Unmarshal(response.Body, &objectResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debugf("Successfully created %T: %v", dummy, objectResponse)
	return &objectResponse, nil
}

func (api *NetboxAPI) BulkDeleteObjects(objectPath string, idSet map[int]bool) error {
	// Netbox API supports only JSON request body in the following format:
	// [ {"id": 1}, {"id": 2}, {"id": 3} ]
	netboxFormatIds := make([]map[string]int, 0)
	for id := range idSet {
		netboxFormatIds = append(netboxFormatIds, map[string]int{"id": id})
	}

	requestBody, err := json.Marshal(netboxFormatIds)
	if err != nil {
		return err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodDelete, objectPath, requestBodyBuffer)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	api.Logger.Debug("Successfully deleted objects")
	return nil
}
