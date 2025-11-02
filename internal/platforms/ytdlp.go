package platforms

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Laky-64/gologging"
	"github.com/amarnathcjd/gogram/telegram"

	"github.com/immortal-music/maythusharmusicversion/internal/cookies"
	"github.com/immortal-music/maythusharmusicversion/internal/state"
)

type YtDlpPlatform struct{}

func init() {
	addPlatform(70, state.PlatformYtDlp, &YtDlpPlatform{})
}

func (*YtDlpPlatform) Name() state.PlatformName {
	return state.PlatformYtDlp
}

func (*YtDlpPlatform) IsValid(query string) bool {
	return false
}

func (*YtDlpPlatform) GetTracks(query string) ([]*state.Track, error) {
	return nil, errors.New("YtDlp is a download-only platform")
}

func (*YtDlpPlatform) IsDownloadSupported(source state.PlatformName) bool {
	return source == state.PlatformYouTube
}

func (p *YtDlpPlatform) Download(ctx context.Context, track *state.Track, _ *telegram.NewMessage) (string, error) {
	if path, err := p.checkDownloadedFile(track.ID); err == nil {
		return path, nil
	}

	os.MkdirAll("downloads", os.ModePerm)
	filePath := filepath.Join("downloads", track.ID+".webm")

	args := []string{
		"-f", "bestaudio[ext=m4a]/bestaudio/best",
		"--no-playlist",
		"-o", filePath,
		"--geo-bypass",
		"--no-warnings",
		"--no-overwrites",
		"--ignore-errors",
		"--no-check-certificate",
		"-q",
		"--extractor-args", "youtube:player_js_version=actual",
	}

	cookieFile, err := cookies.GetRandomCookieFile()
	if err != nil {
		gologging.Error("Failed to get cookie file: " + err.Error())
		return "", err
	}
	if cookieFile != "" {
		args = append(args, "--cookies", cookieFile)
	}

	args = append(args, track.URL)
	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		outStr := stdout.String()
		errStr := stderr.String()
		gologging.ErrorF(
			"yt-dlp download failed for %s: %v\nSTDOUT:\n%s\nSTDERR:\n%s",
			track.URL, err, outStr, errStr,
		)
		return "", fmt.Errorf("yt-dlp error: %v\nstdout: %s\nstderr: %s", err, outStr, errStr)
	}

	return filePath, nil
}

func (f *YtDlpPlatform) checkDownloadedFile(videoId string) (string, error) {
	outputDir := "./downloads"
	pattern := filepath.Join(outputDir, videoId+".*")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search files: %v", err)
	}

	if len(matches) == 0 {
		return "", errors.New("❌ file not found")
	}

	return matches[0], nil
}
