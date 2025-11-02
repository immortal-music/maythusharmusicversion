package modules

import (
	"fmt"
	"html"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func positionHandler(m *telegram.NewMessage) error {
	return handlePosition(m, false)
}

func cpositionHandler(m *telegram.NewMessage) error {
	return handlePosition(m, true)
}

func handlePosition(m *telegram.NewMessage, cplay bool) error {
	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.EndGroup
	}

	if !r.IsActiveChat() || r.Track == nil {
		m.Reply("⚠️ <b>No active playback.</b>\nNothing is playing right now.")
		return telegram.EndGroup
	}

	r.Parse()

	progress := fmt.Sprintf(
		"🎵 <b>%s</b>\n📍 <code>%s / %s</code>\n⚙️ Speed: <b>%.2fx</b>",
		html.EscapeString(utils.ShortTitle(r.Track.Title, 25)),
		formatDuration(r.Position),
		formatDuration(r.Track.Duration),
		r.Speed,
	)

	m.Reply(progress)
	return telegram.EndGroup
}
