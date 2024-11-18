package ovipicker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

var thisModel helper.OVI

func modelAddWindow() {
	// Create window
	App := fyne.CurrentApp()
	window := App.NewWindow("Add model")
	window.SetOnClosed(func() {
		returnedText = "picker"
		waitPressGroup.Done()
	})

	// Create title
	title := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	title.TextSize = theme.Size(theme.SizeNameHeadingText)
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	// Create name form
	nameTitle := widget.NewLabel("Name:")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("My OVI")
	nameEntry.OnChanged = func(s string) {
		// Make sure we have a valid value
		if s == "" {
			thisModel.Name = "My OVI"
			return
		}

		// Update name
		thisModel.Name = s
	}

	// Create description form
	descriptionTitle := widget.NewLabel("Description:")
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("My beloved OVI")
	descriptionEntry.OnChanged = func(s string) {
		// Make sure we have a valid value
		if s == "" {
			thisModel.Description = "My beloved OVI"
			return
		}

		// Update description
		thisModel.Description = s
	}

	// Fetch products
	var productsNames []string
	func() {
		productsDirPath := path.Join(WD, "models")
		productsDir, err := os.ReadDir(productsDirPath)
		if err != nil {
			log.Println("Failed to read", productsDirPath, ". Reason:", err)
			// Attempt to create the directory
			err := os.Mkdir(productsDirPath, 0666)
			if err != nil {
				log.Println("Failed to create", productsDirPath)
				return
			}
			log.Println("Created", productsDirPath)
		}
		for _, val := range productsDir {
			if val.IsDir() {
				productsNames = append(productsNames, val.Name())
			}
		}
	}()

	// Create product form
	productTitle := widget.NewLabel("Product:")
	productEntry := widget.NewSelect(append(productsNames, "Get more online..."), func(s string) {})

	// Create form
	form := container.New(
		layout.NewFormLayout(),
		nameTitle, nameEntry,
		descriptionTitle, descriptionEntry,
		productTitle, productEntry,
	)

	// Make model options
	modelOptions := container.New(
		layout.NewFormLayout(),
	)

	// Make a submit button
	submitBTN := widget.NewButtonWithIcon("Submit", theme.Icon(theme.IconNameLogin), func() {
		b, _ := json.Marshal(thisModel)
		fmt.Println(string(b))
		window.Close()
	})

	// Make content
	content := container.New(
		layout.NewVBoxLayout(),
		title,
		layout.NewSpacer(),
		form,
		modelOptions,
		layout.NewSpacer(),
		submitBTN,
	)

	// Update model options based on the product select Select
	productEntry.OnChanged = func(s string) {
		if s == "Get more online..." {
			return
		}

		modelOptions.Objects = returnModelSpecificForm(s)
		thisModel.ProductName = s
		content.Refresh()

	}

	// Set window content
	window.SetContent(content)
	window.Show()
}
