package utils

import (
	"fmt"

	"github.com/nawfay/didban/internal/models"

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

func FetchCover(url string, outputPath string) error {
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

func TagTrackWithMetadata(tmpPath string, trackPath string, id string, track *models.Track) error {

	coverPath := fmt.Sprintf("%s/%s.jpg", tmpPath, id)
	err := FetchCover(track.Album.Cover, coverPath)
	if err != nil {
		os.Remove(coverPath)
		return fmt.Errorf("failed to fetch cover image: %w", err)
	}

	err = TagMP3(trackPath, coverPath, track.Title, track.Artist.Name, track.Album.Title, fmt.Sprintf("%d", track.Duration))
	if err != nil {
		os.Remove(coverPath)
		return fmt.Errorf("failed to tag MP3: %w", err)
	}
	os.Remove(coverPath)

	return nil
}

func GenerateTrackTitle(track *models.Track) string {
	return fmt.Sprintf("%s - %s", track.Artist.Name, track.Title)
}
