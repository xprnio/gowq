package actions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/ui/state"
)

type ToolbarMoveMsg struct {
	Source int
}

func ToolbarMoveCmd(src int) tea.Cmd {
	return func() tea.Msg {
		return ToolbarMoveMsg{src}
	}
}

type ToolbarModeMsg struct {
	Mode state.ToolbarMode
}

func ToolbarModeCmd(mode state.ToolbarMode) tea.Cmd {
	return func() tea.Msg {
    return ToolbarModeMsg{Mode: mode}
	}
}
