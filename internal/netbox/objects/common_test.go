// This file contains all objects that are common to all Netbox objects.
package objects

import (
	"fmt"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func TestChoice_String(t *testing.T) {
	tests := []struct {
		name string
		c    Choice
		want string
	}{
		{
			name: "Test choice correct string",
			c: Choice{
				Value: "test value",
				Label: "test label",
			},
			want: "test value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Choice.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxObject_String(t *testing.T) {
	tests := []struct {
		name string
		n    NetboxObject
		want string
	}{
		{
			name: "Test netbox object correct string",
			n: NetboxObject{
				ID: 1,
				Tags: []*Tag{
					{Name: "Test tag1"}, {Name: "Test tag2"},
				},
				Description: "Test description",
				CustomFields: map[string]string{
					constants.CustomFieldSourceName: "test",
				},
			},
			want: fmt.Sprintf("Id: %d, Tags: %s, Description: %s", 1, []*Tag{
				{Name: "Test tag1"}, {Name: "Test tag2"},
			}, "Test description"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("NetboxObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
