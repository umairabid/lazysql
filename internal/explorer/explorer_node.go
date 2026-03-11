package explorer

import (
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type SetChildrenMsg []ExplorerNodeModel

type ExplorerNodeModel struct {
	Title    string
	Children []ExplorerNodeModel
	Selected bool
	Expanded bool
	Parent   *ExplorerNodeModel
	Type     string
}

func (m *ExplorerNodeModel) setNodes(nodes []ExplorerNodeModel) tea.Cmd {
	return func() tea.Msg {
		return SetChildrenMsg(nodes)
	}
}

func (m ExplorerNodeModel) Init() tea.Cmd {
	return nil
}

func (m ExplorerNodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case SetChildrenMsg:
		m.Children = msg
	}
	return m, cmd
}

func (m ExplorerNodeModel) View() string {
	var childrenView string
	var symbol string
	symbol = "+"
	if m.Expanded {
		symbol = "-"
	}
	for _, child := range m.Children {
		childrenView += "  " + child.View()
	}
	if m.Parent != nil {
		text := symbol + " " + m.Title + "\n" + childrenView
		if m.Selected {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(text)
		} else {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(text)
		}
	} else {
		return childrenView
	}
}
