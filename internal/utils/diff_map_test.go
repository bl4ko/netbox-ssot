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
			outputDiff, err := JSONDiffMapExceptID(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptID() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptID() = %v, want %v", outputDiff, tt.expectedDiff)
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
			outputDiff, err := JSONDiffMapExceptID(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptID() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptID() = %v, want %v", outputDiff, tt.expectedDiff)
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
						ID: 1,
					},
				},
			},
			existingStruct: &objects.Device{
				DeviceType: &objects.DeviceType{
					NetboxObject: objects.NetboxObject{
						ID: 2,
					},
				},
				DeviceRole: &objects.DeviceRole{
					NetboxObject: objects.NetboxObject{
						ID: 3,
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
						ID: 1,
					},
				},
			},
			existingStruct: &objects.Device{
				DeviceType: &objects.DeviceType{
					NetboxObject: objects.NetboxObject{
						ID: 2,
					},
				},
				DeviceRole: &objects.DeviceRole{
					NetboxObject: objects.NetboxObject{
						ID: 3,
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
			outputDiff, err := JSONDiffMapExceptID(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptID() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptID() = %v, want %v", outputDiff, tt.expectedDiff)
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
					{NetboxObject: objects.NetboxObject{ID: 1}},
					{NetboxObject: objects.NetboxObject{ID: 2}},
				},
			},
			existingStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{ID: 3}},
					{NetboxObject: objects.NetboxObject{ID: 4}},
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
					{NetboxObject: objects.NetboxObject{ID: 1}},
					{NetboxObject: objects.NetboxObject{ID: 2}},
				},
			},
			existingStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{ID: 3}},
					{NetboxObject: objects.NetboxObject{ID: 4}},
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
					{NetboxObject: objects.NetboxObject{ID: 1}},
					{NetboxObject: objects.NetboxObject{ID: 2}},
				},
			},
			existingStruct: &objects.Interface{
				TaggedVlans: []*objects.Vlan{
					{NetboxObject: objects.NetboxObject{ID: 1}},
					{NetboxObject: objects.NetboxObject{ID: 2}},
				},
				Mode: &objects.InterfaceModeAccess,
			},
			expectedDiff: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JSONDiffMapExceptID(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptID() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptID() = %v, want %v", outputDiff, tt.expectedDiff)
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
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldHostCPUCoresName: "10 cpu cores",
						constants.CustomFieldHostMemoryName:   "10 GB",
						constants.CustomFieldSourceIDName:     "123456789",
					},
				},
			},
			existingStruct: &objects.Device{
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldHostCPUCoresName: "5 cpu cores",
						"existing_tag1":                       "existing_tag1",
						"existing_tag2":                       "existing_tag2",
					},
				},
			},
			expectedDiff: map[string]interface{}{
				"custom_fields": map[string]interface{}{
					constants.CustomFieldHostCPUCoresName: "10 cpu cores",
					constants.CustomFieldHostMemoryName:   "10 GB",
					constants.CustomFieldSourceIDName:     "123456789",
					"existing_tag1":                       "existing_tag1",
					"existing_tag2":                       "existing_tag2",
				},
			},
		},
		{
			name:        "Map no diff with reset",
			resetFields: true,
			newStruct: &objects.Device{
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldHostCPUCoresName: "10 cpu cores",
						constants.CustomFieldHostMemoryName:   "10 GB",
					},
				},
			},

			existingStruct: &objects.Device{
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldHostCPUCoresName: "10 cpu cores",
						constants.CustomFieldHostMemoryName:   "10 GB",
						"existing_tag1":                       "existing_tag1",
						"existing_tag2":                       "existing_tag2",
					},
				},
			},
			expectedDiff: map[string]interface{}{},
		},
		{
			name:        "Map single diff with reset",
			resetFields: true,
			newStruct: &objects.Device{
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldHostCPUCoresName: "5 cpu cores",
						constants.CustomFieldHostMemoryName:   "10 GB",
					},
				},
			},
			existingStruct: &objects.Device{
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldHostCPUCoresName: "10 cpu cores",
						constants.CustomFieldHostMemoryName:   "10 GB",
						"existing_tag1":                       "existing_tag1",
						"existing_tag2":                       "existing_tag2",
					},
				},
			},
			expectedDiff: map[string]interface{}{
				"custom_fields": map[string]interface{}{
					constants.CustomFieldHostCPUCoresName: "5 cpu cores",
					"existing_tag1":                       "existing_tag1",
					"existing_tag2":                       "existing_tag2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JSONDiffMapExceptID(tt.newStruct, tt.existingStruct, tt.resetFields, nil)
			if err != nil {
				t.Errorf("JsonDiffMapExceptID() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptID() = %v, want %v", outputDiff, tt.expectedDiff)
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
					Tags: []*objects.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
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
					Tags: []*objects.Tag{
						{ID: 2, Name: "Tag1"},
						{ID: 3, Name: "Tag2"},
					},
				},
			},
			sourcePriority: map[string]int{
				"test1": 0,
				"test2": 1,
			},
			expectedDiff: map[string]interface{}{
				"name": "Vlan1000",
				"custom_fields": map[string]interface{}{
					constants.CustomFieldSourceName: "test1",
				},
				"tags": []int{1, 2},
			},
		},
		{
			name:        "Second object has higher priority",
			resetFields: false,
			newStruct: &objects.Vlan{
				Name:     "Vlan1000",
				Vid:      1000,
				Comments: "Added comment",
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]string{
						constants.CustomFieldSourceName: "test1",
					},
					Tags: []*objects.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
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
					Tags: []*objects.Tag{
						{ID: 2, Name: "Tag1"},
						{ID: 3, Name: "Tag2"},
					},
				},
			},
			sourcePriority: map[string]int{
				"test1": 1,
				"test2": 0,
			},
			expectedDiff: map[string]interface{}{
				"comments": "Added comment",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDiff, err := JSONDiffMapExceptID(tt.newStruct, tt.existingStruct, tt.resetFields, tt.sourcePriority)
			if err != nil {
				t.Errorf("JsonDiffMapExceptID() error = %v", err)
			}
			if !reflect.DeepEqual(outputDiff, tt.expectedDiff) {
				t.Errorf("JsonDiffMapExceptID() = %v, want %v", outputDiff, tt.expectedDiff)
			}
		})
	}
}

func Test_isChoiceEmbedded(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isChoiceEmbedded(tt.args.v); got != tt.want {
				t.Errorf("isChoiceEmbedded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_choiceValue(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := choiceValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("choiceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasPriorityOver(t *testing.T) {
	type args struct {
		newObj          reflect.Value
		existingObj     reflect.Value
		source2priority map[string]int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasPriorityOver(tt.args.newObj, tt.args.existingObj, tt.args.source2priority); got != tt.want {
				t.Errorf("hasPriorityOver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONDiffMapExceptID(t *testing.T) {
	type args struct {
		newObj          interface{}
		existingObj     interface{}
		resetFields     bool
		source2priority map[string]int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONDiffMapExceptID(tt.args.newObj, tt.args.existingObj, tt.args.resetFields, tt.args.source2priority)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONDiffMapExceptID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONDiffMapExceptID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addSliceDiff(t *testing.T) {
	type args struct {
		newSlice      reflect.Value
		existingSlice reflect.Value
		jsonTag       string
		hasPriority   bool
		diffMap       map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addSliceDiff(tt.args.newSlice, tt.args.existingSlice, tt.args.jsonTag, tt.args.hasPriority, tt.args.diffMap); (err != nil) != tt.wantErr {
				t.Errorf("addSliceDiff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_addStructDiff(t *testing.T) {
	type args struct {
		newObj      reflect.Value
		existingObj reflect.Value
		jsonTag     string
		hasPriority bool
		diffMap     map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addStructDiff(tt.args.newObj, tt.args.existingObj, tt.args.jsonTag, tt.args.hasPriority, tt.args.diffMap); (err != nil) != tt.wantErr {
				t.Errorf("addStructDiff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_addMapDiff(t *testing.T) {
	type args struct {
		newMap      reflect.Value
		existingMap reflect.Value
		jsonTag     string
		hasPriority bool
		diffMap     map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addMapDiff(tt.args.newMap, tt.args.existingMap, tt.args.jsonTag, tt.args.hasPriority, tt.args.diffMap); (err != nil) != tt.wantErr {
				t.Errorf("addMapDiff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_addPrimaryDiff(t *testing.T) {
	type args struct {
		newField      reflect.Value
		existingField reflect.Value
		jsonTag       string
		hasPriority   bool
		diffMap       map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			addPrimaryDiff(tt.args.newField, tt.args.existingField, tt.args.jsonTag, tt.args.hasPriority, tt.args.diffMap)
		})
	}
}
