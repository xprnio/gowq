package toolbar

import "github.com/charmbracelet/lipgloss"

var baseStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#262630")).
	Foreground(lipgloss.Color("#ffffff"))

var modeStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#BCCAEB")).
	Foreground(lipgloss.Color("#2A212E")).
	PaddingLeft(1).PaddingRight(1).
	Bold(true)

var inputStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#20202A"))

var glyphLeftStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#BCCAEB")).
	Background(lipgloss.Color("#404C5B"))

var glyphRightStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#404C5B")).
	Background(lipgloss.Color("#262630"))
