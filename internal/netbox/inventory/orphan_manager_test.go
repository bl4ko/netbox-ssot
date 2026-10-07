package inventory

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

func TestRegisterManagedSourceTagsKeepsTagsOfRemovedSources(t *testing.T) {
	t.Cleanup(func() { utils.SetManagedSourceTagNames(nil) })
	nbi := &NetboxInventory{
		tagsIndexByName: map[string]*objects.Tag{
			"Source: proxmox-a": {
				ID: 2, Name: "Source: proxmox-a",
				Description: "Automatically created tag by netbox-ssot for source proxmox-a",
			},
			"Source: proxmox-old": {
				ID: 5, Name: "Source: proxmox-old",
				Description: "Automatically created tag by netbox-ssot for source proxmox-old",
			},
			"Source: vcenter.example": {
				ID: 3, Name: "Source: vcenter.example",
				Description: "Marks objects synced by netbox-sync",
			},
		},
	}

	nbi.RegisterManagedSourceTags([]string{"Source: proxmox-a"})

	for _, name := range []string{"Source: proxmox-a", "Source: proxmox-old"} {
		if !utils.IsManagedSourceTag(name) {
			t.Errorf("IsManagedSourceTag(%q) = false, want true: netbox-ssot created it", name)
		}
	}
	if !utils.IsForeignSourceTag("Source: vcenter.example") {
		t.Errorf("IsForeignSourceTag(%q) = false, want true: another tool created it", "Source: vcenter.example")
	}
}

func TestOrphanManagerAddItemSkipsObjectsOwnedByAnotherTool(t *testing.T) {
	utils.SetManagedSourceTagNames([]string{"Source: proxmox-a"})
	t.Cleanup(func() { utils.SetManagedSourceTagNames(nil) })

	ssotTag := &objects.Tag{ID: 1, Name: constants.SsotTagName}
	tests := []struct {
		name    string
		item    *objects.Platform
		managed bool
	}{
		{
			name: "object tagged by netbox-ssot only",
			item: &objects.Platform{
				NetboxObject: objects.NetboxObject{ID: 10, Tags: []*objects.Tag{ssotTag, {ID: 2, Name: "Source: proxmox-a"}}},
				Name:         "Debian 12",
			},
			managed: true,
		},
		{
			name: "object also tagged by another tool",
			item: &objects.Platform{
				NetboxObject: objects.NetboxObject{
					ID:   11,
					Tags: []*objects.Tag{ssotTag, {ID: 3, Name: "Source: vcenter.example"}, {ID: 4, Name: "NetBox-synced"}},
				},
				Name: "Debian 13",
			},
			managed: false,
		},
		{
			name: "object without netbox-ssot tag",
			item: &objects.Platform{
				NetboxObject: objects.NetboxObject{ID: 12, Tags: []*objects.Tag{{ID: 2, Name: "Source: proxmox-a"}}},
				Name:         "Debian 11",
			},
			managed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orphanManager := NewOrphanManager(mockLogger)
			orphanManager.AddItem(tt.item)
			_, got := orphanManager.Items[constants.PlatformsAPIPath][tt.item.ID]
			if got != tt.managed {
				t.Errorf("AddItem(%s) managed = %t, want %t", tt.item.Name, got, tt.managed)
			}
		})
	}
}
