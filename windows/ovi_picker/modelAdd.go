package ovipicker

import (
	"encoding/json"
	"fmt"
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
)

var thisModel helper.OVI

func modelAddWindow() {
	// Init
	thisModel.Others = make(map[string]string)

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
		// Update name
		if s == "" {
			thisModel.Name = "My OVI"
		} else {
			thisModel.Name = s
		}
		fmt.Println(thisModel.Name)
	}
	thisModel.Name = "My OVI"

	// Create description form
	descriptionTitle := widget.NewLabel("Description:")
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("My beloved OVI")
	descriptionEntry.OnChanged = func(s string) {
		// Update description
		if s == "" {
			thisModel.Description = "My beloved OVI"
		} else {
			thisModel.Description = s
		}
	}
	thisModel.Description = "My beloved OVI"

	// Fetch products
	var productsNames []string
	func() {
		productsDirPath := path.Join(WD, "models")
		productsDir, err := os.ReadDir(productsDirPath)
		if err != nil {
			log.Println("Failed to read", productsDirPath, ". Reason:", err)
			// Attempt to create the directory
			err := os.Mkdir(productsDirPath, 0766)
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
		// Make sure a product is selected
		if thisModel.ProductName == "" || thisModel.ProductName == "Get more online..." {
			return
		}

		// Validate model specific options
		for _, val := range modelOptions.Objects {
			// Validate entries
			if v, ok := val.(*widget.Entry); ok {
				if err := v.Validate(); err != nil {
					log.Println("Couldn't submit form. Reason:", err)
					return
				}
			}

			// Validate selects
			if v, ok := val.(*widget.Select); ok {
				if v.Selected == "" {
					log.Println("Can't save an empty select value")
					return
				}
			}
		}

		saveThisModel()
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

func saveThisModel() error {
	// Make the directory to save all models
	saveDirPath := path.Join(WD, "myModels", time.Now().Format(tFormat))
	err := os.Mkdir(saveDirPath, 0766)
	if err != nil {
		return err
	}

	// Marshall model
	b, err := json.Marshal(thisModel)
	if err != nil {
		return err
	}

	// Save config file
	err = os.WriteFile(path.Join(saveDirPath, "config.json"), b, 0766)

	return err
}
