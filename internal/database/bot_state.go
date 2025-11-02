package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UsersChats struct {
	Users []int64 `bson:"users"`
	Chats []int64 `bson:"chats"`
}

type BotState struct {
	ID            string     `bson:"_id"`
	Served        UsersChats `bson:"served"`
	Sudoers       []int64    `bson:"sudoers"`
	AutoLeave     bool       `bson:"autoleave"`
	LoggerEnabled bool       `bson:"logger"`
	Maintenance   bool       `bson:"maintenance"`
	MaintReason   string     `bson:"maint_reason"`
}

const cacheKey = "bot_state"

func defaultState() *BotState {
	return &BotState{
		ID:            "global",
		Served:        UsersChats{Users: []int64{}, Chats: []int64{}},
		Sudoers:       []int64{},
		LoggerEnabled: true,
	}
}

func getBotState(ctx context.Context) (*BotState, error) {
	if cached, found := dbCache.Get(cacheKey); found {
		if state, ok := cached.(*BotState); ok {
			return state, nil
		}
	}

	var state BotState
	err := settingsColl.FindOne(ctx, bson.M{"_id": "global"}).Decode(&state)
	if err == mongo.ErrNoDocuments {
		def := defaultState()
		dbCache.Set(cacheKey, def)
		return def, nil
	} else if err != nil {
		logger.ErrorF("Failed to get bot state: %v", err)
		return nil, err
	}

	dbCache.Set(cacheKey, &state)
	return &state, nil
}

func updateBotState(ctx context.Context, newState *BotState) error {
	opts := options.UpdateOne().SetUpsert(true)

	_, err := settingsColl.UpdateOne(ctx, bson.M{"_id": "global"}, bson.M{"$set": newState}, opts)
	if err != nil {
		logger.ErrorF("Failed to update bot state: %v", err)
		return err
	}

	dbCache.Set(cacheKey, newState)
	return nil
}
