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
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			if m.Selected == true {
				m.Expanded = true
				cmd = m.setNodes([]ExplorerNodeModel{
					ExplorerNodeModel{Title: "users", Type: "tables", Parent: &m},
					ExplorerNodeModel{Title: "accounts", Type: "tables", Parent: &m},
					ExplorerNodeModel{Title: "pages", Type: "tables", Parent: &m},
					ExplorerNodeModel{Title: "transactions", Type: "tables", Parent: &m},
					ExplorerNodeModel{Title: "categories", Type: "tables", Parent: &m},
					ExplorerNodeModel{Title: "tokens", Type: "tables", Parent: &m},
				})
			}
		case "h":
			if m.Selected == true {
				m.Expanded = false
				cmd = m.setNodes([]ExplorerNodeModel{})
			}
		}
	case SetChildrenMsg:
		if m.Selected == true {
			m.Children = msg
			return m, cmd
		}
	}

	var childCmds []tea.Cmd
	for i, child := range m.Children {
		updatedChild, childCmd := child.Update(msg)
		m.Children[i] = updatedChild.(ExplorerNodeModel)
		childCmds = append(childCmds, childCmd)
	}

	return m, tea.Batch(cmd, tea.Batch(childCmds...))
}

func (m ExplorerNodeModel) View() string {
	var childrenView string
	for _, child := range m.Children {
		childrenView += child.View()
	}
	if m.Parent != nil {
		var symbol string
		symbol = "+"
		if m.Expanded {
			symbol = "-"
		}
		text := symbol + " " + m.Title
		var textLine string
		if m.Selected {
			textLine = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(text)
		} else {
			textLine = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(text)
		}
		return textLine + "\n" + childrenView
	} else {
		return childrenView
	}
}
