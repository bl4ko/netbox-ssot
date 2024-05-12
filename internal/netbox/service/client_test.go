package service

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	servicetest "github.com/bl4ko/netbox-ssot/internal/netbox/service/testing"
)

func TestNewNetboxClient(t *testing.T) {
	type args struct {
		logger       *logger.Logger
		baseURL      string
		apiToken     string
		validateCert bool
		timeout      int
		caCert       string
	}
	tests := []struct {
		name    string
		args    args
		want    *NetboxClient
		wantErr bool
	}{
		{
			name: "test new API creation without ssl verify",
			args: args{
				logger:       &logger.Logger{Logger: log.Default()},
				baseURL:      "netbox.example.com",
				apiToken:     "apitoken",
				validateCert: false,
				timeout:      constants.DefaultAPITimeout,
				caCert:       "",
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
				logger:       &logger.Logger{Logger: log.Default()},
				baseURL:      "netbox.example.com",
				apiToken:     "apitoken",
				validateCert: true,
				timeout:      constants.DefaultAPITimeout,
				caCert:       "",
			},
			want: &NetboxClient{
				Logger:   &logger.Logger{Logger: log.Default()},
				BaseURL:  "netbox.example.com",
				APIToken: "apitoken",
				HTTPClient: &http.Client{Transport: &http.Transport{
					TLSClientConfig: &tls.Config{},
				}},
				Timeout: constants.DefaultAPITimeout,
			},
		},
		{
			name: "test newNetboxClient creation error with worng caCert path",
			args: args{
				logger:       &logger.Logger{Logger: log.Default()},
				baseURL:      "netbox.example.com",
				apiToken:     "apitoken",
				validateCert: true,
				timeout:      constants.DefaultAPITimeout,
				caCert:       "wrong path",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNetboxClient(tt.args.logger, tt.args.baseURL, tt.args.apiToken, tt.args.validateCert, tt.args.timeout, tt.args.caCert)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNetboxClient() error = %v", err)
				return
			}
			if tt.wantErr == false {
				if got.BaseURL != tt.want.BaseURL || got.APIToken != tt.want.APIToken || got.Timeout != tt.want.Timeout {
					t.Errorf("NewNetboxClient() got = %v, want %v", got, tt.want)
				}
				if got.HTTPClient == nil {
					t.Errorf("HTTPClient was not initialized")
				}
			}
		})
	}
}

func TestNetboxAPI_doRequest(t *testing.T) {
	// Setup common elements
	netboxMock := servicetest.CreateMockServer()
	defer netboxMock.Close()

	logger := &logger.Logger{Logger: log.Default()}
	baseURL := netboxMock.URL
	apiToken := "testtoken"

	// Common client for all sub-tests
	netboxClient, err := NewNetboxClient(logger, baseURL, apiToken, false, 30, "")
	if err != nil {
		t.Fatalf("Failed to create netbox client: %s", err)
	}

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "Valid GET request to /api/status",
			method:     MethodGet,
			path:       "/api/status",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:    "Invalid method type",
			method:  "ž",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := netboxClient.doRequest(tt.method, tt.path, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("doRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && response.StatusCode != tt.wantStatus {
				t.Errorf("Wrong status code: got %v, want %v", response.StatusCode, tt.wantStatus)
			}
		})
	}

	// Test for network failure simulation
	t.Run("Network failure simulation", func(t *testing.T) {
		failingNetboxClient := &NetboxClient{
			HTTPClient: &http.Client{Transport: &failingHTTPClientNetwork{}},
			Logger:     logger,
			BaseURL:    baseURL,
			APIToken:   apiToken,
			Timeout:    30,
		}

		_, err := failingNetboxClient.doRequest(MethodGet, "/", nil)
		if err == nil {
			t.Error("Expected error but got none")
		}
	})

	t.Run("Body read failure simulation", func(t *testing.T) {
		failingNetboxClient := &NetboxClient{
			HTTPClient: &http.Client{Transport: &failingHTTPClientRead{}},
			Logger:     logger,
			BaseURL:    baseURL,
			APIToken:   apiToken,
			Timeout:    30,
		}

		_, err := failingNetboxClient.doRequest(MethodGet, "/", nil)
		if err == nil {
			t.Error("Expected error but got none")
		}
	})
}

// Mock for failing HTTP client for triggering network failure
type failingHTTPClientNetwork struct{}

func (f *failingHTTPClientNetwork) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("network failure")
}

// Mock for failing read body HTTP client for triggering failure on body read
type failingHTTPClientRead struct{}

func (m *failingHTTPClientRead) RoundTrip(_ *http.Request) (*http.Response, error) {
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

func TestNetboxClient_doRequest(t *testing.T) {
	type args struct {
		method string
		path   string
		body   io.Reader
	}
	tests := []struct {
		name    string
		api     *NetboxClient
		args    args
		want    *APIResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.api.doRequest(tt.args.method, tt.args.path, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxClient.doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxClient.doRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
