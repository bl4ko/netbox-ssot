package proxmox

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
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
		// Names are the VMware guest full names that netbox-sync writes for vCenter VMs
		// (Broadcom KB 321876), so that both tools share the same platforms.
		{name: "linux 2.6 kernel", osType: "l26", want: "Other 2.6.x Linux (64-bit)"},
		{name: "windows xp", osType: "wxp", want: "Microsoft Windows XP (32-bit)"},
		{name: "windows 2000", osType: "w2k", want: "Microsoft Windows 2000 Server"},
		{name: "windows 2003", osType: "w2k3", want: "Microsoft Windows Server 2003 Standard (32-bit)"},
		{name: "windows 2008", osType: "w2k8", want: "Microsoft Windows Server 2008 (64-bit)"},
		{name: "windows vista", osType: "wvista", want: "Microsoft Windows Vista (32-bit)"},
		{name: "windows 7", osType: "win7", want: "Microsoft Windows 7 (64-bit)"},
		{name: "windows 8", osType: "win8", want: "Microsoft Windows 8 (64-bit)"},
		{name: "windows 10", osType: "win10", want: "Microsoft Windows 10 (64-bit)"},
		{name: "windows 11", osType: "win11", want: "Microsoft Windows 11 (64-bit)"},
		{name: "unspecified os", osType: "other", want: ""},
		{name: "linux 2.4 kernel without known name", osType: "l24", want: ""},
		{name: "solaris without known version", osType: "solaris", want: ""},
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

// newFakeProxmoxAPI serves one node "n1" with the given VMs (vmid -> status) and
// records the VMIDs whose guest agent is queried.
func newFakeProxmoxAPI(t *testing.T, vmStatuses map[int]string) (*httptest.Server, func() []int) {
	t.Helper()
	var mu sync.Mutex
	agentCalls := []int{}
	respond := func(w http.ResponseWriter, data string) {
		_, _ = fmt.Fprintf(w, `{"data":%s}`, data)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api2/json")
		var vmid int
		switch {
		case path == "/access/ticket":
			respond(w, `{"ticket":"t","CSRFPreventionToken":"c","username":"root@pam"}`)
		case path == "/cluster/status":
			respond(w, `[]`)
		case path == "/nodes":
			respond(w, `[{"node":"n1","status":"online"}]`)
		case path == "/nodes/n1/status":
			respond(w, `{}`)
		case path == "/nodes/n1/network", path == "/nodes/n1/lxc":
			respond(w, `[]`)
		case path == "/nodes/n1/qemu":
			vms := make([]string, 0, len(vmStatuses))
			for id, status := range vmStatuses {
				vms = append(vms, fmt.Sprintf(`{"vmid":%d,"name":"vm%d","status":%q}`, id, id, status))
			}
			respond(w, "["+strings.Join(vms, ",")+"]")
		case scanPath(path, "/nodes/n1/qemu/%d/status/current", &vmid):
			respond(w, fmt.Sprintf(`{"vmid":%d,"name":"vm%d","status":%q}`, vmid, vmid, vmStatuses[vmid]))
		case scanPath(path, "/nodes/n1/qemu/%d/config", &vmid):
			respond(w, `{}`)
		case scanPath(path, "/nodes/n1/qemu/%d/agent/network-get-interfaces", &vmid):
			mu.Lock()
			agentCalls = append(agentCalls, vmid)
			mu.Unlock()
			respond(w, `{"result":[{"name":"eth0"}]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	return server, func() []int {
		mu.Lock()
		defer mu.Unlock()
		out := make([]int, len(agentCalls))
		copy(out, agentCalls)
		return out
	}
}

// scanPath reports whether path matches format, an fmt.Sscanf pattern with one %d.
func scanPath(path, format string, vmid *int) bool {
	n, err := fmt.Sscanf(path, format, vmid)
	return err == nil && n == 1 && fmt.Sprintf(format, *vmid) == path
}

func TestInitQueriesGuestAgentOfRunningVMsOnly(t *testing.T) {
	server, agentCalls := newFakeProxmoxAPI(t, map[int]string{101: "running", 102: "stopped"})
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
		Timeout:    5,
	})

	if err := ps.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if got := agentCalls(); len(got) != 1 || got[0] != 101 {
		t.Errorf("guest agent queried for VMIDs %v, want only the running VM 101", got)
	}
	if _, known := ps.VMIfaces[101]; !known {
		t.Errorf("interfaces of running VM 101 are unknown, want them read from the agent")
	}
	if _, known := ps.VMIfaces[102]; known {
		t.Errorf("interfaces of stopped VM 102 are known, want them unknown")
	}
}
