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

func TestOrphanManagerAddItemTracksNetboxSsotObjects(t *testing.T) {
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
			// Tracked like any netbox-ssot object; Deletable leaves it out.
			name: "object also tagged by another tool",
			item: &objects.Platform{
				NetboxObject: objects.NetboxObject{
					ID:   11,
					Tags: []*objects.Tag{ssotTag, {ID: 3, Name: "Source: vcenter.example"}, {ID: 4, Name: "NetBox-synced"}},
				},
				Name: "Debian 13",
			},
			managed: true,
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

// Init loads every object into the orphan manager before the managed source tags are
// registered: ownership must be decided at deletion time, not when the item is added.
func TestOrphanOwnershipIsDecidedAtDeletionTime(t *testing.T) {
	utils.SetManagedSourceTagNames(nil)
	t.Cleanup(func() { utils.SetManagedSourceTagNames(nil) })
	ssotTag := &objects.Tag{ID: 1, Name: constants.SsotTagName}
	own := &objects.Platform{
		NetboxObject: objects.NetboxObject{ID: 10, Tags: []*objects.Tag{ssotTag, {ID: 2, Name: "Source: proxmox-a"}}},
		Name:         "Debian 12",
	}
	shared := &objects.Platform{
		NetboxObject: objects.NetboxObject{ID: 11, Tags: []*objects.Tag{ssotTag, {ID: 3, Name: "Source: vcenter.example"}}},
		Name:         "Debian 13",
	}
	orphanManager := NewOrphanManager(mockLogger)
	orphanManager.AddItem(own)
	orphanManager.AddItem(shared)

	utils.SetManagedSourceTagNames([]string{"Source: proxmox-a"})

	deletable := orphanManager.Deletable(constants.PlatformsAPIPath)
	if len(deletable) != 1 || deletable[0].GetID() != 10 {
		t.Errorf("Deletable() = %v, want only platform 10 (platform 11 is also tagged by another tool)", deletable)
	}
}

func TestRegisterManagedSourceTagsKeepsCustomNamedTagsOfRemovedSources(t *testing.T) {
	t.Cleanup(func() { utils.SetManagedSourceTagNames(nil) })
	nbi := &NetboxInventory{
		tagsIndexByName: map[string]*objects.Tag{
			// source.tag accepts any name: the description is the only ownership mark.
			"pve-old": {
				ID: 6, Name: "pve-old",
				Description: "Automatically created tag by netbox-ssot for source proxmox-old",
			},
			"prod": {ID: 7, Name: "prod", Description: "Production workloads"},
		},
	}

	nbi.RegisterManagedSourceTags(nil)

	if !utils.IsManagedSourceTag("pve-old") {
		t.Errorf("IsManagedSourceTag(%q) = false, want true: netbox-ssot created it", "pve-old")
	}
	if utils.IsManagedSourceTag("prod") {
		t.Errorf("IsManagedSourceTag(%q) = true, want false: a user tag", "prod")
	}
}
