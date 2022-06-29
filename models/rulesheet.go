package models

import (
	"time"

	"github.com/bancodobrasil/featws-api/dtos"
	"gorm.io/gorm"
)

// Rulesheet ...
type Rulesheet struct {
	gorm.Model
	Name          string
	Description   string
	HasStringRule bool
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// NewRulesheetV1 ...
func NewRulesheetV1(dto dtos.Rulesheet) (entity Rulesheet, err error) {

	entity = Rulesheet{
		Model: gorm.Model{
			ID: dto.ID,
		},
		Name:        dto.Name,
		Description: dto.Description,
	}

	return
}
