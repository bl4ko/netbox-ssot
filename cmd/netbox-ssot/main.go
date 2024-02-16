package main

import (
	"fmt"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source"
)

func main() {
	startTime := time.Now()

	// Parse configuration
	fmt.Printf("Netbox-SSOT has started at %s\n", startTime.Format(time.RFC3339))

	config, err := parser.ParseConfig("config.yaml")
	if err != nil {
		fmt.Println("Parser:", err)
		return
	}
	// Initialize Logger
	mainLogger, err := logger.New(config.Logger.Dest, config.Logger.Level, "main")
	if err != nil {
		fmt.Println("Logger:", err)
		return
	}
	mainLogger.Debug("Parsed Logger config: ", config.Logger)
	mainLogger.Debug("Parsed Netbox config: ", config.Netbox)
	mainLogger.Debug("Parsed Source config: ", config.Sources)

	inventoryLogger, err := logger.New(config.Logger.Dest, config.Logger.Level, "netboxInventory")
	if err != nil {
		mainLogger.Errorf("inventoryLogger: %s", err)
	}
	netboxInventory := inventory.NewNetboxInventory(inventoryLogger, config.Netbox)
	mainLogger.Debug("Netbox inventory: ", netboxInventory)

	mainLogger.Info("Starting initializing netbox inventory")
	err = netboxInventory.Init()
	if err != nil {
		mainLogger.Error(err)
		return
	}
	mainLogger.Debug("Netbox inventory initialized: ", netboxInventory)

	// Go through all sources and sync data
	for i := range config.Sources {
		sourceConfig := &config.Sources[i]
		mainLogger.Info("Processing source ", sourceConfig.Name, "...")

		sourceLogger, err := logger.New(config.Logger.Dest, config.Logger.Level, sourceConfig.Name)
		if err != nil {
			mainLogger.Errorf("source logger: %s", err)
		}
		source, err := source.NewSource(sourceConfig, sourceLogger, netboxInventory)
		if err != nil {
			sourceLogger.Error(err)
			return
		}
		mainLogger.Infof("Successfully created source %s", constants.CheckMark)
		sourceLogger.Debugf("Source content: %s", source)

		sourceLogger.Info("Initializing source")
		err = source.Init()
		if err != nil {
			sourceLogger.Error(err)
			return
		}
		sourceLogger.Infof("Successfully initialized source %s", constants.CheckMark)

		// Source synchronization
		sourceLogger.Info("Syncing source...")
		err = source.Sync(netboxInventory)
		if err != nil {
			sourceLogger.Error(err)
			return
		}
		sourceLogger.Infof("Source synced successfully %s", constants.CheckMark)
	}

	// Orphan manager cleanup
	if config.Netbox.RemoveOrphans {
		mainLogger.Info("Cleaning up orphaned objects...")
		err = netboxInventory.DeleteOrphans()
		if err != nil {
			mainLogger.Error(err)
			return
		}
		mainLogger.Infof("%s Successfully removed orphans", constants.CheckMark)
	} else {
		mainLogger.Info("Skipping removing orphaned objects...")
	}

	duration := time.Since(startTime)
	minutes := int(duration.Minutes())
	seconds := int((duration - time.Duration(minutes)*time.Minute).Seconds())
	mainLogger.Infof("%s Syncing took %d min %d sec in total", constants.Rocket, minutes, seconds)
}
