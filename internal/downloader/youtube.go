package downloader

import (
	"fmt"
	"io"
	"os"

	"github.com/kkdai/youtube/v2"
	"github/nawfay/didban/internal/utils"
)


type Client struct {

	Client youtube.Client
	
}

func NewClient() *Client {
	return &Client{
		Client: youtube.Client{},
	}
}

func DownloadVideoYt(client *Client, videoID string, tmpPath string, output string) (bool, error) {

	// Fetch video metadata
	video, err := client.Client.GetVideo(videoID)
	if err != nil {
		return false, fmt.Errorf("failed to get video info: %w", err)
	}

	video.Formats.Type("mp4")
	// Find best audio-only stream (highest bitrate)
	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.Client.GetStream(video, &formats[0])
	
	if err != nil {
		return false,fmt.Errorf("failed to get audio stream: %w", err)
	}
	defer stream.Close() // ✅ Close the stream when done


	tmpFile := fmt.Sprintf("%s/%s.tmp_audio",tmpPath, videoID)

	out, err := os.Create(tmpFile)
	if err != nil {
		return false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, stream)
	if err != nil {
		return false, fmt.Errorf("failed to write audio stream: %w", err)
	}	

	utils.ConvertToMP4(tmpFile, output)
	os.Remove(tmpFile)

	return true, nil
}




func ExampleClient() {
	videoID := "0k00mxCZtcQ"
	
	client := NewClient()
	DownloadVideoYt(client, videoID, "/Users/nawaf/Documents/GitHub/Sorn/didban/tmp", "output.mp3")
}
