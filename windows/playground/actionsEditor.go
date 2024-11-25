package playground

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

func actionsEditor(content *fyne.Container, act *helper.Action) {
	// Make navbar
	editorNavbar := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameNavigateBack), func() { actionsPage(content) }),
			layout.NewSpacer(),
		),
	)

	// Make title
	titleObj := canvas.NewText(act.Name, theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Make subtitle
	subtitleObj := canvas.NewText(act.Description, theme.Color(theme.ColorNameForeground))
	subtitleObj.Alignment = fyne.TextAlignCenter

	// Set content
	content.Layout = &helper.StackWithNavbar{}
	content.Objects = []fyne.CanvasObject{
		container.New(
			layout.NewVBoxLayout(),
			editorNavbar,
			titleObj,
			subtitleObj,
		),
	}

}
