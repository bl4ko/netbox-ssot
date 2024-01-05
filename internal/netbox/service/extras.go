package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

type TagResponse struct {
	Count    int           `json:"count,omitempty"`
	Next     *string       `json:"next,omitempty"`
	Previous *string       `json:"previous,omitempty"`
	Results  []objects.Tag `json:"results,omitempty"`
}

// GET /api/extras/tags/?limit=0
func (api *NetboxAPI) GetAllTags() ([]*objects.Tag, error) {
	api.Logger.Debug("Getting all tags from Netbox")

	response, err := api.doRequest(MethodGet, "/api/extras/tags/?limit=0", nil)
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

	tags := make([]*objects.Tag, len(tagResponse.Results))
	for i := range tagResponse.Results {
		tags[i] = &tagResponse.Results[i]
	}
	api.Logger.Debug("Tags: ", tagResponse.Results)

	return tags, nil
}

// GET /api/extras/tags?name={tag_name}
func (api *NetboxAPI) GetTagByName(name string) (*objects.Tag, error) {
	api.Logger.Debug("Getting tag by name from Netbox")

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
func (api *NetboxAPI) CreateTag(tag *objects.Tag) (*objects.Tag, error) {
	api.Logger.Debug("Creating tag in Netbox: ", tag)

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

	var tagResponse objects.Tag
	err = json.Unmarshal(response.Body, &tagResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Tag: ", tagResponse)

	return &tagResponse, nil
}

// PUT /api/extras/tags/{tag_id}/ -d '{"name": "netbox-ssot", "slug": "netbox-ssot", ...}'
func (api *NetboxAPI) UpdateTag(tag *objects.Tag) (*objects.Tag, error) {
	api.Logger.Debug("Updating tag in netbox with data: ", tag)

	// Remove ID of the tag from the request body
	requestBody, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPut, fmt.Sprintf("/api/extras/tags/%d/", tag.Id), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var tagResponse objects.Tag
	err = json.Unmarshal(response.Body, &tagResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Updated tag: ", tagResponse)
	return &tagResponse, nil
}

// PATCH /api/extras/tags/{tag_id}/ -d '{"name": "netbox-ssot", "slug": "netbox-ssot", ...}'
func (api *NetboxAPI) PatchTag(diffMap map[string]interface{}, tagId int) (*objects.Tag, error) {
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

	var tagResponse objects.Tag
	err = json.Unmarshal(response.Body, &tagResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Patched tag: ", tagResponse)
	return &tagResponse, nil
}

type CustomFieldResponse struct {
	Count    int                   `json:"count,omitempty"`
	Next     *string               `json:"next,omitempty"`
	Previous *string               `json:"previous,omitempty"`
	Results  []objects.CustomField `json:"results,omitempty"`
}

// GET /api/extras/custom-fields/?limit=0
func (api *NetboxAPI) GetAllCustomFields() ([]*objects.CustomField, error) {
	api.Logger.Debug("Getting all custom fields from Netbox")

	response, err := api.doRequest(MethodGet, "/api/extras/custom-fields/?limit=0", nil)
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

	customFields := make([]*objects.CustomField, len(customFieldResponse.Results))
	for i := range customFieldResponse.Results {
		customFields[i] = &customFieldResponse.Results[i]
	}
	api.Logger.Debug("Custom fields: ", customFieldResponse.Results)

	return customFields, nil
}

// PATCH /api/extras/custom-fields/{custom_field_id}/ -d '{...}'
func (api *NetboxAPI) PatchCustomField(diffMap map[string]interface{}, customFieldId int) (*objects.CustomField, error) {
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

	var customFieldResponse *objects.CustomField
	err = json.Unmarshal(response.Body, customFieldResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Patched custom field: ", customFieldResponse)
	return customFieldResponse, nil
}

// CREATE /api/extras/custom-fields/ -d '{...}'
func (api *NetboxAPI) CreateCustomField(customField *objects.CustomField) (*objects.CustomField, error) {
	api.Logger.Debug("Creating custom field in Netbox: ", customField)

	requestBody, err := utils.NetboxJsonMarshal(customField)
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

	var customFieldResponse *objects.CustomField
	err = json.Unmarshal(response.Body, customFieldResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Custom field: ", customFieldResponse)

	return customFieldResponse, nil
}
