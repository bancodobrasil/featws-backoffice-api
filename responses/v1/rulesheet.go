package v1

import "github.com/bancodobrasil/featws-api/models"

// Rulesheet ...
type Rulesheet struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// NewRulesheet ...
func NewRulesheet(entity models.Rulesheet) Rulesheet {
	return Rulesheet{
		ID:   entity.ID.Hex(),
		Name: entity.Name,
	}
}
