package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
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
func openExplorer(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// On Windows, use 'explorer'
		cmd = exec.Command("explorer", path)
	case "darwin":
		// On macOS, use 'open'
		cmd = exec.Command("open", path)
	case "linux":
		// On Linux, use 'xdg-open' or 'gio open' based on availability
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
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
func confirmWindow(title string, subtitle string, yes func(), no func()) {
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

	// Fetch models
	objects := []fyne.CanvasObject{}
	modelsDir, err := os.ReadDir(path.Join(workingDirectory, "/myModels"))
	if err != nil {
		fatalError(err)
	}
	for _, val := range modelsDir {
		// Fetch
		rawModelCard, err := os.ReadFile(path.Join(workingDirectory, "/myModels/", val.Name(), "/model.json"))
		if err != nil {
			log.Println("Failed to load model:", val.Name()+". Skipping...")
			continue
		}
		var modelCard model
		err = json.Unmarshal(rawModelCard, &modelCard)
		if err != nil {
			fatalError(err)
		}

		// Make title
		title := canvas.NewText(modelCard.Name, theme.Color(theme.ColorNameForeground))
		title.TextSize = 18

		// Make trash button
		deleteBTN := widget.NewButton("", func() {
			confirmWindow("Confirm", "You are about to delete "+modelCard.Name, func() {
				err := os.RemoveAll(path.Join(workingDirectory, "/myModels/", val.Name()))
				if err != nil {
					fatalError(err)
				}
				time.Sleep(time.Millisecond * 100)
				modelSelectWindow(window)
			}, func() {

			})
		})
		deleteBTN.Icon = theme.Icon(theme.IconNameDelete)

		// Make a use button
		useBTN := widget.NewButton("", func() {
			thisModel = modelCard
			OVIPlayground(window)
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
		subtitle.Text = saveSelector(form)
		subtitle.Refresh()
		if subtitle.Text == "" {
			addWindow.Close()
			OVIPlayground(window)
		}
	})

	selectOption.OnChanged = func(s string) {
		form.RemoveAll()
		buildForm(form)
		formSelector(s, form)
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

func helpWindow() {
	// Make window
	window := App.NewWindow("Help")

	// Add logo
	ICON_robot_hand_round.SetMinSize(fyne.NewSquareSize(100))

	// Add title
	title := canvas.NewText("Machina Maestro", theme.Color(theme.ColorNameForeground))
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 30

	// Add subtitle
	subtitle := widget.NewLabel("A comprehensive UI application for controlling OVI products.")
	subtitle.Wrapping = fyne.TextWrapWord

	// Add github button
	btn := widget.NewButtonWithIcon("Github", theme.Icon(theme.IconNameComputer), func() {
		url, err := url.Parse("https://github.com/Project-Ovi/Machina-Maestro")
		if err != nil {
			fatalError(err)
		}
		App.OpenURL(url)
	})

	// Assemble UI
	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			ICON_robot_hand_round,
			container.New(
				layout.NewVBoxLayout(),
				title,
				subtitle,
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
		btn,
	)

	// Show window
	window.SetContent(content)
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(500, 200))
	window.Show()
}

func playgroundNavbar(window fyne.Window, sidebar fyne.CanvasObject) *fyne.Container {
	// Add sidebar button
	sidebarBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameList), func() {
		if sidebar.Visible() {
			sidebar.Hide()
			window.Canvas().Content().Refresh()
		} else {
			sidebar.Show()
			window.Canvas().Content().Refresh()
		}
	})

	// Add home button
	homeBTN := widget.NewButtonWithIcon("Home", theme.Icon(theme.IconNameHome), func() {
		confirmWindow("Are you sure you want to exit?", "You are about to exit to the main menu", func() {
			landingPage(window)
		}, func() {})
	})

	// Add file button
	fileBTN := widget.NewButtonWithIcon("Reveal File", theme.Icon(theme.IconNameFile), func() {
		err := openExplorer(path.Join(workingDirectory, "/myModels/", thisModel.Name))
		if err != nil {
			fatalError(err)
		}
	})

	// Add run button
	runBTN := widget.NewButtonWithIcon("Run", theme.Icon(theme.IconNameMailSend), func() {
		fmt.Println("Pressed run button!")
		//TODO: Make this work
	})

	// Add tools button
	toolsBTN := widget.NewButtonWithIcon("Tools", theme.Icon(theme.IconNameComputer), func() {
		fmt.Println("Pressed tools button!")
		//TODO: Make this work
	})

	// Add settings button
	settingsBTN := widget.NewButtonWithIcon("Settings", theme.Icon(theme.IconNameSettings), func() {
		fmt.Println("Pressed settings button!")
		//TODO: Make this work
	})

	// Add help button
	helpBTN := widget.NewButtonWithIcon("Help", theme.Icon(theme.IconNameHelp), helpWindow)

	// Build navbar
	navbar := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		container.New(
			layout.NewHBoxLayout(),
			sidebarBTN,
			fileBTN,
			toolsBTN,
			runBTN,
			settingsBTN,
			layout.NewSpacer(),
			homeBTN,
			helpBTN,
		),
	)

	return navbar
}

func sidebarOverview(content **fyne.Container) {
}

func sidebarInfo(content **fyne.Container) {
	// Display a loading
	progressbar := widget.NewProgressBarInfinite()
	progressbar.Start()
	(*content).RemoveAll()
	(*content).Add(
		container.New(
			layout.NewVBoxLayout(),
			widget.NewLabel("Loading..."),
			progressbar,
			layout.NewSpacer(),
		),
	)

	// Get markdown URL
	mdURL := thisModel.Website

	// Download markdown
	markdown := ""
	resp, err := http.Get(mdURL)
	if err != nil {
		markdown = "# Unable to connect to the internet"
	} else {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			markdown = "# Failed to decode response body"
		} else {
			markdown = string(body)
		}
	}

	// Display markdown
	markdownWidget := widget.NewRichTextFromMarkdown(markdown)
	markdownWidget.Wrapping = fyne.TextWrapWord
	(*content).RemoveAll()
	(*content).Add(container.NewVScroll(markdownWidget))
}

func playgroundSidebar(window fyne.Window, mainContent **fyne.Container) fyne.CanvasObject {
	// Add a overview button
	overviewBTN := widget.NewButtonWithIcon("Overview", theme.Icon(theme.IconNameHome), func() {})
	overviewBTN.Alignment = widget.ButtonAlignLeading

	// Open this tab and invalidate the button
	sidebarOverview(mainContent)
	overviewBTN.Disable()

	// Add info button
	infoBTN := widget.NewButtonWithIcon("Info", theme.Icon(theme.IconNameInfo), func() {})
	infoBTN.Alignment = widget.ButtonAlignLeading

	// Add actions button
	actionsBTN := widget.NewButtonWithIcon("Actions", theme.Icon(theme.IconNameMediaPlay), func() {})
	actionsBTN.Alignment = widget.ButtonAlignLeading

	// Add routines button
	routinesBTN := widget.NewButtonWithIcon("Routines", theme.Icon(theme.IconNameMediaReplay), func() {})
	routinesBTN.Alignment = widget.ButtonAlignLeading

	// Add function to the buttons
	overviewBTN.OnTapped = func() {
		// Invalidate this button
		overviewBTN.Disable()

		// Validate everything else
		infoBTN.Enable()
		actionsBTN.Enable()
		routinesBTN.Enable()

		// Open the menu
		sidebarOverview(mainContent)
	}
	infoBTN.OnTapped = func() {
		// Invalidate this button
		infoBTN.Disable()

		// Validate everything else
		overviewBTN.Enable()
		actionsBTN.Enable()
		routinesBTN.Enable()

		// Open menu
		sidebarInfo(mainContent)
	}
	actionsBTN.OnTapped = func() {
		// Invalidate this button
		actionsBTN.Disable()

		// Validate everything else
		overviewBTN.Enable()
		infoBTN.Enable()
		routinesBTN.Enable()

		// Open menu
		//TODO: Add a menu here
	}
	routinesBTN.OnTapped = func() {
		// Invalidate this button
		routinesBTN.Disable()

		// Validate everything else
		overviewBTN.Enable()
		actionsBTN.Enable()
		infoBTN.Enable()

		// Open menu
		// TODO: Add a menu here
	}

	// Assemble buttons
	content := container.New(
		layout.NewVBoxLayout(),
		overviewBTN,
		infoBTN,
		actionsBTN,
		routinesBTN,
	)

	// Tidy up and return
	buttons := container.NewVScroll(content)
	return container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		buttons,
	)
}

func OVIPlayground(window fyne.Window) {
	// Load functions for this model
	go loadSelector()

	// Make a main content
	mainContent := container.New(layout.NewStackLayout())

	// Get navbar and sidebar
	sidebar := playgroundSidebar(window, &mainContent)
	navbar := playgroundNavbar(window, sidebar)

	// Display contents
	content := container.NewBorder(navbar, nil, sidebar, nil, mainContent)
	window.SetContent(content)
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

type model struct {
	Name    string            `json:"name"`
	Model   string            `json:"model_name"`
	Website string            `json:"website"`
	Other   map[string]string `json:"other"`
}

var thisModel model

// * Selectors
func formSelector(s string, form *fyne.Container) {
	log.Println("Creating form for model", s+"...")
	switch s {
	case "OVI MK2":
		form_OVI_MK2(form)
	case "OVI MK3":
		form_OVI_MK3(form)
	case "OVI Mini":
		form_OVI_Mini(form)
	}
}
func saveSelector(form *fyne.Container) string {
	s := form.Objects[3].(*widget.Select).Selected
	log.Println("Saving model", s+"...")
	switch s {
	case "OVI MK2":
		return save_OVI_MK2(form)
	case "OVI MK3":
		return save_OVI_MK3(form)
	case "OVI Mini":
		return save_OVI_Mini(form)
	}

	return "Please select a valid model"
}
func loadSelector() {
	log.Println("Loading model", thisModel.Model+"...")
	switch thisModel.Model {
	case "OVI MK2":
		load_OVI_MK2()
	case "OVI MK3":
		load_OVI_MK3()
	case "OVI Mini":
		load_OVI_Mini()
	}
}

// * OVI MK2
func form_OVI_MK2(form *fyne.Container) {
	// Add IP entry
	ipName := canvas.NewText("IP", theme.Color(theme.ColorNameForeground))
	ipEntry := widget.NewEntry()
	ipEntry.Text = "192.168.4.1"
	ipEntry.Validator = validation.NewRegexp("^((https?|ftp):\\/\\/)?((([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,6})|(localhost)|(\\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\b)|(\\b(([0-9a-fA-F]{1,4}:){7}([0-9a-fA-F]{1,4}|:)|(::([0-9a-fA-F]{1,4}:){0,6}([0-9a-fA-F]{1,4}|:)))\\b))(:(\\d{1,5}))?(\\/[a-zA-Z0-9_.~%-]*)*(\\?[a-zA-Z0-9_.~%-&=]*)?(\\#[a-zA-Z0-9_-]*)?$", "Invalid IP/URL")
	ipEntry.Refresh()

	form.Objects = append(form.Objects,
		ipName, ipEntry,
	)
}
func save_OVI_MK2(form *fyne.Container) string {
	//
	canContinue := true
	errorMsg := ""

	// Get model name and check if it is unique
	name := form.Objects[1].(*widget.Entry).Text
	models, err := os.ReadDir("myModels")
	if err != nil {
		fatalError(err)
	}
	for _, val := range models {
		if val.IsDir() && val.Name() == name {
			// The name is not unique
			form.Objects[0].(*canvas.Text).Color = colornames.Red
			canContinue = false
			errorMsg = "A model with this name already exists"
		}
	}
	if canContinue {
		form.Objects[0].(*canvas.Text).Color = theme.Color(theme.ColorNameForeground)
	}

	// Get model IP and check it
	IP := form.Objects[7].(*widget.Entry).Text
	form.Refresh()
	if form.Objects[7].(*widget.Entry).Validate() != nil {
		form.Objects[6].(*canvas.Text).Color = colornames.Red
		canContinue = false
		errorMsg = "Invalid IP address"
	} else {
		form.Objects[6].(*canvas.Text).Color = theme.Color(theme.ColorNameForeground)
	}

	// If something is wrong, quit
	if !canContinue {
		return errorMsg
	}

	// Parse model information
	thisModel.Name = name
	thisModel.Model = "OVI MK2"
	thisModel.Website = "https://raw.githubusercontent.com/Project-Ovi/OVI-MK2/refs/heads/main/README.md"
	thisModel.Other = map[string]string{"IP": IP}

	// Save model to file
	os.Mkdir(path.Join(workingDirectory, "/myModels/", name), os.ModePerm)
	information, err := json.Marshal(thisModel)
	if err != nil {
		errorMsg = err.Error()
	}
	err = os.WriteFile(path.Join(workingDirectory, "/myModels/", name, "/model.json"), information, os.ModePerm)
	if err != nil {
		errorMsg = err.Error()
	}

	// Return any error messages
	return errorMsg
}
func load_OVI_MK2() {
	// * Helper functions
	post := func(url string, args map[string]string) error {
		payload := []byte(`{"key1":"value1", "key2":"value2"}`)

		// Create a new POST request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		// Add custom headers from the args map
		for key, value := range args {
			req.Header.Set(key, value)
		}

		// Create an HTTP client
		client := &http.Client{}

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		return nil
	}

	// * Load commands
	// Load rotate command
	rotCom := command{
		DisplayName: "Set rotation speed",
		Arguments: []argument{
			{Name: "Speed", ArgType: "int"},
		},
		f: func(a []argument) error {
			var R1, R2 int
			if a[0].Value.(int) > 0 {
				R1 = a[0].Value.(int)
				R2 = 0
			} else {
				R1 = 0
				R2 = -a[0].Value.(int)
			}

			if R1 > 255 || R2 > 255 {
				return fmt.Errorf("value out of range. allowed: -256 < speed < 256")
			}

			return post(thisModel.Other["IP"], map[string]string{
				"R1": fmt.Sprint(R1),
				"R2": fmt.Sprint(R2),
			})
		},
	}

	// Load move up command
	moveupCom := command{
		DisplayName: "Move up with speed",
		Arguments: []argument{
			{Name: "Speed", ArgType: "int"},
		},
		f: func(a []argument) error {
			var U1, U2 int
			if a[0].Value.(int) > 0 {
				U1 = a[0].Value.(int)
				U2 = 0
			} else {
				U1 = 0
				U2 = -a[0].Value.(int)
			}

			if U1 > 255 || U2 > 255 {
				return fmt.Errorf("value out of range. allowed: -256 < speed < 256")
			}

			return post(thisModel.Other["IP"], map[string]string{
				"U1": fmt.Sprint(U1),
				"U2": fmt.Sprint(U2),
			})
		},
	}

	// Load extend forward command
	extendCom := command{
		DisplayName: "Extend forward with speed",
		Arguments: []argument{
			{Name: "Speed", ArgType: "int"},
		},
		f: func(a []argument) error {
			var E1, E2 int
			if a[0].Value.(int) > 0 {
				E1 = a[0].Value.(int)
				E2 = 0
			} else {
				E1 = 0
				E2 = -a[0].Value.(int)
			}

			if E1 > 255 || E2 > 255 {
				return fmt.Errorf("value out of range. allowed: -256 < speed < 256")
			}

			return post(thisModel.Other["IP"], map[string]string{
				"E1": fmt.Sprint(E1),
				"E2": fmt.Sprint(E2),
			})
		},
	}

	// Load grip command
	gripCom := command{
		DisplayName: "Set gripper state",
		Arguments: []argument{
			{Name: "Grip", ArgType: "bool"},
		},
		f: func(a []argument) error {
			var G1 = int(0)

			if a[0].Value.(bool) {
				G1 = int(255)
			}

			return post(thisModel.Other["IP"], map[string]string{
				"G1": fmt.Sprint(G1),
			})
		},
	}

	// Load wait time command
	waitCom := command{
		DisplayName: "Wait",
		Arguments: []argument{
			{Name: "Milliseconds", ArgType: "int"},
		},
		f: func(a []argument) error {
			time.Sleep(time.Millisecond * time.Duration(a[0].Value.(int)))
			return nil
		},
	}

	// Append Commands to the default commands list
	defaultCommands = append(defaultCommands, rotCom, moveupCom, extendCom, gripCom, waitCom)
}

// * OVI MK3
func form_OVI_MK3(form *fyne.Container) {
}
func save_OVI_MK3(form *fyne.Container) string {
	fmt.Println("Saved a valid OVI MK3")
	return "Unimplemented"
}
func load_OVI_MK3() {
}

// * OVI Mini
func form_OVI_Mini(form *fyne.Container) {
}
func save_OVI_Mini(form *fyne.Container) string {
	fmt.Println("Saved a valid OVI Mini")
	return "Unimplemented"
}
func load_OVI_Mini() {
}

// * Action definitions
type action struct {
	Name        string
	Description string
	Commands    []command
}

type command struct {
	DisplayName string
	Arguments   []argument
	f           func([]argument) error
}

type argument struct {
	Name    string
	ArgType string
	Value   interface{}
}

func (a action) Run() error {
	for _, val := range a.Commands {
		if err := val.f(val.Arguments); err != nil {
			return err
		}
	}

	return nil
}

var defaultCommands []command
var actionCollection []action
