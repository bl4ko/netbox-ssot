package common

import (
	"errors"
	"reflect"
	"strings"
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

func TestJoinErrors(t *testing.T) {
	first := errors.New("vm1 failed")
	second := errors.New("vm2 failed")
	tests := []struct {
		name     string
		errs     []error
		wantNil  bool
		wantMsgs []string
	}{
		{name: "no error", errs: nil, wantNil: true},
		{name: "only nil errors", errs: []error{nil, nil}, wantNil: true},
		{name: "single error", errs: []error{first}, wantMsgs: []string{"vm1 failed"}},
		{name: "every error is kept", errs: []error{first, nil, second}, wantMsgs: []string{"vm1 failed", "vm2 failed"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errChan := make(chan error, len(tt.errs))
			for _, err := range tt.errs {
				errChan <- err
			}
			close(errChan)
			got := JoinErrors(errChan)
			if tt.wantNil {
				if got != nil {
					t.Fatalf("JoinErrors() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatalf("JoinErrors() = nil, want %v", tt.wantMsgs)
			}
			for _, msg := range tt.wantMsgs {
				if !strings.Contains(got.Error(), msg) {
					t.Errorf("JoinErrors() = %q, missing %q", got.Error(), msg)
				}
			}
			if !errors.Is(got, first) && len(tt.wantMsgs) > 0 && tt.wantMsgs[0] == "vm1 failed" {
				t.Errorf("JoinErrors() does not wrap %v", first)
			}
		})
	}
}
