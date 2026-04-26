package inventory

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
)

var MockExistingTags = map[string]*objects.Tag{
	"existing_tag1": {
		Name:        "existing_tag1",
		Description: "Test exististing tag1",
		Slug:        "existing_tag1",
	},
	"existing_tag2": {
		Name:        "existing_tag2",
		Description: "Test exististing tag2",
		Slug:        "existing_tag2",
	},
}

var MockExistingTenants = map[string]*objects.Tenant{
	"existing_tenant1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_tenant1",
		Slug: "existing_tenant1",
	},
	"existing_tenant2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_tenant2",
		Slug: "existing_tenant2",
	},
}

var MockExistingSites = map[string]*objects.Site{
	"existing_site1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_site1",
		Slug: "existing_site1",
	},
	"existing_site2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_site2",
		Slug: "existing_site2",
	},
}

// MockExistingPrefixes simulates prefixes fetched from the NetBox API.
// The "10.0.0.0/24" prefix has an object-type custom field (site_ref)
// returned as a nested object — this is the read format from the API.
var MockExistingPrefixes = map[string]map[int]*objects.Prefix{
	"10.0.0.0/24": {
		0: {
			NetboxObject: objects.NetboxObject{
				ID:   1,
				Tags: []*objects.Tag{service.MockDefaultSsotTag},
				CustomFields: map[string]interface{}{
					"source":           "test",
					"orphan_last_seen": nil,
					"site_ref": map[string]interface{}{
						"id":      float64(1),
						"display": "LCL",
						"url":     "https://netbox/api/dcim/sites/1/",
						"name":    "LCL",
						"slug":    "lcl",
					},
				},
			},
			Prefix: "10.0.0.0/24",
		},
	},
}

var MockExistingContactRoles = map[string]*objects.ContactRole{
	"existing_contact_role1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_contact_role1",
		Slug: "existing_contact_role1",
	},
	"existing_contact_role2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_contact_role2",
		Slug: "existing_contact_role2",
	},
}

var MockExistingContactGroups = map[string]*objects.ContactGroup{
	"existing_contact_group1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_contact_group1",
		Slug: "existing_contact_group1",
	},
	"existing_contact_group2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_contact_group2",
		Slug: "existing_contact_group2",
	},
}

var MockExistingContacts = map[string]*objects.Contact{
	"existing_contact1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_contact1",
	},
	"existing_contact2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_contact2",
	},
}

var MockExistingContactAssignments = map[constants.ContentType]map[int]map[int]map[int]*objects.ContactAssignment{
	constants.ContentTypeDcimDevice: {
		1: {
			1: {
				1: {
					NetboxObject: objects.NetboxObject{
						ID:   1,
						Tags: []*objects.Tag{service.MockDefaultSsotTag},
					},
					ModelType: constants.ContentTypeDcimDevice,
					ObjectID:  1,
					Contact:   &objects.Contact{NetboxObject: objects.NetboxObject{ID: 1}, Name: "existing_contact1"},
					Role:      &objects.ContactRole{NetboxObject: objects.NetboxObject{ID: 1}, Name: "existing_contact_role1"},
				},
			},
		},
	},
}

var MockExistingCustomFields = map[string]*objects.CustomField{
	"existing_cf1": {
		ID:   1,
		Name: "existing_cf1",
		Type: objects.CustomFieldTypeText,
	},
	"existing_cf2": {
		ID:   2, //nolint:mnd
		Name: "existing_cf2",
		Type: objects.CustomFieldTypeText,
	},
}

var MockExistingClusterGroups = map[string]*objects.ClusterGroup{
	"existing_cluster_group1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_cluster_group1",
		Slug: "existing_cluster_group1",
	},
	"existing_cluster_group2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_cluster_group2",
		Slug: "existing_cluster_group2",
	},
}

var MockExistingClusterTypes = map[string]*objects.ClusterType{
	"existing_cluster_type1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_cluster_type1",
		Slug: "existing_cluster_type1",
	},
	"existing_cluster_type2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_cluster_type2",
		Slug: "existing_cluster_type2",
	},
}

var MockExistingClusters = map[string]*objects.Cluster{
	"existing_cluster1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_cluster1",
	},
	"existing_cluster2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_cluster2",
	},
}

var MockExistingDeviceRoles = map[string]*objects.DeviceRole{
	"existing_device_role1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name:  "existing_device_role1",
		Slug:  "existing_device_role1",
		Color: constants.Color(constants.DeviceRoleServerColor),
	},
	"existing_device_role2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name:  "existing_device_role2",
		Slug:  "existing_device_role2",
		Color: constants.Color(constants.DeviceRoleServerColor),
	},
}

var MockExistingManufacturers = map[string]*objects.Manufacturer{
	"existing_manufacturer1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_manufacturer1",
		Slug: "existing_manufacturer1",
	},
	"existing_manufacturer2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_manufacturer2",
		Slug: "existing_manufacturer2",
	},
}

var MockExistingDeviceTypes = map[string]*objects.DeviceType{
	"existing_device_type1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Model: "existing_device_type1",
		Slug:  "existing_device_type1",
	},
	"existing_device_type2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Model: "existing_device_type2",
		Slug:  "existing_device_type2",
	},
}

var MockExistingPlatforms = map[string]*objects.Platform{
	"existing_platform1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_platform1",
		Slug: "existing_platform1",
	},
	"existing_platform2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_platform2",
		Slug: "existing_platform2",
	},
}

var mockSite1 = &objects.Site{
	NetboxObject: objects.NetboxObject{ID: 1},
	Name:         "site1",
	Slug:         "site1",
}

var mockDeviceRole1 = &objects.DeviceRole{
	NetboxObject: objects.NetboxObject{ID: 1},
	Name:         "role1",
	Slug:         "role1",
	Color:        "aa1409",
}

var mockDeviceType1 = &objects.DeviceType{
	NetboxObject: objects.NetboxObject{ID: 1},
	Model:        "type1",
	Slug:         "type1",
}

var MockExistingDevices = map[string]map[int]*objects.Device{
	"existing_device1": {
		1: {
			NetboxObject: objects.NetboxObject{
				ID:   1,
				Tags: []*objects.Tag{service.MockDefaultSsotTag},
			},
			Name:       "existing_device1",
			Site:       mockSite1,
			DeviceRole: mockDeviceRole1,
			DeviceType: mockDeviceType1,
		},
	},
}

var mockDevice1 = &objects.Device{
	NetboxObject: objects.NetboxObject{ID: 1},
	Name:         "existing_device1",
	Site:         mockSite1,
}

var MockExistingDevicesByID = map[int]*objects.Device{
	1: MockExistingDevices["existing_device1"][1],
}

var MockExistingVDCs = map[string]map[int]*objects.VirtualDeviceContext{
	"existing_vdc1": {
		1: {
			NetboxObject: objects.NetboxObject{
				ID:   1,
				Tags: []*objects.Tag{service.MockDefaultSsotTag},
			},
			Name:   "existing_vdc1",
			Device: mockDevice1,
		},
	},
}

var MockExistingVlanGroups = map[string]*objects.VlanGroup{
	"existing_vlan_group1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_vlan_group1",
		Slug: "existing_vlan_group1",
	},
	"existing_vlan_group2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_vlan_group2",
		Slug: "existing_vlan_group2",
	},
}

var mockVlanGroup1 = &objects.VlanGroup{
	NetboxObject: objects.NetboxObject{ID: 1},
	Name:         "vlan_group1",
	Slug:         "vlan_group1",
}

var MockExistingVlans = map[int]map[int]*objects.Vlan{
	1: {
		100: {
			NetboxObject: objects.NetboxObject{
				ID:   1,
				Tags: []*objects.Tag{service.MockDefaultSsotTag},
			},
			Name:  "existing_vlan100",
			Vid:   100, //nolint:mnd
			Group: mockVlanGroup1,
		},
	},
}

var MockExistingInterfaces = map[int]map[string]*objects.Interface{
	1: {
		"eth0": {
			NetboxObject: objects.NetboxObject{
				ID:   1,
				Tags: []*objects.Tag{service.MockDefaultSsotTag},
			},
			Name:   "eth0",
			Device: mockDevice1,
			Type:   &objects.VirtualInterfaceType,
		},
	},
}

var MockExistingInterfacesByID = map[int]*objects.Interface{
	1: MockExistingInterfaces[1]["eth0"],
}

var mockCluster1 = &objects.Cluster{
	NetboxObject: objects.NetboxObject{ID: 1},
	Name:         "cluster1",
}

var mockVM1 = &objects.VM{
	NetboxObject: objects.NetboxObject{ID: 1},
	Name:         "existing_vm1",
	Cluster:      mockCluster1,
}

var MockExistingVMs = map[string]map[int]*objects.VM{
	"existing_vm1": {
		1: {
			NetboxObject: objects.NetboxObject{
				ID:   1,
				Tags: []*objects.Tag{service.MockDefaultSsotTag},
			},
			Name:    "existing_vm1",
			Cluster: mockCluster1,
		},
	},
}

var MockExistingVMsByID = map[int]*objects.VM{
	1: MockExistingVMs["existing_vm1"][1],
}

var MockExistingVMInterfaces = map[int]map[string]*objects.VMInterface{
	1: {
		"vmeth0": {
			NetboxObject: objects.NetboxObject{
				ID:   1,
				Tags: []*objects.Tag{service.MockDefaultSsotTag},
			},
			Name: "vmeth0",
			VM:   mockVM1,
		},
	},
}

var MockExistingVMInterfacesByID = map[int]*objects.VMInterface{
	1: MockExistingVMInterfaces[1]["vmeth0"],
}

var MockExistingIPAddresses = map[constants.ContentType]map[string]map[string]map[string]*objects.IPAddress{
	constants.ContentTypeDcimDevice: {
		"eth0": {
			"existing_device1": {
				"10.0.0.1/24": {
					NetboxObject: objects.NetboxObject{
						ID:   1,
						Tags: []*objects.Tag{service.MockDefaultSsotTag},
					},
					Address:            "10.0.0.1/24",
					AssignedObjectType: constants.ContentTypeDcimInterface,
					AssignedObjectID:   1,
				},
			},
		},
	},
}

var MockExistingMACAddresses = map[constants.ContentType]map[string]map[string]map[string]*objects.MACAddress{
	constants.ContentTypeDcimDevice: {
		"eth0": {
			"existing_device1": {
				"AA:BB:CC:DD:EE:FF": {
					NetboxObject: objects.NetboxObject{
						ID:   1,
						Tags: []*objects.Tag{service.MockDefaultSsotTag},
					},
					MAC:                "AA:BB:CC:DD:EE:FF",
					AssignedObjectType: constants.ContentTypeDcimInterface,
					AssignedObjectID:   1,
				},
			},
		},
	},
}

var MockExistingWirelessLANs = map[string]*objects.WirelessLAN{
	"existing_wlan1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		SSID: "existing_wlan1",
	},
	"existing_wlan2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		SSID: "existing_wlan2",
	},
}

var MockExistingWirelessLANGroups = map[string]*objects.WirelessLANGroup{
	"existing_wlan_group1": {
		NetboxObject: objects.NetboxObject{
			ID:   1,
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_wlan_group1",
		Slug: "existing_wlan_group1",
	},
	"existing_wlan_group2": {
		NetboxObject: objects.NetboxObject{
			ID:   2, //nolint:mnd
			Tags: []*objects.Tag{service.MockDefaultSsotTag},
		},
		Name: "existing_wlan_group2",
		Slug: "existing_wlan_group2",
	},
}

var MockExistingVirtualDisks = map[int]map[string]*objects.VirtualDisk{
	1: {
		"existing_disk1": {
			NetboxObject: objects.NetboxObject{
				ID:   1,
				Tags: []*objects.Tag{service.MockDefaultSsotTag},
			},
			Name: "existing_disk1",
			Size: 100, //nolint:mnd
			VM:   mockVM1,
		},
	},
}

var mockLogger = &logger.Logger{Logger: log.New(os.Stdout, "", log.LstdFlags)}

var MockInventory = &NetboxInventory{
	Logger:                           mockLogger,
	tagsIndexByName:                  MockExistingTags,
	tagsLock:                         sync.Mutex{},
	tenantsIndexByName:               MockExistingTenants,
	tenantsLock:                      sync.Mutex{},
	sitesIndexByName:                 MockExistingSites,
	sitesLock:                        sync.Mutex{},
	prefixesIndexByPrefix:            MockExistingPrefixes,
	prefixesLock:                     sync.Mutex{},
	contactRolesIndexByName:          MockExistingContactRoles,
	contactRolesLock:                 sync.Mutex{},
	contactGroupsIndexByName:         MockExistingContactGroups,
	contactGroupsLock:                sync.Mutex{},
	contactsIndexByName:              MockExistingContacts,
	contactsLock:                     sync.Mutex{},
	contactAssignmentsIndex:          MockExistingContactAssignments,
	contactAssignmentsLock:           sync.Mutex{},
	customFieldsIndexByName:          MockExistingCustomFields,
	customFieldsLock:                 sync.Mutex{},
	clusterGroupsIndexByName:         MockExistingClusterGroups,
	clusterGroupsLock:                sync.Mutex{},
	clusterTypesIndexByName:          MockExistingClusterTypes,
	clusterTypesLock:                 sync.Mutex{},
	clustersIndexByName:              MockExistingClusters,
	clustersLock:                     sync.Mutex{},
	deviceRolesIndexByName:           MockExistingDeviceRoles,
	deviceRolesLock:                  sync.Mutex{},
	manufacturersIndexByName:         MockExistingManufacturers,
	manufacturersLock:                sync.Mutex{},
	deviceTypesIndexByModel:          MockExistingDeviceTypes,
	deviceTypesLock:                  sync.Mutex{},
	platformsIndexByName:             MockExistingPlatforms,
	platformsLock:                    sync.Mutex{},
	devicesIndexByNameAndSiteID:      MockExistingDevices,
	devicesIndexByID:                 MockExistingDevicesByID,
	devicesLock:                      sync.Mutex{},
	virtualDeviceContextsIndex:       MockExistingVDCs,
	virtualDeviceContextsLock:        sync.Mutex{},
	vlanGroupsIndexByName:            MockExistingVlanGroups,
	vlanGroupsLock:                   sync.Mutex{},
	vlansIndexByVlanGroupIDAndVID:    MockExistingVlans,
	vlansLock:                        sync.Mutex{},
	interfacesIndexByDeviceIDAndName: MockExistingInterfaces,
	interfacesIndexByID:              MockExistingInterfacesByID,
	interfacesLock:                   sync.Mutex{},
	vmsIndexByNameAndClusterID:       MockExistingVMs,
	vmsIndexByID:                     MockExistingVMsByID,
	vmsLock:                          sync.Mutex{},
	vmInterfacesIndexByVMIdAndName:   MockExistingVMInterfaces,
	vmInterfacesIndexByID:            MockExistingVMInterfacesByID,
	vmInterfacesLock:                 sync.Mutex{},
	ipAddressesIndex:                 MockExistingIPAddresses,
	ipAddressesLock:                  sync.Mutex{},
	macAddressesIndex:                MockExistingMACAddresses,
	macAddressesLock:                 sync.Mutex{},
	wirelessLANsIndexBySSID:          MockExistingWirelessLANs,
	wirelessLANsLock:                 sync.Mutex{},
	wirelessLANGroupsIndexByName:     MockExistingWirelessLANGroups,
	wirelessLANGroupsLock:            sync.Mutex{},
	virtualDisksIndexByVMIDAndName:   MockExistingVirtualDisks,
	virtualDisksLock:                 sync.Mutex{},
	vrfsIndexByName:                  map[string]*objects.VRF{},
	vrfsLock:                         sync.Mutex{},
	locationsIndexByName:             map[string]*objects.Location{},
	locationsLock:                    sync.Mutex{},
	siteGroupsIndexByName:            map[string]*objects.SiteGroup{},
	siteGroupsLock:                   sync.Mutex{},
	NetboxAPI:                        service.MockNetboxClient,
	OrphanManager:                    NewOrphanManager(mockLogger),
	SourcePriority:                   map[string]int{},
	Ctx: context.WithValue(
		context.Background(),
		constants.CtxSourceKey,
		"testInventory",
	),
	SsotTag: &objects.Tag{
		ID:          0,
		Name:        "netbox-ssot",
		Slug:        "netbox-ssot",
		Description: "default netbox-ssot tag",
		Color:       "ffffff",
	},
}
