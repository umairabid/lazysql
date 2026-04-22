package utils

type ViewerTableData [][]string
type ActiveViewChanged string
type LayoutUpdated ConnectionContainerLayout

type ConnectionManagerLayout struct {
	ScreenWidth         int
	ScreenHeight        int
	WinHeight           int
	WinWidth            int
	HeaderHeight        int
	BodyHeight          int
	ConnectionListWidth int
	ConnectionFormWidth int
	FooterHeight        int
	HelpWidth           int
	HelpHeight          int
}

type ConnectionContainerLayout struct {
	ScreenWidth    int
	ScreenHeight   int
	EditorWidth    int
	EditorHeight   int
	ViewerWidth    int
	ViewerHeight   int
	ExplorerWidth  int
	ExplorerHeight int
}
