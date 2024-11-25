package playground

import (
	"fmt"
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

func actionsPage(content *fyne.Container) {
	// Make navbar
	navbarAddBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), func() {})
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
