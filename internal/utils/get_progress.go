package utils

import (
	"fmt"
	"math"

	"github.com/amarnathcjd/gogram/telegram"
)

func GetProgress(mystic *telegram.NewMessage) *telegram.ProgressManager {
	pm := telegram.NewProgressManager(2)

	if mystic == nil {
		return pm
	}

	var opts *telegram.SendOptions
	if replyMarkup := mystic.ReplyMarkup(); replyMarkup != nil {
		opts = &telegram.SendOptions{ReplyMarkup: *replyMarkup}
	}

	pm.WithEdit(func(total, current int64) {
		if mystic == nil {
			return
		}

		text := fmt.Sprintf(
			"📥 Downloading your track...\n\nProgress: %.1f%%\nETA: %s | Speed: %s",
			pm.GetProgress(current),
			pm.GetETA(current),
			pm.GetSpeed(current),
		)

		if opts != nil {
			mystic.Edit(text, *opts)
		} else {
			mystic.Edit(text)
		}
	})

	return pm
}

func GetProgressBar(playedSec, durationSec int) string {
	if durationSec == 0 || playedSec <= 0 {
		return "◉—————————"
	}

	percentage := (float64(playedSec) / float64(durationSec)) * 100
	umm := math.Floor(percentage)

	var bar string

	switch {
	case umm > 0 && umm <= 10:
		bar = "◉—————————"
	case umm > 10 && umm < 20:
		bar = "—◉————————"
	case umm >= 20 && umm < 30:
		bar = "——◉———————"
	case umm >= 30 && umm < 40:
		bar = "———◉——————"
	case umm >= 40 && umm < 50:
		bar = "————◉—————"
	case umm >= 50 && umm < 60:
		bar = "—————◉————"
	case umm >= 60 && umm < 70:
		bar = "——————◉———"
	case umm >= 70 && umm < 80:
		bar = "———————◉——"
	case umm >= 80 && umm < 90:
		bar = "————————◉—"
	case umm >= 90 && umm <= 100:
		bar = "—————————◉"
	default:
		bar = "—————————◉"
	}

	return bar
}
