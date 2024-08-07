package wq

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

var InvalidIndexErr = errors.New("invalid index")

type ErrorMsg struct {
	Err error
}

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{err}
	}
}
