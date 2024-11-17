package splashscreen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func Launch(App fyne.App) {
	window := App.NewWindow("Launching...")

	window.SetContent(widget.NewLabel("Hello, World!"))

	window.Show()
}
