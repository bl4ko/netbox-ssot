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

func (n NetboxObject) String() string {
	return fmt.Sprintf("Id: %d, Tags: %s, Description: %s", n.Id, n.Tags, n.Description)
}

type Color string

const (
	COLOR_DARK_RED    = "aa1409"
	COLOR_RED         = "f44336"
	COLOR_PINK        = "e91e63"
	COLOR_ROSE        = "ffe4e1"
	COLOR_FUCHSIA     = "ff66ff"
	COLOR_PURPLE      = "9c27b0"
	COLOR_DARK_PURPLE = "673ab7"
	COLOR_INDIGO      = "3f51b5"
	COLOR_BLUE        = "2196f3"
	COLOR_LIGHT_BLUE  = "03a9f4"
	COLOR_CYAN        = "00bcd4"
	COLOR_TEAL        = "009688"
	COLOR_AQUA        = "00ffff"
	COLOR_DARK_GREEN  = "2f6a31"
	COLOR_GREEN       = "4caf50"
	COLOR_LIGHT_GREEN = "8bc34a"
	COLOR_LIME        = "cddc39"
	COLOR_YELLOW      = "ffeb3b"
	COLOR_AMBER       = "ffc107"
	COLOR_ORANGE      = "ff9800"
	COLOR_DARK_ORANGE = "ff5722"
	COLOR_BROWN       = "795548"
	COLOR_LIGHT_GREY  = "c0c0c0"
	COLOR_GREY        = "9e9e9e"
	COLOR_DARK_GREY   = "607d8b"
	COLOR_BLACK       = "111111"
	COLOR_WHITE       = "ffffff"
)
