package utils

import (
	"fmt"
	"html"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

func ShortTitle(title string, max ...int) string {
	limit := 25
	if len(max) > 0 {
		limit = max[0]
	}
	runes := []rune(title)
	if len(runes) <= limit {
		return title
	}
	return string(runes[:limit]) + "..."
}

func CleanURL(raw string) string {
	parts := strings.SplitN(raw, "?", 2)
	return parts[0]
}

func MentionHTML(u *telegram.UserObj) string {
	if u == nil {
		return "Unknown"
	}

	fullName := u.FirstName
	if u.LastName != "" {
		fullName += " " + u.LastName
	}

	if fullName == "" {
		fullName = "User"
	}
	fullName = html.EscapeString(ShortTitle(fullName, 15))

	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, u.ID, fullName)
}
