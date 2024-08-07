package toolbar

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/wq"
)

type Model struct {
	width int
	mode  ToolbarMode
	style lipgloss.Style
	err   error

	ShowCompleted bool

	input textinput.Model
}

func New() *Model {
	t := &Model{}
	t.mode = ModeNormal{}
	t.style = baseStyle

	t.input = textinput.New()
	return t
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
		t.err = msg.Err
		return t, nil

	case ToolbarModeMsg:
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

func (t *Model) updateToolbarMode(msg ToolbarModeMsg) tea.Cmd {
	t.mode = msg.Mode
	t.err = nil

	switch mode := t.mode.(type) {
	case ModeAdd:
		t.input.Prompt = "name: "
		t.input.SetValue("")
		t.input.Focus()
	case ModeEdit:
		if mode.Index < 0 {
			t.input.Prompt = "edit item number: "
			t.input.SetValue("")
			t.input.Focus()
			break
		}

		t.input.Prompt = "name: "
		t.input.SetValue(mode.Name)
		t.input.Focus()
	case ModeMove:
		t.input.Prompt = "move item number: "
		t.input.SetValue("")
		t.input.Focus()
	case ModeComplete:
		t.input.Prompt = "complete item number: "
		t.input.SetValue("")
		t.input.Focus()
	case ModeDelete:
		t.input.Prompt = "delete item number: "
		t.input.SetValue("")
		t.input.Focus()
	default:
		t.input.SetValue("")
		t.input.Blur()
	}

	return nil
}
