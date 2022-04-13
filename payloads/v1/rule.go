package v1

// Rule contains all input for rule execution
// type Rule struct {
// 	Context map[string]interface{} `json:"context"`
// 	Load    []string               `json:"load"`
// }
type Rule struct {
	ID      string                 `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name    string                 `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Type    string                 `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Options map[string]interface{} `json:"options,omitempty" validate:"required" mapstructure:"options"`
}
