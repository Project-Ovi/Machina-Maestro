package helper

import (
	"path"

	"fyne.io/fyne/v2/canvas"
)

var ICON_robot_hand_round *canvas.Image

func LoadIcons(pathToIcons string) {
	ICON_robot_hand_round = canvas.NewImageFromFile(path.Join(pathToIcons, "robot_hand_round.png"))
}
