package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ButtonDialogButtons struct {
	Text string
	Icon fyne.Resource
	F    func()
}

func ButtonDialog(title string, subtitle string, btnsData []ButtonDialogButtons) {
	// Get app
	App := fyne.CurrentApp()

	// Make a new window
	window := App.NewWindow(title)
	window.SetOnClosed(func() {
		btnsData[0].F()
	})

	// Make title
	titleObj := canvas.NewText(title, theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Make subtitle
	subtitleObj := canvas.NewText(subtitle, theme.Color(theme.ColorNameForeground))
	subtitleObj.TextSize = theme.Size(theme.SizeNameSubHeadingText)
	subtitleObj.Alignment = fyne.TextAlignCenter
	subtitleObj.TextStyle.Bold = false

	// Make buttons
	BTNContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer())
	for _, val := range btnsData {
		var btn fyne.CanvasObject
		if val.Icon == nil {
			btn = widget.NewButton(val.Text, val.F)
		} else {
			btn = widget.NewButtonWithIcon(val.Text, val.Icon, val.F)
		}
		BTNContainer.Add(btn)
		BTNContainer.Add(layout.NewSpacer())
	}

	// Make content
	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		titleObj,
		subtitleObj,
		layout.NewSpacer(),
		BTNContainer,
		layout.NewSpacer(),
	)

	// Set window content
	window.SetContent(content)
	window.Show()
}
