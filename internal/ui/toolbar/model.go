package toolbar

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/actions"
	"github.com/xprnio/work-queue/internal/state"
	"github.com/xprnio/work-queue/internal/wq"
)

type Model struct {
	input textinput.Model

	state state.ToolbarState

	style lipgloss.Style
	width int
}

func New() *Model {
	return &Model{
		input: textinput.New(),
		state: state.NewToolbarState(),
		style: baseStyle,
	}
}

func (t *Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
	)
}

func (t *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.style = t.style.Width(t.width)

	case wq.ErrorMsg:
		t.state.Err = msg.Err
		return t, nil

	case actions.ToolbarModeMsg:
		if cmd := t.updateToolbarMode(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		if cmd := t.handleInput(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	var cmd tea.Cmd
	if t.input, cmd = t.input.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
		cmd = nil
	}

	return t, tea.Batch(cmds...)
}

func (t *Model) View() string {
	if t.width == 0 {
		return ""
	}

	return t.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				t.viewMode(),
				t.viewSeparator(),
				t.viewContextKeys(),
			),
			t.viewInput(),
		),
	)
}

func (t *Model) updateToolbarMode(msg actions.ToolbarModeMsg) tea.Cmd {
	t.state.Mode = msg.Mode
	t.state.Err = nil

	switch mode := t.state.Mode.(type) {
	case state.ToolbarModeAdd:
		t.input.Prompt = "name: "
		t.input.SetValue("")
		t.input.Focus()
	case state.ToolbarModeEdit:
		if mode.Index < 0 {
			t.input.Prompt = "edit item number: "
			t.input.SetValue("")
			t.input.Focus()
			break
		}

		t.input.Prompt = "name: "
		t.input.SetValue(mode.Name)
		t.input.Focus()
	case state.ToolbarModeTag:
		if mode.Index < 0 {
			t.input.Prompt = "tag item number: "
			t.input.SetValue("")
			t.input.Focus()
			break
		}

		t.input.Prompt = "tags: "
		t.input.SetValue(mode.Tags)
		t.input.Focus()
	case state.ToolbarModeMove:
		t.input.Prompt = "move item number: "
		t.input.SetValue("")
		t.input.Focus()
	case state.ToolbarModeComplete:
		t.input.Prompt = "complete item number: "
		t.input.SetValue("")
		t.input.Focus()
	case state.ToolbarModeDelete:
		t.input.Prompt = "delete item number: "
		t.input.SetValue("")
		t.input.Focus()
	default:
		t.input.SetValue("")
		t.input.Blur()
	}

	return nil
}
