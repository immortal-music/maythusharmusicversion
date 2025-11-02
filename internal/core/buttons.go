package core

import (
	"fmt"

	tg "github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/config"
	"github.com/immortal-music/maythusharmusicversion/internal/utils"
)

func AddMeMarkup(username string) tg.ReplyMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.URL("⚡ Add Me to Your startgroup",
				"https://t.me/"+username+"?startgroup=true",
			),
		).
		Build()
}

func GetPlayMarkup(r *RoomState, queued bool) tg.ReplyMarkup {
	btn := tg.NewKeyboard()
	prefix := "room:"
	if r.IsCPlay() {
		prefix = "croom:"
	}
	progress := utils.GetProgressBar(r.Position, r.Track.Duration)
	progress = formatDuration(r.Position) + " " + progress + " " + formatDuration(r.Track.Duration)

	if !queued {
		btn.AddRow(
			tg.Button.Data(progress, "progress"),
		)
	}
	btn.AddRow(
		tg.Button.Data("▷", prefix+"resume"),
		tg.Button.Data("II", prefix+"pause"),
		tg.Button.Data("‣‣I", prefix+"skip"),
		tg.Button.Data("▢", prefix+"stop"),
	)

	btn.AddRow(
		tg.Button.Data("↩ 15s", "room:seekback_15"),
		tg.Button.Data("⟳", "room:replay"),
		tg.Button.Data("15s ↪", "room:seek_15"),
	)

	btn.AddRow(
		tg.Button.Data("Close", "close"),
	)

	return btn.Build()
}

func GetGroupHelpKeyboard() *tg.ReplyInlineMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.URL("📒 Commands", "https://t.me/"+BUser.Username+"?start=help"),
		).
		Build()
}

func GetStartMarkup() tg.ReplyMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.URL("⚡ Add Me to Your startgroup",
				"https://t.me/"+BUser.Username+"?startgroup=true",
			),
		).
		AddRow(
			tg.Button.Data("❓ Help & Commands", "help_cb"),
		//	tg.Button.URL("💻 Source", config.RepoURL),
		).
		AddRow(
			tg.Button.URL("📢 Updates", config.SupportChannel),
			tg.Button.URL("💬 Support", config.SupportChat),
		).
		Build()
}

func GetHelpKeyboard() *tg.ReplyInlineMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.Data("🛠 Admins", "help:admins"),
			tg.Button.Data("🌍 Public", "help:public"),
		).
		AddRow(
			tg.Button.Data("👑 Owner", "help:owner"),
			tg.Button.Data("⚡ Sudoers", "help:sudoers"),
		).
		AddRow(tg.Button.Data("⬅️ Back", "start")).
		Build()
}

func GetBackKeyboard() *tg.ReplyInlineMarkup {
	return tg.NewKeyboard().
		AddRow(tg.Button.Data("⬅️ Back", "help:main")).
		Build()
}

func formatDuration(sec int) string {
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60

	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s) // HH:MM:SS
	}
	return fmt.Sprintf("%02d:%02d", m, s) // MM:SS
}
