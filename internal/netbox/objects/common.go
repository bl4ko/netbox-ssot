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

type Color string

const (
	ColorDarkRed    = "aa1409"
	ColorRed        = "f44336"
	ColorPink       = "e91e63"
	ColorRose       = "ffe4e1"
	ColorFuchsia    = "ff66ff"
	ColorPurple     = "9c27b0"
	ColorDarkPurple = "673ab7"
	ColorIndigo     = "3f51b5"
	ColorBlue       = "2196f3"
	ColorLightBlue  = "03a9f4"
	ColorCyan       = "00bcd4"
	ColorTeal       = "009688"
	ColorAqua       = "00ffff"
	ColorDarkGreen  = "2f6a31"
	ColorGreen      = "4caf50"
	ColorLightGreen = "8bc34a"
	ColorLime       = "cddc39"
	ColorYellow     = "ffeb3b"
	ColorAmber      = "ffc107"
	ColorOrange     = "ff9800"
	ColorDarkOrange = "ff5722"
	ColorBrown      = "795548"
	ColorLightGrey  = "c0c0c0"
	ColorGrey       = "9e9e9e"
	ColorDarkGrey   = "607d8b"
	ColorBlack      = "111111"
	ColorWhite      = "ffffff"
)
