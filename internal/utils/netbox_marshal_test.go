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
				{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{Id: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Status: objects.ClusterStatusActive,
		Name:   "Test",
		Type: &objects.ClusterType{
			NetboxObject: objects.NetboxObject{
				Id: 2,
				Tags: []*objects.Tag{
					{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				},
			},
			Name: "oVirt",
			Slug: "ovirt",
		},
		Group: &objects.ClusterGroup{
			NetboxObject: objects.NetboxObject{
				Id: 4,
				Tags: []*objects.Tag{
					{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
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
					{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
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
					{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
					{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
					{Id: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
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
				{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
				{Id: 3, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
				{Id: 4, Name: "Test3", Slug: "test3", Color: "000000", Description: "Test tag 3"},
			},
		},
		Name: "Test device",
		DeviceRole: &objects.DeviceRole{
			NetboxObject: objects.NetboxObject{
				Id: 1,
				Tags: []*objects.Tag{
					{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
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
					{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
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
					{Id: 1, Name: "Test", Slug: "test", Color: "000000", Description: "Test tag"},
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
