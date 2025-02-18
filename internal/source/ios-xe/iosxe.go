package iosxe

import (
	"fmt"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/options"
)

//nolint:revive
type IOSXESource struct {
	common.Config

	// IOSXE fetched data. Initialized in init functions.
	HardwareInfo hardwareReply
	SystemInfo   systemReply
	Interfaces   map[string]iface
	ArpEntries   []arpEntry

	// IOSXE synced data. Created in sync functions.
	NBDevice     *objects.Device
	NBInterfaces map[string]*objects.Interface // interfaceName -> netboxInterface
}

func (is *IOSXESource) Init() error {
	d, err := netconf.NewDriver(
		is.SourceConfig.Hostname,
		options.WithAuthUsername(is.SourceConfig.Username),
		options.WithAuthPassword(is.SourceConfig.Password),
		options.WithPort(is.SourceConfig.Port),
		options.WithAuthNoStrictKey(), // inside container we can't confirm ssh key
	)
	if err != nil {
		return fmt.Errorf("failed to create driver: %s", err)
	}
	err = d.Open()
	if err != nil {
		return fmt.Errorf("failed to open driver: %s", err)
	}
	defer d.Close()

	// Initialize items from vsphere API to local storage
	initFunctions := []func(*netconf.Driver) error{
		is.initDeviceInfo,
		is.initDeviceHardwareInfo,
		is.initInterfaces,
		is.initArpData,
	}

	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(d); err != nil {
			return fmt.Errorf("iosxe initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		is.Logger.Infof(
			is.Ctx,
			"Successfully initialized %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"),
			duration.Seconds(),
		)
	}
	return nil
}

func (is *IOSXESource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		is.syncDevice,
		is.syncInterfaces,
		is.syncArpTable,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		is.Logger.Infof(
			is.Ctx,
			"Successfully synced %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"),
			duration.Seconds(),
		)
	}
	return nil
}
