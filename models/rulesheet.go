package models

import (
	"time"

	"github.com/bancodobrasil/featws-api/dtos"
	"gorm.io/gorm"
)

// Rulesheet represents a model with fields for name, description, slug, boolean flag, and timestamps.
//
// Property:
//   - `gorm.Model`: This is a struct that provides some common fields for db models such as `ID`, `CreatedAt`, `UpdatedAt`, and `DeletedAt`.
//   - Name: is a string that represents the name of a rulesheet. The maximum length of 255 characters and is indexed as unique, two rulesheets can't have the same name.
//   - Description: provides additional information or details about the Rulesheet. It can be used to describe the purpose or function of the Rulesheet, or any other relevant information that may be useful to users or developers.
//   - Slug: a unique identifier for the Rulesheet. It is typically a short, human-readable string that is used in URLs.
//   - HasStringRule: a boolean property that indicates whether the Rulesheet has a string rule or not. It is likely used in the logic of the application to determine how to handle the Rulesheet object.
//   - CreatedAt: represents the timestamp of when the Rulesheet was created. It is of type *time.Time, which is a pointer to a time. This property is automatically set by the GORM library when a new Rules.
//   - UpdatedAt: represents the timestamp of the last time the `Rulesheet` was updated in the database. This property is useful for tracking when a `Rulesheet` was last modified and can be used in various ways within the application logic.
type Rulesheet struct {
	gorm.Model
	Name          string `gorm:"type:varchar(255);uniqueIndex"`
	Description   string
	Slug          string `gorm:"unique_index"`
	HasStringRule bool
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// NewRulesheetV1 creates a new Rulesheet entity from a DTO in Go.
func NewRulesheetV1(dto dtos.Rulesheet) (entity Rulesheet, err error) {

	entity = Rulesheet{
		Model: gorm.Model{
			ID: dto.ID,
		},
		Name:        dto.Name,
		Description: dto.Description,
		Slug:        dto.Slug,
	}

	return
}
