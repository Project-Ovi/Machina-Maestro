package playground

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
	"github.com/Project-Ovi/Machina-Maestro/windows/confirm"
	"github.com/Project-Ovi/Machina-Maestro/windows/fatalerror"
)

// System variables
var workingDirectory string
var logger bytes.Buffer

// Fyne elements
var MainWindow fyne.Window
var App fyne.App

// Collections
var defaultCommands []helper.Command //! Uninitialized
var actionCollection []helper.Action //! Uninitialized
var ThisModel *helper.Model

func Load(wd string, loggerObj bytes.Buffer, masterWindow fyne.Window, app fyne.App) {
	// Load system variables
	workingDirectory = wd
	logger = loggerObj

	// Load fyne elements
	MainWindow = masterWindow
	App = app
}

func Launch(window fyne.Window) {
	// Reset defaults
	defaultCommands = []helper.Command{}

	// Fetch actions
	b, err := os.ReadFile(path.Join(workingDirectory, "/myModels/", (*ThisModel).Name, "/actions.json"))
	if err != nil {
		fatalerror.Show(err, logger, MainWindow, App)
	}
	err = json.Unmarshal(b, &actionCollection)
	if err != nil {
		fatalerror.Show(err, logger, MainWindow, App)
	}

	// Load functions for this model
	go LoadSelector()

	// Make a main content
	mainContent := container.New(layout.NewStackLayout())

	// Get navbar and sidebar
	sidebar := playgroundSidebar(&mainContent)
	navbar := playgroundNavbar(window, sidebar)

	// Display contents
	content := container.NewBorder(navbar, nil, sidebar, nil, mainContent)
	window.SetContent(content)
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
		confirm.Show(App, "Are you sure you want to exit?", "You are about to exit to the main menu", func() {
			// landing.Launch(window)
			log.Println("The home button is not implemented yet")
		}, func() {})
	})

	// Add file button
	fileBTN := widget.NewButtonWithIcon("Reveal File", theme.Icon(theme.IconNameFile), func() {
		err := helper.OpenExplorer(path.Join(workingDirectory, "/myModels/", ThisModel.Name))
		if err != nil {
			fatalerror.Show(err, logger, MainWindow, App)
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
	helpBTN := widget.NewButtonWithIcon("Help", theme.Icon(theme.IconNameHelp), HelpWindow)

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
	mdURL := ThisModel.Website

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
			confirm.Show(App, "Confirm", "You are about to delete "+val.Name, func() {
				actionCollection = append(actionCollection[:i], actionCollection[i+1:]...)
				sidebarActions(content)
			}, func() {})
		})
		playBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameMediaPlay), func() {
			actionCollection[i].Running = true
			go func() {
				val.Run(defaultCommands)
				actionCollection[i].Running = false
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
		if !val.Running {
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
		this := helper.Action{
			Name:        name,
			Description: desc,
			Commands:    []helper.Command{},
			Running:     false,
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
func actionEditor(act *helper.Action, content **fyne.Container) {
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
