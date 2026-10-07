package vmware

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/vmware/govmomi/simulator"
)

// newSimulatedVCenter starts a vcsim vCenter and returns a source configured to reach it.
func newSimulatedVCenter(t *testing.T, sourceConfig *parser.SourceConfig) (*VmwareSource, *simulator.Service) {
	t.Helper()
	model := simulator.VPX()
	if err := model.Create(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(model.Remove)
	server := model.Service.NewServer()
	t.Cleanup(server.Close)
	port, err := strconv.Atoi(server.URL.Port())
	if err != nil {
		t.Fatal(err)
	}
	password, _ := server.URL.User.Password()
	sourceConfig.HTTPScheme = parser.HTTPScheme(server.URL.Scheme)
	sourceConfig.Hostname = server.URL.Hostname()
	sourceConfig.Port = port
	sourceConfig.Username = server.URL.User.Username()
	sourceConfig.Password = password
	return newTestSource(t, sourceConfig), model.Service
}

func TestInitTimesOutOnUnresponsiveTagsAPI(t *testing.T) {
	vc, service := newSimulatedVCenter(t, &parser.SourceConfig{Timeout: 1})
	release := make(chan struct{})
	defer close(release)
	service.ServeMux.HandleFunc("/rest/", func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	})

	done := make(chan error, 1)
	go func() { done <- vc.Init() }()
	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "tags") || !strings.Contains(err.Error(), "Client.Timeout") {
			t.Errorf("Init() error = %v, want a timeout error from the tags API", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatalf("Init() did not return within 15s with a 1s timeout")
	}
}
