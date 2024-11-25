package playground

import (
	"fmt"
	"log"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

var actionsCollection []helper.Action
var actionsPlayQueue []int

func actionCreateDialog() {
	// Make sure the function doesn't exit prematurely
	var dialogWait sync.WaitGroup
	dialogWait.Add(1)

	// Create a window
	App := fyne.CurrentApp()
	window := App.NewWindow("Create an Action")
	window.SetOnClosed(dialogWait.Done)

	// Create a title
	titleObj := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Create name entry
	nameTitleObj := widget.NewLabel("Name:")
	nameEntryObj := widget.NewEntry()

	// Create description entry
	descriptionTitleObj := widget.NewLabel("Description:")
	descriptionEntryObj := widget.NewEntry()

	// Create create button
	createBTN := widget.NewButtonWithIcon("Create", theme.Icon(theme.IconNameContentAdd), func() {})
	createBTN.OnTapped = func() {
		createBTN.Disable()

		// Validate name
		if nameEntryObj.Text == "" {
			log.Println("Action name can't be empty")
			return
		} else if err := nameEntryObj.Validate(); err != nil {
			log.Println(err)
			return
		}

		// Make sure the name is unique
		for _, val := range actionsCollection {
			if val.Name == nameEntryObj.Text {
				log.Println("An action with this name already exists:", val.Name)
				return
			}
		}

		// Validate description
		if err := descriptionEntryObj.Validate(); err != nil {
			log.Println(err)
			return
		}

		// Create action
		actionsCollection = append(actionsCollection, helper.Action{
			Name:        nameEntryObj.Text,
			Description: descriptionEntryObj.Text,
		})

		// Close window
		window.Close()
	}

	// Make content
	content := container.New(
		layout.NewVBoxLayout(),
		titleObj,
		layout.NewSpacer(),
		container.New(
			layout.NewFormLayout(),
			nameTitleObj, nameEntryObj,
			descriptionTitleObj, descriptionEntryObj,
		),
		layout.NewSpacer(),
		createBTN,
		layout.NewSpacer(),
	)

	// Show window
	window.SetContent(content)
	window.Show()

	// Await an exit signal
	dialogWait.Wait()
}

func actionsPage(content *fyne.Container) {
	// Make navbar
	navbarAddBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), func() {})
	navbarAddBTN.OnTapped = func() {
		navbarAddBTN.Disable()
		actionCreateDialog()
		actionsPage(content)
	}
	actionsNavbar := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			navbarAddBTN,
		),
	)

	// Make title
	titleObj := canvas.NewText("Actions", theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Make action list
	actionsList := container.New(layout.NewVBoxLayout())
	for i, val := range actionsCollection {
		// Make action title
		actionTitleObj := canvas.NewText(val.Name, theme.Color(theme.ColorNameForeground))

		// Make delete button
		deleteBTN := widget.NewButtonWithIcon("Delete", theme.Icon(theme.IconNameDelete), func() {})

		// Make edit button
		editBTN := widget.NewButtonWithIcon("Edit", theme.Icon(theme.IconNameSettings), func() {})

		// Make play button
		var playBTN *widget.Button
		if !val.Running {
			// Make a normal play button
			playBTN = widget.NewButtonWithIcon("Start", theme.Icon(theme.IconNameMediaPlay), func() {})
		} else {
			// Make a stop button
			playBTN = widget.NewButtonWithIcon("Stop", theme.Icon(theme.IconNameMediaStop), func() {})

			// Show place in queue
			go func() {
				for len(actionsPlayQueue) > 1 && actionsPlayQueue[0] != i {
					time.Sleep(time.Millisecond * 500)
					for j, val := range actionsPlayQueue {
						if val == i {
							playBTN.SetText(fmt.Sprint("In queue", j, "/", len(actionsPlayQueue)))
						}
					}
				}
			}()
		}

		// Add to actions list
		actionsList.Add(container.New(
			layout.NewStackLayout(),
			canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
			container.New(
				layout.NewHBoxLayout(),
				actionTitleObj,
				layout.NewSpacer(),
				deleteBTN,
				editBTN,
				playBTN,
			),
		))
	}

	// Make content
	content.Layout = &helper.StackWithNavbar{}
	content.Objects = []fyne.CanvasObject{
		actionsNavbar,
		container.New(
			&helper.StackWithNavbar{},
			titleObj,
			container.NewVScroll(actionsList),
		),
	}
}
