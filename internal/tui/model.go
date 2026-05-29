package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"password-manager-cli/internal/core"
	"password-manager-cli/internal/storage"
)

// Custom Keys for Help Menu
type listKeyMap struct {
	add         key.Binding
	edit        key.Binding
	delete      key.Binding
	auditLocal  key.Binding
	auditOnline key.Binding
}

var customKeys = listKeyMap{
	add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	auditLocal: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "audit (local)"),
	),
	auditOnline: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "audit (online)"),
	),
}

// UI Styles
var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("205"))
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	focusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type state int

const (
	stateLogin state = iota
	stateList
	stateView
	stateMessage
	stateForm
	stateConfirmDelete
	stateAudit
)

type item struct {
	service   string
	username  string
	password  string
	notes     string
	createdAt string
	updatedAt string
}

func (i item) Title() string       { return i.service }
func (i item) Description() string { return i.username }
func (i item) FilterValue() string { return i.service }

type model struct {
	state       state
	vaultPath   string
	masterPw    string
	vault       *core.Vault
	
	passwordInput textinput.Model
	servicesList  list.Model
	
	formInputs []textinput.Model
	focusIndex int
	isEditing  bool
	
	selectedItem item
	msg          string
	isError      bool
	auditReport  string
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
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{customKeys.add, customKeys.edit, customKeys.delete, customKeys.auditLocal, customKeys.auditOnline}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{customKeys.add, customKeys.edit, customKeys.delete, customKeys.auditLocal, customKeys.auditOnline}
	}

	return model{
		state:         stateLogin,
		vaultPath:     vaultPath,
		passwordInput: ti,
		servicesList:  l,
	}
}

func (m *model) setupForm(service, username, password, notes string, isEdit bool) {
	m.formInputs = make([]textinput.Model, 4)
	
	var t textinput.Model
	for i := range m.formInputs {
		t = textinput.New()
		t.Cursor.Style = focusStyle
		t.CharLimit = 128

		switch i {
		case 0:
			t.Placeholder = "Service Name (e.g. github.com)"
			t.SetValue(service)
			t.Focus()
			t.PromptStyle = focusStyle
			t.TextStyle = focusStyle
			if isEdit {
				t.Blur() // Disable editing service name
			}
		case 1:
			t.Placeholder = "Username"
			t.SetValue(username)
		case 2:
			t.Placeholder = "Password"
			t.SetValue(password)
		case 3:
			t.Placeholder = "Notes (optional)"
			t.SetValue(notes)
		}

		m.formInputs[i] = t
	}
	m.focusIndex = 0
	if isEdit {
		m.focusIndex = 1
		m.formInputs[0].PromptStyle = blurStyle
		m.formInputs[0].TextStyle = blurStyle
		m.formInputs[1].Focus()
		m.formInputs[1].PromptStyle = focusStyle
		m.formInputs[1].TextStyle = focusStyle
	}
	m.isEditing = isEdit
}

func (m *model) updateList() {
	var items []list.Item
	for s, e := range m.vault.Entries {
		items = append(items, item{
			service:   s,
			username:  e.Username,
			password:  e.Password,
			notes:     e.Notes,
			createdAt: e.CreatedAt,
			updatedAt: e.UpdatedAt,
		})
	}
	m.servicesList.SetItems(items)
}

func (m *model) runAudit(online bool) {
	var report strings.Builder
	report.WriteString(titleStyle.Render(fmt.Sprintf("Password Audit Report (Online: %v)", online)))
	report.WriteString("\n\n")

	weakCount := 0
	reusedCount := 0
	pwnedCount := 0

	// Check for reuse
	pwMap := make(map[string][]string)
	for service, entry := range m.vault.Entries {
		pwMap[entry.Password] = append(pwMap[entry.Password], service)
	}

	for _, services := range pwMap {
		if len(services) > 1 {
			reusedCount++
			report.WriteString(fmt.Sprintf("[!] REUSED: The password for %s is used across %d services.\n", services[0], len(services)))
		}
	}

	for service, entry := range m.vault.Entries {
		// Basic strength check
		if len(entry.Password) < 8 {
			weakCount++
			report.WriteString(fmt.Sprintf("[!] WEAK: Password for '%s' is too short (under 8 chars).\n", service))
		} else if !strings.ContainsAny(entry.Password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") ||
			!strings.ContainsAny(entry.Password, "0123456789") {
			weakCount++
			report.WriteString(fmt.Sprintf("[!] WEAK: Password for '%s' lacks numbers or uppercase letters.\n", service))
		}

		// Online HIBP check
		if online {
			pwned, err := core.CheckPwned(entry.Password)
			if err != nil {
				report.WriteString(fmt.Sprintf("[-] HIBP API Check failed for '%s': %v\n", service, err))
			} else if pwned {
				pwnedCount++
				report.WriteString(fmt.Sprintf("[!!!] PWNED: Password for '%s' has been found in data breaches!\n", service))
			}
		}
	}

	report.WriteString("\n----------------------------------------\n")
	report.WriteString(fmt.Sprintf("Audit Complete! Total passwords checked: %d\n", len(m.vault.Entries)))
	report.WriteString(fmt.Sprintf("Weak: %d | Reused: %d | Pwned: %d\n", weakCount, reusedCount, pwnedCount))
	if !online {
		report.WriteString("(Run with Shift+r (R) to check online for data breaches)\n")
	}
	report.WriteString("\n[esc] Back to list")

	m.auditReport = report.String()
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
			if m.state == stateView || m.state == stateMessage || m.state == stateForm || m.state == stateConfirmDelete || m.state == stateAudit {
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
				m.masterPw = pw
				m.vault = vault
				m.updateList()
				m.state = stateList
			}
		}
	case stateList:
		m.servicesList, cmd = m.servicesList.Update(msg)
		cmds = append(cmds, cmd)

		if msg, ok := msg.(tea.KeyMsg); ok {
			// Only allow shortcuts when not filtering
			if m.servicesList.FilterState() != list.Filtering {
				switch msg.String() {
				case "enter":
					if i, ok := m.servicesList.SelectedItem().(item); ok {
						m.selectedItem = i
						m.state = stateView
					}
				case "a":
					m.setupForm("", "", "", "", false)
					m.state = stateForm
				case "e":
					if i, ok := m.servicesList.SelectedItem().(item); ok {
						m.selectedItem = i
						m.setupForm(i.service, i.username, i.password, i.notes, true)
						m.state = stateForm
					}
				case "d":
					if i, ok := m.servicesList.SelectedItem().(item); ok {
						m.selectedItem = i
						m.state = stateConfirmDelete
					}
				case "r":
					m.runAudit(false)
					m.state = stateAudit
				case "R":
					m.runAudit(true)
					m.state = stateAudit
				}
			}
		}
	case stateForm:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()
				
				if s == "enter" && m.focusIndex == len(m.formInputs)-1 {
					// Save
					service := m.formInputs[0].Value()
					username := m.formInputs[1].Value()
					password := m.formInputs[2].Value()
					notes := m.formInputs[3].Value()

					if service == "" {
						break
					}
					
					entry, exists := m.vault.Entries[service]
					if m.isEditing {
						entry.Username = username
						entry.Password = password
						entry.Notes = notes
						entry.UpdatedAt = time.Now().Format(time.RFC3339)
					} else {
						if exists {
							m.msg = "Service already exists!"
							m.isError = true
							m.state = stateMessage
							return m, nil
						}
						entry = core.Entry{
							Username:  username,
							Password:  password,
							Notes:     notes,
							CreatedAt: time.Now().Format(time.RFC3339),
						}
					}
					
					m.vault.Entries[service] = entry
					if err := storage.SaveVault(m.vaultPath, m.masterPw, m.vault); err != nil {
						m.msg = "Failed to save vault: " + err.Error()
						m.isError = true
					} else {
						m.msg = "Credential saved successfully!"
						m.isError = false
						m.updateList()
					}
					m.state = stateMessage
					return m, nil
				}

				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}

				if m.focusIndex > len(m.formInputs)-1 {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.formInputs) - 1
				}
				
				if m.isEditing && m.focusIndex == 0 {
					if s == "up" || s == "shift+tab" {
						m.focusIndex = len(m.formInputs) - 1
					} else {
						m.focusIndex = 1
					}
				}

				for i := 0; i <= len(m.formInputs)-1; i++ {
					if i == m.focusIndex {
						m.formInputs[i].Focus()
						m.formInputs[i].PromptStyle = focusStyle
						m.formInputs[i].TextStyle = focusStyle
						continue
					}
					m.formInputs[i].Blur()
					m.formInputs[i].PromptStyle = blurStyle
					m.formInputs[i].TextStyle = blurStyle
				}
				return m, textinput.Blink
			}
		}
		
		for i := range m.formInputs {
			m.formInputs[i], cmd = m.formInputs[i].Update(msg)
			cmds = append(cmds, cmd)
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
			case "d":
				m.state = stateConfirmDelete
			case "e":
				m.setupForm(m.selectedItem.service, m.selectedItem.username, m.selectedItem.password, m.selectedItem.notes, true)
				m.state = stateForm
			}
		}
	case stateConfirmDelete:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "y", "Y":
				delete(m.vault.Entries, m.selectedItem.service)
				if err := storage.SaveVault(m.vaultPath, m.masterPw, m.vault); err != nil {
					m.msg = "Failed to delete: " + err.Error()
					m.isError = true
				} else {
					m.msg = "Deleted successfully!"
					m.isError = false
					m.updateList()
				}
				m.state = stateMessage
			case "n", "N":
				m.state = stateList
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
			"%s\n\nService: %s\nUsername: %s\nPassword: %s\nNotes: %s\nCreated: %s\nUpdated: %s\n\n[c] Copy Password  [e] Edit  [d] Delete  [esc] Back",
			titleStyle.Render("View Credential"),
			m.selectedItem.service,
			m.selectedItem.username,
			m.selectedItem.password,
			m.selectedItem.notes,
			m.selectedItem.createdAt,
			m.selectedItem.updatedAt,
		)
	case stateForm:
		title := "Add Credential"
		if m.isEditing {
			title = "Edit Credential"
		}
		s = titleStyle.Render(title) + "\n\n"
		for i := range m.formInputs {
			s += m.formInputs[i].View() + "\n"
		}
		s += "\n\n" + helpStyle.Render("tab/up/down: Move | enter (on last): Save | esc: Cancel")
	case stateConfirmDelete:
		s = fmt.Sprintf(
			"\nAre you sure you want to delete '%s'? (y/N)",
			m.selectedItem.service,
		)
	case stateMessage:
		style := infoStyle
		if m.isError {
			style = errorStyle
		}
		s = fmt.Sprintf("\n%s\n\nPress Enter to continue.", style.Render(m.msg))
	case stateAudit:
		s = m.auditReport
	}
	return appStyle.Render(s)
}
