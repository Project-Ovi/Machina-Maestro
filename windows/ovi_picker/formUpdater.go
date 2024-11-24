package ovipicker

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

func ReturnModelSpecificForm(name string) []fyne.CanvasObject {
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
			if !isCalledFromOutside() {
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
			if !isCalledFromOutside() {
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

func isCalledFromOutside() bool {
	// Get the call stack
	pc := make([]uintptr, 10) // limit stack to 10 frames
	n := runtime.Callers(2, pc)
	if n == 0 {
		return false // no callers
	}

	frames := runtime.CallersFrames(pc[:n])
	mypkg := "github.com/Project-Ovi/Machina-Maestro/windows/ovi_picker"

	for {
		frame, more := frames.Next()
		// Check the package name in the function name
		if !strings.Contains(frame.Function, mypkg) {
			// The first frame not in `mypackage` indicates an external caller
			return true
		}
		if !more {
			break
		}
	}
	return false // All callers are from `mypackage`
}
