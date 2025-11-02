package state

import (
	"context"

	"github.com/amarnathcjd/gogram/telegram"
)

type (
	PlatformName string

	Platform interface {
		Name() PlatformName
		IsValid(query string) bool
		GetTracks(query string) ([]*Track, error)
		Download(ctx context.Context, track *Track, mystic *telegram.NewMessage) (string, error)
		IsDownloadSupported(source PlatformName) bool
	}

	Track struct {
		ID       string
		Title    string
		Duration int
		Artwork  string
		URL      string
		BY       string
		Source   PlatformName
	}
)

const (
	PlatformYouTube   PlatformName = "YouTube"
	PlatformTelegram  PlatformName = "Telegram"
	PlatformFallenApi PlatformName = "FallenApi"
	PlatformYtDlp     PlatformName = "YtDlp"
)

var DownloadCancels = make(map[int64]context.CancelFunc)
