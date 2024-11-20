package helper

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type StackWithSidebar struct{}

func (l *StackWithSidebar) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// Make variables
	sidebarMinSize := objects[0].MinSize()
	stackMinSize := fyne.NewSize(0, 0)

	// Calculate stack min size
	for _, val := range objects[1:] {
		stackMinSize = fyne.NewSize(max(stackMinSize.Width, val.MinSize().Width), max(stackMinSize.Height, val.MinSize().Height))
	}

	// Return
	return fyne.NewSize(sidebarMinSize.Width+stackMinSize.Width+theme.Size(theme.SizeNameLineSpacing), max(sidebarMinSize.Height, stackMinSize.Height))
}

func (l *StackWithSidebar) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	// Position sidebar
	objects[0].Resize(fyne.NewSize(objects[0].MinSize().Width, containerSize.Height))
	objects[0].Move(fyne.NewPos(0, 0))

	// Stack other elements
	for _, val := range objects[1:] {
		val.Resize(fyne.NewSize(containerSize.Width-objects[0].MinSize().Width-theme.Size(theme.SizeNameLineSpacing), containerSize.Height))
		val.Move(fyne.NewPos(objects[0].MinSize().Width+theme.Size(theme.SizeNameLineSpacing), 0))
	}
}
