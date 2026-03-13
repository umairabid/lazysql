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
	l.Selected.Expanded = true

	for i := range l.Selected.Children {
		child := &l.Selected.Children[i]
		child.Parent = l.Selected

		if i == 0 {
			child.Previous = l.Selected
		} else {
			child.Previous = &l.Selected.Children[i-1]
			l.Selected.Children[i-1].Next = child
		}
	}

	if len(children) > 0 {
		l.Selected.Children[len(children)-1].Next = l.Selected.Next
	}
	l.Selected = &l.Selected.Children[0]
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
