package toolbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/database"
)

type ToolbarMode interface{}

type ModeNormal struct{}
type ModeMove struct {
	Item *database.WorkItem
}
type ModeAdd struct{}
type ModeEdit struct {
	Index int
	Name  string
}
type ModeComplete struct{}
type ModeDelete struct{}

type ToolbarModeMsg struct {
	Mode ToolbarMode
}

func ToolbarModeCmd(mode ToolbarMode) tea.Cmd {
	return func() tea.Msg {
		return ToolbarModeMsg{mode}
	}
}
