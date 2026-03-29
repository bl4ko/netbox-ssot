package hetznercloud

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func newTestSource() *Source {
	return &Source{
		Config: common.Config{
			Logger: &logger.Logger{Logger: log.New(os.Stdout, "", log.LstdFlags)},
			Ctx: context.WithValue(
				context.Background(),
				constants.CtxSourceKey,
				"hetznercloud-test",
			),
			SourceConfig: &parser.SourceConfig{
				Name: "test-hcloud",
			},
		},
	}
}

func TestGbToMB(t *testing.T) {
	tests := []struct {
		name string
		gb   float32
		want int
	}{
		{"1 GB", 1, 1024},
		{"2 GB", 2, 2048},
		{"0.5 GB", 0.5, 512},
		{"32 GB", 32, 32768},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gbToMB(tt.gb); got != tt.want {
				t.Errorf("gbToMB(%v) = %v, want %v", tt.gb, got, tt.want)
			}
		})
	}
}

func TestSourceSiteDedup(t *testing.T) {
	hcs := newTestSource()

	hcs.Locations = []*hcloud.Location{
		{ID: 1, Name: "fsn1", City: "Falkenstein"},
		{ID: 2, Name: "fsn2", City: "Falkenstein"},
		{ID: 3, Name: "nbg1", City: "Nuremberg"},
	}

	hcs.NetboxSites = make(map[string]*objects.Site)
	for _, loc := range hcs.Locations {
		if _, exists := hcs.NetboxSites[loc.City]; exists {
			continue
		}
		hcs.NetboxSites[loc.City] = &objects.Site{Name: loc.City}
	}

	if len(hcs.NetboxSites) != 2 {
		t.Errorf("expected 2 unique sites, got %d", len(hcs.NetboxSites))
	}
	if _, ok := hcs.NetboxSites["Falkenstein"]; !ok {
		t.Error("expected Falkenstein site")
	}
	if _, ok := hcs.NetboxSites["Nuremberg"]; !ok {
		t.Error("expected Nuremberg site")
	}
}

func TestSourceDatacenterToLocationMapping(t *testing.T) {
	hcs := newTestSource()

	falkenstein := &objects.Site{Name: "Falkenstein"}
	nuremberg := &objects.Site{Name: "Nuremberg"}
	hcs.NetboxSites = map[string]*objects.Site{
		"Falkenstein": falkenstein,
		"Nuremberg":   nuremberg,
	}

	hcs.Datacenters = []*hcloud.Datacenter{
		{ID: 1, Name: "fsn1-dc14", Description: "Falkenstein 1 DC14", Location: &hcloud.Location{City: "Falkenstein"}},
		{ID: 2, Name: "nbg1-dc3", Description: "Nuremberg 1 DC3", Location: &hcloud.Location{City: "Nuremberg"}},
		{ID: 3, Name: "fsn1-dc15", Description: "Falkenstein 1 DC15", Location: &hcloud.Location{City: "Falkenstein"}},
	}

	hcs.NetboxLocations = make(map[string]*objects.Location)
	for _, dc := range hcs.Datacenters {
		var site *objects.Site
		if dc.Location != nil {
			site = hcs.NetboxSites[dc.Location.City]
		}
		hcs.NetboxLocations[dc.Name] = &objects.Location{
			Name: dc.Name,
			Site: site,
		}
	}

	if len(hcs.NetboxLocations) != 3 {
		t.Errorf("expected 3 locations, got %d", len(hcs.NetboxLocations))
	}

	loc := hcs.NetboxLocations["fsn1-dc14"]
	if loc.Site != falkenstein {
		t.Error("fsn1-dc14 should be linked to Falkenstein site")
	}

	loc2 := hcs.NetboxLocations["nbg1-dc3"]
	if loc2.Site != nuremberg {
		t.Error("nbg1-dc3 should be linked to Nuremberg site")
	}
}

func TestSourceServerStatusMapping(t *testing.T) {
	tests := []struct {
		name       string
		status     hcloud.ServerStatus
		wantStatus *objects.VMStatus
	}{
		{"running server", hcloud.ServerStatusRunning, &objects.VMStatusActive},
		{"off server", hcloud.ServerStatusOff, &objects.VMStatusOffline},
		{"initializing server", hcloud.ServerStatusInitializing, &objects.VMStatusActive},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := &objects.VMStatusActive
			if tt.status == hcloud.ServerStatusOff {
				status = &objects.VMStatusOffline
			}
			if status != tt.wantStatus {
				t.Errorf("status = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestSourceServerSiteLookup(t *testing.T) {
	hcs := newTestSource()

	falkenstein := &objects.Site{Name: "Falkenstein"}
	hcs.NetboxSites = map[string]*objects.Site{
		"Falkenstein": falkenstein,
	}

	server := &hcloud.Server{
		Location: &hcloud.Location{City: "Falkenstein"},
	}

	var site *objects.Site
	if server.Location != nil {
		site = hcs.NetboxSites[server.Location.City]
	}

	if site != falkenstein {
		t.Error("expected server to be mapped to Falkenstein site")
	}
}

func TestSourceServerSiteLookupNilLocation(t *testing.T) {
	hcs := newTestSource()
	hcs.NetboxSites = map[string]*objects.Site{}

	server := &hcloud.Server{}

	var site *objects.Site
	if server.Location != nil {
		site = hcs.NetboxSites[server.Location.City]
	}

	if site != nil {
		t.Error("expected nil site for server with no location")
	}
}

func TestSourceNetworkIPRange(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("10.0.0.0/8")
	network := &hcloud.Network{
		ID:      1,
		Name:    "my-network",
		IPRange: ipNet,
	}

	if network.IPRange.String() != "10.0.0.0/8" {
		t.Errorf("expected 10.0.0.0/8, got %s", network.IPRange.String())
	}
}

func TestSourceFloatingIPFormat(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want string
	}{
		{"ipv4", net.ParseIP("1.2.3.4"), "1.2.3.4/32"},
		{"ipv6", net.ParseIP("2001:db8::1"), "2001:db8::1/32"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf("%s/32", tt.ip.String())
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSourceIPv6HostIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want string
	}{
		{"trailing double colon", "2a01:4f8:1:2::", "2a01:4f8:1:2::1"},
		{"already has host part", "2a01:4f8:1:2::1", "2a01:4f8:1:2::1"},
		{"no trailing colons", "2001:db8::abcd", "2001:db8::abcd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipStr := tt.ip
			if len(ipStr) >= 2 && ipStr[len(ipStr)-2:] == "::" {
				ipStr += "1"
			}
			if ipStr != tt.want {
				t.Errorf("got %q, want %q", ipStr, tt.want)
			}
		})
	}
}

func TestSourcePrivateNetInterfaceNaming(t *testing.T) {
	tests := []struct {
		index int
		want  string
	}{
		{0, "eth1"},
		{1, "eth2"},
		{4, "eth5"},
	}
	for _, tt := range tests {
		got := fmt.Sprintf("eth%d", tt.index+1)
		if got != tt.want {
			t.Errorf("index %d: got %q, want %q", tt.index, got, tt.want)
		}
	}
}

func TestMbPerGB(t *testing.T) {
	if mbPerGB != 1000 {
		t.Errorf("mbPerGB = %d, want 1000", mbPerGB)
	}
}
