package database

func IsAuthUser(chatID, userID int64) (bool, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	settings, err := getChatSettings(ctx, chatID)
	if err != nil {
		return false, err
	}

	for _, user := range settings.AuthUsers {
		if user == userID {
			return true, nil
		}
	}
	return false, nil
}

func AddAuthUser(chatID, userID int64) error {
	ctx, cancel := mongoCtx()
	defer cancel()

	settings, err := getChatSettings(ctx, chatID)
	if err != nil {
		return err
	}

	for _, user := range settings.AuthUsers {
		if user == userID {
			return nil // User already authorized
		}
	}

	settings.AuthUsers = append(settings.AuthUsers, userID)
	return updateChatSettings(ctx, settings)
}

func RemoveAuthUser(chatID, userID int64) error {
	ctx, cancel := mongoCtx()
	defer cancel()

	settings, err := getChatSettings(ctx, chatID)
	if err != nil {
		return err
	}

	var newAuthUsers []int64
	var found bool
	for _, user := range settings.AuthUsers {
		if user == userID {
			found = true
			continue
		}
		newAuthUsers = append(newAuthUsers, user)
	}

	if !found {
		return nil // User not in the auth list
	}

	settings.AuthUsers = newAuthUsers
	return updateChatSettings(ctx, settings)
}

func GetAuthUsers(chatID int64) ([]int64, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	settings, err := getChatSettings(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return settings.AuthUsers, nil
}
