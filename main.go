package didban

import (
	"fmt"

	"github.com/nawfay/didban/didban/downloader"
	"github.com/nawfay/didban/didban/logic"
	"github.com/nawfay/didban/didban/models"

	"github.com/kkdai/youtube/v2"
)

type State struct {
	YtClient *youtube.Client
	ARL      string
	TmpPath  string
}

var appState *State

func Init(arl string, tmpPath string) error {
	// Set global state
	appState = &State{
		YtClient: &youtube.Client{},
		ARL:      arl,
		TmpPath:  tmpPath,
	}

	// Set up ARL session (deezer cookie for downloading)
	if err := downloader.SetARLCookie(arl) 
	err != nil {
		return err
	}

	return nil
}


func DownloadTracks(item models.QueueItem) error {
	if item.Youtube && item.DeezerID != "" {
		track, err := logic.FetchTrack(item.DeezerID)
		if err != nil {
			return err
		}

		youtubeId, err := logic.DeezerToYtResolver(item.DeezerID)
		if err != nil {
			return err
		}
		
		finished, err := downloader.DownloadTrackYt(appState.YtClient, youtubeId, appState.TmpPath, item.Path, track)
		if err != nil {
			return err
		}
		if !finished {
			return fmt.Errorf("track %s not found on YouTube", item.Title)
		}

		return nil
	}

	track, err := logic.FetchTrack(item.DeezerID)
	if err != nil {
		return err
	}

	finished, err := downloader.DownloadTrackDeezer(track, appState.TmpPath, item.Path)
	if err != nil {
		return err
	}
	if !finished {
		return fmt.Errorf("track %s not found on Deezer", item.Title)
	}

	return nil
}
