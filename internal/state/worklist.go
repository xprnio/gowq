package state

type WorkListState struct {
	ShowNumbers   bool
	ShowCompleted bool

  Focused *workListFocusedState
	Moving *workListMovingState
}

type workListFocusedState struct {
	Index int
}

type workListMovingState struct {
	Source int
	Active int
}

func WorkListFocusedState(index int) *workListFocusedState {
	return &workListFocusedState{
		Index: index,
	}
}

func WorkListMovingState(source int) *workListMovingState {
	return &workListMovingState{
		Source: source,
		Active: source,
	}
}
