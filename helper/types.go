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
	Running     bool
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

func (a *Action) Run(defaultCommands []Command) error {
	// Fix function
	a.Fix(defaultCommands)

	// Run
	for _, val := range (*a).Commands {
		if err := val.f(val.Arguments); err != nil {
			return err
		}
	}

	// Return
	return nil
}

func (a *Action) Fix(defaultCommands []Command) {
	for i, val1 := range (*a).Commands {
		for _, val2 := range defaultCommands {
			if val1.DisplayName == val2.DisplayName {
				(*a).Commands[i].f = val2.f
			}
		}
	}
}
