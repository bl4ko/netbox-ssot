package service

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/logger"
)

// NetboxAPI is a service used for communicating with the Netbox API.
// It is created via constructur func newNetboxAPI()
type NetboxAPI struct {
	Logger     *logger.Logger
	HttpClient *http.Client
	BaseURL    string
	ApiKey     string
}

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
	MethodPatch  = "PATCH"
)

// APIResponse is a struct that represents a response from the Netbox API.
type APIResponse struct {
	StatusCode int
	Body       []byte
}

// Constructor function for creating a new netBoxAPI instance.
func NewNetBoxAPI(logger *logger.Logger, baseURL string, apiToken string, validateCert bool) *NetboxAPI {
	var client *http.Client
	if validateCert {
		client = &http.Client{}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	return &NetboxAPI{
		HttpClient: client,
		Logger:     logger,
		BaseURL:    baseURL,
		ApiKey:     apiToken,
	}
}

func (api *NetboxAPI) doRequest(method string, path string, body io.Reader) (*APIResponse, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelCtx()

	req, err := http.NewRequestWithContext(ctx, method, api.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	// We add necessary headers to the request
	req.Header.Add("Authorization", "Token "+api.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := api.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &APIResponse{
		StatusCode: resp.StatusCode,
		Body:       responseBody,
	}, nil
}
