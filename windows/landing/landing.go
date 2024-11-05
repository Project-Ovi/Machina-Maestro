package landing

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/Project-Ovi/Machina-Maestro/helper"
	"golang.org/x/image/colornames"
)

func Launch(window fyne.Window, ICON_robot_hand_round *canvas.Image) {
	// Make title
	title := canvas.NewText("Machina Maestro", colornames.Orange)
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 32

	// Make icon
	ICON_robot_hand_round.SetMinSize(fyne.NewSquareSize(512))
	ICON_robot_hand_round.FillMode = canvas.ImageFillContain

	// Make button
	btn := helper.BigBTN("Select your OVI", 40, 20, func() {
		modelSelectWindow(window)
	})

	// Assemble UI
	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), title, layout.NewSpacer()),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), ICON_robot_hand_round, layout.NewSpacer()),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), btn, layout.NewSpacer()),
		layout.NewSpacer(),
	)
	window.SetContent(content)
	window.Show()
	window.CenterOnScreen()
}
