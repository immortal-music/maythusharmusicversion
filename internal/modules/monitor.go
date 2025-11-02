package modules

import (
	"time"

	"github.com/Laky-64/gologging"
	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/core"
)

var logger = gologging.GetLogger("monitor")

func MonitorRooms() {
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	sem := make(chan struct{}, 20)

	for range ticker.C {
		for _, chatID := range core.GetAllRoomIDs() {
			sem <- struct{}{}
			go func(id int64) {
				defer func() { <-sem }()

				r, ok := core.GetRoom(id)
				if !ok || !r.IsActiveChat() || r.IsPaused() {
					return
				}

				r.Parse()
				mystic := r.GetMystic()
				if mystic == nil {
					return
				}

				markup := core.GetPlayMarkup(r, false)
				opts := telegram.SendOptions{
					ReplyMarkup: markup,
					Entities:    mystic.Message.Entities,
				}
				mystic.Edit(mystic.Text(), opts)
			}(chatID)
		}
	}
}
