package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
)

var (
	configPath = flag.String("config", "config.yaml", "Path to the configuration file")
	dryRun     = flag.Bool("dry-run", false, "Preview changes without writing to Netbox")
)

// Build variables provided with ldflags.
var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Print build information
	fmt.Printf("Running version %s built on %s (commit %s)\n\n", version, date, commit)

	startTime := time.Now()

	// Parse configuration
	fmt.Printf("Netbox-SSOT has started at %s\n", startTime.Format(time.RFC3339))
	flag.Parse()
	config, err := parser.ParseConfig(*configPath)
	if err != nil {
		fmt.Println("Parser:", err)
		os.Exit(1)
	}

	// Create our main context
	mainCtx := context.Background()
	mainCtx = context.WithValue(mainCtx, constants.CtxSourceKey, "main")

	// Initialize Logger
	ssotLogger, err := logger.New(config.Logger.Dest, config.Logger.Level)
	if err != nil {
		fmt.Println("Logger:", err)
		os.Exit(1)
	}
	ssotLogger.Debug(mainCtx, "Parsed Logger config: ", config.Logger)
	ssotLogger.Debug(mainCtx, "Parsed Netbox config: ", config.Netbox)
	ssotLogger.Debug(mainCtx, "Parsed Source config: ", config.Sources)

	inventoryLogger, err := logger.New(config.Logger.Dest, config.Logger.Level)
	if err != nil {
		ssotLogger.Errorf(mainCtx, "inventoryLogger: %s", err)
		os.Exit(1)
	}
	if *dryRun {
		ssotLogger.Info(mainCtx, "DRY-RUN MODE ENABLED: No changes will be written to Netbox")
	}

	inventoryCtx := context.WithValue(context.Background(), constants.CtxSourceKey, "inventory")
	netboxInventory := inventory.NewNetboxInventory(inventoryCtx, inventoryLogger, config.Netbox, *dryRun)
	ssotLogger.Debug(mainCtx, "Netbox inventory: ", netboxInventory)

	ssotLogger.Info(mainCtx, "Starting initializing netbox inventory")
	err = netboxInventory.Init()
	if err != nil {
		ssotLogger.Error(mainCtx, err)
		os.Exit(1)
	}
	ssotLogger.Debug(mainCtx, "Netbox inventory initialized: ", netboxInventory)

	// Register the source tags owned by netbox-ssot, so tags of other tools are left alone
	sourceTagNames := make([]string, 0, len(config.Sources))
	for _, sourceConfig := range config.Sources {
		sourceTagNames = append(sourceTagNames, sourceConfig.Tag)
	}
	netboxInventory.RegisterManagedSourceTags(sourceTagNames)

	// Create every source first, then sync them all in parallel
	sources := make([]namedSource, 0, len(config.Sources))
	for i := range config.Sources {
		sourceConfig := &config.Sources[i]
		ssotLogger.Info(mainCtx, "Processing source ", sourceConfig.Name, "...")
		sourceCtx := context.WithValue(mainCtx, constants.CtxSourceKey, sourceConfig.Name)
		source, err := source.NewSource(sourceCtx, sourceConfig, ssotLogger, netboxInventory)
		if err != nil {
			ssotLogger.Error(sourceCtx, err)
			os.Exit(1)
		}
		ssotLogger.Infof(sourceCtx, "Successfully created source %s", constants.CheckMark)
		ssotLogger.Debugf(sourceCtx, "Source content: %s", source)
		sources = append(sources, namedSource{name: sourceConfig.Name, ctx: sourceCtx, source: source})
	}
	// Failed sources by name. If any failed, we don't remove orphans.
	encounteredErrors := syncSources(ssotLogger, netboxInventory, sources)
	successfullRun := len(encounteredErrors) == 0

	// Orphan manager cleanup on successful run and if enabled
	if successfullRun {
		ssotLogger.Info(mainCtx, "Cleaning up orphaned objects...")
		err = netboxInventory.DeleteOrphans(config.Netbox.RemoveOrphans)
		if err != nil {
			ssotLogger.Error(mainCtx, err)
			os.Exit(1)
		}
		ssotLogger.Infof(mainCtx, "%s Successfully removed orphans", constants.CheckMark)
	} else {
		ssotLogger.Info(mainCtx, "Skipping removing orphaned objects because run failed...")
	}

	duration := time.Since(startTime)
	minutes := int(duration.Minutes())
	seconds := int((duration - time.Duration(minutes)*time.Minute).Seconds())
	if *dryRun {
		ssotLogger.Info(mainCtx, "DRY-RUN COMPLETE: Review the log above for [DRY-RUN] entries to see what would change")
	}

	if successfullRun {
		ssotLogger.Infof(
			mainCtx,
			"%s Syncing took %d min %d sec in total",
			constants.Rocket,
			minutes,
			seconds,
		)
	} else {
		for source, err := range encounteredErrors {
			ssotLogger.Infof(mainCtx, "%s syncing of source %s failed with: %v", constants.WarningSign, source, err)
		}
		os.Exit(1)
	}
}

// namedSource is a source ready to be synced, with its name and logging context.
type namedSource struct {
	name   string
	ctx    context.Context
	source common.Source
}

// syncSources initializes and syncs all sources in parallel, and returns the
// error of each source that failed, by source name.
func syncSources(
	ssotLogger *logger.Logger,
	netboxInventory *inventory.NetboxInventory,
	sources []namedSource,
) map[string]error {
	encounteredErrors := map[string]error{}
	var errorsLock sync.Mutex
	addError := func(sourceName string, err error) {
		errorsLock.Lock()
		defer errorsLock.Unlock()
		encounteredErrors[sourceName] = err
	}
	var wg sync.WaitGroup
	for _, s := range sources {
		wg.Add(1)
		go func(s namedSource) {
			defer wg.Done()
			// Source initialization
			ssotLogger.Info(s.ctx, "Initializing source")
			if err := s.source.Init(); err != nil {
				ssotLogger.Error(s.ctx, err)
				addError(s.name, err)
				return
			}
			ssotLogger.Infof(s.ctx, "Successfully initialized source %s", constants.CheckMark)

			// Source synchronization
			ssotLogger.Info(s.ctx, "Syncing source...")
			if err := s.source.Sync(netboxInventory); err != nil {
				ssotLogger.Error(s.ctx, err)
				addError(s.name, err)
				return
			}
			ssotLogger.Infof(s.ctx, "Source synced successfully %s", constants.CheckMark)
		}(s)
	}
	wg.Wait()
	return encounteredErrors
}
