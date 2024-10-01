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
	"strconv"
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

// * Main page
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

// * Model selection
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

// * Playground additional windows

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

// * Playground elements
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

func playgroundSidebar(mainContent **fyne.Container) fyne.CanvasObject {
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
		sidebarActions(mainContent)
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

// * Sidebar elements
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

func sidebarActions(content **fyne.Container) {
	// Save actions
	saveAllActions()

	// Add title
	title := canvas.NewText("Actions", theme.Color(theme.ColorNameForeground))
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 32

	// Add navbar
	btn := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), func() {})
	btn.OnTapped = func() {
		actionCreate(btn, content)
	}
	navbar := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			btn,
		),
	)

	// Build topbar
	topbar := container.New(
		layout.NewVBoxLayout(),
		navbar,
		title,
	)

	// Get actions list
	actionsDisp := container.New(
		layout.NewVBoxLayout(),
	)
	for i, val := range actionCollection {
		// Make buttons
		editBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDocumentCreate), func() { actionEditor(&actionCollection[i], content) })
		trashBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDelete), func() {
			confirmWindow("Confirm", "You are about to delete "+val.Name, func() {
				actionCollection = append(actionCollection[:i], actionCollection[i+1:]...)
				sidebarActions(content)
			}, func() {})
		})
		playBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameMediaPlay), func() {
			actionCollection[i].running = true
			go func() {
				val.Run()
				actionCollection[i].running = false
				sidebarActions(content)
			}()
			time.Sleep(time.Millisecond * 5)
			sidebarActions(content)
		})
		buttons := container.New(
			layout.NewHBoxLayout(),
			editBTN,
			trashBTN,
			playBTN,
		)

		// Make backdrop
		var backdrop fyne.CanvasObject
		if !val.running {
			backdrop = canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground))
		} else {
			backdrop = widget.NewProgressBarInfinite()
			backdrop.(*widget.ProgressBarInfinite).Start()
		}

		actionsDisp.Add(container.New(
			layout.NewStackLayout(),
			backdrop,
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel(val.Name),
				layout.NewSpacer(),
				buttons,
			),
		))
	}

	// Build actions list
	actionsList := container.NewVScroll(actionsDisp)

	// Build UI
	myContent := container.NewBorder(topbar, nil, nil, nil, actionsList)

	// Display contents
	(*content).RemoveAll()
	(*content).Add(myContent)
}

// * Additional windows
func actionCreate(summoner *widget.Button, parent **fyne.Container) {
	// Make sure we can't create new windows while this one is open
	summoner.Disable()

	// Create a window
	window := App.NewWindow("Create Action")
	window.SetFixedSize(true)

	// Allow further creation after this window is closed
	window.SetCloseIntercept(func() {
		summoner.Enable()
		window.Close()
	})

	// Create a title
	title := canvas.NewText("Create a new action", theme.Color(theme.ColorNameForeground))
	title.TextSize = 32
	title.Alignment = fyne.TextAlignCenter

	// Create an entry for name
	nameText := widget.NewLabel("Name")
	nameEntry := widget.NewEntry()
	nameEntry.Validator = validation.NewRegexp("^.+$", "String can't be empty")

	// Create an entry for description
	descText := widget.NewLabel("Description")
	descEntry := widget.NewEntry()
	descEntry.Validator = validation.NewRegexp("^.+$", "String can't be empty")

	// Build form
	form := container.New(
		layout.NewFormLayout(),
		nameText, nameEntry,
		descText, descEntry,
	)

	// Add create button
	btn := widget.NewButtonWithIcon("Create", theme.Icon(theme.IconNameDocumentCreate), func() {
		// Get data
		name := nameEntry.Text
		desc := descEntry.Text

		// Make sure the data is correct
		if nameEntry.Validate() != nil || descEntry.Validate() != nil {
			return
		}

		// Make a  new action
		this := action{
			Name:        name,
			Description: desc,
			Commands:    []command{},
			running:     false,
		}

		// Add the action to the list of actions
		actionCollection = append(actionCollection, this)

		// Close this window
		summoner.Enable()
		sidebarActions(parent)
		window.Close()
	})

	// Build UI
	content := container.New(
		layout.NewVBoxLayout(),
		title,
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		btn,
		layout.NewSpacer(),
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 200))
	window.Show()
}

// * Action editor
func actionEditor(act *action, content **fyne.Container) {
	saveAllActions()
	// Make a navbar
	navbar := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameNavigateBack), func() { sidebarActions(content) }),
		),
	)

	// Make a topbar
	title := canvas.NewText((*act).Name, theme.Color(theme.ColorNameForeground))
	subtitle := canvas.NewText((*act).Description, theme.Color(theme.ColorNameForeground))
	title.Alignment = fyne.TextAlignCenter
	subtitle.Alignment = fyne.TextAlignCenter
	title.TextSize = 32
	subtitle.TextSize = 18
	topbar := container.New(
		layout.NewVBoxLayout(),
		navbar,
		title,
		subtitle,
	)

	// Make actions editor
	commandsDisplay := []fyne.CanvasObject{}
	for i, val := range (*act).Commands {
		// Make name
		functionsNames := []string{}
		for _, com := range defaultCommands {
			functionsNames = append(functionsNames, com.DisplayName)
		}
		name := widget.NewSelect(functionsNames, func(s string) {
			// Find the function associated and replace it with it
			for j, com := range defaultCommands {
				if s != com.DisplayName {
					continue
				}
				(*act).Commands[i] = defaultCommands[j]
			}

			actionEditor(act, content)
		})
		name.Selected = val.DisplayName

		// Get arguments
		argsDisplay := []fyne.CanvasObject{}
		for j, arg := range val.Arguments {
			// Make arg label
			argName := widget.NewLabel(arg.Name)

			// Make arg entry
			var argEntry fyne.CanvasObject
			switch arg.ArgType {
			case "int":
				this := widget.NewEntry()
				this.Validator = validation.NewRegexp("^-?\\d+$", "Invalid integer")
				if arg.Value != nil {
					this.Text = fmt.Sprint(arg.Value)
				}
				this.OnChanged = func(s string) {
					if this.Validate() != nil {
						return
					}
					num, err := strconv.Atoi(s)
					if err != nil {
						return
					}

					(*act).Commands[i].Arguments[j].Value = num
					saveAllActions()
				}
				argEntry = this
			case "float":
				this := widget.NewEntry()
				this.Validator = validation.NewRegexp("^-?\\d+(\\.\\d+)?$", "Invalid integer")
				if arg.Value != nil {
					this.Text = fmt.Sprint(arg.Value)
				}
				this.OnChanged = func(s string) {
					if this.Validate() != nil {
						return
					}
					num, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return
					}

					(*act).Commands[i].Arguments[j].Value = num
					saveAllActions()
				}
				argEntry = this
			case "string":
				this := widget.NewEntry()
				if arg.Value != nil {
					this.Text = fmt.Sprint(arg.Value)
				}
				this.OnChanged = func(s string) {
					if this.Validate() != nil {
						return
					}

					(*act).Commands[i].Arguments[j].Value = s
					saveAllActions()
				}
				argEntry = this
			case "bool":
				this := widget.NewCheck("", func(b bool) {})
				if arg.Value != nil && arg.Value.(bool) {
					this.Checked = true
				}
				this.OnChanged = func(b bool) {
					(*act).Commands[i].Arguments[j].Value = b
					saveAllActions()
				}
				argEntry = this
			}

			// Assemble arg element
			argElement := container.New(
				layout.NewFormLayout(),
				argName, argEntry,
			)
			argsDisplay = append(argsDisplay, argElement)
		}

		// Make command toolbar
		toolbar := container.New(
			layout.NewHBoxLayout(),
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentRemove), func() {
				(*act).Commands = append((*act).Commands[:i], (*act).Commands[i+1:]...)
				actionEditor(act, content)
			}),
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameMoveUp), func() {
				if i > 0 {
					(*act).Commands[i], (*act).Commands[i-1] = (*act).Commands[i-1], (*act).Commands[i]
				}
				actionEditor(act, content)
			}),
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameMoveDown), func() {
				if i < len((*act).Commands)-1 {
					(*act).Commands[i], (*act).Commands[i+1] = (*act).Commands[i+1], (*act).Commands[i]
				}
				actionEditor(act, content)
			}),
		)

		// Assemble command arguments
		argsAssembled := container.New(
			layout.NewHBoxLayout(),
			argsDisplay...,
		)

		// Assemble command display
		commandAssembled := container.New(
			layout.NewStackLayout(),
			canvas.NewRectangle(theme.Color(theme.ColorNameButton)),
			container.New(
				layout.NewHBoxLayout(),
				name,
				layout.NewSpacer(),
				argsAssembled,
				layout.NewSpacer(),
				toolbar,
			),
		)

		commandsDisplay = append(commandsDisplay, commandAssembled)
	}

	// Make a add button
	addBTN := widget.NewButtonWithIcon("New Instruction", theme.Icon(theme.IconNameContentAdd), func() {
		(*act).Commands = append((*act).Commands, defaultCommands[0])
		actionEditor(act, content)
	})

	// Make editor context
	editor := container.New(
		layout.NewVBoxLayout(),
		commandsDisplay...,
	)
	editor.Add(addBTN)
	scrollingEditor := container.NewVScroll(editor)

	// Display
	(*content).RemoveAll()
	(*content).Add(container.NewBorder(
		topbar, nil, nil, nil,
		scrollingEditor,
	))
}

// * Playground
func OVIPlayground(window fyne.Window) {
	// Reset defaults
	defaultCommands = []command{}

	// Fetch actions
	b, err := os.ReadFile(path.Join(workingDirectory, "/myModels/", thisModel.Name, "/actions.json"))
	if err != nil {
		fatalError(err)
	}
	err = json.Unmarshal(b, &actionCollection)
	if err != nil {
		fatalError(err)
	}

	// Load functions for this model
	go loadSelector()

	// Make a main content
	mainContent := container.New(layout.NewStackLayout())

	// Get navbar and sidebar
	sidebar := playgroundSidebar(&mainContent)
	navbar := playgroundNavbar(window, sidebar)

	// Display contents
	content := container.NewBorder(navbar, nil, sidebar, nil, mainContent)
	window.SetContent(content)
}

// ! MAIN
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

		startup(1)

		landingPage(MainWindow)
	}()

	// Display app
	MainWindow.ShowAndRun()
}

// * Model definitions
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
	if name == "" {
		name = "Unnamed"
	}
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
			{Name: "Speed", ArgType: "int", Value: int(0)},
		},
		f: func(a []argument) error {
			// Fix up data
			if val, ok := a[0].Value.(float64); ok {
				a[0].Value = int(val)
			}

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
			{Name: "Speed", ArgType: "int", Value: int(0)},
		},
		f: func(a []argument) error {
			// Fix up data
			if val, ok := a[0].Value.(float64); ok {
				a[0].Value = int(val)
			}

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
			{Name: "Speed", ArgType: "int", Value: int(0)},
		},
		f: func(a []argument) error {
			// Fix up data
			if val, ok := a[0].Value.(float64); ok {
				a[0].Value = int(val)
			}

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
			{Name: "Grip", ArgType: "bool", Value: bool(false)},
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
			{Name: "Milliseconds", ArgType: "int", Value: int(0)},
		},
		f: func(a []argument) error {
			// Fix up data
			if val, ok := a[0].Value.(float64); ok {
				a[0].Value = int(val)
			}

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
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Commands    []command `json:"commands"`
	running     bool
}

type command struct {
	DisplayName string     `json:"name"`
	Arguments   []argument `json:"args"`
	f           func([]argument) error
}

type argument struct {
	Name    string      `json:"name"`
	ArgType string      `json:"type"`
	Value   interface{} `json:"val"`
}

func (a *action) Run() error {
	// Fix function
	a.Fix()

	// Run
	for _, val := range (*a).Commands {
		if err := val.f(val.Arguments); err != nil {
			return err
		}
	}

	// Return
	return nil
}

func (a *action) Fix() {
	for i, val1 := range (*a).Commands {
		for _, val2 := range defaultCommands {
			if val1.DisplayName == val2.DisplayName {
				(*a).Commands[i].f = val2.f
			}
		}
	}
}

func saveAllActions() error {
	// Marshall
	b, err := json.Marshal(actionCollection)
	if err != nil {
		return err
	}

	// Save to file
	return os.WriteFile(path.Join(workingDirectory, "/myModels/", thisModel.Name, "/actions.json"), b, os.ModePerm)
}

var defaultCommands []command
var actionCollection []action
