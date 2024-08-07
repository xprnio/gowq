package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/ui/header"
	"github.com/xprnio/work-queue/internal/ui/toolbar"
	"github.com/xprnio/work-queue/internal/ui/worklist"
	"github.com/xprnio/work-queue/internal/wq"
)

type Model struct {
	manager *wq.Manager

	header  *header.Model
	toolbar *toolbar.Model
	list    *worklist.Model

	width, height int
}

func New(manager *wq.Manager) *Model {
	return &Model{
		manager: manager,
		header:  header.New(manager),
		list:    worklist.New(manager),
		toolbar: toolbar.New(),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.toolbar.Init(),
		m.list.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if _, cmd := m.header.Update(msg); cmd != nil {

	}

	if _, cmd := m.toolbar.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	if _, cmd := m.list.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.header.View(),
		m.list.View(),
		m.toolbar.View(),
	)
}
