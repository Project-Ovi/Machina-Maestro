package playground

import (
	"encoding/json"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/Project-Ovi/Machina-Maestro/windows/fatalerror"
)

func OVIPlayground(window fyne.Window) {
	// Reset defaults
	defaultCommands = []command{}

	// Fetch actions
	b, err := os.ReadFile(path.Join(workingDirectory, "/myModels/", thisModel.Name, "/actions.json"))
	if err != nil {
		fatalerror.Show(err, logger, MainWindow, App)
	}
	err = json.Unmarshal(b, &actionCollection)
	if err != nil {
		fatalerror.Show(err, logger, MainWindow, App)
	}

	// Load functions for this model
	go loadSelector()

	// Make a main content
	mainContent := container.New(layout.NewStackLayout())

	// Get navbar and sidebar
	sidebar := playgroundSidebar(&mainContent)
	navbar := playgroundNavbar(window, sidebar)

	// Display contents
	content := container.NewBorder(navbar, nil, sidebar, nil, mainContent)
	window.SetContent(content)
}
