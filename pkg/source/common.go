// Common structs and interfaces for all sources
package source

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/pkg/logger"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/extras"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/pkg/parser"
)

// Source is an interface for all sources (e.g. oVirt, VMware, etc.)
type Source interface {
	// Init initializes the source
	Init() error
	// Sync syncs the source to NetBox inventory
	Sync(*inventory.NetBoxInventory) error
}

// CommonConfig is a common configuration that all sources share
type CommonConfig struct {
	Logger       *logger.Logger
	SourceConfig *parser.SourceConfig
	SourceTag    *extras.Tag
}

// NewSource creates a Source from the given configuration
func NewSource(config *parser.SourceConfig, logger *logger.Logger, sourceTag *extras.Tag) (Source, error) {
	commonConfig := CommonConfig{
		Logger:       logger,
		SourceConfig: config,
		SourceTag:    sourceTag,
	}
	switch config.Type {
	case "ovirt":
		return &OVirtSource{CommonConfig: commonConfig}, nil
	default:
		return nil, fmt.Errorf("unsupported source type: %s", config.Type)
	}
}
