package dnac

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	dnac "github.com/cisco-en-programmability/dnacenter-go-sdk/v5/sdk"
)

type DnacSource struct {
	common.CommonConfig

	Sites map[string]dnac.ResponseSitesGetSiteResponse // SiteId -> Site
}

func (ds *DnacSource) Init() error {
	dnacUrl := fmt.Sprintf("%s://%s:%d", ds.CommonConfig.SourceConfig.HTTPScheme, ds.CommonConfig.SourceConfig.Hostname, ds.CommonConfig.SourceConfig.Port)
	Client, err := dnac.NewClientWithOptions(dnacUrl, ds.SourceConfig.Username, ds.SourceConfig.Password, "false", strconv.FormatBool(ds.SourceConfig.ValidateCert), nil)
	if err != nil {
		return fmt.Errorf("creating dnac client: %s", err)
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize items from vsphere API to local storage
	initFunctions := []func(context.Context, *dnac.Client) error{
		ds.InitSites,
		ds.InitDevices,
	}

	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(ctx, Client); err != nil {
			return fmt.Errorf("dnac initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		ds.Logger.Infof("Successfully initialized %s in %f seconds", utils.ExtractFunctionName(initFunc), duration.Seconds())
	}

	return nil
}

func (ds *DnacSource) Sync(nbi *inventory.NetBoxInventory) error {
	return nil
}
