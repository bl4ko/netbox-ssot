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
					CustomFields: map[string]interface{}{
						constants.CustomFieldHostCPUCoresName: "10 cpu cores",
						constants.CustomFieldHostMemoryName:   "10 GB",
					},
				},
			},
			existingStruct: &objects.Device{
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
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
					CustomFields: map[string]interface{}{
						constants.CustomFieldHostCPUCoresName: "10 cpu cores",
						constants.CustomFieldHostMemoryName:   "10 GB",
					},
				},
			},

			existingStruct: &objects.Device{
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
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
					CustomFields: map[string]interface{}{
						constants.CustomFieldHostCPUCoresName: "5 cpu cores",
						constants.CustomFieldHostMemoryName:   "10 GB",
					},
				},
			},
			existingStruct: &objects.Device{
				NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
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
					CustomFields: map[string]interface{}{
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
					CustomFields: map[string]interface{}{
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
					CustomFields: map[string]interface{}{
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
					CustomFields: map[string]interface{}{
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
		{
			name: "When custom fields on object are missing return true",
			args: args{
				newObj:          reflect.ValueOf(objects.Tag{Name: "NewDevice"}),
				existingObj:     reflect.ValueOf(objects.Tag{Name: "ExistingDevice"}),
				source2priority: map[string]int{},
			},
			want: true,
		},
		{
			name: "IP address representing arp entry has always lower priority than standard IP address",
			args: args{
				newObj: reflect.ValueOf(objects.IPAddress{NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: false,
						constants.CustomFieldSourceName:   "source1",
					},
				}}),
				existingObj: reflect.ValueOf(objects.IPAddress{NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: true,
						constants.CustomFieldSourceName:   "source2",
					},
				}}),
				source2priority: map[string]int{
					"source1": 1, "source2": 2,
				},
			},
			want: true,
		},
		{
			name: "IP address representing arp entry has always lower priority than standard IP address",
			args: args{
				newObj: reflect.ValueOf(objects.IPAddress{NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: true,
						constants.CustomFieldSourceName:   "source1",
					},
				}}),
				existingObj: reflect.ValueOf(objects.IPAddress{NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: true,
						constants.CustomFieldSourceName:   "source2",
					},
				}}),
				source2priority: map[string]int{
					"source1": 1, "source2": 2,
				},
			},
			want: true,
		},
		{
			name: "IP address representing arp entry has always lower priority than standard IP address",
			args: args{
				newObj: reflect.ValueOf(objects.IPAddress{NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: true,
						constants.CustomFieldSourceName:   "source2",
					},
				}}),
				existingObj: reflect.ValueOf(objects.IPAddress{NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: true,
						constants.CustomFieldSourceName:   "source1",
					},
				}}),
				source2priority: map[string]int{
					"source1": 1, "source2": 2,
				},
			},
			want: false,
		},
		{
			name: "IP address that has set arpentryname to nil and existing is arp data has priority over",
			args: args{
				newObj: reflect.ValueOf(objects.IPAddress{NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: nil,
						constants.CustomFieldSourceName:   "source2",
					},
				}}),
				existingObj: reflect.ValueOf(objects.IPAddress{NetboxObject: objects.NetboxObject{
					CustomFields: map[string]interface{}{
						constants.CustomFieldArpEntryName: true,
						constants.CustomFieldSourceName:   "source1",
					},
				}}),
				source2priority: map[string]int{
					"source1": 1, "source2": 2,
				},
			},
			want: true,
		},
		{
			name: "IP address representing arp entry has always lower priority than standard IP address",
			args: args{
				newObj:          reflect.ValueOf(objects.Tag{Name: "NewDevice"}),
				existingObj:     reflect.ValueOf(objects.Tag{Name: "ExistingDevice"}),
				source2priority: map[string]int{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasPriorityOver(tt.args.newObj, tt.args.existingObj, tt.args.source2priority); got != tt.want {
				t.Errorf("hasPriorityOver() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testStructWithTestAttribute struct {
	Test string
}

type testStructWithStringIDAttribute struct {
	ID string
}

type testStructWithIntIDAttribute struct {
	ID int
}

type testStruct2 struct {
	Test int
}

type testSliceStruct struct {
	Test []testStructWithTestAttribute
}

type wrongIDField struct {
	ID string
}

type testStructWithWrongIDField struct {
	SubStruct wrongIDField
}

type structWithIntMapAttribute struct {
	Map map[int]int
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
		{
			name: "Test diff map with objects of different kinds",
			args: args{
				newObj:          objects.Device{Name: "TestDevice"},
				existingObj:     &objects.Interface{Name: "TestInterface"},
				resetFields:     true,
				source2priority: map[string]int{},
			},
			wantErr: true,
		},
		{
			name: "New object is not a struct",
			args: args{
				newObj:          "TestDevice",
				existingObj:     "Existing device",
				resetFields:     true,
				source2priority: map[string]int{},
			},
			wantErr: true,
		},
		{
			name: "Trigger error for fieldName == \"NetboxObject\"",
			args: args{
				newObj: &objects.Device{
					Name: "NewDevice",
				},
				existingObj: &objects.Tag{
					Name: "ExistingTag",
				},
				resetFields: true,
				source2priority: map[string]int{
					"InvalidPrioritySetting": -1,
				},
			},
			wantErr: true,
		},
		{
			name: "New object has a field with no json",
			args: args{
				newObj:          testStructWithTestAttribute{Test: "Test"},
				existingObj:     testStructWithTestAttribute{Test: "Test2"},
				resetFields:     true,
				source2priority: map[string]int{},
			},
			want: map[string]interface{}{
				"Test": "Test",
			},
			wantErr: false,
		},
		{
			name: "Fail with both fields are not of the same type",
			args: args{
				newObj:          testStructWithTestAttribute{Test: "test"},
				existingObj:     testStruct2{Test: 1},
				resetFields:     false,
				source2priority: map[string]int{},
			},
			wantErr: true,
		},
		{
			name: "Fail with case reflect.Slice",
			args: args{
				newObj:          testSliceStruct{Test: []testStructWithTestAttribute{{Test: "test1"}, {Test: "test2"}}},
				existingObj:     testSliceStruct{Test: []testStructWithTestAttribute{}},
				resetFields:     false,
				source2priority: map[string]int{},
			},
			wantErr: true,
		},
		{
			name: "Fail with case reflect.Struct",
			args: args{
				newObj:          testStructWithWrongIDField{SubStruct: wrongIDField{ID: "wrong"}},
				existingObj:     testStructWithWrongIDField{SubStruct: wrongIDField{ID: "should be int"}},
				resetFields:     false,
				source2priority: map[string]int{},
			},
			wantErr: true,
		},
		{
			name: "Fail with case reflect.Map",
			args: args{
				newObj:          structWithIntMapAttribute{Map: map[int]int{1: 1}},
				existingObj:     structWithIntMapAttribute{Map: map[int]int{1: 1}},
				resetFields:     false,
				source2priority: map[string]int{},
			},
			wantErr: true,
		},
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
		name        string
		args        args
		wantErr     bool
		wantDiffMap map[string]interface{}
	}{
		{
			name: "new slice has priority, but is empty or nil",
			args: args{
				newSlice:      reflect.ValueOf([]int{}),
				existingSlice: reflect.ValueOf([]int{1, 2, 3}),
				jsonTag:       "test",
				hasPriority:   true,
				diffMap:       map[string]interface{}{},
			},
			wantDiffMap: map[string]interface{}{"test": []interface{}{}},
		},
		{
			name: "Has priority newSlice different length",
			args: args{
				newSlice:      reflect.ValueOf([]string{"pineapple", "strawberry"}),
				existingSlice: reflect.ValueOf([]string{"pineapple"}),
				jsonTag:       "test",
				hasPriority:   true,
				diffMap:       map[string]interface{}{},
			},
			wantDiffMap: map[string]interface{}{"test": []string{"pineapple", "strawberry"}},
		},
		{
			name: "Has priority same lengths",
			args: args{
				newSlice:      reflect.ValueOf([]string{"pineapple", "strawberry"}),
				existingSlice: reflect.ValueOf([]string{"pineapple", "pie"}),
				jsonTag:       "test",
				hasPriority:   true,
				diffMap:       map[string]interface{}{},
			},
			wantDiffMap: map[string]interface{}{"test": []string{"pineapple", "strawberry"}},
		},
		{
			name: "No priority newSlice different length",
			args: args{
				newSlice:      reflect.ValueOf([]string{"pineapple", "strawberry"}),
				existingSlice: reflect.ValueOf([]string{"pineapple"}),
				jsonTag:       "test",
				hasPriority:   false,
				diffMap:       map[string]interface{}{},
			},
			wantDiffMap: map[string]interface{}{},
		},
		{
			name: "Has priority. Different length. Slice with structs with ID attributes.",
			args: args{
				newSlice:      reflect.ValueOf([]testStructWithIntIDAttribute{{ID: 1}, {ID: 2}}),
				existingSlice: reflect.ValueOf([]testStructWithIntIDAttribute{{ID: 1}}),
				jsonTag:       "test",
				hasPriority:   true,
				diffMap:       map[string]interface{}{},
			},
			wantErr:     false,
			wantDiffMap: map[string]interface{}{"test": []int{1, 2}},
		},
		{
			name: "Has priority. Same length but different. Slice with structs with ID attributes.",
			args: args{
				newSlice:      reflect.ValueOf([]testStructWithIntIDAttribute{{ID: 1}, {ID: 2}}),
				existingSlice: reflect.ValueOf([]testStructWithIntIDAttribute{{ID: 1}, {ID: 3}}),
				jsonTag:       "test",
				hasPriority:   true,
				diffMap:       map[string]interface{}{},
			},
			wantErr:     false,
			wantDiffMap: map[string]interface{}{"test": []int{1, 2}},
		},
		{
			name: "Has priority. Same length and the same. Slice with structs with ID attributes.",
			args: args{
				newSlice:      reflect.ValueOf([]testStructWithIntIDAttribute{{ID: 1}, {ID: 2}}),
				existingSlice: reflect.ValueOf([]testStructWithIntIDAttribute{{ID: 1}, {ID: 2}}),
				jsonTag:       "test",
				hasPriority:   false,
				diffMap:       map[string]interface{}{},
			},
			wantErr:     false,
			wantDiffMap: map[string]interface{}{},
		},
		{
			name: "Test interface slices of same length. Fails because struct elements don't have an ID attribute",
			args: args{
				newSlice:      reflect.ValueOf([]testStructWithIntIDAttribute{{ID: 1}}),
				existingSlice: reflect.ValueOf([]testStructWithTestAttribute{{Test: "1"}}),
				jsonTag:       "test",
				hasPriority:   true,
				diffMap:       map[string]interface{}{},
			},
			wantErr:     true,
			wantDiffMap: map[string]interface{}{},
		},
		{
			name: "Test interface slices of same length. Fails because struct elements don't have an ID attribute",
			args: args{
				existingSlice: reflect.ValueOf([]testStructWithTestAttribute{{Test: "1"}, {Test: "test"}}),
				newSlice:      reflect.ValueOf([]testStructWithTestAttribute{{Test: "1"}}),
				jsonTag:       "test",
				hasPriority:   true,
				diffMap:       map[string]interface{}{},
			},
			wantErr:     true,
			wantDiffMap: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addSliceDiff(tt.args.newSlice, tt.args.existingSlice, tt.args.jsonTag, tt.args.hasPriority, tt.args.diffMap); (err != nil) != tt.wantErr {
				t.Errorf("addSliceDiff() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				if !reflect.DeepEqual(tt.wantDiffMap, tt.args.diffMap) {
					t.Errorf("got diffmap: %s, want diffmap: %s", tt.args.diffMap, tt.wantDiffMap)
				}
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
		{
			name: "Skip if new object is invalid",
			args: args{
				newObj:      reflect.ValueOf(nil),
				existingObj: reflect.ValueOf("string"),
				jsonTag:     "Test",
				diffMap:     map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "Existing obj is invalid",
			args: args{
				newObj:      reflect.ValueOf(testStructWithTestAttribute{Test: "new"}),
				existingObj: reflect.ValueOf(nil),
				jsonTag:     "Test",
				diffMap:     map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "New object has priority and is different",
			args: args{
				newObj:      reflect.ValueOf(testStructWithTestAttribute{Test: "new string"}),
				existingObj: reflect.ValueOf(testStructWithTestAttribute{Test: "old string"}),
				jsonTag:     "Test",
				diffMap:     map[string]interface{}{},
				hasPriority: true,
			},
			wantErr: false,
		},
		{
			name: "New object doesn't have priority and is different",
			args: args{
				newObj:      reflect.ValueOf(testStructWithTestAttribute{Test: "new string"}),
				existingObj: reflect.ValueOf(testStructWithTestAttribute{Test: "old string"}),
				jsonTag:     "Test",
				diffMap:     map[string]interface{}{},
				hasPriority: false,
			},
			wantErr: false,
		},
		{
			name: "New object has priority and is different and has choice attribute",
			args: args{
				newObj:      reflect.ValueOf(objects.InterfaceModeAccess),
				existingObj: reflect.ValueOf(objects.InterfaceModeTagged),
				jsonTag:     "Test",
				diffMap:     map[string]interface{}{},
				hasPriority: true,
			},
			wantErr: false,
		},
		{
			name: "ID field is of type string",
			args: args{
				newObj:      reflect.ValueOf(testStructWithStringIDAttribute{ID: "1"}),
				existingObj: reflect.Value{},
				jsonTag:     "Test",
				diffMap:     map[string]interface{}{},
				hasPriority: true,
			},
			wantErr: true,
		},
		{
			name: "Existing object is not valid",
			args: args{
				newObj:      reflect.ValueOf(testStructWithIntIDAttribute{ID: 1}),
				existingObj: reflect.Value{},
				jsonTag:     "Test",
				diffMap:     map[string]interface{}{},
				hasPriority: true,
			},
			wantErr: false,
		},
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
		name        string
		args        args
		wantErr     bool
		wantDiffMap map[string]interface{}
	}{
		{
			name: "skip if new map is not valid",
			args: args{
				newMap:      reflect.ValueOf(nil),
				existingMap: reflect.ValueOf(map[string]string{}),
				jsonTag:     "test",
				hasPriority: false,
				diffMap:     map[string]interface{}{},
			},
			wantDiffMap: map[string]interface{}{},
		},
		{
			name: "map diff test",
			args: args{
				newMap: reflect.ValueOf(map[string]interface{}{
					constants.CustomFieldArpIPLastSeenName: "2024-04-18 10:59:17",
					constants.CustomFieldSourceIDName:      nil,
				}),
				existingMap: reflect.ValueOf(map[string]interface{}{
					constants.CustomFieldArpIPLastSeenName: "2024-04-18 10:29:30",
					constants.CustomFieldSourceIDName:      nil,
				}),
				hasPriority: true,
				jsonTag:     "CustomFields",
				diffMap:     map[string]interface{}{},
			},
			wantDiffMap: map[string]interface{}{
				"CustomFields": map[string]interface{}{constants.CustomFieldArpIPLastSeenName: "2024-04-18 10:59:17"},
			},
		},
		{
			name: "map diff test",
			args: args{
				newMap: reflect.ValueOf(map[string]interface{}{
					constants.CustomFieldArpIPLastSeenName: "2024-04-18 10:59:17",
				}),
				existingMap: reflect.ValueOf(map[string]interface{}{
					constants.CustomFieldArpIPLastSeenName: "2024-04-18 10:29:30",
					constants.CustomFieldSourceIDName:      nil,
				}),
				hasPriority: true,
				jsonTag:     "CustomFields",
				diffMap:     map[string]interface{}{},
			},
			wantDiffMap: map[string]interface{}{
				"CustomFields": map[string]interface{}{constants.CustomFieldArpIPLastSeenName: "2024-04-18 10:59:17"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addMapDiff(tt.args.newMap, tt.args.existingMap, tt.args.jsonTag, tt.args.hasPriority, tt.args.diffMap); (err != nil) != tt.wantErr {
				t.Errorf("addMapDiff() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.diffMap, tt.wantDiffMap) {
				t.Errorf("diffMap: %s, wantDiffMap: %s", tt.args.diffMap, tt.wantDiffMap)
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

func Test_convertSliceToComparableSlice(t *testing.T) {
	type args struct {
		slice reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertSliceToComparableSlice(tt.args.slice)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertSliceToComparableSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertSliceToComparableSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sliceToSet(t *testing.T) {
	type args struct {
		slice reflect.Value
	}
	tests := []struct {
		name string
		args args
		want map[interface{}]bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceToSet(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceToSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
