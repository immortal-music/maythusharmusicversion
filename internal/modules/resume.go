package modules

import (
	"fmt"
	"html"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func resumeHandler(m *telegram.NewMessage) error {
	return handleResume(m, false)
}

func cresumeHandler(m *telegram.NewMessage) error {
	return handleResume(m, true)
}

func handleResume(m *telegram.NewMessage, cplay bool) error {
	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.EndGroup
	}
	if !r.IsActiveChat() {
		m.Reply("⚠️ <b>No active playback.</b>\nNothing is playing right now.")
		return telegram.EndGroup
	}
	if !r.IsPaused() {
		m.Reply("ℹ️ <b>Already Playing</b>\nThe music is already playing in this chat.\nWould you like to pause it?")
		return telegram.EndGroup
	}
	if _, err := r.Resume(); err != nil {
		m.Reply(fmt.Sprintf("❌ <b>Playback Resume Failed</b>\nError: <code>%v</code>", err))
	} else {
		title := html.EscapeString(utils.ShortTitle(r.Track.Title, 25))
		pos := formatDuration(r.Position)
		total := formatDuration(r.Track.Duration)
		mention := utils.MentionHTML(m.Sender)
		msg := fmt.Sprintf("▶️ Resuming playback:\n\n <b>Title: </b>\"%s\"\n📍 Position: %s / %s\nResumed by: %s", title, pos, total, mention)

		if sp := r.GetSpeed(); sp != 1.0 {
			msg += fmt.Sprintf("\n⚙️ Speed: <b>%.2fx</b>", sp)
		}
		m.Reply(msg)
	}
	return telegram.EndGroup
}
