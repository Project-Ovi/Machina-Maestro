package helper

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type StackWithFooter struct{}

func (l *StackWithFooter) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// Make sure we have enough elements
	if len(objects) < 1 {
		return fyne.NewSize(0, 0)
	}

	// Variables
	navbarWidth := objects[0].MinSize().Width
	navbarHeight := objects[0].MinSize().Height
	w, h := float32(0), float32(0)

	for _, val := range objects {
		w = max(w, val.MinSize().Width)
		h = max(h, val.MinSize().Height)
	}

	return fyne.NewSize(max(navbarWidth, w), navbarHeight+h)
}

func (l *StackWithFooter) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	// Position footer
	objects[0].Resize(fyne.NewSize(containerSize.Width, objects[0].MinSize().Height))
	objects[0].Move(fyne.NewPos(0, containerSize.Height-objects[0].MinSize().Height))

	// Position other elements
	for i, val := range objects {
		// Skip footer
		if i == 0 {
			continue
		}

		// Stack them here
		val.Resize(fyne.NewSize(containerSize.Width, containerSize.Height-objects[0].MinSize().Height-theme.LineSpacing()))
		val.Move(fyne.NewPos(0, 0))
	}
}
