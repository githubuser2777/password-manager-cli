package tui

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"password-manager-cli/internal/sys"
	"password-manager-cli/internal/vault"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Custom Keys for Help Menu
type listKeyMap struct {
	enter       key.Binding
	copy        key.Binding
	add         key.Binding
	edit        key.Binding
	delete      key.Binding
	auditLocal  key.Binding
	auditOnline key.Binding
	exportCsv   key.Binding
	importCsv   key.Binding
}

var customKeys = listKeyMap{
	enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view"),
	),
	copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy"),
	),
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
	exportCsv: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "export csv"),
	),
	importCsv: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "import csv"),
	),
}

// Premium UI Styles
var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("141"))
	appStyle        = lipgloss.NewStyle().Padding(1, 2)
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	infoStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Bold(true)
	focusStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	blurStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	detailCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("141")).
			Padding(1, 2).
			MarginLeft(2)
	activeTabStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("212")).
			Foreground(lipgloss.Color("16")).
			Padding(0, 2).
			Bold(true)
	inactiveTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(0, 2)
)

type state int

const (
	stateLogin state = iota
	stateDecrypting
	stateList
	stateView
	stateMessage
	stateForm
	stateConfirmDelete
	stateConfirmReset
	stateInit
	stateAudit
	stateAuditing
	stateFilePickerImport
	stateFilePickerExport
	stateExportFilename
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
	state     state
	vaultPath string
	masterPw  []byte
	vault     *vault.Vault

	passwordInput  textinput.Model
	servicesList   list.Model
	spinner        spinner.Model
	fp             filepicker.Model
	auditList      list.Model

	confirmInput    textinput.Model
	initFocusIndex  int
	initError       string
	initSuggestedPw string

	exportNameInput textinput.Model

	activeTab      int
	genLengthInput textinput.Model
	genUseNumbers  bool
	genUseSymbols  bool
	genFocusIndex  int

	formInputs []textinput.Model
	focusIndex int
	isEditing  bool

	selectedItem item
	msg          string
	isError      bool
	auditReport  string

	width  int
	height int
}

type auditItem struct {
	title string
	desc  string
}

func (i auditItem) Title() string       { return i.title }
func (i auditItem) Description() string { return i.desc }
func (i auditItem) FilterValue() string { return i.title }

func initialModel(vaultPath string) model {
	ti := textinput.New()
	ti.Placeholder = "Master Password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.Focus()

	ci := textinput.New()
	ci.Placeholder = "Confirm Password"
	ci.EchoMode = textinput.EchoPassword
	ci.EchoCharacter = '•'

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Password Vault"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{customKeys.enter, customKeys.copy, customKeys.add, customKeys.edit, customKeys.delete}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{customKeys.enter, customKeys.copy, customKeys.add, customKeys.edit, customKeys.delete}
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = titleStyle

	fp := filepicker.New()
	homeDir, _ := os.UserHomeDir()
	defaultDir := fmt.Sprintf("%s%c%s%c%s", homeDir, os.PathSeparator, "Documents", os.PathSeparator, "passmgr-cli")
	os.MkdirAll(defaultDir, 0700)
	fp.CurrentDirectory = defaultDir
	fp.ShowHidden = true

	eni := textinput.New()
	eni.Placeholder = "vault_export.csv"
	eni.SetValue("vault_export.csv")

	genTi := textinput.New()
	genTi.Placeholder = "16"
	genTi.SetValue("16")
	genTi.CharLimit = 3
	genTi.PromptStyle = focusStyle
	genTi.TextStyle = focusStyle

	st := stateLogin
	var suggestedPw string
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		st = stateInit
		suggestedPw, _ = vault.GenerateRandomPassword(16, true, true)
	}

	return model{
		state:           st,
		vaultPath:       vaultPath,
		passwordInput:   ti,
		confirmInput:    ci,
		servicesList:    l,
		spinner:         s,
		fp:              fp,
		initSuggestedPw: suggestedPw,
		exportNameInput: eni,
		activeTab:      0,
		genLengthInput: genTi,
		genUseNumbers:  true,
		genUseSymbols:  true,
		genFocusIndex:  0,
		auditList:      list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
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

func (m model) runAuditCmd(online bool) tea.Cmd {
	return func() tea.Msg {
		report := vault.RunAudit(m.vault, online)

		var items []list.Item
		for _, msg := range report.Messages {
			desc := "Severity: " + msg.Level
			items = append(items, auditItem{title: msg.Message, desc: desc})
		}

		titleStr := fmt.Sprintf("Audit: %d Checked | %d Weak | %d Reused | %d Pwned", report.TotalChecked, report.WeakCount, report.ReusedCount, report.PwnedCount)
		return auditResultMsg{items: items, title: titleStr}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.fp.Init())
}

type decryptResultMsg struct {
	vault *vault.Vault
	err   error
}

type auditResultMsg struct {
	items []list.Item
	title string
}

func (m model) decryptVaultCmd() tea.Cmd {
	return func() tea.Msg {
		vault, err := vault.LoadVault(m.vaultPath, m.masterPw)
		return decryptResultMsg{vault: vault, err: err}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.state == stateInit {
				return m, tea.Quit
			}
			if m.state == stateView || m.state == stateMessage || m.state == stateForm || m.state == stateConfirmDelete || m.state == stateConfirmReset || m.state == stateAudit || m.state == stateFilePickerImport || m.state == stateFilePickerExport || m.state == stateExportFilename {
				if m.vault != nil {
					m.state = stateList
					return m, nil
				}
				m.state = stateLogin
				return m, nil
			}
			// Zero out key on exit
			if m.masterPw != nil {
				vault.ZeroBytes(m.masterPw)
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.fp.Height = m.height - 10
		h, v := appStyle.GetFrameSize()
		m.servicesList.SetSize(msg.Width-h, msg.Height-v)
		m.auditList.SetSize(msg.Width-h, msg.Height-v)
	case spinner.TickMsg:
		if m.state == stateDecrypting || m.state == stateAuditing {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case decryptResultMsg:
		if msg.err != nil {
			m.msg = "Invalid Master Password or Vault not found.\nError: " + msg.err.Error()
			m.isError = true
			m.state = stateMessage
			if m.masterPw != nil {
				vault.ZeroBytes(m.masterPw)
				m.masterPw = nil
			}
		} else {
			m.vault = msg.vault
			m.updateList()
			m.state = stateList
		}
		return m, nil
	case auditResultMsg:
		m.auditList.Title = msg.title
		m.auditList.Styles.Title = titleStyle
		m.auditList.SetItems(msg.items)
		m.auditList.SetShowStatusBar(true)
		m.auditList.SetShowHelp(true)
		m.auditList.AdditionalShortHelpKeys = func() []key.Binding { return nil }
		m.auditList.AdditionalFullHelpKeys = func() []key.Binding { return nil }
		m.state = stateAudit
		return m, nil
	}

	switch m.state {
	case stateLogin:
		m.passwordInput, cmd = m.passwordInput.Update(msg)
		cmds = append(cmds, cmd)

		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.Type {
			case tea.KeyEnter:
				rawPw := []byte(m.passwordInput.Value())
				m.masterPw = make([]byte, len(rawPw))
				copy(m.masterPw, rawPw)

				// Zero raw values
				vault.ZeroBytes(rawPw)
				m.passwordInput.SetValue("")

				m.state = stateDecrypting
				return m, tea.Batch(m.spinner.Tick, m.decryptVaultCmd())
			case tea.KeyCtrlR:
				m.state = stateConfirmReset
				return m, nil
			}
		}
	case stateDecrypting:
		// Let the spinner update handle this
	case stateList:
		if msg, ok := msg.(tea.KeyMsg); ok {
			isFiltering := m.activeTab == 0 && m.servicesList.FilterState() == list.Filtering
			isGenTyping := m.activeTab == 1 && m.genFocusIndex == 0

			if !isFiltering && !isGenTyping {
				switch msg.String() {
				case "1":
					m.activeTab = 0
					return m, nil
				case "2":
					m.activeTab = 1
					m.genFocusIndex = 0
					m.genLengthInput.Focus()
					return m, nil
				case "left", "right", "h", "l":
					m.activeTab = 1 - m.activeTab
					if m.activeTab == 1 {
						m.genFocusIndex = 0
						m.genLengthInput.Focus()
					}
					return m, nil
				}
			}

			if !isFiltering {
				switch msg.String() {
				case "tab", "shift+tab":
					m.activeTab = 1 - m.activeTab
					if m.activeTab == 1 {
						m.genFocusIndex = 0
						m.genLengthInput.Focus()
					}
					return m, nil
				}
			}

			if m.activeTab == 1 {
				switch msg.String() {
				case "up":
					m.genFocusIndex--
					if m.genFocusIndex < 0 {
						m.genFocusIndex = 3
					}
				case "down":
					m.genFocusIndex++
					if m.genFocusIndex > 3 {
						m.genFocusIndex = 0
					}
				case "enter", " ":
					if m.genFocusIndex == 1 {
						m.genUseNumbers = !m.genUseNumbers
					} else if m.genFocusIndex == 2 {
						m.genUseSymbols = !m.genUseSymbols
					} else if m.genFocusIndex == 3 && msg.String() == "enter" {
						l := 16
						fmt.Sscanf(m.genLengthInput.Value(), "%d", &l)
						if l < 4 { l = 4 }
						pw, _ := vault.GenerateRandomPassword(l, m.genUseNumbers, m.genUseSymbols)
						_ = sys.WriteClipboard(pw)
						m.msg = "Generated password copied to clipboard! (Auto-clears in 30s)"
						m.isError = false
						time.AfterFunc(30*time.Second, func() { sys.WriteClipboard("") })
						m.state = stateMessage
					}
				case "r":
					m.state = stateAuditing
					return m, tea.Batch(m.spinner.Tick, m.runAuditCmd(false))
				case "R":
					m.state = stateAuditing
					return m, tea.Batch(m.spinner.Tick, m.runAuditCmd(true))
				case "x":
					m.fp.DirAllowed = true
					m.fp.FileAllowed = false
					m.state = stateFilePickerExport
				case "i":
					m.fp.DirAllowed = false
					m.fp.FileAllowed = true
					m.fp.AllowedTypes = []string{".csv"}
					m.state = stateFilePickerImport
				}

				if m.genFocusIndex == 0 {
					m.genLengthInput.Focus()
					m.genLengthInput, cmd = m.genLengthInput.Update(msg)
					cmds = append(cmds, cmd)
				} else {
					m.genLengthInput.Blur()
				}
				return m, tea.Batch(cmds...)
			}
		}

		if m.activeTab == 0 {
			m.servicesList, cmd = m.servicesList.Update(msg)
			cmds = append(cmds, cmd)

			if msg, ok := msg.(tea.KeyMsg); ok && m.servicesList.FilterState() != list.Filtering {
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
				case "c":
					if i, ok := m.servicesList.SelectedItem().(item); ok {
						if err := sys.WriteClipboard(i.password); err != nil {
							m.msg = "Failed to copy password: " + err.Error()
							m.isError = true
						} else {
							m.msg = "Password copied to clipboard! (Auto-clears in 30s)"
							m.isError = false
							time.AfterFunc(30*time.Second, func() { sys.WriteClipboard("") })
						}
						m.state = stateMessage
					}
				}
			}
		}
	case stateFilePickerImport, stateFilePickerExport:
		var fpCmd tea.Cmd
		m.fp, fpCmd = m.fp.Update(msg)
		cmds = append(cmds, fpCmd)

		if didSelect, path := m.fp.DidSelectFile(msg); didSelect && m.state == stateFilePickerImport {
			f, err := os.Open(path)
			if err == nil {
				r := csv.NewReader(f)
				records, _ := r.ReadAll()
				count := 0
				for idx, rec := range records {
					if idx == 0 || len(rec) < 6 { continue }
					m.vault.Entries[rec[0]] = vault.Entry{
						Username: rec[1], Password: rec[2], Notes: rec[3], CreatedAt: rec[4], UpdatedAt: rec[5],
					}
					count++
				}
				f.Close()
				vault.SaveVault(m.vaultPath, m.masterPw, m.vault)
				m.updateList()
				m.msg = fmt.Sprintf("Imported %d entries from %s", count, path)
				m.isError = false
			} else {
				m.msg = "Import failed: " + err.Error()
				m.isError = true
			}
			m.state = stateMessage
		}

		if didSelect, path := m.fp.DidSelectFile(msg); didSelect && m.state == stateFilePickerExport {
			m.fp.CurrentDirectory = path
			m.exportNameInput.Focus()
			m.state = stateExportFilename
		}
	case stateExportFilename:
		var cmd tea.Cmd
		m.exportNameInput, cmd = m.exportNameInput.Update(msg)
		cmds = append(cmds, cmd)

		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			filename := m.exportNameInput.Value()
			if filename == "" { filename = "vault_export.csv" }
			if !strings.HasSuffix(filename, ".csv") { filename += ".csv" }

			exportPath := fmt.Sprintf("%s%c%s", m.fp.CurrentDirectory, os.PathSeparator, filename)
			f, err := os.Create(exportPath)
			if err == nil {
				w := csv.NewWriter(f)
				w.Write([]string{"Service", "Username", "Password", "Notes", "CreatedAt", "UpdatedAt"})
				for s, e := range m.vault.Entries {
					w.Write([]string{s, e.Username, e.Password, e.Notes, e.CreatedAt, e.UpdatedAt})
				}
				w.Flush()
				f.Close()
				m.msg = fmt.Sprintf("Exported to %s", exportPath)
				m.isError = false
			} else {
				m.msg = "Export failed: " + err.Error()
				m.isError = true
			}
			m.state = stateMessage
		}
	case stateForm:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				if s == "enter" && m.focusIndex == len(m.formInputs)-1 {
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
						entry = vault.Entry{
							Username:  username,
							Password:  password,
							Notes:     notes,
							CreatedAt: time.Now().Format(time.RFC3339),
						}
					}

					m.vault.Entries[service] = entry
					if err := vault.SaveVault(m.vaultPath, m.masterPw, m.vault); err != nil {
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
				if err := sys.WriteClipboard(m.selectedItem.password); err != nil {
					m.msg = "Failed to copy password: " + err.Error()
					m.isError = true
				} else {
					m.msg = "Password copied to clipboard! (Auto-clears in 30s)"
					m.isError = false
					time.AfterFunc(30*time.Second, func() { sys.WriteClipboard("") })
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
				if err := vault.SaveVault(m.vaultPath, m.masterPw, m.vault); err != nil {
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
	case stateConfirmReset:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "y", "Y":
				_ = os.Remove(m.vaultPath)
				m.passwordInput.SetValue("")
				m.confirmInput.SetValue("")
				m.initFocusIndex = 0
				m.passwordInput.Focus()
				m.confirmInput.Blur()
				pw, _ := vault.GenerateRandomPassword(16, true, true)
				m.initSuggestedPw = pw
				m.state = stateInit
			case "n", "N":
				m.state = stateLogin
			}
		}
	case stateInit:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "tab", "shift+tab", "up", "down":
				if m.initFocusIndex == 0 {
					m.initFocusIndex = 1
					m.passwordInput.Blur()
					m.confirmInput.Focus()
				} else {
					m.initFocusIndex = 0
					m.passwordInput.Focus()
					m.confirmInput.Blur()
				}
				return m, textinput.Blink
			case "enter":
				if m.initFocusIndex == 0 {
					m.initFocusIndex = 1
					m.passwordInput.Blur()
					m.confirmInput.Focus()
					return m, textinput.Blink
				}
				pw := m.passwordInput.Value()
				conf := m.confirmInput.Value()
				if pw != conf {
					m.initError = "Passwords do not match."
					return m, nil
				}
				rawPw := []byte(pw)
				if err := vault.ValidateMasterPassword(rawPw); err != nil {
					m.initError = "Weak password: " + err.Error()
					return m, nil
				}
				salt, err := vault.GenerateSalt(16)
				if err != nil {
					m.initError = "Failed to generate salt: " + err.Error()
					return m, nil
				}
				v := &vault.Vault{
					Salt:    salt,
					Entries: make(map[string]vault.Entry),
				}
				m.masterPw = make([]byte, len(rawPw))
				copy(m.masterPw, rawPw)
				vault.ZeroBytes(rawPw)
				if err := vault.SaveVault(m.vaultPath, m.masterPw, v); err != nil {
					m.initError = "Failed to save vault: " + err.Error()
					return m, nil
				}
				m.vault = v
				m.passwordInput.SetValue("")
				m.confirmInput.SetValue("")
				m.initError = ""
				m.updateList()
				m.state = stateList
				return m, nil
			default:
				m.initError = ""
			}
		}
		if m.initFocusIndex == 0 {
			m.passwordInput, cmd = m.passwordInput.Update(msg)
		} else {
			m.confirmInput, cmd = m.confirmInput.Update(msg)
		}
		cmds = append(cmds, cmd)

	case stateMessage:
		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			if m.vault != nil {
				m.state = stateList
			} else {
				m.state = stateLogin
				m.passwordInput.SetValue("")
			}
		}
	case stateAuditing:
		// wait for audit to finish
	case stateAudit:
		m.auditList, cmd = m.auditList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) renderStrengthMeter() string {
	pwVal := m.formInputs[2].Value()
	if len(pwVal) == 0 {
		return ""
	}

	score, _ := vault.EvaluatePasswordStrength(pwVal)

	var bar string
	var label string
	var hint string
	var color lipgloss.Color

	if score <= 2 {
		bar = "██░░░░"
		label = "Weak"
		color = lipgloss.Color("203")
		hint = " (Too short or missing types)"
	} else if score <= 4 {
		bar = "████░░"
		label = "Medium"
		color = lipgloss.Color("228")
		hint = " (Add symbols/numbers)"
	} else {
		bar = "██████"
		label = "Strong"
		color = lipgloss.Color("120")
		hint = ""
	}

	style := lipgloss.NewStyle().Foreground(color)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	return fmt.Sprintf("  %s %s%s", style.Render("Strength: ["+bar+"]"), style.Render(label), hintStyle.Render(hint))
}

func (m model) View() string {
	var s string
	switch m.state {
	case stateLogin:
		s = fmt.Sprintf(
			"%s\n\n%s\n\n(esc to quit)\n[ctrl+r] Forgot Password? (Reset Vault)",
			titleStyle.Render("Unlock Vault"),
			m.passwordInput.View(),
		)
	case stateInit:
		errText := ""
		if m.initError != "" {
			errText = "\n\n" + errorStyle.Render(m.initError)
		}
		suggestText := ""
		if m.initSuggestedPw != "" {
			suggestText = "\n\n" + helpStyle.Render("Suggested strong password: ") + lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(m.initSuggestedPw)
		}
		s = fmt.Sprintf(
			"%s%s%s\n\n%s\n%s\n\n(tab: switch | enter: submit | esc: quit)",
			titleStyle.Render("Initialize New Vault"),
			suggestText,
			errText,
			m.passwordInput.View(),
			m.confirmInput.View(),
		)
	case stateDecrypting:
		s = fmt.Sprintf(
			"\n\n  %s Deriving keys and decrypting vault...\n\n",
			m.spinner.View(),
		)
	case stateList:
		var tabs string
		if m.activeTab == 0 {
			tabs = lipgloss.JoinHorizontal(lipgloss.Top, activeTabStyle.Render("Vault"), inactiveTabStyle.Render("Tools"))
		} else {
			tabs = lipgloss.JoinHorizontal(lipgloss.Top, inactiveTabStyle.Render("Vault"), activeTabStyle.Render("Tools"))
		}
		s = tabs + "\n\n"

		if m.activeTab == 0 {
			var listContent string
			if m.width >= 80 {
				listWidth := 38
				m.servicesList.SetSize(listWidth, m.height-6)

				var detailContent string
				if selected, ok := m.servicesList.SelectedItem().(item); ok {
					notes := selected.notes
					if notes == "" { notes = "(none)" }
					detailContent = fmt.Sprintf(
						"%s\n\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n\n%s",
						titleStyle.Render(selected.service),
						infoStyle.Render("Username"), selected.username,
						infoStyle.Render("Password"), "•••••••• (press Enter to view / copy)",
						infoStyle.Render("Notes"), notes,
						infoStyle.Render("Created"), selected.createdAt,
						infoStyle.Render("Updated"), selected.updatedAt,
						helpStyle.Render("[a] Add | [enter] View Full | [c] Copy\n[e] Edit | [d] Delete"),
					)
				} else {
					detailContent = "\n\n  No credentials saved yet.\n  Press [a] to add."
				}

				cardWidth := m.width - listWidth - 10
				if cardWidth > 60 {
					cardWidth = 60
				}
				cardHeight := m.height - 6
				if cardHeight > 20 {
					cardHeight = 20
				}
				card := detailCardStyle.Width(cardWidth).Height(cardHeight).Render(detailContent)

				listContent = lipgloss.JoinHorizontal(lipgloss.Top, m.servicesList.View(), card)
			} else {
				m.servicesList.SetSize(m.width-4, m.height-6)
				listContent = m.servicesList.View()
			}
			s += listContent
		} else {
			lenStyle := blurStyle
			if m.genFocusIndex == 0 { lenStyle = focusStyle }
			numStyle := blurStyle
			if m.genFocusIndex == 1 { numStyle = focusStyle }
			symStyle := blurStyle
			if m.genFocusIndex == 2 { symStyle = focusStyle }
			btnStyle := blurStyle
			if m.genFocusIndex == 3 { btnStyle = focusStyle }

			cbNum := "[ ]"
			if m.genUseNumbers { cbNum = "[x]" }
			cbSym := "[ ]"
			if m.genUseSymbols { cbSym = "[x]" }

			genStr := titleStyle.Render("Password Generator") + "\n\n"
			genStr += fmt.Sprintf("  %s %s\n", lenStyle.Render("Length: "), m.genLengthInput.View())
			genStr += fmt.Sprintf("  %s %s\n", numStyle.Render("Numbers:"), numStyle.Render(cbNum))
			genStr += fmt.Sprintf("  %s %s\n\n", symStyle.Render("Symbols:"), symStyle.Render(cbSym))
			genStr += fmt.Sprintf("  %s\n", btnStyle.Render("[ Generate & Copy ]"))

			auditStr := titleStyle.Render("Security Audit") + "\n\n  Press [r] to run Local Audit\n  Press [R] to run Online Audit"

			syncStr := titleStyle.Render("Data Sync") + "\n\n  Press [x] to Export CSV\n  Press [i] to Import CSV"

			s += lipgloss.JoinVertical(lipgloss.Left, genStr, "\n", auditStr, "\n", syncStr)
		}

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
			if i == 2 {
				s += m.renderStrengthMeter() + "\n"
			}
		}
		s += "\n\n" + helpStyle.Render("tab/up/down: Move | enter (on last): Save | esc: Cancel")
	case stateConfirmDelete:
		s = fmt.Sprintf(
			"\nAre you sure you want to delete '%s'? (y/N)",
			m.selectedItem.service,
		)
	case stateConfirmReset:
		s = fmt.Sprintf(
			"\n%s\n\nAre you sure you want to PERMANENTLY DELETE your vault and start over? (y/N)",
			errorStyle.Render("WARNING: ALL PASSWORDS WILL BE LOST!"),
		)
	case stateMessage:
		style := infoStyle
		if m.isError {
			style = errorStyle
		}
		s = fmt.Sprintf("\n%s\n\nPress Enter to continue.", style.Render(m.msg))
	case stateAuditing:
		s = fmt.Sprintf(
			"\n\n  %s Running Security Audit...\n\n  (Online checks may take a while depending on vault size)\n",
			m.spinner.View(),
		)
	case stateAudit:
		s = m.auditList.View() + "\n\n(esc to return)"
	case stateFilePickerImport:
		s = fmt.Sprintf("%s\n\n%s\n%s\n\n(esc to cancel)", titleStyle.Render("Select CSV file to import:"), infoStyle.Render("Path: "+m.fp.CurrentDirectory), m.fp.View())
	case stateFilePickerExport:
		s = fmt.Sprintf("%s\n\n%s\n%s\n\n(esc to cancel)", titleStyle.Render("Select directory to save export:"), infoStyle.Render("Path: "+m.fp.CurrentDirectory), m.fp.View())
	case stateExportFilename:
		s = fmt.Sprintf("%s\n\n%s\n\n(enter to save | esc to cancel)", titleStyle.Render(fmt.Sprintf("Enter filename to save in %s:", m.fp.CurrentDirectory)), m.exportNameInput.View())
	}
	return appStyle.Render(s)
}
