package playground

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

var thisModel helper.OVI

func Launch(window fyne.Window, chosenModelDirName string) {
	// Load model
	loadModel(chosenModelDirName)

	// Make content
	content := container.New(
		&helper.StackWithNavbar{},
		navbar(),
	)

	window.SetContent(content)
}
