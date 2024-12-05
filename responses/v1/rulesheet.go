package v1

import "github.com/bancodobrasil/featws-api/dtos"

// Rulesheet type is a struct that contains various fields related to a set of rules, including its
// ID, name, description, slug, version, features, parameters, and rules.
//
// Property:
//
//   - FindResult: This is an embedded struct that contains fields related to the result of a search operation.
//   - ID: is an unsigned integer that represents the unique identifier of a Rulesheet.
//   - Name: The name of the rulesheet.
//   - Description - The `Rulesheet` struct is a data structure in Go programming language that represents a set of rules for a system or application. It contains various properties such as `ID`, `Name`, `Description`, `Slug`, `Version`, `Features`, `Parameters`, and `Rules`.
//   - Slug: is a string property in the Rulesheet struct that represents a unique identifier for the rulesheet. It is typically used in URLs to identify and access a specific rulesheet.
//   - Version: Represents the version number of the rulesheet.
//   - Features: It is a pointer to a slice of maps that represent the features of the rulesheet. Each map contains key-value pairs, where the key represents the feature name as a string, and the value is an interface that allows for flexibility in defining different types of feature values.
//   - Parameters: It is a pointer to a slice of maps, where each map represents a parameter used in the rules defined within the "Rules" property. Each map consists of key-value pairs, with the key being a string representing the parameter name, and the value being an interface. This design allows for flexibility in defining various types of parameter values.
//   - Rules: a pointer to a map of string keys and interface values. This is likely where the actual rules for the rulesheet are stored. The keys in the map would likely correspond to some sort of rule identifier or name, and the values would contain the logic or conditions for.
type Rulesheet struct {
	FindResult
	ID             uint                      `json:"id,omitempty"`
	Name           string                    `json:"name,omitempty"`
	Description    string                    `json:"description,omitempty"`
	Slug           string                    `json:"slug,omitempty"`
	Version        string                    `json:"version,omitempty"`
	PipelineStatus string                    `json:"status,omitempty"`
	WebURL         string                    `json:"web_url,omitempty"`
	Features       *[]map[string]interface{} `json:"features,omitempty"`
	Parameters     *[]map[string]interface{} `json:"parameters,omitempty"`
	Rules          *map[string]interface{}   `json:"rules,omitempty"`
}

// NewRulesheet creates a new Rulesheet object by copying data from a DTO object.
func NewRulesheet(dto *dtos.Rulesheet) Rulesheet {
	return Rulesheet{
		ID:             dto.ID,
		Name:           dto.Name,
		Description:    dto.Description,
		Slug:           dto.Slug,
		Version:        dto.Version,
		PipelineStatus: dto.PipelineStatus,
		WebURL:         dto.WebURL,
		Features:       dto.Features,
		Parameters:     dto.Parameters,
		Rules:          dto.Rules,
	}
}
