package actions

import tea "github.com/charmbracelet/bubbletea"

type ToolbarMoveMsg struct {
	Source int
}

func ToolbarMoveCmd(src int) tea.Cmd {
	return func() tea.Msg {
		return ToolbarMoveMsg{src}
	}
}
