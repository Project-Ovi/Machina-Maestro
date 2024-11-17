package ovipicker

import (
	"io/fs"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
