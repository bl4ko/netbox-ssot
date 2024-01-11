// Common structs and interfaces for all sources
package source

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/source/ovirt"
	"github.com/bl4ko/netbox-ssot/internal/source/vmware"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// NewSource creates a Source from the given configuration
func NewSource(config *parser.SourceConfig, logger *logger.Logger, netboxInventory *inventory.NetBoxInventory) (common.Source, error) {
	// First we create default tags for the source
	sourceTag, err := netboxInventory.AddTag(&objects.Tag{
		Name:        config.Tag,
		Slug:        utils.Slugify("source-" + config.Name),
		Color:       objects.Color(config.TagColor),
		Description: fmt.Sprintf("Automatically created tag by netbox-ssot for source %s", config.Name),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sourceTag: %s", err)
	}
	sourceTypeTag, err := netboxInventory.AddTag(&objects.Tag{
		Name:        string(config.Type),
		Slug:        utils.Slugify("type-" + string(config.Type)),
		Color:       objects.Color(parser.SourceTypeToTagColorMap[config.Type]),
		Description: fmt.Sprintf("Automatically created tag by netbox-ssot for source type %s", config.Type),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sourceTypeTag: %s", err)
	}
	commonConfig := common.CommonConfig{
		Logger:       logger,
		SourceConfig: config,
		SourceTags:   []*objects.Tag{sourceTag, sourceTypeTag},
	}
	switch config.Type {
	case "ovirt":
		return &ovirt.OVirtSource{CommonConfig: commonConfig}, nil
	case "vmware":
		return &vmware.VmwareSource{CommonConfig: commonConfig}, nil
	default:
		return nil, fmt.Errorf("unsupported source type: %s", config.Type)
	}
}
