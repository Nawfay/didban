package downloader

import (
	"fmt"
	"io"
	"os"

	"github.com/nawfay/didban/didban/logic"
	"github.com/nawfay/didban/didban/models"
	"github.com/nawfay/didban/didban/utils"

	"github.com/kkdai/youtube/v2"
)

func DownloadTrackYt(client *youtube.Client, videoID string, tmpPath string, path string, track *models.Track) (bool, error) {

	// Fetch video metadata
	video, err := client.GetVideo(videoID)
	if err != nil {
		return false, fmt.Errorf("failed to get video info: %w", err)
	}

	video.Formats.Type("mp4")
	// Find best audio-only stream (highest bitrate)
	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])

	if err != nil {
		return false, fmt.Errorf("failed to get audio stream: %w", err)
	}
	defer stream.Close() // Close the stream when done

	tmpFile := fmt.Sprintf("%s/%s.tmp_audio", tmpPath, videoID)

	out, err := os.Create(tmpFile)
	if err != nil {
		return false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, stream)
	if err != nil {
		return false, fmt.Errorf("failed to write audio stream: %w", err)
	}

	trackPath := fmt.Sprintf("%s/%s.mp3", path, utils.GenerateTrackTitle(track))

	utils.ConvertToMP4(tmpFile, trackPath, fmt.Sprintf("%d", video.Duration))
	os.Remove(tmpFile)

	err = utils.TagTrackWithMetadata(tmpPath, trackPath, videoID, track)
	if err != nil {
		os.Remove(trackPath)
		return false, fmt.Errorf("failed to tag MP3: %w", err)
	}

	return true, nil
}

func ExampleClient() {

	fmt.Println("This is an example function for the Client struct.")

	client := youtube.Client{}
	track, _ := logic.FetchTrack("77450636")
	videoID, _ := logic.DeezerToYtResolver("77450636")
	err1, err2 := DownloadTrackYt(&client, videoID, "/home/nawaf/Documents/GitHub/didban/tmp", fmt.Sprintf("%s - %s.mp3", track.Title, track.Artist.Name), track)

	if err1 != false {
		fmt.Println("Error downloading video:", err1)
	} else {
		fmt.Println("Video downloaded successfully:", err2)
	}
	// Output: Video downloaded successfully: <nil>
	// fmt.Println("Video downloaded successfully:", err2)
}
