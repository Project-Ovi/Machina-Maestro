package main

import (
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	homescreen "github.com/Project-Ovi/Machina-Maestro/windows/home_screen"
	splashscreen "github.com/Project-Ovi/Machina-Maestro/windows/splash_screen"
)

var App fyne.App
var MainWindow fyne.Window

var logFile *os.File

func preInit() {
	// Set logger output
	var err error
	logFile, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	log.SetOutput(logFile)
}

func Init() {
	time.Sleep(time.Second * 3)
}

func maini() {
	splashscreen.Launch(Init)

	homescreen.Launch(MainWindow)
}

func main() {
	// Run preinitialization
	preInit()

	// Create App
	App = app.New()
	MainWindow = App.NewWindow("Machina Maestro")
	MainWindow.SetMaster()
	MainWindow.Show()

	// Launch main function
	go maini()

	// Run app
	App.Run()
}
