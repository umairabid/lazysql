package utils

import "testing"

func TestExplorerListInitialize(t *testing.T) {
	list := ExplorerList{}
	list.Initialize()
	t.Run("Root node should be initialized correctly", func(t *testing.T) {
		if list.Root == nil {
			t.Fatal("Expected Root to be initialized, got nil")
		}
		if list.Root.Title != "" {
			t.Errorf("Expected Root Title to be empty, got '%s'", list.Root.Title)
		}
		if list.Root.Type != "root" {
			t.Errorf("Expected Root Type to be 'root', got '%s'", list.Root.Type)
		}
	})
}

func TestExplorerListChildren(t *testing.T) {
	list := ExplorerList{}
	list.Initialize()

	children := []ExplorerNode{
		{Title: "Child 1", Type: "child"},
		{Title: "Child 2", Type: "child"},
		{Title: "Child 3", Type: "child"},
	}

	t.Run("Should set children correctly", func(t *testing.T) {
		list.Expand(children)
		if len(list.Selected.Children) != 2 {
			t.Fatalf("Expected 2 children, got %d", len(list.Selected.Children))
		}
	})

	t.Run("Should set selected as first node previous", func(t *testing.T) {
		if list.Selected.Children[0].Previous != list.Selected {
			t.Errorf("Expected first child Previous to be selected node, got %v", list.Selected.Children[0].Previous)
		}
	})

	t.Run("Should set second child as next of first child", func(t *testing.T) {
		if list.Selected.Children[0].Next != &list.Selected.Children[1] {
			t.Errorf("Expected first child Next to be second child, got %v", list.Selected.Children[0].Next)
		}
	})

	t.Run("Should set next of last child to next of selected node", func(t *testing.T) {
		if list.Selected.Children[2].Next != list.Selected.Next {
			t.Errorf("Expected last child Next to be selected node Next, got %v", list.Selected.Children[1].Next)
		}
	})

	t.Run("Should set previous of first child to previous of selected node", func(t *testing.T) {
		if list.Selected.Children[0].Previous != list.Selected.Previous {
			t.Errorf("Expected first child Previous to be selected node Previous, got %v", list.Selected.Children[0].Previous)
		}
	})
}
