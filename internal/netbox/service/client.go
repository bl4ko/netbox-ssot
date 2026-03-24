package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// NetboxClient is a service used for communicating with the Netbox API.
// It is created via constructor func newNetboxAPI().
type NetboxClient struct {
	Logger     *logger.Logger
	HTTPClient *http.Client
	BaseURL    string
	APIToken   string
	Timeout    int // in seconds
	MaxRetires int
	DryRun     bool

	nextFakeID     int64
	nextFakeIDLock sync.Mutex
}

const dryRunFakeIDStart = 100_000_000

// APIResponse is a struct that represents a response from the Netbox API.
type APIResponse struct {
	StatusCode int
	Body       []byte
}

// Constructor function for creating a new netBoxAPI instance.
func NewNetboxClient(
	logger *logger.Logger,
	baseURL string,
	apiToken string,
	validateCert bool,
	timeout int,
	caCert string,
	dryRun bool,
) (*NetboxClient, error) {
	httpClient, err := utils.NewHTTPClient(validateCert, caCert)
	if err != nil {
		return nil, fmt.Errorf("create new HTTP client: %s", err)
	}
	return &NetboxClient{
		HTTPClient: httpClient,
		Logger:     logger,
		BaseURL:    baseURL,
		APIToken:   apiToken,
		Timeout:    timeout,
		DryRun:     dryRun,
		nextFakeID: dryRunFakeIDStart,
	}, nil
}

func (api *NetboxClient) generateFakeID() int {
	api.nextFakeIDLock.Lock()
	defer api.nextFakeIDLock.Unlock()
	id := api.nextFakeID
	api.nextFakeID++
	return int(id)
}

func (api *NetboxClient) doRequest(
	method string,
	path string,
	body io.Reader,
) (*APIResponse, error) {
	ctx, cancelCtx := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(api.Timeout),
	)
	defer cancelCtx()

	req, err := http.NewRequestWithContext(ctx, method, api.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	// We add necessary headers to the request
	req.Header.Add("Authorization", "Token "+api.APIToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := api.HTTPClient.Do(req)
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
