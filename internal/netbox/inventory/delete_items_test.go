package inventory

import (
	"context"
	"testing"
)

func TestNetboxInventory_DeleteOrphans(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		nbi     *NetboxInventory
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nbi.DeleteOrphans(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("NetboxInventory.DeleteOrphans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
