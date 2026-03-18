package service

import (
	"context"
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
				ctx: context.WithValue(
					context.Background(),
					constants.CtxSourceKey,
					"test",
				),
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
		want    *objects.Tag
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
			want:    &MockTagPatchResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
		MockNetboxClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			response, err := Patch[objects.Tag](
				tt.args.ctx,
				tt.args.api,
				tt.args.objectID,
				tt.args.body,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Errorf("marshal tag patch response: %s", err)
			}
			if !reflect.DeepEqual(response, tt.want) {
				t.Errorf("Patch() = %v, want %v", response, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		ctx    context.Context
		api    *NetboxClient
		object *objects.Tag
	}
	tests := []struct {
		name    string
		args    args
		want    *objects.Tag
		wantErr bool
	}{
		{
			name: "Test create tag",
			args: args{
				ctx:    context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				api:    MockNetboxClient,
				object: &MockTagCreateResponse,
			},
			// See predefined values in api_test for mockserver
			want:    &MockTagCreateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
		MockNetboxClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			response, err := Create(tt.args.ctx, tt.args.api, tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(response, tt.want) {
				t.Errorf("Patch() = %v, want %v", response, tt.want)
			}
		})
	}
}

func TestGetAll_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
	_, err := GetAll[objects.Tag](ctx, FailingMockNetboxClient, "")
	if err == nil {
		t.Error("expected error from failing client, got nil")
	}
}

func TestPatch_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
	_, err := Patch[objects.Tag](ctx, FailingMockNetboxClient, 1, map[string]interface{}{"name": "x"})
	if err == nil {
		t.Error("expected error from failing client, got nil")
	}
}

func TestCreate_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
	tag := &objects.Tag{Name: "test"}
	_, err := Create(ctx, FailingMockNetboxClient, tag)
	if err == nil {
		t.Error("expected error from failing client, got nil")
	}
}

func TestBulkDeleteObjects_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
	err := FailingMockNetboxClient.BulkDeleteObjects(ctx, constants.TagsAPIPath, map[int]bool{1: true})
	if err == nil {
		t.Error("expected error from failing client, got nil")
	}
}

func TestGetAll_NotFound(t *testing.T) {
	mockServer := CreateMockServer()
	defer mockServer.Close()
	MockNetboxClient.BaseURL = mockServer.URL
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
	_, err := GetAll[objects.Site](ctx, MockNetboxClient, "")
	if err != nil {
		t.Errorf("GetAll for sites should succeed with mock server, got %v", err)
	}
}

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name         string
		netboxClient *NetboxClient
		wantErr      bool
	}{
		{
			name:         "success returns version string",
			netboxClient: MockNetboxClient,
			wantErr:      false,
		},
		{
			name:         "failing client returns error",
			netboxClient: FailingMockNetboxClient,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
		MockNetboxClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			_, err := GetVersion(ctx, tt.netboxClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteObject(t *testing.T) {
	tests := []struct {
		name         string
		netboxClient *NetboxClient
		item         objects.IDItem
		wantErr      bool
	}{
		{
			name:         "success deleting a tag",
			netboxClient: MockNetboxClient,
			item:         &objects.Tag{ID: 1},
			wantErr:      false,
		},
		{
			name:         "failing client returns error",
			netboxClient: FailingMockNetboxClient,
			item:         &objects.Tag{ID: 1},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
		MockNetboxClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")
			err := tt.netboxClient.DeleteObject(ctx, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBulkDeleteObjects(t *testing.T) {
	type args struct {
		ctx        context.Context
		objectPath constants.APIPath
		idSet      map[int]bool
		api        *NetboxClient
	}
	tests := []struct {
		name    string
		args    args
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
