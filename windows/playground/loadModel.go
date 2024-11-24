package playground

import (
	"encoding/json"
	"log"
	"os"
	"path"

	lua "github.com/yuin/gopher-lua"
)

var WD string
var LuaVM *lua.LState

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

	// Create the Lua VM
	LuaVM = lua.NewState()

	// Execute model-specific loader
	err = LuaVM.DoFile(path.Join(WD, "models", thisModel.ProductName, "loader.lua"))
	if err != nil {
		log.Println("Failed to run Lua loader for", thisModel.ProductName, ". Reason: ", err)
	}
}
