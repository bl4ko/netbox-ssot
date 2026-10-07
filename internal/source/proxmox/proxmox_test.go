package proxmox

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/parser"
)

func TestProxmoxOSTypeToPlatformName(t *testing.T) {
	tests := []struct {
		name   string
		osType string
		want   string
	}{
		{name: "linux 2.6 kernel", osType: "l26", want: "Other 2.6.x Linux (64-bit)"},
		{name: "windows 11", osType: "win11", want: "Windows 11"},
		{name: "windows 10", osType: "win10", want: "Windows 10"},
		{name: "empty type", osType: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := proxmoxOSTypeToPlatformName(tt.osType); got != tt.want {
				t.Errorf("proxmoxOSTypeToPlatformName(%q) = %q, want %q", tt.osType, got, tt.want)
			}
		})
	}
}

func TestInitTimesOutOnUnresponsiveAPI(t *testing.T) {
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

	ps := newTestSource(t, &parser.SourceConfig{
		HTTPScheme: parser.HTTP,
		Hostname:   serverURL.Hostname(),
		Port:       port,
		Username:   "root@pam",
		Password:   "secret",
		Timeout:    1,
	})
	done := make(chan error, 1)
	go func() { done <- ps.Init() }()
	select {
	case err := <-done:
		if err == nil {
			t.Errorf("Init() = nil, want a timeout error")
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("Init() did not return within 10s with a 1s timeout")
	}
}
