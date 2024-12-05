package dtos

import (
	"encoding/json"

	v1 "github.com/bancodobrasil/featws-api/payloads/v1"
)

// Rule type is a struct with three fields: "Condition", "Value", and "Type", all of which are optional and have specific JSON tags.
//
// Property:
//
//   - Condition: is a string that represents the condition that needs to be met for the rule to be applied. It could be a comparison operator like ">", "<", ">=","<=", "=", or a logical operator like "AND", "OR", "NOT".
//   - Value: is an interface type, which means it can hold any type of value. It's used to store the value of a rule, which can be a string, number, boolean, or any other data type depending on the specific rule being defined.
//   - Type: is a string that specifies the type of the value. It can be used to indicate whether the value is a string, number, boolean, or any other data type. This property is optional and can be omitted.
type Rule struct {
	Condition string      `json:"condition,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	Dynamic   string      `json:"dynamic,omitempty"`
	Type      string      `json:"type,omitempty"`
}

// Rulesheet type represents a set of rules and parameters for a system, including features and a version number.
//
// Property:
//   - ID: is an unsigned integer that represents the unique identifier of a Rulesheet.
//   - Name: is a string property of the Rulesheet struct. It represents the name of the rulesheet.
//   - Description: is a string that provides a brief explanation or summary of what the Rulesheet is about. It can be used to give context to the Rulesheet and help users understand its purpose.
//   - Slug: is a string that represents a unique identifier for the Rulesheet. It is typically used in URLs to identify and access a specific Rulesheet.
//   - HasStringRule - HasStringRule is a boolean property that indicates whether or not the rulesheet contains a string rule. A string rule is a rule that involves comparing or manipulating strings.
//   - Version - The version of the rulesheet. It could be a string or a number that represents the version number of the rulesheet. This is useful for tracking changes and updates to the rulesheet over time.
//   - Features - Features is a pointer to a slice of maps, where each map represents a feature of the rulesheet. Each map contains key-value pairs where the key is a string representing the name of the feature and the value is an interface{} representing the value of the feature.
//   - Parameters: a pointer to a slice of maps, where each map represents a parameter that can be used in the rules defined in the `Rules` property. Each mapcontains key-value pairs where the key is a string representing the name of the parameter and the value is an interface{}
//   - Rules: property is a pointer to a map of string keys and interface values. This map represents the set of rules that are associated with the rulesheet. Each key in the map represents a unique rule identifier, and the corresponding value is an interface that can be usedto store any type of data. The use of `interface` allows for flexibility in the type of data that can be stored in the map.
type Rulesheet struct {
	ID             uint
	Name           string
	Description    string
	Slug           string
	HasStringRule  bool
	Version        string
	PipelineStatus string
	WebURL         string
	Features       *[]map[string]interface{}
	Parameters     *[]map[string]interface{}
	Rules          *map[string]interface{}
}

// NewRulesheetV1 takes in a payload of rulesheet and returns a DTO with the rules converted to a
// specific format.
func NewRulesheetV1(payload v1.Rulesheet) (dto Rulesheet, err error) {

	dto = Rulesheet{
		ID:          payload.ID,
		Name:        payload.Name,
		Description: payload.Description,
		Slug:        payload.Slug,
		Version:     payload.Version,
		Features:    payload.Features,
		Parameters:  payload.Parameters,
	}

	isRule := false

	if payload.Rules == nil {
		return
	}

	for _, v := range *payload.Rules {
		_, isString := v.(string)
		if !isString {
			isRule = true
		}
	}

	dto.HasStringRule = !isRule

	// FIXME - Remover restricao de exclusividade entre regras string e complexas
	if dto.HasStringRule {
		return
	}

	for k, v := range *payload.Rules {

		(*payload.Rules)[k], err = buildRule(v)
		if err != nil {
			return
		}
	}

	dto.Rules = payload.Rules

	return
}

// The function takes in an interface and recursively builds a rule based on its type.
func buildRule(v interface{}) (interface{}, error) {
	switch value := v.(type) {
	case []interface{}:
		//                                                                                                                                                                                                                                                                                                                                                                                                                                                                        fmt.Println("LIST", value)
		list := make([]interface{}, 0)
		for _, item := range value {
			itemRule, err := buildRule(item)
			if err != nil {
				return nil, err
			}
			list = append(list, itemRule)
		}
		return list, nil
	case map[string]interface{}:
		//fmt.Println("MAP INTERFACE", value)

		if _, ok := value["value"]; !ok {
			mapp := make(map[string]interface{}, 0)
			for k, item := range value {
				itemRule, err := buildRule(item)
				if err != nil {
					return nil, err
				}
				mapp[k] = itemRule
			}
			return mapp, nil
		}

		jsonData, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		regra := &Rule{}

		err = json.Unmarshal(jsonData, regra)
		if err != nil {
			return nil, err
		}

		return regra, nil

	default:
		// fmt.Errorf("DEFAULT [%v] {%t}", value, value)
		// return value, fmt.Errorf("DEFAULT [%v] {%t}", value, value)
		return value, nil
	}
}
