package paloalto

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/eth"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"
	"github.com/PaloAltoNetworks/pango/netw/zone"
	"github.com/PaloAltoNetworks/pango/vsys"
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
	SystemInfo          map[string]string         // Map storing system information
	VirtualSystems      map[string]vsys.Entry     // VirtualSystem name -> VirtualSystem
	SecurityZones       map[string]zone.Entry     // SecurityZone name -> SecurityZone
	Iface2SecurityZone  map[string]string         // Iface name -> SecurityZone name
	Iface2VirtualRouter map[string]string         // Iface name -> VirtualRouter name
	Ifaces              map[string]eth.Entry      // Iface name -> Iface
	Iface2SubIfaces     map[string][]layer3.Entry // Iface name -> SubIfaces
	VirtualRouters      map[string]router.Entry   // VirtualRouter name -> VirutalRouter
	ArpData             []ArpEntry                // Array of arp entreies

	// NBFirewall representing paloalto firewall created in syncDevice func.
	NBFirewall *objects.Device
}

func (pas *PaloAltoSource) Init() error {
	var transport *http.Transport
	var err error
	if pas.Config.CAFile != "" {
		transport, err = utils.LoadExtraCertInTransportConfig(pas.Config.CAFile)
		if err != nil {
			return fmt.Errorf("load extra cert in transport config: %s", err)
		}
	}
	c := &pango.Firewall{Client: pango.Client{
		Hostname:          pas.SourceConfig.Hostname,
		Username:          pas.SourceConfig.Username,
		Password:          pas.SourceConfig.Password,
		Logging:           pango.LogAction | pango.LogOp,
		VerifyCertificate: pas.SourceConfig.ValidateCert,
		Port:              uint(pas.SourceConfig.Port), //nolint:gosec
		Timeout:           constants.DefaultAPITimeout,
		Protocol:          string(pas.SourceConfig.HTTPScheme),
		Transport:         transport,
	}}

	if err := c.Initialize(); err != nil {
		return fmt.Errorf("paloalto failed to initialize client: %s", err)
	}

	initFunctions := []func(*pango.Firewall) error{
		pas.initArpData,
		pas.initSystemInfo,
		pas.initVirtualSystems,
		pas.initInterfaces,
		pas.initVirtualRouters,
	}
	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(c); err != nil {
			return fmt.Errorf("paloalto initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		pas.Logger.Infof(
			pas.Ctx,
			"Successfully initialized %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"),
			duration.Seconds(),
		)
	}
	return nil
}

func (pas *PaloAltoSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		pas.syncDevice,
		pas.syncSecurityZones,
		pas.syncInterfaces,
		pas.syncArpTable,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		pas.Logger.Infof(
			pas.Ctx,
			"Successfully synced %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"),
			duration.Seconds(),
		)
	}
	return nil
}
