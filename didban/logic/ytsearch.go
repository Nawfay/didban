package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// extractVideoIDs is a recursive helper to parse YouTube's initial data blob
func extractVideoIDs(data interface{}, ids *[]string) {
	switch val := data.(type) {
	case map[string]interface{}:
		if renderer, ok := val["videoRenderer"]; ok {
			if videoMap, ok := renderer.(map[string]interface{}); ok {
				if videoId, ok := videoMap["videoId"].(string); ok {
					*ids = append(*ids, videoId)
				}
			}
		}
		for _, v := range val {
			extractVideoIDs(v, ids)
		}
	case []interface{}:
		for _, v := range val {
			extractVideoIDs(v, ids)
		}
	}
}

// SearchYouTube returns a list of YouTube video IDs for a given query
func SearchYouTube(query string, maxResults int) ([]string, error) {
	url := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", strings.ReplaceAll(query, " ", "+"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/114.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	body := string(bodyBytes)

	start := strings.Index(body, "var ytInitialData =")
	if start == -1 {
		return nil, fmt.Errorf("ytInitialData not found")
	}
	start += len("var ytInitialData = ")
	end := strings.Index(body[start:], ";</script>")
	if end == -1 {
		return nil, fmt.Errorf("end of ytInitialData not found")
	}
	jsonStr := body[start : start+end]

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("JSON parse failed: %w", err)
	}

	var videoIDs []string
	extractVideoIDs(parsed, &videoIDs)

	if len(videoIDs) > maxResults {
		videoIDs = videoIDs[:maxResults]
	}

	return videoIDs, nil
}
