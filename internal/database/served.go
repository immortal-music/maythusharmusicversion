package database

func GetServed(user ...bool) ([]int64, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return nil, err
	}

	if len(user) > 0 && user[0] {
		return state.Served.Users, nil
	}
	return state.Served.Chats, nil
}

func IsServed(id int64, user ...bool) (bool, error) {
	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return false, err
	}

	target := state.Served.Chats
	if len(user) > 0 && user[0] {
		target = state.Served.Users
	}

	for _, v := range target {
		if v == id {
			return true, nil
		}
	}
	return false, nil
}

func AddServed(id int64, user ...bool) error {
	exists, err := IsServed(id, user...)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return err
	}

	target := &state.Served.Chats
	if len(user) > 0 && user[0] {
		target = &state.Served.Users
	}

	*target = append(*target, id)
	return updateBotState(ctx, state)
}

func DeleteServed(id int64, user ...bool) error {
	exists, err := IsServed(id, user...)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	ctx, cancel := mongoCtx()
	defer cancel()

	state, err := getBotState(ctx)
	if err != nil {
		return err
	}

	target := &state.Served.Chats
	if len(user) > 0 && user[0] {
		target = &state.Served.Users
	}

	newSlice := make([]int64, 0, len(*target))
	for _, v := range *target {
		if v != id {
			newSlice = append(newSlice, v)
		}
	}
	*target = newSlice
	return updateBotState(ctx, state)
}
