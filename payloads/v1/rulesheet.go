package v1

// Rulesheet contains all input for rulesheet execution
// type Rulesheet struct {
// 	Context map[string]interface{} `json:"context"`
// 	Load    []string               `json:"load"`
// }
type Rulesheet struct {
	ID   string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name string `json:"name,omitempty" validate:"required" mapstructure:"name"`
}
