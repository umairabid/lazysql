package utils

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestEscapeCell(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello\nworld", "hello\\nworld"},
		{"hello\rworld", "hello\\rworld"},
		{"hello\tworld", "hello\\tworld"},
		{"\n\r\t", "\\n\\r\\t"},
	}

	for _, test := range tests {
		if got := escapeCell(test.input); got != test.expected {
			t.Errorf("escapeCell(%q) = %q; want %q", test.input, got, test.expected)
		}
	}
}

func TestInitTable(t *testing.T) {
	data := [][]string{
		{"ID", "Name", "Age"},
		{"1", "Alice", "30"},
		{"2", "Bob", "25"},
	}
	width, height := 40, 20
	table := InitTable(data, width, height)

	t.Run("Initialization", func(t *testing.T) {
		if len(table.Columns) != 3 {
			t.Errorf("Expected 3 columns, got %d", len(table.Columns))
		}
		if len(table.Rows) != 2 {
			t.Errorf("Expected 2 rows, got %d", len(table.Rows))
		}
		if table.Width != width || table.Height != height {
			t.Errorf("Expected size %dx%d, got %dx%d", width, height, table.Width, table.Height)
		}
		if table.SelectedRow != 0 || table.SelectedColumn != 0 {
			t.Errorf("Expected initial selection (0,0), got (%d,%d)", table.SelectedRow, table.SelectedColumn)
		}
	})

	t.Run("Empty Data", func(t *testing.T) {
		emptyTable := InitTable([][]string{}, width, height)
		if len(emptyTable.Columns) != 0 || len(emptyTable.Rows) != 0 {
			t.Errorf("Expected empty table, got %d cols and %d rows", len(emptyTable.Columns), len(emptyTable.Rows))
		}
		if emptyTable.HasData() {
			t.Error("HasData() should be false for empty table")
		}
	})
}

func TestCalculateColumnWidths(t *testing.T) {
	cols := []string{"ID", "Name"}
	rows := [][]string{
		{"1", "Alice"},
		{"2", "A very long name that exceeds standard width"},
	}

	widths := calculateColumnWidths(cols, rows)

	if len(widths) != 2 {
		t.Fatalf("Expected 2 widths, got %d", len(widths))
	}

	// ID: len("ID")+2 = 4, len("1")+2 = 3. Max is 4.
	if widths[0] != 4 {
		t.Errorf("Expected ID column width 4, got %d", widths[0])
	}

	// Name: len("Name")+2 = 6, len("Alice")+2 = 7, len(long)+2 = many.
	// But it should be capped at MIN_COLUMN_WIDTH (100)
	if widths[1] > 100 {
		t.Errorf("Expected Name column width capped at 100, got %d", widths[1])
	}
}

func TestTableUpdate(t *testing.T) {
	data := [][]string{
		{"Col1", "Col2"},
		{"R1C1", "R1C2"},
		{"R2C1", "R2C2"},
	}
	table := InitTable(data, 40, 20)

	t.Run("Move Down", func(t *testing.T) {
		table, _ = table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		if table.SelectedRow != 1 {
			t.Errorf("Expected selected row 1, got %d", table.SelectedRow)
		}
		// Move down again - should stay at 1 since only 2 rows (index 0, 1)
		table, _ = table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		if table.SelectedRow != 1 {
			t.Errorf("Expected selected row 1 after second 'j', got %d", table.SelectedRow)
		}
	})

	t.Run("Move Up", func(t *testing.T) {
		table.SelectedRow = 1
		table, _ = table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
		if table.SelectedRow != 0 {
			t.Errorf("Expected selected row 0, got %d", table.SelectedRow)
		}
		// Move up again - should stay at 0
		table, _ = table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
		if table.SelectedRow != 0 {
			t.Errorf("Expected selected row 0 after second 'k', got %d", table.SelectedRow)
		}
	})

	t.Run("Move Right", func(t *testing.T) {
		table.SelectedColumn = 0
		table, _ = table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
		if table.SelectedColumn != 1 {
			t.Errorf("Expected selected column 1, got %d", table.SelectedColumn)
		}
		// Move right again - should stay at 1
		table, _ = table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
		if table.SelectedColumn != 1 {
			t.Errorf("Expected selected column 1 after second 'l', got %d", table.SelectedColumn)
		}
	})

	t.Run("Move Left", func(t *testing.T) {
		table.SelectedColumn = 1
		table, _ = table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
		if table.SelectedColumn != 0 {
			t.Errorf("Expected selected column 0, got %d", table.SelectedColumn)
		}
		// Move left again - should stay at 0
		table, _ = table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
		if table.SelectedColumn != 0 {
			t.Errorf("Expected selected column 0 after second 'h', got %d", table.SelectedColumn)
		}
	})
}

func TestTableHasData(t *testing.T) {
	tests := []struct {
		data     [][]string
		expected bool
	}{
		{[][]string{{"Col1"}, {"Row1"}}, true},
		{[][]string{{"Col1"}}, false},
		{[][]string{}, false},
	}

	for i, test := range tests {
		table := InitTable(test.data, 10, 10)
		if got := table.HasData(); got != test.expected {
			t.Errorf("Test case %d: HasData() = %v; want %v", i, got, test.expected)
		}
	}
}
