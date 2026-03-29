package utils

func CalculateConnectionManagerLayout(width int, height int) ConnectionManagerLayout {
	headerHeight := 10
	footerHeight := 10

	listWidth := width / 3
	formWidth := width - listWidth

	return ConnectionManagerLayout{
		WinWidth:            width,
		WinHeight:           height,
		HeaderHeight:        headerHeight,
		BodyHeight:          height - (headerHeight + footerHeight),
		ConnectionListWidth: listWidth,
		ConnectionFormWidth: formWidth,
		FooterHeight:        footerHeight,
	}
}
