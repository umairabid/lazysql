package conn_manager

import (
	"fmt"
	"os"
	"slices"

	"app.lazygit/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	
	tea "github.com/charmbracelet/bubbletea"
	postgres "app.lazygit/internal/services/postgres"
)

var MIN_WIDTH = 80
var MIN_HEIGHT = 24

type ConnectionManager struct {
	width              int
	height             int
	list               tea.Model
	form               tea.Model
	connections        []postgres.Connection
	selectedConnection int
	editingConnection  bool
	connecting         bool
	connectionError    string
}

type SelectedConnectionMsg postgres.Connection
type EditConnectionMsg bool
type ConnectionErrorMsg string
type ConnectedMsg bool

func initializeNewConnection(host string) postgres.Connection {
	return postgres.Connection{
		Name:     "New Connection",
		Host:     host,
		Port:     "5432",
		Username: "user",
		Password: "password",
	}
}

func InitConnectionManager() ConnectionManager {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = MIN_WIDTH
		height = MIN_HEIGHT
	}
	connections := []postgres.Connection{initializeNewConnection("localhost"), initializeNewConnection("pocalhost"), initializeNewConnection("totalhost")}
	selectedConnection := connections[0]
	return ConnectionManager{
		width:             width,
		height:            height,
		list:              InitConnectionList(connections),
		form:              InitConnForm(selectedConnection),
		connections:       connections,
		editingConnection: false,
		connecting:        false,
		connectionError:   "",
	}
}

func (m ConnectionManager) establishConnection() tea.Cmd {
	form := m.form.(ConnectionForm)
	connection := postgres.Connection{
		Host:     form.inputs[0].Value(),
		Port:     form.inputs[1].Value(),
		Username: form.inputs[2].Value(),
		Password: form.inputs[3].Value(),
		Driver:   "pgx",
	}
	return func() tea.Msg {
		_, err := postgres.ConnectWithDatabase(connection)
		if err != nil {
			return ConnectionErrorMsg(fmt.Sprintf("Failed to connect: %s", err))
		}
		return ConnectedMsg(true)
	}
}

func (m ConnectionManager) toggleConnectionEdit() tea.Cmd {
	return func() tea.Msg {
		return EditConnectionMsg(m.editingConnection)
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
		m.width = msg.Width
		m.height = msg.Height
	case ConnectionErrorMsg:
		m.connectionError = string(msg)
		m.connecting = false
	}

	m.list, listCmd = m.list.Update(msg)
	m.form, formCmd = m.form.Update(msg)
	cmd := tea.Batch(listCmd, formCmd, command)
	return m, cmd
}

func (m ConnectionManager) View() string {
	widths := []int{MIN_WIDTH, m.width / 3}
	heights := []int{MIN_HEIGHT, m.height / 3}
	width := slices.Max(widths)
	height := slices.Max(heights)

	header := lipgloss.NewStyle().Width(width).Padding(1).Render("Connection Manager")
	footer := m.buildFooter()
	bodyHeight := height - (lipgloss.Height(header) + lipgloss.Height(footer))

	listView := utils.RightBorder().Width(width / 2).Height(bodyHeight).Render(m.list.View())
	formView := lipgloss.NewStyle().Width(width/2).Height(bodyHeight).Padding(1, 2).Render(m.form.View())
	listAndFormView := lipgloss.JoinHorizontal(lipgloss.Top, listView, formView)
	body := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.NormalBorder(), true, false, true, false).
		Height(bodyHeight - 3).
		Render(listAndFormView)

	container := utils.Border().Width(width).Height(height).Render(
		fmt.Sprintf("%s\n%s\n%s", header, body, footer),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, container)
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
	return lipgloss.NewStyle().Width(m.width).Padding(1).Render(fmt.Sprintf("%s", footerContent))

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
