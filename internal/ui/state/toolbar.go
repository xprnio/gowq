package state

import (
	"github.com/xprnio/work-queue/internal/database"
)

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

