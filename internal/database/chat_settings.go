package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ChatSettings struct {
	ChatID    int64   `bson:"_id"`
	CPlayID   int64   `bson:"cplay_id"`
	AuthUsers []int64 `bson:"auth_users"`
}

func defaultChatSettings(chatID int64) *ChatSettings {
	return &ChatSettings{
		ChatID:    chatID,
		CPlayID:   0,
		AuthUsers: []int64{},
	}
}

func getChatSettings(ctx context.Context, chatID int64) (*ChatSettings, error) {
	cacheKey := fmt.Sprintf("chat_settings_%d", chatID)
	if cached, found := dbCache.Get(cacheKey); found {
		if settings, ok := cached.(*ChatSettings); ok {
			return settings, nil
		}
	}

	var settings ChatSettings
	err := chatSettingsColl.FindOne(ctx, bson.M{"_id": chatID}).Decode(&settings)
	if err == mongo.ErrNoDocuments {
		def := defaultChatSettings(chatID)
		dbCache.Set(cacheKey, def)
		return def, nil
	} else if err != nil {
		logger.ErrorF("Failed to get chat settings for chat %d: %v", chatID, err)
		return nil, err
	}

	dbCache.Set(cacheKey, &settings)

	// Proactively cache the cplayID -> chatID mapping
	if settings.CPlayID != 0 {
		cplayCacheKey := fmt.Sprintf("cplayid_%d", settings.CPlayID)
		dbCache.Set(cplayCacheKey, settings.ChatID)
	}

	return &settings, nil
}

func updateChatSettings(ctx context.Context, newSettings *ChatSettings) error {
	cacheKey := fmt.Sprintf("chat_settings_%d", newSettings.ChatID)
	opts := options.UpdateOne().SetUpsert(true)

	_, err := chatSettingsColl.UpdateOne(ctx, bson.M{"_id": newSettings.ChatID}, bson.M{"$set": newSettings}, opts)
	if err != nil {
		logger.ErrorF("Failed to update chat settings for chat %d: %v", newSettings.ChatID, err)
		return err
	}

	dbCache.Set(cacheKey, newSettings)
	return nil
}
