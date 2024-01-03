package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/objects"
)

type TenantResponse struct {
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []objects.Tenant `json:"results"`
}

// GET /api/tenancy/tenants/?limit=0
func (api *NetboxAPI) GetAllTenants() ([]*objects.Tenant, error) {
	api.Logger.Debug("Getting all tenants from Netbox")

	response, err := api.doRequest(MethodGet, "/api/tenancy/tenants/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var tenantResponse TenantResponse
	err = json.Unmarshal(response.Body, &tenantResponse)
	if err != nil {
		return nil, err
	}

	tenants := make([]*objects.Tenant, len(tenantResponse.Results))
	for i := range tenantResponse.Results {
		tenants[i] = &tenantResponse.Results[i]
	}
	api.Logger.Debug("Tenants: ", tenantResponse.Results)

	return tenants, nil
}
