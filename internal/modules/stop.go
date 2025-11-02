package modules

import (
	"github.com/amarnathcjd/gogram/telegram"
)

func stopHandler(m *telegram.NewMessage) error {
	return handleStop(m, false)
}

func cstopHandler(m *telegram.NewMessage) error {
	return handleStop(m, true)
}

func handleStop(m *telegram.NewMessage, cplay bool) error {
	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.EndGroup
	}
	if !r.IsActiveChat() {
		m.Reply("⚠️ <b>No active playback.</b>\nNothing is playing right now.")
		return telegram.EndGroup
	}
	r.Destroy()
	m.Reply("⏹️ <b>Playback stopped and cleared.</b>")
	return telegram.EndGroup
}
