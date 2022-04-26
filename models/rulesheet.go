package models

import (
	v1 "github.com/bancodobrasil/featws-api/payloads/v1"
	"gorm.io/gorm"
)

// Rulesheet ...
type Rulesheet struct {
	gorm.Model
	Name       string             
	Version    string             `gorm:"-"`
	Features   *[]interface{}     `gorm:"-"`
	Parameters *[]interface{}     `gorm:"-"`
	Rules      *map[string]string `gorm:"-"`
}

// NewRulesheetV1 ...
func NewRulesheetV1(payload v1.Rulesheet) (entity Rulesheet, err error) {

	entity = Rulesheet{
		Model: gorm.Model{
			ID: payload.ID,
		},
		Name:       payload.Name,
		Version:    payload.Version,
		Features:   payload.Features,
		Parameters: payload.Parameters,
		Rules:      payload.Rules,
	}
	return
}
