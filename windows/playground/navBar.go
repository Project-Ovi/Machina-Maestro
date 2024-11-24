package playground

import (
	"log"
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

func settingsWindow() {
	// Make window
	App := fyne.CurrentApp()
	window := App.NewWindow("Settings")

	// Make title
	titleObj := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Make a container for all the forms
	forms := container.New(
		layout.NewFormLayout(),
	)

	// Make name form
	nameTitle := widget.NewLabel("Name:")
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "My OVI"
	nameEntry.Text = thisModel.Name
	forms.Add(nameTitle)
	forms.Add(nameEntry)

	// Make description form
	descriptionTitle := widget.NewLabel("Description:")
	descriptionEntry := widget.NewEntry()
	descriptionEntry.PlaceHolder = "My beloved OVI"
	descriptionEntry.Text = thisModel.Description
	forms.Add(descriptionTitle)
	forms.Add(descriptionEntry)

	// Make other forms
	otherForms := returnModelSpecificForm(thisModel.ProductName)
	for i := 0; i < len(otherForms); i += 2 {
		// Get title
		var title *widget.Label
		if val, ok := otherForms[i].(*widget.Label); ok {
			title = val
		} else {
			continue
		}

		// Get current value
		currentValue := thisModel.Others[title.Text]

		// Populate entry
		if val, ok := otherForms[i+1].(*widget.Entry); ok {
			val.SetText(currentValue)
			val.OnChanged = func(s string) {}
		}

		// Populate select
		if val, ok := otherForms[i+1].(*widget.Select); ok {
			val.SetSelected(currentValue)
			val.OnChanged = func(s string) {}
		}

		// Append to the form
		forms.Add(otherForms[i])
		forms.Add(otherForms[i+1])
	}

	// Make buttons
	cancelBTN := widget.NewButtonWithIcon("Cancel", theme.Icon(theme.IconNameCancel), func() { window.Close() })
	saveBTN := widget.NewButtonWithIcon("Save", theme.Icon(theme.IconNameDocumentSave), func() {})
	saveBTN.OnTapped = func() {
		// Disable buttons
		saveBTN.Disable()
		cancelBTN.Disable()

		// Save settings
		for i := 0; i < len(forms.Objects); i += 2 {
			// Save name
			if i == 0 {
				nameEntry := forms.Objects[i+1].(*widget.Entry)
				thisModel.Name = nameEntry.Text
				if thisModel.Name == "" {
					thisModel.Name = nameEntry.PlaceHolder
				}
				continue
			}

			// Save description
			if i == 2 {
				descriptionEntry := forms.Objects[i+1].(*widget.Entry)
				thisModel.Description = descriptionEntry.Text
				if thisModel.Description == "" {
					thisModel.Description = descriptionEntry.PlaceHolder
				}
				continue
			}

			// If code execution got this far, it means we are in the model specific form section

			// Get title
			var title string
			if val, ok := forms.Objects[i].(*widget.Label); ok {
				title = val.Text
			} else {
				continue
			}

			// Check entry
			if val, ok := forms.Objects[i+1].(*widget.Entry); ok {
				err := val.Validate()
				if err != nil {
					log.Println("Failed to validate form:", err)
					return
				}

				thisModel.Others[title] = val.Text
			}

			// Check select
			if val, ok := forms.Objects[i+1].(*widget.Select); ok {
				if val.Selected == "" {
					log.Println("Empty select. Unable to save", title)
					return
				}

				thisModel.Others[title] = val.Selected
			}
		}

		// Save model
		//TODO

		// Close window
		window.Close()
	}

	// Display window
	content := container.New(
		layout.NewVBoxLayout(),
		titleObj,
		layout.NewSpacer(),
		forms,
		layout.NewSpacer(),
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			cancelBTN,
			layout.NewSpacer(),
			saveBTN,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	window.SetContent(content)
	window.Show()
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
	settingsBTN := widget.NewButtonWithIcon("Settings", theme.Icon(theme.IconNameSettings), settingsWindow)

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
