//go:build test

package testing

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

// Hardcoded mock API responses for tags endpoint.
const (
	MockTagGetResponse = `{
    "count": 3,
    "next": null,
    "previous": null,
    "results": [
        {
            "id": 2,
            "url": "http://ubtest.home:8000/api/extras/tags/2/",
            "display": "Source: proxmox",
            "name": "Source: proxmox",
            "slug": "source-proxmox",
            "color": "9e9e9e",
            "description": "Automatically created tag by netbox-ssot for source proxmox",
            "object_types": [],
            "tagged_items": 61,
            "created": "2024-05-09T14:07:10.376096Z",
            "last_updated": "2024-05-09T14:07:10.376108Z"
        },
        {
            "id": 1,
            "url": "http://ubtest.home:8000/api/extras/tags/1/",
            "display": "netbox-ssot",
            "name": "netbox-ssot",
            "slug": "netbox-ssot",
            "color": "00add8",
            "description": "Tag used by netbox-ssot to mark devices that are managed by it",
            "object_types": [],
            "tagged_items": 78,
            "created": "2024-05-09T14:07:08.027579Z",
            "last_updated": "2024-05-09T14:07:08.027593Z"
        },
        {
            "id": 3,
            "url": "http://ubtest.home:8000/api/extras/tags/3/",
            "display": "proxmox",
            "name": "proxmox",
            "slug": "type-proxmox",
            "color": "9e9e9e",
            "description": "Automatically created tag by netbox-ssot for source type proxmox",
            "object_types": [],
            "tagged_items": 61,
            "created": "2024-05-09T14:07:10.423769Z",
            "last_updated": "2024-05-09T14:07:10.423782Z"
        }
    ]
    }`
	MockTagPatchResponse = `{
    "id": 2,
    "url": "http://ubtest.home:8000/api/extras/tags/2/",
    "display": "Source: proxmox",
    "name": "Source: proxmox",
    "slug": "source-proxmox",
    "color": "9e9e9e",
    "description": "patched description",
    "object_types": [],
    "tagged_items": 61,
    "created": "2024-05-09T14:07:10.376096Z",
    "last_updated": "2024-05-10T10:39:22.750217Z"
}`
	MockTagCreateResponse = `{
    "id": 4,
    "url": "http://ubtest.home:8000/api/extras/tags/4/",
    "display": "newtag",
    "name": "newtag",
    "slug": "newtag",
    "color": "9e9e9e",
    "description": "",
    "object_types": [],
    "created": "2024-05-10T10:40:06.272285Z",
    "last_updated": "2024-05-10T10:40:06.272298Z"
  }`
)

const MockStatusResponse = `{
    "django-version": "5.0.5",
    "installed-apps": {
        "debug_toolbar": "4.3.0",
        "django_filters": "24.2",
        "django_prometheus": "2.3.1",
        "django_rq": "2.10.2",
        "django_tables2": "2.7.0",
        "drf_spectacular": "0.27.2",
        "drf_spectacular_sidecar": "2024.5.1",
        "mptt": "0.16.0",
        "rest_framework": "3.15.1",
        "social_django": "5.4.1",
        "taggit": "5.0.1",
        "timezone_field": "6.1.0"
    },
    "netbox-version": "4.0.0",
    "plugins": {},
    "python-version": "3.11.6",
    "rq-workers-running": 1
}`

func CreateMockServer() *httptest.Server {
	handler := http.NewServeMux()
	// Define handler for a specific path e.g., "/api/path"
	handler.HandleFunc("/api/status/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, MockStatusResponse)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	handler.HandleFunc(constants.TagsAPIPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			_, err := io.WriteString(w, MockTagPatchResponse)
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, err := io.WriteString(w, MockTagGetResponse)
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, err := io.WriteString(w, MockTagCreateResponse)
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			log.Printf("Wrong http method: %v", r.Method)
		}
	})

	// dummy endpoint for producing api error (http.StatusInternalServerError)
	handler.HandleFunc("/api/read-error", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
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
