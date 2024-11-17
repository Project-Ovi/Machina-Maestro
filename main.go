package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	splashscreen "github.com/Project-Ovi/Machina-Maestro/windows/splash_screen"
)

var App fyne.App

func main() {
	App = app.New()

	splashscreen.Launch(App)

	App.Run()
}
