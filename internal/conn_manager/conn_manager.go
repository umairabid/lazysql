package conn_manager

import (
	"fmt"
	"os"
	"slices"

	"app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var MIN_WIDTH = 80
var MIN_HEIGHT = 24

type ConnectionManager struct {
	width              int
	height             int
	list               tea.Model
	form               tea.Model
	connections        []Connection
	selectedConnection int
	editingConnection  bool
}

type SelectedConnectionMsg Connection
type EditConnectionMsg bool

func initializeNewConnection(host string) Connection {
	return Connection{
		name:     "New Connection",
		host:     host,
		port:     "5432",
		username: "user",
		password: "password",
	}
}

func InitConnectionManager() ConnectionManager {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = 80
		height = 24
	}
	connections := []Connection{initializeNewConnection("localhost"), initializeNewConnection("pocalhost"), initializeNewConnection("totalhost")}
	selectedConnection := connections[0]
	return ConnectionManager{
		width:             width,
		height:            height,
		list:              InitConnectionList(connections),
		form:              InitConnForm(selectedConnection),
		connections:       connections,
		editingConnection: false,
	}
}

func (m ConnectionManager) establishConnection() tea.Cmd {
	form := m.form.(ConnectionForm)
	connection := Connection{
		host:     form.inputs[0].Value(),
		port:     form.inputs[1].Value(),
		username: form.inputs[2].Value(),
		password: form.inputs[3].Value(),
		driver:   "pgx",
	}
	return func() tea.Msg {
		_, err := connectWithDatabase(connection)
		if err != nil {
			return fmt.Sprintf("Failed to connect: %s", err)
		}
		return "Connection established successfully!"
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
			command = m.establishConnection()
		case "e":
			if !m.editingConnection {
				m.editingConnection = true
			}
			command = m.toggleConnectionEdit()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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

	header := utils.BottomBorder().Width(width).Padding(1).Render("Connection Manager")
	footer := buildFooter(width)
	bodyHeight := height - (lipgloss.Height(header) + lipgloss.Height(footer))

	listView := utils.RightBorder().Width(width / 2).Height(bodyHeight).Render(m.list.View())
	formView := lipgloss.NewStyle().Width(width/2).Height(bodyHeight).Padding(1, 2).Render(m.form.View())
	listAndFormView := lipgloss.JoinHorizontal(lipgloss.Top, listView, formView)
	body := lipgloss.NewStyle().Height(bodyHeight - 3).Render(listAndFormView)

	container := utils.Border().Width(width).Height(height).Render(
		fmt.Sprintf("%s\n%s\n%s", header, body, footer),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, container)
}

func buildFooter(width int) string {
	return utils.TopBorder().Width(width).Padding(1).Render(fmt.Sprintf("%s, %s, %s, %s",
		"Connect (enter)",
		"Edit (e)",
		"Save (s)",
		"Navigate (j,k)",
	))
}
