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
