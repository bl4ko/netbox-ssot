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

// Standard response format from Netbox's API.
type Response[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

var type2path = map[reflect.Type]string{
	reflect.TypeOf((*objects.VlanGroup)(nil)).Elem():         VlanGroupsAPIPath,
	reflect.TypeOf((*objects.Vlan)(nil)).Elem():              VlansAPIPath,
	reflect.TypeOf((*objects.IPAddress)(nil)).Elem():         IPAddressesAPIPath,
	reflect.TypeOf((*objects.ClusterType)(nil)).Elem():       ClusterTypesAPIPath,
	reflect.TypeOf((*objects.ClusterGroup)(nil)).Elem():      ClusterGroupsAPIPath,
	reflect.TypeOf((*objects.Cluster)(nil)).Elem():           ClustersAPIPath,
	reflect.TypeOf((*objects.VM)(nil)).Elem():                VirtualMachinesAPIPath,
	reflect.TypeOf((*objects.VMInterface)(nil)).Elem():       VMInterfacesAPIPath,
	reflect.TypeOf((*objects.Device)(nil)).Elem():            DevicesAPIPath,
	reflect.TypeOf((*objects.DeviceRole)(nil)).Elem():        DeviceRolesAPIPath,
	reflect.TypeOf((*objects.DeviceType)(nil)).Elem():        DeviceTypesAPIPath,
	reflect.TypeOf((*objects.Interface)(nil)).Elem():         InterfacesAPIPath,
	reflect.TypeOf((*objects.Site)(nil)).Elem():              SitesAPIPath,
	reflect.TypeOf((*objects.Manufacturer)(nil)).Elem():      ManufacturersAPIPath,
	reflect.TypeOf((*objects.Platform)(nil)).Elem():          PlatformsAPIPath,
	reflect.TypeOf((*objects.Tenant)(nil)).Elem():            TenantsAPIPath,
	reflect.TypeOf((*objects.ContactGroup)(nil)).Elem():      ContactGroupsAPIPath,
	reflect.TypeOf((*objects.ContactRole)(nil)).Elem():       ContactRolesAPIPath,
	reflect.TypeOf((*objects.Contact)(nil)).Elem():           ContactsAPIPath,
	reflect.TypeOf((*objects.CustomField)(nil)).Elem():       CustomFieldsAPIPath,
	reflect.TypeOf((*objects.Tag)(nil)).Elem():               TagsAPIPath,
	reflect.TypeOf((*objects.ContactAssignment)(nil)).Elem(): ContactAssignmentsAPIPath,
	reflect.TypeOf((*objects.Prefix)(nil)).Elem():            PrefixesAPIPath,
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
		api.Logger.Debugf("Getting %T with limit=%d and offset=%d", dummy, limit, offset)
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
// Path of the object (must contain the id), for example /api/dcim/devices/1/.
func Patch[T any](api *NetboxAPI, objectID int, body map[string]interface{}) (*T, error) {
	var dummy T // dummy variable for printf
	path := type2path[reflect.TypeOf(dummy)]
	path = fmt.Sprintf("%s%d/", path, objectID)
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

	requestBody, err := utils.NetboxJSONMarshal(object)
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

// Function that deletes object on path objectPath.
// It deletes objects in pages of 50 so we don't stress
// the API too much.
func (api *NetboxAPI) BulkDeleteObjects(objectPath string, idSet map[int]bool) error {
	const pageSize = 50

	// Convert the map to a slice for easier slicing.
	ids := make([]int, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	for i := 0; i < len(ids); i += pageSize {
		api.Logger.Debugf("Deleting %s with pagesize=%d and offset=%d", objectPath, pageSize, i)
		end := i + pageSize
		if end > len(ids) {
			end = len(ids)
		}

		// Netbox API supports only JSON request body in the following format:
		// [ {"id": 1}, {"id": 2}, {"id": 3} ]
		netboxFormatIDs := make([]map[string]int, 0, end-i)
		for _, id := range ids[i:end] {
			netboxFormatIDs = append(netboxFormatIDs, map[string]int{"id": id})
		}

		requestBody, err := json.Marshal(netboxFormatIDs)
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
	}
	api.Logger.Debugf("Successfully deleted all objects of path %s", objectPath)

	return nil
}
