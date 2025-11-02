package database

func GetAutoLeave() (bool, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return false, err
	}
	return state.AutoLeave, nil
}

func SetAutoLeave(value bool) error {
	ctx, cancel := mongoCtx()
	defer cancel()
	current, err := GetAutoLeave()
	if err != nil {
		logger.ErrorF("Failed to get current AutoEnd: %v", err)
		return err
	}

	if current == value {
		return nil
	}

	state, err := getBotState(ctx)
	if err != nil {
		logger.ErrorF("Failed to get bot state for setting AutoEnd: %v", err)
		return err
	}

	state.AutoLeave = value
	if err := updateBotState(ctx, state); err != nil {
		logger.ErrorF("Failed to update AutoEnd: %v", err)
		return err
	}
	return nil
}
