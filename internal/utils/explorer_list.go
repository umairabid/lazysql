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
	l.Root = &ExplorerNode{Title: "Root", Type: "root", Expanded: true}
	l.Selected = l.Root
	l.Root.Next = l.Root
	l.Root.Previous = l.Root
}

func (l *ExplorerList) Expand(children []ExplorerNode) {
	if len(children) == 0 {
		return
	}

	l.Selected.Children = children
	l.Selected.Expanded = true


	for i := range l.Selected.Children {
		child := &l.Selected.Children[i]
		child.Parent = l.Selected
		isLastChild := i == len(l.Selected.Children)-1

		if i == 0 {
			child.Previous = l.Selected
		} else {
			child.Previous = &l.Selected.Children[i-1]
		}

		if isLastChild {
			child.Next = l.Selected.Next
		} else {
			child.Next = &l.Selected.Children[i+1]
		}
	}

	l.Selected.Next.Previous = &l.Selected.Children[len(l.Selected.Children)-1]
	l.Selected.Next = &l.Selected.Children[0]

	if l.Selected.Type == "root" {
		l.Selected = &l.Selected.Children[0]
	}
}

func (l *ExplorerList) Contract() {
	if !l.Selected.Expanded || l.Selected.Type == "root" {
		return
	}
	l.Selected.Next = l.Selected.Children[len(l.Selected.Children)-1].Next
	l.Selected.Next.Previous = l.Selected
	l.Selected.Expanded = false
	l.Selected.Children = nil
}

func (l *ExplorerList) ContractParent() {
	if l.Selected.Parent == nil || l.Selected.Parent.Type == "root" {
		return
	}
	l.Selected = l.Selected.Parent
	l.Contract()
}

func (l *ExplorerList) MoveDown() {
	next := l.Selected.Next
	if next.Type == "root" {
		next = next.Next
	}
	l.Selected = next
}

func (l *ExplorerList) MoveUp() {
	previous := l.Selected.Previous
	if previous.Type == "root" {
		previous = previous.Previous
	}
	l.Selected = previous
}
