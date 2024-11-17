package ovipicker

import (
	"fmt"

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
	window.SetOnClosed(func() {
		returnedText = "picker"
		waitPressGroup.Done()
	})

	// Create title
	title := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	title.TextSize = theme.Size(theme.SizeNameHeadingText)
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	// Create name form
	nameTitle := widget.NewLabel("Name:")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("My OVI")

	// Create description form
	descriptionTitle := widget.NewLabel("Description:")
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("My beloved OVI")

	// Create product form
	productTitle := widget.NewLabel("Product:")
	productEntry := widget.NewSelect([]string{"Option 1", "Option 2", "Option 3"}, func(s string) {
		fmt.Println("Changed to", s)
	})

	// Create form
	form := container.New(
		layout.NewFormLayout(),
		nameTitle, nameEntry,
		descriptionTitle, descriptionEntry,
		productTitle, productEntry,
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
