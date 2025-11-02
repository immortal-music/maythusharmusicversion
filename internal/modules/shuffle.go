package modules

import (
	"fmt"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func shuffleHandler(m *telegram.NewMessage) error {
	return handleShuffle(m, false)
}

func cshuffleHandler(m *telegram.NewMessage) error {
	return handleShuffle(m, true)
}

func handleShuffle(m *telegram.NewMessage, cplay bool) error {
	arg := strings.ToLower(m.Args())

	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.EndGroup
	}

	if r.Track == nil {
		m.Reply("⚠️ No active playback found.")
		return telegram.EndGroup
	}

	r.Parse()

	if arg == "" {
		state := "disabled ❌"
		cmd := getCommand(m) + " on"
		if r.Shuffle {
			state = "enabled ✅"
			cmd = getCommand(m) + " off"
		}

		m.Reply(fmt.Sprintf(
			"🔀 Currently shuffle is <b>%s</b> for this chat.\n\nUse <code>%s</code> to toggle it.",
			state, cmd,
		))
		return telegram.EndGroup
	}

	var newState bool
	if arg == "on" || arg == "enable" || arg == "true" || arg == "1" {
		newState = true
	} else if arg == "off" || arg == "disable" || arg == "false" || arg == "0" {
		newState = false
	}

	r.SetShuffle(newState)

	state := "disabled ❌"
	if newState {
		state = "enabled ✅"
	}

	m.Reply(fmt.Sprintf("🔀 Shuffle <b>%s</b> by %s.", state, utils.MentionHTML(m.Sender)))
	return telegram.EndGroup
}
