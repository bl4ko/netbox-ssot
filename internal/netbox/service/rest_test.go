package service

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestGetAll(t *testing.T) {
	type args struct {
		ctx         context.Context
		api         *NetboxClient
		extraParams string
	}
	tests := []struct {
		name    string
		args    args
		want    []objects.Tag
		wantErr bool
	}{
		{
			name: "TestGetAll_Success",
			args: args{
				ctx:         context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				api:         MockNetboxClient,
				extraParams: "",
			},
			// See predefined values in api_test for mockserver
			want: []objects.Tag{
				{
					ID:          0,
					Name:        "Source: proxmox",
					Slug:        "source-proxmox",
					Color:       "9e9e9e",
					Description: "Automatically created tag by netbox-ssot for source proxmox",
				},
				{
					ID:          1,
					Name:        "netbox-ssot",
					Slug:        "netbox-ssot",
					Color:       "00add8",
					Description: "Tag used by netbox-ssot to mark devices that are managed by it",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
		MockNetboxClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			response, err := GetAll[objects.Tag](tt.args.ctx, tt.args.api, tt.args.extraParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Parse the object
			if !reflect.DeepEqual(response, tt.want) {
				t.Errorf("GetAll() = %v, want %v", response, tt.want)
			}
		})
	}
}

func TestPatch(t *testing.T) {
	type args struct {
		ctx      context.Context
		api      *NetboxClient
		objectID int
		body     map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test patch tag",
			args: args{
				ctx:      context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				api:      MockNetboxClient,
				objectID: 1,
				body: map[string]interface{}{
					"description": "new description",
				},
			},
			// See predefined values in api_test for mockserver
			want:    TagPatchResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
		MockNetboxClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			response, err := Patch[objects.Tag](tt.args.ctx, tt.args.api, tt.args.objectID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Parse the object
			var wantTag objects.Tag
			err = json.Unmarshal([]byte(tt.want), &wantTag)
			if err != nil {
				t.Errorf("marshal tag patch response: %s", err)
			}
			if !reflect.DeepEqual(response, &wantTag) {
				t.Errorf("Patch() = %v, want %v", response, wantTag)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		ctx    context.Context
		api    *NetboxClient
		object string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test create tag",
			args: args{
				ctx:    context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				api:    MockNetboxClient,
				object: TagCreateResponse,
			},
			// See predefined values in api_test for mockserver
			want:    TagPatchResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
		MockNetboxClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			var newTag objects.Tag
			err := json.Unmarshal([]byte(tt.args.object), &newTag)
			if err != nil {
				t.Errorf("unmarshal tag: %s", err)
			}
			response, err := Create[objects.Tag](tt.args.ctx, tt.args.api, &newTag)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(response, &newTag) {
				t.Errorf("Patch() = %v, want %v", response, newTag)
			}
		})
	}
}

func TestBulkDeleteObjects(t *testing.T) {
	type args struct {
		ctx        context.Context
		objectPath string
		idSet      map[int]bool
		api        *NetboxClient
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test bulk delete tags",
			args: args{
				ctx:        context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				objectPath: constants.TagsAPIPath,
				idSet:      map[int]bool{0: true, 1: true},
				api:        MockNetboxClient,
			},
			// See predefined values in api_test for mockserver
			want:    TagPatchResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
		MockNetboxClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.api.BulkDeleteObjects(tt.args.ctx, tt.args.objectPath, tt.args.idSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
