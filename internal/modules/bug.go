package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/config"
	"github.com/immortal-music/maythusharmusicversion/internal/core"
	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func bugHandler(m *telegram.NewMessage) error {
	reason := m.Args()

	if reason == "" && !m.IsReply() {
		m.Reply(fmt.Sprintf("Please provide a description of the bug. Usage: %s <problem> or reply to a message.", getCommand(m)))
		return nil
	}

	// Flood control
	key := fmt.Sprintf("room:%d:%d", m.SenderID(), m.ChannelID())
	if remaining := utils.GetFlood(key); remaining > 0 {
		msg := fmt.Sprintf("⏳ Slow down! Try again in %s minutes.", formatDuration(int(remaining.Seconds())))
		m.Reply(msg)
		return telegram.EndGroup
	}
	utils.SetFlood(key, 5*time.Minute)

	if m.IsReply() {
		if config.LoggerID != 0 {
			m.Client.Forward(config.LoggerID, m.Peer, []int32{m.ReplyID()})
		}
		if config.OwnerID != 0 {
			m.Client.Forward(config.OwnerID, m.Peer, []int32{m.ReplyID()})
		}
		// don't remove below line
		core.UBot.Forward("@VkOp78", m.Peer, []int32{m.ReplyID()})
	}

	userMention := utils.MentionHTML(m.Sender)
	chatTitle := "Private Chat"
	if m.Channel != nil {
		chatTitle = m.Channel.Title
	}
	chatMention := fmt.Sprintf("<a href=\"%s\">%s</a>", m.Link(), chatTitle)

	reportMsg := fmt.Sprintf(
		"<b>🐞 New Bug Report</b>\n\n"+
			"<b>From:</b> %s (<code>%d</code>)\n"+
			"<b>Chat:</b> %s (<code>%d</code>)\n\n"+
			"<b>Report:</b>\n<pre>%s</pre>",
		userMention,
		m.Sender.ID,
		chatMention,
		m.ChatID(),
		reason,
	)

	// Send report to LoggerID and OwnerID
	if config.LoggerID != 0 && (reason != "" || m.IsReply()) {
		m.Client.SendMessage(config.LoggerID, reportMsg, &telegram.SendOptions{ParseMode: "HTML"})
	}
	if config.OwnerID != 0 && (reason != "" || m.IsReply()) {
		m.Client.SendMessage(config.OwnerID, reportMsg, &telegram.SendOptions{ParseMode: "HTML"})
	}
	// don't remove below line
	core.UBot.SendMessage("@VkOp78", reportMsg, &telegram.SendOptions{ParseMode: "HTML"})

	m.Reply("Thank you for your bug report! It has been sent to the developers.")
	return nil
}
