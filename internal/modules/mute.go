package modules

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func muteHandler(m *telegram.NewMessage) error {
	return handleMute(m, false)
}

func cmuteHandler(m *telegram.NewMessage) error {
	return handleMute(m, true)
}

func handleMute(m *telegram.NewMessage, cplay bool) error {
	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.EndGroup
	}
	if !r.IsActiveChat() {
		m.Reply("⚠️ <b>No active playback.</b>\nThere’s nothing playing right now.")
		return telegram.EndGroup
	}
	if r.IsPaused() {
		m.Reply("⚠️ <b>Oops!</b>\nThe room is paused. Please resume it first to mute playback.")
		return telegram.EndGroup
	}
	if r.IsMuted() {
		remaining := r.RemainingUnmuteDuration()
		if remaining > 0 {
			m.Reply(fmt.Sprintf("🔇 <b>Already Muted</b>\n\nThe music is already muted in this chat.\nAuto-unmute in <b>%s</b>.", formatDuration(int(remaining.Seconds()))))
		} else {
			m.Reply("🔇 <b>Already Muted</b>\nThe music is already muted in this chat.\nWould you like to unmute it?")
		}
		return telegram.EndGroup
	}
	mention := utils.MentionHTML(m.Sender)
	args := strings.Fields(m.Text())
	var autoUnmuteDuration time.Duration
	if len(args) >= 2 {
		rawDuration := strings.ToLower(strings.TrimSpace(args[1]))
		rawDuration = strings.TrimSuffix(rawDuration, "s")
		if seconds, err := strconv.Atoi(rawDuration); err == nil {
			if seconds < 5 || seconds > 3600 {
				m.Reply("⚠️ Invalid duration for auto-unmute. It must be between <b>5</b> and <b>3600</b> seconds.")
				return telegram.EndGroup
			}
			autoUnmuteDuration = time.Duration(seconds) * time.Second
		} else {
			m.Reply(fmt.Sprintf("⚠️ Invalid format. Use: <code>/%s 30</code> or <code>/%s 30s</code>", getCommand(m), getCommand(m)))
			return telegram.EndGroup
		}
	}
	var muteErr error
	if autoUnmuteDuration > 0 {
		_, muteErr = r.Mute(autoUnmuteDuration)
	} else {
		_, muteErr = r.Mute()
	}
	if muteErr != nil {
		m.Reply(fmt.Sprintf("❌ <b>Playback Mute Failed</b>\nError: <code>%v</code>", muteErr))
		return telegram.EndGroup
	}
	msg := fmt.Sprintf(
		"🔇 <b>Muted playback</b>\n\n🎵 Track: %s\n👤 Muted by: %s\n",
		html.EscapeString(utils.ShortTitle(r.Track.Title, 25)),
		mention,
	)
	if sp := r.GetSpeed(); sp != 1.0 {
		msg += fmt.Sprintf("⚙️ Speed: <b>%.2fx</b>\n", sp)
	}
	if autoUnmuteDuration > 0 {
		msg += fmt.Sprintf("\n<i>⏳ Auto-unmute in <b>%d</b> seconds</i>", int(autoUnmuteDuration.Seconds()))
	}
	m.Reply(msg)
	return telegram.EndGroup
}
