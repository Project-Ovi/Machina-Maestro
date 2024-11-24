package playground

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

func errorPage(content *fyne.Container, err error) {
	// Make  title
	titleObj := canvas.NewText("Something went wrong", theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Make subtitle
	subtitleObj := canvas.NewText(err.Error(), theme.Color(theme.ColorNameForeground))
	subtitleObj.Alignment = fyne.TextAlignCenter

	// Get image
	imgObj := canvas.NewImageFromResource(theme.Icon(theme.IconNameWarning))
	imgObj.FillMode = canvas.ImageFillContain

	// Make content
	content.Layout = layout.NewVBoxLayout()
	content.Objects = []fyne.CanvasObject{
		titleObj,
		subtitleObj,
		imgObj,
		layout.NewSpacer(),
	}
}
