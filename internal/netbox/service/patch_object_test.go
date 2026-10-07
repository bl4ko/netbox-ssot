package service

import (
	"context"
	"log"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

func TestPatchObjectDryRunAppliesOnlyThePatchedFields(t *testing.T) {
	client := &NetboxClient{Logger: &logger.Logger{Logger: log.Default()}, DryRun: true}
	current := &objects.Device{
		NetboxObject: objects.NetboxObject{ID: 12, Description: "set by a source with priority"},
		Name:         "dry-run-merge",
		SerialNumber: "OLD-SERIAL",
	}
	target := &objects.Device{
		NetboxObject: objects.NetboxObject{Description: "lower priority description"},
		Name:         "dry-run-merge",
		SerialNumber: "NEW-SERIAL",
	}

	body := map[string]interface{}{"serial": "NEW-SERIAL"}
	got, err := PatchObject(context.Background(), client, 12, body, current, target)
	if err != nil {
		t.Fatalf("PatchObject() error = %v", err)
	}
	if got.ID != 12 || got.SerialNumber != "NEW-SERIAL" || got.Description != "set by a source with priority" {
		t.Errorf("PatchObject() = %+v, want device 12 with serial NEW-SERIAL and its former description", got)
	}
}
