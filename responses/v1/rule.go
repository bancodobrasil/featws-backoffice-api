package v1

import "github.com/bancodobrasil/featws-api/models"

// Rule ...
type Rule struct {
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"name,omitempty"`
	Type    string                 `json:"type,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// NewRule ...
func NewRule(entity models.Rule) Rule {
	return Rule{
		ID:      entity.ID.Hex(),
		Name:    entity.Name,
		Type:    entity.Type,
		Options: entity.Options,
	}
}
