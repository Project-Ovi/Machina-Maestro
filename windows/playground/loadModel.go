package playground

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

var WD string

func loadModel(dirName string) {
	// Get working directory
	var err error
	WD, err = os.Getwd()
	if err != nil {
		log.Panic("Failed to get working directory:", err)
	}

	// Create dir path
	dirPath := path.Join(WD, "myModels", dirName)

	// Read config file
	f, err := os.ReadFile(path.Join(dirPath, "config.json"))
	if err != nil {
		log.Panic("Failed to read config file:", err)
	}

	// Unmarshal config file
	if err := json.Unmarshal(f, &thisModel); err != nil {
		log.Panic("Failed to unmarshal config file:", err)
	}
}
