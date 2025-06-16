package utils


import (
	"regexp"
	"strings"
)

// NormalizeName cleans up a name to be search/query friendly (e.g. for Google or YouTube)
func NormalizeString(name string) string {
	name = strings.ToLower(name)

	// Allow only alphanumerics, spaces, dots, underscores, and hyphens
	clean := regexp.MustCompile(`[^a-z0-9 ._-]+`)
	name = clean.ReplaceAllString(name, "")

	// Collapse multiple spaces
	space := regexp.MustCompile(`\s+`)
	name = space.ReplaceAllString(name, " ")

	name = strings.TrimSpace(name)

	// Replace spaces with +
	name = strings.ReplaceAll(name, " ", "+")

	return name
}
