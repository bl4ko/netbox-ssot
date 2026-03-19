package f5

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

//nolint:revive
type F5Source struct {
	common.Config
	// F5 BIG-IP data. Initialized in init functions.
	VirtualServers []VirtualServerResponse
}

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Username   string
	Password   string
	APIToken   string
}

func NewAPIClient(
	username, password, apiToken, baseURL string,
	httpClient *http.Client,
) *Client {
	return &Client{
		HTTPClient: httpClient,
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		APIToken:   apiToken,
	}
}

func (c *Client) MakeRequest(
	ctx context.Context,
	method, path string,
	body io.Reader,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.BaseURL, path), body)
	if err != nil {
		return nil, err
	}
	if c.APIToken != "" {
		req.Header.Set("X-F5-Auth-Token", c.APIToken)
	} else {
		req.SetBasicAuth(c.Username, c.Password)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.HTTPClient.Do(req)
}

func (fs *F5Source) Init() error {
	httpClient, err := utils.NewHTTPClient(fs.SourceConfig.ValidateCert, fs.CAFile)
	if err != nil {
		return fmt.Errorf("create new http client: %s", err)
	}
	c := NewAPIClient(
		fs.SourceConfig.Username,
		fs.SourceConfig.Password,
		fs.SourceConfig.APIToken,
		fmt.Sprintf(
			"%s://%s:%d/mgmt/tm",
			fs.SourceConfig.HTTPScheme,
			fs.SourceConfig.Hostname,
			fs.SourceConfig.Port,
		),
		httpClient,
	)
	ctx := context.Background()
	defer ctx.Done()

	initFunctions := []func(context.Context, *Client) error{
		fs.initVirtualServers,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(ctx, c); err != nil {
			return fmt.Errorf("f5 initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		fs.Logger.Infof(
			fs.Ctx,
			"Successfully initialized %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"),
			duration.Seconds(),
		)
	}
	return nil
}

func (fs *F5Source) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		fs.syncVirtualServers,
	}

	var encounteredErrors []error
	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		funcName := utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync")
		err := syncFunc(nbi)
		if err != nil {
			if fs.SourceConfig.ContinueOnError {
				fs.Logger.Errorf(
					fs.Ctx,
					"Error syncing %s: %s (continuing due to continueOnError flag)",
					funcName,
					err,
				)
				encounteredErrors = append(encounteredErrors, fmt.Errorf("%s: %w", funcName, err))
			} else {
				return err
			}
		} else {
			duration := time.Since(startTime)
			fs.Logger.Infof(
				fs.Ctx,
				"Successfully synced %s in %f seconds",
				funcName,
				duration.Seconds(),
			)
		}
	}
	if len(encounteredErrors) > 0 {
		return fmt.Errorf("encountered %d errors during sync: %v", len(encounteredErrors), encounteredErrors)
	}
	return nil
}
