package service

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	servicetest "github.com/bl4ko/netbox-ssot/internal/netbox/service/testing"
)

// func TestGetAll(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		api         *NetboxClient
// 		extraParams string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []objects.Tag
// 		wantErr bool
// 	}{
// 		{
// 			name: "TestGetAll_Success",
// 			args: args{
// 				ctx:         context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
// 				api:         servicetest.MockNetboxClient,
// 				extraParams: "",
// 			},
// 			// See predefined values in api_test for mockserver
// 			want: []objects.Tag{
// 				{
// 					ID:          0,
// 					Name:        "Source: proxmox",
// 					Slug:        "source-proxmox",
// 					Color:       "9e9e9e",
// 					Description: "Automatically created tag by netbox-ssot for source proxmox",
// 				},
// 				{
// 					ID:          1,
// 					Name:        "netbox-ssot",
// 					Slug:        "netbox-ssot",
// 					Color:       "00add8",
// 					Description: "Tag used by netbox-ssot to mark devices that are managed by it",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		mockServer := CreateMockServer()
// 		defer mockServer.Close()
// 		MockNetboxClient.BaseURL = mockServer.URL
// 		t.Run(tt.name, func(t *testing.T) {
// 			response, err := GetAll[objects.Tag](tt.args.ctx, tt.args.api, tt.args.extraParams)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			// Parse the object
// 			if !reflect.DeepEqual(response, tt.want) {
// 				t.Errorf("GetAll() = %v, want %v", response, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPatch(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		api      *NetboxClient
// 		objectID int
// 		body     map[string]interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *objects.Tag
// 		wantErr bool
// 	}{
// 		{
// 			name: "Test patch tag",
// 			args: args{
// 				ctx:      context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
// 				api:      servicetest.MockNetboxClient,
// 				objectID: 1,
// 				body: map[string]interface{}{
// 					"description": "new description",
// 				},
// 			},
// 			// See predefined values in api_test for mockserver
// 			want:    &servicetest.MockTagPatchResponse,
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		mockServer := CreateMockServer()
// 		defer mockServer.Close()
// 		MockNetboxClient.BaseURL = mockServer.URL
// 		t.Run(tt.name, func(t *testing.T) {
// 			response, err := Patch[objects.Tag](tt.args.ctx, tt.args.api, tt.args.objectID, tt.args.body)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if err != nil {
// 				t.Errorf("marshal tag patch response: %s", err)
// 			}
// 			if !reflect.DeepEqual(response, tt.want) {
// 				t.Errorf("Patch() = %v, want %v", response, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCreate(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		api    *NetboxClient
// 		object *objects.Tag
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *objects.Tag
// 		wantErr bool
// 	}{
// 		{
// 			name: "Test create tag",
// 			args: args{
// 				ctx:    context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
// 				api:    servicetest.MockNetboxClient,
// 				object: &servicetest.MockTagCreateResponse,
// 			},
// 			// See predefined values in api_test for mockserver
// 			want:    &servicetest.MockTagCreateResponse,
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		mockServer := servicetest.CreateMockServer()
// 		defer mockServer.Close()
// 		servicetest.MockNetboxClient.BaseURL = mockServer.URL
// 		t.Run(tt.name, func(t *testing.T) {
// 			response, err := Create(tt.args.ctx, tt.args.api, tt.args.object)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(response, tt.want) {
// 				t.Errorf("Patch() = %v, want %v", response, tt.want)
// 			}
// 		})
// 	}
// }

func TestBulkDeleteObjects(t *testing.T) {
	type args struct {
		ctx        context.Context
		objectPath string
		idSet      map[int]bool
		nbClient   *NetboxClient
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
				nbClient: &NetboxClient{
					Logger:     &logger.Logger{Logger: log.Default()},
					HTTPClient: &http.Client{},
					Timeout:    constants.DefaultAPITimeout,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := servicetest.CreateMockServer()
		defer mockServer.Close()
		tt.args.nbClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.nbClient.BulkDeleteObjects(tt.args.ctx, tt.args.objectPath, tt.args.idSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteObject(t *testing.T) {
	type args struct {
		ctx        context.Context
		objectPath string
		id         int
		nbClient   *NetboxClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test delete tag",
			args: args{
				ctx:        context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				objectPath: constants.TagsAPIPath,
				id:         1,
				nbClient: &NetboxClient{
					Logger:     &logger.Logger{Logger: log.Default()},
					HTTPClient: &http.Client{},
					Timeout:    constants.DefaultAPITimeout,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := servicetest.CreateMockServer()
		defer mockServer.Close()
		tt.args.nbClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.nbClient.DeleteObject(tt.args.ctx, tt.args.objectPath, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	type args struct {
		ctx         context.Context
		objectPath  string
		nbClient    *NetboxClient
		extraParams string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test get all tags",
			args: args{
				ctx:        context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				objectPath: constants.TagsAPIPath,
				nbClient: &NetboxClient{
					Logger:     &logger.Logger{Logger: log.Default()},
					HTTPClient: &http.Client{},
					Timeout:    constants.DefaultAPITimeout,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := servicetest.CreateMockServer()
		defer mockServer.Close()
		tt.args.nbClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAll[objects.Tag](tt.args.ctx, tt.args.nbClient, tt.args.objectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		ctx      context.Context
		nbClient *NetboxClient
		newTag   *objects.Tag
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test get all tags",
			args: args{
				ctx:    context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
				newTag: &objects.Tag{Name: "New Tag"},
				nbClient: &NetboxClient{
					Logger:     &logger.Logger{Logger: log.Default()},
					HTTPClient: &http.Client{},
					Timeout:    constants.DefaultAPITimeout,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mockServer := servicetest.CreateMockServer()
		defer mockServer.Close()
		tt.args.nbClient.BaseURL = mockServer.URL
		t.Run(tt.name, func(t *testing.T) {
			_, err := Create[objects.Tag](tt.args.ctx, tt.args.nbClient, tt.args.newTag)
			if (err != nil) != tt.wantErr {
				t.Errorf("create new tag error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	netboxMock := servicetest.CreateMockServer()
	defer netboxMock.Close()

	ctx := context.Background()
	defer ctx.Done()

	netboxClient, err := NewNetboxClient(&logger.Logger{Logger: log.Default()}, netboxMock.URL, "", false, 10, "")
	if err != nil {
		t.Errorf("new netbox client: %s", err)
	}

	version, err := GetVersion(ctx, netboxClient)
	if err != nil {
		t.Errorf("get version: %s", err)
	}

	expectedVersion := "4.0.0"
	if version != expectedVersion {
		t.Errorf("want: %s, got: %s", expectedVersion, version)
	}

}
