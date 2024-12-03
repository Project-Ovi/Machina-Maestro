package playground

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

func actionsEditor(content *fyne.Container, act *helper.Action) {
	log.Printf("Pointer address: %p\n", act)
	// Make navbar
	editorNavbar := container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameNavigateBack), func() { actionsPage(content) }),
			layout.NewSpacer(),
		),
	)

	// Make title
	titleObj := canvas.NewText(act.Name, theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Make subtitle
	subtitleObj := canvas.NewText(act.Description, theme.Color(theme.ColorNameForeground))
	subtitleObj.Alignment = fyne.TextAlignCenter

	// Make functions
	funcsContainerObj := container.New(layout.NewVBoxLayout())
	for findex, val := range act.Functions {
		// Make func name
		funcnameObj := widget.NewLabel(val.FunctionName)
		funcnameObj.TextStyle.Bold = true

		// Make form
		argsForm := container.New(layout.NewFormLayout())
		args := make(map[string]string)
		for _, val2 := range BuiltinFunctions {
			if val2.Name == val.FunctionName {
				for keys, values := range val2.Arguments {
					args[keys] = values
				}
				break
			}
		}
		for argName := range args {
			// Make argument title
			argTitleObj := widget.NewLabel(argName)

			// Extract argument value
			argValueName, ArgValueValue, _ := strings.Cut(args[argName], "/")
			log.Println("ArgValueName: ", argValueName)

			// Make argument entry
			var argEntryObj fyne.CanvasObject
			argEntryObj = widget.NewLabel("Invalid config value")
			switch argValueName {
			case "Number":
				// Make entry
				argEntryObj = widget.NewEntry()
				argEntryObj.(*widget.Entry).Validator = validation.NewRegexp("(?:\\d+(?:\\.\\d*)?|\\.\\d+)", "Not a number")
				argEntryObj.(*widget.Entry).OnChanged = func(s string) {
					act.Functions[findex].Arguments[argName] = s
				}
				argEntryObj.(*widget.Entry).Text = act.Functions[findex].Arguments[argName]
			case "Slider":
				// Split min-max
				vals := strings.Split(ArgValueValue, "-")
				if len(vals) == 2 {
					// Parse min
					minVal, err := strconv.ParseFloat(vals[0], 64)
					if err == nil {
						// Parse max
						maxVal, err := strconv.ParseFloat(vals[1], 64)
						if err == nil {
							// Make widget
							argEntryObj = widget.NewSlider(minVal, maxVal)
							argEntryObj.(*widget.Slider).OnChanged = func(f float64) { act.Functions[findex].Arguments[argName] = fmt.Sprint(f) }
						}
					}
				}
			case "Bool":
				argEntryObj = widget.NewCheck("", func(b bool) {
					act.Functions[findex].Arguments[argName] = fmt.Sprint(b)
				})
			case "String":
				argEntryObj = widget.NewEntry()
				argEntryObj.(*widget.Entry).OnChanged = func(s string) {
					act.Functions[findex].Arguments[argName] = s
				}
			default:
				log.Println("Failed to match arg value name:", argValueName)
			}

			// Add values to the form
			argsForm.Add(argTitleObj)
			argsForm.Add(argEntryObj)
		}

		// Make delete button
		deleteBTN := container.New(
			layout.NewVBoxLayout(),
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDelete), func() {
				(*act).Functions = append((*act).Functions[:findex], (*act).Functions[findex+1:]...)
				actionsEditor(content, act)
			}),
			layout.NewSpacer(),
		)

		// Assemble layout
		funcsContainerObj.Add(
			container.New(
				layout.NewStackLayout(),
				canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
				container.New(
					layout.NewHBoxLayout(),
					funcnameObj,
					layout.NewSpacer(),
					argsForm,
					layout.NewSpacer(),
					deleteBTN,
				),
			),
		)
	}

	// Add function add button
	functionCreateBTNContainer := container.New(layout.NewHBoxLayout())
	functionCreateBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), func() {
		// Remove original button
		functionCreateBTNContainer.RemoveAll()

		// Create a dropdown
		var availableFunctionNames []string
		for _, val := range BuiltinFunctions {
			availableFunctionNames = append(availableFunctionNames, val.Name)
		}
		dropdownObj := widget.NewSelect(availableFunctionNames, func(s string) {})

		// Create create button
		createButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), func() {
			// Discard unselected
			if dropdownObj.Selected == "" {
				return
			}

			// Make the new function
			this := helper.ModelFunction{}
			this.FunctionName = dropdownObj.Selected
			for _, val := range BuiltinFunctions {
				if val.Name == this.FunctionName {
					// this.Arguments = val.Arguments
					this.Arguments = make(map[string]string)
					for keys, values := range val.Arguments {
						this.Arguments[keys] = values
					}
					break
				}
			}

			// Add function to the action
			act.Functions = append(act.Functions, this)

			// Reload
			actionsEditor(content, act)
		})

		// Make Layout
		functionCreateBTNContainer.Add(layout.NewSpacer())
		functionCreateBTNContainer.Add(dropdownObj)
		functionCreateBTNContainer.Add(createButton)
		functionCreateBTNContainer.Add(layout.NewSpacer())
	})
	functionCreateBTNContainer.Add(functionCreateBTN)
	funcsContainerObj.Add(container.New(
		layout.NewStackLayout(),
		canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		functionCreateBTNContainer,
	))

	// Set content
	content.Layout = &helper.StackWithNavbar{}
	content.Objects = []fyne.CanvasObject{
		container.New(
			layout.NewVBoxLayout(),
			editorNavbar,
			titleObj,
			subtitleObj,
		),
		container.NewVScroll(funcsContainerObj),
	}

}
