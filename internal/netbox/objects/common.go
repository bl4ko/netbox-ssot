// This file contains all objects that are common to all Netbox objects.
package objects

import "fmt"

// Choice represents a choice in a Netbox's choice field.
// This struct is used as an embedded struct in other structs that represent Choice fields.
type Choice struct {
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
}

func (c Choice) String() string {
	return c.Value
}

const (
	MaxDescriptionLength = 200
)

// Struct representing attributes that are common to all objects in Netbox.
// We can this struct as an embedded struct in other structs that represent
// Netbox objects.
type NetboxObject struct {
	// Netbox's ID of the object.
	ID int `json:"id,omitempty"`
	// List of tags assigned to this object.
	Tags []*Tag `json:"tags,omitempty"`
	// Description represents custom description of the object.
	Description string `json:"description,omitempty"`
	// Array of custom fields, in format customFieldLabel: customFieldValue
	CustomFields map[string]string `json:"custom_fields,omitempty"`
}

func (n NetboxObject) String() string {
	return fmt.Sprintf("Id: %d, Tags: %s, Description: %s", n.ID, n.Tags, n.Description)
}
