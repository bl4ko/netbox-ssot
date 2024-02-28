package main

import (
	"context"
	"fmt"
	"sync"
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

	// Create our main context
	mainCtx := context.Background()
	mainCtx = context.WithValue(mainCtx, constants.CtxSourceKey, "main")

	// Initialize Logger
	ssotLogger, err := logger.New(config.Logger.Dest, config.Logger.Level)
	if err != nil {
		fmt.Println("Logger:", err)
		return
	}
	ssotLogger.Debug(mainCtx, "Parsed Logger config: ", config.Logger)
	ssotLogger.Debug(mainCtx, "Parsed Netbox config: ", config.Netbox)
	ssotLogger.Debug(mainCtx, "Parsed Source config: ", config.Sources)

	inventoryLogger, err := logger.New(config.Logger.Dest, config.Logger.Level)
	if err != nil {
		ssotLogger.Errorf(mainCtx, "inventoryLogger: %s", err)
	}
	inventoryCtx := context.WithValue(context.Background(), constants.CtxSourceKey, "inventory")
	netboxInventory := inventory.NewNetboxInventory(inventoryCtx, inventoryLogger, config.Netbox)
	ssotLogger.Debug(mainCtx, "Netbox inventory: ", netboxInventory)

	ssotLogger.Info(mainCtx, "Starting initializing netbox inventory")
	err = netboxInventory.Init()
	if err != nil {
		ssotLogger.Error(mainCtx, err)
		return
	}
	ssotLogger.Debug(mainCtx, "Netbox inventory initialized: ", netboxInventory)

	// Go through all sources and sync data
	var wg sync.WaitGroup
	for i := range config.Sources {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sourceConfig := &config.Sources[i]
			ssotLogger.Info(mainCtx, "Processing source ", sourceConfig.Name, "...")

			sourceCtx := context.WithValue(context.Background(), constants.CtxSourceKey, sourceConfig.Name)
			source, err := source.NewSource(sourceCtx, sourceConfig, ssotLogger, netboxInventory)
			if err != nil {
				ssotLogger.Error(sourceCtx, err)
				return
			}
			ssotLogger.Infof(sourceCtx, "Successfully created source %s", constants.CheckMark)
			ssotLogger.Debugf(sourceCtx, "Source content: %s", source)

			// Run each source in parallel
			ssotLogger.Info(sourceCtx, "Initializing source")
			err = source.Init()
			if err != nil {
				ssotLogger.Error(sourceCtx, err)
				return
			}
			ssotLogger.Infof(sourceCtx, "Successfully initialized source %s", constants.CheckMark)

			// Source synchronization
			ssotLogger.Info(sourceCtx, "Syncing source...")
			err = source.Sync(netboxInventory)
			if err != nil {
				ssotLogger.Error(sourceCtx, err)
				return
			}
			ssotLogger.Infof(sourceCtx, "Source synced successfully %s", constants.CheckMark)
		}(i)
	}
	wg.Wait()

	// Orphan manager cleanup
	if config.Netbox.RemoveOrphans {
		ssotLogger.Info(mainCtx, "Cleaning up orphaned objects...")
		err = netboxInventory.DeleteOrphans(mainCtx)
		if err != nil {
			ssotLogger.Error(mainCtx, err)
			return
		}
		ssotLogger.Infof(mainCtx, "%s Successfully removed orphans", constants.CheckMark)
	} else {
		ssotLogger.Info(mainCtx, "Skipping removing orphaned objects...")
	}

	duration := time.Since(startTime)
	minutes := int(duration.Minutes())
	seconds := int((duration - time.Duration(minutes)*time.Minute).Seconds())
	ssotLogger.Infof(mainCtx, "%s Syncing took %d min %d sec in total", constants.Rocket, minutes, seconds)
}
