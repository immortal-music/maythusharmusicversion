package modules

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Laky-64/gologging"
	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/config"
	"github.com/immortal-music/maythusharmusicversion/internal/core"
	"github.com/immortal-music/maythusharmusicversion/internal/database"
	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func handleAddSudo(m *telegram.NewMessage) error {
	if m.Args() == "" && !m.IsReply() {
		m.Reply("⚠️ Please provide a user — use:\n" + getCommand(m) + " [user_id]</code> or reply to a user's message.")
		return telegram.EndGroup
	}
	targetID, err := utils.ExtractUser(m)
	if err != nil {
		m.Reply("❌ Failed to extract user: " + err.Error())
		return telegram.EndGroup
	}

	if targetID == config.OwnerID {
		m.Reply("😂 Haha, you’re the boss already! Why are you even trying to add yourself?")
		return telegram.EndGroup
	}

	if targetID == core.BUser.ID {
		m.Reply("🤖 Haha, I’m the bot! I can’t add myself to the sudo list — even the owner can’t cheat me 😎")
		return telegram.EndGroup
	}

	user, err := m.Client.GetUser(targetID)
	if err != nil {
		m.Reply("❌ Failed to fetch user info. Maybe the user is inaccessible.")
		gologging.Error("Failed to get user: " + err.Error())
		return telegram.EndGroup
	}

	if user.Bot {
		m.Reply("🤖 You can’t add a bot to the sudo list — sudoers must be human!")
		return telegram.EndGroup
	}

	uname := utils.MentionHTML(user)
	if user.Username != "" {
		uname = "@" + user.Username
	}

	exists, err := database.IsSudo(targetID)
	if err != nil {
		m.Reply("❌ Failed to check sudo existence: " + err.Error())
		return telegram.EndGroup
	}

	if exists {
		m.Reply("⚠️ User " + uname + " (ID:<code>" + strconv.FormatInt(targetID, 10) + "</code>) is already a sudoer.")
		return telegram.EndGroup
	}

	if err := database.AddSudo(targetID); err != nil {
		m.Reply("❌ Failed to add sudo: " + err.Error())
		return telegram.EndGroup
	}

	m.Reply("✅ Added " + uname + " (<code>" + strconv.FormatInt(targetID, 10) + "</code>) to sudoers.")

	sudoCommands := append(AllCommands.PrivateUserCommands, AllCommands.PrivateSudoCommands...)

	if _, err := m.Client.BotsSetBotCommands(&telegram.BotCommandScopePeer{
		Peer: &telegram.InputPeerUser{UserID: targetID, AccessHash: 0},
	},
		"",
		sudoCommands,
	); err != nil {
		gologging.Error("Failed to set PrivateSudoCommands " + err.Error())
	}
	return telegram.EndGroup
}

func handleDelSudo(m *telegram.NewMessage) error {
	if m.Args() == "" && !m.IsReply() {
		m.Reply("⚠️ Please provide a user — use:\n" + getCommand(m) + " [user_id]</code> or reply to a user's message.")
		return telegram.EndGroup
	}
	targetID, err := utils.ExtractUser(m)
	if err != nil {
		m.Reply("❌ Failed to extract user: " + err.Error())
		return telegram.EndGroup
	}

	if targetID == config.OwnerID {
		m.Reply("😎 Nice try! You can’t remove yourself from being the owner — you’re untouchable.")
		return telegram.EndGroup
	}

	if targetID == core.BUser.ID {
		m.Reply("😂 I can’t remove myself from the sudo list.")
		return telegram.EndGroup
	}

	user, err := m.Client.GetUser(targetID)
	if err != nil {
		m.Reply("❌ Failed to fetch user info. Maybe the user is hidden or inaccessible.")
		return telegram.EndGroup
	}

	uname := utils.MentionHTML(user)
	if user.Username != "" {
		uname = "@" + user.Username
	}

	exists, err := database.IsSudo(targetID)
	if err != nil {
		m.Reply("❌ Failed to check sudo existence: " + err.Error())
		return telegram.EndGroup
	}

	if !exists {
		m.Reply("⚠️ User " + uname + " (<code>" + strconv.FormatInt(targetID, 10) + "</code>) is not a sudoer.")
		return telegram.EndGroup
	}

	if _, err := m.Client.BotsResetBotCommands(&telegram.BotCommandScopePeer{
		Peer: &telegram.InputPeerUser{UserID: targetID, AccessHash: 0},
	}, ""); err != nil {
		gologging.Error("Failed to reset sudo commands: " + err.Error())
	}

	if err := database.DeleteSudo(targetID); err != nil {
		m.Reply("❌ Failed to remove sudo: " + err.Error())
		return telegram.EndGroup
	}

	m.Reply("🗑️ Removed " + uname + " (ID: <code>" + strconv.FormatInt(targetID, 10) + "</code>) from sudoers.")
	return telegram.EndGroup
}

func handleGetSudoers(m *telegram.NewMessage) error {
	floodKey := fmt.Sprintf("sudoers:%d%d", m.ChannelID(), m.SenderID())
	if remaining := utils.GetFlood(floodKey); remaining > 0 {
		return m.E(m.Reply("⏳ Please wait " + strconv.Itoa(int(remaining.Seconds())) + "seconds before using this command again."))
	}
	utils.SetFlood(floodKey, 30*time.Second)

	mystic, _ := m.Reply("⏳ Fetching sudoers list...")
	list, err := database.GetSudoers()
	if err != nil {
		utils.EOR(mystic, "❌ Failed to get sudoers: "+err.Error())
		return telegram.EndGroup
	}

	var sb strings.Builder
	sb.WriteString("👑 <b>Current Sudoers:</b>\n\n")

	// First, show owner
	ownerStr := "<code>" + strconv.FormatInt(config.OwnerID, 10) + "</code>"
	user, err := m.Client.GetUser(config.OwnerID)
	if err == nil {
		if user.Username != "" {
			ownerStr = "@" + user.Username + " (ID: <code>" + strconv.FormatInt(config.OwnerID, 10) + "</code>)"
		} else {
			ownerStr = utils.MentionHTML(user) + " (ID: <code>" + strconv.FormatInt(config.OwnerID, 10) + "</code>)"
		}
	}
	sb.WriteString("1. " + ownerStr + " — <b>Owner</b>\n")

	// Then list other sudoers
	idx := 2
	for _, id := range list {
		if id == config.OwnerID {
			continue // skip owner since already listed
		}

		userStr := "<code>" + strconv.FormatInt(id, 10) + "</code>" // fallback
		user, err := m.Client.GetUser(id)
		if err == nil {
			if user.Username != "" {
				userStr = "@" + user.Username + " (ID: <code>" + strconv.FormatInt(id, 10) + "</code>)"
			} else {
				userStr = utils.MentionHTML(user) + " (ID: <code>" + strconv.FormatInt(id, 10) + "</code>)"
			}
		}

		sb.WriteString(strconv.Itoa(idx) + ". " + userStr + "\n")
		idx++
	}

	if idx == 2 {
		sb.WriteString("⚠️ No additional sudoers found.\n")
	}

	utils.EOR(mystic, sb.String())
	return telegram.EndGroup
}
