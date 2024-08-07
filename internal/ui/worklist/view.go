package worklist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/database"
)

func (l *Model) viewContent() string {
	if l.width <= 0 {
		return ""
	}

	var rows []string
	for i, item := range l.queue {
		rows = append(rows, l.viewItem(item, i))
	}

	return l.style.Render(
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)
}

func (l *Model) viewItem(item database.WorkItem, index int) string {
	return l.itemStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			l.viewItemContent(item, index),
		),
	)
}

func (m *Model) viewItemContent(item database.WorkItem, index int) string {
	style := m.viewItemStyle(item, index)

	prefix := m.viewItemPrefix(item, index)
	tags := m.viewItemTags(item)
	width := lipgloss.Width(prefix) + lipgloss.Width(item.Name) + lipgloss.Width(tags)
	spacing := strings.Repeat(" ", max(0, m.width-width))
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		style.Render(prefix),
		style.Render(item.Name),
		style.Render(spacing),
		style.Faint(true).Render(tags),
	)
}

func (l *Model) viewItemTags(item database.WorkItem) string {
	if len(item.Tags) == 0 {
		return fmt.Sprintf(" [-] ")
	}

  tags := ""
  for i, tag := range item.Tags {
    if i > 0 {
      tags += ", "
    }

    tags += tag
  }
	return fmt.Sprintf(" [%s] ", tags)
}

func (m *Model) viewItemPrefix(item database.WorkItem, index int) string {
	if m.state.ShowNumbers {
		return fmt.Sprintf(" %d ", index+1)
	}

	if item.IsCompleted {
		return " 󰱒 "
	}

	return " 󰄱 "
}

func (m *Model) viewItemStyle(item database.WorkItem, index int) lipgloss.Style {
	style := m.itemStyle.UnsetWidth().Faint(item.IsCompleted)

	if focused := m.state.Focused; focused != nil {
		style = style.Faint(index != focused.Index)
	}

	if moving := m.state.Moving; moving != nil {
		style = style.Faint(index != moving.Active)
	}

	return style
}
