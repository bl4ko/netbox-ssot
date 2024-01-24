package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Standard response format from Netbox's API
type Response[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

// GetAll queries all objects of type T from Netbox's API.
// It is querying objects via pagination of limit=100.
//
// extraParams in a format of: &extraParam1=...&extraParam2=...
func GetAll[T any](api *NetboxAPI, path string, extraParams string) ([]T, error) {
	var allResults []T
	limit := 100
	offset := 0

	var zeroValueT T
	api.Logger.Debugf("Getting all %T from Netbox", zeroValueT)

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

	api.Logger.Debugf("Successfully received all %T: %v", zeroValueT, allResults)

	return allResults, nil
}
