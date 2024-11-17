package ovipicker

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

func modelAddWindow() {
	// Create window
	App := fyne.CurrentApp()
	window := App.NewWindow("Add model")

	// Create title
	title := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	title.TextSize = theme.Size(theme.SizeNameHeadingText)
	title.TextStyle.Bold = true

	// Create name form
	nameTitle := widget.NewLabel("Name:")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("My OVI")

	// Create form
	form := container.New(
		layout.NewFormLayout(),
		nameTitle, nameEntry,
	)

	// Make content
	content := container.New(
		&helper.StackWithNavbar{},
		title,
		form,
	)

	// Set window content
	window.SetContent(content)
	window.Show()
}
