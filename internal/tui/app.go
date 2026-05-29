package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// StartApp initializes and runs the Bubble Tea TUI.
func StartApp(vaultPath string) error {
	m := initialModel(vaultPath)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running tui: %w", err)
	}
	return nil
}
