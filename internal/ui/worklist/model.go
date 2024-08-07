package worklist

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/database"
	"github.com/xprnio/work-queue/internal/wq"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/actions"
	"github.com/xprnio/work-queue/internal/state"
)

type Model struct {
	manager *wq.Manager
	viewport viewport.Model

	state state.WorkListState
	queue []database.WorkItem

	style         lipgloss.Style
	itemStyle     lipgloss.Style
	width, height int
}

func New(man *wq.Manager) *Model {
	l := &Model{}
	l.manager = man

	l.style = baseStyle
	l.itemStyle = itemStyle
	return l
}

func (l *Model) Init() tea.Cmd {
	return func() tea.Msg {
		l.queue = l.RefreshQueue()
		return nil
	}
}

func (l *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: get rid of magic number 3

		if l.height == 0 {
			l.viewport = viewport.New(msg.Width, msg.Height-3)
			l.viewport.SetContent(l.viewContent())
		}

		l.width = msg.Width
		l.height = msg.Height - 3
		l.style = baseStyle.Width(l.width).Height(l.height)
		l.itemStyle = itemStyle.Width(l.width)

		l.viewport.Width = l.width
		l.viewport.Height = l.height
		return l, actions.RefreshWorkListCmd()
	case actions.RefreshWorkListMsg:
		l.viewport.SetContent(l.viewContent())

		// refresh viewport offset
		maxOffset := l.viewport.TotalLineCount() - l.height
		switch {
		case l.viewport.YOffset < 0:
			l.viewport.YOffset = 0
		case l.viewport.YOffset > maxOffset:
			l.viewport.YOffset = maxOffset
		}

		return l, nil
	case actions.ScrollWorkListMsg:
		l.viewport.YOffset += msg.Direction
		return l, actions.RefreshWorkListCmd()
	case actions.ToolbarModeMsg:
		switch mode := msg.Mode.(type) {
		case state.ToolbarModeEdit:
			if mode.Index >= 0 && mode.Name == "" {
				item := l.manager.Get(mode.Index)
				if item == nil {
					return l, tea.Sequence(
						actions.ToggleNumbersCmd(false),
						actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
						actions.RefreshWorkListCmd(),
						wq.ErrorCmd(wq.InvalidIndexErr),
					)
				}

				return l, actions.ToolbarModeCmd(state.ToolbarModeEdit{
					Index: mode.Index,
					Name:  item.Name,
				})
			}
		}
	case actions.WorkEditedMsg:
		item := l.manager.Get(msg.Index)
		if item == nil {
			return l, tea.Sequence(
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
				actions.RefreshWorkListCmd(),
				wq.ErrorCmd(wq.InvalidIndexErr),
			)
		}

		if err := l.manager.Edit(msg.Index, msg.Name); err != nil {
			return l, tea.Sequence(
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
				actions.RefreshWorkListCmd(),
				wq.ErrorCmd(err),
			)
		}

		l.queue = l.RefreshQueue()
	case actions.SetWorkListFocusMsg:
		if msg.Index >= 0 && msg.Index < l.Len() {
			l.state.Focused = state.WorkListFocusedState(msg.Index)
			break
		}

		l.state.Focused = nil
	case actions.ToggleCompletedVisibilityMsg:
		l.state.ShowCompleted = msg.Visible
		l.queue = l.RefreshQueue()
	case actions.ToggleNumbersActionMsg:
		l.state.ShowNumbers = msg.Visible
		l.queue = l.RefreshQueue()
	case actions.ToolbarMoveMsg:
		item := l.manager.Get(msg.Source)
		if item == nil {
			return l, tea.Sequence(
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
				actions.RefreshWorkListCmd(),
				wq.ErrorCmd(wq.InvalidIndexErr),
			)
		}

		l.state.Moving = state.WorkListMovingState(msg.Source)
		return l, tea.Sequence(
			actions.ToggleNumbersCmd(false),
			actions.RefreshWorkListCmd(),
			actions.ToolbarModeCmd(state.ToolbarModeMove{Item: item}),
		)
	case actions.MoveWorkMsg:
		switch msg.Direction {
		case actions.MovementDirUp:
			if s := l.state.Moving; s != nil {
				l.queue = wq.Move(l.queue, s.Active, s.Active-1)
				s.Active = max(0, s.Active-1)
			}
		case actions.MovementDirDown:
			if s := l.state.Moving; s != nil {
				l.queue = wq.Move(l.queue, s.Active, s.Active+1)
				s.Active = min(s.Active+1, l.Len()-1)
			}
		}
	case actions.FinishMovingWorkMsg:
		if s := l.state.Moving; s != nil {
			if msg.Commit {
				l.manager.Move(s.Source, s.Active)
			}

			// disable moving
			l.state.Moving = nil
		}

		l.queue = l.RefreshQueue()
	case actions.WorkAddedMsg:
		// TODO: Implement adding to top of queue
		l.manager.AddToBottom(msg.Work)
		l.queue = l.RefreshQueue()
	case actions.WorkCompletedMsg:
		l.manager.Complete(msg.Index)
		l.queue = l.RefreshQueue()
	case actions.WorkDeletedMsg:
		if err := l.manager.Delete(msg.Index); err != nil {
			return l, wq.ErrorCmd(err)
		}
		l.queue = l.RefreshQueue()
	}

	return l, actions.RefreshWorkListCmd()
}

func (l *Model) View() string {
	return l.viewport.View()
}

func (l *Model) RefreshQueue() []database.WorkItem {
	queue := make([]database.WorkItem, 0)
	for _, item := range l.manager.Queue {
		if !l.state.ShowCompleted && item.IsCompleted {
			continue
		}

		queue = append(queue, item)
	}

	return queue
}

func (l *Model) Len() int {
	return len(l.queue)
}
