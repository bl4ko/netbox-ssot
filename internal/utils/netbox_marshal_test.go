package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

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
				ID: 2,
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
				ID: 4,
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
				ID: 2,
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
				ID: 1,
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
	expectedJSONMap := map[string]interface{}{
		"description": "Test Description",
		"tags":        []int{1, 3, 4},
		"name":        "Test",
		"type":        2,
		"status":      "active",
		"site":        2,
		"group":       4,
		"tenant":      1,
	}
	expectedJSON, err := json.Marshal(expectedJSONMap)
	if err != nil {
		t.Errorf("NetboxMarshal() error = %v", err)
	}
	responseJSON, err := NetboxJSONMarshal(newObj)
	if err != nil {
		t.Errorf("NetboxMarshal() error = %v", err)
	}
	if !reflect.DeepEqual(responseJSON, expectedJSON) {
		t.Errorf("NetboxMarshal() = %s\nwant %s", string(responseJSON), string(expectedJSON))
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
				ID: 1,
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
				ID: 1,
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
				ID:          1,
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
	expectedJSON, err := json.Marshal(expectedMap)
	if err != nil {
		t.Errorf("NetboxMarshal() error = %v", err)
	}
	responseJSON, err := NetboxJSONMarshal(device)
	fmt.Println(string(responseJSON))
	if err != nil {
		t.Errorf("NetboxMarshal() error = %v", err)
	}
	if !reflect.DeepEqual(expectedJSON, responseJSON) {
		t.Errorf("NetboxMarshal() = %s\nwant %s", string(responseJSON), string(expectedJSON))
	}
}

// func TestNetboxJsonMarshalComplex(t *testing.T) {
// 	testDevice := objects.Interface{
// 		NetboxObject: objects.NetboxObject{
// 			Tags: []*objects.Tag{
// 				&objects.Tag{
// 					ID: 4,
// 				},
// 				&objects.Tag{
// 					ID: 22,
// 				},
// 				&objects.Tag{
// 					ID: 14,
// 				},
// 			},
// 			Description: "10GB/s pNIC (vSwitch0)",
// 		},
// 		Name:   "vmnic0",
// 		Status: true,
// 		Type:   &objects.OtherInterfaceType,
// 		Speed:  10000,
// 		MTU:    1500,
// 		Mode:   &objects.InterfaceModeTagged,
// 		TaggedVlans: ,
// 	}
// }

type testStruct struct {
	Test string `json:"test"`
}

type testStructWithStructAttribute struct {
	Test testStruct `json:"test"`
}

type structWithSliceAttribute struct {
	Test []string `json:"test"`
}

type structWithSliceAttributeOfStructs struct {
	Test []testStruct `json:"test"`
}

func TestNetboxJSONMarshal(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test when slices attribute doesn't have ids",
			args: args{
				obj: structWithSliceAttribute{Test: []string{"one", "two", "three"}},
			},
			want: []byte("{\"test\":[\"one\",\"two\",\"three\"]}"),
		},
		{
			name: "Struct with slice attribute of structs with no ids",
			args: args{
				obj: structWithSliceAttributeOfStructs{Test: []testStruct{{Test: "one"}, {Test: "two"}}},
			},
			want: []byte("{\"test\":[{\"test\":\"one\"},{\"test\":\"two\"}]}"),
		},
		{
			name: "Struct that has no ID attribute marshal",
			args: args{
				obj: testStructWithStructAttribute{Test: testStruct{Test: "test"}},
			},
			want: []byte("{\"test\":{\"test\":\"test\"}}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NetboxJSONMarshal(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetboxJSONMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxJSONMarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructToNetboxJSONMap(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Skip when slice has length of zero",
			args: args{
				obj: interface{}(objects.NetboxObject{
					Tags: []*objects.Tag{},
				}),
			},
			want: map[string]interface{}{},
		},
		{
			name: "From slice attributes with ids extract only ids",
			args: args{
				obj: interface{}(objects.NetboxObject{
					Tags: []*objects.Tag{{ID: 1}, {ID: 2}, {ID: 3}},
				}),
			},
			want: map[string]interface{}{"tags": []interface{}{int64(1), int64(2), int64(3)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StructToNetboxJSONMap(tt.args.obj)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StructToNetboxJSONMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
