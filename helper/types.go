package helper

type OVI struct {
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	ProductName         string            `json:"product"`
	MarkdownRefrenceURL string            `json:"mdsite"`
	Others              map[string]string `json:"modelSpecific"`
}

type BuiltinFunction struct {
	LuaFileName string            `json:"File"`
	Name        string            `json:"Name"`
	DisplayName string            `json:"Display"`
	Arguments   map[string]string `json:"Inputs"`
}

type ModelFunction struct {
	FunctionName string            `json:"fname"`
	Arguments    map[string]string `json:"arguments"`
}

type Action struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Functions   []ModelFunction `json:"functions"`
	Running     bool
}
