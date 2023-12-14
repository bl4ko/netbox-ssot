package parser

// Default mappings of sources to colors (for tags)
var DefaultSourceToTagColorMap = map[SourceType]string{
	Ovirt:  "07426b",
	Vmware: "0000ff",
}

// Object for mapping source type to tag color
var SourceTypeToTagColorMap = map[SourceType]string{
	Ovirt:  "ff0000",
	Vmware: "0000ff",
}
