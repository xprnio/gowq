package header

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/wq"
)

type Model struct {
	manager *wq.Manager
	style   lipgloss.Style
	width   int
	ready   bool
}

func New(man *wq.Manager) *Model {
	return &Model{
		manager: man,
		style:   baseStyle,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.style = m.style.Width(m.width)
		m.ready = true
	}
	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return ""
	}

	headerLeft := " Work Queue "
	headerRight := fmt.Sprintf("%d items to go ", m.manager.LenActive())
	headerSpacing := strings.Repeat(" ", m.width-len(headerLeft)-len(headerRight))
	return m.style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftStyle.Render(headerLeft),
			headerSpacing,
			headerRight,
		),
	)
}
