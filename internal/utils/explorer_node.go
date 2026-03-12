package utils

type ExplorerNode struct {
	Title    string
	Children []ExplorerNode
	Expanded bool
	Parent   *ExplorerNode
	Type     string
	Next     *ExplorerNode
	Previous *ExplorerNode
}

type ExplorerList struct {
	Root     *ExplorerNode
	Selected *ExplorerNode
}

func (l *ExplorerList) Initialize() {
	l.Root = &ExplorerNode{Title: "", Type: "root"}
	l.Selected = l.Root
}

func (l *ExplorerList) Expand(children []ExplorerNode) {
	l.Selected.Children = children
}

func (l *ExplorerList) Contract() {
	l.Selected.Children = nil
}

func (l *ExplorerList) MoveDown() {
	next := l.Selected.Next
	l.Selected = next
}

func (l *ExplorerList) MoveUp() {
	previous := l.Selected.Previous
	l.Selected = previous
}
