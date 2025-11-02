package modules

import (
	"fmt"
	"log"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/config"
	"github.com/immortal-music/maythusharmusicversion/internal/core"
	"github.com/immortal-music/maythusharmusicversion/internal/database"
	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

var startMSG = "⚡️Pika Pika, %s!\n⚡️  Welcome to <b>%s</b> \n🎶  I’m here to help you play, stream, and manage music right here on Telegram. 🎵"

func startHandler(m *telegram.NewMessage) error {
	if m.ChatType() != telegram.EntityUser {
		database.AddServed(m.ChannelID())
		m.Reply("🎶 I'm all set!\n▶️ Drop a command to light up the chat with music.")
		return telegram.EndGroup
	}

	arg := m.Args()
	database.AddServed(m.ChannelID(), true)

	switch arg {

	case "help":
		helpHandler(m)

	default:

		caption := fmt.Sprintf(startMSG, utils.MentionHTML(m.Sender), utils.MentionHTML(core.BUser))

		if _, err := m.RespondMedia(config.StartImage, telegram.MediaOptions{
			Caption:     caption,
			NoForwards:  true,
			ReplyMarkup: core.GetStartMarkup(),
		}); err != nil {
			log.Printf("Error responding start in chat: %v", err)
			return err
		}
	}

	return telegram.EndGroup
}

func startCB(c *telegram.CallbackQuery) error {
	c.Answer("")

	caption := fmt.Sprintf(startMSG, utils.MentionHTML(c.Sender), utils.MentionHTML(core.BUser))

	opt := &telegram.SendOptions{
		ReplyMarkup: core.GetStartMarkup(),
		NoForwards:  true,
	}

	if config.StartImage != "" {
		opt.Media = config.StartImage
	}
	c.Edit(caption, opt)
	return telegram.EndGroup
}
