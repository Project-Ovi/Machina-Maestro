package confirm

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Show(App fyne.App, title string, subtitle string, yes func(), no func()) {
	// Make window
	confirmWindow := App.NewWindow(title)
	confirmWindow.SetFixedSize(true)
	confirmWindow.SetCloseIntercept(func() {
		confirmWindow.RequestFocus()
	})

	// Add title
	titleWidget := canvas.NewText(title, theme.Color(theme.ColorNameForeground))
	titleWidget.TextSize = 25
	titleWidget.Alignment = fyne.TextAlignCenter

	// Add subtitle
	subtitleWidget := canvas.NewText(subtitle, theme.Color(theme.ColorNameForeground))
	subtitleWidget.TextSize = 16
	subtitleWidget.Alignment = fyne.TextAlignCenter

	// Create buttons
	yesBTN := widget.NewButtonWithIcon("YES", theme.Icon(theme.IconNameConfirm), func() {
		yes()
		confirmWindow.Close()
	})
	noBTN := widget.NewButtonWithIcon("NO", theme.Icon(theme.IconNameCancel), func() {
		no()
		confirmWindow.Close()
	})
	buttons := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		yesBTN,
		layout.NewSpacer(),
		noBTN,
		layout.NewSpacer(),
	)

	// Assemble UI
	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		titleWidget,
		subtitleWidget,
		layout.NewSpacer(),
		buttons,
		layout.NewSpacer(),
	)

	confirmWindow.SetContent(content)
	confirmWindow.Resize(fyne.NewSize(500, 200))
	confirmWindow.Show()
}
