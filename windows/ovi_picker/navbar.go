package ovipicker

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func navbar() fyne.CanvasObject {
	// Make background
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground))
	log.Println("Picker navbar: Loaded backgoround")

	// Make back button
	backBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameNavigateBack), func() {})
	backBTN.OnTapped = func() {
		backBTN.Disable()
		returnedText = "home"
		waitPressGroup.Done()
	}
	log.Println("Picker navbar: Loaded back button")

	// Make add btn
	addBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), func() {})
	addBTN.OnTapped = func() {
		addBTN.Disable()
	}
	log.Println("Picker navbar: Loaded add button")

	// Build Layout
	content := container.New(
		layout.NewHBoxLayout(),
		backBTN,
		layout.NewSpacer(),
		addBTN,
	)
	log.Println("Picker navbar: Built layout")

	// Finish the layout
	return container.New(
		layout.NewStackLayout(),
		bg,
		content,
	)
}
