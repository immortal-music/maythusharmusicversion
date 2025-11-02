package modules

import (
	"strconv"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/core"
)

func activeHandler(m *telegram.NewMessage) error {
	chats := len(core.GetAllRoomIDs())
	m.Reply("Active Chats info: " + strconv.Itoa(chats))
	return telegram.EndGroup
}
