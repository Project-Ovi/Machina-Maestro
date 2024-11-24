package playground

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func productPage(content *fyne.Container) {
	//! Make loading page
	// Make title
	titleObj := canvas.NewText("Loading...", theme.Color(theme.ColorNameForeground))
	titleObj.TextSize = theme.Size(theme.SizeNameHeadingText)
	titleObj.Alignment = fyne.TextAlignCenter
	titleObj.TextStyle.Bold = true

	// Make subtitle
	subtitleObj := canvas.NewText("Please be patient while the requested page is being downloaded...", theme.Color(theme.ColorNameForeground))
	subtitleObj.Alignment = fyne.TextAlignCenter

	// Make loading bar
	loadingBar := widget.NewProgressBarInfinite()
	loadingBar.Start()

	// Display the loading
	content.Layout = layout.NewVBoxLayout()
	content.Objects = []fyne.CanvasObject{
		layout.NewSpacer(),
		titleObj,
		subtitleObj,
		loadingBar,
		layout.NewSpacer(),
	}
	content.Refresh()

	//! Display product page
	// Build url to markdown file
	mdFileURL, err := url.Parse(thisModel.MarkdownRefrenceURL)
	if err != nil {
		//TODO
	}

	// Make a HTTP GET request
	log.Println("Making a HTPP GET request to:", mdFileURL.String())
	resp, err := http.Get(mdFileURL.String())
	if err != nil {
		//TODO
	}
	defer resp.Body.Close()

	// Read response body
	markdownText, err := io.ReadAll(resp.Body)
	if err != nil {
		//TODO
	}

	// Create a widget to display the markdown text
	mdObj := widget.NewRichTextFromMarkdown(string(markdownText))

	// Display the text
	content.Layout = layout.NewStackLayout()
	content.Objects = []fyne.CanvasObject{mdObj}
	content.Refresh()
}
