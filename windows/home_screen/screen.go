package homescreen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

func Launch(window fyne.Window) {
	// Make title
	title := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 35
	title.TextStyle.Bold = true

	// Make Logo
	logo := getLogo()
	logo.SetMinSize(fyne.NewSquareSize(120))

	// Make content
	content := container.New(
		layout.NewVBoxLayout(),
		title,
		layout.NewSpacer(),
		logo,
		layout.NewSpacer(),
	)

	// Show window
	window.SetContent(content)
	window.Show()
}
