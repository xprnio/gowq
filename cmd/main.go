package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/database"
	"github.com/xprnio/work-queue/internal/ui/application"
	"github.com/xprnio/work-queue/internal/wq"
)

func main() {
  db, err := database.NewDatabase("/home/ragnar/.config/wq/database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	manager := wq.NewManager(db)
	if err := manager.Init(); err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(
		application.New(manager),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
