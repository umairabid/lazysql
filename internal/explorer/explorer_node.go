package explorer

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ExplorerNodeModel struct {
	Title string
	Children []ExplorerNodeModel
	Selected bool
	Expanded bool
}

func (m ExplorerNodeModel) Init() tea.Cmd {
	return nil
}

func (m ExplorerNodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ExplorerNodeModel) View() string {
	indicator := " "
	if m.Selected {
		indicator = ">"
	}
	return indicator + " " + m.Title
}
