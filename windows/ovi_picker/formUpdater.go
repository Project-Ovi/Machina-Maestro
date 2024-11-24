package ovipicker

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

func ReturnModelSpecificForm(name string, internallyCalled bool) []fyne.CanvasObject {
	// Read the options file
	f, err := os.ReadFile(path.Join(WD, "models", name, "options.json"))
	if err != nil {
		log.Println("Failed to open options file for", name, err)
		return nil
	}

	// Decode file
	var options map[string]interface{}
	err = json.Unmarshal(f, &options)
	if err != nil {
		log.Println("Failed to decode options file for", name, err)
		return nil
	}

	// Create elements
	var toReturn []fyne.CanvasObject
	for key, value := range options {
		// Get option title
		objTitle := widget.NewLabel(key)

		// Try to make an entry
		var objEntry fyne.CanvasObject
		if v, ok := value.(string); ok {
			objEntry = widget.NewEntry()
			objEntry.(*widget.Entry).Validator = validation.NewRegexp(v, "Invalid text")
			if !internallyCalled {
				objEntry.(*widget.Entry).OnChanged = func(s string) {
					thisModel.Others[key] = s
				}
			}
		}

		// Try to make a select
		if v, ok := value.([]interface{}); ok {
			selectOptions := make([]string, len(v))
			for i, val := range v {
				if vali, ok := val.(string); ok {
					selectOptions[i] = vali
				} else {
					log.Println("Unexpected select option:", val)
				}
			}

			objEntry = widget.NewSelect(selectOptions, func(s string) {})
			if !internallyCalled {
				objEntry.(*widget.Entry).OnChanged = func(s string) {
					thisModel.Others[key] = s
				}
			}
		}

		// Make sure we have at leas one entry
		if objEntry == nil {
			log.Println("Invalid option entry:", key)
			log.Println("Value of", key, ":", reflect.TypeOf(value))
			continue
		}

		// Append elements
		toReturn = append(toReturn, objTitle, objEntry)
	}

	return toReturn
}
