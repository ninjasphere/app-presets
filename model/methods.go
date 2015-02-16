package model

// given a thing state, produce a new thing state with the same
// channel states as the receiver, but the undo states matching the states
// of the specified undo thing state.
func (m *ThingState) MergeUndoState(u *ThingState) *ThingState {
	result := &ThingState{
		ID:       m.ID,
		Channels: make([]ChannelState, len(m.Channels)),
	}

	tmp := make(map[string]*ChannelState)
	if u != nil {
		for _, ch := range u.Channels {
			tmp[ch.ID] = &ch
		}
	}

	for i, ch := range m.Channels {
		var undo *ChannelState = nil
		ok := false

		if undo, ok = tmp[ch.ID]; !ok {
			undo = nil
		}
		result.Channels[i] = *ch.MergeUndoState(undo)
	}
	return result
}

// given a channel state, produce a new channel state with the same
// state as the receiver, but with the undo state matching the state
// of the specified undo channel state.
func (m *ChannelState) MergeUndoState(u *ChannelState) *ChannelState {
	result := &ChannelState{
		ID:    m.ID,
		State: m.State,
	}
	if u != nil {
		result.UndoState = u.State
	}
	return result
}
