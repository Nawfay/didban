package downloader

import (
	"bytes"
	"crypto/cipher"
	"crypto/md5"
	// "encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/blowfish"
)

var Config = struct {
	ARLCookie      string
	LicenseToken   string
	BlowfishSecret string
	BlowfishIV     string
}{
	ARLCookie:      "7ed297081270075398f7eaa1745457959452149c482f7ae4842be8f49dc7c3f127080b5da34ec0e2717f7bd5742d4211e4c71149ff3c0f8b7a8cdcc0854460b553b1d5b9bf841207771a0f4942ab7ed511c322fc6520a1b800da88fe67d3ff61",
	LicenseToken:   "AAAAAmfW8uBn6rlgFo2SPQeagYjuPVCJto7VYRu51OK45xKOAzOjJZYbDG7nFTHL5wH9_Vljve40tlcG7fR8wY80vEI2dapH4TYvkAHo-qirmXkwtbNIGFA8xzc",
	BlowfishSecret: "g4el58wc0zvf9na1",
	BlowfishIV:     "0001020304050607",
}

type Song struct {
	ID         int    json:"id"
	TrackToken string json:"track_token"
	Title      string json:"title"
	Artist     struct {
		Name string json:"name"
	} json:"artist"
}

type MediaResponse struct {
	Data []struct {
		Media []struct {
			Sources []struct {
				URL      string json:"url"
				Provider string json:"provider" // Added field for provider
			} json:"sources"
		} json:"media"
		Errors []struct {
			Code    int    json:"code"
			Message string json:"message"
		} json:"errors"
	} json:"data"
	Errors []struct {
		Code    int    json:"code"
		Message string json:"message"
	} json:"errors"
}

func DownloadTrack(trackID string, outputPath string) error {
	song, err := getTrackMetadata(trackID)
	if err != nil {
		return fmt.Errorf("metadata error: %w", err)
	}

	media, err := getMediaURL(song)
	if err != nil {
		return fmt.Errorf("media error: %w", err)
	}

	fmt.Println(media)

	if len(media.Data) == 0 || len(media.Data[0].Media) == 0 || len(media.Data[0].Media[0].Sources) == 0 {
		return fmt.Errorf("no media URL found")
	}

	// Choose the URL, prioritizing the source with Provider "ak"
	url := media.Data[0].Media[0].Sources[0].URL
	for _, source := range media.Data[0].Media[0].Sources {
		if source.Provider == "ak" {
			url = source.URL
			break
		}
	}

	fmt.Println(url)

	return downloadAndDecrypt(
		url,
		outputPath,
		song.ID,
	)
}


func getTrackMetadata(trackID string) (*Song, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.deezer.com/track/%s", trackID), nil)
	req.Header.Add("Cookie", "arl="+Config.ARLCookie)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d status", resp.StatusCode)
	}

	var song Song
	if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}

	if song.TrackToken == "" {
		return nil, fmt.Errorf("invalid track token in response")
	}

	return &song, nil
}

func getMediaURL(song *Song) (*MediaResponse, error) {
	payload := fmt.Sprintf({
		"license_token": "%s",
		"media": [{"type": "FULL", "formats": [{"cipher":"BF_CBC_STRIPE","format":"FLAC"},{"cipher":"BF_CBC_STRIPE","format":"MP3_320"},{"cipher":"BF_CBC_STRIPE","format":"MP3_128"}]}],
		"track_tokens": ["%s"]
	}, Config.LicenseToken, song.TrackToken)

	reqBody := []byte(payload)
	resp, err := http.Post("https://media.deezer.com/v1/get_url", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var media MediaResponse
	if err := json.Unmarshal(body, &media); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}

	if len(media.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s (code %d)", media.Errors[0].Message, media.Errors[0].Code)
	}

	return &media, nil
}

const ChunkSize = 2048

func downloadAndDecrypt(url, path string, trackID int) error {
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("download failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP error %d", resp.StatusCode)
    }

    file, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("file creation failed: %w", err)
    }
    defer file.Close()

    // Define the chunk size (must be divisible by blowfish.BlockSize, e.g., 8)
    // const ChunkSize = 2048

	trackIDStr := fmt.Sprintf("%d", trackID)
	// key := generateBlowfishKey(trackIDStr)


    key := generateBlowfishKey(trackIDStr)
    buffer := make([]byte, ChunkSize)

	for chunk := 0; ; chunk++ {
		totalRead := 0
		for totalRead < ChunkSize {
			n, err := resp.Body.Read(buffer[totalRead:])
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				return err
			}

			if n > 0 {
				totalRead += n
			}
		}

		if totalRead == 0 {
			break
		}

		if chunk%3 == 0 && totalRead == ChunkSize {
			buffer, err = DecryptBlowfish(buffer, key)
			if err != nil {
				return err
			}
		}

		_, err = file.Write(buffer[:totalRead])
		if err != nil {
			return err
		}

		if totalRead < ChunkSize {
			break
		}
	}

	return nil
}

func DecryptBlowfish(data, key []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(Config.BlowfishIV)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	return decrypted, nil
}
func generateBlowfishKey(songID string) []byte {
	hash := md5.Sum([]byte(songID))
	hashHex := fmt.Sprintf("%x", hash)

	key := []byte(Config.BlowfishSecret)
	for i := 0; i < len(hash); i++ {
		key[i] = key[i] ^ hashHex[i] ^ hashHex[i+16]
	}

	fmt.Println(key)

	return key
}


