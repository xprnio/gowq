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
	s := l.itemStyle.UnsetWidth().Faint(item.IsCompleted)
	if l.FocusedItem >= 0 {
		s = s.Faint(index != l.FocusedItem)
	}
	if l.IsMoving {
		s = s.Faint(index != l.MoveActive)
	}

	return l.itemStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			s.Render(l.viewItemPrefix(item, index)),
			s.Strikethrough(item.IsCompleted).Render(item.Name),
		),
	)
}

func (l *Model) viewItemPrefix(item database.WorkItem, index int) string {
	if l.ShowNumbers {
		return fmt.Sprintf(" %d ", index+1)
	}

	if item.IsCompleted {
		return " 󰱒 "
	}

	return " 󰄱 "
}

