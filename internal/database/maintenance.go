package database

// SetMaintenance sets the maintenance mode.
// If enabling, you can provide an optional reason.
// If disabling, it clears any existing reason.
func SetMaintenance(enabled bool, reason ...string) error {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return err
	}

	if state.Maintenance == enabled {
		// already in desired state, just update reason if enabling
		if enabled && len(reason) > 0 {
			state.MaintReason = reason[0]
			return updateBotState(ctx, state)
		} else if !enabled && state.MaintReason != "" {
			// clear reason when disabling
			state.MaintReason = ""
			return updateBotState(ctx, state)
		}
		return nil
	}

	state.Maintenance = enabled
	if enabled && len(reason) > 0 {
		state.MaintReason = reason[0]
	} else if !enabled {
		state.MaintReason = ""
	}

	return updateBotState(ctx, state)
}

// GetMaintReason retrieves the current maintenance reason
func GetMaintReason() (string, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return "", err
	}
	return state.MaintReason, nil
}

// IsMaintenance returns whether maintenance mode is enabled
func IsMaintenance() (bool, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return false, err
	}
	return state.Maintenance, nil
}
