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
	userClosed := true
	window.SetOnClosed(func() {
		if len(btnsData) > 0 && userClosed {
			btnsData[0].F()
		}
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
		// Make the function for this buton
		function := func() {
			// Convert title
			window.SetTitle("Loading...")
			titleObj.Text = window.Title()

			// Make a loading bar
			loadingBar := widget.NewProgressBarInfinite()
			loadingBar.Start()

			// Make content for the loading
			content := container.New(
				layout.NewVBoxLayout(),
				layout.NewSpacer(),
				titleObj,
				loadingBar,
				layout.NewSpacer(),
			)

			// Set window content
			window.SetContent(content)

			// Run function dedicated to this button
			val.F()

			// Close window
			userClosed = false
			window.Close()
		}

		var btn fyne.CanvasObject
		if val.Icon == nil {
			btn = widget.NewButton(val.Text, function)
		} else {
			btn = widget.NewButtonWithIcon(val.Text, val.Icon, function)
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
