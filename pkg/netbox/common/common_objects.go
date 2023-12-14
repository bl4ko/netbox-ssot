// This module contains all objects that are shared between multiple modules.
// With this apporach we can avoid circular imports.
// For example, we have virtualization.cluster that has a field common.Site
// On the other hand dcim.Device has a field virtualization.Cluster.

package common

import "fmt"

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
type Choice struct {
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
}

// Struct representing attributes that are common to all objects in NetBox.
// We can use this struct as an embedded struct in other structs that represent
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
	StatusActive          = SiteStatus{Choice{Value: "active", Label: "Active"}}
	StatusOffline         = SiteStatus{Choice{Value: "offline", Label: "Offline"}}
	StatusPlanned         = SiteStatus{Choice{Value: "planned", Label: "Planned"}}
	StatusStaged          = SiteStatus{Choice{Value: "staged", Label: "Staged"}}
	StatusFailed          = SiteStatus{Choice{Value: "failed", Label: "Failed"}}
	StatusInventory       = SiteStatus{Choice{Value: "inventory", Label: "Inventory"}}
	StatusDecommissioning = SiteStatus{Choice{Value: "decommissioning", Label: "Decommissioning"}}
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
	Status SiteStatus `json:"status,omitempty"`
}

type Platform struct {
	NetboxObject
	// Name of the platform. This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
	// Manufacturer is the manufacturer of the platform.
	Manafacturer *Manafacturer `json:"manufacturer,omitempty"`
}

type Manafacturer struct {
	NetboxObject
	// Name of the manufacturer (e.g. Cisco). This field is required.
	Name string `json:"name,omitempty"`
	// URL-friendly unique shorthand. This field is required.
	Slug string `json:"slug,omitempty"`
}
