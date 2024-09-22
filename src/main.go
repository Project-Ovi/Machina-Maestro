package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

var App fyne.App
var MainWindow fyne.Window

func initialize() {
	//statusTxt := "Loading application"

	// Render a splash screen while the app is loading
	var splashScreen fyne.Window
	var status *canvas.Text
	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		// Create splash screen
		splashScreen = drv.CreateSplashWindow()
		splashScreen.RequestFocus()

		// Draw the background
		background := canvas.NewImageFromFile("splash.png")
		background.Resize(fyne.NewSize(1024, 512)) // Resize image to desired size
		background.Move(fyne.NewPos(0, 0))         // Set the position for the image

		// Draw text
		status = canvas.NewText("Loading application...", colornames.Orangered)
		status.TextSize = 20
		status.Move(fyne.NewPos(10, 512-35))

		// Draw rectangle
		rectangle := canvas.NewRectangle(color.RGBA{R: 9, G: 9, B: 9, A: 255})
		rectangle.Resize(fyne.NewSize(1024, 40))
		rectangle.Move(fyne.NewPos(0, 512-40))

		// Assemble UI
		content := container.NewWithoutLayout(background, rectangle, status)

		// Show splash screen
		splashScreen.SetContent(content)
		splashScreen.Resize(fyne.NewSize(1024, 512))
		splashScreen.Show()
		MainWindow.Hide()
	}

	// Some init
	i := 0
	for range time.Tick(time.Millisecond * 50) {
		status.Text = "Loading... " + fmt.Sprint(i) + "%"
		status.Refresh()
		if i == 100 {
			break
		}
		i += 1
	}

	// Hide splash screen
	status.Text = "Finished loading!"
	status.Refresh()
	time.Sleep(time.Second * 2)
	splashScreen.Hide()
	MainWindow.Show()

}

func main() {

	// Start fyne
	App = app.New()
	MainWindow = App.NewWindow("Hello World")
	MainWindow.SetOnClosed(func() {
		os.Exit(0)
	})

	// Render app
	go func() {
		initialize()

		clock := widget.NewLabel("")
		MainWindow.SetContent(clock)
		for range time.Tick(time.Second) {
			clock.SetText(time.Now().Format("Time: 03:04:05"))
		}
	}()

	// Display app
	MainWindow.ShowAndRun()
}
