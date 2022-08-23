package v1

import "github.com/bancodobrasil/featws-api/dtos"

// Rulesheet ...
type Rulesheet struct {
	FindResult
	ID          uint                    `json:"id,omitempty"`
	Name        string                  `json:"name,omitempty"`
	Description string                  `json:"description,omitempty"`
	Slug        string                  `json:"slug,omitempty"`
	Version     string                  `json:"version,omitempty"`
	Features    *[]interface{}          `json:"features,omitempty"`
	Parameters  *[]interface{}          `json:"parameters,omitempty"`
	Rules       *map[string]interface{} `json:"rules,omitempty"`
}

// NewRulesheet ...
func NewRulesheet(dto *dtos.Rulesheet) Rulesheet {
	return Rulesheet{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Slug:        dto.Slug,
		Version:     dto.Version,
		Features:    dto.Features,
		Parameters:  dto.Parameters,
		Rules:       dto.Rules,
	}
}
