package platform

import (
	"fmt"
	"sync"
)

var (
	registry = make(map[string]Platform)
	mu       sync.RWMutex
)

// Register adds a platform to the global registry.
func Register(p Platform) {
	mu.Lock()
	defer mu.Unlock()
	registry[p.Name()] = p
}

// GetPlatform returns a registered platform by name.
func GetPlatform(name string) (Platform, bool) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := registry[name]
	return p, ok
}

// List returns all registered platform names.
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// All returns all registered platforms.
func All() []Platform {
	mu.RLock()
	defer mu.RUnlock()
	platforms := make([]Platform, 0, len(registry))
	for _, p := range registry {
		platforms = append(platforms, p)
	}
	return platforms
}

// MustGet returns a platform by name or panics.
func MustGet(name string) Platform {
	p, ok := GetPlatform(name)
	if !ok {
		panic(fmt.Sprintf("platform %q not registered", name))
	}
	return p
}