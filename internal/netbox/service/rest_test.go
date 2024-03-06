package service

import (
	"context"
	"reflect"
	"testing"

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
		// {
		// 	name: "TestGetAll_Success",
		// 	args: args{
		// 		ctx:         context.TODO(),
		// 		api:         MockNetboxClient,
		// 		extraParams: "",
		// 	},
		// 	// See predefined values in api_test for mockserver
		// 	want: []objects.Tag{
		// 		{
		// 			ID:          34,
		// 			Name:        "Source: proxmox",
		// 			Slug:        "source-proxmox",
		// 			Color:       "9e9e9e",
		// 			Description: "Automatically created tag by netbox-ssot for source proxmox",
		// 		},
		// 		{
		// 			ID:          1,
		// 			Name:        "netbox-ssot",
		// 			Slug:        "netbox-ssot",
		// 			Color:       "00add8",
		// 			Description: "Tag used by netbox-ssot to mark devices that are managed by it",
		// 		},
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		mockServer := CreateMockServer()
		defer mockServer.Close()
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
