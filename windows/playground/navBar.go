package playground

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/dialogs"
)

func switchMeasurementSystem() {
	switch thisMeasurementSystem {
	case metricSystem:
		thisMeasurementSystem = imperialSystem
	case imperialSystem:
		thisMeasurementSystem = SISystem
	case SISystem:
		thisMeasurementSystem = metricSystem
	}
}

func helpWindow() {
	// Make window
	App := fyne.CurrentApp()
	window := App.NewWindow("Help")

	// Get logo
	logo := getLogo()
	logo.SetMinSize(fyne.NewSquareSize(150))

	// Make a title
	title := canvas.NewText("Machina Maestro", theme.Color(theme.ColorNameForeground))
	title.TextSize = theme.Size(theme.SizeNameHeadingText)
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	// Make a subtitle
	subtitle := canvas.NewText("The unified control interface for all of your OVIs", theme.Color(theme.ColorNameForeground))
	subtitle.TextSize = theme.Size(theme.SizeNameSubHeadingText)
	subtitle.Alignment = fyne.TextAlignCenter
	subtitle.TextStyle.Bold = false

	// Make a github button
	githubBTN := widget.NewButtonWithIcon("Github", theme.Icon(theme.IconNameComputer), func() {
		githubURL, _ := url.Parse("https://github.com/Project-Ovi/Machina-Maestro")
		App.OpenURL(githubURL)
	})

	// Make the content
	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			logo,
			container.New(
				layout.NewVBoxLayout(),
				layout.NewSpacer(),
				title,
				subtitle,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		),
		githubBTN,
		layout.NewSpacer(),
	)

	// Display window
	window.SetContent(content)
	window.Resize(fyne.NewSize(window.Canvas().Size().Width*1.1, 0))
	window.SetFixedSize(true)
	window.Show()
}

func navbar() *fyne.Container {
	var content *fyne.Container

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
	unitsBTN.OnTapped = func() {
		switchMeasurementSystem()
		unitsBTN.SetText(thisMeasurementSystem.Name())
	}

	// Add settings button
	settingsBTN := widget.NewButtonWithIcon("Settings", theme.Icon(theme.IconNameSettings), func() {})

	// Add help button
	helpBTN := widget.NewButtonWithIcon("Help", theme.Icon(theme.IconNameHelp), helpWindow)

	// Add exit button
	exitBTN := widget.NewButtonWithIcon("Exit", theme.Icon(theme.IconNameCancel), func() {
		dialogs.ButtonDialog("Exit", "Are you sure you want to exit?", []dialogs.ButtonDialogButtons{
			{
				Text: "Cancel",
				Icon: theme.Icon(theme.IconNameCancel),
				F:    func() {},
			},
			{
				Text: "Exit",
				Icon: theme.Icon(theme.IconNameLogout),
				F: func() {
					toReturn = "home"
					playgroundExitWaitGroup.Done()
				},
			},
		})
	})

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
	content = container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		BTNGroup,
	)

	return content
}
