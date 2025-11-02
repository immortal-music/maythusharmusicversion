package modules

import (
	"fmt"
	"html"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func unmuteHandler(m *telegram.NewMessage) error {
	return handleUnmute(m, false)
}

func cunmuteHandler(m *telegram.NewMessage) error {
	return handleUnmute(m, true)
}

func handleUnmute(m *telegram.NewMessage, cplay bool) error {
	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.EndGroup
	}
	if !r.IsActiveChat() {
		m.Reply("⚠️ <b>No active playback.</b>\nThere’s nothing playing right now.")
		return telegram.EndGroup
	}
	if !r.IsMuted() {
		m.Reply("ℹ️ <b>Already Unmuted</b>\nThe music is not muted in this chat.")
		return telegram.EndGroup
	}
	mention := utils.MentionHTML(m.Sender)
	trackTitle := html.EscapeString(utils.ShortTitle(r.Track.Title, 25))
	if _, err := r.Unmute(); err != nil {
		m.Reply(fmt.Sprintf("❌ <b>Playback Unmute Failed</b>\nError: <code>%v</code>", err))
		return telegram.EndGroup
	}
	msg := fmt.Sprintf(
		"🔊 <b>Unmuted playback</b>\n\n🎵 Track: %s\n👤 Unmuted by: %s",
		trackTitle, mention,
	)
	if sp := r.GetSpeed(); sp != 1.0 {
		msg += fmt.Sprintf("\n⚙️ Speed: <b>%.2fx</b>", sp)
	}
	m.Reply(msg)
	return telegram.EndGroup
}
