package objects

import (
	"testing"
)

func TestTag_String(t *testing.T) {
	tests := []struct {
		name string
		tr   Tag
		want string
	}{
		{
			name: "Test tag correct string",
			tr: Tag{
				Name:        "Test tag",
				Description: "Test tag description",
			},
			want: "Tag{Name: Test tag}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("Tag.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomField_String(t *testing.T) {
	tests := []struct {
		name string
		cf   CustomField
		want string
	}{
		{
			name: "Test custom field correct string",
			cf: CustomField{
				ID:   10,
				Name: "host_cpu_cores",
			},
			want: "CustomField{ID: 10, Name: host_cpu_cores}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cf.String(); got != tt.want {
				t.Errorf("CustomField.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTag_GetID(t *testing.T) {
	tests := []struct {
		name string
		tr   *Tag
		want int
	}{
		{
			name: "Test tag get id",
			tr: &Tag{
				ID: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetID(); got != tt.want {
				t.Errorf("Tag.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomField_GetID(t *testing.T) {
	tests := []struct {
		name string
		cf   *CustomField
		want int
	}{
		{
			name: "Test Custom Field Get ID",
			cf: &CustomField{
				ID: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cf.GetID(); got != tt.want {
				t.Errorf("CustomField.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}
