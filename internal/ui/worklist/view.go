package worklist

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/database"
)

func (l *Model) viewContent() string {
	var rows []string
	for i, item := range l.queue {
		rows = append(rows, l.viewItem(item, i))
	}

	return l.style.Render(
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)
}

func (l *Model) viewItem(item database.WorkItem, index int) string {
	style := l.itemStyle.UnsetWidth().Faint(item.IsCompleted)

	if focused := l.state.Focused; focused != nil {
		style = style.Faint(index != focused.Index)
	}

	if moving := l.state.Moving; moving != nil {
		style = style.Faint(index != moving.Active)
	}

	return l.itemStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			style.Render(l.viewItemPrefix(item, index)),
			style.Strikethrough(item.IsCompleted).Render(item.Name),
		),
	)
}

func (l *Model) viewItemPrefix(item database.WorkItem, index int) string {
	if l.state.ShowNumbers {
		return fmt.Sprintf(" %d ", index+1)
	}

	if item.IsCompleted {
		return " 󰱒 "
	}

	return " 󰄱 "
}
