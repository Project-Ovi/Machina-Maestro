package playground

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func navbar() *fyne.Container {
	// Add sidebar toggle
	sidebarBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameList), func() {})

	// Add start button
	startBTN := widget.NewButtonWithIcon("Start", theme.Icon(theme.IconNameMediaPlay), func() {})

	// Add stop button
	stopBTN := widget.NewButtonWithIcon("Stop", theme.Icon(theme.IconNameMediaStop), func() {})

	// Add tools button
	toolsBTN := widget.NewButtonWithIcon("Tools", theme.Icon(theme.IconNameComputer), func() {})

	// Add units button
	unitsBTN := widget.NewButtonWithIcon(thisMeasurementSystem.Name(), theme.Icon(theme.IconNameSettings), func() {})

	// Add settings button
	settingsBTN := widget.NewButtonWithIcon("Settings", theme.Icon(theme.IconNameSettings), func() {})

	// Add help button
	helpBTN := widget.NewButtonWithIcon("Help", theme.Icon(theme.IconNameHelp), func() {})

	// Add exit button
	exitBTN := widget.NewButtonWithIcon("Exit", theme.Icon(theme.IconNameCancel), func() {})

	// Make button group
	BTNGroup := container.New(
		layout.NewHBoxLayout(),
		sidebarBTN,
		startBTN,
		stopBTN,
		toolsBTN,
		layout.NewSpacer(),
		unitsBTN,
		settingsBTN,
		helpBTN,
		exitBTN,
	)

	// Add a background
	content := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		BTNGroup,
	)

	return content
}
