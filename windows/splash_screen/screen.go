package splashscreen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func Launch(App fyne.App) {
	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		w := drv.CreateSplashWindow()
		w.SetContent(widget.NewLabel("Hello World"))
		w.Show()
	}
}
