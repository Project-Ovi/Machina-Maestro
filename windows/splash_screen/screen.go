package splashscreen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func Launch() {
	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		w := drv.CreateSplashWindow()

		img := getBGIMG()

		w.SetContent(img)
		w.Show()
	}
}
