package playground

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/Project-Ovi/Machina-Maestro/helper"
	lua "github.com/yuin/gopher-lua"
)

var WD string
var LuaVM *lua.LState
var BuiltinFunctions []helper.BuiltinFunction

func loadModel(dirName string) {
	//! Misc
	// Get working directory
	var err error
	WD, err = os.Getwd()
	if err != nil {
		log.Panic("Failed to get working directory:", err)
	}

	// Create dir path
	dirPath := path.Join(WD, "myModels", dirName)

	//!----------------------------------------!
	//!Config file
	// Read config file
	f, err := os.ReadFile(path.Join(dirPath, "config.json"))
	if err != nil {
		log.Panic("Failed to read config file:", err)
	}

	// Unmarshal config file
	if err := json.Unmarshal(f, &thisModel); err != nil {
		log.Panic("Failed to unmarshal config file:", err)
	}

	//!Actions file
	// Read actions file
	f, err = os.ReadFile(path.Join(dirPath, "actions.json"))
	if err != nil {
		log.Println("Failed to read actions file:", err)
		log.Println("Attempting to create actions file")
		err = os.WriteFile(path.Join(dirPath, "config.json"), []byte("[]"), 0766)
		if err != nil {
			log.Panic(err)
		}
	} else {
		err = json.Unmarshal(f, &actionsCollection)
		if err != nil {
			log.Panic("Failed to unmarshal actions file:", err)
		}
	}

	//!----------------------------------------!
	//!Functions file
	// Read functions file
	f, err = os.ReadFile(path.Join(WD, "models", thisModel.ProductName, "functions", "functions.json"))
	if err != nil {
		log.Panic("Failed to read functions file:", err)
	}

	// Unmarshal functions file
	if err := json.Unmarshal(f, &BuiltinFunctions); err != nil {
		log.Panic(err)
	}

	// Cleanup function file paths
	for i, val := range BuiltinFunctions {
		BuiltinFunctions[i].LuaFileName = path.Join(WD, "models", thisModel.ProductName, "functions", val.LuaFileName)
	}

	//!Custom Lua loader
	// Create the Lua VM
	LuaVM = lua.NewState()

	// Execute model-specific loader
	err = LuaVM.DoFile(path.Join(WD, "models", thisModel.ProductName, "loader.lua"))
	if err != nil {
		log.Println("Failed to run Lua loader for", thisModel.ProductName, ". Reason: ", err)
	}
}
