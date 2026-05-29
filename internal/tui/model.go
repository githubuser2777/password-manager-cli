package tui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"password-manager-cli/internal/core"
	"password-manager-cli/internal/storage"
)

// UI Styles
var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("205"))
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

type state int

const (
	stateLogin state = iota
	stateList
	stateView
	stateMessage
)

type item struct {
	service  string
	username string
	password string
}

func (i item) Title() string       { return i.service }
func (i item) Description() string { return i.username }
func (i item) FilterValue() string { return i.service }

type model struct {
	state       state
	vaultPath   string
	vault       *core.Vault
	
	passwordInput textinput.Model
	servicesList  list.Model
	
	selectedItem item
	msg          string
	isError      bool
}

func initialModel(vaultPath string) model {
	ti := textinput.New()
	ti.Placeholder = "Master Password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.Focus()

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Password Vault"
	l.Styles.Title = titleStyle

	return model{
		state:         stateLogin,
		vaultPath:     vaultPath,
		passwordInput: ti,
		servicesList:  l,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.state == stateView || m.state == stateMessage {
				if m.vault != nil {
					m.state = stateList
					return m, nil
				}
				m.state = stateLogin
				return m, nil
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.servicesList.SetSize(msg.Width-h, msg.Height-v)
	}

	switch m.state {
	case stateLogin:
		m.passwordInput, cmd = m.passwordInput.Update(msg)
		cmds = append(cmds, cmd)

		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			pw := m.passwordInput.Value()
			vault, err := storage.LoadVault(m.vaultPath, pw)
			if err != nil {
				m.msg = "Invalid Master Password or Vault not found.\nError: " + err.Error()
				m.isError = true
				m.state = stateMessage
			} else {
				m.vault = vault
				// Populate list
				var items []list.Item
				for s, e := range vault.Entries {
					items = append(items, item{service: s, username: e.Username, password: e.Password})
				}
				m.servicesList.SetItems(items)
				m.state = stateList
			}
		}
	case stateList:
		m.servicesList, cmd = m.servicesList.Update(msg)
		cmds = append(cmds, cmd)

		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			if i, ok := m.servicesList.SelectedItem().(item); ok {
				m.selectedItem = i
				m.state = stateView
			}
		}
	case stateView:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "c":
				if err := clipboard.WriteAll(m.selectedItem.password); err != nil {
					m.msg = "Failed to copy password: " + err.Error()
					m.isError = true
				} else {
					m.msg = "Password copied to clipboard!"
					m.isError = false
				}
				m.state = stateMessage
			}
		}
	case stateMessage:
		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			if m.vault != nil {
				m.state = stateList
			} else {
				m.state = stateLogin
				m.passwordInput.SetValue("")
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	switch m.state {
	case stateLogin:
		s = fmt.Sprintf(
			"%s\n\n%s\n\n(esc to quit)",
			titleStyle.Render("Unlock Vault"),
			m.passwordInput.View(),
		)
	case stateList:
		s = m.servicesList.View()
	case stateView:
		s = fmt.Sprintf(
			"%s\n\nService: %s\nUsername: %s\nPassword: %s\n\n[c] Copy Password  [esc] Back",
			titleStyle.Render("View Credential"),
			m.selectedItem.service,
			m.selectedItem.username,
			m.selectedItem.password,
		)
	case stateMessage:
		style := infoStyle
		if m.isError {
			style = errorStyle
		}
		s = fmt.Sprintf("\n%s\n\nPress Enter to continue.", style.Render(m.msg))
	}
	return appStyle.Render(s)
}
