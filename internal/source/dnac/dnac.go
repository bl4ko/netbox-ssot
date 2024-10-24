package dnac

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	dnac "github.com/cisco-en-programmability/dnacenter-go-sdk/v5/sdk"
)

//nolint:revive
type DnacSource struct {
	common.Config

	// Dnac fetched data. Initialized in init functions.
	Sites                           map[string]dnac.ResponseSitesGetSiteResponse                // SiteID -> Site
	Devices                         map[string]dnac.ResponseDevicesGetDeviceListResponse        // DeviceID -> Device
	Interfaces                      map[string]dnac.ResponseDevicesGetAllInterfacesResponse     // InterfaceID -> Interface
	Vlans                           map[int]dnac.ResponseDevicesGetDeviceInterfaceVLANsResponse // VlanID -> Vlan
	WirelessLANInterfaceName2VlanID map[string]int                                              // InterfaceName -> VlanID
	SSID2WirelessProfileDetails     map[string]dnac.ResponseItemWirelessGetWirelessProfileProfileDetailsSSIDDetails
	SSID2WlanGroupName              map[string]string                                                // SSID -> WirelessLANGroup name
	SSID2SecurityDetails            map[string]dnac.ResponseItemWirelessGetEnterpriseSSIDSSIDDetails // WirelessLANName -> SSIDDetails

	// Relations between dnac data. Initialized in init functions.
	Site2Devices          map[string]map[string]bool // Site ID - > set of device IDs
	Device2Site           map[string]string          // Device ID -> Site ID
	DeviceID2InterfaceIDs map[string][]string        // DeviceID -> []InterfaceID

	// Netbox related data for easier access. Initialized in sync functions.
	VID2nbVlan              sync.Map // VlanID -> nbVlan
	SiteID2nbSite           sync.Map // SiteID -> nbSite
	DeviceID2nbDevice       sync.Map // DeviceID -> nbDevice
	InterfaceID2nbInterface sync.Map // InterfaceID -> nbInterface

	// User defined relations
	HostTenantRelations map[string]string
	VlanGroupRelations  map[string]string
	VlanTenantRelations map[string]string
	WlanTenantRelations map[string]string
}

func (ds *DnacSource) Init() error {
	dnacURL := fmt.Sprintf("%s://%s:%d", ds.Config.SourceConfig.HTTPScheme, ds.Config.SourceConfig.Hostname, ds.Config.SourceConfig.Port)
	Client, err := dnac.NewClientWithOptions(dnacURL, ds.SourceConfig.Username, ds.SourceConfig.Password, "false", strconv.FormatBool(ds.SourceConfig.ValidateCert), nil)
	if err != nil {
		return fmt.Errorf("creating dnac client: %s", err)
	}

	// Initialize regex relations for this source
	ds.VlanGroupRelations = utils.ConvertStringsToRegexPairs(ds.SourceConfig.VlanGroupRelations)
	ds.Logger.Debugf(ds.Ctx, "VlanGroupRelations: %s", ds.VlanGroupRelations)
	ds.VlanTenantRelations = utils.ConvertStringsToRegexPairs(ds.SourceConfig.VlanTenantRelations)
	ds.Logger.Debugf(ds.Ctx, "VlanTenantRelations: %s", ds.VlanTenantRelations)
	ds.HostTenantRelations = utils.ConvertStringsToRegexPairs(ds.SourceConfig.HostTenantRelations)
	ds.Logger.Debugf(ds.Ctx, "HostTenantRelations: %s", ds.HostTenantRelations)
	ds.WlanTenantRelations = utils.ConvertStringsToRegexPairs(ds.SourceConfig.WlanTenantRelations)

	// Initialize items from vsphere API to local storage
	initFunctions := []func(*dnac.Client) error{
		ds.initSites,
		ds.initMemberships,
		ds.initDevices,
		ds.initInterfaces,
		ds.initWirelessLANs,
	}

	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(Client); err != nil {
			return fmt.Errorf("dnac initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		ds.Logger.Infof(ds.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"), duration.Seconds())
	}
	return nil
}

func (ds *DnacSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		ds.syncSites,
		ds.syncVlans,
		ds.syncDevices,
		ds.syncDeviceInterfaces,
		ds.syncWirelessLANs,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		ds.Logger.Infof(ds.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"), duration.Seconds())
	}
	return nil
}
