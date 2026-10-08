package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
)

// failingSource is a source whose Init always fails.
type failingSource struct{}

func (failingSource) Init() error                           { return errors.New("init failed") }
func (failingSource) Sync(*inventory.NetboxInventory) error { return nil }

func TestSyncSourcesCollectsConcurrentFailures(t *testing.T) {
	ssotLogger, err := logger.New("", logger.ERROR)
	if err != nil {
		t.Fatal(err)
	}
	const sourceCount = 64
	sources := make([]namedSource, 0, sourceCount)
	for i := range sourceCount {
		name := fmt.Sprintf("source-%d", i)
		ctx := context.WithValue(context.Background(), constants.CtxSourceKey, name)
		sources = append(sources, namedSource{name: name, ctx: ctx, source: failingSource{}})
	}

	encounteredErrors := syncSources(ssotLogger, nil, sources)

	if len(encounteredErrors) != sourceCount {
		t.Errorf("syncSources() reported %d failed sources, want %d", len(encounteredErrors), sourceCount)
	}
}
