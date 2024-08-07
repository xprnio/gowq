package worklist

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/database"
	"github.com/xprnio/work-queue/internal/ui/state"
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
	if l.FocusedItem >= 0 {
		style = style.Faint(index != l.FocusedItem)
	}

	if s, isMoving := l.state.(state.WorkListMovingState); isMoving {
		style = style.Faint(index != s.Active)
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
	if l.ShowNumbers {
		return fmt.Sprintf(" %d ", index+1)
	}

	if item.IsCompleted {
		return " 󰱒 "
	}

	return " 󰄱 "
}
