package main

import (
	"context"
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

const (
	NumberOfSites         = 10
	NumberOfManufacturers = 100
	NumberOfPlatforms     = 1000
	NumberOfContacts      = 1000
)

func main() {
	config, err := parser.ParseConfig("config.yaml")
	if err != nil {
		fmt.Println("Parser:", err)
		return
	}
	benchmarkCtx := context.WithValue(context.Background(), constants.CtxSourceKey, "benchmark")
	// Initialize Logger
	mainLogger, err := logger.New(config.Logger.Dest, config.Logger.Level)
	if err != nil {
		fmt.Println("Logger:", err)
		return
	}
	inventoryLogger, err := logger.New(config.Logger.Dest, config.Logger.Level)
	if err != nil {
		mainLogger.Errorf(benchmarkCtx, "inventoryLogger: %s", err)
	}
	nbi := inventory.NewNetboxInventory(benchmarkCtx, inventoryLogger, config.Netbox)
	mainLogger.Debug(benchmarkCtx, "Netbox inventory: ", nbi)

	err = nbi.Init()
	if err != nil {
		mainLogger.Error(benchmarkCtx, err)
		return
	}
	initSites(benchmarkCtx, NumberOfSites, nbi)
	InitManufacturers(benchmarkCtx, NumberOfManufacturers, nbi)
	InitPlatforms(benchmarkCtx, NumberOfPlatforms, nbi)
	initContacts(benchmarkCtx, NumberOfContacts, nbi)
}

func initSites(ctx context.Context, n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		siteName := fmt.Sprintf("Site %d", i)
		_, err := nbi.AddSite(ctx, &objects.Site{
			Name: siteName,
			Slug: utils.Slugify(siteName),
		})
		if err != nil {
			fmt.Printf("Adding site: %s", err)
		}
	}
}

func initContacts(ctx context.Context, n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		contactName := fmt.Sprintf("Contact %d", i)
		_, err := nbi.AddContact(ctx, &objects.Contact{
			Name:  contactName,
			Email: fmt.Sprintf("user%d@example.com", i),
		})
		if err != nil {
			fmt.Printf("Adding contact: %s", err)
		}
	}
}

func InitManufacturers(ctx context.Context, n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		manufacturerName := fmt.Sprintf("Manufacturer %d", i)
		_, err := nbi.AddManufacturer(ctx, &objects.Manufacturer{
			Name: manufacturerName,
			Slug: utils.Slugify(manufacturerName),
		})
		if err != nil {
			fmt.Printf("Adding manufacturer: %s", err)
		}
	}
}

func InitPlatforms(ctx context.Context, n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		platformName := fmt.Sprintf("Platform %d", i)
		_, err := nbi.AddPlatform(ctx, &objects.Platform{
			Name:         platformName,
			Slug:         utils.Slugify(platformName),
			Manufacturer: nbi.ManufacturersIndexByName[fmt.Sprintf("Manufacturer %d", i%NumberOfManufacturers)],
		})
		if err != nil {
			fmt.Printf("Adding platform: %s", err)
		}
	}
}

func InitVMs(ctx context.Context, n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		vmName := fmt.Sprintf("VM %d", i)
		_, err := nbi.AddVM(ctx, &objects.VM{
			Name: vmName,
		})
		if err != nil {
			fmt.Printf("Adding VM: %s", err)
		}
	}
}
