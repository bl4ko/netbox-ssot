package utils

import (
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestPrimaryAttributesDiff(t *testing.T) {
	tests := []struct {
		name           string
		newStruct      interface{}
		existingStruct interface{}
		resetFields    bool
		expectedDiff   map[string]interface{}
	}{
		{
			name:        "Addition with resetFields=true",
			resetFields: true,
			newStruct: &objects.Contact{
				Name:  "New Contact",
				Email: "newcontact@example.com",
			},
			existingStruct: &objects.Contact{
				Name:  "Existing Contact",
				Phone: "123456789",
			},
			expectedDiff: map[string]interface{}{
				"name":  "New Contact",
				"email": "newcontact@example.com",
				"phone": "",
			},
		},
		{
			name:        "Addition with resetFields=false",
			resetFields: false,
			newStruct: &objects.Contact{
				Name:  "New Contact",
				Email: "newcontact@example.com",
			},
			existingStruct: &objects.Contact{
				Name:  "Existing Contact",
				Phone: "123456789",
			},
			expectedDiff: map[string]interface{}{
				"name":  "New Contact",
				"email": "newcontact@example.com",
			},
		},
		{
			name:        "NoAddition with resetFields=true",
			resetFields: true,
			newStruct: &objects.Contact{
				Name:  "Existing Contact",
				Phone: "123456789",
			},
			existingStruct: &objects.Contact{
				Name:  "Existing Contact",
				Email: "newcontact@example.com",
				Phone: "123456789",
			},
			expectedDiff: map[string]interface{}{
				"email": "",
			},
		},
		{
			name:        "NoAddition with resetFields=false",
			resetFields: false,
			newStruct: &objects.Contact{
				Name:  "Existing Contact",
				Phone: "123456789",
			},
			existingStruct: &objects.Contact{
				Name:  "Existing Contact",
				Email: "newcontact@example.com",
				Phone: "123456789",
			},
			expectedDiff: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JsonDiffMapExceptId(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptId() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptId() = %v, want %v", outputDiff, tt.expectedDiff)
			}
		})
	}
}

func TestChoicesAttributesDiff(t *testing.T) {
	tests := []struct {
		name           string
		newStruct      interface{}
		existingStruct interface{}
		resetFields    bool
		expectedDiff   map[string]interface{}
	}{
		{
			name:        "Choices new attr (with rf=true)",
			resetFields: true,
			newStruct: &objects.Device{
				Airflow: &objects.FrontToRear,
				Status:  &objects.DeviceStatusActive,
			},
			existingStruct: &objects.Device{
				Status: &objects.DeviceStatusOffline,
			},
			expectedDiff: map[string]interface{}{
				"airflow": objects.FrontToRear.Value,
				"status":  objects.DeviceStatusActive.Value,
			},
		},
		{
			name:        "Choices attr removal with resetFields=true",
			resetFields: true,
			newStruct: &objects.Device{
				Status: &objects.DeviceStatusActive,
			},
			existingStruct: &objects.Device{
				Status:  &objects.DeviceStatusOffline,
				Airflow: &objects.FrontToRear,
			},
			expectedDiff: map[string]interface{}{
				"airflow": nil,
				"status":  objects.DeviceStatusActive.Value,
			},
		},
		{
			name:        "Removal with resetFields=false",
			resetFields: false,
			newStruct: &objects.Device{
				Status: &objects.DeviceStatusActive,
			},
			existingStruct: &objects.Device{
				Status:  &objects.DeviceStatusOffline,
				Airflow: &objects.FrontToRear,
			},
			expectedDiff: map[string]interface{}{
				"status": objects.DeviceStatusActive.Value,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JsonDiffMapExceptId(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptId() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptId() = %v, want %v", outputDiff, tt.expectedDiff)
			}
		})
	}
}

func TestStructAttributeDiff(t *testing.T) {
	tests := []struct {
		name           string
		newStruct      interface{}
		existingStruct interface{}
		resetFields    bool
		expectedDiff   map[string]interface{}
	}{
		{
			name:        "Struct diff with reset",
			resetFields: true,
			newStruct: &objects.Device{
				DeviceType: &objects.DeviceType{
					NetboxObject: objects.NetboxObject{
						Id: 1,
					},
				},
			},
			existingStruct: &objects.Device{
				DeviceType: &objects.DeviceType{
					NetboxObject: objects.NetboxObject{
						Id: 2,
					},
				},
				DeviceRole: &objects.DeviceRole{
					NetboxObject: objects.NetboxObject{
						Id: 3,
					},
				},
			},
			expectedDiff: map[string]interface{}{
				"device_type": IDObject{ID: 1},
				"role":        nil,
			},
		},
		{
			name:        "Struct diff without reset",
			resetFields: false,
			newStruct: &objects.Device{
				DeviceType: &objects.DeviceType{
					NetboxObject: objects.NetboxObject{
						Id: 1,
					},
				},
			},
			existingStruct: &objects.Device{
				DeviceType: &objects.DeviceType{
					NetboxObject: objects.NetboxObject{
						Id: 2,
					},
				},
				DeviceRole: &objects.DeviceRole{
					NetboxObject: objects.NetboxObject{
						Id: 3,
					},
				},
			},
			expectedDiff: map[string]interface{}{
				"device_type": IDObject{ID: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JsonDiffMapExceptId(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptId() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptId() = %v, want %v", outputDiff, tt.expectedDiff)
			}
		})
	}
}

func TestSliceAttributeDiff(t *testing.T) {
	tests := []struct {
		name           string
		newStruct      interface{}
		existingStruct interface{}
		resetFields    bool
		expectedDiff   map[string]interface{}
	}{
		{
			name:        "Slice diff with reset",
			resetFields: true,
			newStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{Id: 1}},
					{NetboxObject: objects.NetboxObject{Id: 2}},
				},
			},
			existingStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{Id: 3}},
					{NetboxObject: objects.NetboxObject{Id: 4}},
				},
				Mode: &objects.InterfaceModeAccess,
			},
			expectedDiff: map[string]interface{}{
				"tagged_vlans": []int{1, 2},
				"mode":         nil,
			},
		},
		{
			name:        "Slice diff without reset",
			resetFields: false,
			newStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{Id: 1}},
					{NetboxObject: objects.NetboxObject{Id: 2}},
				},
			},
			existingStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{Id: 3}},
					{NetboxObject: objects.NetboxObject{Id: 4}},
				},
				Mode: &objects.InterfaceModeAccess,
			},
			expectedDiff: map[string]interface{}{
				"tagged_vlans": []int{1, 2},
			},
		},
		{
			name:        "Slices no diff",
			resetFields: false,
			newStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{Id: 1}},
					{NetboxObject: objects.NetboxObject{Id: 2}},
				},
			},
			existingStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{Id: 1}},
					{NetboxObject: objects.NetboxObject{Id: 2}},
				},
				Mode: &objects.InterfaceModeAccess,
			},
			expectedDiff: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JsonDiffMapExceptId(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptId() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptId() = %v, want %v", outputDiff, tt.expectedDiff)
			}
		})
	}
}

func TestMapAttributeDiff(t *testing.T) {
	tests := []struct {
		name           string
		newStruct      interface{}
		existingStruct interface{}
		resetFields    bool
		expectedDiff   map[string]interface{}
	}{
		{
			name:        "Map diff with reset",
			resetFields: true,
			newStruct: &objects.Device{
				CustomFields: map[string]string{
					"host_cpu_cores": "10 cpu cores",
					"host_mem":       "10 GB",
					"host_id":        "123456789",
				},
			},
			existingStruct: &objects.Device{
				CustomFields: map[string]string{
					"host_cpu_cores": "5 cpu cores",
					"existing_tag1":  "existing_tag1",
					"existing_tag2":  "existing_tag2",
				},
			},
			expectedDiff: map[string]interface{}{
				"custom_fields": map[string]interface{}{
					"host_cpu_cores": "10 cpu cores",
					"host_mem":       "10 GB",
					"host_id":        "123456789",
					"existing_tag1":  "existing_tag1",
					"existing_tag2":  "existing_tag2",
				},
			},
		},
		{
			name:        "Map no diff with reset",
			resetFields: true,
			newStruct: &objects.Device{
				CustomFields: map[string]string{
					"host_cpu_cores": "10 cpu cores",
					"host_mem":       "10 GB",
				},
			},
			existingStruct: &objects.Device{
				CustomFields: map[string]string{
					"host_cpu_cores": "10 cpu cores",
					"host_mem":       "10 GB",
					"existing_tag1":  "existing_tag1",
					"existing_tag2":  "existing_tag2",
				},
			},
			expectedDiff: map[string]interface{}{},
		},
		{
			name:        "Map single diff with reset",
			resetFields: true,
			newStruct: &objects.Device{
				CustomFields: map[string]string{
					"host_cpu_cores": "5 cpu cores",
					"host_mem":       "10 GB",
				},
			},
			existingStruct: &objects.Device{
				CustomFields: map[string]string{
					"host_cpu_cores": "10 cpu cores",
					"host_mem":       "10 GB",
					"existing_tag1":  "existing_tag1",
					"existing_tag2":  "existing_tag2",
				},
			},
			expectedDiff: map[string]interface{}{
				"custom_fields": map[string]interface{}{
					"host_cpu_cores": "5 cpu cores",
					"existing_tag1":  "existing_tag1",
					"existing_tag2":  "existing_tag2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JsonDiffMapExceptId(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptId() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptId() = %v, want %v", outputDiff, tt.expectedDiff)
			}
		})
	}
}

func TestPriorityMergeDiff(t *testing.T) {
	tests := []struct {
		name           string
		newStruct      interface{}
		existingStruct interface{}
		resetFields    bool
		sourcePriority map[string]int
		expectedDiff   map[string]interface{}
	}{
		{
			name:        "First object has higher priority",
			resetFields: false,
			newStruct: &objects.Vlan{
				Name: "Vlan1000",
				Vid:  1000,
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: "test1",
					},
				},
			},
			existingStruct: &objects.Vlan{
				Name: "1000Vlan",
				Vid:  1000,
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: "test2",
					},
				},
			},
			sourcePriority: map[string]int{
				"test1": 0,
				"test2": 1,
			},
			expectedDiff: map[string]interface{}{
				"name": "Vlan1000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JsonDiffMapExceptId(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptId() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptId() = %v, want %v", outputDiff, tt.expectedDiff)
			}
		})
	}
}
