package playground

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

var thisModel helper.OVI

var toReturn string
var playgroundExitWaitGroup sync.WaitGroup
var playgroundMainContent *fyne.Container

func Launch(window fyne.Window, chosenModelDirName string) string {
	// Load
	playgroundExitWaitGroup.Add(1)
	loadModel(chosenModelDirName)

	// Make content
	playgroundMainContent = container.NewWithoutLayout()
	content := container.New(
		&helper.StackWithNavbar{},
		navbar(),
		container.New(
			&helper.StackWithSidebar{},
			sidebar(),
			playgroundMainContent,
		),
	)

	// Display contents to window
	window.SetContent(content)
	window.Show()

	// Wait
	playgroundExitWaitGroup.Wait()
	return toReturn
}
