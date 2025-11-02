package database

func IsLoggerEnabled() (bool, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return false, err
	}
	return state.LoggerEnabled, nil
}

func SetLoggerEnabled(enabled bool) error {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil || state.LoggerEnabled == enabled {
		return err
	}

	state.LoggerEnabled = enabled
	return updateBotState(ctx, state)
}
