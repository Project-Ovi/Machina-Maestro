package startup

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Project-Ovi/Machina-Maestro/windows/fatalerror"
	"golang.org/x/image/colornames"
)

func Show(minLoadingTime int, logger bytes.Buffer, MainWindow fyne.Window, App fyne.App, Init func()) {
	// Initialize logger
	log.SetOutput(&logger)
	log.Println("Initialized a new logger")

	// Render a splash screen while the app is loading
	log.Println("Rendering splash screen...")
	var splashScreen fyne.Window
	var status *canvas.Text
	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		// Create splash screen
		splashScreen = drv.CreateSplashWindow()
		splashScreen.RequestFocus()

		// Draw the background
		background := canvas.NewImageFromFile("splash.png")
		background.Resize(fyne.NewSize(1024, 512)) // Resize image to desired size

		// Draw text
		status = canvas.NewText("Loading application...", colornames.Orangered)
		status.TextSize = 20

		// Draw rectangle
		rectangle := canvas.NewRectangle(color.RGBA{R: 9, G: 9, B: 9, A: 255})
		rectangle.Resize(fyne.NewSize(1024, 40))

		// Assemble UI
		content := container.NewWithoutLayout(background, rectangle, status)
		background.Move(fyne.NewPos(0, 0))
		status.Move(fyne.NewPos(10, 512-35))
		rectangle.Move(fyne.NewPos(0, 512-40))

		// Show splash screen
		splashScreen.SetContent(content)
		splashScreen.Resize(fyne.NewSize(1024, 512))
		rectangle.Refresh()
		status.Refresh()
		splashScreen.CenterOnScreen()
		splashScreen.Show()
		MainWindow.Hide()
	} else {
		log.Println("Failed to show loading splash screen")
		fatalerror.Show(fmt.Errorf("unexpected error getting desktop driver"), logger, MainWindow, App)
	}
	log.Println("Started loading")
	startTime := time.Now()

	// Update status in real time
	go func() {
		prevLog := ""
		for {
			if prevLog == logger.String() {
				continue
			}
			prevLog = logger.String()

			entries := strings.Split(logger.String(), "\n")
			status.Text = entries[len(entries)-2]
			status.Refresh()
			fmt.Println(entries[len(entries)-2])
		}
	}()

	// Some init
	Init()

	// Hide splash screen
	deltaTime := time.Since(startTime)
	status.Text = fmt.Sprintf("Finished loading in %s!", deltaTime.String())
	log.Println(status.Text)
	status.Refresh()
	time.Sleep(time.Second*time.Duration(minLoadingTime) - deltaTime)
	splashScreen.Hide()
}
