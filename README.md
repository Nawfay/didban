# Didban – Music Fetching Library

**Didban** (Persian: دیدبان – *watchman*) is a Go library for downloading and tagging music tracks from Deezer and YouTube.

## ✨ Features

- 🎵 **Multi-source downloads**: Deezer (FLAC/MP3_320/MP3_128) and YouTube
- 🔐 **Deezer decryption**: Blowfish decryption for encrypted streams  
- 🏷️ **Complete metadata**: ID3 tags with artist, title, album, and cover art
- 🎯 **Auto-resolution**: YouTube search from Deezer metadata

## �* Installation

```bash
go get github.com/nawfay/didban
```

## 🚀 Quick Start

```go
package main

import (
    "log"
    "github.com/nawfay/didban"
    "github.com/nawfay/didban/didban/models"
)

func main() {
    // Initialize with Deezer ARL cookie
    err := didban.Init("your-deezer-arl-cookie", "/tmp")
    if err != nil {
        log.Fatal(err)
    }

    // Download from Deezer
    track := models.QueueItem{
        DeezerID: "3135556",
        Title:    "Harder Better Faster Stronger",
        Artist:   "Daft Punk",
        Path:     "./downloads",
        Youtube:  false, // true for YouTube
    }

    err = didban.DownloadTracks(track)
    if err != nil {
        log.Fatal(err)
    }
}
```

## 📖 Usage

### Deezer Download
```go
track := models.QueueItem{
    DeezerID: "1013324012",
    Path:     "./downloads",
    Youtube:  false,
}
```

### YouTube Download
```go
track := models.QueueItem{
    DeezerID: "1013324012", // Still needed for metadata
    Path:     "./downloads", 
    Youtube:  true,
}
```

### Batch Downloads
```go
tracks := []models.QueueItem{...}
for _, track := range tracks {
    err := didban.DownloadTracks(track)
    // handle error
}
```

## 🔧 API

### `didban.Init(arl, tmpPath string) error`
Initialize with Deezer ARL cookie and temp directory.

### `didban.DownloadTracks(models.QueueItem) error`
Download track based on queue item configuration.

### `models.QueueItem`
```go
type QueueItem struct {
    DeezerID string // Required: Deezer track ID
    Title    string // Track metadata
    Artist   string
    Album    string
    Path     string // Output directory
    Youtube  bool   // true=YouTube, false=Deezer
    Status   string // pending/downloading/completed/failed
}
```

## 🔐 Authentication

Get your Deezer ARL cookie:
1. Login to Deezer in browser
2. Open DevTools → Application → Cookies
3. Copy `arl` cookie value from `deezer.com`

```go
// Use environment variable
arl := os.Getenv("DEEZER_ARL")
```

## 🎵 How It Works

**Deezer**: ARL auth → encrypted stream → Blowfish decrypt → ID3 tag → MP3
**YouTube**: Deezer metadata → YouTube search → audio stream → convert → tag → MP3

## ⚠️ Legal Notice

Educational/personal use only. Comply with Deezer/YouTube ToS and copyright laws.

## 📄 License

MIT License