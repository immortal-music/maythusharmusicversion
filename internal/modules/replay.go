package modules

import (
	"fmt"
	"html"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func replayHandler(m *telegram.NewMessage) error {
	return handleReplay(m, false)
}

func creplayHandler(m *telegram.NewMessage) error {
	return handleReplay(m, true)
}

func handleReplay(m *telegram.NewMessage, cplay bool) error {
	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.EndGroup
	}

	if !r.IsActiveChat() {
		m.Reply("⚠️ <b>No active playback.</b>\nNothing is playing right now.")
		return telegram.EndGroup
	}

	if err := r.Replay(); err != nil {
		m.Reply(fmt.Sprintf("❌ <b>Replay Failed</b>\nError: <code>%v</code>", err))
	} else {
		trackTitle := html.EscapeString(utils.ShortTitle(r.Track.Title, 25))
		totalDuration := formatDuration(r.Track.Duration)
		m.Reply(fmt.Sprintf("🔁 Now replaying:\n\n<b>Title: </b>%s\n🎵 Duration: <code>%s</code>\n⏱️ Speed: %.2fx", trackTitle, totalDuration, r.Speed))
	}

	return telegram.EndGroup
}
