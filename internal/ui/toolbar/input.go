package toolbar

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/ui/actions"
	"github.com/xprnio/work-queue/internal/ui/state"
	"github.com/xprnio/work-queue/internal/wq"
)

func (t *Model) handleInput(msg tea.KeyMsg) tea.Cmd {
	switch mode := t.state.Mode.(type) {
	case state.ToolbarModeNormal:
		switch msg.String() {
		case "a":
			return tea.Sequence(
				actions.ToolbarModeCmd(state.ToolbarModeAdd{}),
			)
		case "e":
			return tea.Sequence(
				actions.ToggleNumbersCmd(true),
				actions.ToolbarModeCmd(state.ToolbarModeEdit{Index: -1}),
			)
		case "up", "k":
			return actions.ScrollWorkListCmd(-1)
		case "down", "j":
			return actions.ScrollWorkListCmd(1)
		case "h":
			t.state.ShowCompleted = !t.state.ShowCompleted
			return actions.ToggleCompletedVisibilityCmd(t.state.ShowCompleted)
		case "m":
			return tea.Sequence(
				actions.ToggleNumbersCmd(true),
				actions.ToggleCompletedVisibilityCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeMove{}),
			)
		case "c":
			return tea.Sequence(
				actions.ToggleNumbersCmd(true),
				actions.ToolbarModeCmd(state.ToolbarModeComplete{}),
			)
		case "d":
			return tea.Sequence(
				actions.ToggleNumbersCmd(true),
				actions.ToolbarModeCmd(state.ToolbarModeDelete{}),
			)
		case "q":
			return tea.Quit
		}
	case state.ToolbarModeAdd:
		if msg.Type == tea.KeyEscape {
			return actions.ToolbarModeCmd(state.ToolbarModeNormal{})
		}

		if msg.Type == tea.KeyEnter {
			name := strings.TrimSpace(t.input.Value())
			return tea.Sequence(
				actions.WorkAddedCmd(name),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
			)
		}
	case state.ToolbarModeEdit:
		if msg.Type == tea.KeyEscape {
			return tea.Sequence(
				actions.SetWorkListFocusCmd(-1),
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
			)
		}

		if msg.Type == tea.KeyEnter {
			value := t.input.Value()

			if mode.Index < 0 {
				i, err := strconv.Atoi(value)
				if err != nil {
					return tea.Sequence(
						actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
						wq.ErrorCmd(wq.InvalidIndexErr),
					)
				}

				return tea.Sequence(
					actions.SetWorkListFocusCmd(i-1),
					actions.ToggleNumbersCmd(false),
					actions.ToolbarModeCmd(state.ToolbarModeEdit{Index: i - 1}),
				)
			}

			return tea.Sequence(
				actions.SetWorkListFocusCmd(-1),
				actions.ToggleNumbersCmd(false),
				actions.WorkEditedCmd(mode.Index, value),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
			)
		}

	case state.ToolbarModeMove:
		if msg.Type == tea.KeyEscape {
			return tea.Sequence(
				actions.FinishMovingWorkCmd(false),
				actions.ToggleNumbersCmd(false),
				actions.ToggleCompletedVisibilityCmd(t.state.ShowCompleted),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
			)
		}

		if msg.Type == tea.KeyEnter {
			if mode.Item == nil {
				value := strings.TrimSpace(t.input.Value())
				src, err := strconv.Atoi(value)
				if err != nil {
					return tea.Sequence(
						actions.ToggleNumbersCmd(false),
						actions.ToggleCompletedVisibilityCmd(t.state.ShowCompleted),
						actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
						wq.ErrorCmd(wq.InvalidIndexErr),
					)
				}

				return actions.ToolbarMoveCmd(src - 1)
			}

			return tea.Sequence(
				actions.FinishMovingWorkCmd(true),
				actions.ToggleCompletedVisibilityCmd(t.state.ShowCompleted),
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
			)
		}

		switch msg.String() {
		case "up", "k":
			return actions.MoveWorkCmd(actions.MovementDirUp)
		case "down", "j":
			return actions.MoveWorkCmd(actions.MovementDirDown)
		}
	case state.ToolbarModeComplete:
		if msg.Type == tea.KeyEscape {
			return tea.Sequence(
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
			)
		}

		if msg.Type == tea.KeyEnter {
			value := strings.TrimSpace(t.input.Value())
			index, err := strconv.Atoi(value)
			if err != nil {
				return tea.Sequence(
					actions.ToggleNumbersCmd(false),
					actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
					wq.ErrorCmd(wq.InvalidIndexErr),
				)
			}

			return tea.Sequence(
				actions.WorkCompletedCmd(index-1),
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
			)
		}
	case state.ToolbarModeDelete:
		if msg.Type == tea.KeyEscape {
			return tea.Sequence(
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
			)
		}

		if msg.Type == tea.KeyEnter {
			value := strings.TrimSpace(t.input.Value())
			index, err := strconv.Atoi(value)
			if err != nil {
				return tea.Sequence(
					actions.ToggleNumbersCmd(false),
					actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
					wq.ErrorCmd(wq.InvalidIndexErr),
				)
			}

			return tea.Sequence(
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
				actions.WorkDeletedCmd(index-1),
			)
		}
	}

	return nil
}
