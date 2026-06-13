package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackwener/aiview/tui"
	"github.com/spf13/cobra"
)

// NewTUICmd creates the TUI command
func NewTUICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive TUI",
		Long: `Launch the interactive Terminal User Interface for AIView.

The TUI provides a keyboard-driven interface for:
  • Browsing hot/trending content across platforms
  • Navigating with arrow keys and vim-style bindings
  • Viewing detailed information about items

Examples:
  aiview tui`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("TUI error: %w", err)
			}
			return nil
		},
	}

	return cmd
}
