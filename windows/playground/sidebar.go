package playground

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func sidebar() *fyne.Container {
	// Make home button
	homeBTN := widget.NewButtonWithIcon("Home", theme.Icon(theme.IconNameHome), func() {})
	homeBTN.Alignment = widget.ButtonAlignLeading

	// Make overview button
	overviewBTN := widget.NewButtonWithIcon("Overview", theme.Icon(theme.IconNameVisibility), func() {})
	overviewBTN.Alignment = widget.ButtonAlignLeading

	// Make product button
	productBTN := widget.NewButtonWithIcon("Product", theme.Icon(theme.IconNameInfo), func() {})
	productBTN.Alignment = widget.ButtonAlignLeading

	// Make actions button
	actionsBTN := widget.NewButtonWithIcon("Actions", theme.Icon(theme.IconNameMediaPlay), func() {})
	actionsBTN.Alignment = widget.ButtonAlignLeading

	// Make scripting button
	scriptingBTN := widget.NewButtonWithIcon("Scripting", theme.Icon(theme.IconNameFile), func() {})
	scriptingBTN.Alignment = widget.ButtonAlignLeading

	// Make routines button
	routinesBTN := widget.NewButtonWithIcon("Routines", theme.Icon(theme.IconNameMediaReplay), func() {})
	routinesBTN.Alignment = widget.ButtonAlignLeading

	// Build final sidebar
	content := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		container.NewVScroll(container.New(
			layout.NewVBoxLayout(),
			homeBTN,
			overviewBTN,
			productBTN,
			actionsBTN,
			scriptingBTN,
			routinesBTN,
			layout.NewSpacer(),
		)),
	)
	return content
}
