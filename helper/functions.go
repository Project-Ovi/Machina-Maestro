package helper

import (
	"fmt"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func bigBTN(text string, extraPadding float32, size float32, f func()) *fyne.Container {
	btnText := canvas.NewText(text, theme.Color(theme.ColorNameForeground))
	btnText.TextSize = size
	btnBG := widget.NewButton(" ", f)
	btnBGHeightIncrements := 1
	for {
		desiredValue := fyne.MeasureText(btnText.Text, btnText.TextSize, btnText.TextStyle)
		actualValue := fyne.MeasureText(btnBG.Text, theme.TextSize(), btnText.TextStyle)
		if desiredValue.Width+extraPadding > actualValue.Width {
			btnBG.Text += " "
			btnBG.Refresh()
			continue
		}

		if desiredValue.Height > actualValue.Height*float32(btnBGHeightIncrements) {
			btnBG.Text += "\n"
			btnBGHeightIncrements++
			btnBG.Refresh()
			// continue
		}

		break
	}
	btn := container.New(layout.NewCenterLayout(), btnBG, btnText)
	return btn
}
func openExplorer(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// On Windows, use 'explorer'
		cmd = exec.Command("explorer", path)
	case "darwin":
		// On macOS, use 'open'
		cmd = exec.Command("open", path)
	case "linux":
		// On Linux, use 'xdg-open' or 'gio open' based on availability
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
