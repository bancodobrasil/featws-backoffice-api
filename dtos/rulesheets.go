package dtos

import (
	"encoding/json"

	v1 "github.com/bancodobrasil/featws-api/payloads/v1"
)

// Rule ...
type Rule struct {
	Condition string      `json:"condition,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	Type      string      `json:"type,omitempty"`
}

// Rulesheet ...
type Rulesheet struct {
	ID            uint
	Name          string
	Description   string
	Slug          string
	HasStringRule bool
	Version       string
	Features      *[]map[string]interface{}
	Parameters    *[]map[string]interface{}
	Rules         *map[string]interface{}
}

// NewRulesheetV1 ...
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
