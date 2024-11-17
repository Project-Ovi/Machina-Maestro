package splashscreen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
)

func Launch() {
	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		w := drv.CreateSplashWindow()

		img := getBGIMG()

		content := container.New(
			layout.NewBorderLayout(nil, nil, nil, nil),
			img,
		)

		w.SetContent(content)
		w.Show()
	}
}
