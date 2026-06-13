package ovirt

import (
	"context"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func newTestOVirtSource(t *testing.T, interfaceFilter string) *OVirtSource {
	t.Helper()
	testLogger, err := logger.New("", 1)
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	return &OVirtSource{
		Config: common.Config{
			Logger: testLogger,
			Ctx:    context.WithValue(context.Background(), constants.CtxSourceKey, "test"),
			SourceConfig: &parser.SourceConfig{
				Name:            "test",
				InterfaceFilter: interfaceFilter,
			},
		},
	}
}

func newTestNic(name, id, mac string) *ovirtsdk4.Nic {
	return ovirtsdk4.NewNicBuilder().
		Name(name).
		Id(id).
		Mac(ovirtsdk4.NewMacBuilder().Address(mac).MustBuild()).
		MustBuild()
}

func TestCollectVMNicData(t *testing.T) {
	o := newTestOVirtSource(t, "docker.*")
	vm := ovirtsdk4.NewVmBuilder().
		Name("test-vm").
		NicsOfAny(
			newTestNic("nic1", "nic-id-1", "56:6f:be:6a:03:21"),
			newTestNic("eth1", "nic-id-2", "56:6f:be:6a:03:26"),
			newTestNic("docker0", "nic-id-3", "56:6f:be:6a:03:27"),
		).
		MustBuild()

	nicsData, err := o.collectVMNicData(nil, vm)
	if err != nil {
		t.Fatalf("collectVMNicData() error = %v", err)
	}
	if len(nicsData) != 2 {
		t.Fatalf("collectVMNicData() returned %d nics, want 2 (docker0 filtered out)", len(nicsData))
	}
	want := []struct {
		name string
		id   string
		mac  string
	}{
		{"nic1", "nic-id-1", "56:6F:BE:6A:03:21"},
		{"eth1", "nic-id-2", "56:6F:BE:6A:03:26"},
	}
	for i, w := range want {
		if nicsData[i].name != w.name {
			t.Errorf("nicsData[%d].name = %q, want %q", i, nicsData[i].name, w.name)
		}
		if nicsData[i].id != w.id {
			t.Errorf("nicsData[%d].id = %q, want %q", i, nicsData[i].id, w.id)
		}
		if nicsData[i].mac != w.mac {
			t.Errorf("nicsData[%d].mac = %q, want %q (uppercase)", i, nicsData[i].mac, w.mac)
		}
	}
}

func TestCollectVMNicDataNoNics(t *testing.T) {
	o := newTestOVirtSource(t, "")
	vm := ovirtsdk4.NewVmBuilder().Name("test-vm").MustBuild()

	nicsData, err := o.collectVMNicData(nil, vm)
	if err != nil {
		t.Fatalf("collectVMNicData() error = %v", err)
	}
	if len(nicsData) != 0 {
		t.Fatalf("collectVMNicData() returned %d nics, want 0", len(nicsData))
	}
}
