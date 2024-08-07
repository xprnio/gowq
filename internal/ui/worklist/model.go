package worklist

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/xprnio/work-queue/internal/database"
	"github.com/xprnio/work-queue/internal/wq"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/ui/actions"
	"github.com/xprnio/work-queue/internal/ui/state"
)

type Model struct {
	Manager *wq.Manager

	Viewport viewport.Model

	FocusedItem int

	state state.WorkListState

	// item numbers
	ShowNumbers bool

	// completion visibility
	ShowCompleted bool

	width, height int

	queue     []database.WorkItem
	style     lipgloss.Style
	itemStyle lipgloss.Style
}

func New(man *wq.Manager) *Model {
	l := &Model{}
	l.Manager = man
	l.FocusedItem = -1

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
			l.Viewport = viewport.New(msg.Width, msg.Height-3)
			l.Viewport.SetContent(l.viewContent())
		}

		l.width = msg.Width
		l.height = msg.Height - 3
		l.style = baseStyle.Width(l.width).Height(l.height)
		l.itemStyle = itemStyle.Width(l.width)

		l.Viewport.Width = l.width
		l.Viewport.Height = l.height
		return l, actions.RefreshWorkListCmd()
	case actions.RefreshWorkListMsg:
		l.Viewport.SetContent(l.viewContent())

		// refresh viewport offset
		maxOffset := l.Viewport.TotalLineCount() - l.height
		switch {
		case l.Viewport.YOffset < 0:
			l.Viewport.YOffset = 0
		case l.Viewport.YOffset > maxOffset:
			l.Viewport.YOffset = maxOffset
		}

		return l, nil
	case actions.ScrollWorkListMsg:
		l.Viewport.YOffset += msg.Direction
		return l, actions.RefreshWorkListCmd()
	case actions.ToolbarModeMsg:
		switch mode := msg.Mode.(type) {
		case state.ToolbarModeEdit:
			if mode.Index >= 0 && mode.Name == "" {
				item := l.Manager.Get(mode.Index)
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
		item := l.Manager.Get(msg.Index)
		if item == nil {
			return l, tea.Sequence(
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
				actions.RefreshWorkListCmd(),
				wq.ErrorCmd(wq.InvalidIndexErr),
			)
		}

		if err := l.Manager.Edit(msg.Index, msg.Name); err != nil {
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
			l.FocusedItem = msg.Index
			break
		}

		l.FocusedItem = -1
	case actions.ToggleCompletedVisibilityMsg:
		l.ShowCompleted = msg.Visible
		l.queue = l.RefreshQueue()
	case actions.ToggleNumbersActionMsg:
		l.ShowNumbers = msg.Visible
		l.queue = l.RefreshQueue()
	case actions.ToolbarMoveMsg:
		item := l.Manager.Get(msg.Source)
		if item == nil {
			return l, tea.Sequence(
				actions.ToggleNumbersCmd(false),
				actions.ToolbarModeCmd(state.ToolbarModeNormal{}),
				actions.RefreshWorkListCmd(),
				wq.ErrorCmd(wq.InvalidIndexErr),
			)
		}

		l.state = state.WorkListMovingState{
			Source: msg.Source,
			Active: msg.Source,
		}
		return l, tea.Sequence(
			actions.ToggleNumbersCmd(false),
			actions.RefreshWorkListCmd(),
			actions.ToolbarModeCmd(state.ToolbarModeMove{Item: item}),
		)
	case actions.MoveWorkMsg:
		switch msg.Direction {
		case actions.MovementDirUp:
			if s, ok := l.state.(state.WorkListMovingState); ok {
				l.queue = wq.Move(l.queue, s.Active, s.Active-1)

				s.Active = max(0, s.Active-1)
				l.state = s
			}
		case actions.MovementDirDown:
			if s, ok := l.state.(state.WorkListMovingState); ok {
				l.queue = wq.Move(l.queue, s.Active, s.Active+1)

				s.Active = min(s.Active+1, l.Len()-1)
				l.state = s
			}
		}
	case actions.FinishMovingWorkMsg:
		if s, ok := l.state.(state.WorkListMovingState); ok {
			if msg.Commit {
				l.Manager.Move(s.Source, s.Active)
			}

			l.state = nil
		}

		l.queue = l.RefreshQueue()
	case actions.WorkAddedMsg:
		// TODO: Implement adding to top of queue
		l.Manager.AddToBottom(msg.Work)
		l.queue = l.RefreshQueue()
	case actions.WorkCompletedMsg:
		l.Manager.Complete(msg.Index)
		l.queue = l.RefreshQueue()
	case actions.WorkDeletedMsg:
		if err := l.Manager.Delete(msg.Index); err != nil {
			return l, wq.ErrorCmd(err)
		}
		l.queue = l.RefreshQueue()
	}

	return l, actions.RefreshWorkListCmd()
}

func (l *Model) View() string {
	return l.Viewport.View()
}

func (l *Model) RefreshQueue() []database.WorkItem {
	queue := make([]database.WorkItem, 0)
	for _, item := range l.Manager.Queue {
		if !l.ShowCompleted && item.IsCompleted {
			continue
		}

		queue = append(queue, item)
	}

	return queue
}

func (l *Model) Len() int {
	return len(l.queue)
}
