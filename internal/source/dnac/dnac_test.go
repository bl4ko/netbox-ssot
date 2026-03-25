package dnac

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	dnac "github.com/cisco-en-programmability/dnacenter-go-sdk/v8/sdk"
)

func newTestDnacSource() *DnacSource {
	return &DnacSource{
		Config: common.Config{
			Logger: &logger.Logger{Logger: log.New(os.Stdout, "", log.LstdFlags)},
			Ctx: context.WithValue(
				context.Background(),
				constants.CtxSourceKey,
				"dnac-test",
			),
			SourceConfig: &parser.SourceConfig{},
		},
	}
}

func TestGetInterfaceDuplex(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name   string
		duplex string
		want   *objects.InterfaceDuplex
	}{
		{"empty string", "", nil},
		{"FullDuplex", "FullDuplex", &objects.DuplexFull},
		{"AutoNegotiate", "AutoNegotiate", &objects.DuplexAuto},
		{"HalfDuplex", "HalfDuplex", &objects.DuplexHalf},
		{"unknown value", "SomeOther", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ds.getInterfaceDuplex(tt.duplex)
			if got != tt.want {
				t.Errorf("getInterfaceDuplex(%q) = %v, want %v", tt.duplex, got, tt.want)
			}
		})
	}
}

func TestGetInterfaceStatus(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name    string
		status  string
		want    bool
		wantErr bool
	}{
		{"down", "down", false, false},
		{"up", "up", true, false},
		{"unknown", "error", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ds.getInterfaceStatus(tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInterfaceStatus(%q) error = %v, wantErr %v", tt.status, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getInterfaceStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestGetInterfaceType(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name    string
		iType   string
		speed   int
		want    *objects.InterfaceType
		wantErr bool
	}{
		{
			name:  "Physical with known speed",
			iType: "Physical",
			speed: int(objects.GBPS1),
			want:  &objects.GE1FixedInterfaceType,
		},
		{
			name:  "Physical with unknown speed",
			iType: "Physical",
			speed: 999,
			want:  &objects.OtherInterfaceType,
		},
		{
			name:  "Virtual",
			iType: "Virtual",
			speed: 0,
			want:  &objects.VirtualInterfaceType,
		},
		{
			name:    "Unknown type",
			iType:   "Wireless",
			speed:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ds.getInterfaceType(tt.iType, tt.speed)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInterfaceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getInterfaceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateInterfaceName(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name    string
		iName   string
		iID     string
		filter  string
		wantErr bool
	}{
		{"empty name", "", "iface-1", "", true},
		{"valid name no filter", "eth0", "iface-1", "", false},
		{"filtered name", "eth0", "iface-1", "^eth.*", true},
		{"non-matching filter", "GigabitEthernet0/0", "iface-2", "^eth.*", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds.SourceConfig.InterfaceFilter = tt.filter
			err := ds.validateInterfaceName(tt.iName, tt.iID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInterfaceName(%q) error = %v, wantErr %v", tt.iName, err, tt.wantErr)
			}
		})
	}
}

func TestGetVlanModeAndAccessVlan(t *testing.T) {
	ds := newTestDnacSource()

	tests := []struct {
		name     string
		portMode string
		vlanID   string
		wantMode *objects.InterfaceMode
		wantErr  bool
	}{
		{
			name:     "access mode",
			portMode: "access",
			vlanID:   "100",
			wantMode: &objects.InterfaceModeAccess,
		},
		{
			name:     "trunk mode",
			portMode: "trunk",
			vlanID:   "1",
			wantMode: &objects.InterfaceModeTagged,
		},
		{
			name:     "dynamic_auto mode",
			portMode: "dynamic_auto",
			vlanID:   "1",
			wantMode: nil,
		},
		{
			name:     "routed mode",
			portMode: "routed",
			vlanID:   "1",
			wantMode: nil,
		},
		{
			name:     "unknown mode",
			portMode: "unknown",
			vlanID:   "1",
			wantErr:  true,
		},
		{
			name:     "invalid vlan ID",
			portMode: "access",
			vlanID:   "abc",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, _, err := ds.getVlanModeAndAccessVlan(tt.portMode, tt.vlanID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVlanModeAndAccessVlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && mode != tt.wantMode {
				t.Errorf("getVlanModeAndAccessVlan() mode = %v, want %v", mode, tt.wantMode)
			}
		})
	}
}

func TestGetDevice(t *testing.T) {
	ds := newTestDnacSource()

	expectedDevice := &objects.Device{
		NetboxObject: objects.NetboxObject{Tags: []*objects.Tag{{Name: "test"}}},
		Name:         "switch-01",
	}
	ds.DeviceID2nbDevice.Store("device-1", expectedDevice)
	ds.DeviceID2nbDevice.Store("device-bad-type", "not-a-device")

	tests := []struct {
		name     string
		deviceID string
		wantName string
		wantErr  bool
	}{
		{"existing device", "device-1", "switch-01", false},
		{"non-existent device", "device-999", "", true},
		{"wrong type in map", "device-bad-type", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ds.getDevice(tt.deviceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDevice(%q) error = %v, wantErr %v", tt.deviceID, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.wantName {
				t.Errorf("getDevice(%q).Name = %q, want %q", tt.deviceID, got.Name, tt.wantName)
			}
		})
	}
}

func TestGetVlanModeAndAccessVlan_WithStoredVlan(t *testing.T) {
	ds := newTestDnacSource()

	storedVlan := &objects.Vlan{
		NetboxObject: objects.NetboxObject{Tags: []*objects.Tag{{Name: "test"}}},
		Name:         "VLAN100",
		Vid:          100,
	}
	ds.VID2nbVlan.Store(100, storedVlan)

	mode, vlan, err := ds.getVlanModeAndAccessVlan("access", "100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != &objects.InterfaceModeAccess {
		t.Errorf("mode = %v, want access", mode)
	}
	if vlan != storedVlan {
		t.Errorf("vlan = %v, want stored vlan %v", vlan, storedVlan)
	}
}

func TestGetVlanModeAndAccessVlan_AccessWithoutStoredVlan(t *testing.T) {
	ds := newTestDnacSource()

	mode, vlan, err := ds.getVlanModeAndAccessVlan("access", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != &objects.InterfaceModeAccess {
		t.Errorf("mode = %v, want access", mode)
	}
	if vlan != nil {
		t.Errorf("vlan = %v, want nil (no stored vlan for VID 999)", vlan)
	}
}

func TestDnacSourceSiteRelations(t *testing.T) {
	ds := newTestDnacSource()

	sites := map[string]dnac.ResponseSitesGetSiteResponse{
		"site-1": {ID: "site-1", Name: "HQ", ParentID: ""},
		"site-2": {ID: "site-2", Name: "Branch", ParentID: "site-1"},
		"site-3": {ID: "site-3", Name: "Floor1", ParentID: "site-2"},
	}

	ds.Sites = sites
	ds.Site2Parent = make(map[string]string, len(sites))
	for _, site := range sites {
		ds.Site2Parent[site.ID] = site.ParentID
	}

	tests := []struct {
		siteID         string
		wantParentID   string
		wantParentName string
	}{
		{"site-1", "", ""},
		{"site-2", "site-1", "HQ"},
		{"site-3", "site-2", "Branch"},
	}
	for _, tt := range tests {
		t.Run(tt.siteID, func(t *testing.T) {
			parentID := ds.Site2Parent[tt.siteID]
			if parentID != tt.wantParentID {
				t.Errorf("Site2Parent[%q] = %q, want %q", tt.siteID, parentID, tt.wantParentID)
			}
			if tt.wantParentID != "" {
				parent := ds.Sites[parentID]
				if parent.Name != tt.wantParentName {
					t.Errorf("parent name = %q, want %q", parent.Name, tt.wantParentName)
				}
			}
		})
	}
}

func TestDnacSourceDeviceToSiteMapping(t *testing.T) {
	ds := newTestDnacSource()

	ds.Site2Devices = map[string]map[string]bool{
		"site-1": {"dev-a": true, "dev-b": true},
		"site-2": {"dev-c": true},
	}
	ds.Device2Site = map[string]string{
		"dev-a": "site-1",
		"dev-b": "site-1",
		"dev-c": "site-2",
	}

	tests := []struct {
		name     string
		deviceID string
		wantSite string
	}{
		{"dev-a in site-1", "dev-a", "site-1"},
		{"dev-b in site-1", "dev-b", "site-1"},
		{"dev-c in site-2", "dev-c", "site-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			siteID := ds.Device2Site[tt.deviceID]
			if siteID != tt.wantSite {
				t.Errorf("Device2Site[%q] = %q, want %q", tt.deviceID, siteID, tt.wantSite)
			}
			if !ds.Site2Devices[siteID][tt.deviceID] {
				t.Errorf("Site2Devices[%q] does not contain %q", siteID, tt.deviceID)
			}
		})
	}
}

func TestDnacSourceVlanStorage(t *testing.T) {
	ds := newTestDnacSource()

	vlanNum10 := 10
	vlanNum20 := 20
	ds.Vlans = map[int]dnac.ResponseDevicesGetDeviceInterfaceVLANsResponse{
		10: {
			InterfaceName:  "Vlan10",
			VLANType:       "access",
			VLANNumber:     &vlanNum10,
			NetworkAddress: "10.0.10.0",
			Prefix:         "24",
		},
		20: {
			InterfaceName:  "Vlan20",
			VLANType:       "trunk",
			VLANNumber:     &vlanNum20,
			NetworkAddress: "10.0.20.0",
			Prefix:         "24",
		},
	}

	if len(ds.Vlans) != 2 {
		t.Fatalf("expected 2 vlans, got %d", len(ds.Vlans))
	}

	vlan10 := ds.Vlans[10]
	if vlan10.InterfaceName != "Vlan10" {
		t.Errorf("Vlan 10 name = %q, want Vlan10", vlan10.InterfaceName)
	}
	if vlan10.NetworkAddress != "10.0.10.0" {
		t.Errorf("Vlan 10 network = %q, want 10.0.10.0", vlan10.NetworkAddress)
	}
	if vlan10.Prefix != "24" {
		t.Errorf("Vlan 10 prefix = %q, want 24", vlan10.Prefix)
	}
}

func TestDnacSourceInterfaceToDeviceMapping(t *testing.T) {
	ds := newTestDnacSource()

	ds.Interfaces = map[string]dnac.ResponseDevicesGetAllInterfacesResponse{
		"if-1": {ID: "if-1", DeviceID: "dev-a", PortName: "GigabitEthernet0/0"},
		"if-2": {ID: "if-2", DeviceID: "dev-a", PortName: "GigabitEthernet0/1"},
		"if-3": {ID: "if-3", DeviceID: "dev-b", PortName: "TenGigabitEthernet1/0"},
	}
	ds.DeviceID2InterfaceIDs = map[string][]string{
		"dev-a": {"if-1", "if-2"},
		"dev-b": {"if-3"},
	}

	if len(ds.DeviceID2InterfaceIDs["dev-a"]) != 2 {
		t.Errorf("dev-a should have 2 interfaces, got %d", len(ds.DeviceID2InterfaceIDs["dev-a"]))
	}
	if len(ds.DeviceID2InterfaceIDs["dev-b"]) != 1 {
		t.Errorf("dev-b should have 1 interface, got %d", len(ds.DeviceID2InterfaceIDs["dev-b"]))
	}

	iface := ds.Interfaces["if-3"]
	if iface.PortName != "TenGigabitEthernet1/0" {
		t.Errorf("if-3 PortName = %q, want TenGigabitEthernet1/0", iface.PortName)
	}
}

func TestDnacSourceWirelessLANMappings(t *testing.T) {
	ds := newTestDnacSource()

	isEnabled := true
	ds.WirelessLANInterfaceName2VlanID = map[string]int{
		"wlan-iface-1": 100,
		"wlan-iface-2": 200,
	}
	ds.SSID2WlanGroupName = map[string]string{
		"Corp-WiFi":  "Corporate",
		"Guest-WiFi": "Guest",
	}
	ds.SSID2SecurityDetails = map[string]dnac.ResponseItemWirelessGetEnterpriseSSIDSSIDDetails{
		"Corp-WiFi": {
			Name:          "Corp-WiFi",
			SecurityLevel: "wpa2_enterprise",
			IsEnabled:     &isEnabled,
		},
	}

	if ds.WirelessLANInterfaceName2VlanID["wlan-iface-1"] != 100 {
		t.Errorf("expected vlan 100 for wlan-iface-1, got %d", ds.WirelessLANInterfaceName2VlanID["wlan-iface-1"])
	}
	if ds.SSID2WlanGroupName["Corp-WiFi"] != "Corporate" {
		t.Errorf("expected Corporate group for Corp-WiFi, got %q", ds.SSID2WlanGroupName["Corp-WiFi"])
	}
	if ds.SSID2SecurityDetails["Corp-WiFi"].SecurityLevel != "wpa2_enterprise" {
		t.Errorf("expected wpa2_enterprise security for Corp-WiFi")
	}
}

func TestDnacSourceDeviceID2isMissingPrimaryIP(t *testing.T) {
	ds := newTestDnacSource()

	ds.DeviceID2isMissingPrimaryIP.Store("dev-1", true)
	ds.DeviceID2isMissingPrimaryIP.Store("dev-2", true)
	ds.DeviceID2isMissingPrimaryIP.Store("dev-3", false)

	tests := []struct {
		name     string
		deviceID string
		want     bool
	}{
		{"device with missing primary IP", "dev-1", true},
		{"another missing", "dev-2", true},
		{"device with primary IP set", "dev-3", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := ds.DeviceID2isMissingPrimaryIP.Load(tt.deviceID)
			if !ok {
				t.Fatalf("device %q not found", tt.deviceID)
			}
			got, ok := val.(bool)
			if !ok {
				t.Fatalf("type assertion failed for device %q", tt.deviceID)
			}
			if got != tt.want {
				t.Errorf("DeviceID2isMissingPrimaryIP[%q] = %v, want %v", tt.deviceID, got, tt.want)
			}
		})
	}
}

func TestDnacSourceConcurrentDeviceMapAccess(t *testing.T) {
	ds := newTestDnacSource()

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			device := &objects.Device{Name: "device"}
			ds.DeviceID2nbDevice.Store(id, device)
			ds.DeviceID2nbDevice.Load(id)
		}(i)
	}
	wg.Wait()

	count := 0
	ds.DeviceID2nbDevice.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	if count != 100 {
		t.Errorf("expected 100 devices stored, got %d", count)
	}
}

func TestDnacSourceSiteAdditionalInfo(t *testing.T) {
	site := dnac.ResponseSitesGetSiteResponse{
		ID:       "site-1",
		Name:     "Headquarters",
		ParentID: "",
		AdditionalInfo: []dnac.ResponseSitesGetSiteResponseAdditionalInfo{
			{
				Namespace: "Location",
				Attributes: dnac.ResponseSitesGetSiteResponseAdditionalInfoAttributes{
					Address:   "123 Main St, City",
					Latitude:  "46.0569",
					Longitude: "14.5058",
					Country:   "Slovenia",
				},
			},
			{
				Namespace: "System",
				Attributes: dnac.ResponseSitesGetSiteResponseAdditionalInfoAttributes{
					Type: "building",
				},
			},
		},
	}

	var locationInfo *dnac.ResponseSitesGetSiteResponseAdditionalInfo
	for i := range site.AdditionalInfo {
		if site.AdditionalInfo[i].Namespace == "Location" {
			locationInfo = &site.AdditionalInfo[i]
			break
		}
	}

	if locationInfo == nil {
		t.Fatal("expected Location namespace in AdditionalInfo")
	}
	if locationInfo.Attributes.Address != "123 Main St, City" {
		t.Errorf("address = %q, want '123 Main St, City'", locationInfo.Attributes.Address)
	}
	if locationInfo.Attributes.Latitude != "46.0569" {
		t.Errorf("latitude = %q, want '46.0569'", locationInfo.Attributes.Latitude)
	}
	if locationInfo.Attributes.Longitude != "14.5058" {
		t.Errorf("longitude = %q, want '14.5058'", locationInfo.Attributes.Longitude)
	}
}

func TestDnacSourceDeviceStatusMapping(t *testing.T) {
	tests := []struct {
		name               string
		reachabilityStatus string
		wantStatus         *objects.DeviceStatus
	}{
		{"reachable device", "Reachable", &objects.DeviceStatusActive},
		{"unreachable device", "Unreachable", &objects.DeviceStatusOffline},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deviceStatus := &objects.DeviceStatusActive
			if tt.reachabilityStatus == "Unreachable" {
				deviceStatus = &objects.DeviceStatusOffline
			}
			if deviceStatus != tt.wantStatus {
				t.Errorf("status = %v, want %v", deviceStatus, tt.wantStatus)
			}
		})
	}
}

func TestDnacSourceDeviceDescriptionTruncation(t *testing.T) {
	tests := []struct {
		name            string
		description     string
		wantDescription string
		wantComments    string
	}{
		{
			name:            "short description",
			description:     "A short description",
			wantDescription: "A short description",
			wantComments:    "",
		},
		{
			name:            "empty description",
			description:     "",
			wantDescription: "",
			wantComments:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var description, comments string
			if tt.description != "" {
				description = tt.description
			}
			if len(description) > objects.MaxDescriptionLength {
				comments = description
				description = "See comments"
			}
			if description != tt.wantDescription {
				t.Errorf("description = %q, want %q", description, tt.wantDescription)
			}
			if comments != tt.wantComments {
				t.Errorf("comments = %q, want %q", comments, tt.wantComments)
			}
		})
	}
}
