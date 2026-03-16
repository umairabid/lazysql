package utils

import (
	"fmt"
	"testing"
)

func TestExplorerListInitialize(t *testing.T) {
	list := ExplorerList{}
	list.Initialize()
	t.Run("Root node should be initialized correctly", func(t *testing.T) {
		if list.Root == nil {
			t.Fatal("Expected Root to be initialized, got nil")
		}
		if list.Root.Title != "Root" {
			t.Errorf("Expected Root Title to be Root, got '%s'", list.Root.Title)
		}
		if list.Root.Type != "root" {
			t.Errorf("Expected Root Type to be 'root', got '%s'", list.Root.Type)
		}
	})
}

func TestExplorerListExpandNonRoot(t *testing.T) {
	list := ExplorerList{}
	list.Initialize()

	children := []ExplorerNode{
		{Title: "Child 1", Type: "child"},
		{Title: "Child 2", Type: "child"},
		{Title: "Child 3", Type: "child"},
	}
	list.Expand(children)

	t.Run("Should set children correctly", func(t *testing.T) {
		if len(list.Root.Children) != 3 {
			t.Fatalf("Expected 3 children, got %d", len(list.Selected.Children))
		}
	})

	t.Run("Should set selected as first node previous", func(t *testing.T) {
		if list.Selected.Previous != list.Root {
			t.Errorf("Expected first child Previous to be selected node, got %v", list.Selected.Previous)
		}
	})

	t.Run("Should set second child as next of first child", func(t *testing.T) {
		if list.Selected.Next != &list.Root.Children[1] {
			t.Errorf("Expected first child Next to be second child, got %v", list.Root.Children[0].Next)
		}
	})

	t.Run("Should set next of last child to next of selected node", func(t *testing.T) {
		if list.Root.Children[2].Next != list.Root {
			t.Errorf("Expected last child Next to be selected root, got %v", list.Root.Children[2].Next)
		}
	})

	t.Run("Should set previous of first child to selected node", func(t *testing.T) {
		if list.Root.Children[0].Previous != list.Root {
			t.Errorf("Expected first child Previous to be root node, got %v", list.Root.Children[0].Previous)
		}
	})

	t.Run("Should set next of selected node to first child", func(t *testing.T) {
		if list.Root.Next != &list.Root.Children[0] {
			t.Errorf("Expected selected node Next to be first child, got %v", list.Root.Next)
		}
	})

	t.Run("Should set previous of next node to last child", func(t *testing.T) {
		if list.Root.Next.Next.Next == list.Root {
			t.Errorf("Expected next node Previous to be root, got %v", list.Root.Next.Next.Next)
		}
	})
}

func TestExplorerListExpandRoot(t *testing.T) {
	list := ExplorerList{}
	list.Initialize()

	children := []ExplorerNode{
		{Title: "Child 1", Type: "child"},
		{Title: "Child 2", Type: "child"},
	}

	list.Expand(children)

	t.Run("Should set first child as selected", func(t *testing.T) {
		if &list.Root.Children[0] != list.Selected {
			t.Errorf("Expected last child Next to wrap back to root, got %v", list.Root.Children[1].Next)
		}
	})
}

func TestExplorerListCollapse(t *testing.T) {
	list := ExplorerList{}
	list.Initialize()

	rootChildren := []ExplorerNode{
		{Title: "Parent 1", Type: "child"},
		{Title: "Parent 2", Type: "child"},
	}

	paren1Children := []ExplorerNode{
		{Title: "Child 1", Type: "child"},
		{Title: "Child 2", Type: "child"},
	}

	list.Expand(rootChildren)
	list.Expand(paren1Children)

	t.Run("Should remove children when collapsed", func(t *testing.T) {
		list.Contract()
		if len(list.Selected.Children) != 0 {
			t.Fatalf("Expected no children after collapse, got %d", len(list.Selected.Children))
		}
	})

	t.Run("Should make selected node not expanded", func(t *testing.T) {
		if list.Selected.Expanded {
			t.Errorf("Expected selected node to be not expanded, got %v", list.Selected.Expanded)
		}
	})

	t.Run("Should change previous of node after last child to selected node", func(t *testing.T) {
		if rootChildren[1].Previous != list.Selected {
			t.Errorf("Expected previous of node after last child to be selected node, got %v", rootChildren[1].Previous)
		}
	})

	t.Run("Should next of selected node to sibling", func(t *testing.T) {
		if list.Selected.Next != &rootChildren[1] {
			t.Errorf("Expected next of selected node to be sibling, got %v", list.Selected.Next)
		}
	})

	t.Run("Should not allow collapsing root node", func(t *testing.T) {
		list.Selected = list.Root
		list.Contract()
		if len(list.Selected.Children) != 2 {
			t.Fatalf("Expected root node to still have children after collapse, got %d", len(list.Selected.Children))
		}
	})

	t.Run("Should be a no-op on non-expanded node", func(t *testing.T) {
		list.Selected = &rootChildren[1]
		list.Contract()
		if list.Selected.Expanded {
			t.Errorf("Expected non-expanded node to remain not expanded after Contract")
		}
		if len(list.Selected.Children) != 0 {
			t.Errorf("Expected non-expanded node to have no children after Contract, got %d", len(list.Selected.Children))
		}
	})
}

func TestExplorerListLastDescendant(t *testing.T) {
	list := ExplorerList{}
	list.Initialize()

	rootChildren := []ExplorerNode{
		{Title: "Parent 1", Type: "child"},
		{Title: "Parent 2", Type: "child"},
		{Title: "Parent 3", Type: "child"},
	}

	list.Expand(rootChildren)
	list.MoveDown()
	list.MoveDown()

	fmt.Println("Selected node:", list.Selected.Title)

	list.MoveDown()
	fmt.Println("Selected node after moving down:", list.Selected.Title)

	parentChildren := []ExplorerNode{
		{Title: "Child 1", Type: "child"},
		{Title: "Child 2", Type: "child"},
		{Title: "Child 3", Type: "child"},
	}

	list.Expand(parentChildren)
	list.MoveDown()
	list.MoveDown()
	list.MoveDown()
	fmt.Println("Selected node after moving down through children:", list.Selected.Title)
}
