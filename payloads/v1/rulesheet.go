package v1

// Rulesheet contains all input for rulesheet execution
// type Rulesheet struct {
// 	Context map[string]interface{} `json:"context"`
// 	Load    []string               `json:"load"`
// }
type Rulesheet struct {
	ID            uint                    `json:"id,omitempty"`
	Name          string                  `json:"name,omitempty" validate:"required"`
	Description   string                  `json:"description,omitempty"`
	Slug          string                  `json:"slug,omitempty"`
	Version       string                  `json:"version,omitempty"`
	HasStringRule bool                    `json:"hasStringRule,omitempty"`
	Features      *[]interface{}          `json:"features,omitempty"`
	Parameters    *[]interface{}          `json:"parameters,omitempty"`
	Rules         *map[string]interface{} `json:"rules,omitempty"`
}
