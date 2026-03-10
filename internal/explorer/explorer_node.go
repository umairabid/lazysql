package explorer

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
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
	fmt.Println("Setting nodes for", m.Title)
	return func() tea.Msg {
		return SetChildrenMsg(nodes)
	}
}

func (m ExplorerNodeModel) Init() tea.Cmd {
	return nil
}

func (m ExplorerNodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Println(msg)
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case SetChildrenMsg:
		m.Children = msg
	}
	return m, cmd
}

func (m ExplorerNodeModel) View() string {
	fmt.Println("Rendering node:", m.Title)
	indicator := "L"
	if m.Selected {
		indicator = ">"
	}
	var childrenView string
		for _, child := range m.Children {
			childrenView += "  " + child.View() + "\n"
		}
	return indicator + " " + m.Title + "\n" + childrenView
}
