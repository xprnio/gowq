package state

type WorkListState interface{}

type WorkListMovingState struct {
	Source int
	Active int
}
