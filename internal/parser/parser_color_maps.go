package parser

import "github.com/bl4ko/netbox-ssot/internal/netbox/objects"

// Default mappings of sources to colors (for tags)
var DefaultSourceToTagColorMap = map[SourceType]string{
	Ovirt:  objects.COLOR_DARK_RED,
	Vmware: objects.COLOR_LIGHT_GREEN,
}

// Object for mapping source type to tag color
var SourceTypeToTagColorMap = map[SourceType]string{
	Ovirt:  objects.COLOR_RED,
	Vmware: objects.COLOR_GREEN,
}
