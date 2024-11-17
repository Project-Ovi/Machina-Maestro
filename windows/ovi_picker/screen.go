package ovipicker

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path"
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

var waitPressGroup sync.WaitGroup
var returnedText string

const tFormat = "20060102150405000"

func Launch(window fyne.Window) string {
	// Initialize wait group
	waitPressGroup.Add(1)

	// Make the navbar
	navbar := navbar()

	// Make scroll
	scrollable := container.NewVScroll(container.New(
		layout.NewVBoxLayout(),
	))

	// Update scrollable items
	go func() {
		// Get working directory
		WD, err := os.Getwd()
		if err != nil {
			log.Println(err)
			log.Println("Aborting loading other OVI models")
			return
		}

		// Read dir
		var dir []fs.DirEntry
		dirName := path.Join(WD, "myModels")
		dir, err = os.ReadDir(dirName)
		if err != nil {
			log.Println("Failed to open", dirName)
			err = os.Mkdir("myModels", 0666)
			if err != nil {
				log.Println("Failed to create", dirName)
				log.Println("Aborting loading other OVI models")
				return
			}
			dir, err = os.ReadDir(dirName)
			if err != nil {
				log.Println(err)
				log.Println("Aborting loading other OVI models")
				return
			}
		}

		// Loop through the dir elements
		for _, val := range dir {
			// Skip non-dirs
			if !val.IsDir() {
				log.Println(val.Name(), "is not a dir. Skipping")
				continue
			}

			// Initialize variables
			var name string
			var productName string
			var rawDirName string
			var dateCreated time.Time

			// Find the dirName
			rawDirName = val.Name()

			// Find the dateCreated
			dateCreated, err = time.Parse(tFormat, rawDirName)
			if err != nil {
				log.Println("Failed to parse: ", rawDirName)
				continue
			}

			// Open model's config file
			configLocation := path.Join(dirName, rawDirName, "config.json")
			b, err := os.ReadFile(configLocation)
			if err != nil {
				log.Println("Failed to read config file at", configLocation, ":", err)
			}

			// Decode file and extract name
			var modelData helper.OVI
			err = json.Unmarshal(b, &modelData)
			if err != nil {
				log.Println("Failed to decode JSON file at:", configLocation)
				continue
			}
			name = modelData.Name
			productName = modelData.ProductName

			// Make name text
			nameText := canvas.NewText(name, theme.Color(theme.ColorNameForeground))
			nameText.TextSize = theme.Size(theme.SizeNameText)

			// Make product name
			productNameText := canvas.NewText(productName, theme.Color(theme.ColorNameForeground))
			productNameText.TextSize = theme.Size(theme.SizeNameText)

			// Make date created text
			dateCreatedText := canvas.NewText(dateCreated.Format("02/01/2006"), theme.Color(theme.ColorNameForeground))
			dateCreatedText.TextSize = theme.Size(theme.SizeNameText)

			// Make delete button
			deleteBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDelete), func() {})
			deleteBTN.OnTapped = func() {
				// Disable button to prevent further pressing
				deleteBTN.Disable()

				// Delete instance
				err := os.RemoveAll(path.Join(dirName, rawDirName))
				if err != nil {
					log.Println("Failed to remove instance", rawDirName)
				}

				// Restart this window
				returnedText = "picker"
				waitPressGroup.Done()
			}

			// Make play button
			playBTN := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameMediaPlay), func() {})
			playBTN.OnTapped = func() {
				// Disable button to prevent further pressing
				playBTN.Disable()

				// Return the chosen model
				returnedText = rawDirName
				waitPressGroup.Done()
			}

			// Make content
			content := container.New(
				layout.NewHBoxLayout(),
				nameText,
				layout.NewSpacer(),
				productNameText,
				layout.NewSpacer(),
				dateCreatedText,
				layout.NewSpacer(),
				container.New(
					layout.NewHBoxLayout(),
					deleteBTN,
					playBTN,
				),
			)

			// Create model banner (add a bg to content)
			modelBanner := container.New(
				layout.NewStackLayout(),
				canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
				content,
			)

			// Append the object to the scrollable list
			scrollable.Content.(*fyne.Container).Add(modelBanner)

		}
	}()

	// Make content
	content := container.New(
		layout.NewBorderLayout(navbar, nil, nil, nil),
		scrollable,
	)

	// Set content
	window.SetContent(content)
	window.Show()

	waitPressGroup.Wait()
	return returnedText
}
