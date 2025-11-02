package modules

import (
	"fmt"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/database"
	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func handleLogger(m *telegram.NewMessage) error {
	args := strings.Fields(m.Text())
	current, dbErr := database.IsLoggerEnabled()

	if len(args) < 2 {
		if dbErr == nil {
			status := "🟢 Enabled"
			if !current {
				status = "🔴 Disabled"
			}
			m.Reply(
				fmt.Sprintf("⚙️ Usage: <code>%s [enable|disable]</code> - To enable or disable the logger\n\n📜 Current status: %s", getCommand(m), status),
			)
		} else {
			m.Reply(fmt.Sprintf("⚙️ Usage: <code>/%s [enable|disable]</code> - To enable or disable the logger", getCommand(m)))
		}
		return telegram.EndGroup
	}

	enable, err := utils.ParseBool(args[1])
	if err != nil {
		m.Reply("⚠️ Invalid option. Use 'enable' or 'disable'.")
		return telegram.EndGroup
	}

	if dbErr != nil {
		m.Reply("❌ Failed to check logger status: " + dbErr.Error())
		return telegram.EndGroup
	}

	if current == enable {
		status := "enabled"
		if !enable {
			status = "disabled"
		}
		m.Reply("ℹ️ Logger is already " + status + ".")
		return telegram.EndGroup
	}

	if err := database.SetLoggerEnabled(enable); err != nil {
		m.Reply("❌ Failed to update logger setting: " + err.Error())
		return telegram.EndGroup
	}

	status := "disabled"
	if enable {
		status = "enabled"
	}
	m.Reply("✅ Logger has been " + status + " successfully.")
	return telegram.EndGroup
}
