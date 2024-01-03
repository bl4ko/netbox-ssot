// This file contains all objects that are common to all Netbox objects.
package objects

// Choice represents a choice in a Netbox's choice field.
// This struct is used as an embedded struct in other structs that represent Choice fields.
type Choice struct {
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
}

// Struct representing attributes that are common to all objects in Netbox.
// We can this struct as an embedded struct in other structs that represent
// Netbox objects.
type NetboxObject struct {
	// Netbox's Id of the object.
	Id int `json:"id,omitempty"`
	// List of tags assigned to this object.
	Tags []*Tag `json:"tags,omitempty"`
	// Description represents custom description of the object.
	Description string `json:"description,omitempty"`
}
