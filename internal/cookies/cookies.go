package cookies

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Laky-64/gologging"

	"github.com/immortal-music/maythusharmusicversion/config"
)

var logger = gologging.GetLogger("cookies")

func Init() {
	if err := os.MkdirAll("internal/cookies", 0o755); err != nil {
		logger.Fatal("Failed to create cookies directory:", err)
	}

	urls := strings.Fields(config.CookiesLink)
	for _, url := range urls {
		if err := downloadCookieFile(url); err != nil {
			logger.WarnF("Failed to download cookie file from %s: %v", url, err)
		}
	}
}

func downloadCookieFile(url string) error {
	id := filepath.Base(url)
	rawURL := "https://batbin.me/raw/" + id

	resp, err := http.Get(rawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filePath := filepath.Join("internal/cookies", id+".txt")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func GetRandomCookieFile() (string, error) {
	files, err := filepath.Glob("internal/cookies/*.txt")
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", nil // No cookie files found
	}
	return files[rand.Intn(len(files))], nil
}
