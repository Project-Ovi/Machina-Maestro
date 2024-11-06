package playground

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	ovimk2machinamaestroplugin "github.com/Project-Ovi/ovi-mk2-machina-maestro-plugin"
)

func FormSelector(s string, form *fyne.Container) {
	log.Println("Creating form for model", s+"...")
	switch s {
	case "OVI MK2":
		ovimk2machinamaestroplugin.Form(form)
	}
}
func SaveSelector(form *fyne.Container) string {
	s := form.Objects[3].(*widget.Select).Selected
	log.Println("Saving model", s+"...")
	switch s {
	case "OVI MK2":
		str_err, model := ovimk2machinamaestroplugin.Save(form, workingDirectory)
		ThisModel = model

		return str_err
	}

	return "Please select a valid model"
}
func LoadSelector() {
	log.Println("Loading model", ThisModel.Model+"...")
	switch ThisModel.Model {
	case "OVI MK2":
		ovimk2machinamaestroplugin.Load()
	}
}

var SupportedModels []string = []string{
	"OVI MK2",
}
