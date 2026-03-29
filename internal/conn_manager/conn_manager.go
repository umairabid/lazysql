package conn_manager

import (
	"fmt"
	"os"
	"slices"

	"app.lazygit/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	adapters "app.lazygit/internal/adapters"
	tea "github.com/charmbracelet/bubbletea"
)

var MIN_WIDTH = 80
var MIN_HEIGHT = 24

type ConnectionManager struct {
	layout utils.ConnectionManagerLayout
	list               tea.Model
	form               tea.Model
	connections        []adapters.DbConnection
	selectedConnection int
	editingConnection  bool
	connecting         bool
	connectionError    string
}

type SelectedConnectionMsg adapters.DbConnection
type EditConnectionMsg bool
type ConnectionErrorMsg string
type ConnectedMsg adapters.Database
type LayoutUpdated utils.ConnectionManagerLayout

func initializeNewConnection(host string) adapters.DbConnection {
	return adapters.DbConnection{
		Name:     "New Connection New Connection New Connection New Connection New Connection New Connection New Connection New Connection New Connection New Connection New Connection New Connection New Connection ",
		Host:     host,
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
	}
}

func InitConnectionManager() ConnectionManager {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = MIN_WIDTH
		height = MIN_HEIGHT
	}
	connections := []adapters.DbConnection{
		initializeNewConnection("localhost"), 
		initializeNewConnection("pocalhost"), 
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
		initializeNewConnection("totalhost"),
	}
	selectedConnection := connections[0]
	layout := utils.CalculateConnectionManagerLayout(width, height)
	return ConnectionManager{
		layout: layout,
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

func (m ConnectionManager) setLayout(width int, height int) tea.Cmd {
	return func() tea.Msg {
	widths := []int{MIN_WIDTH, width / 3}
	heights := []int{MIN_HEIGHT, height / 3}
	width = slices.Max(widths)
	height = slices.Max(heights)
	return LayoutUpdated(utils.CalculateConnectionManagerLayout(width, height))
	}
}

func (m ConnectionManager) Init() tea.Cmd {
	return tea.Batch(m.list.Init(), m.form.Init())
}

func (m ConnectionManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var listCmd, formCmd, command tea.Cmd
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
		}
	case tea.WindowSizeMsg:
		command = m.setLayout(msg.Width, msg.Height)
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

func (m ConnectionManager) View() string {

	header := lipgloss.NewStyle().Width(m.layout.WinWidth).Padding(1).Render("Connection Manager")
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
		fmt.Sprintf("%s\n%s\n%s", header, body, footer),
	)

	return lipgloss.Place(m.layout.WinWidth, m.layout.WinHeight, lipgloss.Center, lipgloss.Center, container)
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
	return lipgloss.NewStyle().Width(m.layout.WinWidth).Padding(1).Render(fmt.Sprintf("%s", footerContent))
}

func editFooter() string {
	return fmt.Sprintf("%s, %s, %s",
		"Save (enter)",
		"Cancel (esc)",
		"Navigate (tab, shift+tab)",
	)
}

func normalFooter() string {
	return fmt.Sprintf("%s, %s, %s",
		"Connect (enter)",
		"Edit (e)",
		"Navigate (j,k)",
	)
}

func errorFooter(errorMessage string) string {
	error_message := lipgloss.NewStyle().Foreground(lipgloss.Color("161")).Render(errorMessage)
	return fmt.Sprintf("%s\n%s", error_message, "Press 'e' to edit connection details.")
}

func connectingFooter() string {
	return "Connecting..."
}
