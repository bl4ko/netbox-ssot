package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/common"
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
			newObj: &common.Tag{
				Name:        "Test",
				Slug:        "test",
				Color:       "000000",
				Description: "Test tag",
			},
			existingObj: &common.Tag{
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
			newObj: &common.Tag{
				Name:        "Test Changed",
				Slug:        "test-changed",
				Color:       "000000",
				Description: "Changed tag",
			},
			existingObj: &common.Tag{
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
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: common.NetboxObject{
					Tags: []*common.Tag{
						{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
						{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
					},
					Description: "New group",
				},
			},
			existingObj: &virtualization.ClusterGroup{
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: common.NetboxObject{
					Tags: []*common.Tag{
						{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					},
					Description: "New group",
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
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: common.NetboxObject{
					Description: "New group",
					Tags: []*common.Tag{
						{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
						{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
					},
				},
			},
			existingObj: &virtualization.ClusterGroup{
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: common.NetboxObject{
					Description: "New group",
					Tags: []*common.Tag{
						{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					},
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
		Name: "Hosting",
		Type: &virtualization.ClusterType{
			NetboxObject: common.NetboxObject{ID: 2},
			Name:         "oVirt",
			Slug:         "ovirt",
		},
		Group: &virtualization.ClusterGroup{
			NetboxObject: common.NetboxObject{ID: 4},
			Name:         "New Cluster Group",
			Slug:         "new-cluster-group",
		},
		Status: virtualization.ClusterStatusActive,
		NetboxObject: common.NetboxObject{
			Description: "New Description",
			Tags: []*common.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{ID: 4, Name: "TestX", Slug: "test2", Color: "000000", Description: "Test tag 2"},
			},
		},
	}
	existingObj := &virtualization.Cluster{
		NetboxObject: common.NetboxObject{
			ID:          7,
			Description: "Hosting cluster",
			Tags: []*common.Tag{
				{
					ID:    2,
					Name:  "Netbox-synced",
					Slug:  "netbox-synced",
					Color: "9e9e9e",
				},
			},
		},
		Name: "Hosting",
		Type: &virtualization.ClusterType{
			NetboxObject: common.NetboxObject{ID: 2},
			Name:         "oVirt",
			Slug:         "ovirt",
		},
		Group: &virtualization.ClusterGroup{
			NetboxObject: common.NetboxObject{ID: 3},
			Name:         "Hosting",
			Slug:         "hosting",
		},
		Status: virtualization.ClusterStatusActive,
		Tenant: &tenancy.Tenant{
			NetboxObject: common.NetboxObject{ID: 1},
			Name:         "Default",
			Slug:         "default",
		},
		Site: &common.Site{
			NetboxObject: common.NetboxObject{ID: 2},
			Name:         "New York",
			Slug:         "new-york",
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

func TestJsonDiffMapWithChoiceAttr(t *testing.T) {
	newObj := &dcim.Device{
		Name: "Test device",
		DeviceRole: &dcim.DeviceRole{
			NetboxObject: common.NetboxObject{
				ID: 1,
			},
			Name:  "Test device role",
			Slug:  "test-device-role",
			Color: "000000",
		},
		DeviceType: &dcim.DeviceType{
			NetboxObject: common.NetboxObject{
				ID: 1,
			},
			Model: "Test device model",
			Slug:  "test-device-type",
		},
		Status: &dcim.DeviceStatusActive,
	}

	existingObj := &dcim.Device{
		NetboxObject: common.NetboxObject{
			ID:          1,
			Description: "Test device",
			Tags: []*common.Tag{
				{ID: 2, Name: "Netbox-synced"},
			},
		},
		Name: "Test device",
		DeviceRole: &dcim.DeviceRole{
			NetboxObject: common.NetboxObject{ID: 1},
			Name:         "Test device role",
			Slug:         "test-device-role",
			Color:        "000000",
		},
		DeviceType: &dcim.DeviceType{
			NetboxObject: common.NetboxObject{ID: 1},
			Model:        "test-model",
			Slug:         "test-device-type",
		},
		Status: &dcim.DeviceStatusOffline,
	}
	expectedDiffMap := map[string]interface{}{
		"description": "",
		"tags":        nil,
		"status":      "active",
	}

	respDiffMap, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	if !reflect.DeepEqual(respDiffMap, expectedDiffMap) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", respDiffMap, expectedDiffMap)
	}

}

func TestJsonDiffMapWithMapAttr(t *testing.T) {
	newObj := &dcim.Device{
		Name: "Test device",
		DeviceRole: &dcim.DeviceRole{
			NetboxObject: common.NetboxObject{
				ID: 1,
			},
			Name:  "Test device role",
			Slug:  "test-device-role",
			Color: "000000",
		},
		DeviceType: &dcim.DeviceType{
			NetboxObject: common.NetboxObject{
				ID: 1,
			},
			Model: "Test device model",
			Slug:  "test-device-type",
		},
		Status: &dcim.DeviceStatusActive,
		CustomFields: map[string]string{
			"host_cpu_cores": "10 cpu cores",
			"host_mem":       "10 GB",
		},
	}

	existingObj := &dcim.Device{
		NetboxObject: common.NetboxObject{
			ID:          1,
			Description: "Test device",
			Tags: []*common.Tag{
				{ID: 2, Name: "Netbox-synced"},
			},
		},
		Name: "Test device",
		DeviceRole: &dcim.DeviceRole{
			NetboxObject: common.NetboxObject{ID: 1},
			Name:         "Test device role",
			Slug:         "test-device-role",
			Color:        "000000",
		},
		DeviceType: &dcim.DeviceType{
			NetboxObject: common.NetboxObject{ID: 1},
			Model:        "test-model",
			Slug:         "test-device-type",
		},
		Status: &dcim.DeviceStatusOffline,
		CustomFields: map[string]string{
			"host_cpu_cores":  "10 cpu cores",
			"host_mem":        "10 GB",
			"extra_from":      "before",
			"don't remove me": "please",
		},
	}
	expectedDiffMap := map[string]interface{}{
		"description": "",
		"tags":        nil,
		"status":      "active",
	}

	respDiffMap, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	if !reflect.DeepEqual(respDiffMap, expectedDiffMap) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", respDiffMap, expectedDiffMap)
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
			expected: "test_string",
		},
		{
			name:     "String with trailing spaces",
			input:    "    Te st    ",
			expected: "te_st",
		},
		{
			name:     "String with special characters",
			input:    "Test@#String$%^",
			expected: "teststring",
		},
		{
			name:     "String with mixed case letters",
			input:    "TeSt StRiNg",
			expected: "test_string",
		},
		{
			name:     "String with numbers",
			input:    "Test123 String456",
			expected: "test123_string456",
		},
		{
			name:     "String with underscores",
			input:    "Test_String",
			expected: "test_string",
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
		NetboxObject: common.NetboxObject{
			Description: "Test Description",
			Tags: []*common.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{ID: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Status: virtualization.ClusterStatusActive,
		Name:   "Test",
		Type: &virtualization.ClusterType{
			NetboxObject: common.NetboxObject{
				ID: 2,
				Tags: []*common.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			Name: "oVirt",
			Slug: "ovirt",
		},
		Group: &virtualization.ClusterGroup{
			NetboxObject: common.NetboxObject{
				ID: 4,
				Tags: []*common.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
				Description: "New cluster group",
			},
			Name: "New Cluster Group",
			Slug: "new-cluster-group",
		},
		Site: &common.Site{
			NetboxObject: common.NetboxObject{
				ID: 2,
				Tags: []*common.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			Name:   "New York",
			Slug:   "new-york",
			Status: &common.SiteStatusActive,
		},
		Tenant: &tenancy.Tenant{
			NetboxObject: common.NetboxObject{
				ID: 1,
				Tags: []*common.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
					{ID: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			Name: "Default",
			Slug: "default",
		},
	}
	expectedJsonMap := map[string]interface{}{
		"description": "Test Description",
		"tags":        []int{1, 3, 4},
		"name":        "Test",
		"type":        2,
		"status":      "active",
		"site":        2,
		"group":       4,
		"tenant":      1,
	}
	expectedJson, _ := json.Marshal(expectedJsonMap)
	responseJson, err := NetboxJsonMarshal(newObj)
	if err != nil {
		t.Errorf("NetboxMarshal() error = %v", err)
	}
	if !reflect.DeepEqual(responseJson, expectedJson) {
		t.Errorf("NetboxMarshal() = %s\nwant %s", string(responseJson), string(expectedJson))
	}
}

func TestNetboxJsonMarshalWithChoiceAttr(t *testing.T) {
	device := &dcim.Device{
		NetboxObject: common.NetboxObject{
			Description: "Test Description",
			Tags: []*common.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{ID: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Name: "Test device",
		DeviceRole: &dcim.DeviceRole{
			NetboxObject: common.NetboxObject{
				ID: 1,
				Tags: []*common.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				},
				Description: "Test device role",
			},
			Name:  "Test device role",
			Slug:  "test-device-role",
			Color: "000000",
		},
		DeviceType: &dcim.DeviceType{
			NetboxObject: common.NetboxObject{
				ID: 1,
				Tags: []*common.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				},
				Description: "Test device type",
			},
		},
		Airflow: &dcim.FrontToRear,
		Status:  &dcim.DeviceStatusActive,
		Site: &common.Site{
			NetboxObject: common.NetboxObject{
				ID:          1,
				Description: "Test site",
				Tags: []*common.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				},
			},
			Name:   "Test site",
			Slug:   "test-site",
			Status: &common.SiteStatusActive,
		},
	}

	expectedMap := map[string]interface{}{
		"description": "Test Description",
		"tags":        []int{1, 3, 4},
		"name":        "Test device",
		"role":        1,
		"device_type": 1,
		"airflow":     "front-to-rear",
		"status":      "active",
		"site":        1,
	}
	expectedJson, _ := json.Marshal(expectedMap)
	responseJson, err := NetboxJsonMarshal(device)
	fmt.Println(string(responseJson))
	if err != nil {
		t.Errorf("NetboxMarshal() error = %v", err)
	}
	if !reflect.DeepEqual(expectedJson, responseJson) {
		t.Errorf("NetboxMarshal() = %s\nwant %s", string(responseJson), string(expectedJson))
	}
}
