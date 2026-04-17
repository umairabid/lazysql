package utils

import (
	"slices"
)

var CONNECTION_MANAGER_MIN_WIDTH = 80
var CONNECTION_MANAGER_MIN_HEIGHT = 24
var EXPLORER_MIN_WIDTH = 20

func CalculateConnectionManagerLayout(width int, height int) ConnectionManagerLayout {
	headerHeight := 3
	footerHeight := 3

	widths := []int{CONNECTION_MANAGER_MIN_WIDTH, width / 3}
	heights := []int{CONNECTION_MANAGER_MIN_HEIGHT, height / 3}
	winWidth := slices.Max(widths)
	winHeight := slices.Max(heights)
	listWidth := winWidth / 3
	formWidth := winWidth - listWidth

	return ConnectionManagerLayout{
		ScreenWidth:         width,
		ScreenHeight:        height,
		WinWidth:            winWidth,
		WinHeight:           winHeight,
		HeaderHeight:        headerHeight,
		BodyHeight:          winHeight - (headerHeight + footerHeight),
		ConnectionListWidth: listWidth,
		ConnectionFormWidth: formWidth,
		FooterHeight:        footerHeight,
	}
}

func CalculateConnectionContainerLayout(width int, height int) ConnectionContainerLayout {
	explorerWidths := []int{EXPLORER_MIN_WIDTH, width / 4}
	explorerWidth := slices.Max(explorerWidths)

	editorWidth := width - explorerWidth
	viewerWidth := editorWidth

	explorerHeight := height
	editorHeight := height / 3
	viewerHeight := height - editorHeight

	return ConnectionContainerLayout{
		ScreenWidth:    width,
		ScreenHeight:   height,
		ExplorerWidth:  explorerWidth,
		ExplorerHeight: explorerHeight,
		EditorWidth:    editorWidth,
		EditorHeight:   editorHeight,
		ViewerWidth:    viewerWidth,
		ViewerHeight:   viewerHeight,
	}
}
