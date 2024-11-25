package playground

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"time"
)

func SaveAll() error {
	//! Save config
	// Marshal config
	b, err := json.Marshal(thisModel)
	if err != nil {
		return err
	}

	// Save to file
	err = os.WriteFile(path.Join(loadedModelDirPath, "config.json"), b, 0766)
	if err != nil {
		return err
	}

	//!Save actions
	// Marshal actions
	// Marshal config
	b, err = json.Marshal(actionsCollection)
	if err != nil {
		return err
	}

	// Save to file
	err = os.WriteFile(path.Join(loadedModelDirPath, "actions.json"), b, 0766)
	if err != nil {
		return err
	}

	return nil
}

func enableAutoSave(ms int) {
	go func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(ms))
			err := SaveAll()
			if err != nil {
				log.Println("Error while saving model:", err)
			}
		}
	}()
}
