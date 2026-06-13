package tui

import "github.com/charmbracelet/lipgloss"

// Styles defines the visual styles for the TUI.
var (
	// Title style for headers and titles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginBottom(1)

	// Subtitle style for secondary headers
	Subtitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06B6D4")).
		Bold(true)

	// SelectedItem style for the currently selected item
	SelectedItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true).
		SetString("> ")

	// NormalItem style for unselected items
	NormalItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		SetString("  ")

	// HotValue style for hot/trending values
	HotValue = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Italic(true)

	// Error style for error messages
	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Bold(true)

	// Info style for informational messages
	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6"))

	// Help style for help text
	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)

	// Border style for boxes
	Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1)

	// Status Bar style
	StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Background(lipgloss.Color("#1F2937")).
		Padding(0, 1)
)
