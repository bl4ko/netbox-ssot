package fmc

import (
	"fmt"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/source/fmc/client"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Source representing Cisco Firewall Management Center.
//
//nolint:revive
type FMCSource struct {
	common.Config

	// FMC data. Initialized in init functions.
	Domains              map[string]client.Domain
	Devices              map[string]*client.DeviceInfo
	DevicePhysicalIfaces map[string][]*client.PhysicalInterfaceInfo
	DeviceVlanIfaces     map[string][]*client.VLANInterfaceInfo

	// Netbox devices representing firewalls.
	NBDevices map[string]*objects.Device
}

func (fmcs *FMCSource) Init() error {
	httpClient, err := utils.NewHTTPClient(fmcs.SourceConfig.ValidateCert, fmcs.CAFile)
	if err != nil {
		return fmt.Errorf("create new http client: %s", err)
	}

	c, err := client.NewFMCClient(
		fmcs.Ctx,
		fmcs.SourceConfig.Username,
		fmcs.SourceConfig.Password,
		string(fmcs.SourceConfig.HTTPScheme),
		fmcs.SourceConfig.Hostname,
		fmcs.SourceConfig.Port,
		httpClient,
		fmcs.Logger,
	)
	if err != nil {
		return fmt.Errorf("create FMC client: %s", err)
	}

	initFunctions := []func(*client.FMCClient) error{
		fmcs.initObjects,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(c); err != nil {
			return fmt.Errorf("fmc initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		fmcs.Logger.Infof(
			fmcs.Ctx,
			"Successfully initialized %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"),
			duration.Seconds(),
		)
	}
	return nil
}

func (fmcs *FMCSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		fmcs.syncDevices,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		fmcs.Logger.Infof(
			fmcs.Ctx,
			"Successfully synced %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"),
			duration.Seconds(),
		)
	}
	return nil
}
