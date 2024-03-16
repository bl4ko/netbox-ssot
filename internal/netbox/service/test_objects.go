package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

var MockDefaultSsotTag = &objects.Tag{
	ID:   0,
	Name: constants.DefaultSourceName,
}

// Hardcoded mock API responses for tags endpoint.
var (
	MockTagsGetResponse = Response[objects.Tag]{
		Count:    2, //nolint:gomnd
		Next:     nil,
		Previous: nil,
		Results: []objects.Tag{
			{
				ID:          0,
				Name:        "Source: proxmox",
				Slug:        "source-proxmox",
				Color:       "9e9e9e",
				Description: "Automatically created tag by netbox-ssot for source proxmox",
			},
			{
				ID:          1,
				Name:        "netbox-ssot",
				Slug:        "netbox-ssot",
				Color:       "00add8",
				Description: "Tag used by netbox-ssot to mark devices that are managed by it",
			},
		},
	}
	MockTagPatchResponse = objects.Tag{
		ID:          1,
		Name:        "netbox-ssot",
		Slug:        "netbox-ssot",
		Color:       "00add8",
		Description: "patched description",
	}
	MockTagCreateResponse = objects.Tag{
		ID:          1,
		Name:        "netbox-ssot",
		Slug:        "netbox-ssot",
		Color:       "00add8",
		Description: "created description",
	}
)

// Hardcoded mock api return values for tenant endpoint.
var (
	MockTenantsGetResponse = Response[objects.Tenant]{
		Count:    2, //nolint:gomnd
		Next:     nil,
		Previous: nil,
		Results: []objects.Tenant{
			{
				NetboxObject: objects.NetboxObject{
					ID: 1,
					Tags: []*objects.Tag{
						MockDefaultSsotTag,
					},
				},
				Name: "MockTenant1",
				Slug: "mock-tenant-1",
			},
			{
				NetboxObject: objects.NetboxObject{
					ID: 2, //nolint:gomnd
					Tags: []*objects.Tag{
						MockDefaultSsotTag,
					},
				},
				Name: "MockTenant2",
				Slug: "mock-tenant-2",
			},
		},
	}
	MockTenantCreateResponse = objects.Tenant{
		NetboxObject: objects.NetboxObject{
			ID: 3, //nolint:gomnd
		},
		Name: "MockTenant3",
		Slug: "mock-tenant-3",
	}
	MockTenantPatchResponse = objects.Tenant{
		NetboxObject: objects.NetboxObject{
			ID: 1,
		},
		Name: "MockPatched",
		Slug: "mock-patched-tenant",
	}
)

// Hardcoded mock api return values for site endpoint.
var (
	MockSitesGetResponse = Response[objects.Site]{
		Count:    2, //nolint:gomnd
		Next:     nil,
		Previous: nil,
		Results: []objects.Site{
			{
				NetboxObject: objects.NetboxObject{
					ID: 1,
					Tags: []*objects.Tag{
						MockDefaultSsotTag,
					},
				},
				Name: "MockSite1",
				Slug: "mock-site-1",
			},
			{
				NetboxObject: objects.NetboxObject{
					ID: 2, //nolint:gomnd
					Tags: []*objects.Tag{
						MockDefaultSsotTag,
					},
				},
				Name: "MockSite2",
				Slug: "mock-site-2",
			},
		},
	}
	MockSiteCreateResponse = objects.Site{
		NetboxObject: objects.NetboxObject{
			ID: 3, //nolint:gomnd
		},
		Name: "MockSite3",
		Slug: "mock-site-3",
	}
	MockSitePatchResponse = objects.Site{
		NetboxObject: objects.NetboxObject{
			ID: 1,
		},
		Name: "MockSitePatched",
		Slug: "mock-site-patched",
	}
)

const (
	MockVersionResponseJSON = "{\"django-version\": \"4.2.10\"}"
)

//nolint:gocyclo
func CreateMockServer() *httptest.Server {
	handler := http.NewServeMux()
	// Define handler for a specific path e.g., "/api/path"
	handler.HandleFunc("/api/status/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, MockVersionResponseJSON) // Mock JSON Response
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	handler.HandleFunc(constants.TagsAPIPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			tagStr, err := json.Marshal(MockTagPatchResponse)
			if err != nil {
				log.Printf("Error marshaling tag patch response: %v", err)
			}
			_, err = io.WriteString(w, string(tagStr))
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			tagsResponseStr, err := json.Marshal(MockTagsGetResponse)
			if err != nil {
				log.Printf("Error marshaling tags response: %v", err)
			}
			_, err = io.WriteString(w, string(tagsResponseStr))
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			tagStr, err := json.Marshal(MockTagCreateResponse)
			if err != nil {
				log.Printf("Error marshaling tag create response: %v", err)
			}
			_, err = io.WriteString(w, string(tagStr))
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			log.Printf("Wrong http method: %v", r.Method)
		}
	})

	handler.HandleFunc(constants.TenantsAPIPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			tenantStr, err := json.Marshal(MockTenantPatchResponse)
			if err != nil {
				log.Printf("Error marshaling tenant patch response: %v", err)
			}
			_, err = io.WriteString(w, string(tenantStr))
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			tenantsResponseStr, err := json.Marshal(MockTenantsGetResponse)
			if err != nil {
				log.Printf("Error marshaling tenants response: %v", err)
			}
			_, err = io.WriteString(w, string(tenantsResponseStr))
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			tenantStr, err := json.Marshal(MockTenantCreateResponse)
			if err != nil {
				log.Printf("Error marshaling tenant create response: %v", err)
			}
			_, err = io.WriteString(w, string(tenantStr))
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			log.Printf("Wrong http method: %v", r.Method)
		}
	})

	handler.HandleFunc(constants.SitesAPIPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			siteStr, err := json.Marshal(MockSitePatchResponse)
			if err != nil {
				log.Printf("Error marshaling site patch response: %v", err)
			}
			_, err = io.WriteString(w, string(siteStr))
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			siteResponseStr, err := json.Marshal(MockSitesGetResponse)
			if err != nil {
				log.Printf("Error marshaling sites response: %v", err)
			}
			_, err = io.WriteString(w, string(siteResponseStr))
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			siteStr, err := json.Marshal(MockSiteCreateResponse)
			if err != nil {
				log.Printf("Error marshaling site create response: %v", err)
			}
			_, err = io.WriteString(w, string(siteStr))
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			log.Printf("Wrong http method: %v", r.Method)
		}
	})

	handler.HandleFunc("/api/read-error", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // or any relevant status
		//nolint:all
		w.(http.Flusher).Flush() // Flush the headers to client
		// Do not write any body, let the client read from the FaultyReader
	})
	// Wildcard handler for all other paths
	handler.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error": "page not found"}`)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})
	// More handlers can be added here for different paths or methods
	return httptest.NewServer(handler)
}

var MockNetboxClient = &NetboxClient{
	HTTPClient: &http.Client{},
	Logger:     &logger.Logger{Logger: log.Default()},
	BaseURL:    "",
	APIToken:   "testtoken",
	Timeout:    constants.DefaultAPITimeout,
}

var FailingMockNetboxClient = &NetboxClient{
	HTTPClient: &http.Client{Transport: &FailingHTTPClient{}},
	Logger:     &logger.Logger{Logger: log.Default()},
	BaseURL:    "",
	APIToken:   "testtoken",
	Timeout:    constants.DefaultAPITimeout,
}

type FailingHTTPClient struct{}

func (m *FailingHTTPClient) RoundTrip(_ *http.Request) (*http.Response, error) {
	// Return an error to simulate a failure in the HTTP request
	return nil, fmt.Errorf("mock error")
}

var MockNetboxClientWithReadError = &NetboxClient{
	HTTPClient: &http.Client{Transport: &FailingHTTPClientRead{}},
	Logger:     &logger.Logger{Logger: log.Default()},
	BaseURL:    "",
	APIToken:   "testtoken",
	Timeout:    constants.DefaultAPITimeout,
}

type FailingHTTPClientRead struct{}

func (m *FailingHTTPClientRead) RoundTrip(_ *http.Request) (*http.Response, error) {
	// Simulate a response with a FaultyReader as its Body
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(&FaultyReader{}),
		Header:     make(http.Header),
	}, nil
}

type FaultyReader struct{}

func (m *FaultyReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("mock read error")
}
