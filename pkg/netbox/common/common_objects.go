// This module contains all objects that are shared between multiple modules.
// With this apporach we can avoid circular imports.
// For example, we have virtualization.cluster that has a field common.Site
// On the other hand dcim.Device has a field virtualization.Cluster.

package common

import (
	"fmt"
)

type Tag struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

func (t Tag) String() string {
	return fmt.Sprintf("Tag{ID: %d, Name: %s, Slug: %s, Color: %s, Description: %s}", t.ID, t.Name, t.Slug, t.Color, t.Description)
}

// Choice represents a choice in a Netbox's choice field.
// This struct is used as an embedded struct in other structs that represent Choice fields.
type Choice struct {
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
}

// Struct representing attributes that are common to all objects in NetBox.
// We can this struct as an embedded struct in other structs that represent
// NetBox objects.
type NetboxObject struct {
	// Netbox's ID of the object.
	ID int `json:"id,omitempty"`
	// List of tags assigned to this object.
	Tags []*Tag `json:"tags,omitempty"`
	// Description represents custom description of the object.
	Description string `json:"description,omitempty"`
}

type SiteStatus struct {
	Choice
}

var (
	SiteStatusActive         = SiteStatus{Choice{Value: "active", Label: "Active"}}
	SiteStatusOffline        = SiteStatus{Choice{Value: "offline", Label: "Offline"}}
	SiteStatusPlanned        = SiteStatus{Choice{Value: "planned", Label: "Planned"}}
	SiteStatusStaged         = SiteStatus{Choice{Value: "staged", Label: "Staged"}}
	SiteStatusFailed         = SiteStatus{Choice{Value: "failed", Label: "Failed"}}
	SiteStatusInventory      = SiteStatus{Choice{Value: "inventory", Label: "Inventory"}}
	SiteStatucDecommisioning = SiteStatus{Choice{Value: "decommissioning", Label: "Decommissioning"}}
)

// Site ares used for functional groupings.
// A site usually represents a building within a region.
type Site struct {
	NetboxObject
	// Full name for the site. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Site status. This field is required.
	Status *SiteStatus `json:"status,omitempty"`
}

// Manufacturer represents a hardware manufacturer (e.g. Cisco, HP, ...).
type Manufacturer struct {
	NetboxObject
	// Name of the manufacturer (e.g. Cisco). This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
}

var ManufacturerMap = map[string]string{
	"^AMD$":               "AMD",
	".*Broadcom.*":        "Broadcom",
	".*Cisco.*":           "Cisco",
	".*Dell.*":            "Dell",
	"FTS Corp":            "Fujitsu",
	".*Fujitsu.*":         "Fujitsu",
	".*HiSilicon.*":       "HiSilicon",
	"^HP$":                "HPE",
	"^HPE$":               "HPE",
	".*Huawei.*":          "Huawei",
	".*Hynix.*":           "Hynix",
	".*Inspur.*":          "Inspur",
	".*Intel.*":           "Intel",
	"LEN":                 "Lenovo",
	".*Lenovo.*":          "Lenovo",
	".*Micron.*":          "Micron",
	".*Nvidea.*":          "Nvidia",
	".*Samsung.*":         "Samsung",
	".*Supermicro.*":      "Supermicro",
	".*Toshiba.*":         "Toshiba",
	"^WD$":                "Western Digital",
	".*Western Digital.*": "Western Digital",
}

// Platform represents an operating system or other software platform which may be running on a device.
type Platform struct {
	NetboxObject
	// Name of the platform. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Manufacturer is the manufacturer of the platform.
	Manafacturer *Manufacturer `json:"manufacturer,omitempty"`
}
