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

// Hardcoded api return values.
const (
	VersionResponse   = "{\"django-version\": \"4.2.10\"}"
	TagsResponse      = "{\"count\":2,\"next\":null,\"previous\":null,\"results\":[{\"id\":0,\"url\":\"http://netbox.example.com/api/extras/tags/34/\",\"display\":\"Source: proxmox\",\"name\":\"Source: proxmox\",\"slug\":\"source-proxmox\",\"color\":\"9e9e9e\",\"description\":\"Automatically created tag by netbox-ssot for source proxmox\",\"object_types\":[],\"tagged_items\":115,\"created\":\"2024-02-25T16:52:51.691729Z\",\"last_updated\":\"2024-02-25T16:52:51.691743Z\"},{\"id\":1,\"url\":\"http://netbox.example.com/api/extras/tags/1/\",\"display\":\"netbox-ssot\",\"name\":\"netbox-ssot\",\"slug\":\"netbox-ssot\",\"color\":\"00add8\",\"description\":\"Tag used by netbox-ssot to mark devices that are managed by it\",\"object_types\":[],\"tagged_items\":134,\"created\":\"2024-02-11T16:23:17.082244Z\",\"last_updated\":\"2024-02-11T16:23:17.082257Z\"}]}"
	TagPatchResponse  = "{\"id\":1,\"name\":\"netbox-ssot\",\"slug\":\"netbox-ssot\",\"color\":\"00add8\",\"description\":\"patched description\"}"
	TagCreateResponse = "{\"id\":1,\"name\":\"netbox-ssot\",\"slug\":\"netbox-ssot\",\"color\":\"00add8\",\"description\":\"patched description\"}"
)

func CreateMockServer() *httptest.Server {
	handler := http.NewServeMux()
	// Define handler for a specific path e.g., "/api/path"
	handler.HandleFunc("/api/status/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, VersionResponse) // Mock JSON Response
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	handler.HandleFunc(constants.TagsAPIPath, func(w http.ResponseWriter, r *http.Request) {
		tagsIndex := map[int]objects.Tag{}
		var tags []objects.Tag
		err := json.Unmarshal([]byte(TagsResponse), &tags)
		if err != nil {
			log.Printf("error unmarshalling tags response: %s", err)
		}
		for _, tag := range tags {
			tagsIndex[tag.ID] = tag
		}
		switch r.Method {
		case http.MethodPatch:
			// TODO: add logic
			_, err = io.WriteString(w, TagPatchResponse)
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, err := io.WriteString(w, TagsResponse)
			if err != nil {
				log.Printf("Error writing response")
			}
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, err := io.WriteString(w, TagCreateResponse)
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
