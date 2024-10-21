package helper

// * Model definitions
type Model struct {
	Name    string            `json:"name"`
	Model   string            `json:"model_name"`
	Website string            `json:"website"`
	Other   map[string]string `json:"other"`
}

// * Action definitions
type Action struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Commands    []Command `json:"commands"`
	running     bool
}

type Command struct {
	DisplayName string     `json:"name"`
	Arguments   []Argument `json:"args"`
	f           func([]Argument) error
}

type Argument struct {
	Name    string      `json:"name"`
	ArgType string      `json:"type"`
	Value   interface{} `json:"val"`
}
