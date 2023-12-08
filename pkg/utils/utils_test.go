package utils

import (
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/extras"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/virtualization"
)

// Assuming your structs and JsonDiffMapExceptId function are in the same package
// If not, import the package where they are defined

func TestJsonDiffMapExceptId(t *testing.T) {
	tests := []struct {
		name        string
		newObj      interface{}
		existingObj interface{}
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "No difference in Tag",
			newObj: &extras.Tag{
				Name:        "Test",
				Slug:        "test",
				Color:       "000000",
				Description: "Test tag",
			},
			existingObj: &extras.Tag{
				ID:          1,
				Name:        "Test",
				Slug:        "test",
				Color:       "000000",
				Description: "Test tag",
			},
			expected:    map[string]interface{}{},
			expectError: false,
		},
		{
			name: "Different fields in Tag",
			newObj: &extras.Tag{
				Name:        "Test Changed",
				Slug:        "test-changed",
				Color:       "000000",
				Description: "Changed tag",
			},
			existingObj: &extras.Tag{
				ID:          1,
				Name:        "Test",
				Slug:        "test",
				Color:       "000000",
				Description: "Test tag",
			},
			expected: map[string]interface{}{
				"name":        "Test Changed",
				"slug":        "test-changed",
				"description": "Changed tag",
			},
			expectError: false,
		},
		{
			name: "Different number of Tags in ClusterGroup",
			newObj: &virtualization.ClusterGroup{
				Name:        "New Group",
				Slug:        "new-group",
				Description: "New group",
				Tags: []*extras.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			existingObj: &virtualization.ClusterGroup{
				Name:        "New Group",
				Slug:        "new-group",
				Description: "New group",
				Tags: []*extras.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
				},
			},
			expected: map[string]interface{}{
				"tags": []*extras.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			expectError: false,
		},
		{
			name: "Different Tags in ClusterGroup",
			newObj: &virtualization.ClusterGroup{
				Name:        "New Group",
				Slug:        "new-group",
				Description: "New group",
				Tags: []*extras.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 4, Name: "TestX", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			existingObj: &virtualization.ClusterGroup{
				Name:        "New Group",
				Slug:        "new-group",
				Description: "New group",
				Tags: []*extras.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			expected: map[string]interface{}{
				"tags": []*extras.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 4, Name: "TestX", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := JsonDiffMapExceptId(tt.newObj, tt.existingObj)
			if (err != nil) != tt.expectError {
				t.Errorf("JsonDiffMapExceptId() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !reflect.DeepEqual(diff, tt.expected) {
				t.Errorf("JsonDiffMapExceptId() = %v, want %v", diff, tt.expected)
			}
		})
	}
}
