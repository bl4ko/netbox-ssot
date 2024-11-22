// This file contains all objects that are common to all Netbox objects.
package objects

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func TestChoice_String(t *testing.T) {
	tests := []struct {
		name string
		c    Choice
		want string
	}{
		{
			name: "Test choice correct string",
			c: Choice{
				Value: "test value",
				Label: "test label",
			},
			want: "test value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Choice.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxObject_String(t *testing.T) {
	tests := []struct {
		name string
		n    NetboxObject
		want string
	}{
		{
			name: "Test netbox object correct string",
			n: NetboxObject{
				ID: 1,
				Tags: []*Tag{
					{Name: "Test tag1"}, {Name: "Test tag2"},
				},
				Description: "Test description",
				CustomFields: map[string]interface{}{
					constants.CustomFieldSourceName: "test",
				},
			},
			want: fmt.Sprintf("ID: %d, Tags: %s, Description: %s", 1, []*Tag{
				{Name: "Test tag1"}, {Name: "Test tag2"},
			}, "Test description"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("NetboxObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxObject_GetID(t *testing.T) {
	tests := []struct {
		name string
		n    *NetboxObject
		want int
	}{
		{
			name: "Test netbox object get id",
			n: &NetboxObject{
				ID: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.GetID(); got != tt.want {
				t.Errorf("NetboxObject.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxObject_GetCustomField(t *testing.T) {
	type args struct {
		label string
	}
	tests := []struct {
		name string
		n    *NetboxObject
		args args
		want interface{}
	}{
		{
			name: "Test netbox object get custom field",
			n: &NetboxObject{
				CustomFields: map[string]interface{}{
					"test": "test",
				},
			},
			args: args{
				label: "test",
			},
			want: "test",
		},
		{
			name: "Test netbox object get custom field",
			n:    &NetboxObject{},
			args: args{
				label: "test",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.GetCustomField(tt.args.label); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetboxObject.GetCustomField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxObject_SetCustomField(t *testing.T) {
	type args struct {
		label string
		value interface{}
	}
	tests := []struct {
		name string
		n    *NetboxObject
		args args
	}{
		{
			name: "Test netbox object set custom field",
			n:    &NetboxObject{},
			args: args{
				label: "test",
				value: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tt.n.SetCustomField(tt.args.label, tt.args.value)
		})
	}
}

func TestNetboxObject_AddTag(t *testing.T) {
	type args struct {
		newTag *Tag
	}
	tests := []struct {
		name string
		n    *NetboxObject
		args args
	}{
		{
			name: "Test netbox object add tag",
			n: &NetboxObject{
				Tags: []*Tag{
					{Name: "Test tag1"},
					{Name: "Test tag2"},
				},
			},
			args: args{
				newTag: &Tag{Name: "Test tag3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tt.n.AddTag(tt.args.newTag)
		})
	}
}

func TestNetboxObject_HasTag(t *testing.T) {
	type args struct {
		tag *Tag
	}
	tests := []struct {
		name string
		n    *NetboxObject
		args args
		want bool
	}{
		{
			name: "Test netbox object has tag",
			n: &NetboxObject{
				Tags: []*Tag{
					{Name: "Test tag1"},
					{Name: "Test tag2"},
				},
			},
			args: args{
				tag: &Tag{Name: "Test tag1"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.HasTag(tt.args.tag); got != tt.want {
				t.Errorf("NetboxObject.HasTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxObject_HasTagByName(t *testing.T) {
	type args struct {
		tagName string
	}
	tests := []struct {
		name string
		n    *NetboxObject
		args args
		want bool
	}{
		{
			name: "Test netbox object has tag by name",
			n: &NetboxObject{
				Tags: []*Tag{
					{Name: "Test tag1"},
					{Name: "Test tag2"},
				},
			},
			args: args{
				tagName: "Test tag1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.HasTagByName(tt.args.tagName); got != tt.want {
				t.Errorf("NetboxObject.HasTagByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxObject_RemoveTag(t *testing.T) {
	type args struct {
		tag *Tag
	}
	tests := []struct {
		name string
		n    *NetboxObject
		args args
	}{
		{
			name: "Test netbox object remove tag",
			n: &NetboxObject{
				Tags: []*Tag{
					{Name: "Test tag1"},
					{Name: "Test tag2"},
				},
			},
			args: args{
				tag: &Tag{Name: "Test tag1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tt.n.RemoveTag(tt.args.tag)
		})
	}
}
