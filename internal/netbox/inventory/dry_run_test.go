package inventory

import (
	"context"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
)

func TestDryRunPatchKeepsTheUpdatedObjectIndexed(t *testing.T) {
	ssotTag := &objects.Tag{ID: 1, Name: constants.SsotTagName}
	site := &objects.Site{NetboxObject: objects.NetboxObject{ID: 3}, Name: "dry-run-site"}
	cluster := &objects.Cluster{NetboxObject: objects.NetboxObject{ID: 4}, Name: "dry-run-cluster"}
	oldDevice := &objects.Device{
		NetboxObject: objects.NetboxObject{ID: 10, Tags: []*objects.Tag{ssotTag}},
		Name:         "dry-run-host",
		Site:         site,
		Cluster:      cluster,
		SerialNumber: "OLD-SERIAL",
	}
	nbi := &NetboxInventory{
		Logger:                      mockLogger,
		SsotTag:                     ssotTag,
		NetboxAPI:                   &service.NetboxClient{Logger: mockLogger, DryRun: true},
		OrphanManager:               NewOrphanManager(mockLogger),
		devicesIndexByNameAndSiteID: map[string]map[int]*objects.Device{"dry-run-host": {3: oldDevice}},
		devicesIndexByID:            map[int]*objects.Device{10: oldDevice},
	}
	ctx := context.WithValue(context.Background(), constants.CtxSourceKey, "test")

	got, err := nbi.AddDevice(ctx, &objects.Device{
		NetboxObject: objects.NetboxObject{Tags: []*objects.Tag{ssotTag}},
		Name:         "dry-run-host",
		Site:         site,
		Cluster:      cluster,
		SerialNumber: "NEW-SERIAL",
	})
	if err != nil {
		t.Fatalf("AddDevice() error = %v", err)
	}

	indexed, _ := nbi.GetDevice("dry-run-host", site.ID)
	for name, device := range map[string]*objects.Device{"returned": got, "indexed": indexed} {
		if device == nil || device.ID != 10 || device.Name != "dry-run-host" || device.Cluster != cluster ||
			device.SerialNumber != "NEW-SERIAL" {
			t.Errorf("%s device = %+v, want device 10 dry-run-host in cluster 4 with serial NEW-SERIAL", name, device)
		}
	}
}
