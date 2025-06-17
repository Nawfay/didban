// main.go
package main

import (
	"log"

	"github/nawfay/didban/internal/downloader"
)

func main() {
	// 1) Hard-code your ARL here for testing:
	const myARL = "7ed297081270075398f7eaa1745457959452149c482f7ae4842be8f49dc7c3f127080b5da34ec0e2717f7bd5742d4211e4c71149ff3c0f8b7a8cdcc0854460b553b1d5b9bf841207771a0f4942ab7ed511c322fc6520a1b800da88fe67d3ff61"

	// 2) Bootstrap auth (fetches a fresh license_token):
	if err := downloader.SetARLCookie(myARL); err != nil {
		log.Fatalf("Auth setup failed: %v", err)
	}

	// 3) Download a track (FLAC/320/128 all supported):
	trackID := "2864040362"
	output  := "mytrack.mp3"

	if err := downloader.DownloadTrack(trackID, output); err != nil {
		log.Fatalf("Download failed: %v", err)
	}
	log.Println("Download complete:", output)
}
