package common

import (
	"context"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
)

// Source is an interface for all sources (e.g. oVirt, VMware, etc.).
type Source interface {
	// Init initializes the source
	Init() error
	// Sync syncs the source to Netbox inventory
	Sync(*inventory.NetboxInventory) error
}

// Config is a common configuration that all sources share.
type Config struct {
	Logger       *logger.Logger
	SourceConfig *parser.SourceConfig
	SourceTags   []*objects.Tag
	Ctx          context.Context //nolint:containedctx
}
