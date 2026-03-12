package explorer

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type SetChildrenMsg []ExplorerNodeModel
type SetSelectedMsg bool

type ExplorerNodeModel struct {
	Title    string
	Children []ExplorerNodeModel
	Selected bool
	Expanded bool
	Parent   *ExplorerNodeModel
	Type     string
	Index    int
}

func (m ExplorerNodeModel) setNodes(nodes []ExplorerNodeModel) tea.Cmd {
	return func() tea.Msg {
		return SetChildrenMsg(nodes)
	}
}

func (m ExplorerNodeModel) setSelected(selected bool) tea.Cmd {
	return func() tea.Msg {
		fmt.Printf("Setting selected to %v for node %s\n", selected, m.Title)
		return SetSelectedMsg(selected)
	}
}

func (m ExplorerNodeModel) Init() tea.Cmd {
	return nil
}

func (m ExplorerNodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.Selected == true {
		cmd = m.handleKeyboardActions(msg)
	}
	switch msg := msg.(type) {
	case SetChildrenMsg:
		if m.Selected == true {
			m.Children = msg
			return m, cmd
		}
	case SetSelectedMsg:
		m.Selected = bool(msg)
	}

	var childCmds []tea.Cmd
	for i, child := range m.Children {
		updatedChild, childCmd := child.Update(msg)
		m.Children[i] = updatedChild.(ExplorerNodeModel)
		childCmds = append(childCmds, childCmd)
	}

	return m, tea.Batch(cmd, tea.Batch(childCmds...))
}

func (m ExplorerNodeModel) handleKeyboardActions(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			if m.Selected == true {
				m.Expanded = true
				cmd = m.setNodes([]ExplorerNodeModel{
					ExplorerNodeModel{Title: "users", Type: "tables", Parent: &m, Index: 0},
					ExplorerNodeModel{Title: "accounts", Type: "tables", Parent: &m, Index: 1},
					ExplorerNodeModel{Title: "pages", Type: "tables", Parent: &m, Index: 2},
					ExplorerNodeModel{Title: "transactions", Type: "tables", Parent: &m, Index: 3},
					ExplorerNodeModel{Title: "categories", Type: "tables", Parent: &m, Index: 4},
					ExplorerNodeModel{Title: "tokens", Type: "tables", Parent: &m, Index: 5},
				})
			}
		case "h":
			if m.Selected == true {
				m.Expanded = false
				cmd = m.setNodes([]ExplorerNodeModel{})
			}
		case "j":
			if m.Expanded {
				child := m.Children[0]
				m.Selected = false
				child.Selected = true
				cmd = m.setNodes(m.Children)
			} else {
				siblings := m.Parent.Children
				nextSiblingIndex := m.Index + 1
				nextSibling := siblings[nextSiblingIndex]
				nextSibling.Selected = true
				m.Selected = false
				cmd = m.Parent.setNodes(m.Parent.Children)
			}
		}
	}
	return cmd
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
