package state

import (
	"github.com/xprnio/work-queue/internal/database"
)

type ToolbarState struct {
	ShowCompleted bool

	Mode ToolbarMode
	Err  error
}

func NewToolbarState() ToolbarState {
	return ToolbarState{
		Mode: ToolbarModeNormal{},
	}
}

type ToolbarMode interface{}

type ToolbarModeNormal struct{}
type ToolbarModeMove struct {
	Item *database.WorkItem
}
type ToolbarModeAdd struct{}
type ToolbarModeEdit struct {
	Index int
	Name  string
}
type ToolbarModeComplete struct{}
type ToolbarModeDelete struct{}
