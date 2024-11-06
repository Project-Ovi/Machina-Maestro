package modelselect

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
	"github.com/Project-Ovi/Machina-Maestro/windows/playground"
	"golang.org/x/image/colornames"
)

func modelAddWindow(window fyne.Window, btnAdd *widget.Button) {
	// Make a new window
	addWindow := App.NewWindow("Add a new model")
	addWindow.SetOnClosed(func() { btnAdd.Enable() })

	// Make a navbar
	btnBack := widget.NewButton("", func() { addWindow.Close() })
	btnBack.SetIcon(theme.Icon(theme.IconNameNavigateBack))
	navbar := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameButton)),
		container.New(
			layout.NewHBoxLayout(),
			btnBack,
			layout.NewSpacer(),
		),
	)

	// Add a title
	title := canvas.NewText("Add a new model", theme.Color(theme.ColorNameForeground))
	title.TextSize = 32
	title.Alignment = fyne.TextAlignCenter

	// Add a subtitle
	subtitle := canvas.NewText("", colornames.Red)
	subtitle.TextSize = 28
	subtitle.Alignment = fyne.TextAlignCenter

	// Combine title and subtitle
	top := container.New(layout.NewVBoxLayout(), title, subtitle)

	// Add name field
	nameTitle := canvas.NewText("Name", theme.Color(theme.ColorNameForeground))
	nameOption := widget.NewEntry()
	nameOption.Text = "My OVI"

	// Add model select button
	selectTitle := canvas.NewText("Model", theme.Color(theme.ColorNameForeground))
	selectOption := widget.NewSelect(playground.SupportedModels, func(s string) {})

	// Make the form
	buildForm := func(form *fyne.Container) {
		form.Objects = append(form.Objects,
			nameTitle, nameOption,
			selectTitle, selectOption,
			widget.NewSeparator(), widget.NewSeparator(),
		)
	}
	form := container.New(layout.NewFormLayout())
	buildForm(form)

	// Add a submit button
	btn := helper.BigBTN("Submit", 60, 15, func() {
		subtitle.Text = playground.SaveSelector(form)
		subtitle.Refresh()
		if subtitle.Text == "" {
			addWindow.Close()
			playground.Launch(window)
		}
	})

	selectOption.OnChanged = func(s string) {
		form.RemoveAll()
		buildForm(form)
		playground.FormSelector(s, form)
	}

	content := container.New(
		layout.NewVBoxLayout(),
		navbar,
		top,
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		btn,
	)
	addWindow.SetContent(content)
	addWindow.Resize(fyne.NewSquareSize(512))
	addWindow.Show()
}
