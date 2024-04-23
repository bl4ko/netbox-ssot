package fmc

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Source representing Cisco Firewall Management Center.
//
//nolint:revive
type FMCSource struct {
	common.Config

	// FMC data. Initialized in init functions.
	Domains              map[string]*Domain
	Devices              map[string]*DeviceInfo
	DevicePhysicalIfaces map[string][]*PhysicalInterfaceInfo
	DeviceVlanIfaces     map[string][]*VLANInterfaceInfo

	// Netbox devices representing firewalls.
	NBDevices map[string]*objects.Device

	// User defined relation
	HostTenantRelations map[string]string
	HostSiteRelations   map[string]string
	VlanGroupRelations  map[string]string
	VlanTenantRelations map[string]string
}

func (fmcs *FMCSource) Init() error {
	HTTPClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !fmcs.SourceConfig.ValidateCert,
			},
		},
	}

	c, err := newFMCClient(fmcs.SourceConfig.Username, fmcs.SourceConfig.Password, string(fmcs.SourceConfig.HTTPScheme), fmcs.SourceConfig.Hostname, fmcs.SourceConfig.Port, HTTPClient)
	if err != nil {
		return fmt.Errorf("create FMC client: %s", err)
	}

	// Initialize regex relations for this source
	fmcs.VlanGroupRelations = utils.ConvertStringsToRegexPairs(fmcs.SourceConfig.VlanGroupRelations)
	fmcs.Logger.Debugf(fmcs.Ctx, "VlanGroupRelations: %s", fmcs.VlanGroupRelations)
	fmcs.VlanTenantRelations = utils.ConvertStringsToRegexPairs(fmcs.SourceConfig.VlanTenantRelations)
	fmcs.Logger.Debugf(fmcs.Ctx, "VlanTenantRelations: %s", fmcs.VlanTenantRelations)
	fmcs.HostTenantRelations = utils.ConvertStringsToRegexPairs(fmcs.SourceConfig.HostTenantRelations)
	fmcs.Logger.Debugf(fmcs.Ctx, "HostTenantRelations: %s", fmcs.HostTenantRelations)
	fmcs.HostSiteRelations = utils.ConvertStringsToRegexPairs(fmcs.SourceConfig.HostSiteRelations)
	fmcs.Logger.Debugf(fmcs.Ctx, "HostSiteRelations: %s", fmcs.HostSiteRelations)

	initFunctions := []func(*fmcClient) error{
		fmcs.initDevices,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(c); err != nil {
			return fmt.Errorf("fmc initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		fmcs.Logger.Infof(fmcs.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionName(initFunc), duration.Seconds())
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
		fmcs.Logger.Infof(fmcs.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionName(syncFunc), duration.Seconds())
	}
	return nil
}
