package utils

import (
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestExtractJSONTagsFromStruct(t *testing.T) {
	type args struct {
		inputStruct interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Extract fields from tag",
			args: args{
				inputStruct: objects.Tag{},
			},
			want: []string{"id", "name", "slug", "color", "description"},
		},
		{
			name: "Extract fields from custom field",
			args: args{
				inputStruct: objects.CustomField{},
			},
			want: []string{
				"id",
				"name",
				"label",
				"type",
				"object_types",
				"description",
				"search_weight",
				"filter_logic",
				"ui_visible",
				"ui_editable",
				"weight",
				"default",
				"required",
			},
		},
		{
			name: "Extract json fields from device",
			args: args{
				inputStruct: objects.Device{},
			},
			want: []string{
				"id",
				"tags",
				"description",
				"custom_fields",
				"name",
				"role",
				"device_type",
				"airflow",
				"serial",
				"asset_tag",
				"site",
				"location",
				"status",
				"platform",
				"primary_ip4",
				"primary_ip6",
				"cluster",
				"tenant",
				"comments",
			},
		},
		{
			name: "Extract fields from VMInterface",
			args: args{
				inputStruct: objects.VMInterface{},
			},
			want: []string{
				"id",
				"tags",
				"description",
				"custom_fields",
				"virtual_machine",
				"name",
				"primary_mac_address",
				"mtu",
				"enabled",
				"parent",
				"bridge",
				"mode",
				"tagged_vlans",
				"untagged_vlan",
			},
		},
		{
			name: "Extract fields from interface",
			args: args{
				inputStruct: objects.Interface{},
			},
			want: []string{
				"id",
				"tags",
				"description",
				"custom_fields",
				"device",
				"name",
				"enabled",
				"type",
				"speed",
				"parent",
				"bridge",
				"lag",
				"mtu",
				"primary_mac_address",
				"duplex",
				"mode",
				"tagged_vlans",
				"untagged_vlan",
				"vdcs",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractJSONTagsFromStruct(tt.args.inputStruct); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("ExtractStructJSONFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractJSONTagsFromStructIntoString(t *testing.T) {
	type args struct {
		inputStruct interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Extract fields from tag",
			args: args{
				inputStruct: objects.Tag{},
			},
			want: "id,name,slug,color,description",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractJSONTagsFromStructIntoString(tt.args.inputStruct); got != tt.want {
				t.Errorf("ExtractJSONTagsFromStructIntoString() = %v, want %v", got, tt.want)
			}
		})
	}
}
