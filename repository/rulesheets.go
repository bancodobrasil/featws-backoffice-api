package repository

import (
	"context"
	"errors"

	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/featws-api/models"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// Rulesheets ...
type Rulesheets struct {
	model *gorm.DB
}

var instanceRulesheets = Rulesheets{}

// GetRulesheetsRepository ...
func GetRulesheetsRepository() Rulesheets {
	if instanceRulesheets.model == nil {
		database.GetConn().AutoMigrate(&models.Rulesheet{})
		instanceRulesheets.model = database.GetModel(&models.Rulesheet{})
	}
	return instanceRulesheets
}

// Create ...
func (r Rulesheets) Create(ctx context.Context, rulesheet *models.Rulesheet) error {

	result := r.model.Create(&rulesheet)
	if result.Error != nil {
		log.Errorf("error on insert the result into model: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected != 1 {
		err := errors.New("error on insert not inserted")
		log.Error(err)
		return err
	}

	//rulesheet.ID = result. InsertedID.(primitive.ObjectID)

	return nil
}

// Find ...
func (r Rulesheets) Find(ctx context.Context, filter interface{}) (list []*models.Rulesheet, err error) {

	result := r.model.Find(&list)

	err = result.Error
	if err != nil {
		log.Errorf("Error on find: %v", err)
		return
	}

	return
}

// Get ...
func (r Rulesheets) Get(ctx context.Context, id string) (rulesheet *models.Rulesheet, err error) {

	result := r.model.First(&rulesheet, "id = ? or name = ?", id, id)

	err = result.Error
	if err != nil {
		log.Errorf("Error on find one result into collection: %v", err)
		return
	}

	return
}

// Update ...
func (r Rulesheets) Update(ctx context.Context, entity models.Rulesheet) (updated *models.Rulesheet, err error) {

	result := r.model.Save(&entity)

	err = result.Error
	if err != nil {
		log.Errorf("Error on update into collection: %v", err)
		return
	}

	updated = &entity

	return
}

// Delete ...
func (r Rulesheets) Delete(ctx context.Context, id string) (deleted bool, err error) {

	result := r.model.Delete(id)

	err = result.Error
	if err != nil {
		log.Errorf("Error on delete from collection: %v", err)
		return
	}

	deleted = result.RowsAffected == 1

	return
}
