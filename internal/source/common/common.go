package common

import (
	"context"
	"crypto/x509"

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
	Logger         *logger.Logger
	SourceConfig   *parser.SourceConfig
	CustomCertPool *x509.CertPool
	SourceNameTag  *objects.Tag
	SourceTypeTag  *objects.Tag
	Ctx            context.Context //nolint:containedctx
	CAFile         string          // path to the ca file
}

func (c Config) GetSourceTags() []*objects.Tag {
	return []*objects.Tag{c.SourceNameTag, c.SourceTypeTag}
}
