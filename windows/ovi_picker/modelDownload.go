package ovipicker

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Project-Ovi/Machina-Maestro/helper"
	"github.com/go-git/go-git/v5"
)

type modelPresentation struct {
	name       string
	link       string
	downloaded bool
}

const modelRepository = string("github.com/Project-Ovi/Maestro-Plugins")
const modelRepositoryBranch = string("main")

func fetchAvailableToDownloadModels() []modelPresentation {
	// Build link to model repository markdown file
	modelRepoMDOrigin, err := url.JoinPath("http://raw."+modelRepository, "refs", "heads", modelRepositoryBranch, "README.md")
	if err != nil {
		log.Println("Error in converting model add to model download. Failed to join URL: ", err)
		return nil
	}

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
		// fmt.Println(val)
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

	checkModelsDownloaded(availableModelsToDownload)
	return availableModelsToDownload
}

func checkModelsDownloaded(models []modelPresentation) {
	// Read directory
	downloadedDirs, err := os.ReadDir(path.Join(WD, "models"))
	if err != nil {
		log.Println(err)
		return
	}

	// Compare the models to the downloaded models
	for i, requested := range models {
		for _, available := range downloadedDirs {
			if requested.name == available.Name() && available.IsDir() {
				log.Println("Found downloaded model:", requested.name)
				models[i].downloaded = true
				break
			}
		}
	}
}

func waitingPage(window fyne.Window) {
	// Create title
	window.SetTitle("Please wait...")
	title := canvas.NewText(window.Title(), theme.Color(theme.ColorNameForeground))
	title.TextSize = theme.Size(theme.SizeNameHeadingText)
	title.Alignment = fyne.TextAlignCenter

	// Create progressbar
	progressbar := widget.NewProgressBarInfinite()
	progressbar.Start()

	// Create content
	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		title,
		progressbar,
		layout.NewSpacer(),
	)

	// Set window content
	window.SetContent(content)
}

func convertModelAddToModelDownload(window fyne.Window) {
	// Display waiting page
	waitingPage(window)

	// Fetch models
	modelOptions := fetchAvailableToDownloadModels()
	userSelectedModelOptions := make([]modelPresentation, len(modelOptions))
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

	// Create submit button
	submitBTN := widget.NewButtonWithIcon("Submit and download", theme.Icon(theme.IconNameDownload), func() {})
	submitBTNWrapper := container.New(
		layout.NewStackLayout(),
		submitBTN,
	)

	// Assemble UI
	content := container.New(
		&helper.StackWithNavbar{},
		titleObj,
		container.New(
			&helper.StackWithFooter{},
			submitBTNWrapper,
			container.New(
				layout.NewHBoxLayout(),
				layout.NewSpacer(),
				scrollModelsObj,
				layout.NewSpacer(),
			),
		),
	)

	submitBTN.OnTapped = func() {
		// Make a progress bar
		progressBarObj := widget.NewProgressBar()

		// Replace the button
		submitBTNWrapper.RemoveAll()
		submitBTNWrapper.Add(progressBarObj)
		content.Refresh()

		// Start loading
		for index, value := range userSelectedModelOptions {
			// Skip unchanged values
			if modelOptions[index].downloaded == value.downloaded {
				continue
			}

			// Manage models
			if !value.downloaded {
				err := os.RemoveAll(path.Join(WD, "models", value.name))
				if err != nil {
					log.Println("Failed to remove model", value.name, ". Reason:", err)
				}
			} else {
				// Clone repository
				git.PlainClone(path.Join(WD, "models", value.name), false, &git.CloneOptions{
					URL:      value.link,
					Progress: log.Writer(),
				})
			}

			// Update progress bar
			progressBarObj.SetValue(float64(index+1) / float64(len(userSelectedModelOptions)) * 100)
		}

		// Close window
		window.Close()
	}

	// Display content
	window.SetContent(content)
}
