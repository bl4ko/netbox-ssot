package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *NetboxAPI) BulkDeleteObjects(objectPath string, idSet map[int]bool) error {
	api.Logger.Debug("Deleting objects with path ", objectPath, " from Netbox with ids: ", idSet)

	// Netbox API supports only JSON request body in the following format:
	// [ {"id": 1}, {"id": 2}, {"id": 3} ]
	netboxFormatIds := make([]map[string]int, len(idSet))
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
