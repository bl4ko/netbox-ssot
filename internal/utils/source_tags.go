package utils

import (
	"strings"
	"sync"
)

// SourceTagPrefix is the prefix of the per-source tags. Other tools (e.g. netbox-sync)
// use the same convention, so the prefix alone does not tell who owns a tag.
const SourceTagPrefix = "Source: "

// managedSourceTags holds the names of the source tags configured for this run.
var managedSourceTags = struct {
	sync.RWMutex
	names map[string]bool
}{}

// SetManagedSourceTagNames registers the names of the source tags that netbox-ssot
// owns in this run. It must be called with the configured sources before syncing.
func SetManagedSourceTagNames(names []string) {
	set := make(map[string]bool, len(names))
	for _, name := range names {
		set[name] = true
	}
	managedSourceTags.Lock()
	defer managedSourceTags.Unlock()
	managedSourceTags.names = set
}

// IsManagedSourceTag reports whether tagName is a source tag owned by netbox-ssot.
func IsManagedSourceTag(tagName string) bool {
	managedSourceTags.RLock()
	defer managedSourceTags.RUnlock()
	return managedSourceTags.names[tagName]
}

// IsForeignSourceTag reports whether tagName is a source tag owned by another tool.
func IsForeignSourceTag(tagName string) bool {
	return strings.HasPrefix(tagName, SourceTagPrefix) && !IsManagedSourceTag(tagName)
}
