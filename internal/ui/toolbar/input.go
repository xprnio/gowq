package toolbar

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/ui/actions"
	"github.com/xprnio/work-queue/internal/wq"
)

func (t *Model) handleInput(msg tea.KeyMsg) tea.Cmd {
	switch mode := t.mode.(type) {
	case ModeNormal:
		switch msg.String() {
		case "a":
			return tea.Sequence(
				ToolbarModeCmd(ModeAdd{}),
			)
		case "e":
			return tea.Sequence(
				actions.ToggleNumbersCmd(true),
				ToolbarModeCmd(ModeEdit{Index: -1}),
			)
		case "up", "k":
			return actions.ScrollWorkListCmd(-1)
		case "down", "j":
			return actions.ScrollWorkListCmd(1)
		case "h":
			t.ShowCompleted = !t.ShowCompleted
			return actions.ToggleCompletedVisibilityCmd(t.ShowCompleted)
		case "m":
			return tea.Sequence(
				actions.ToggleNumbersCmd(true),
				actions.ToggleCompletedVisibilityCmd(false),
				ToolbarModeCmd(ModeMove{}),
			)
		case "c":
			return tea.Sequence(
				actions.ToggleNumbersCmd(true),
				ToolbarModeCmd(ModeComplete{}),
			)
		case "d":
			return tea.Sequence(
				actions.ToggleNumbersCmd(true),
				ToolbarModeCmd(ModeDelete{}),
			)
		case "q":
			return tea.Quit
		}
	case ModeAdd:
		if msg.Type == tea.KeyEscape {
			return ToolbarModeCmd(ModeNormal{})
		}

		if msg.Type == tea.KeyEnter {
			name := strings.TrimSpace(t.input.Value())
			return tea.Sequence(
				actions.WorkAddedCmd(name),
				ToolbarModeCmd(ModeNormal{}),
			)
		}
	case ModeEdit:
		if msg.Type == tea.KeyEscape {
			return tea.Sequence(
				actions.SetWorkListFocusCmd(-1),
				actions.ToggleNumbersCmd(false),
				ToolbarModeCmd(ModeNormal{}),
			)
		}

		if msg.Type == tea.KeyEnter {
			value := t.input.Value()

			if mode.Index < 0 {
				i, err := strconv.Atoi(value)
				if err != nil {
					return tea.Sequence(
						ToolbarModeCmd(ModeNormal{}),
						wq.ErrorCmd(wq.InvalidIndexErr),
					)
				}

				return tea.Sequence(
					actions.SetWorkListFocusCmd(i-1),
					actions.ToggleNumbersCmd(false),
					ToolbarModeCmd(ModeEdit{Index: i - 1}),
				)
			}

			return tea.Sequence(
				actions.SetWorkListFocusCmd(-1),
				actions.ToggleNumbersCmd(false),
				actions.WorkEditedCmd(mode.Index, value),
				ToolbarModeCmd(ModeNormal{}),
			)
		}

	case ModeMove:
		if msg.Type == tea.KeyEscape {
			return tea.Sequence(
				actions.FinishMovingWorkCmd(false),
				actions.ToggleNumbersCmd(false),
				actions.ToggleCompletedVisibilityCmd(t.ShowCompleted),
				ToolbarModeCmd(ModeNormal{}),
			)
		}

		if msg.Type == tea.KeyEnter {
			if mode.Item == nil {
				value := strings.TrimSpace(t.input.Value())
				src, err := strconv.Atoi(value)
				if err != nil {
					return tea.Sequence(
						actions.ToggleNumbersCmd(false),
						actions.ToggleCompletedVisibilityCmd(t.ShowCompleted),
						ToolbarModeCmd(ModeNormal{}),
						wq.ErrorCmd(wq.InvalidIndexErr),
					)
				}

				return actions.ToolbarMoveCmd(src - 1)
			}

			return tea.Sequence(
				actions.FinishMovingWorkCmd(true),
				actions.ToggleCompletedVisibilityCmd(t.ShowCompleted),
				actions.ToggleNumbersCmd(false),
				ToolbarModeCmd(ModeNormal{}),
			)
		}

		switch msg.String() {
		case "up", "k":
			return actions.MoveWorkCmd(actions.MovementDirUp)
		case "down", "j":
			return actions.MoveWorkCmd(actions.MovementDirDown)
		}
	case ModeComplete:
		if msg.Type == tea.KeyEscape {
			return tea.Sequence(
				actions.ToggleNumbersCmd(false),
				ToolbarModeCmd(ModeNormal{}),
			)
		}

		if msg.Type == tea.KeyEnter {
			value := strings.TrimSpace(t.input.Value())
			index, err := strconv.Atoi(value)
			if err != nil {
				return tea.Sequence(
					actions.ToggleNumbersCmd(false),
					ToolbarModeCmd(ModeNormal{}),
					wq.ErrorCmd(wq.InvalidIndexErr),
				)
			}

			return tea.Sequence(
				actions.WorkCompletedCmd(index-1),
				actions.ToggleNumbersCmd(false),
				ToolbarModeCmd(ModeNormal{}),
			)
		}
	case ModeDelete:
		if msg.Type == tea.KeyEscape {
			return tea.Sequence(
				actions.ToggleNumbersCmd(false),
				ToolbarModeCmd(ModeNormal{}),
			)
		}

		if msg.Type == tea.KeyEnter {
			value := strings.TrimSpace(t.input.Value())
			index, err := strconv.Atoi(value)
			if err != nil {
				return tea.Sequence(
					actions.ToggleNumbersCmd(false),
					ToolbarModeCmd(ModeNormal{}),
					wq.ErrorCmd(wq.InvalidIndexErr),
				)
			}

			return tea.Sequence(
				actions.ToggleNumbersCmd(false),
				ToolbarModeCmd(ModeNormal{}),
				actions.WorkDeletedCmd(index-1),
			)
		}
	}

	return nil
}
