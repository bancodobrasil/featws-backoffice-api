package repository

import (
	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/featws-api/models"
	"gorm.io/gorm"
)

// Rulesheets is defining an interface that has a single method signature `Repository[models.Rulesheet]` that is defined in repository.go.
type Rulesheets interface {
	Repository[models.Rulesheet]
}

// rulesheets contains an array of "Rulesheet" objects within a "repository" field.
//
// Property:
//   - repository: is a field of the `rulesheets` struct which is an array of `models.Rulesheet`. It's used to store multiple rulesheets in a single instance of the `rulesheets` struct.
type rulesheets struct {
	repository[models.Rulesheet]
}

var instanceRulesheets Rulesheets

// GetRulesheets returns an instance of the Rulesheets struct, creating it if it doesn't already exist.
func GetRulesheets() Rulesheets {
	if instanceRulesheets == nil {
		i, err := newRulesheets()
		if err != nil {
			panic(err)
		}
		instanceRulesheets = i
	}
	return instanceRulesheets
}

// newRulesheets creates a new instance of Rulesheets and returns it along with any errors encountered.
func newRulesheets() (Rulesheets, error) {
	db := database.GetConn()
	return NewRulesheetsWithDB(db)
}

// NewRulesheetsWithDB creates a new instance of Rulesheets with a given db connection and performs db migration.
func NewRulesheetsWithDB(db *gorm.DB) (Rulesheets, error) {
	err := db.AutoMigrate(&models.Rulesheet{})
	if err != nil {
		return nil, err
	}
	return &rulesheets{
		repository[models.Rulesheet]{
			db: db,
		},
	}, err
}
