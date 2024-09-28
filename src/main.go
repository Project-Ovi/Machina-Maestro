package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

// Images
var ICON_robot_hand_round *canvas.Image

// Variables
var logger bytes.Buffer
var workingDirectory string

// Fyne variables
var App fyne.App
var MainWindow fyne.Window

func bigBTN(text string, extraPadding float32, size float32, f func()) *fyne.Container {
	btnText := canvas.NewText(text, theme.Color(theme.ColorNameForeground))
	btnText.TextSize = size
	btnBG := widget.NewButton(" ", f)
	btnBGHeightIncrements := 1
	for {
		desiredValue := fyne.MeasureText(btnText.Text, btnText.TextSize, btnText.TextStyle)
		actualValue := fyne.MeasureText(btnBG.Text, theme.TextSize(), btnText.TextStyle)
		// fmt.Println("Desired: ", desiredValue, "\nActual: ", actualValue)
		if desiredValue.Width+extraPadding > actualValue.Width {
			btnBG.Text += " "
			btnBG.Refresh()
			continue
		}

		if desiredValue.Height > actualValue.Height*float32(btnBGHeightIncrements) {
			btnBG.Text += "\n"
			btnBGHeightIncrements++
			btnBG.Refresh()
			// continue
		}

		break
	}
	btn := container.New(layout.NewCenterLayout(), btnBG, btnText)
	return btn
}

// * Particular windows
func startup(minLoadingTime int) {
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
		fatalError(fmt.Errorf("unexpected error getting desktop driver"))
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
func fatalError(err error) {
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

// Init function
func Init() {
	var err error

	// Get working directory
	log.Println("Fetching working directory...")
	workingDirectory, err = os.Getwd()
	if err != nil {
		fatalError(err)
	}

	// Reading Icons
	ICON_robot_hand_round = canvas.NewImageFromFile(path.Join(workingDirectory, "/assets/robot_hand_round.png"))
}

// Main page
func landingPage(window fyne.Window) {
	// Make title
	title := canvas.NewText("Machina Maestro", colornames.Orange)
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 32

	// Make icon
	ICON_robot_hand_round.SetMinSize(fyne.NewSquareSize(512))
	ICON_robot_hand_round.FillMode = canvas.ImageFillContain

	// Make button
	btn := bigBTN("Select your OVI", 40, 20, func() {
		modelSelectWindow(window)
	})

	// Assemble UI
	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), title, layout.NewSpacer()),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), ICON_robot_hand_round, layout.NewSpacer()),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), btn, layout.NewSpacer()),
		layout.NewSpacer(),
	)
	window.SetContent(content)
	window.Show()
	window.CenterOnScreen()
}

func modelSelectWindow(window fyne.Window) {
	// Make a navbar
	btnBack := widget.NewButton("", func() { landingPage(window) })
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

	// Arrange and draw
	content := container.New(
		layout.NewVBoxLayout(),
		navbar,
		title,
	)
	window.SetContent(content)
}

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

	// Add name field
	nameTitle := canvas.NewText("Name", theme.Color(theme.ColorNameForeground))
	nameOption := widget.NewEntry()
	nameOption.Text = "My OVI"

	// Add model select button
	selectTitle := canvas.NewText("Model", theme.Color(theme.ColorNameForeground))
	selectOption := widget.NewSelect(models, func(s string) {})

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
	btn := bigBTN("Submit", 60, 15, func() {
		saveSelector(form)
	})

	selectOption.OnChanged = func(s string) {
		form.RemoveAll()
		buildForm(form)
		formSelector(s, form)
	}

	content := container.New(
		layout.NewVBoxLayout(),
		navbar,
		title,
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		btn,
	)
	addWindow.SetContent(content)
	addWindow.Resize(fyne.NewSquareSize(512))
	addWindow.Show()
}

func main() {
	// Start fyne
	App = app.New()
	MainWindow = App.NewWindow("Machina Maestro")
	MainWindow.SetOnClosed(func() {
		App.Quit()
	})

	// Make sure we don't crash
	defer func() {
		if r := recover(); r != nil {
			fatalError(fmt.Errorf("unhandled crash"))
		}

		if err := os.WriteFile("log.txt", logger.Bytes(), os.ModePerm); err != nil {
			log.Println(err)
			fmt.Println(logger.String())
			panic(err)
		}
	}()

	// Render app
	go func() {
		startup(1)

		landingPage(MainWindow)
	}()

	// Display app
	MainWindow.ShowAndRun()
}

var models []string = []string{
	"OVI MK2",
	"OVI MK3",
	"OVI Mini",
}

func formSelector(s string, form *fyne.Container) {
	switch s {
	case "OVI MK2":
		form_OVI_MK2(form)
	case "OVI MK3":
		form_OVI_MK3(form)
	case "OVI Mini":
		form_OVI_Mini(form)
	}
}

func saveSelector(form *fyne.Container) {
	s := form.Objects[3].(*widget.Select).Selected
	switch s {
	case "OVI MK2":
		save_OVI_MK2(form)
	case "OVI MK3":
		save_OVI_MK3(form)
	case "OVI Mini":
		save_OVI_Mini(form)
	}
}

func form_OVI_MK2(form *fyne.Container) {
	// Add IP entry
	ipName := canvas.NewText("IP", theme.Color(theme.ColorNameForeground))
	ipEntry := widget.NewEntry()
	ipEntry.Text = "192.168.4.1"

	form.Objects = append(form.Objects,
		ipName, ipEntry,
	)

}

func save_OVI_MK2(form *fyne.Container) {
	fmt.Println("Saved a valid OVI MK2")
}

func form_OVI_MK3(form *fyne.Container) {
}

func save_OVI_MK3(form *fyne.Container) {
	fmt.Println("Saved a valid OVI MK3")
}

func form_OVI_Mini(form *fyne.Container) {
}

func save_OVI_Mini(form *fyne.Container) {
	fmt.Println("Saved a valid OVI Mini")
}
