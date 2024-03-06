package constants

import "github.com/bl4ko/netbox-ssot/internal/netbox/objects"

type SourceType string

const (
	Ovirt   SourceType = "ovirt"
	Vmware  SourceType = "vmware"
	Dnac    SourceType = "dnac"
	Proxmox SourceType = "proxmox"
)

const DefaultSourceName = "netbox-ssot"

const (
	DefaultOSName       string = "Generic OS"
	DefaultOSVersion    string = "Generic Version"
	DefaultManufacturer string = "Generic Manufacturer"
	DefaultModel        string = "Generic Model"
)

// Default mappings of sources to colors (for tags).
var DefaultSourceToTagColorMap = map[SourceType]string{
	Ovirt:  objects.ColorDarkRed,
	Vmware: objects.ColorLightGreen,
	Dnac:   objects.ColorLightBlue,
}

// Object for mapping source type to tag color.
var SourceTypeToTagColorMap = map[SourceType]string{
	Ovirt:  objects.ColorRed,
	Vmware: objects.ColorGreen,
	Dnac:   objects.ColorBlue,
}

const (
	// API timeout in seconds.
	DefaultAPITimeout = 30
)

// Magic numbers for dealing with bytes.
const (
	B   = 1
	KB  = 1000 * B
	MB  = 1000 * KB
	GB  = 1000 * MB
	TB  = 1000 * GB
	KiB = 1024 * B
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
)

// Magic numbers for dealing with IP addresses.
const (
	IPv4 = 4
	IPv6 = 6
)

const (
	HTTPSDefaultPort = 443
)

// Names used for netbox objects custom fields attribute.
const (
	// Custom Field for matching object with a source. This custom field is important
	// for priority diff.
	CustomFieldSourceName        = "source"
	CustomFieldSourceLabel       = "Source"
	CustomFieldSourceDescription = "Name of the source from which the object was collected"

	// Custom field for adding source ID for each object.
	CustomFieldSourceIDName        = "source_id"
	CustomFieldSourceIDLabel       = "Source ID"
	CustomFieldSourceIDDescription = "ID of the object on the source API"

	// Custom field dcim.device, so we can add number of cpu cores for each server.
	CustomFieldHostCPUCoresName        = "host_cpu_cores"
	CustomFieldHostCPUCoresLabel       = "Host CPU cores"
	CustomFieldHostCPUCoresDescription = "Number of CPU cores on the host"

	// Custom field for dcim.device, so we can add number of ram for each server.
	CustomFieldHostMemoryName        = "host_memory"
	CustomFieldHostMemoryLabel       = "Host memory"
	CustomFieldHostMemoryDescription = "Amount of memory on the host"
)

// Constants used for variables in our contexts.
type CtxKey int

const (
	CtxSourceKey CtxKey = iota
)

const (
	UntaggedVID = 0
	DefaultVID  = 1
	MaxVID      = 4094
	TaggedVID   = 4095
)
