package modelselect

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
	"github.com/Project-Ovi/Machina-Maestro/windows/confirm"
	"github.com/Project-Ovi/Machina-Maestro/windows/fatalerror"
	"github.com/Project-Ovi/Machina-Maestro/windows/playground"
)

var workingDirectory string
var logger bytes.Buffer
var MainWindow fyne.Window
var App fyne.App

func Launch(window fyne.Window, wd string, loggerObj bytes.Buffer, MasterWindow fyne.Window, MasterApp fyne.App) {
	// Load data into variables
	workingDirectory = wd
	logger = loggerObj
	MainWindow = MasterWindow
	App = MasterApp

	// Make a navbar
	btnBack := widget.NewButton("", func() {
		// landing.Launch(window)
		log.Println("Back button not implemented yet!")
	})
	btnBack.SetIcon(theme.Icon(theme.IconNameNavigateBack))
	btnAdd := widget.NewButton("", func() {})
	btnAdd.OnTapped = func() {
		modelAddWindow(window, btnAdd)
		btnAdd.Disable()
	}
	btnAdd.SetIcon(theme.Icon(theme.IconNameContentAdd))
	navbar := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameButton)),
		container.New(
			layout.NewHBoxLayout(),
			btnBack,
			layout.NewSpacer(),
			btnAdd,
		),
	)

	// Make a title
	title := canvas.NewText("Select an OVI", theme.Color(theme.ColorNameForeground))
	title.TextSize = 32
	title.Alignment = fyne.TextAlignCenter

	// Fetch models
	objects := []fyne.CanvasObject{}
	modelsDir, err := os.ReadDir(path.Join(workingDirectory, "/myModels"))
	if err != nil {
		fatalerror.Show(err, logger, MainWindow, App)
	}
	for _, val := range modelsDir {
		// Fetch
		rawModelCard, err := os.ReadFile(path.Join(workingDirectory, "/myModels/", val.Name(), "/model.json"))
		if err != nil {
			log.Println("Failed to load model:", val.Name()+". Skipping...")
			continue
		}
		var modelCard helper.Model
		err = json.Unmarshal(rawModelCard, &modelCard)
		if err != nil {
			fatalerror.Show(err, logger, MainWindow, App)
		}

		// Make title
		title := canvas.NewText(modelCard.Name, theme.Color(theme.ColorNameForeground))
		title.TextSize = 18

		// Make trash button
		deleteBTN := widget.NewButton("", func() {
			confirm.Show(App, "Confirm", "You are about to delete "+modelCard.Name, func() {
				err := os.RemoveAll(path.Join(workingDirectory, "/myModels/", val.Name()))
				if err != nil {
					fatalerror.Show(err, logger, MainWindow, App)
				}
				time.Sleep(time.Millisecond * 100)
				Launch(window, workingDirectory, logger, MainWindow, App)
			}, func() {

			})
		})
		deleteBTN.Icon = theme.Icon(theme.IconNameDelete)

		// Make a use button
		useBTN := widget.NewButton("", func() {
			playground.ThisModel = &modelCard
			playground.Launch(window)
		})
		useBTN.Icon = theme.Icon(theme.IconNameMediaPlay)

		// Build sidebar
		sidebar := container.New(
			layout.NewHBoxLayout(),
			deleteBTN,
			useBTN,
		)

		thisObject := container.New(
			layout.NewStackLayout(),
			canvas.NewRectangle(theme.Color(theme.ColorNameButton)),
			container.New(
				layout.NewHBoxLayout(),
				title,
				layout.NewSpacer(),
				sidebar,
			),
		)

		objects = append(objects, thisObject)

	}

	// Display available models
	modelEntries := container.New(
		layout.NewVBoxLayout(),
		objects...,
	)
	modelsList := container.NewVScroll(modelEntries)
	modelsList.SetMinSize(fyne.NewSize(0, 500))

	// Arrange and draw
	content := container.New(
		layout.NewVBoxLayout(),
		navbar,
		title,
		layout.NewSpacer(),
		modelsList,
		layout.NewSpacer(),
	)
	window.SetContent(content)
}
