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

var mockLogger = &logger.Logger{Logger: log.New(os.Stdout, "", log.LstdFlags)}

var MockInventory = &NetboxInventory{
	Logger:                 mockLogger,
	tagsIndexByName:        MockExistingTags,
	tagsLock:               sync.Mutex{},
	tenantsIndexByName:     MockExistingTenants,
	tenantsLock:            sync.Mutex{},
	sitesIndexByName:       MockExistingSites,
	sitesLock:              sync.Mutex{},
	prefixesIndexByPrefix:  MockExistingPrefixes,
	prefixesLock:           sync.Mutex{},
	deviceRolesIndexByName: map[string]*objects.DeviceRole{},
	deviceRolesLock:        sync.Mutex{},
	vlanGroupsIndexByName:  map[string]*objects.VlanGroup{},
	vlanGroupsLock:         sync.Mutex{},
	vrfsIndexByName:        map[string]*objects.VRF{},
	vrfsLock:               sync.Mutex{},
	NetboxAPI:              service.MockNetboxClient,
	OrphanManager:          NewOrphanManager(mockLogger),
	SourcePriority:         map[string]int{},
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
