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

	// User defined relation
	HostTenantRelations map[string]string
	HostSiteRelations   map[string]string
	VlanGroupRelations  map[string]string
	VlanTenantRelations map[string]string
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

	// Initialize regex relations for the sourcce
	// Initialize regex relations for this source
	fs.VlanGroupRelations = utils.ConvertStringsToRegexPairs(fs.SourceConfig.VlanGroupRelations)
	fs.Logger.Debugf(fs.Ctx, "VlanGroupRelations: %s", fs.VlanGroupRelations)
	fs.VlanTenantRelations = utils.ConvertStringsToRegexPairs(fs.SourceConfig.VlanTenantRelations)
	fs.Logger.Debugf(fs.Ctx, "VlanTenantRelations: %s", fs.VlanTenantRelations)
	fs.HostTenantRelations = utils.ConvertStringsToRegexPairs(fs.SourceConfig.HostTenantRelations)
	fs.Logger.Debugf(fs.Ctx, "HostTenantRelations: %s", fs.HostTenantRelations)
	fs.HostSiteRelations = utils.ConvertStringsToRegexPairs(fs.SourceConfig.HostSiteRelations)
	fs.Logger.Debugf(fs.Ctx, "HostSiteRelations: %s", fs.HostSiteRelations)

	initFunctions := []func(context.Context, *FortiClient) error{
		fs.InitSystemInfo,
		fs.InitInterfaces,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(ctx, c); err != nil {
			return fmt.Errorf("fortigate initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		fs.Logger.Infof(fs.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionName(initFunc), duration.Seconds())
	}
	return nil
}

func (fs *FortigateSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		fs.syncDevice,
		fs.SyncInterfaces,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		fs.Logger.Infof(fs.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionName(syncFunc), duration.Seconds())
	}
	return nil
}
