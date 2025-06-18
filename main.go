// main.go
package didban

import (
	"github.com/nawfay/didban/internal/downloader"
)

type Didban struct {
	Client *downloader.Client
	ARL    string
}

func NewClient(arl string) *Didban {
	return &Didban{
		Client: downloader.NewYtClient(),
		ARL:    arl,
	}
}
