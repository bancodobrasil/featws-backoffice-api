package v1

// Rulesheet contains all input for rulesheet execution
// type Rulesheet struct {
// 	Context map[string]interface{} `json:"context"`
// 	Load    []string               `json:"load"`
// }

// Rulesheet defines a set of properties for a rulesheet object in Go, including its ID, name, description, slug, version, features, parameters, and rules.
//
// Property:
//   - ID: is an unsigned integer that represents the unique identifier of a Rulesheet.
//   - Name: The name of the rulesheet.
//   - Description - The `Rulesheet` struct is a data structure in Go programming language that represents a set of rules for a system or application. It contains various properties such as `ID`, `Name`, `Description`, `Slug`, `Version`, `Features`, `Parameters`, and `Rules`.
//   - Slug: is a string property in the Rulesheet struct that represents a unique identifier for the rulesheet. It is typically used in URLs to identify and access a specific rulesheet.
//   - Version: Represents the version number of the rulesheet.
//   - HasStringRule - The HasStringRule property is a boolean value that indicates whether the rulesheet contains a string rule or not. If it is true, it means that the rulesheet has at least one rule that involves a string value. If it is false, it means that all the rules in the rulesheet involve
//   - Features: It is a pointer to a slice of maps that represent the features of the rulesheet. Each map contains key-value pairs, where the key represents the feature name as a string, and the value is an interface that allows for flexibility in defining different types of feature values.
//   - Parameters: It is a pointer to a slice of maps, where each map represents a parameter used in the rules defined within the "Rules" property. Each map consists of key-value pairs, with the key being a string representing the parameter name, and the value being an interface. This design allows for flexibility in defining various types of parameter values.
//   - Rules: a pointer to a map of string keys and interface values. This is likely where the actual rules for the rulesheet are stored. The keys in the map would likely correspond to some sort of rule identifier or name, and the values would contain the logic or conditions for.
type Rulesheet struct {
	ID            uint                      `json:"id,omitempty"`
	Name          string                    `json:"name,omitempty" validate:"required"`
	Description   string                    `json:"description,omitempty"`
	Slug          string                    `json:"slug,omitempty"`
	Version       string                    `json:"version,omitempty"`
	HasStringRule bool                      `json:"hasStringRule,omitempty"`
	Features      *[]map[string]interface{} `json:"features,omitempty"`
	Parameters    *[]map[string]interface{} `json:"parameters,omitempty"`
	Rules         *map[string]interface{}   `json:"rules,omitempty"`
}
