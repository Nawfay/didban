package utils

import (
	"net/http"
	"os"
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


func FetchCover(url string, outputPath string) (error) {
	// fmt.Println("Downloading image:", url, outputPath)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.Body)
	return err
}