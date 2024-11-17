package helper

type OVI struct {
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	ProductName         string            `json:"product"`
	MarkdownRefrenceURL string            `json:"mdsite"`
	Others              map[string]string `json:"modelSpecific"`
}
