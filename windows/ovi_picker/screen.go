package ovipicker

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var waitPressGroup sync.WaitGroup
var returnedText string

func Launch(window fyne.Window) string {
	// Initialize wait group
	waitPressGroup.Add(1)

	// Make the navbar
	navbar := navbar()

	// Make content
	content := container.New(
		layout.NewVBoxLayout(),
		navbar,
		layout.NewSpacer(),
	)

	// Set content
	window.SetContent(content)
	window.Show()

	waitPressGroup.Wait()
	return returnedText
}
