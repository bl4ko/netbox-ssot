// This file contains all objects that are common to all Netbox objects.
package objects

import (
	"fmt"
	"slices"
)

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
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

func (n NetboxObject) String() string {
	return fmt.Sprintf("ID: %d, Tags: %s, Description: %s", n.ID, n.Tags, n.Description)
}

func (n *NetboxObject) GetID() int {
	return n.ID
}

func (n *NetboxObject) GetCustomField(label string) interface{} {
	if n.CustomFields == nil {
		return nil
	}
	return n.CustomFields[label]
}

func (n *NetboxObject) SetCustomField(label string, value interface{}) {
	if n.CustomFields == nil {
		n.CustomFields = make(map[string]interface{})
	}
	n.CustomFields[label] = value
}

// AddTag adds a tag to the NetboxObject if
// it doesn't have it already. If the tag is already present,
// nothing happens.
func (n *NetboxObject) AddTag(newTag *Tag) {
	if slices.IndexFunc(n.Tags, func(t *Tag) bool {
		return t.Name == newTag.Name
	}) == -1 {
		n.Tags = append(n.Tags, newTag)
	}
}

// HasTag checks if the NetboxObject has a tag.
// It returns true if the object has the tag, otherwise false.
func (n *NetboxObject) HasTag(tag *Tag) bool {
	return slices.IndexFunc(n.Tags, func(t *Tag) bool {
		return t.Name == tag.Name
	}) >= 0
}

// HasTagByName checks if the NetboxObject has a tag by name.
// It returns true if the object has the tag, otherwise false.
func (n *NetboxObject) HasTagByName(tagName string) bool {
	return slices.IndexFunc(n.Tags, func(t *Tag) bool {
		return t.Name == tagName
	}) >= 0
}

// RemoveTag removes a tag from the NetboxObject.
// If the tag is not present, nothing happens.
func (n *NetboxObject) RemoveTag(tag *Tag) {
	index := slices.IndexFunc(n.Tags, func(t *Tag) bool {
		return t.Name == tag.Name
	})
	if index >= 0 {
		n.Tags = append(n.Tags[:index], n.Tags[index+1:]...)
	}
}
