package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xprnio/work-queue/internal/database"
	"github.com/xprnio/work-queue/internal/ui/app"
	"github.com/xprnio/work-queue/internal/wq"
)

func main() {
	dbp, err := databasePath()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDatabase(dbp)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	manager := wq.NewManager(db)
	if err := manager.Init(); err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(
		app.New(manager),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// databasePath returns the first acceptable database path found.
// the following paths are checked, in order:
// - $XDG_CONFIG_HOME/wq.sqlite
// - $XDG_CONFIG_HOME/wq/database.sqlite
// - $HOME/.wq.sqlite
// - $HOME/.wq/database.sqlite
// - $HOME/.config/wq.sqlite
// - $HOME/.config/wq/database.sqlite
// - $(pwd)/wq.sqlite
func databasePath() (string, error) {
	if home, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		paths := []string{
			path.Join(home, "wq.sqlite"),
			path.Join(home, "wq", "database.sqlite"),
		}
		if path, err := tryPaths(paths); path != "" || err != nil {
			return path, err
		}
	}

	if home, ok := os.LookupEnv("HOME"); ok {
		paths := []string{
			path.Join(home, ".wq.sqlite"),
			path.Join(home, "wq", "database.sqlite"),
			path.Join(home, ".config", "wq.sqlite"),
			path.Join(home, ".config", "wq", "database.sqlite"),
		}
		if path, err := tryPaths(paths); path != "" || err != nil {
			return path, err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(wd, "wq.sqlite"), nil
}

func tryPaths(paths []string) (string, error) {
	for _, p := range paths {
		if exists, err := databaseExists(p); exists || err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return p, err
		}
	}

	return "", nil
}

func databaseExists(path string) (bool, error) {
	log.Printf("trying: %s\n", path)
	stat, err := os.Lstat(path)
	if err != nil {
		log.Printf("error: %s\n", err)
		return false, err
	}

	if stat.IsDir() {
		log.Printf("error: is a directory\n")
		return true, fmt.Errorf("path exists but is directory: %s", path)
	}

	return true, nil
}
