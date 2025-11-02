package modules

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/core"
	"github.com/immortal-music/maythusharmusicversion/internal/database"
	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

var maintenanceCancel = struct {
	sync.Mutex
	cancel bool
}{}

func handleMaintenance(m *telegram.NewMessage) error {
	args := strings.Fields(m.Text())
	current, err := database.IsMaintenance()
	if err != nil {
		m.Reply("❌ Failed to check maintenance status: " + err.Error())
		return telegram.EndGroup
	}

	if len(args) < 2 {

		status := "🔴 Disabled"
		if current {
			if reason, rerr := database.GetMaintReason(); rerr == nil && reason != "" {
				status = fmt.Sprintf("🟢 Enabled\n📝 Reason: %s", reason)
			} else {
				status = "🟢 Enabled"
			}
		}

		m.Reply(fmt.Sprintf(
			"⚙️ Usage: %s [<code>enable</code>|<code>disable</code>] [reason]\n\n📜 Current status: %s",
			getCommand(m),
			status,
		))
		return telegram.EndGroup
	}

	enable, err := utils.ParseBool(args[1])
	if err != nil {
		m.Reply("⚠️ Invalid option. Use 'enable' or 'disable'.")
		return telegram.EndGroup
	}
	reason := strings.Join(args[2:], " ")

	oldReason, _ := database.GetMaintReason()

	if current == enable {
		if enable {
			switch {
			case reason == oldReason:
				m.Reply("ℹ️ Maintenance mode is already enabled with the same reason.")
				return telegram.EndGroup
			case reason == "" && oldReason != "":
				_ = database.SetMaintenance(true, "")
				m.Reply("✅ Maintenance reason removed successfully.")
				return telegram.EndGroup
			case reason != "" && reason != oldReason:
				_ = database.SetMaintenance(true, reason)
				m.Reply(fmt.Sprintf("✅ Maintenance reason updated successfully.\n📝 Reason: %s", reason))
				return telegram.EndGroup
			default:
				m.Reply("ℹ️ Maintenance mode is already enabled 🟢.")
				return telegram.EndGroup
			}
		} else {
			m.Reply("ℹ️ Maintenance mode is already disabled 🔴.")
			return telegram.EndGroup
		}
	}

	_ = database.SetMaintenance(enable, reason)
	logger.InfoF("User %d set maintenance mode to %v. Reason: %s", m.SenderID(), enable, reason)

	if enable {
		maintenanceCancel.Lock()
		maintenanceCancel.cancel = false
		maintenanceCancel.Unlock()

		go func(c *telegram.Client, reason string) {
			for _, id := range core.GetAllRoomIDs() {
				maintenanceCancel.Lock()
				if maintenanceCancel.cancel {
					maintenanceCancel.Unlock()
					break
				}
				maintenanceCancel.Unlock()

				if r, ok := core.GetRoom(id); ok {

					r.Destroy()
					if reason != "" {
						c.SendMessage(id, "⚠️ Bot is entering maintenance mode.\n📝 Reason: "+reason)

						time.Sleep(1 * time.Second)
					}
				}
			}
		}(m.Client, reason)

		msg := "✅ Maintenance mode enabled successfully."
		if reason != "" {
			msg += fmt.Sprintf("\n📝 Reason: %s", reason)
		}
		m.Reply(msg)
		return telegram.EndGroup
	}

	maintenanceCancel.Lock()
	maintenanceCancel.cancel = true
	maintenanceCancel.Unlock()

	m.Reply("✅ Maintenance mode disabled successfully.")
	return telegram.EndGroup
}
