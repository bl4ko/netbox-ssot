package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/extras"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
)

type TagResponse struct {
	Count    int          `json:"count,omitempty"`
	Next     int          `json:"next,omitempty"`
	Previous int          `json:"previous,omitempty"`
	Results  []extras.Tag `json:"results,omitempty"`
}

// /api/extras/tags
func (api *NetboxAPI) GetAllTags() ([]*extras.Tag, error) {
	api.Logger.Debug("Getting all tags from NetBox")

	response, err := api.doRequest(MethodGet, "/api/extras/tags/", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var tagResponse TagResponse
	err = json.Unmarshal(response.Body, &tagResponse)
	if err != nil {
		return nil, err
	}

	tags := make([]*extras.Tag, len(tagResponse.Results))
	for i := range tagResponse.Results {
		tags[i] = &tagResponse.Results[i]
	}
	api.Logger.Debug("Tags: ", tagResponse.Results)

	return tags, nil
}

// GET /api/extras/tags?name={tag_name}
func (api *NetboxAPI) GetTagByName(name string) (*extras.Tag, error) {
	api.Logger.Debug("Getting tag by name from NetBox")

	response, err := api.doRequest(MethodGet, fmt.Sprintf("/api/extras/tags?name=%s", name), nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var tagResponse TagResponse
	err = json.Unmarshal(response.Body, &tagResponse)
	if err != nil {
		return nil, err
	}

	if len(tagResponse.Results) == 0 {
		return nil, nil
	}

	tag := tagResponse.Results[0]
	api.Logger.Debug("Tag: ", tag)

	return &tag, nil
}

// POST /api/extras/tags/ -d '{"name": "netbox-ssot", "slug": "netbox-ssot"}'
func (api *NetboxAPI) CreateTag(tag *extras.Tag) (*extras.Tag, error) {
	api.Logger.Debug("Creating tag in NetBox: ", tag)

	requestBody, err := utils.NetboxJsonMarshal(tag)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/extras/tags/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var tagResponse extras.Tag
	err = json.Unmarshal(response.Body, &tagResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Tag: ", tagResponse)

	return &tagResponse, nil
}

// PUT /api/extras/tags/{tag_id}/ -d '{"name": "netbox-ssot", "slug": "netbox-ssot", ...}'
func (api *NetboxAPI) UpdateTag(tag *extras.Tag) (*extras.Tag, error) {
	api.Logger.Debug("Updating tag in netbox with data: ", tag)

	// Remove ID of the tag from the request body
	requestBody, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPut, fmt.Sprintf("/api/extras/tags/%d/", tag.ID), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var tagResponse extras.Tag
	err = json.Unmarshal(response.Body, &tagResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Updated tag: ", tagResponse)
	return &tagResponse, nil
}

// PATCH /api/extras/tags/{tag_id}/ -d '{"name": "netbox-ssot", "slug": "netbox-ssot", ...}'
func (api *NetboxAPI) PatchTag(diffMap map[string]interface{}, tagId int) (*extras.Tag, error) {
	api.Logger.Debug("Patching tag[", fmt.Sprintf("%d", tagId), "] in netbox with data: ", diffMap)

	// Remove ID of the tag from the request body
	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/extras/tags/%d/", tagId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var tagResponse extras.Tag
	err = json.Unmarshal(response.Body, &tagResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Patched tag: ", tagResponse)
	return &tagResponse, nil
}

type CustomFieldResponse struct {
	Count    int                  `json:"count,omitempty"`
	Next     int                  `json:"next,omitempty"`
	Previous int                  `json:"previous,omitempty"`
	Results  []extras.CustomField `json:"results,omitempty"`
}

// GET /api/extras/custom-fields/
func (api *NetboxAPI) GetAllCustomFields() ([]*extras.CustomField, error) {
	api.Logger.Debug("Getting all custom fields from NetBox")

	response, err := api.doRequest(MethodGet, "/api/extras/custom-fields/", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var customFieldResponse CustomFieldResponse
	err = json.Unmarshal(response.Body, &customFieldResponse)
	if err != nil {
		return nil, err
	}

	customFields := make([]*extras.CustomField, len(customFieldResponse.Results))
	for i := range customFieldResponse.Results {
		customFields[i] = &customFieldResponse.Results[i]
	}
	api.Logger.Debug("Custom fields: ", customFieldResponse.Results)

	return customFields, nil
}

// PATCH /api/extras/custom-fields/{custom_field_id}/ -d '{...}'
func (api *NetboxAPI) PatchCustomField(diffMap map[string]interface{}, customFieldId int) (*extras.CustomField, error) {
	api.Logger.Debug("Patching custom field[", fmt.Sprintf("%d", customFieldId), "] in netbox with data: ", diffMap)

	// Remove ID of the custom field from the request body
	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/extras/custom-fields/%d/", customFieldId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var customFieldResponse extras.CustomField
	err = json.Unmarshal(response.Body, &customFieldResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Patched custom field: ", customFieldResponse)
	return &customFieldResponse, nil
}

// CREATE /api/extras/custom-fields/ -d '{...}'
func (api *NetboxAPI) CreateCustomField(customField *extras.CustomField) (*extras.CustomField, error) {
	api.Logger.Debug("Creating custom field in NetBox: ", customField)

	requestBody, err := json.Marshal(customField)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/extras/custom-fields/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var customFieldResponse extras.CustomField
	err = json.Unmarshal(response.Body, &customFieldResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Custom field: ", customFieldResponse)

	return &customFieldResponse, nil
}
