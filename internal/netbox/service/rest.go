package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

type VersionResponse struct {
	NetboxVersion string `json:"netbox-version"`
}

// Standard response format from Netbox's API.
type Response[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

var type2path = map[reflect.Type]string{
	reflect.TypeOf((*objects.VlanGroup)(nil)).Elem():            constants.VlanGroupsAPIPath,
	reflect.TypeOf((*objects.Vlan)(nil)).Elem():                 constants.VlansAPIPath,
	reflect.TypeOf((*objects.IPAddress)(nil)).Elem():            constants.IPAddressesAPIPath,
	reflect.TypeOf((*objects.ClusterType)(nil)).Elem():          constants.ClusterTypesAPIPath,
	reflect.TypeOf((*objects.ClusterGroup)(nil)).Elem():         constants.ClusterGroupsAPIPath,
	reflect.TypeOf((*objects.Cluster)(nil)).Elem():              constants.ClustersAPIPath,
	reflect.TypeOf((*objects.VM)(nil)).Elem():                   constants.VirtualMachinesAPIPath,
	reflect.TypeOf((*objects.VMInterface)(nil)).Elem():          constants.VMInterfacesAPIPath,
	reflect.TypeOf((*objects.Device)(nil)).Elem():               constants.DevicesAPIPath,
	reflect.TypeOf((*objects.VirtualDeviceContext)(nil)).Elem(): constants.VirtualDeviceContextsAPIPath,
	reflect.TypeOf((*objects.DeviceRole)(nil)).Elem():           constants.DeviceRolesAPIPath,
	reflect.TypeOf((*objects.DeviceType)(nil)).Elem():           constants.DeviceTypesAPIPath,
	reflect.TypeOf((*objects.Interface)(nil)).Elem():            constants.InterfacesAPIPath,
	reflect.TypeOf((*objects.Site)(nil)).Elem():                 constants.SitesAPIPath,
	reflect.TypeOf((*objects.Manufacturer)(nil)).Elem():         constants.ManufacturersAPIPath,
	reflect.TypeOf((*objects.Platform)(nil)).Elem():             constants.PlatformsAPIPath,
	reflect.TypeOf((*objects.Tenant)(nil)).Elem():               constants.TenantsAPIPath,
	reflect.TypeOf((*objects.ContactGroup)(nil)).Elem():         constants.ContactGroupsAPIPath,
	reflect.TypeOf((*objects.ContactRole)(nil)).Elem():          constants.ContactRolesAPIPath,
	reflect.TypeOf((*objects.Contact)(nil)).Elem():              constants.ContactsAPIPath,
	reflect.TypeOf((*objects.CustomField)(nil)).Elem():          constants.CustomFieldsAPIPath,
	reflect.TypeOf((*objects.Tag)(nil)).Elem():                  constants.TagsAPIPath,
	reflect.TypeOf((*objects.ContactAssignment)(nil)).Elem():    constants.ContactAssignmentsAPIPath,
	reflect.TypeOf((*objects.Prefix)(nil)).Elem():               constants.PrefixesAPIPath,
}

// Function that queries and returns netbox version on success.
func GetVersion(ctx context.Context, netboxClient *NetboxClient) (string, error) {
	var versionResponse VersionResponse
	netboxClient.Logger.Debugf(ctx, "Getting netbox's version")
	response, err := netboxClient.doRequest(MethodGet, "/api/status", nil)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d: %s", response.StatusCode, response.Body)
	}

	err = json.Unmarshal(response.Body, &versionResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling body: %s", err)
	}
	return versionResponse.NetboxVersion, nil
}

// GetAll queries all objects of type T from Netbox's API.
// It is querying objects via pagination of limit=100.
//
// extraParams in a string format of: &extraParam1=...&extraParam2=...
func GetAll[T any](ctx context.Context, netboxClient *NetboxClient, extraParams string) ([]T, error) {
	var allResults []T
	var dummy T // Dummy variable for extracting type of generic
	path := type2path[reflect.TypeOf(dummy)]
	limit := 100
	offset := 0

	netboxClient.Logger.Debugf(ctx, "Getting all %T from Netbox", dummy)

	for {
		netboxClient.Logger.Debugf(ctx, "Getting %T with limit=%d and offset=%d", dummy, limit, offset)
		queryPath := fmt.Sprintf("%s?limit=%d&offset=%d%s", path, limit, offset, extraParams)
		response, err := netboxClient.doRequest(MethodGet, queryPath, nil)
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

	netboxClient.Logger.Debugf(ctx, "Successfully received all %T: %v", dummy, allResults)

	return allResults, nil
}

// Patch func patches the object of type T, with the given api path and body.
// Path of the object (must contain the id), for example /api/dcim/devices/1/.
func Patch[T any](ctx context.Context, netboxClient *NetboxClient, objectID int, body map[string]interface{}) (*T, error) {
	var dummy T // dummy variable for printf
	path := type2path[reflect.TypeOf(dummy)]
	path = fmt.Sprintf("%s%d/", path, objectID)
	netboxClient.Logger.Debugf(ctx, "Patching %T with path %s with data: %v", dummy, path, body)

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := netboxClient.doRequest(MethodPatch, path, requestBodyBuffer)
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

	netboxClient.Logger.Debugf(ctx, "Successfully patched %T: %v", dummy, objectResponse)
	return &objectResponse, nil
}

// Create func creates the new NetboxObject of type T, with the given api path and body.
func Create[T any](ctx context.Context, netboxClient *NetboxClient, object *T) (*T, error) {
	var dummy T // dummy variable for printf
	path := type2path[reflect.TypeOf(dummy)]
	netboxClient.Logger.Debugf(ctx, "Creating %T with path %s with data: %v", dummy, path, object)

	requestBody, err := utils.NetboxJSONMarshal(object)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := netboxClient.doRequest(MethodPost, path, requestBodyBuffer)
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

	netboxClient.Logger.Debugf(ctx, "Successfully created %T: %v", dummy, objectResponse)
	return &objectResponse, nil
}

// Function that deletes object on path objectPath.
// It deletes objects in pages of 50 so we don't stress
// the API too much.
func (nbClient *NetboxClient) BulkDeleteObjects(ctx context.Context, objectPath string, idSet map[int]bool) error {
	const pageSize = 50

	// Convert the map to a slice for easier slicing.
	ids := make([]int, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	for i := 0; i < len(ids); i += pageSize {
		nbClient.Logger.Debugf(ctx, "Deleting %s with pagesize=%d and offset=%d", objectPath, pageSize, i)
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
		response, err := nbClient.doRequest(MethodDelete, objectPath, requestBodyBuffer)
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusNoContent {
			return fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
		}
	}
	nbClient.Logger.Debugf(ctx, "Successfully deleted all objects of path %s", objectPath)

	return nil
}

// Function that deletes objectas on path objectPath.
// It deletes a single object at a time. It is alternative to bulk delete
// because if one delete fails other still go.
func (nbClient *NetboxClient) DeleteObject(ctx context.Context, objectPath string, id int) error {
	nbClient.Logger.Debugf(ctx, "Deleting object with id %d on route %s", id, objectPath)

	response, err := nbClient.doRequest(MethodDelete, fmt.Sprintf("%s%d/", objectPath, id), nil)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}
	return nil
}
