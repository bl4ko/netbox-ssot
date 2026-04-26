package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

var MockDefaultSsotTag = &objects.Tag{
	ID:   0,
	Name: constants.SsotTagName,
}

// Hardcoded mock API responses for tags endpoint.
var (
	MockTagsGetResponse = Response[objects.Tag]{
		Count:    2, //nolint:mnd
		Next:     nil,
		Previous: nil,
		Results: []objects.Tag{
			{
				ID:          0,
				Name:        "Source: proxmox",
				Slug:        "source-proxmox",
				Color:       "9e9e9e",
				Description: "Automatically created tag by netbox-ssot for source proxmox",
			},
			{
				ID:          1,
				Name:        "netbox-ssot",
				Slug:        "netbox-ssot",
				Color:       "00add8",
				Description: "Tag used by netbox-ssot to mark devices that are managed by it",
			},
		},
	}
	MockTagPatchResponse = objects.Tag{
		ID:          1,
		Name:        "netbox-ssot",
		Slug:        "netbox-ssot",
		Color:       "00add8",
		Description: "patched description",
	}
	MockTagCreateResponse = objects.Tag{
		ID:          1,
		Name:        "netbox-ssot",
		Slug:        "netbox-ssot",
		Color:       "00add8",
		Description: "created description",
	}
)

// Hardcoded mock api return values for tenant endpoint.
var (
	MockTenantsGetResponse = Response[objects.Tenant]{
		Count:    2, //nolint:mnd
		Next:     nil,
		Previous: nil,
		Results: []objects.Tenant{
			{
				NetboxObject: objects.NetboxObject{
					ID: 1,
					Tags: []*objects.Tag{
						MockDefaultSsotTag,
					},
				},
				Name: "MockTenant1",
				Slug: "mock-tenant-1",
			},
			{
				NetboxObject: objects.NetboxObject{
					ID: 2, //nolint:mnd
					Tags: []*objects.Tag{
						MockDefaultSsotTag,
					},
				},
				Name: "MockTenant2",
				Slug: "mock-tenant-2",
			},
		},
	}
	MockTenantCreateResponse = objects.Tenant{
		NetboxObject: objects.NetboxObject{
			ID: 3, //nolint:mnd
		},
		Name: "MockTenant3",
		Slug: "mock-tenant-3",
	}
	MockTenantPatchResponse = objects.Tenant{
		NetboxObject: objects.NetboxObject{
			ID: 1,
		},
		Name: "MockPatched",
		Slug: "mock-patched-tenant",
	}
)

// Hardcoded mock api return values for site endpoint.
var (
	MockSitesGetResponse = Response[objects.Site]{
		Count:    2, //nolint:mnd
		Next:     nil,
		Previous: nil,
		Results: []objects.Site{
			{
				NetboxObject: objects.NetboxObject{
					ID: 1,
					Tags: []*objects.Tag{
						MockDefaultSsotTag,
					},
				},
				Name: "MockSite1",
				Slug: "mock-site-1",
			},
			{
				NetboxObject: objects.NetboxObject{
					ID: 2, //nolint:mnd
					Tags: []*objects.Tag{
						MockDefaultSsotTag,
					},
				},
				Name: "MockSite2",
				Slug: "mock-site-2",
			},
		},
	}
	MockSiteCreateResponse = objects.Site{
		NetboxObject: objects.NetboxObject{
			ID: 3, //nolint:mnd
		},
		Name: "MockSite3",
		Slug: "mock-site-3",
	}
	MockSitePatchResponse = objects.Site{
		NetboxObject: objects.NetboxObject{
			ID: 1,
		},
		Name: "MockSitePatched",
		Slug: "mock-site-patched",
	}
)

// Hardcoded mock api return values for VlanGroup endpoint.
var (
	MockVlanGroupsGetResponse = Response[objects.VlanGroup]{
		Count:    1,
		Next:     nil,
		Previous: nil,
		Results: []objects.VlanGroup{
			{
				NetboxObject: objects.NetboxObject{
					ID:   1,
					Tags: []*objects.Tag{MockDefaultSsotTag},
				},
				Name: "MockVlanGroup1",
				Slug: "mock-vlan-group-1",
			},
		},
	}
	MockVlanGroupCreateResponse = objects.VlanGroup{
		NetboxObject: objects.NetboxObject{
			ID: 1,
		},
		Name: "MockVlanGroup1",
		Slug: "mock-vlan-group-1",
	}
	MockVlanGroupPatchResponse = objects.VlanGroup{
		NetboxObject: objects.NetboxObject{
			ID: 1,
		},
		Name: "MockVlanGroupPatched",
		Slug: "mock-vlan-group-patched",
	}
)

// Hardcoded mock api return values for DeviceRole endpoint.
var (
	MockDeviceRolesGetResponse = Response[objects.DeviceRole]{
		Count:    1,
		Next:     nil,
		Previous: nil,
		Results: []objects.DeviceRole{
			{
				NetboxObject: objects.NetboxObject{
					ID:   1,
					Tags: []*objects.Tag{MockDefaultSsotTag},
				},
				Name:  "MockDeviceRole1",
				Slug:  "mock-device-role-1",
				Color: constants.Color(constants.DeviceRoleServerColor),
			},
		},
	}
	MockDeviceRoleCreateResponse = objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			ID: 2, //nolint:mnd
		},
		Name:  "MockDeviceRole2",
		Slug:  "mock-device-role-2",
		Color: constants.Color(constants.DeviceRoleServerColor),
	}
	MockDeviceRolePatchResponse = objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			ID: 1,
		},
		Name:  "MockDeviceRolePatched",
		Slug:  "mock-device-role-patched",
		Color: constants.Color(constants.DeviceRoleServerColor),
	}
)

// Hardcoded mock api return values for prefix endpoint.
// MockPrefixGetResponse simulates NetBox returning prefixes with object-type
// custom fields as nested objects (the read format from the NetBox REST API).
var (
	MockPrefixGetResponse = Response[objects.Prefix]{
		Count:    1,
		Next:     nil,
		Previous: nil,
		Results: []objects.Prefix{
			{
				NetboxObject: objects.NetboxObject{
					ID:   1,
					Tags: []*objects.Tag{MockDefaultSsotTag},
					// NetBox returns object-type custom fields as nested objects.
					// This is the read format — write format expects just the ID.
					CustomFields: map[string]interface{}{
						"source":           "test",
						"orphan_last_seen": nil,
						"site_ref": map[string]interface{}{
							"id":      float64(1),
							"display": "LCL",
							"url":     "https://netbox/api/dcim/sites/1/",
							"name":    "LCL",
							"slug":    "lcl",
						},
					},
				},
				Prefix: "10.0.0.0/24",
			},
		},
	}
	MockPrefixCreateResponse = objects.Prefix{
		NetboxObject: objects.NetboxObject{
			ID: 2, //nolint:mnd
		},
		Prefix: "192.168.1.0/24",
	}
	MockPrefixPatchResponse = objects.Prefix{
		NetboxObject: objects.NetboxObject{
			ID: 1,
			Tags: []*objects.Tag{
				MockDefaultSsotTag,
			},
			CustomFields: map[string]interface{}{
				"source":           "test",
				"orphan_last_seen": nil,
				"site_ref":         float64(1),
			},
		},
		Prefix: "10.0.0.0/24",
	}
)

// Mock responses for ContactRole endpoint.
var (
	MockContactRolesGetResponse = Response[objects.ContactRole]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.ContactRole{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockContactRole1",
				Slug:         "mock-contact-role-1",
			},
		},
	}
	MockContactRoleCreateResponse = objects.ContactRole{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockContactRole3",
		Slug:         "mock-contact-role-3",
	}
	MockContactRolePatchResponse = objects.ContactRole{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockContactRolePatched",
		Slug:         "mock-contact-role-patched",
	}
)

// Mock responses for ContactGroup endpoint.
var (
	MockContactGroupsGetResponse = Response[objects.ContactGroup]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.ContactGroup{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockContactGroup1",
				Slug:         "mock-contact-group-1",
			},
		},
	}
	MockContactGroupCreateResponse = objects.ContactGroup{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockContactGroup3",
		Slug:         "mock-contact-group-3",
	}
	MockContactGroupPatchResponse = objects.ContactGroup{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockContactGroupPatched",
		Slug:         "mock-contact-group-patched",
	}
)

// Mock responses for Contact endpoint.
var (
	MockContactsGetResponse = Response[objects.Contact]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.Contact{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockContact1",
			},
		},
	}
	MockContactCreateResponse = objects.Contact{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockContact3",
	}
	MockContactPatchResponse = objects.Contact{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockContactPatched",
	}
)

// Mock responses for ContactAssignment endpoint.
var (
	MockContactAssignmentsGetResponse = Response[objects.ContactAssignment]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.ContactAssignment{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				ModelType:    constants.ContentTypeDcimDevice,
				ObjectID:     1,
			},
		},
	}
	MockContactAssignmentCreateResponse = objects.ContactAssignment{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		ModelType:    constants.ContentTypeDcimDevice,
		ObjectID:     2, //nolint:mnd
	}
	MockContactAssignmentPatchResponse = objects.ContactAssignment{
		NetboxObject: objects.NetboxObject{ID: 1},
		ModelType:    constants.ContentTypeDcimDevice,
		ObjectID:     1,
	}
)

// Mock responses for CustomField endpoint.
var (
	MockCustomFieldsGetResponse = Response[objects.CustomField]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.CustomField{
			{
				ID:   1,
				Name: "mock_custom_field",
				Type: objects.CustomFieldTypeText,
			},
		},
	}
	MockCustomFieldCreateResponse = objects.CustomField{
		ID:   3, //nolint:mnd
		Name: "mock_custom_field_new",
		Type: objects.CustomFieldTypeText,
	}
	MockCustomFieldPatchResponse = objects.CustomField{
		ID:   1,
		Name: "mock_custom_field_patched",
		Type: objects.CustomFieldTypeText,
	}
)

// Mock responses for Location endpoint.
var (
	MockLocationsGetResponse = Response[objects.Location]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.Location{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockLocation1",
				Slug:         "mock-location-1",
			},
		},
	}
	MockLocationCreateResponse = objects.Location{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockLocation3",
		Slug:         "mock-location-3",
	}
	MockLocationPatchResponse = objects.Location{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockLocationPatched",
		Slug:         "mock-location-patched",
	}
)

// Mock responses for SiteGroup endpoint.
var (
	MockSiteGroupsGetResponse = Response[objects.SiteGroup]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.SiteGroup{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockSiteGroup1",
				Slug:         "mock-site-group-1",
			},
		},
	}
	MockSiteGroupCreateResponse = objects.SiteGroup{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockSiteGroup3",
		Slug:         "mock-site-group-3",
	}
	MockSiteGroupPatchResponse = objects.SiteGroup{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockSiteGroupPatched",
		Slug:         "mock-site-group-patched",
	}
)

// Mock responses for Manufacturer endpoint.
var (
	MockManufacturersGetResponse = Response[objects.Manufacturer]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.Manufacturer{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockManufacturer1",
				Slug:         "mock-manufacturer-1",
			},
		},
	}
	MockManufacturerCreateResponse = objects.Manufacturer{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockManufacturer3",
		Slug:         "mock-manufacturer-3",
	}
	MockManufacturerPatchResponse = objects.Manufacturer{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockManufacturerPatched",
		Slug:         "mock-manufacturer-patched",
	}
)

// Mock responses for DeviceType endpoint.
var (
	MockDeviceTypesGetResponse = Response[objects.DeviceType]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.DeviceType{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Model:        "MockDeviceType1",
				Slug:         "mock-device-type-1",
			},
		},
	}
	MockDeviceTypeCreateResponse = objects.DeviceType{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Model:        "MockDeviceType3",
		Slug:         "mock-device-type-3",
	}
	MockDeviceTypePatchResponse = objects.DeviceType{
		NetboxObject: objects.NetboxObject{ID: 1},
		Model:        "MockDeviceTypePatched",
		Slug:         "mock-device-type-patched",
	}
)

// Mock responses for Platform endpoint.
var (
	MockPlatformsGetResponse = Response[objects.Platform]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.Platform{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockPlatform1",
				Slug:         "mock-platform-1",
			},
		},
	}
	MockPlatformCreateResponse = objects.Platform{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockPlatform3",
		Slug:         "mock-platform-3",
	}
	MockPlatformPatchResponse = objects.Platform{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockPlatformPatched",
		Slug:         "mock-platform-patched",
	}
)

// Mock responses for Device endpoint.
var (
	MockDevicesGetResponse = Response[objects.Device]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.Device{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockDevice1",
			},
		},
	}
	MockDeviceCreateResponse = objects.Device{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockDevice3",
	}
	MockDevicePatchResponse = objects.Device{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockDevicePatched",
	}
)

// Mock responses for Interface endpoint.
var (
	MockInterfacesGetResponse = Response[objects.Interface]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.Interface{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockInterface1",
			},
		},
	}
	MockInterfaceCreateResponse = objects.Interface{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockInterface3",
	}
	MockInterfacePatchResponse = objects.Interface{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockInterfacePatched",
	}
)

// Mock responses for VirtualDeviceContext endpoint.
var (
	MockVirtualDeviceContextsGetResponse = Response[objects.VirtualDeviceContext]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.VirtualDeviceContext{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockVDC1",
			},
		},
	}
	MockVirtualDeviceContextCreateResponse = objects.VirtualDeviceContext{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockVDC3",
	}
	MockVirtualDeviceContextPatchResponse = objects.VirtualDeviceContext{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockVDCPatched",
	}
)

// Mock responses for MACAddress endpoint.
var (
	MockMACAddressesGetResponse = Response[objects.MACAddress]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.MACAddress{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				MAC:          "00:11:22:33:44:55",
			},
		},
	}
	MockMACAddressCreateResponse = objects.MACAddress{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		MAC:          "00:11:22:33:44:66",
	}
	MockMACAddressPatchResponse = objects.MACAddress{
		NetboxObject: objects.NetboxObject{ID: 1},
		MAC:          "00:11:22:33:44:77",
	}
)

// Mock responses for Vlan endpoint.
var (
	MockVlansGetResponse = Response[objects.Vlan]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.Vlan{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockVlan1",
				Vid:          100, //nolint:mnd
			},
		},
	}
	MockVlanCreateResponse = objects.Vlan{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockVlan3",
		Vid:          300, //nolint:mnd
	}
	MockVlanPatchResponse = objects.Vlan{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockVlanPatched",
		Vid:          100, //nolint:mnd
	}
)

// Mock responses for IPAddress endpoint.
var (
	MockIPAddressesGetResponse = Response[objects.IPAddress]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.IPAddress{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Address:      "10.0.0.1/24",
			},
		},
	}
	MockIPAddressCreateResponse = objects.IPAddress{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Address:      "10.0.0.3/24",
	}
	MockIPAddressPatchResponse = objects.IPAddress{
		NetboxObject: objects.NetboxObject{ID: 1},
		Address:      "10.0.0.1/24",
	}
)

// Mock responses for ClusterType endpoint.
var (
	MockClusterTypesGetResponse = Response[objects.ClusterType]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.ClusterType{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockClusterType1",
				Slug:         "mock-cluster-type-1",
			},
		},
	}
	MockClusterTypeCreateResponse = objects.ClusterType{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockClusterType3",
		Slug:         "mock-cluster-type-3",
	}
	MockClusterTypePatchResponse = objects.ClusterType{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockClusterTypePatched",
		Slug:         "mock-cluster-type-patched",
	}
)

// Mock responses for ClusterGroup endpoint.
var (
	MockClusterGroupsGetResponse = Response[objects.ClusterGroup]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.ClusterGroup{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockClusterGroup1",
				Slug:         "mock-cluster-group-1",
			},
		},
	}
	MockClusterGroupCreateResponse = objects.ClusterGroup{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockClusterGroup3",
		Slug:         "mock-cluster-group-3",
	}
	MockClusterGroupPatchResponse = objects.ClusterGroup{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockClusterGroupPatched",
		Slug:         "mock-cluster-group-patched",
	}
)

// Mock responses for Cluster endpoint.
var (
	MockClustersGetResponse = Response[objects.Cluster]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.Cluster{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockCluster1",
			},
		},
	}
	MockClusterCreateResponse = objects.Cluster{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockCluster3",
	}
	MockClusterPatchResponse = objects.Cluster{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockClusterPatched",
	}
)

// Mock responses for VM endpoint.
var (
	MockVMsGetResponse = Response[objects.VM]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.VM{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockVM1",
			},
		},
	}
	MockVMCreateResponse = objects.VM{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockVM3",
	}
	MockVMPatchResponse = objects.VM{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockVMPatched",
	}
)

// Mock responses for VMInterface endpoint.
var (
	MockVMInterfacesGetResponse = Response[objects.VMInterface]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.VMInterface{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockVMInterface1",
			},
		},
	}
	MockVMInterfaceCreateResponse = objects.VMInterface{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockVMInterface3",
	}
	MockVMInterfacePatchResponse = objects.VMInterface{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockVMInterfacePatched",
	}
)

// Mock responses for VirtualDisk endpoint.
var (
	MockVirtualDisksGetResponse = Response[objects.VirtualDisk]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.VirtualDisk{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockVirtualDisk1",
				Size:         100, //nolint:mnd
			},
		},
	}
	MockVirtualDiskCreateResponse = objects.VirtualDisk{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockVirtualDisk3",
		Size:         300, //nolint:mnd
	}
	MockVirtualDiskPatchResponse = objects.VirtualDisk{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockVirtualDiskPatched",
		Size:         100, //nolint:mnd
	}
)

// Mock responses for WirelessLAN endpoint.
var (
	MockWirelessLANsGetResponse = Response[objects.WirelessLAN]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.WirelessLAN{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				SSID:         "MockWirelessLAN1",
			},
		},
	}
	MockWirelessLANCreateResponse = objects.WirelessLAN{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		SSID:         "MockWirelessLAN3",
	}
	MockWirelessLANPatchResponse = objects.WirelessLAN{
		NetboxObject: objects.NetboxObject{ID: 1},
		SSID:         "MockWirelessLANPatched",
	}
)

// Mock responses for WirelessLANGroup endpoint.
var (
	MockWirelessLANGroupsGetResponse = Response[objects.WirelessLANGroup]{
		Count: 1, Next: nil, Previous: nil,
		Results: []objects.WirelessLANGroup{
			{
				NetboxObject: objects.NetboxObject{ID: 1, Tags: []*objects.Tag{MockDefaultSsotTag}},
				Name:         "MockWirelessLANGroup1",
				Slug:         "mock-wireless-lan-group-1",
			},
		},
	}
	MockWirelessLANGroupCreateResponse = objects.WirelessLANGroup{
		NetboxObject: objects.NetboxObject{ID: 3}, //nolint:mnd
		Name:         "MockWirelessLANGroup3",
		Slug:         "mock-wireless-lan-group-3",
	}
	MockWirelessLANGroupPatchResponse = objects.WirelessLANGroup{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "MockWirelessLANGroupPatched",
		Slug:         "mock-wireless-lan-group-patched",
	}
)

const (
	MockVersionResponseJSON = "{\"django-version\": \"4.2.10\"}"
)

// mockEndpointHandler creates a generic HTTP handler for a mock NetBox API endpoint.
// It supports GET (returns getResp), POST (echoes request body with injected createID, 201),
// PATCH (validates custom_fields, returns patchResp), and DELETE (204 no content).
// POST echoes the request body back with the create response's ID injected, so callers
// see their own field values (matching the real NetBox API behavior).
func mockEndpointHandler(getResp interface{}, createID int, patchResp interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			resp, err := json.Marshal(getResp)
			if err != nil {
				log.Printf("Error marshaling GET response: %v", err)
			}
			_, _ = w.Write(resp)
		case http.MethodPost:
			body, _ := io.ReadAll(r.Body)
			if err := validateCustomFieldsPayload(body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errResp := map[string][]string{"custom_fields": {err.Error()}}
				resp, _ := json.Marshal(errResp)
				_, _ = w.Write(resp)
				return
			}
			// Echo request body with injected ID (mimics real NetBox behavior).
			var obj map[string]interface{}
			if err := json.Unmarshal(body, &obj); err == nil {
				obj["id"] = createID
			}
			w.WriteHeader(http.StatusCreated)
			resp, err := json.Marshal(obj)
			if err != nil {
				log.Printf("Error marshaling POST response: %v", err)
			}
			_, _ = w.Write(resp)
		case http.MethodPatch:
			body, _ := io.ReadAll(r.Body)
			if err := validateCustomFieldsPayload(body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errResp := map[string][]string{"custom_fields": {err.Error()}}
				resp, _ := json.Marshal(errResp)
				_, _ = w.Write(resp)
				return
			}
			w.WriteHeader(http.StatusOK)
			resp, err := json.Marshal(patchResp)
			if err != nil {
				log.Printf("Error marshaling PATCH response: %v", err)
			}
			_, _ = w.Write(resp)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			log.Printf("Wrong http method: %q", r.Method) //nolint:gosec
		}
	}
}

// mockEndpoint groups mock response data for a single API endpoint.
type mockEndpoint struct {
	path     constants.APIPath
	get      interface{}
	createID int
	patch    interface{}
}

func CreateMockServer() *httptest.Server {
	handler := http.NewServeMux()

	endpoints := []mockEndpoint{
		// Extras
		{constants.TagsAPIPath, MockTagsGetResponse, 1, MockTagPatchResponse},
		{constants.CustomFieldsAPIPath, MockCustomFieldsGetResponse, 3, MockCustomFieldPatchResponse},
		// Tenancy
		{constants.TenantsAPIPath, MockTenantsGetResponse, 3, MockTenantPatchResponse},
		{constants.ContactRolesAPIPath, MockContactRolesGetResponse, 3, MockContactRolePatchResponse},
		{constants.ContactGroupsAPIPath, MockContactGroupsGetResponse, 3, MockContactGroupPatchResponse},
		{constants.ContactsAPIPath, MockContactsGetResponse, 3, MockContactPatchResponse},
		{constants.ContactAssignmentsAPIPath, MockContactAssignmentsGetResponse, 3, MockContactAssignmentPatchResponse},
		// DCIM
		{constants.SitesAPIPath, MockSitesGetResponse, 3, MockSitePatchResponse},
		{constants.SiteGroupsAPIPath, MockSiteGroupsGetResponse, 3, MockSiteGroupPatchResponse},
		{constants.LocationsAPIPath, MockLocationsGetResponse, 3, MockLocationPatchResponse},
		{constants.ManufacturersAPIPath, MockManufacturersGetResponse, 3, MockManufacturerPatchResponse},
		{constants.DeviceTypesAPIPath, MockDeviceTypesGetResponse, 3, MockDeviceTypePatchResponse},
		{constants.DeviceRolesAPIPath, MockDeviceRolesGetResponse, 2, MockDeviceRolePatchResponse},
		{constants.PlatformsAPIPath, MockPlatformsGetResponse, 3, MockPlatformPatchResponse},
		{constants.DevicesAPIPath, MockDevicesGetResponse, 3, MockDevicePatchResponse},
		{constants.InterfacesAPIPath, MockInterfacesGetResponse, 3, MockInterfacePatchResponse},
		{
			constants.VirtualDeviceContextsAPIPath,
			MockVirtualDeviceContextsGetResponse, 3, MockVirtualDeviceContextPatchResponse,
		},
		{constants.MACAddressesAPIPath, MockMACAddressesGetResponse, 3, MockMACAddressPatchResponse},
		// IPAM
		{constants.VlanGroupsAPIPath, MockVlanGroupsGetResponse, 1, MockVlanGroupPatchResponse},
		{constants.VlansAPIPath, MockVlansGetResponse, 3, MockVlanPatchResponse},
		{constants.IPAddressesAPIPath, MockIPAddressesGetResponse, 3, MockIPAddressPatchResponse},
		{constants.PrefixesAPIPath, MockPrefixGetResponse, 2, MockPrefixPatchResponse},
		// Virtualization
		{constants.ClusterTypesAPIPath, MockClusterTypesGetResponse, 3, MockClusterTypePatchResponse},
		{constants.ClusterGroupsAPIPath, MockClusterGroupsGetResponse, 3, MockClusterGroupPatchResponse},
		{constants.ClustersAPIPath, MockClustersGetResponse, 3, MockClusterPatchResponse},
		{constants.VirtualMachinesAPIPath, MockVMsGetResponse, 3, MockVMPatchResponse},
		{constants.VMInterfacesAPIPath, MockVMInterfacesGetResponse, 3, MockVMInterfacePatchResponse},
		{constants.VirtualDisksAPIPath, MockVirtualDisksGetResponse, 3, MockVirtualDiskPatchResponse},
		// Wireless
		{constants.WirelessLANsAPIPath, MockWirelessLANsGetResponse, 3, MockWirelessLANPatchResponse},
		{constants.WirelessLANGroupsAPIPath, MockWirelessLANGroupsGetResponse, 3, MockWirelessLANGroupPatchResponse},
	}

	for _, ep := range endpoints {
		handler.HandleFunc(string(ep.path), mockEndpointHandler(ep.get, ep.createID, ep.patch))
	}

	// Special handlers
	handler.HandleFunc("/api/status/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, MockVersionResponseJSON)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	handler.HandleFunc("/api/read-error", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:all
		w.(http.Flusher).Flush()
	})

	// Wildcard handler for all other paths
	handler.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error": "page not found"}`)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	return httptest.NewServer(handler)
}

var MockNetboxClient = &NetboxClient{
	HTTPClient: &http.Client{},
	Logger:     &logger.Logger{Logger: log.Default()},
	BaseURL:    "",
	APIToken:   "testtoken",
	Timeout:    constants.DefaultAPITimeout,
}

var FailingMockNetboxClient = &NetboxClient{
	HTTPClient: &http.Client{Transport: &FailingHTTPClient{}},
	Logger:     &logger.Logger{Logger: log.Default()},
	BaseURL:    "",
	APIToken:   "testtoken",
	Timeout:    constants.DefaultAPITimeout,
}

type FailingHTTPClient struct{}

func (m *FailingHTTPClient) RoundTrip(_ *http.Request) (*http.Response, error) {
	// Return an error to simulate a failure in the HTTP request
	return nil, fmt.Errorf("mock error")
}

var MockNetboxClientWithReadError = &NetboxClient{
	HTTPClient: &http.Client{Transport: &FailingHTTPClientRead{}},
	Logger:     &logger.Logger{Logger: log.Default()},
	BaseURL:    "",
	APIToken:   "testtoken",
	Timeout:    constants.DefaultAPITimeout,
}

type FailingHTTPClientRead struct{}

func (m *FailingHTTPClientRead) RoundTrip(_ *http.Request) (*http.Response, error) {
	// Simulate a response with a FaultyReader as its Body
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(&FaultyReader{}),
		Header:     make(http.Header),
	}, nil
}

// validateCustomFieldsPayload mimics NetBox 4.2.x REST API validation:
// object-type custom field values must be sent as scalar IDs, not as nested
// objects. If a custom_fields value is a map containing "display", NetBox
// returns 400: "Cannot resolve keyword 'display' into field".
func validateCustomFieldsPayload(body []byte) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil //nolint:nilerr
	}
	cf, ok := payload["custom_fields"]
	if !ok {
		return nil
	}
	cfMap, ok := cf.(map[string]interface{})
	if !ok {
		return nil
	}
	for key, val := range cfMap {
		if nested, isMap := val.(map[string]interface{}); isMap {
			if _, hasDisplay := nested["display"]; hasDisplay {
				return fmt.Errorf(
					"cannot resolve keyword 'display' into field for custom_field '%s'", key, //nolint:perfsprint
				)
			}
		}
	}
	return nil
}

type FaultyReader struct{}

func (m *FaultyReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("mock read error")
}
