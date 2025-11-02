package platforms

import (
	"sort"
	"sync"

	"github.com/immortal-music/maythusharmusicversion/internal/state"
)

type regEntry struct {
	platform state.Platform
	priority int
}

var (
	registry = make(map[state.PlatformName]regEntry)
	regLock  sync.RWMutex
)

func addPlatform(priority int, name state.PlatformName, p state.Platform) {
	regLock.Lock()
	defer regLock.Unlock()
	registry[name] = regEntry{
		platform: p,
		priority: priority,
	}
}

func getPlatform(name state.PlatformName) (state.Platform, bool) {
	entry, ok := registry[name]
	if !ok {
		return nil, false
	}
	return entry.platform, true
}

func getOrderedPlatforms() []state.Platform {
	platforms := make([]regEntry, 0, len(registry))
	for _, entry := range registry {
		platforms = append(platforms, entry)
	}

	sort.Slice(platforms, func(i, j int) bool {
		return platforms[i].priority > platforms[j].priority
	})

	result := make([]state.Platform, len(platforms))
	for i, entry := range platforms {
		result[i] = entry.platform
	}
	return result
}
