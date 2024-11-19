package ovipicker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
)

type modelPresentation struct {
	name       string
	link       string
	downloaded bool
}

const modelRepository = string("https://github.com/Project-Ovi/Maestro-Plugins")
const modelRepositoryBranch = string("main")

func fetchAvailableToDownloadModels() []modelPresentation {
	// Build link to model repository markdown file
	modelRepoMDOrigin, err := url.JoinPath(modelRepository, "refs", "heads", modelRepositoryBranch, "README.md")
	if err != nil {
		log.Println("Error in converting model add to model download. Failed to join URL: ", err)
		return nil
	}
	fmt.Println(modelRepoMDOrigin)

	// Make a HTTP GET request
	log.Println("Making HTTP GET request to:", modelRepoMDOrigin)
	httpGETResponse, err := http.Get(modelRepoMDOrigin)
	log.Println(httpGETResponse.Status)
	if err != nil {
		log.Println("Failed to make GET request:", err)
		return nil
	}
	defer httpGETResponse.Body.Close()

	// Read response body
	markdownText, err := io.ReadAll(httpGETResponse.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return nil
	}

	// Process response
	var availableModelsToDownload []modelPresentation
	for _, val := range strings.Split(string(markdownText), "\n") {
		fmt.Println(val)
		// Discard invalid items
		if len(val) < 1 || val[0] != '-' {
			continue
		}

		// Process
		fields := strings.SplitN(strings.TrimPrefix(val, "- "), " - ", 2)
		availableModelsToDownload = append(availableModelsToDownload, modelPresentation{
			name: fields[0],
			link: fields[1],
		})
	}

	return availableModelsToDownload
}

func convertModelAddToModelDownload(window fyne.Window) {
	// Fetch models
	modelOptions := fetchAvailableToDownloadModels()
	var userSelectedModelOptions []modelPresentation
	copy(userSelectedModelOptions, modelOptions)

	// Configure window
	window.SetTitle("Download OVIs")

	// Make a title
	titleObj := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.TextStyle.Bold = true
	titleObj.Alignment = fyne.TextAlignCenter

	// Make a list of available models
	modelListObjs := container.New(layout.NewVBoxLayout())

	// Append to that list
	for i, val := range modelOptions {
		// Make check
		checkObj := widget.NewCheck(val.name, func(b bool) {
			userSelectedModelOptions[i].downloaded = b
		})
		checkObj.Checked = val.downloaded

		// Append check
		modelListObjs.Add(checkObj)
	}

	// Add scroll functionality to the model list
	scrollModelsObj := container.NewVScroll(modelListObjs)

	// Assemble UI
	content := container.New(
		&helper.StackWithNavbar{},
		titleObj,
		scrollModelsObj,
	)

	// Display content
	window.SetContent(content)
}
