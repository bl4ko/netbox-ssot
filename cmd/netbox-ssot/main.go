package main

import (
	"fmt"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source"
)

func main() {
	starttime := time.Now()

	// Parse configuration
	fmt.Printf("Netbox-SSOT has started at %s\n", starttime.Format(time.RFC3339))

	config, err := parser.ParseConfig("config.yaml")
	if err != nil {
		fmt.Println("Parser:", err)
		return
	}
	// Initialise Logger
	logger, err := logger.New(config.Logger.Dest, config.Logger.Level)
	if err != nil {
		fmt.Println("Logger:", err)
		return
	}
	logger.Debug("Parsed Logger config: ", config.Logger)
	logger.Debug("Parsed Netbox config: ", config.Netbox)
	logger.Debug("Parsed Source config: ", config.Sources)

	netboxInventory := inventory.NewNetboxInventory(logger, config.Netbox)
	logger.Debug("Netbox inventory: ", netboxInventory)
	err = netboxInventory.Init()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug("Netbox inventory initialized: ", netboxInventory)

	// Go through all sources and sync data
	for _, sourceConfig := range config.Sources {
		logger.Info("Processing source ", sourceConfig.Name, "...")

		// Source initialization
		logger.Info("Creating new source...")
		source, err := source.NewSource(&sourceConfig, logger, netboxInventory)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Info("Source initialized successfully: ", source)
		err = source.Init()
		if err != nil {
			logger.Error(err)
			return
		}

		// Source synchronization
		logger.Info("Syncing source...")
		err = source.Sync(netboxInventory)
		if err != nil {
			logger.Error(err)
			return
		}

		logger.Info("Source ", sourceConfig.Name, " successfully")
	}
}
