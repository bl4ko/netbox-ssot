package dnac

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	dnac "github.com/cisco-en-programmability/dnacenter-go-sdk/v5/sdk"
)

type DnacSource struct {
	common.CommonConfig

	// Dnac fetched data. Initialized in init functions.
	Sites      map[string]dnac.ResponseSitesGetSiteResponse                // SiteId -> Site
	Devices    map[string]dnac.ResponseDevicesGetDeviceListResponse        // DeviceId -> Device
	Interfaces map[string]dnac.ResponseDevicesGetAllInterfacesResponse     // InterfaceId -> Interface
	Vlans      map[int]dnac.ResponseDevicesGetDeviceInterfaceVLANsResponse // VlanId -> Vlan
	// Relations between dnac data. Initialized in init functions.
	Site2Devices          map[string]map[string]bool // Site Id - > set of device Ids
	Device2Site           map[string]string          // Device Id -> Site Id
	DeviceId2InterfaceIds map[string][]string        // DeviceId -> []InterfaceId

	// Netbox related data for easier access. Initialized in sync functions.
	SiteId2nbSite           map[string]*objects.Site   // SiteId -> nbSite
	DeviceId2nbDevice       map[string]*objects.Device // DeviceId -> nbDevice
	InterfaceId2nbInterface map[string]*objects.Interface

	// User defined relations
	VlanGroupRelations  map[string]string
	VlanTenantRelations map[string]string
}

func (ds *DnacSource) Init() error {
	dnacUrl := fmt.Sprintf("%s://%s:%d", ds.CommonConfig.SourceConfig.HTTPScheme, ds.CommonConfig.SourceConfig.Hostname, ds.CommonConfig.SourceConfig.Port)
	Client, err := dnac.NewClientWithOptions(dnacUrl, ds.SourceConfig.Username, ds.SourceConfig.Password, "false", strconv.FormatBool(ds.SourceConfig.ValidateCert), nil)
	if err != nil {
		return fmt.Errorf("creating dnac client: %s", err)
	}

	// Initialize regex relations for this source
	ds.VlanGroupRelations = utils.ConvertStringsToRegexPairs(ds.SourceConfig.VlanGroupRelations)
	ds.Logger.Debug("VlanGroupRelations: ", ds.VlanGroupRelations)
	ds.VlanTenantRelations = utils.ConvertStringsToRegexPairs(ds.SourceConfig.VlanTenantRelations)
	ds.Logger.Debug("VlanTenantRelations: ", ds.VlanTenantRelations)

	// Initialize items from vsphere API to local storage
	initFunctions := []func(*dnac.Client) error{
		ds.InitSites,
		ds.InitMemberships,
		ds.InitDevices,
		ds.InitInterfaces,
	}

	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(Client); err != nil {
			return fmt.Errorf("dnac initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		ds.Logger.Infof("Successfully initialized %s in %f seconds", utils.ExtractFunctionName(initFunc), duration.Seconds())
	}

	return nil
}

func (ds *DnacSource) Sync(nbi *inventory.NetboxInventory) error {
	// initialize variables, that are shared between sync functions
	ds.SiteId2nbSite = make(map[string]*objects.Site)
	ds.DeviceId2nbDevice = make(map[string]*objects.Device)
	ds.InterfaceId2nbInterface = make(map[string]*objects.Interface)

	syncFunctions := []func(*inventory.NetboxInventory) error{
		ds.SyncSites,
		ds.SyncVlans,
		ds.SyncDevices,
		ds.SyncDeviceInterfaces,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		ds.Logger.Infof("Successfully synced %s in %f seconds", utils.ExtractFunctionName(syncFunc), duration.Seconds())
	}
	return nil
}
