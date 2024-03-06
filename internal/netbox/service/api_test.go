package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
)

// Hardcoded api return values.
const (
	VersionResponse = "{\"django-version\": \"4.2.10\"}"
	TagsResponse    = "{\"count\":2,\"next\":null,\"previous\":null,\"results\":[{\"id\":34,\"url\":\"http://netbox.example.com/api/extras/tags/34/\",\"display\":\"Source: proxmox\",\"name\":\"Source: proxmox\",\"slug\":\"source-proxmox\",\"color\":\"9e9e9e\",\"description\":\"Automatically created tag by netbox-ssot for source proxmox\",\"object_types\":[],\"tagged_items\":115,\"created\":\"2024-02-25T16:52:51.691729Z\",\"last_updated\":\"2024-02-25T16:52:51.691743Z\"},{\"id\":1,\"url\":\"http://netbox.example.com/api/extras/tags/1/\",\"display\":\"netbox-ssot\",\"name\":\"netbox-ssot\",\"slug\":\"netbox-ssot\",\"color\":\"00add8\",\"description\":\"Tag used by netbox-ssot to mark devices that are managed by it\",\"object_types\":[],\"tagged_items\":134,\"created\":\"2024-02-11T16:23:17.082244Z\",\"last_updated\":\"2024-02-11T16:23:17.082257Z\"}]}"
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

	handler.HandleFunc(TagsAPIPath, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, TagsResponse)
		if err != nil {
			log.Printf("Error writing response")
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

func TestNewNetBoxAPI(t *testing.T) {
	type args struct {
		ctx          context.Context
		logger       *logger.Logger
		baseURL      string
		apiToken     string
		validateCert bool
		timeout      int
	}
	tests := []struct {
		name string
		args args
		want *NetboxClient
	}{
		{
			name: "test new API creation without ssl verify",
			args: args{
				ctx:          context.Background(),
				logger:       &logger.Logger{Logger: log.Default()},
				baseURL:      "netbox.example.com",
				apiToken:     "apitoken",
				validateCert: false,
				timeout:      constants.DefaultAPITimeout,
			},
			want: &NetboxClient{
				Logger:     &logger.Logger{Logger: log.Default()},
				BaseURL:    "netbox.example.com",
				APIToken:   "apitoken",
				HTTPClient: &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
				Timeout:    constants.DefaultAPITimeout,
			},
		},
		{
			name: "test new API creation with ssl verify",
			args: args{
				ctx:          context.Background(),
				logger:       &logger.Logger{Logger: log.Default()},
				baseURL:      "netbox.example.com",
				apiToken:     "apitoken",
				validateCert: true,
				timeout:      constants.DefaultAPITimeout,
			},
			want: &NetboxClient{
				Logger:     &logger.Logger{Logger: log.Default()},
				BaseURL:    "netbox.example.com",
				APIToken:   "apitoken",
				HTTPClient: &http.Client{},
				Timeout:    constants.DefaultAPITimeout,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNetboxClient(tt.args.ctx, tt.args.logger, tt.args.baseURL, tt.args.apiToken, tt.args.validateCert, tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNetBoxAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxAPI_doRequest(t *testing.T) {
	type args struct {
		method string
		path   string
		body   io.Reader
	}
	tests := []struct {
		name         string
		netboxClient *NetboxClient
		args         args
		want         *APIResponse
		wantErr      bool
	}{
		{
			name:         "Test GET /api/status/",
			netboxClient: MockNetboxClient,
			args: args{
				method: MethodGet,
				path:   "/api/status/",
				body:   nil,
			},
			want:    &APIResponse{StatusCode: http.StatusOK, Body: []byte(VersionResponse)},
			wantErr: false,
		},
		{
			name:         "Test Invalid Request",
			netboxClient: MockNetboxClient,
			args: args{
				method: "\n", // Invalid method
				path:   "/api/status",
				body:   nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "Client failure",
			netboxClient: FailingMockNetboxClient,
			args: args{
				method: http.MethodGet,
				path:   "/api/status",
				body:   nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "Test ReadALL Error",
			netboxClient: MockNetboxClientWithReadError,
			args: args{
				method: http.MethodGet,
				path:   "/api/read-error",
				body:   nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	mockServer := CreateMockServer()
	defer mockServer.Close()
	MockNetboxClient.BaseURL = mockServer.URL
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.netboxClient.doRequest(tt.args.method, tt.args.path, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxAPI.doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxAPI.doRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
