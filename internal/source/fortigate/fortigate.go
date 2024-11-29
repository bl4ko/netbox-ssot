package fortigate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

//nolint:revive
type FortigateSource struct {
	common.Config
	// Fortinet data. Initialized in init functions.
	SystemInfo FortiSystemInfo              // Map storing system information
	Ifaces     map[string]InterfaceResponse // iface name -> FortigateInterface

	// NBFirewall representing fortinet firewall created in syncDevice func.
	NBFirewall *objects.Device
}

type FortiSystemInfo struct {
	Hostname string
	Version  string
	Serial   string
}

type FortiClient struct {
	HTTPClient *http.Client
	BaseURL    string
	APIToken   string
}

func NewAPIClient(apiToken string, baseURL string, httpClient *http.Client) *FortiClient {
	return &FortiClient{
		HTTPClient: httpClient,
		BaseURL:    baseURL,
		APIToken:   apiToken,
	}
}

func (c FortiClient) MakeRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.BaseURL, path), body)
	if err != nil {
		return nil, err
	}
	// Set the Authorization header.
	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	return c.HTTPClient.Do(req)
}

func (fs *FortigateSource) Init() error {
	httpClient, err := utils.NewHTTPClient(fs.SourceConfig.ValidateCert, fs.CAFile)
	if err != nil {
		return fmt.Errorf("create new http client: %s", err)
	}
	c := NewAPIClient(fs.SourceConfig.APIToken, fmt.Sprintf("%s://%s:%d/api/v2", fs.SourceConfig.HTTPScheme, fs.SourceConfig.Hostname, fs.SourceConfig.Port), httpClient)
	ctx := context.Background()
	defer ctx.Done()

	initFunctions := []func(context.Context, *FortiClient) error{
		fs.initSystemInfo,
		fs.initInterfaces,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(ctx, c); err != nil {
			return fmt.Errorf("fortigate initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		fs.Logger.Infof(fs.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"), duration.Seconds())
	}
	return nil
}

func (fs *FortigateSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		fs.syncDevice,
		fs.syncInterfaces,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		fs.Logger.Infof(fs.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"), duration.Seconds())
	}
	return nil
}
