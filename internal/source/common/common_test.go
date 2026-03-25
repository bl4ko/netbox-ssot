package common

import (
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestConfig_GetSourceTags(t *testing.T) {
	tag1 := &objects.Tag{Name: "source-name", Slug: "source-name"}
	tag2 := &objects.Tag{Name: "source-type", Slug: "source-type"}

	tests := []struct {
		name string
		cfg  Config
		want []*objects.Tag
	}{
		{
			name: "both tags set",
			cfg:  Config{SourceNameTag: tag1, SourceTypeTag: tag2},
			want: []*objects.Tag{tag1, tag2},
		},
		{
			name: "name tag nil",
			cfg:  Config{SourceNameTag: nil, SourceTypeTag: tag2},
			want: []*objects.Tag{nil, tag2},
		},
		{
			name: "type tag nil",
			cfg:  Config{SourceNameTag: tag1, SourceTypeTag: nil},
			want: []*objects.Tag{tag1, nil},
		},
		{
			name: "both tags nil",
			cfg:  Config{},
			want: []*objects.Tag{nil, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetSourceTags()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.GetSourceTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
