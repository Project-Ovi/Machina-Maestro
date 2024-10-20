package fatalerror

import (
	"bytes"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

func Show(err error, logger bytes.Buffer, MainWindow fyne.Window, App fyne.App) {
	App.SendNotification(fyne.NewNotification("Crash report", "Machina-Maestro ran into a problem and crashed"))
	const popupSize = 512

	// Create window
	errorPopup := App.NewWindow(err.Error())
	errorPopup.SetOnClosed(func() { log.Fatal(err) })
	MainWindow.SetOnClosed(func() { log.Fatal(err) })
	errorPopup.SetFixedSize(true)
	errorPopup.Hide()

	// Make logging widget
	logWidget := widget.NewRichTextWithText(logger.String())
	// logWidget.ShowLineNumbers = true

	// Make titles
	title := canvas.NewText("A fatal error has occured:", colornames.White)
	subtitle := canvas.NewText(err.Error(), colornames.Orange)

	// Wrap everything into containers
	cntScrolling := container.NewScroll(logWidget)
	wincnt := container.NewWithoutLayout(title, subtitle, cntScrolling)
	errorPopup.SetContent(wincnt)
	cntScrolling.ScrollToBottom()

	// Move title into the view
	title.TextSize = 25
	title.Alignment = fyne.TextAlignCenter
	title.Move(fyne.NewPos(popupSize/2, 0))
	title.Refresh()

	// Move subtitle into view
	subtitle.Alignment = fyne.TextAlignCenter
	subtitle.TextSize = 1
	for i := 1; i < 30; i++ {
		if fyne.MeasureText(err.Error(), subtitle.TextSize, subtitle.TextStyle).Width < float32(popupSize) {
			subtitle.TextSize = float32(i)
		} else {
			break
		}
	}
	subtitle.Move(fyne.NewPos(popupSize/2, fyne.MeasureText(title.Text, title.TextSize, title.TextStyle).Height))
	subtitle.Refresh()

	// Move scrollbar into view
	cntScrolling.Move(fyne.NewPos(popupSize*.1/2, 2*(fyne.MeasureText(title.Text, title.TextSize, title.TextStyle).Height+fyne.MeasureText(err.Error(), subtitle.TextSize, subtitle.TextStyle).Height)))
	cntScrolling.Resize(fyne.NewSize(popupSize*.9, popupSize/2))

	// Refresh
	title.Refresh()
	subtitle.Refresh()
	errorPopup.Resize(fyne.NewSquareSize(popupSize))
	errorPopup.Show()
}
