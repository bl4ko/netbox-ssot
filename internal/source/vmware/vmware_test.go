package vmware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/vmware/govmomi/vim25/types"
)

func TestVmwareConnectionStateToDeviceStatus(t *testing.T) {
	tests := []struct {
		name  string
		state types.HostSystemConnectionState
		want  *objects.DeviceStatus
	}{
		{name: "connected", state: "connected", want: &objects.DeviceStatusActive},
		{name: "disconnected", state: "disconnected", want: &objects.DeviceStatusOffline},
		{name: "not responding", state: "notResponding", want: &objects.DeviceStatusOffline},
		{name: "empty", state: "", want: &objects.DeviceStatusOffline},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vmwareConnectionStateToDeviceStatus(tt.state)
			if *got != *tt.want {
				t.Errorf("vmwareConnectionStateToDeviceStatus(%q) = %v, want %v", tt.state, *got, *tt.want)
			}
		})
	}
}

func TestInitTimesOutOnUnresponsiveVCenter(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	defer server.Close()
	defer close(release)
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(serverURL.Port())
	if err != nil {
		t.Fatal(err)
	}

	vc := newTestSource(t, &parser.SourceConfig{
		HTTPScheme: parser.HTTP,
		Hostname:   serverURL.Hostname(),
		Port:       port,
		Username:   "administrator@vsphere.local",
		Password:   "secret",
		Timeout:    1,
	})
	done := make(chan error, 1)
	go func() { done <- vc.Init() }()
	select {
	case err := <-done:
		if err == nil {
			t.Errorf("Init() = nil, want a timeout error")
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("Init() did not return within 10s with a 1s timeout")
	}
}
