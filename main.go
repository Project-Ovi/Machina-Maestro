package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	splashscreen "github.com/Project-Ovi/Machina-Maestro/windows/splash_screen"
)

var App fyne.App

func initi() {
	// Set logger output
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	log.SetOutput(logFile)
}

func main() {
	initi()

	App = app.New()

	splashscreen.Launch()

	App.Run()
}
