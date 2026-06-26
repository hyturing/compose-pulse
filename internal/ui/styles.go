package ui

import "github.com/charmbracelet/lipgloss"

var (
	// State indicator colors.
	colorHealthy   = lipgloss.Color("#22c55e") // green
	colorStarting  = lipgloss.Color("#eab308") // yellow
	colorUnhealthy = lipgloss.Color("#ef4444") // red
	colorPending   = lipgloss.Color("#6b7280") // gray
	colorSelected  = lipgloss.Color("#1e3a5f") // dark blue background for cursor row

	// Tree.
	styleSelected = lipgloss.NewStyle().Background(colorSelected)
	styleName     = lipgloss.NewStyle().Bold(true)
	styleDim      = lipgloss.NewStyle().Faint(true)

	// Log modal.
	styleModal = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6366f1")).
			Padding(1, 2)

	// Search bar inside the modal.
	styleSearchPrompt = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#94a3b8")).
				Bold(true)

	styleSearchInput = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f8fafc"))

	// Status bar at the bottom.
	styleStatusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94a3b8")).
			Faint(true)
)

// spinnerFrames cycles through these for StateStarting.
var spinnerFrames = []string{"◐", "◓", "◑", "◒"}
