package repository

import (
	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/featws-api/models"
	"gorm.io/gorm"
)

// Rulesheets ...
type Rulesheets interface {
	Repository[models.Rulesheet]
}

type rulesheets struct {
	repository[models.Rulesheet]
}

var instanceRulesheets Rulesheets

// GetRulesheets ...
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

func newRulesheets() (Rulesheets, error) {
	db := database.GetConn()
	return NewRulesheetsWithDB(db)
}

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
