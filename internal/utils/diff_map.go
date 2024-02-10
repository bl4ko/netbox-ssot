package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

// Helper function to determine if a given reflect.Value contains an embedded objects.Choice
func isChoiceEmbedded(v reflect.Value) bool {
	vType := v.Type()
	return vType.Field(0).Type == reflect.TypeOf(objects.Choice{})
	// for i := 0; i < v.NumField(); i++ {
	// 	if vType.Field(i).Type == reflect.TypeOf(objects.Choice{}) {
	// 		return true
	// 	}
	// }
	// return false
}

func choiceValue(v reflect.Value) interface{} {
	return v.Field(0).FieldByName("Value").Interface()
}

// Struct used for patching objects, when attributes are structs, or slices
// we only need object with an id of original object for patching.
type IDObject struct {
	ID int `json:"id"`
}

// JsonDiffMapExceptID compares two objects and returns a map of fields
// (represented by their JSON tag names) that are different with their
// values from newObj.
// If resetFields is set to true, the function will also include fields
// that are empty in newObj but might have a value in existingObj.
func JsonDiffMapExceptId(newObj, existingObj interface{}, resetFields bool) (map[string]interface{}, error) {
	diff := make(map[string]interface{})

	newObject := reflect.ValueOf(newObj)
	existingObject := reflect.ValueOf(existingObj)

	// Ensure that both objects are of the same kind (e.g. struct)
	if newObject.Kind() != existingObject.Kind() {
		return nil, fmt.Errorf("arguments are not of the same type")
	}

	// Check if the values are pointers and get the element they point to
	if newObject.Kind() == reflect.Ptr {
		newObject = newObject.Elem()
		existingObject = existingObject.Elem()
	}

	// Ensure that we are dealing with structs
	if newObject.Kind() != reflect.Struct {
		return nil, fmt.Errorf("arguments are not structs")
	}

	for i := 0; i < newObject.NumField(); i++ {
		fieldName := newObject.Type().Field(i).Name
		jsonTag := newObject.Type().Field(i).Tag.Get("json")

		if fieldName == "Id" {
			continue
		}

		// Custom logic for all objects that inherit from NetboxObject
		if fieldName == "NetboxObject" {
			netboxObjectDiffMap, err := JsonDiffMapExceptId(newObject.Field(i).Interface(), existingObject.Field(i).Interface(), resetFields)
			if err != nil {
				return nil, fmt.Errorf("error processing JsonDiffMapExceptId when processing NetboxObject %s", err)
			}
			for k, v := range netboxObjectDiffMap {
				diff[k] = v
			}
			continue
		}

		// Get json tag, so we can get the json field name
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = fieldName
		} else {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Ensure that both fields are of the same kind
		if newObject.Field(i).Kind() != existingObject.Field(i).Kind() {
			return nil, fmt.Errorf("field %s is not of the same type in both objects", jsonTag)
		}

		// Check if elements are pointers, in that case get the elements they are pointing to
		newObjectField := newObject.Field(i)
		existingObjectField := existingObject.Field(i)
		if newObjectField.Kind() == reflect.Ptr {
			newObjectField = newObjectField.Elem()
			existingObjectField = existingObjectField.Elem()
		}

		// If reset is set to false and the newObjectField is empty,
		// we don't do anything (like omitempty), see tests
		if !resetFields && (!newObjectField.IsValid() || newObjectField.IsZero()) {
			continue
		}

		switch newObjectField.Kind() {

		// Reset the field (when it is set to nil),
		// this only happens if flag resetFields is set to true
		case reflect.Invalid:
			if existingObjectField.IsValid() {
				diff[jsonTag] = nil
			}

		case reflect.Slice:
			err := addSliceDiff(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error processing JsonDiffMapExceptId when processing slice %s", err)
			}

		case reflect.Struct:
			err := addStructDiff(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error processing JsonDiffMapExceptId when processing struct %s", err)
			}

		case reflect.Map:
			err := addMapDiff(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error processing JsonDiffMapExceptId when processing map %s", err)
			}

		default:
			if !newObjectField.IsValid() {
				if existingObjectField.IsValid() {
					diff[jsonTag] = nil
				}
				continue
			}
			if newObjectField.Interface() != existingObjectField.Interface() {
				diff[jsonTag] = newObjectField.Interface()
			}
		}
	}

	return diff, nil
}

// Function that takes two objects (of type slice) and returns a map
// that can be easily used with json.Marshal
// To achieve this the map is of the following format:
// map[jsonTag] = [id1, id2, id3] // If the slice contains objects with ID field
// map[jsonTag] = [value1, value2, value3] // If the slice contains strings
func addSliceDiff(newSlice reflect.Value, existingSlice reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If first slice is nil, that means that we reset the value
	if !newSlice.IsValid() || newSlice.Len() == 0 {
		if existingSlice.IsValid() && existingSlice.Len() > 0 {
			diffMap[jsonTag] = []interface{}{}
		}
		return nil
	}

	// There are going to be 2 kinds of comparison.
	// One where slice will contain objects, in that case we
	// will compare ids of the objects, else we will compare
	// the values of the slice
	switch newSlice.Index(0).Kind() {
	case reflect.String:
		strSet := make(map[string]bool)
		for j := 0; j < newSlice.Len(); j++ {
			strSet[newSlice.Index(j).Interface().(string)] = true
		}
		if len(strSet) != existingSlice.Len() {
			diffMap[jsonTag] = newSlice.Interface()
		} else {
			for j := 0; j < existingSlice.Len(); j++ {
				if !strSet[existingSlice.Index(j).Interface().(string)] {
					diffMap[jsonTag] = newSlice.Interface()
					return nil
				}
			}
		}

	default:
		newIdSet := make(map[int]bool, newSlice.Len())
		var id int
		for j := 0; j < newSlice.Len(); j++ {
			element := newSlice.Index(j)
			if newSlice.Index(j).IsNil() {
				continue
			}
			if element.Kind() == reflect.Ptr {
				element = element.Elem()
			}
			id := element.FieldByName("Id").Interface().(int)
			newIdSet[id] = true
		}

		newIdSlice := make([]int, 0, len(newIdSet))
		for id := range newIdSet {
			newIdSlice = append(newIdSlice, id)
		}
		slices.Sort(newIdSlice)

		if len(newIdSet) != existingSlice.Len() {
			diffMap[jsonTag] = newIdSlice
		} else {
			for j := 0; j < existingSlice.Len(); j++ {
				existingSliceEl := existingSlice.Index(j)
				if existingSlice.Index(j).Kind() == reflect.Ptr {
					existingSliceEl = existingSliceEl.Elem()
				}
				id = existingSliceEl.FieldByName("Id").Interface().(int)
				if _, ok := newIdSet[id]; !ok {
					diffMap[jsonTag] = newIdSlice
					return nil
				}
			}
		}
	}
	return nil
}

// Returns json form for patching the difference e.g. { "Id": 1 }
func addStructDiff(newObj reflect.Value, existingObj reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If first struct is nil, that means that we reset the attribute to nil
	if !newObj.IsValid() {
		diffMap[jsonTag] = nil
		return nil
	}

	// We check if struct is a objects.Choice (special netbox struct)
	if isChoiceEmbedded(newObj) {
		if !existingObj.IsValid() || newObj.Interface() != existingObj.Interface() {
			diffMap[jsonTag] = choiceValue(newObj)
		}
		return nil
	}

	// We use ids for comparison between structs, because for patching objects, all we need is id of attribute
	idField := newObj.FieldByName("Id")

	// If objects don't have ID field, compare them by their values
	if !idField.IsValid() {
		if !existingObj.IsValid() {
			diffMap[jsonTag] = newObj.Interface()
		} else {
			if newObj.Interface() != existingObj.Interface() {
				if isChoiceEmbedded(newObj) {
					diffMap[jsonTag] = choiceValue(newObj)
				} else {
					diffMap[jsonTag] = newObj.Interface()
				}
			}
		}
	} else {
		if !existingObj.IsValid() {
			diffMap[jsonTag] = IDObject{ID: idField.Interface().(int)}
		} else {
			// Objects have ID field, compare their ids
			if newObj.FieldByName("Id").Interface() != existingObj.FieldByName("Id").Interface() {
				id := newObj.FieldByName("Id").Interface().(int)
				diffMap[jsonTag] = IDObject{ID: id}
			}
		}
	}
	return nil
}

func addMapDiff(newMap reflect.Value, existingMap reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If the new map is not set, we don't change anything
	if !newMap.IsValid() {
		return nil
	}

	// Go through all keys in new map, and check if they are in existing map
	// If they are not, add them to diff map
	mapsDiff := make(map[string]interface{})
	for _, key := range newMap.MapKeys() {
		// Keys have to be strings
		if key.Kind() != reflect.String {
			return fmt.Errorf("map keys have to be strings. Not implemented for anything else yet")
		}
		if !existingMap.MapIndex(key).IsValid() {
			mapsDiff[key.Interface().(string)] = newMap.MapIndex(key).Interface()
		} else if newMap.MapIndex(key).Interface() != existingMap.MapIndex(key).Interface() {
			mapsDiff[key.Interface().(string)] = newMap.MapIndex(key).Interface()
		}
	}

	if len(mapsDiff) > 0 {
		for _, key := range existingMap.MapKeys() {
			if !newMap.MapIndex(key).IsValid() {
				mapsDiff[key.Interface().(string)] = existingMap.MapIndex(key).Interface()
			}
		}
		diffMap[jsonTag] = mapsDiff
	}
	return nil
}

// Some more complex tests:

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
				Id:          1,
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
				Id:          1,
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
						{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{Id: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
						{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
					},
					Description: "New group",
				},
			},
			existingObj: &objects.ClusterGroup{
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: objects.NetboxObject{
					Tags: []*objects.Tag{
						{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{Id: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
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
						{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{Id: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
						{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
					},
				},
			},
			existingObj: &objects.ClusterGroup{
				Name: "New Group",
				Slug: "new-group",
				NetboxObject: objects.NetboxObject{
					Description: "New group",
					Tags: []*objects.Tag{
						{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
						{Id: 2, Name: "Test2", Slug: "test2", Color: "000000", Description: "Test tag 2"},
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
			diff, err := JsonDiffMapExceptId(tt.newObj, tt.existingObj, true)
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
				{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{Id: 4, Name: "TestX", Slug: "test2", Color: "000000", Description: "Test tag 2"},
			},
		},
	}
	existingObj := &objects.Cluster{
		NetboxObject: objects.NetboxObject{
			Id:          7,
			Description: "Hosting cluster",
			Tags: []*objects.Tag{
				{
					Id:    2,
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
		"tags": []int{1, 3, 4},
	}

	diff, err := JsonDiffMapExceptId(newObj, existingObj, true)
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
		Id:            0,
		Name:          "New Custom field",
		Label:         "New-custom-field",
		Type:          objects.CustomFieldTypeText,
		ContentTypes:  []string{"objects.device, objects.cluster"},
		SearchWeight:  1000,
		FilterLogic:   objects.FilterLogicLoose,
		DisplayWeight: 100,
	}
	existingObj := &objects.CustomField{
		Id:            1,
		Name:          "New Custom field",
		Label:         "New-custom-field",
		Type:          objects.CustomFieldTypeText,
		ContentTypes:  []string{"objects.device"},
		Description:   "New custom field",
		SearchWeight:  1000,
		FilterLogic:   objects.FilterLogicLoose,
		DisplayWeight: 10,
	}
	expectedDiff := map[string]interface{}{
		"content_types": []string{"objects.device, objects.cluster"},
		"description":   "",
		"weight":        100,
	}

	diff, err := JsonDiffMapExceptId(newObj, existingObj, true)
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
				{Id: 2, Name: "Netbox-synced"},
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

	respDiffMap, err := JsonDiffMapExceptId(newObj, existingObj, true)
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
					Id:          18,
					Name:        "Source: olvm",
					Slug:        "source-olvm",
					Color:       "07426b",
					Description: "Automatically created tag by netbox-ssot for source olvm",
				},
				{
					Id:          21,
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
					Id:          15,
					Name:        "existingTag",
					Slug:        "exiting-tag",
					Color:       "07426b",
					Description: "Automatically created tag by netbox-ssot for source olvm",
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
	gotDiffMap, err := JsonDiffMapExceptId(newInterface, existingInterface, true)
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
