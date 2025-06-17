package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/bogem/id3v2/v2"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// ConvertToMP4 rewraps any input video as MP4 container
func ConvertToMP4(input, output string) error {
	return ffmpeg_go.
		Input(input).
		Output(output, ffmpeg_go.KwArgs{
			"c":        "copy",
			"movflags": "+faststart", // for web compatibility
		}).
		OverWriteOutput().
		Run()
}

func TagMP3(mp3Path, coverPath, title, artist, album, year string) error {
    // Open the MP3 file (read/write)
    tag, err := id3v2.Open(mp3Path, id3v2.Options{Parse: true})
    if err != nil {
        return fmt.Errorf("error opening mp3: %w", err)
    }
    defer tag.Close()

    // Set text metadata
    tag.SetTitle(title)
    tag.SetArtist(artist)
    tag.SetAlbum(album)
    tag.AddTextFrame(tag.CommonID("Year"), tag.DefaultEncoding(), year)

    // Load cover image from disk
    picBytes, err := ioutil.ReadFile(coverPath)
    if err != nil {
        return fmt.Errorf("error reading cover image: %w", err)
    }

    // Create and set picture frame
    pic := id3v2.PictureFrame{
        Encoding:    id3v2.EncodingISO, // ISO-8859-1 is standard for images
        MimeType:    "image/jpeg",       // or "image/png"
        PictureType: id3v2.PTFrontCover,
        Description: "Cover",
        Picture:     picBytes,
    }
    tag.AddAttachedPicture(pic)

    // Save changes
    if err = tag.Save(); err != nil {
        return fmt.Errorf("error saving tag: %w", err)
    }
    return nil
}