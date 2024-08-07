package toolbar

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/ui/state"
	"github.com/xprnio/work-queue/internal/wq"
)

func (t *Model) viewMode() string {
	switch t.state.Mode.(type) {
	case state.ToolbarModeNormal:
		return modeStyle.Render("NORMAL")
	case state.ToolbarModeAdd:
		return modeStyle.Render("ADD")
	case state.ToolbarModeEdit:
		return modeStyle.Render("EDIT")
	case state.ToolbarModeComplete:
		return modeStyle.Render("COMPLETE")
	case state.ToolbarModeMove:
		return modeStyle.Render("MOVE")
	case state.ToolbarModeDelete:
		return modeStyle.Render("DELETE")
	default:
		return ""
	}
}

func (t *Model) viewSeparator() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		glyphLeftStyle.Render(wq.GlyphSeparator),
		glyphRightStyle.Render(wq.GlyphSeparator),
	)
}

func (t *Model) viewContextKeys() string {
	keys := make([]string, 0)

	switch mode := t.state.Mode.(type) {
	case state.ToolbarModeNormal:
		keys = append(keys,
			"[a]dd",
			"[e]dit",
			"[c]omplete",
			"[d]elete",
		)

		if !t.state.ShowCompleted {
			keys = append(keys, "s[h]ow completed")
		} else {
			keys = append(keys, "[h]ide completed")
		}

		keys = append(keys, "[q]uit")
	case
	 state.ToolbarModeAdd,
	 state.ToolbarModeComplete,
	 state.ToolbarModeDelete:
		keys = append(keys, "[esc] cancel")
	case state.ToolbarModeEdit:
		if mode.Index >= 0 {
			keys = append(keys, "[enter] confirm")
		}
		keys = append(keys, "[esc] cancel")
	case state.ToolbarModeMove:
		if mode.Item != nil {
			keys = append(keys, "[enter] confirm")
		}

		keys = append(keys, "[esc] cancel")
	}

	return t.viewKeys(keys)
}

func (t *Model) viewInput() string {
	if t.state.Err != nil {
		return t.state.Err.Error()
	}

	switch mode := t.state.Mode.(type) {
	case
	 state.ToolbarModeAdd,
	 state.ToolbarModeComplete,
	 state.ToolbarModeDelete:
		return t.input.View()
	case state.ToolbarModeEdit:
		if mode.Index < 0 {
			return t.input.View()
		}

		if mode.Name != "" {
			return t.input.View()
		}

		return ""
	case state.ToolbarModeMove:
		if mode.Item == nil {
			return t.input.View()
		}

		return fmt.Sprintf("moving: %s", mode.Item.Name)
	default:
		return ""
	}
}

func (t *Model) viewKeys(keys []string) string {
	for i, key := range keys {
		keys[i] = t.viewKey(key)
	}
	s := lipgloss.NewStyle().Background(lipgloss.Color("#262630"))
	return s.Render(
		lipgloss.JoinHorizontal(lipgloss.Top, keys...),
	)
}

func (t *Model) viewKey(key string) string {
	// do not apply special styling in 'normal' mode
	// if t.mode == ToolbarModeNormal {
	// 	return fmt.Sprintf(" %s ", key)
	// }

	// apply faint style if key is not current mode
	// style := lipgloss.NewStyle().Faint(t.mode != mode)
	style := lipgloss.NewStyle()
	return style.Render(
		fmt.Sprintf(" %s ", key),
	)
}
