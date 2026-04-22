package conn_manager

import (
	"fmt"
	"os"

	"app.lazygit/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	adapters "app.lazygit/internal/adapters"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionManager struct {
	layout             utils.ConnectionManagerLayout
	list               tea.Model
	form               tea.Model
	connections        []adapters.DbConnection
	selectedConnection int
	editingConnection  bool
	connecting         bool
	savingConnection   bool
	showHelp           bool
	connectionError    string
}

type SelectedConnectionMsg adapters.DbConnection
type EditConnectionMsg bool
type ConnectionErrorMsg string
type ConnectedMsg adapters.Database
type LayoutUpdated utils.ConnectionManagerLayout

func initializeNewConnection() adapters.DbConnection {
	return adapters.DbConnection{
		Name:     "New Connection",
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
	}
}

func calculateLayout(width int, height int) utils.ConnectionManagerLayout {
	return utils.CalculateConnectionManagerLayout(width, height)
}

func setLayout(width int, height int) tea.Cmd {
	return func() tea.Msg {
		return LayoutUpdated(calculateLayout(width, height))
	}
}

func InitConnectionManager() ConnectionManager {
	var connections []adapters.DbConnection
	connections = append(connections, initializeNewConnection())
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = 1
		height = 1
	}
	selectedConnection := connections[0]
	layout := calculateLayout(width, height)
	return ConnectionManager{
		layout:            layout,
		list:              InitConnectionList(connections, layout),
		form:              InitConnForm(selectedConnection, layout),
		connections:       connections,
		editingConnection: false,
		connecting:        false,
		connectionError:   "",
	}
}

func (m ConnectionManager) establishConnection() tea.Cmd {
	form := m.form.(ConnectionForm)
	connection := adapters.DbConnection{
		Host:     form.inputs[0].Value(),
		Port:     form.inputs[1].Value(),
		Username: form.inputs[2].Value(),
		Password: form.inputs[3].Value(),
		Driver:   "pgx",
	}
	return func() tea.Msg {
		database, err := connection.InitConnection()
		if err != nil {
			return ConnectionErrorMsg(fmt.Sprintf("Failed to connect: %s", err))
		}
		return ConnectedMsg(database)
	}
}

func (m ConnectionManager) toggleConnectionEdit() tea.Cmd {
	return func() tea.Msg {
		return EditConnectionMsg(m.editingConnection)
	}
}

func (m ConnectionManager) saveConnection() tea.Cmd {
	form := m.form.(ConnectionForm)
	connection := adapters.DbConnection{
		Name:     form.inputs[0].Value(),
		Host:     form.inputs[1].Value(),
		Port:     form.inputs[2].Value(),
		Username: form.inputs[3].Value(),
		Password: form.inputs[4].Value(),
	}
	return func() tea.Msg {
		m.connections = append(m.connections, connection)
		return SelectedConnectionMsg(connection)
	}
}

func (m ConnectionManager) Init() tea.Cmd {
	return tea.Batch(m.list.Init(), m.form.Init())
}

func (m ConnectionManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var listCmd, formCmd, command tea.Cmd
	m, command = m.handleKeyboardActions(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		command = setLayout(msg.Width, msg.Height)
	case ConnectionErrorMsg:
		m.connectionError = string(msg)
		m.connecting = false
	case LayoutUpdated:
		m.layout = utils.ConnectionManagerLayout(msg)
	}

	m.list, listCmd = m.list.Update(msg)
	m.form, formCmd = m.form.Update(msg)
	cmd := tea.Batch(listCmd, formCmd, command)
	return m, cmd
}

func (m ConnectionManager) handleKeyboardActions(msg tea.Msg) (ConnectionManager, tea.Cmd) {
	var command tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.connecting = true
			m.editingConnection = false
			connectionCommand := m.establishConnection()
			toggleConnectionCommand := m.toggleConnectionEdit()
			command = tea.Batch(connectionCommand, toggleConnectionCommand)
		case "e":
			if !m.editingConnection {
				m.editingConnection = true
				m.connecting = false
				m.connectionError = ""
				command = m.toggleConnectionEdit()
			}
		case "esc":
			if m.editingConnection {
				m.editingConnection = false
				command = m.toggleConnectionEdit()
			} else if m.showHelp {
				m.showHelp = false
			}
		case "s":
			m.savingConnection = true
			command = m.saveConnection()
		case "?":
			m.showHelp = !m.showHelp
		}
	}

	return m, command
}

func (m ConnectionManager) View() string {
	header := lipgloss.NewStyle().Width(m.layout.WinWidth).Height(m.layout.HeaderHeight).Padding(1).Render("Connection Manager")
	footer := m.buildFooter()

	listView := m.list.View()
	formView := m.form.View()
	listAndFormView := lipgloss.JoinHorizontal(lipgloss.Top, listView, formView)
	body := lipgloss.NewStyle().
		Width(m.layout.WinWidth).
		Border(lipgloss.NormalBorder(), true, false, true, false).
		Height(m.layout.BodyHeight).
		Render(listAndFormView)

	container := utils.Border().Width(m.layout.WinWidth).Height(m.layout.WinHeight).Render(
		lipgloss.JoinVertical(lipgloss.Top, header, body, footer),
	)

	base := lipgloss.Place(m.layout.ScreenWidth, m.layout.ScreenHeight, lipgloss.Center, lipgloss.Center, container)

	if m.showHelp {
		helpView := m.renderHelp()
		return lipgloss.Place(m.layout.ScreenWidth, m.layout.ScreenHeight, lipgloss.Center, lipgloss.Center, helpView)
	}

	return base
}

func (m ConnectionManager) renderHelp() string {
	helpText := `Connection Manager Help
- Name is for connection name that will appear in the list
- Driver is used to establish and find the database server, user
 * pgx for PostgreSQL
- Quit this dialog by hitting "?" or "esc"
	`
	return lipgloss.NewStyle().
		Width(m.layout.HelpWidth).
		Height(m.layout.HelpHeight).
		Padding(1).
		Border(lipgloss.NormalBorder()).
		Render(helpText)
}

func (m ConnectionManager) buildFooter() string {
	var footerContent string
	if m.connectionError != "" {
		footerContent = errorFooter(m.connectionError)
	} else if m.editingConnection {
		footerContent = editFooter()
	} else if m.connecting {
		footerContent = connectingFooter()
	} else if !m.editingConnection {
		footerContent = normalFooter()
	}
	return lipgloss.NewStyle().Width(m.layout.WinWidth).Height(m.layout.FooterHeight).Padding(1).Render(fmt.Sprintf("%s", footerContent))
}

func editFooter() string {
	return fmt.Sprintf("%s, %s, %s",
		"Save (enter)",
		"Cancel (esc)",
		"Navigate (tab, shift+tab)",
	)
}

func normalFooter() string {
	return fmt.Sprintf("%s, %s, %s, %s, %s",
		"Connect (enter)",
		"Edit (e)",
		"Save (s)",
		"Navigate (j,k)",
		"Help (?)",
	)
}

func errorFooter(errorMessage string) string {
	error_message := lipgloss.NewStyle().Foreground(lipgloss.Color("161")).Render(errorMessage)
	return fmt.Sprintf("%s\n%s", error_message, "Press 'e' to edit connection details.")
}

func connectingFooter() string {
	return "Connecting..."
}
