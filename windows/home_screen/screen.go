package homescreen

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Launch(window fyne.Window) {
	// Make title
	title := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 35
	title.TextStyle.Bold = true

	// Make Logo
	logo := getLogo()
	logo.SetMinSize(fyne.NewSquareSize(400))

	// Make a button
	var waitPressGroup sync.WaitGroup
	waitPressGroup.Add(1)
	btn := widget.NewButton("Pick your OVI", func() {
		waitPressGroup.Done()
	})

	// Make content
	content := container.New(
		layout.NewVBoxLayout(),
		title,
		layout.NewSpacer(),
		logo,
		layout.NewSpacer(),
		btn,
		layout.NewSpacer(),
	)

	// Show window
	window.SetContent(content)
	window.Show()

	// Wait the button press before exiting
	waitPressGroup.Wait()
}
