package actions

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ToggleNumbersActionMsg struct {
	Visible bool
}

type ToggleCompletedVisibilityMsg struct {
	Visible bool
}

type SetWorkListFocusMsg struct {
	Index int
}

type RefreshWorkListMsg struct{}

type ScrollWorkListMsg struct {
	Direction int
}

func ToggleNumbersCmd(visible bool) tea.Cmd {
	return func() tea.Msg {
		return ToggleNumbersActionMsg{visible}
	}
}

func ToggleCompletedVisibilityCmd(visible bool) tea.Cmd {
	return func() tea.Msg {
		return ToggleCompletedVisibilityMsg{visible}
	}
}

func SetWorkListFocusCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return SetWorkListFocusMsg{index}
	}
}

func RefreshWorkListCmd() tea.Cmd {
	return func() tea.Msg {
		return RefreshWorkListMsg{}
	}
}

func ScrollWorkListCmd(dir int) tea.Cmd {
	return func() tea.Msg {
		return ScrollWorkListMsg{dir}
	}
}
