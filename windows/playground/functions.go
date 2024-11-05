package playground

import (
	"encoding/json"
	"os"
	"path"
)

func saveAllActions() error {
	// Marshall
	b, err := json.Marshal(actionCollection)
	if err != nil {
		return err
	}

	// Save to file
	return os.WriteFile(path.Join(workingDirectory, "/myModels/", thisModel.Name, "/actions.json"), b, os.ModePerm)
}
