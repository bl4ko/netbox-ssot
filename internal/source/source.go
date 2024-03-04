// Common structs and interfaces for all sources
package source

import (
	"context"
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/source/dnac"
	"github.com/bl4ko/netbox-ssot/internal/source/ovirt"
	"github.com/bl4ko/netbox-ssot/internal/source/proxmox"
	"github.com/bl4ko/netbox-ssot/internal/source/vmware"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// NewSource creates a Source from the given configuration.
func NewSource(ctx context.Context, config *parser.SourceConfig, logger *logger.Logger, netboxInventory *inventory.NetboxInventory) (common.Source, error) {
	// First we create default tags for the source
	sourceTag, err := netboxInventory.AddTag(ctx, &objects.Tag{
		Name:        config.Tag,
		Slug:        utils.Slugify("source-" + config.Name),
		Color:       objects.Color(config.TagColor),
		Description: fmt.Sprintf("Automatically created tag by netbox-ssot for source %s", config.Name),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sourceTag: %s", err)
	}
	sourceTypeTag, err := netboxInventory.AddTag(ctx, &objects.Tag{
		Name:        string(config.Type),
		Slug:        utils.Slugify("type-" + string(config.Type)),
		Color:       objects.Color(constants.SourceTypeToTagColorMap[config.Type]),
		Description: fmt.Sprintf("Automatically created tag by netbox-ssot for source type %s", config.Type),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sourceTypeTag: %s", err)
	}
	commonConfig := common.Config{
		Logger:       logger,
		SourceConfig: config,
		SourceTags:   []*objects.Tag{sourceTag, sourceTypeTag},
		Ctx:          ctx,
	}

	switch config.Type {
	case constants.Ovirt:
		return &ovirt.OVirtSource{Config: commonConfig}, nil
	case constants.Vmware:
		return &vmware.VmwareSource{Config: commonConfig}, nil
	case constants.Dnac:
		return &dnac.DnacSource{Config: commonConfig}, nil
	case constants.Proxmox:
		return &proxmox.ProxmoxSource{Config: commonConfig}, nil
	default:
		return nil, fmt.Errorf("unsupported source type: %s", config.Type)
	}
}
