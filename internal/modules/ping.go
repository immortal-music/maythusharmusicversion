package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/config"
	"github.com/immortal-music/maythusharmusicversion/internal/database"
)

func formatUptime(d time.Duration) string {
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm ", minutes)
	}
	result += fmt.Sprintf("%ds", seconds)
	return result
}

func pingHandler(m *telegram.NewMessage) error {
	if m.IsPrivate() {
		m.Delete()
		database.AddServed(m.ChannelID(), true)
	} else {
		database.AddServed(m.ChannelID())
	}
	start := time.Now()
	reply, err := m.Respond("🏓 Pinging...")
	if err != nil {
		return err
	}

	latency := time.Since(start).Milliseconds()
	uptime := time.Since(config.StartTime)
	uptimeStr := formatUptime(uptime)

	text := fmt.Sprintf(
		"🏓 Pong!\nLatency: %dms\n🤖 I've been running for %s without rest!",
		latency, uptimeStr,
	)

	reply.Edit(text)
	return telegram.EndGroup
}
