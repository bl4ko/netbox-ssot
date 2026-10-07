package vmware

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25/mo"

	// Register the vCenter REST API (tags) on the simulator
	_ "github.com/vmware/govmomi/vapi/simulator"
)

// tagsSessionPath is the first REST call made to read vCenter tags.
const tagsSessionPath = "/rest/com/vmware/cis/session"

// newSimulatedVCenter starts a vcsim vCenter and returns a source configured to reach it.
func newSimulatedVCenter(t *testing.T, sourceConfig *parser.SourceConfig) (*VmwareSource, *simulator.Service) {
	t.Helper()
	model := simulator.VPX()
	if err := model.Create(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(model.Remove)
	model.Service.RegisterEndpoints = true
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
	service.ServeMux.HandleFunc(tagsSessionPath, func(_ http.ResponseWriter, r *http.Request) {
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

func TestInitSkipsTheTagsAPIWithIgnoreTags(t *testing.T) {
	vc, service := newSimulatedVCenter(t, &parser.SourceConfig{Timeout: 5, IgnoreTags: true})
	service.ServeMux.HandleFunc(tagsSessionPath, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "tags API unavailable", http.StatusServiceUnavailable)
	})

	if err := vc.Init(); err != nil {
		t.Errorf("Init() error = %v, want the tags API skipped with ignoreTags", err)
	}
}

func TestInitAndSyncAgainstSimulatedVCenter(t *testing.T) {
	nbi, _ := newTestInventory(t)
	vc, service := newSimulatedVCenter(t, &parser.SourceConfig{Timeout: 5})
	// Expected VMs and their NetBox clusters come from the simulator inventory: a VM is in
	// "<cluster> (<datacenter>)", or in the cluster named after its standalone host.
	simulatorMap := service.Context.Map
	vmClusters := map[string]string{}
	for _, entity := range simulatorMap.All("VirtualMachine") {
		vm := entity.(*simulator.VirtualMachine)
		host := simulatorMap.Get(*vm.Runtime.Host).(*simulator.HostSystem)
		clusterName := host.Name
		if cluster, ok := simulatorMap.Get(*host.Parent).(*simulator.ClusterComputeResource); ok {
			datacenter := simulatorMap.Get(*simulatorMap.Get(*cluster.Parent).(*simulator.Folder).Parent)
			clusterName = cluster.Name + " (" + datacenter.(*simulator.Datacenter).Name + ")"
		}
		vmClusters[vm.Name] = clusterName
	}
	if len(vmClusters) == 0 {
		t.Fatal("the simulator has no VM")
	}

	if err := vc.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := vc.Sync(nbi); err != nil {
		t.Fatalf("Sync() error = %v", err)
	}
	for vmName, clusterName := range vmClusters {
		cluster, ok := nbi.GetCluster(clusterName)
		if !ok {
			t.Errorf("cluster %s of vm %s not synced", clusterName, vmName)
			continue
		}
		if _, ok := nbi.GetVM(vmName, cluster.ID); !ok {
			t.Errorf("vm %s not synced in cluster %s", vmName, clusterName)
		}
	}
}

func TestInitLogsOutWhenItFailsAfterLogin(t *testing.T) {
	vc, service := newSimulatedVCenter(t, &parser.SourceConfig{Timeout: 5})
	service.ServeMux.HandleFunc(tagsSessionPath, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "tags API unavailable", http.StatusServiceUnavailable)
	})
	if err := vc.Init(); err == nil {
		t.Fatal("Init() = nil, want the tags API error")
	}

	ctx := context.Background()
	sc := vc.SourceConfig
	vcURL := &url.URL{
		Scheme: string(sc.HTTPScheme),
		Host:   net.JoinHostPort(sc.Hostname, strconv.Itoa(sc.Port)),
		Path:   "/sdk",
		User:   url.UserPassword(sc.Username, sc.Password),
	}
	client, err := govmomi.NewClient(ctx, vcURL, true)
	if err != nil {
		t.Fatal(err)
	}
	var sessionManager mo.SessionManager
	err = client.RetrieveOne(ctx, *client.ServiceContent.SessionManager, []string{"sessionList"}, &sessionManager)
	if err != nil {
		t.Fatal(err)
	}
	// Only the session of this check remains
	if n := len(sessionManager.SessionList); n != 1 {
		t.Errorf("open vCenter sessions = %d, want 1: Init left its session open", n)
	}
}

func TestInitHandlesPortgroupsWithoutNicTeaming(t *testing.T) {
	vc, service := newSimulatedVCenter(t, &parser.SourceConfig{Timeout: 5, IgnoreTags: true})
	patched := 0
	for _, entity := range service.Context.Map.All("HostSystem") {
		host := entity.(*simulator.HostSystem)
		if host.Config == nil || host.Config.Network == nil {
			continue
		}
		for i := range host.Config.Network.Portgroup {
			host.Config.Network.Portgroup[i].ComputedPolicy.NicTeaming = nil
			patched++
		}
	}
	if patched == 0 {
		t.Fatal("the simulator has no host port group")
	}

	if err := vc.Init(); err != nil {
		t.Errorf("Init() error = %v", err)
	}
}
