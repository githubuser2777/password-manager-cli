package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"password-manager-cli/internal/core"
	"testing"
)

func TestModelUpdate(t *testing.T) {
	var m tea.Model = initialModel("/dummy/path")

	if m.(model).state != stateLogin {
		t.Errorf("Expected stateLogin, got %v", m.(model).state)
	}

	// Type master password
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("pass")})

	// Enter key
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m.(model).state != stateDecrypting {
		t.Errorf("Expected stateDecrypting on enter, got %v", m.(model).state)
	}

	// Simulate successful decrypt
	m, _ = m.Update(decryptResultMsg{
		vault: &core.Vault{Entries: make(map[string]core.Entry)},
	})
	if m.(model).state != stateList {
		t.Errorf("Expected stateList after successful decrypt, got %v", m.(model).state)
	}

	mModel := m.(model)
	mModel.vault.Entries["test"] = core.Entry{Username: "u", Password: "p"}
	mModel.updateList()
	m = mModel

	// Select first item and view
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Go to form via a
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
}

func TestModelView(t *testing.T) {
	m := initialModel("/dummy/path")
	_ = m.View() // Login view

	m.state = stateDecrypting
	_ = m.View()

	m.state = stateList
	_ = m.View()

	m.state = stateForm
	m.setupForm("test", "user", "pass", "notes", false)
	_ = m.View()

	m.state = stateConfirmDelete
	m.selectedItem = item{service: "test"}
	_ = m.View()

	m.state = stateView
	_ = m.View()

	m.state = stateMessage
	_ = m.View()

	m.state = stateAudit
	m.auditReport = "Test Report"
	_ = m.View()
}
