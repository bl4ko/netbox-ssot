package paloalto

import (
	"fmt"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/eth"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"
	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

//nolint:revive
type PaloAltoSource struct {
	common.Config
	// Paloalto data. Initialized in init functions.
	Interface2Router        map[string]string         // Iface name -> Router name
	Interfaces              map[string]eth.Entry      // Iface name -> Iface
	Interface2Subinterfaces map[string][]layer3.Entry // Iface name -> Subinterfaces
	Devices                 map[string]router.Entry   // Router name -> router

	// Netbox related data for easier access. Initialized in sync functions.
	DeviceName2NbDevice map[string]*objects.Device

	// User defined relation
	HostTenantRelations map[string]string
	HostSiteRelations   map[string]string
	VlanGroupRelations  map[string]string
	VlanTenantRelations map[string]string
}

func (pas *PaloAltoSource) Init() error {
	c := &pango.Firewall{Client: pango.Client{
		Hostname:          pas.SourceConfig.Hostname,
		Username:          pas.SourceConfig.Username,
		Password:          pas.SourceConfig.Password,
		Logging:           pango.LogAction | pango.LogOp,
		VerifyCertificate: pas.SourceConfig.ValidateCert,
		Port:              uint(pas.SourceConfig.Port),
		Timeout:           constants.DefaultAPITimeout,
		Protocol:          string(pas.SourceConfig.HTTPScheme),
	}}

	if err := c.Initialize(); err != nil {
		return fmt.Errorf("paloalto failed to initialize client: %s", err)
	}

	// Initialize regex relations for the sourcce
	// Initialize regex relations for this source
	pas.VlanGroupRelations = utils.ConvertStringsToRegexPairs(pas.SourceConfig.VlanGroupRelations)
	pas.Logger.Debugf(pas.Ctx, "VlanGroupRelations: %s", pas.VlanGroupRelations)
	pas.VlanTenantRelations = utils.ConvertStringsToRegexPairs(pas.SourceConfig.VlanTenantRelations)
	pas.Logger.Debugf(pas.Ctx, "VlanTenantRelations: %s", pas.VlanTenantRelations)
	pas.HostTenantRelations = utils.ConvertStringsToRegexPairs(pas.SourceConfig.HostTenantRelations)
	pas.Logger.Debugf(pas.Ctx, "HostTenantRelations: %s", pas.HostTenantRelations)
	pas.HostSiteRelations = utils.ConvertStringsToRegexPairs(pas.SourceConfig.HostSiteRelations)
	pas.Logger.Debugf(pas.Ctx, "HostSiteRelations: %s", pas.HostSiteRelations)

	initFunctions := []func(*pango.Firewall) error{
		pas.InitInterfaces,
		pas.InitDevices,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(c); err != nil {
			return fmt.Errorf("paloalto initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		pas.Logger.Infof(pas.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionName(initFunc), duration.Seconds())
	}
	return nil
}

func (pas *PaloAltoSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		pas.SyncDevices,
		pas.SyncInterfaces,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		pas.Logger.Infof(pas.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionName(syncFunc), duration.Seconds())
	}
	return nil
}
