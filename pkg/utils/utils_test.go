package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/objects"
)

func TestJsonDiffMapExceptIdWithMapsAddition(t *testing.T) {
	newObj := &objects.Device{
		CustomFields: map[string]string{
			"host_cpu_cores": "10 cpu cores",
			"host_mem":       "10 GB",
			"host_id":        "123456789",
		},
	}
	existingObj := &objects.Device{
		CustomFields: map[string]string{
			"host_cpu_cores": "5 cpu cores",
			"existing_tag1":  "existing_tag1",
			"existing_tag2":  "existing_tag2",
		},
	}
	expectedDiff := map[string]interface{}{
		"custom_fields": map[string]interface{}{
			"host_cpu_cores": "10 cpu cores",
			"host_mem":       "10 GB",
			"host_id":        "123456789",
			"existing_tag1":  "existing_tag1",
			"existing_tag2":  "existing_tag2",
		},
	}
	receivedDiff, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	if !reflect.DeepEqual(receivedDiff, expectedDiff) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", receivedDiff, expectedDiff)
	}
	// We need to ensure that the diff is a valid JSON
	_, err = json.Marshal(receivedDiff)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
}

func TestJsonDiffMapExceptIdWithMapsNoAddition(t *testing.T) {
	newObj := &objects.Device{
		CustomFields: map[string]string{
			"host_cpu_cores": "10 cpu cores",
			"host_mem":       "10 GB",
		},
	}
	existingObj := &objects.Device{
		CustomFields: map[string]string{
			"host_cpu_cores": "10 cpu cores",
			"host_mem":       "10 GB",
			"existing_tag1":  "existing_tag1",
			"existing_tag2":  "existing_tag2",
		},
	}
	// We expect no difference, because all new fields are already present in the attribute
	expectedDiff := map[string]interface{}{}
	receivedDiff, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	if !reflect.DeepEqual(receivedDiff, expectedDiff) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", receivedDiff, expectedDiff)
	}
	// We need to ensure that the diff is a valid JSON
	_, err = json.Marshal(receivedDiff)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
}

func TestJsonDiffMapExceptIdWithMapsEmpty(t *testing.T) {
	newObj := &objects.Device{}
	existingObj := &objects.Device{
		CustomFields: map[string]string{
			"host_cpu_cores": "10 cpu cores",
			"host_mem":       "10 GB",
			"existing_tag1":  "existing_tag1",
			"existing_tag2":  "existing_tag2",
		},
	}
	// We expect no difference, because all new fields are already present in the attribute
	expectedDiff := map[string]interface{}{}
	receivedDiff, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	if !reflect.DeepEqual(receivedDiff, expectedDiff) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", receivedDiff, expectedDiff)
	}
	// We need to ensure that the diff is a valid JSON
	_, err = json.Marshal(receivedDiff)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
}

// TODO: determine how to reset attribute maps
// func TestJsonDiffMapExceptIdWithMapsReset(t *testing.T) {
// 	newObj := &objects.Device{}
// 	existingObj := &objects.Device{
// 		CustomFields: map[string]string{},
// 	}
// 	// We expect no difference, because all new fields are already present in the attribute
// 	expectedDiff := map[string]interface{}{}
// 	receivedDiff, err := JsonDiffMapExceptId(newObj, existingObj)
// 	if err != nil {
// 		t.Errorf("JsonDiffMapExceptId() error = %v", err)
// 	}
// 	if !reflect.DeepEqual(receivedDiff, expectedDiff) {
// 		t.Errorf("JsonDiffMapExceptId() = %v, want %v", receivedDiff, expectedDiff)
// 	}
// 	// We need to ensure that the diff is a valid JSON
// 	_, err = json.Marshal(receivedDiff)
// 	if err != nil {
// 		t.Errorf("JsonDiffMapExceptId() error = %v", err)
// 	}
// }

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
			newObj: &objects.Tag{
				Name:        "Test",
				Slug:        "test",
				Color:       "000000",
				Description: "Test tag",
			},
			existingObj: &objects.Tag{
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
			newObj: &objects.Tag{
				Name:        "Test Changed",
				Slug:        "test-changed",
				Color:       "000000",
				Description: "Changed tag",
			},
			existingObj: &objects.Tag{
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
			newObj: &objects.ClusterGroup{
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: objects.NetboxObject{
					Tags: []*objects.Tag{
						{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
						{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
					},
					Description: "New group",
				},
			},
			existingObj: &objects.ClusterGroup{
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: objects.NetboxObject{
					Tags: []*objects.Tag{
						{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					},
					Description: "New group",
				},
			},
			expected: map[string]interface{}{
				"tags": []int{1, 2, 3},
			},
			expectError: false,
		},
		{
			name: "Different tags in ClusterGroup",
			newObj: &objects.ClusterGroup{
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: objects.NetboxObject{
					Description: "New group",
					Tags: []*objects.Tag{
						{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
						{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
					},
				},
			},
			existingObj: &objects.ClusterGroup{
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: objects.NetboxObject{
					Description: "New group",
					Tags: []*objects.Tag{
						{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{ID: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
					},
				},
			},
			expected: map[string]interface{}{
				"tags": []int{1, 2, 3},
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
			// Also ensure that the diff is a valid JSON
			_, err = json.Marshal(diff)
			if err != nil {
				t.Errorf("JsonDiffMapExceptId() error = %v", err)
			}
		})
	}
}

// TestJsonDiffMapComplex is a more complex test case
// Where nested attributes are changed and set to nil
func TestJsonDiffMapComplex(t *testing.T) {
	newObj := &objects.Cluster{
		Name: "Hosting",
		Type: &objects.ClusterType{
			NetboxObject: objects.NetboxObject{Id: 2},
			Name:         "oVirt",
			Slug:         "ovirt",
		},
		Group: &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{Id: 4},
			Name:         "New Cluster Group",
			Slug:         "new-cluster-group",
		},
		Status: objects.ClusterStatusActive,
		NetboxObject: objects.NetboxObject{
			Description: "New Description",
			Tags: []*objects.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{ID: 4, Name: "TestX", Slug: "test2", Color: "000000", Description: "Test tag 2"},
			},
		},
	}
	existingObj := &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			Id:          7,
			Description: "Hosting cluster",
			Tags: []*objects.Tag{
				{
					ID:    2,
					Name:  "Netbox-synced",
					Slug:  "netbox-synced",
					Color: "9e9e9e",
				},
			},
		},
		Name: "Hosting",
		Type: &objects.ClusterType{
			NetboxObject: objects.NetboxObject{Id: 2},
			Name:         "oVirt",
			Slug:         "ovirt",
		},
		Group: &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{Id: 3},
			Name:         "Hosting",
			Slug:         "hosting",
		},
		Status: objects.ClusterStatusActive,
		Tenant: &objects.Tenant{
			NetboxObject: objects.NetboxObject{Id: 1},
			Name:         "Default",
			Slug:         "default",
		},
		Site: &objects.Site{
			NetboxObject: objects.NetboxObject{Id: 2},
			Name:         "New York",
			Slug:         "new-york",
		},
	}
	expectedDiff := map[string]interface{}{
		"description": "New Description",
		"group": IDObject{
			ID: 4,
		},
		"site":   nil,
		"tags":   []int{1, 3, 4},
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
	// Also ensure that the diff is a valid JSON
	_, err = json.Marshal(diff)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
}

func TestJsonDiffMapComplex2(t *testing.T) {
	newObj := &objects.CustomField{
		ID:            0,
		Name:          "New Custom field",
		Label:         "New-custom-field",
		Type:          objects.CustomFieldTypeText,
		ContentTypes:  []string{"objects.device, objects.cluster"},
		SearchWeight:  1000,
		FilterLogic:   objects.FilterLogicLoose,
		UIVisibility:  objects.UIVisibilityReadWrite,
		DisplayWeight: 100,
	}
	existingObj := &objects.CustomField{
		ID:            1,
		Name:          "New Custom field",
		Label:         "New-custom-field",
		Type:          objects.CustomFieldTypeText,
		ContentTypes:  []string{"objects.device"},
		Description:   "New custom field",
		SearchWeight:  1000,
		FilterLogic:   objects.FilterLogicLoose,
		UIVisibility:  objects.UIVisibilityReadWrite,
		DisplayWeight: 10,
	}
	expectedDiff := map[string]interface{}{
		"content_types": []string{"dcim.device, dcim.cluster"},
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
	// Also ensure that the diff is a valid JSON
	_, err = json.Marshal(diff)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
}

func TestJsonDiffMapWithChoiceAttr(t *testing.T) {
	newObj := &objects.Device{
		Name:   "Test device",
		Status: &objects.DeviceStatusActive,
	}

	existingObj := &objects.Device{
		Name: "Test device",
		NetboxObject: objects.NetboxObject{
			Id:          1,
			Description: "Test device",
			Tags: []*objects.Tag{
				{ID: 2, Name: "Netbox-synced"},
			},
		},
		Status: &objects.DeviceStatusOffline,
	}
	expectedDiffMap := map[string]interface{}{
		"description": "",
		"tags":        []interface{}{},
		"status":      "active",
	}

	respDiffMap, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	if !reflect.DeepEqual(respDiffMap, expectedDiffMap) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", respDiffMap, expectedDiffMap)
	}
	// Also ensure that the diff is a valid JSON
	_, err = json.Marshal(respDiffMap)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
}

func TestJsonDiffMapWithMapAttr(t *testing.T) {
	newObj := &objects.Device{
		Name: "Test device",
		DeviceRole: &objects.DeviceRole{
			NetboxObject: objects.NetboxObject{
				Id: 1,
			},
			Name:  "Test device role",
			Slug:  "test-device-role",
			Color: "000000",
		},
		DeviceType: &objects.DeviceType{
			NetboxObject: objects.NetboxObject{
				Id: 1,
			},
			Model: "Test device model",
			Slug:  "test-device-type",
		},
		Status: &objects.DeviceStatusActive,
		CustomFields: map[string]string{
			"host_cpu_cores": "10 cpu cores",
			"host_mem":       "10 GB",
		},
	}

	existingObj := &objects.Device{
		NetboxObject: objects.NetboxObject{
			Id:          1,
			Description: "Test device",
			Tags: []*objects.Tag{
				{ID: 2, Name: "Netbox-synced"},
			},
		},
		Name: "Test device",
		DeviceRole: &objects.DeviceRole{
			NetboxObject: objects.NetboxObject{Id: 1},
			Name:         "Test device role",
			Slug:         "test-device-role",
			Color:        "000000",
		},
		DeviceType: &objects.DeviceType{
			NetboxObject: objects.NetboxObject{Id: 1},
			Model:        "test-model",
			Slug:         "test-device-type",
		},
		Status: &objects.DeviceStatusOffline,
		CustomFields: map[string]string{
			"host_cpu_cores":  "10 cpu cores",
			"host_mem":        "10 GB",
			"extra_from":      "before",
			"don't remove me": "please",
		},
	}
	expectedDiffMap := map[string]interface{}{
		"description": "",
		"tags":        []interface{}{},
		"status":      "active",
	}

	respDiffMap, err := JsonDiffMapExceptId(newObj, existingObj)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	if !reflect.DeepEqual(respDiffMap, expectedDiffMap) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", respDiffMap, expectedDiffMap)
	}
	// Also ensure that the diff is a valid JSON
	_, err = json.Marshal(respDiffMap)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}

}

func TestJsonDiffMapWithMapAttr2(t *testing.T) {
	newInterface := &objects.Interface{
		NetboxObject: objects.NetboxObject{
			Tags: []*objects.Tag{
				{
					ID:          18,
					Name:        "Source: olvm",
					Slug:        "source-olvm",
					Color:       "07426b",
					Description: "Automatically created tag by netbox-ssot for source srcolvm",
				},
				{
					ID:          21,
					Name:        "ovirt",
					Slug:        "type-ovirt",
					Color:       "ff0000",
					Description: "Automatically created tag by netbox-ssot for source type ovirt",
				},
			},
		},
		Name:   "enp4s0f4",
		Status: true,
		CustomFields: map[string]string{
			"source_id": "abcdefghijkl",
			"extra_one": "extra_one",
		},
	}
	existingInterface := &objects.Interface{
		NetboxObject: objects.NetboxObject{
			Tags: []*objects.Tag{
				{
					ID:          15,
					Name:        "existingTag",
					Slug:        "exiting-tag",
					Color:       "07426b",
					Description: "Automatically created tag by netbox-ssot for source srcolvm",
				},
			},
		},
		Name:   "enp4s0f4",
		Status: true,
	}

	expectedDiffMap := map[string]interface{}{
		"tags": []int{18, 21},
		"custom_fields": map[string]interface{}{
			"source_id": "abcdefghijkl",
			"extra_one": "extra_one",
		},
	}
	gotDiffMap, err := JsonDiffMapExceptId(newInterface, existingInterface)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	// try to marshal gotDiffMap
	_, err = json.Marshal(gotDiffMap)
	if err != nil {
		t.Errorf("JsonDiffMapExceptId() error = %v", err)
	}
	if !reflect.DeepEqual(gotDiffMap, expectedDiffMap) {
		t.Errorf("JsonDiffMapExceptId() = %v, want %v", gotDiffMap, expectedDiffMap)
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
	newObj := &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			Description: "Test Description",
			Tags: []*objects.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{ID: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Status: objects.ClusterStatusActive,
		Name:   "Test",
		Type: &objects.ClusterType{
			NetboxObject: objects.NetboxObject{
				Id: 2,
				Tags: []*objects.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			Name: "oVirt",
			Slug: "ovirt",
		},
		Group: &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{
				Id: 4,
				Tags: []*objects.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
				Description: "New cluster group",
			},
			Name: "New Cluster Group",
			Slug: "new-cluster-group",
		},
		Site: &objects.Site{
			NetboxObject: objects.NetboxObject{
				Id: 2,
				Tags: []*objects.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			Name:   "New York",
			Slug:   "new-york",
			Status: &objects.SiteStatusActive,
		},
		Tenant: &objects.Tenant{
			NetboxObject: objects.NetboxObject{
				Id: 1,
				Tags: []*objects.Tag{
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
	device := &objects.Device{
		NetboxObject: objects.NetboxObject{
			Description: "Test Description",
			Tags: []*objects.Tag{
				{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{ID: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{ID: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Name: "Test device",
		DeviceRole: &objects.DeviceRole{
			NetboxObject: objects.NetboxObject{
				Id: 1,
				Tags: []*objects.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				},
				Description: "Test device role",
			},
			Name:  "Test device role",
			Slug:  "test-device-role",
			Color: "000000",
		},
		DeviceType: &objects.DeviceType{
			NetboxObject: objects.NetboxObject{
				Id: 1,
				Tags: []*objects.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				},
				Description: "Test device type",
			},
		},
		Airflow: &objects.FrontToRear,
		Status:  &objects.DeviceStatusActive,
		Site: &objects.Site{
			NetboxObject: objects.NetboxObject{
				Id:          1,
				Description: "Test site",
				Tags: []*objects.Tag{
					{ID: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				},
			},
			Name:   "Test site",
			Slug:   "test-site",
			Status: &objects.SiteStatusActive,
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
