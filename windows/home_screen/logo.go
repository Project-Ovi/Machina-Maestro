package homescreen

import (
	"image"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const brokenSize = 120

func brokenIMG() *canvas.Image {
	imgObj := canvas.NewImageFromResource(theme.BrokenImageIcon())
	imgObj.FillMode = canvas.ImageFillContain
	imgObj.SetMinSize(fyne.NewSize(brokenSize, brokenSize))

	return imgObj
}

func getLogo() *canvas.Image {
	log.Println("Started loading splash screen backgorund image")
	// Read current directory
	dir, err := os.ReadDir("./")
	if err != nil {
		log.Println("Failed to read current dir: ", err)
		return brokenIMG()
	}

	log.Println("Looping through files")
	// Loop through files
	for _, val := range dir {
		// Discard directories
		if val.IsDir() {
			log.Println("Discarded", val.Name()+".", "Reason: IsDir == true")
			continue
		}

		// Load image
		this := strings.Split(val.Name(), ".")
		if len(this) == 2 && this[0] == "logo" {
			// Read file
			f, err := os.Open(val.Name())
			if err != nil {
				log.Println("Failed to read file: ", err)
				continue
			}
			defer f.Close()

			// Decode file
			img, _, err := image.Decode(f)
			if err != nil {
				log.Println("Failed to decode file: ", err)
				continue
			} else {
				log.Println("Loaded file", val.Name())
			}

			// Convert file
			imgObj := canvas.NewImageFromImage(img)
			imgObj.FillMode = canvas.ImageFillContain
			return imgObj
		}
	}
	return brokenIMG()
}
