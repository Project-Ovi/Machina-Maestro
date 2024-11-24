package playground

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func productPage(content *fyne.Container) {
	//! Make loading page
	// Make title
	titleObj := canvas.NewText("Loading...", theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Make subtitle
	subtitleObj := canvas.NewText("Please be patient while the requested page is being downloaded...", theme.Color(theme.ColorNameForeground))
	subtitleObj.Alignment = fyne.TextAlignCenter

	// Make loading bar
	loadingBar := widget.NewProgressBarInfinite()
	loadingBar.Start()

	// Display the loading
	content.Layout = layout.NewVBoxLayout()
	content.Objects = []fyne.CanvasObject{
		layout.NewSpacer(),
		titleObj,
		subtitleObj,
		loadingBar,
		layout.NewSpacer(),
	}
}
