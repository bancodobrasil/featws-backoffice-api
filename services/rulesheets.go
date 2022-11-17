package services

import (
	"context"
	"fmt"

	"github.com/bancodobrasil/featws-api/dtos"
	"github.com/bancodobrasil/featws-api/models"
	"github.com/bancodobrasil/featws-api/repository"
	"github.com/gosimple/slug"
	log "github.com/sirupsen/logrus"
)

// FindOptions ...
type FindOptions struct {
	Limit int
	Page  int
}

// Rulesheets ...
type Rulesheets interface {
	Create(context.Context, *dtos.Rulesheet) error
	Find(ctx context.Context, filter interface{}, options *FindOptions) ([]*dtos.Rulesheet, error)
	Count(ctx context.Context, entity interface{}) (count int64, err error)
	Get(ctx context.Context, id string) (*dtos.Rulesheet, error)
	Update(ctx context.Context, entity dtos.Rulesheet) (*dtos.Rulesheet, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type rulesheets struct {
	gitlabService Gitlab
	repository    repository.Rulesheets
}

// NewRulesheets ...
func NewRulesheets(repository repository.Rulesheets, gitlabService Gitlab) Rulesheets {
	return rulesheets{
		repository:    repository,
		gitlabService: gitlabService,
	}
}

// CreateRulesheet ...
func (rs rulesheets) Create(ctx context.Context, rulesheetDTO *dtos.Rulesheet) (err error) {

	rulesheet, _ := models.NewRulesheetV1(*rulesheetDTO)
	if rulesheet.Slug == "" {
		rulesheet.Slug = slug.Make(rulesheet.Name)
	}

	fmt.Print(rulesheet.Slug)

	err = rs.repository.Create(ctx, &rulesheet)
	if err != nil {
		log.Errorf("Error on create rulesheet into repository: %v", err)
		return
	}
	rulesheetDTO.ID = rulesheet.ID
	rulesheetDTO.Slug = rulesheet.Slug
	err = rs.gitlabService.Save(rulesheetDTO, "[FEATWS BOT] Create Repo")
	if err != nil {
		log.Errorf("Error on save rulesheet into repository: %v", err)
		return
	}

	err = rs.gitlabService.Fill(rulesheetDTO)
	if err != nil {
		log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
		return
	}

	return
}

// FetchRulesheets ...
func (rs rulesheets) Find(ctx context.Context, filter interface{}, options *FindOptions) (result []*dtos.Rulesheet, err error) {

	var opts *repository.FindOptions = nil

	if options != nil {
		opts = &repository.FindOptions{
			Limit: options.Limit,
			Page:  options.Page,
		}
	}

	entities, err := rs.repository.Find(ctx, filter, opts)
	if err != nil {
		log.Errorf("Error on fetch the rulesheets(find): %v", err)
		return
	}

	result = make([]*dtos.Rulesheet, 0)

	for _, entity := range entities {
		result = append(result, newRulesheetDTO(entity))
	}

	return
}

// Count ...
func (rs rulesheets) Count(ctx context.Context, entity interface{}) (count int64, err error) {

	count, err = rs.repository.Count(ctx, entity)
	if err != nil {
		log.Errorf("Error on count the entities(find): %v", err)
		return
	}

	return
}

// FetchRulesheet ...
func (rs rulesheets) Get(ctx context.Context, id string) (result *dtos.Rulesheet, err error) {

	entity, err := rs.repository.Get(ctx, id)
	if err != nil {
		log.Errorf("Error on fetch rulesheet(get): %v", err)
		return
	}

	result = newRulesheetDTO(entity)

	if result != nil {
		err = rs.gitlabService.Fill(result)
		if err != nil {
			log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
			return
		}
	}

	return
}

// UpdateRulesheet ...
func (rs rulesheets) Update(ctx context.Context, rulesheetDTO dtos.Rulesheet) (result *dtos.Rulesheet, err error) {

	entity, _ := models.NewRulesheetV1(rulesheetDTO)

	_, err = rs.repository.Update(ctx, entity)
	if err != nil {
		log.Errorf("Error on update the rulesheet from repository: %v", err)
		return
	}

	err = rs.gitlabService.Save(&rulesheetDTO, "[FEATWS BOT] Update Repo")
	if err != nil {
		log.Errorf("Error on save the rulesheet into repository: %v", err)
		return
	}

	result = &rulesheetDTO

	return
}

func (rs rulesheets) Delete(ctx context.Context, id string) (bool, error) {

	db := rs.repository.GetDB()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// get the specific rulesheet
	rulesheet, err := rs.repository.Get(ctx, id)
	if err != nil {
		tx.Rollback()
		log.Errorf("Error on fetch rulesheet(get): %v", err)
		return false, err
	}

	// update the ruleshet name to deleted
	rulesheet.Name = fmt.Sprintf("%s-deleted-%v", rulesheet.Name, rulesheet.ID)

	// update the rulesheet
	_, err = rs.repository.Update(ctx, *rulesheet)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	_, err = rs.repository.DeleteInTransaction(ctx, db, id)
	if err != nil {
		tx.Rollback()
		log.Errorf("Error on delete the rulesheet from repository: %v", err)
		return false, err
	}

	return true, tx.Commit().Error
}

func newRulesheetDTO(entity *models.Rulesheet) *dtos.Rulesheet {
	return &dtos.Rulesheet{
		ID:            entity.ID,
		Name:          entity.Name,
		Description:   entity.Description,
		Slug:          entity.Slug,
		HasStringRule: entity.HasStringRule,
	}
}
