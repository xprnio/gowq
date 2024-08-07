package wq

import (
	"github.com/sqids/sqids-go"
	"github.com/xprnio/work-queue/internal/database"
)

type Work struct {
	Id        string
	Name      string
	Completed bool
}

var ids *sqids.Sqids

func init() {
	var err error
	if ids, err = sqids.New(); err != nil {
		panic(err)
	}
}

func NewWorkItem(name string) database.WorkItem {
	return database.WorkItem{
		Name: name,
	}
}

type Manager struct {
	db    *database.Database
	Queue []database.WorkItem
}

func NewManager(db *database.Database) *Manager {
	return &Manager{db: db}
}

func (m *Manager) Init() (err error) {
	m.Queue, err = m.db.GetAllItems()
	return err
}

func (m *Manager) Get(i int) *database.WorkItem {
	if i >= 0 && i < m.Len() {
		return &m.Queue[i]
	}

	return nil
}

func (m *Manager) AddToTop(w database.WorkItem) error {
	queue := make([]database.WorkItem, len(m.Queue)+1)
	queue[0] = w

	for i := range m.Queue {
		queue[i+1] = m.Queue[i]
	}

	m.Queue = queue
	return m.Flush()
}

func (m *Manager) AddToBottom(w database.WorkItem) error {
	queue := make([]database.WorkItem, m.Len()+1)

	for i := range queue {
		// copy all active items to the new queue
		if i < m.LenActive() {
			queue[i] = m.Queue[i]
			continue
		}

		if i > m.LenActive() {
			// copy all completed items to i+1
			queue[i] = m.Queue[i-1]
			continue
		}

		queue[i] = w
	}

	m.Queue = queue
	return m.Flush()
}

func (m *Manager) Edit(i int, name string) error {
	if i >= 0 && i < m.Len() {
		m.Queue[i].Name = name
		return m.Flush()
	}

	return nil
}

func (m *Manager) Complete(i int) error {
	if i >= 0 && i < m.Len() {
		m.Queue[i].IsCompleted = true
		m.Move(i, m.Len()-1)
		return m.Flush()
	}

	return nil
}

func (m *Manager) Delete(i int) error {
	if i >= 0 && i < m.Len() {
		m.Queue[i].Deleted()
		return m.Flush()
	}

	return nil
}

func (m *Manager) Move(src, dest int) error {
	m.Queue = Move(m.Queue, src, dest)
	return m.Flush()
}

func (m *Manager) Flush() error {
	m.Balance()
	items, err := m.db.SaveAllItems(m.Queue)

	if err != nil {
		return err
	}

	m.Queue = items
	return nil
}

func (m *Manager) Balance() {
	for i := range m.Queue {
		m.Queue[i].Order = i
	}
}

func (m *Manager) Len() int {
	return len(m.Queue)
}

func (m *Manager) LenActive() int {
	var active int

	for _, item := range m.Queue {
		if !item.IsCompleted {
			active++
		}
	}

	return active
}

func (m *Manager) LenCompleted() int {
	var completed int

	for _, item := range m.Queue {
		if item.IsCompleted {
			completed++
		}
	}

	return completed
}
