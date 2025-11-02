package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetCPlayID(chatID int64) (int64, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	settings, err := getChatSettings(ctx, chatID)
	if err != nil {
		return 0, err
	}
	return settings.CPlayID, nil
}

func SetCPlayID(chatID, cplayID int64) error {
	ctx, cancel := mongoCtx()
	defer cancel()

	settings, err := getChatSettings(ctx, chatID)
	if err != nil {
		return err
	}

	if settings.CPlayID == cplayID {
		return nil
	}

	// Proactively update cache to prevent stale entries
	if oldCPlayID := settings.CPlayID; oldCPlayID != 0 {
		oldCacheKey := fmt.Sprintf("cplayid_%d", oldCPlayID)
		dbCache.Delete(oldCacheKey)
	}

	settings.CPlayID = cplayID
	err = updateChatSettings(ctx, settings)
	if err == nil {
		newCacheKey := fmt.Sprintf("cplayid_%d", cplayID)
		dbCache.Set(newCacheKey, chatID)
	}
	return err
}

func GetChatIDFromCPlayID(cplayID int64) (int64, error) {
	cacheKey := fmt.Sprintf("cplayid_%d", cplayID)
	if cached, found := dbCache.Get(cacheKey); found {
		if chatID, ok := cached.(int64); ok {
			return chatID, nil
		}
	}

	ctx, cancel := mongoCtx()
	defer cancel()

	var settings ChatSettings
	err := chatSettingsColl.FindOne(ctx, bson.M{"cplay_id": cplayID}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, fmt.Errorf("no chat found with cplayID %d", cplayID)
		}
		return 0, err
	}

	dbCache.Set(cacheKey, settings.ChatID)
	return settings.ChatID, nil
}
