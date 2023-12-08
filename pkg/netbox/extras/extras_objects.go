package extras

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
