package service

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
)

func TestNewNetBoxAPI(t *testing.T) {
	type args struct {
		logger       *logger.Logger
		baseURL      string
		apiToken     string
		validateCert bool
		timeout      int
		caCert       string
	}
	tests := []struct {
		name string
		args args
		want *NetboxClient
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
				Logger:   &logger.Logger{Logger: log.Default()},
				BaseURL:  "netbox.example.com",
				APIToken: "apitoken",
				HTTPClient: &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				},
				Timeout: constants.DefaultAPITimeout,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNetboxClient(
				tt.args.logger,
				tt.args.baseURL,
				tt.args.apiToken,
				tt.args.validateCert,
				tt.args.timeout,
				tt.args.caCert,
			)
			if err != nil {
				t.Errorf("NewNetboxClient() error = %v", err)
				return
			}
			// Check non-pointer fields for simplicity or use an interface to mock clients
			if got.BaseURL != tt.want.BaseURL || got.APIToken != tt.want.APIToken ||
				got.Timeout != tt.want.Timeout {
				t.Errorf("NewNetboxClient() got = %v, want %v", got, tt.want)
			}
			// Optionally check if HTTPClient is not nil to confirm it's initialized
			if got.HTTPClient == nil {
				t.Errorf("HTTPClient was not initialized")
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
				method: http.MethodGet,
				path:   "/api/status/",
				body:   nil,
			},
			want:    &APIResponse{StatusCode: http.StatusOK, Body: []byte(MockVersionResponseJSON)},
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
