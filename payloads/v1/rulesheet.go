package v1

// Rulesheet contains all input for rulesheet execution
// type Rulesheet struct {
// 	Context map[string]interface{} `json:"context"`
// 	Load    []string               `json:"load"`
// }
type Rulesheet struct {
	ID         string             `json:"id,omitempty"`
	Name       string             `json:"name,omitempty" validate:"required"`
	Version    string             `json:"version,omitempty"`
	Features   *[]interface{}     `json:"features,omitempty"`
	Parameters *[]interface{}     `json:"parameters,omitempty"`
	Rules      *map[string]string `json:"rules,omitempty"`
}
