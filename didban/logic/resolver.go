package logic

import (
	"fmt"

	"github.com/nawfay/didban/didban/utils"
)

func DeezerToYtResolver(deezerId string) (string, error) {

	song, err := FetchTrack(deezerId)
	if err != nil {
		return "", err
	}

	ytquery := fmt.Sprintf(song.Title + " " + song.Artist.Name + " lyrics")
	ytquery = utils.NormalizeString(ytquery)

	ytString, err := SearchYouTube(ytquery, 1)
	if err != nil {
		return "", err
	}

	return ytString[0], nil
}
