package playground

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
	"github.com/Project-Ovi/Machina-Maestro/windows/fatalerror"
)

func HelpWindow() {
	// Make window
	window := App.NewWindow("Help")

	// Add logo
	ICON_robot_hand_round := helper.ICON_robot_hand_round
	ICON_robot_hand_round.SetMinSize(fyne.NewSquareSize(100))

	// Add title
	title := canvas.NewText("Machina Maestro", theme.Color(theme.ColorNameForeground))
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 30

	// Add subtitle
	subtitle := widget.NewLabel("A comprehensive UI application for controlling OVI products.")
	subtitle.Wrapping = fyne.TextWrapWord

	// Add github button
	btn := widget.NewButtonWithIcon("Github", theme.Icon(theme.IconNameComputer), func() {
		url, err := url.Parse("https://github.com/Project-Ovi/Machina-Maestro")
		if err != nil {
			fatalerror.Show(err, logger, MainWindow, App)
		}
		App.OpenURL(url)
	})

	// Assemble UI
	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			ICON_robot_hand_round,
			container.New(
				layout.NewVBoxLayout(),
				title,
				subtitle,
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
		btn,
	)

	// Show window
	window.SetContent(content)
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(500, 200))
	window.Show()
}
