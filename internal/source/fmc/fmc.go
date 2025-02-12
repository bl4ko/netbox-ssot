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

// FMCSource represents Cisco Firewall Management Center.
//
//nolint:revive
type FMCSource struct {
	common.Config

	// FMC data. Initialized in init functions.
	// Domains is a map of domain UUIDs to Domain objects.
	Domains map[string]client.Domain
	// Devices is a map of device IDs to DeviceInfo objects.
	Devices map[string]*client.DeviceInfo
	// DevicePhysicalIfaces is a map of device IDs to a slice of PhysicalInterfaceInfo objects.
	DevicePhysicalIfaces map[string][]*client.PhysicalInterfaceInfo
	// DeviceVlanIfaces is a map of device IDs to a slice of VLANInterfaceInfo objects.
	DeviceVlanIfaces map[string][]*client.VLANInterfaceInfo
	// DeviceEtherChannelIfaces is a map of device IDs to a slice of EtherChannelInterfaceInfo objects.
	DeviceEtherChannelIfaces map[string][]*client.EtherChannelInterfaceInfo
	// DeviceSubIfaces is a map of device IDs to a slice of SubInterfaceInfo objects.
	DeviceSubIfaces map[string][]*client.SubInterfaceInfo

	// Netbox devices representing firewalls.
	NBDevices map[string]*objects.Device
	// NBInterfaces represents all fmc interfaces that have been synced to netbox.
	// It is a map of interface name to interface, so we can find parents of sub interfaces.
	Name2NBInterface map[string]*objects.Interface
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

	fmcs.Name2NBInterface = make(map[string]*objects.Interface)

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
