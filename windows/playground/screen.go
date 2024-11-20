package playground

import (
	"fyne.io/fyne/v2"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

var thisModel helper.OVI

func Launch(MainWindow fyne.Window, chosenModelDirName string) {
	// Load model
	loadModel(chosenModelDirName)
}
