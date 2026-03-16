package mapper

import (
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func TestType2PathAndPath2TypeConsistency(t *testing.T) {
	t.Run("every Type2Path entry has a reverse in Path2Type", func(t *testing.T) {
		for typ, path := range Type2Path {
			reverseType, ok := Path2Type[path]
			if !ok {
				t.Errorf("Path2Type missing entry for path %q (type %v)", path, typ)
				continue
			}
			if reverseType != typ {
				t.Errorf("Path2Type[%q] = %v, want %v", path, reverseType, typ)
			}
		}
	})

	t.Run("every Path2Type entry has a reverse in Type2Path", func(t *testing.T) {
		for path, typ := range Path2Type {
			reversePath, ok := Type2Path[typ]
			if !ok {
				t.Errorf("Type2Path missing entry for type %v (path %q)", typ, path)
				continue
			}
			if reversePath != path {
				t.Errorf("Type2Path[%v] = %q, want %q", typ, reversePath, path)
			}
		}
	})

	t.Run("maps have the same length", func(t *testing.T) {
		if len(Type2Path) != len(Path2Type) {
			t.Errorf("Type2Path has %d entries, Path2Type has %d entries", len(Type2Path), len(Path2Type))
		}
	})
}

func TestType2PathCompleteness(t *testing.T) {
	expectedPaths := []constants.APIPath{
		constants.VlanGroupsAPIPath,
		constants.VlansAPIPath,
		constants.IPAddressesAPIPath,
		constants.ClusterTypesAPIPath,
		constants.ClusterGroupsAPIPath,
		constants.ClustersAPIPath,
		constants.VirtualMachinesAPIPath,
		constants.VMInterfacesAPIPath,
		constants.DevicesAPIPath,
		constants.MACAddressesAPIPath,
		constants.VirtualDeviceContextsAPIPath,
		constants.DeviceRolesAPIPath,
		constants.DeviceTypesAPIPath,
		constants.InterfacesAPIPath,
		constants.SitesAPIPath,
		constants.SiteGroupsAPIPath,
		constants.ManufacturersAPIPath,
		constants.PlatformsAPIPath,
		constants.TenantsAPIPath,
		constants.ContactGroupsAPIPath,
		constants.ContactRolesAPIPath,
		constants.ContactsAPIPath,
		constants.CustomFieldsAPIPath,
		constants.TagsAPIPath,
		constants.ContactAssignmentsAPIPath,
		constants.PrefixesAPIPath,
		constants.WirelessLANsAPIPath,
		constants.WirelessLANGroupsAPIPath,
		constants.VirtualDisksAPIPath,
		constants.VRFsAPIPath,
	}

	for _, path := range expectedPaths {
		if _, ok := Path2Type[path]; !ok {
			t.Errorf("Path2Type missing expected API path %q", path)
		}
	}
}

func TestReverseMapDoesNotLoseEntries(t *testing.T) {
	// Ensure no duplicate values in Type2Path that would cause reverseMap to drop entries.
	seen := make(map[constants.APIPath]reflect.Type)
	for typ, path := range Type2Path {
		if prev, ok := seen[path]; ok {
			t.Errorf("duplicate API path %q: used by both %v and %v", path, prev, typ)
		}
		seen[path] = typ
	}
}
