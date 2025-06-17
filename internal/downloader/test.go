// downloader/downloader.go
package downloader

import (
	"bytes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/blowfish"
)

// Config holds your ARL cookie (for testing), the fetched license token,
// and the Blowfish secret/IV used to decrypt Deezer’s stream.
var Config = struct {
	ARLCookie      string
	LicenseToken   string
	BlowfishSecret string
	BlowfishIV     string
}{
	ARLCookie:      "",
	LicenseToken:   "",
	BlowfishSecret: "g4el58wc0zvf9na1",
	BlowfishIV:     "0001020304050607",
}

// ─── Auth Bootstrap ────────────────────────────────────────────────────────────

// getUserDataResponse matches Deezer’s getUserData JSON envelope.
type getUserDataResponse struct {
	Results struct {
		User struct {
			Options struct {
				LicenseToken string `json:"license_token"`
			} `json:"options"`
		} `json:"user"`
	} `json:"results"`
}

// fetchLicenseToken posts your ARL to Deezer and returns a fresh license_token + sid.
func fetchLicenseToken(arl string) (licenseToken, sid string, err error) {
	const url = "https://www.deezer.com/ajax/gw-light.php?method=deezer.getUserData&input=3&api_version=1.0&api_token=null"
	req, err := http.NewRequest("POST", url, strings.NewReader("{}"))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "arl", Value: arl})

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var body getUserDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", "", fmt.Errorf("decode getUserData JSON: %w", err)
	}

	// grab the 'sid' cookie if you ever need it
	for _, c := range resp.Cookies() {
		if c.Name == "sid" {
			sid = c.Value
			break
		}
	}
	if sid == "" {
		return "", "", fmt.Errorf("sid cookie not found")
	}

	return body.Results.User.Options.LicenseToken, sid, nil
}

// SetARLCookie hard-codes your ARL and fetches+stores a valid license_token.
func SetARLCookie(arl string) error {
	lt, _, err := fetchLicenseToken(arl)
	if err != nil {
		return fmt.Errorf("unable to fetch license_token: %w", err)
	}
	Config.ARLCookie = arl
	Config.LicenseToken = lt
	return nil
}

// ─── Data Models ──────────────────────────────────────────────────────────────

// Song holds the minimal metadata we need from the public API.
type Song struct {
	ID         int    `json:"id"`
	TrackToken string `json:"track_token"`
	Title      string `json:"title"`
	Artist     struct {
		Name string `json:"name"`
	} `json:"artist"`
}

// MediaResponse models the v1/get_url response.
type MediaResponse struct {
	Data []struct {
		Media []struct {
			Sources []struct {
				URL      string `json:"url"`
				Provider string `json:"provider"`
			} `json:"sources"`
		} `json:"media"`
		Errors []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	} `json:"data"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

// ─── Public API ────────────────────────────────────────────────────────────────

// DownloadTrack retrieves and decrypts a full Deezer track (FLAC/320/128).
// trackID is the Deezer track ID; outputPath is the destination file.
func DownloadTrack(trackID, outputPath string) error {
	// 1) fetch public metadata + track_token
	song, err := getTrackMetadata(trackID)
	if err != nil {
		return fmt.Errorf("metadata error: %w", err)
	}

	// 2) fetch encrypted source URLs
	media, err := getMediaURL(song)
	if err != nil {
		return fmt.Errorf("media error: %w", err)
	}

	// 3) pick a URL (prefer provider "ak")
	if len(media.Data) == 0 ||
		len(media.Data[0].Media) == 0 ||
		len(media.Data[0].Media[0].Sources) == 0 {
		return fmt.Errorf("no media URL found")
	}
	url := media.Data[0].Media[0].Sources[0].URL
	for _, src := range media.Data[0].Media[0].Sources {
		if src.Provider == "ak" {
			url = src.URL
			break
		}
	}

	// 4) download & decrypt
	return downloadAndDecrypt(url, outputPath, song.ID)
}

// ─── Internal Helpers ─────────────────────────────────────────────────────────

// getTrackMetadata calls the public API and reads track_token.
func getTrackMetadata(trackID string) (*Song, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.deezer.com/track/%s", trackID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", "arl="+Config.ARLCookie)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d", resp.StatusCode)
	}

	var song Song
	if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}
	if song.TrackToken == "" {
		return nil, fmt.Errorf("invalid track token")
	}
	return &song, nil
}

// getMediaURL posts to /v1/get_url with your license_token + track_token.
func getMediaURL(song *Song) (*MediaResponse, error) {
	payload := fmt.Sprintf(
		`{"license_token":"%s","media":[{"type":"FULL","formats":[{"cipher":"BF_CBC_STRIPE","format":"FLAC"},{"cipher":"BF_CBC_STRIPE","format":"MP3_320"},{"cipher":"BF_CBC_STRIPE","format":"MP3_128"}]}],"track_tokens":["%s"]}`,
		Config.LicenseToken, song.TrackToken,
	)
	resp, err := http.Post("https://media.deezer.com/v1/get_url", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var media MediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&media); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}
	if len(media.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s (code %d)", media.Errors[0].Message, media.Errors[0].Code)
	}
	return &media, nil
}

const chunkSize = 2048

// downloadAndDecrypt streams in 2 048-byte chunks, decrypting every 3rd one.
func downloadAndDecrypt(url, outPath string, trackID int) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d", resp.StatusCode)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}
	defer f.Close()

	key := generateBlowfishKey(fmt.Sprintf("%d", trackID))
	buf := make([]byte, chunkSize)
	for chunk := 0; ; chunk++ {
		// read up to chunkSize
		n, rerr := io.ReadFull(resp.Body, buf)
		if rerr != nil && rerr != io.ErrUnexpectedEOF && rerr != io.EOF {
			return rerr
		}
		if n == 0 {
			break
		}
		data := buf[:n]

		// decrypt every 3rd full chunk
		if chunk%3 == 0 && n == chunkSize {
			dec, derr := DecryptBlowfish(data, key)
			if derr != nil {
				return derr
			}
			data = dec
		}

		if _, werr := f.Write(data); werr != nil {
			return werr
		}
		if rerr == io.ErrUnexpectedEOF || rerr == io.EOF {
			break
		}
	}

	return nil
}

// DecryptBlowfish decrypts a full chunk using CBC with your IV.
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
	dst := make([]byte, len(data))
	mode.CryptBlocks(dst, data)
	return dst, nil
}

// generateBlowfishKey derives the 16-byte key via md5(songID) ⊕ secret.
func generateBlowfishKey(songID string) []byte {
	sum := md5.Sum([]byte(songID))
	hexStr := fmt.Sprintf("%x", sum[:])
	secret := []byte(Config.BlowfishSecret)
	key := make([]byte, 16)
	for i := 0; i < 16; i++ {
		key[i] = secret[i] ^ hexStr[i] ^ hexStr[i+16]
	}
	return key
}
