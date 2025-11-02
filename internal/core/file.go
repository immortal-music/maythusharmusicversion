package core

import (
	"os"
	"path/filepath"

	"github.com/immortal-music/maythusharmusicversion/internal/state"
)

func (r *RoomState) releaseFile() {
	if r == nil || r.Track == nil || r.FilePath == "" {
		return
	}

	vid := r.Track.ID
	path := r.FilePath

	roomsMu.RLock()
	shouldRemove := true
	for _, room := range rooms {
		if room == nil || room.Track == nil {
			continue
		}

		// same room: if still queued in first 5, don’t remove
		n := len(room.Queue)
		if n > 5 {
			n = 5
		}
		if room.ChatID == r.ChatID {
			for _, q := range room.Queue[:n] {
				if q.ID == vid {
					shouldRemove = false
					break
				}
			}
			continue
		}

		// another room currently playing or queued
		if room.Track.ID == vid {
			shouldRemove = false
			break
		}
		for _, q := range room.Queue[:n] {
			if q.ID == vid {
				shouldRemove = false
				break
			}
		}
		if !shouldRemove {
			break
		}
	}
	roomsMu.RUnlock()

	if shouldRemove {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			logger.ErrorF("failed to remove file %s: %v", path, err)
		} else {
			logger.DebugF("removed unused file: %s", path)
		}
	} else {
		logger.DebugF("file still in use, skipped remove: %s", path)
	}
}

func (r *RoomState) cleanupFile() {
	if r == nil {
		return
	}

	// Collect up to 5 tracks: current + queue
	tracks := []*state.Track{}
	if r.Track != nil {
		tracks = append(tracks, r.Track)
	}
	tracks = append(tracks, r.Queue...)
	if len(tracks) > 5 {
		tracks = tracks[:5]
	}

	var toDelete []string

	roomsMu.RLock()
	for _, t := range tracks {
		if t == nil || t.ID == "" {
			continue
		}

		used := false
		for _, other := range rooms {
			if other == nil || other.Track == nil || other.ChatID == r.ChatID {
				continue
			}

			// playing
			if other.Track.ID == t.ID {
				used = true
				break
			}

			// queued
			n := len(other.Queue)
			if n > 5 {
				n = 5
			}
			for _, q := range other.Queue[:n] {
				if q.ID == t.ID {
					used = true
					break
				}
			}
			if used {
				break
			}
		}

		if !used {
			pattern := filepath.Join("downloads", t.ID+".*")
			matches, err := filepath.Glob(pattern)
			if err != nil {
				logger.ErrorF("glob failed for %s: %v", pattern, err)
				continue
			}
			toDelete = append(toDelete, matches...)
		} else {
			logger.DebugF("track %s still in use, skip delete", t.ID)
		}
	}
	roomsMu.RUnlock()

	for _, f := range toDelete {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			logger.ErrorF("failed to remove file %s: %v", f, err)
		} else {
			logger.DebugF("removed unused file: %s", f)
		}
	}
}
