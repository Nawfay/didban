package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nawfay/didban/internal/models"
)

const DeezerApi = "https://api.deezer.com/"

func FetchTrack(id int) (*models.Track, error) {
	// Construct the API URL
	url := DeezerApi + "track/" + fmt.Sprint(id)
	var trackResponse models.Track

	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 HTTP responses
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch track: " + resp.Status)
	}

	// Decode the response body into the TrackResponse struct
	err = json.NewDecoder(resp.Body).Decode(&trackResponse)
	if err != nil {
		return nil, err
	}

	// Return the tracks from the response
	return &trackResponse, nil
}
