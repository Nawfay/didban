// main.go
package main

import (
	"github.com/nawfay/didban/internal/downloader"
)

type didban struct {
	Client *downloader.Client
	ARL    string
}

func NewClient(arl string) *didban {
	return &didban{
		Client: downloader.NewYtClient(),
		ARL:    arl,
	}
}
