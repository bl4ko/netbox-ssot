package inventory

import (
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// RegisterManagedSourceTags registers the source tags owned by netbox-ssot: the ones
// configured for this run, plus the ones netbox-ssot created for sources that have since
// been removed or renamed, so that their objects are still cleaned up as orphans.
// It must be called after the inventory has loaded the tags from NetBox.
func (nbi *NetboxInventory) RegisterManagedSourceTags(configuredTagNames []string) {
	names := append([]string{}, configuredTagNames...)
	nbi.tagsLock.Lock()
	for name, tag := range nbi.tagsIndexByName {
		if strings.HasPrefix(name, utils.SourceTagPrefix) &&
			strings.HasPrefix(tag.Description, constants.SourceTagDescriptionPrefix) {
			names = append(names, name)
		}
	}
	nbi.tagsLock.Unlock()
	utils.SetManagedSourceTagNames(names)
}
