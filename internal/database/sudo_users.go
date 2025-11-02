package database

// GetSudoers returns all sudoers.
func GetSudoers() ([]int64, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		logger.ErrorF("Failed to get current sudoers: %v", err)
		return nil, err
	}
	return state.Sudoers, nil
}

// IsSudo checks if the given ID is a sudoer.
func IsSudo(id int64) (bool, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		logger.ErrorF("Failed to get current sudoers: %v", err)
		return false, err
	}

	for _, v := range state.Sudoers {
		if v == id {
			return true, nil
		}
	}
	return false, nil
}

// AddSudo adds a new sudoer if not already present.
func AddSudo(id int64) error {
	exists, err := IsSudo(id)
	if err != nil {
		logger.ErrorF("Failed to check sudo existence: %v", err)
		return err
	}
	if exists {
		return nil
	}

	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		logger.ErrorF("Failed to get current sudoers: %v", err)
		return err
	}

	state.Sudoers = append(state.Sudoers, id)
	if err := updateBotState(ctx, state); err != nil {
		logger.ErrorF("Failed to update sudoers: %v", err)
		return err
	}

	return nil
}

// DeleteSudo removes a sudoer by ID.
func DeleteSudo(id int64) error {
	exists, err := IsSudo(id)
	if err != nil {
		logger.ErrorF("Failed to check sudo existence: %v", err)
		return err
	}
	if !exists {
		return nil
	}

	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		logger.ErrorF("Failed to get current sudoers: %v", err)
		return err
	}

	newSudoers := make([]int64, 0, len(state.Sudoers))
	for _, v := range state.Sudoers {
		if v != id {
			newSudoers = append(newSudoers, v)
		}
	}
	state.Sudoers = newSudoers

	if err := updateBotState(ctx, state); err != nil {
		logger.ErrorF("Failed to update sudoers: %v", err)
		return err
	}

	return nil
}
