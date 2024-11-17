package splashscreen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func Launch(loadFunc func()) {
	var w fyne.Window
	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		// Create window
		w = drv.CreateSplashWindow()

		// Load image
		img := getBGIMG()

		// Make a text that will output the logger's status
		descriptiveText := widget.NewLabel("Loading...")

		// Create content
		content := container.New(
			layout.NewBorderLayout(nil, descriptiveText, nil, nil),
			img,
		)

		// Set content
		w.SetContent(content)
		w.Show()
	}

	// Load and exit
	loadFunc()
	w.Hide()
}
