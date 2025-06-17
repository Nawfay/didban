package downloader

import (
	"fmt"
	"io"
	"os"

	"github/nawfay/didban/internal/logic"
	"github/nawfay/didban/internal/models"
	"github/nawfay/didban/internal/utils"

	"github.com/kkdai/youtube/v2"
)


type Client struct {

	Client youtube.Client
	
}

func NewClient() *Client {
	return &Client{
		Client: youtube.Client{},
	}
}

func DownloadVideoYt(client *Client, videoID string, tmpPath string, output string, track models.Track) (bool, error) {

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
	defer stream.Close() // Close the stream when done


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

	utils.ConvertToMP4(tmpFile, output, fmt.Sprintf("%d", video.Duration))
	os.Remove(tmpFile)

	coverPath := fmt.Sprintf("%s/%s.jpg", tmpPath, videoID)
	err = utils.FetchCover(track.Album.Cover, coverPath)
	if err != nil {
		os.Remove(output)
		os.Remove(coverPath)
		return false, fmt.Errorf("failed to fetch cover image: %w", err)
	}

	err = utils.TagMP3(output, coverPath, track.Title, track.Artist.Name, track.Album.Title, fmt.Sprintf("%d", video.Duration))
	if err != nil {
		os.Remove(output)
		os.Remove(coverPath)
		return false, fmt.Errorf("failed to tag MP3: %w", err)
	}
	os.Remove(coverPath)

	return true, nil
}




func ExampleClient() {
// 	videoID := "0k00mxCZtcQ"
	
// 	client := NewClient()


// 	err1, err2 := DownloadVideoYt(client, videoID, "/home/nawaf/Documents/GitHub/didban/tmp", "output.mp3")
// 	if err1 != false {
// 		fmt.Println("Error downloading video:", err1)
// 	} else {
// 		fmt.Println("Video downloaded successfully:", err2)
// 	}
// 	// Output: Video downloaded successfully: <nil>
// 	// fmt.Println("Video downloaded successfully:", err2)
// }

	fmt.Println("This is an example function for the Client struct.")

	client := NewClient()
	track, _ := logic.FetchTrack(3383978371)
	videoID, _ := logic.DeezerToYtResolver(3383978371)
	err1, err2 := DownloadVideoYt(client, videoID, "/home/nawaf/Documents/GitHub/didban/tmp", "output.mp3", *track)
	
	if err1 != false {
		fmt.Println("Error downloading video:", err1)
	} else {
		fmt.Println("Video downloaded successfully:", err2)
	}
	// Output: Video downloaded successfully: <nil>
	// fmt.Println("Video downloaded successfully:", err2)
}




