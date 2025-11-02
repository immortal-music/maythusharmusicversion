package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func loopHandler(m *telegram.NewMessage) error {
	return handleLoop(m, false)
}

func cloopHandler(m *telegram.NewMessage) error {
	return handleLoop(m, true)
}

func handleLoop(m *telegram.NewMessage, cplay bool) error {
	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.EndGroup
	}
	args := strings.Fields(m.Text())
	currentLoop := r.Loop
	if !r.IsActiveChat() {
		m.Reply("⚠️ <b>No active playback.</b>\nThere's nothing playing right now.")
		return telegram.EndGroup
	}
	if len(args) < 2 {
		msg := fmt.Sprintf("🔁 <b>Loop Control</b>\n\nUsage: %s [count]\n• 0 - Disable loop\n• 1-10 - Loop count", getCommand(m))
		if currentLoop > 0 {
			msg += fmt.Sprintf("\n• Current loop: <b>%d</b> time(s)", currentLoop)
		}
		m.Reply(msg)
		return telegram.EndGroup
	}
	newLoop, err := strconv.Atoi(args[1])
	if err != nil || newLoop < 0 || newLoop > 10 {
		m.Reply("⚠️ <b>Invalid loop count.</b>\nUse 0 to disable or 1-10 to set loop count.")
		return telegram.EndGroup
	}
	if newLoop == currentLoop {
		m.Reply(fmt.Sprintf("⚠️ Loop count is already set to <b>%d</b> time(s).", currentLoop))
		return telegram.EndGroup
	}
	r.Lock()
	r.Loop = newLoop
	r.Unlock()
	mention := utils.MentionHTML(m.Sender)
	msg := "🔁 Loop has been <b>disabled</b> by " + mention
	if newLoop > 0 {
		msg = fmt.Sprintf("🔁 Set to loop <b>%d</b> time(s)\n└ Changed by: %s", newLoop, mention)
	}
	m.Reply(msg)
	return telegram.EndGroup
}
