package theme

import (
	"sort"
	"strings"
)

// registry maps lowercase theme names to constructor functions.
var registry = map[string]func() *Theme{ //nolint:gochecknoglobals // theme registry is intentionally global
	"modern":  Modern,
	"default": MermaidDefault,
	"dark":    Dark,
	"forest":  Forest,
	"neutral": Neutral,
}

// ByName returns a theme by its name (case-insensitive).
// Returns the Modern theme if the name is not recognized.
func ByName(name string) *Theme {
	if fn, ok := registry[strings.ToLower(name)]; ok {
		return fn()
	}
	return Modern()
}

// Names returns a sorted list of available theme names.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
