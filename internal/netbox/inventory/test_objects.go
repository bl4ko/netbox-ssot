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

var MockInventory = &NetboxInventory{
	Logger:             &logger.Logger{Logger: log.New(os.Stdout, "", log.LstdFlags)},
	tagsIndexByName:    MockExistingTags,
	tagsLock:           sync.Mutex{},
	tenantsIndexByName: MockExistingTenants,
	tenantsLock:        sync.Mutex{},
	sitesIndexByName:   MockExistingSites,
	sitesLock:          sync.Mutex{},
	NetboxAPI:          service.MockNetboxClient,
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
