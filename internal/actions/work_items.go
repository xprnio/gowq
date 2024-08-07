package actions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/database"
)

type WorkAddedMsg struct {
	Work database.WorkItem
}

type WorkCompletedMsg struct {
	Index int
}

type WorkDeletedMsg struct {
	Index int
}

type WorkEditedMsg struct {
	Index int
	Name  string
}

type FinishMovingWorkMsg struct {
	Commit bool
}

type MovementDir string

const (
	MovementDirUp   MovementDir = "up"
	MovementDirDown MovementDir = "down"
)
type MoveWorkMsg struct {
	Direction MovementDir
}

func WorkAddedCmd(name string) tea.Cmd {
	return func() tea.Msg {
		item := database.WorkItem{Name: name}
		return WorkAddedMsg{item}
	}
}

func WorkCompletedCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return WorkCompletedMsg{index}
	}
}

func WorkDeletedCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return WorkDeletedMsg{index}
	}
}

func WorkEditedCmd(i int, name string) tea.Cmd {
	return func() tea.Msg {
		return WorkEditedMsg{i, name}
	}
}

func FinishMovingWorkCmd(commit bool) tea.Cmd {
	return func() tea.Msg {
		return FinishMovingWorkMsg{commit}
	}
}

func MoveWorkCmd(dir MovementDir) tea.Cmd {
	return func() tea.Msg {
		return MoveWorkMsg{dir}
	}
}

