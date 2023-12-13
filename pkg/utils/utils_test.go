package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/dcim"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/extras"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/tenancy"
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
				"tags": []IDObject{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
			},
			expectError: false,
		},
		{
			name: "Different tags in ClusterGroup",
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
				"tags": []IDObject{
					{ID: 1},
					{ID: 2},
					{ID: 3},
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

// TestJsonDiffMapComplex is a more complex test case
// Where nested attributes are changed and set to nil
func TestJsonDiffMapComplex(t *testing.T) {
	newObj := &virtualization.Cluster{
		Name:        "Hosting",
		Description: "New Description",
		Type: &virtualization.ClusterType{
			ID:   2,
			Name: "oVirt",
			Slug: "ovirt",
		},
		Group: &virtualization.ClusterGroup{
			ID:   4,
			Name: "New Cluster Group",
			Slug: "new-cluster-group",
		},
		Status: &dcim.Status{
			Value: "active",
			Label: "Active",
		},
		Tags: []*extras.Tag{
			{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
			{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			{ID: 4, Name: "TestX", Slug: "test2", Color: "000000", Description: "Test tag 2"},
		},
	}
	existingObj := &virtualization.Cluster{
		ID:   7,
		Name: "Hosting",
		Type: &virtualization.ClusterType{
			ID:   2,
			Name: "oVirt",
			Slug: "ovirt",
		},
		Group: &virtualization.ClusterGroup{
			ID:   3,
			Name: "Hosting",
			Slug: "hosting",
		},
		Status: &dcim.Status{
			Value: "active",
			Label: "Active",
		},
		Tenant: &tenancy.Tenant{
			ID:   1,
			Name: "Default",
			Slug: "default",
		},
		Site: &dcim.Site{
			ID:   2,
			Name: "New York",
			Slug: "new-york",
		},
		Description: "Hosting cluster",
		Tags: []*extras.Tag{
			{
				ID:    2,
				Name:  "Netbox-synced",
				Slug:  "netbox-synced",
				Color: "9e9e9e",
			},
		},
	}
	expectedDiff := map[string]interface{}{
		"description": "New Description",
		"group": IDObject{
			ID: 4,
		},
		"site": nil,
		"tags": []IDObject{
			{ID: 1},
			{ID: 3},
			{ID: 4},
		},
		"tenant": nil,
	}

	diff, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
		return
	}
	if !reflect.DeepEqual(diff, expectedDiff) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", diff, expectedDiff)
	}
}

func TestJsonDiffMapComplex2(t *testing.T) {
	newObj := &extras.CustomField{
		ID:            0,
		Name:          "New Custom field",
		Label:         "New-custom-field",
		Type:          extras.CustomFieldTypeText,
		ContentTypes:  []string{"dcim.device, virtualization.cluster"},
		SearchWeight:  1000,
		FilterLogic:   extras.FilterLogicLoose,
		UIVisibility:  extras.UIVisibilityReadWrite,
		DisplayWeight: 100,
	}
	existingObj := &extras.CustomField{
		ID:            1,
		Name:          "New Custom field",
		Label:         "New-custom-field",
		Type:          extras.CustomFieldTypeText,
		ContentTypes:  []string{"dcim.device"},
		Description:   "New custom field",
		SearchWeight:  1000,
		FilterLogic:   extras.FilterLogicLoose,
		UIVisibility:  extras.UIVisibilityReadWrite,
		DisplayWeight: 10,
	}
	expectedDiff := map[string]interface{}{
		"content_types": []string{"dcim.device, virtualization.cluster"},
		"description":   "",
		"weight":        100,
	}

	diff, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
		return
	}
	if !reflect.DeepEqual(diff, expectedDiff) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", diff, expectedDiff)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string",
			input:    "Test",
			expected: "test",
		},
		{
			name:     "String with spaces",
			input:    "Test String",
			expected: "test-string",
		},
		{
			name:     "String with trailing spaces",
			input:    "    Te st    ",
			expected: "te-st",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := Slugify(tt.input)
			if slug != tt.expected {
				t.Errorf("Slugify() = %v, want %v", slug, tt.expected)
			}
		})
	}
}

func TestNetboxMarshal(t *testing.T) {
	newObj := &virtualization.Cluster{
		Status:      &dcim.Active,
		Name:        "Test",
		Description: "Test Description",
		Type: &virtualization.ClusterType{
			ID:   2,
			Name: "oVirt",
			Slug: "ovirt",
			Tags: []*extras.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Group: &virtualization.ClusterGroup{
			ID:          4,
			Name:        "New Cluster Group",
			Description: "New cluster group",
			Slug:        "new-cluster-group",
			Tags: []*extras.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Site: &dcim.Site{
			ID:   2,
			Name: "New York",
			Slug: "new-york",
			Status: dcim.Status{
				Value: "active",
				Label: "Active",
			},
			Tags: []*extras.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Tenant: &tenancy.Tenant{
			ID:   1,
			Name: "Default",
			Slug: "default",
			Tags: []*extras.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{ID: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Tags: []*extras.Tag{
			{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
			{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			{ID: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
		},
	}
	stringRepresentation := "{\"name\":\"Test\",\"description\":\"Test Description\",\"type\":{\"id\":2},\"group\":{\"id\":4},\"status\":\"active\",\"site\":{\"id\":2},\"tenant\":{\"id\":1},\"tags\":[{\"id\":2},{\"id\":3},{\"id\":4}]}"
	expectedJson := bytes.NewBufferString(stringRepresentation)

	jsonRes, err := NetboxMarshal(newObj)
	stringRes := string(jsonRes)
	fmt.Println(stringRes)
	if err != nil {
		t.Errorf("NetboxMarshal() error = %v", err)
	}
	if !reflect.DeepEqual(jsonRes, expectedJson) {
		t.Errorf("NetboxMarshal() = %v, want %v", jsonRes, expectedJson)
	}

}
