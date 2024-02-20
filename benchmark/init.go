package main

import (
	"fmt"

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
	// Initialize Logger
	mainLogger, err := logger.New(config.Logger.Dest, config.Logger.Level, "main")
	if err != nil {
		fmt.Println("Logger:", err)
		return
	}
	inventoryLogger, err := logger.New(config.Logger.Dest, config.Logger.Level, "netboxInventory")
	if err != nil {
		mainLogger.Errorf("inventoryLogger: %s", err)
	}
	nbi := inventory.NewNetboxInventory(inventoryLogger, config.Netbox)
	mainLogger.Debug("Netbox inventory: ", nbi)

	err = nbi.Init()
	if err != nil {
		mainLogger.Error(err)
		return
	}
	initSites(NumberOfSites, nbi)
	InitManufacturers(NumberOfManufacturers, nbi)
	InitPlatforms(NumberOfPlatforms, nbi)
	initContacts(NumberOfContacts, nbi)
}

func initSites(n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		siteName := fmt.Sprintf("Site %d", i)
		_, err := nbi.AddSite(&objects.Site{
			Name: siteName,
			Slug: utils.Slugify(siteName),
		})
		if err != nil {
			fmt.Printf("Adding site: %s", err)
		}
	}
}

func initContacts(n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		contactName := fmt.Sprintf("Contact %d", i)
		_, err := nbi.AddContact(&objects.Contact{
			Name:  contactName,
			Email: fmt.Sprintf("user%d@example.com", i),
		})
		if err != nil {
			fmt.Printf("Adding contact: %s", err)
		}
	}
}

func InitManufacturers(n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		manufacturerName := fmt.Sprintf("Manufacturer %d", i)
		_, err := nbi.AddManufacturer(&objects.Manufacturer{
			Name: manufacturerName,
			Slug: utils.Slugify(manufacturerName),
		})
		if err != nil {
			fmt.Printf("Adding manufacturer: %s", err)
		}
	}
}

func InitPlatforms(n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		platformName := fmt.Sprintf("Platform %d", i)
		_, err := nbi.AddPlatform(&objects.Platform{
			Name:         platformName,
			Slug:         utils.Slugify(platformName),
			Manufacturer: nbi.ManufacturersIndexByName[fmt.Sprintf("Manufacturer %d", i%NumberOfManufacturers)],
		})
		if err != nil {
			fmt.Printf("Adding platform: %s", err)
		}
	}
}

func InitVMs(n int, nbi *inventory.NetboxInventory) {
	for i := 0; i < n; i++ {
		vmName := fmt.Sprintf("VM %d", i)
		_, err := nbi.AddVM(&objects.VM{
			Name: vmName,
		})
		if err != nil {
			fmt.Printf("Adding VM: %s", err)
		}
	}
}
