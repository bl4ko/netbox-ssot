// Common structs and interfaces for all sources
package source

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/pkg/logger"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/objects"
	"github.com/bl4ko/netbox-ssot/pkg/parser"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
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
	SourceTags   []*objects.Tag
}

// NewSource creates a Source from the given configuration
func NewSource(config *parser.SourceConfig, logger *logger.Logger, netboxInventory *inventory.NetBoxInventory) (Source, error) {
	// First we create default tags for the source
	sourceTag, err := netboxInventory.AddTag(&objects.Tag{
		Name:        config.Tag,
		Slug:        utils.Slugify("source-" + config.Name),
		Color:       config.TagColor,
		Description: fmt.Sprintf("Automatically created tag by netbox-ssot for source %s", config.Name),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sourceTag: %s", err)
	}
	sourceTypeTag, err := netboxInventory.AddTag(&objects.Tag{
		Name:        string(config.Type),
		Slug:        utils.Slugify("type-" + string(config.Type)),
		Color:       parser.SourceTypeToTagColorMap[config.Type],
		Description: fmt.Sprintf("Automatically created tag by netbox-ssot for source type %s", config.Type),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sourceTypeTag: %s", err)
	}
	commonConfig := CommonConfig{
		Logger:       logger,
		SourceConfig: config,
		SourceTags:   []*objects.Tag{sourceTag, sourceTypeTag},
	}
	switch config.Type {
	case "ovirt":
		return &OVirtSource{CommonConfig: commonConfig}, nil
	default:
		return nil, fmt.Errorf("unsupported source type: %s", config.Type)
	}
}
