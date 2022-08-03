package services

import (
	"context"

	"github.com/bancodobrasil/featws-api/dtos"
	"github.com/bancodobrasil/featws-api/models"
	"github.com/bancodobrasil/featws-api/repository"
	log "github.com/sirupsen/logrus"
)

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

	//TODO verifica unicidade do nome
	rulesheet, err := models.NewRulesheetV1(*rulesheetDTO)
	if err != nil {
		log.Errorf("Error on create rulesheet on create model: %v", err)
		return
	}

	err = rs.repository.Create(ctx, &rulesheet)
	if err != nil {
		log.Errorf("Error on create rulesheet into repository: %v", err)
		return
	}

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

	// for _, rulesheet := range result {
	// 	err = fillWithGitlab(rulesheet)
	// 	if err != nil {
	// 		return
	// 	}
	// }

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

	entity, err := models.NewRulesheetV1(rulesheetDTO)
	if err != nil {
		log.Errorf("Error on create rulesheet on create model: %v", err)
		return
	}

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

	// err = rs.gitlabService.Fill(&rulesheetDTO)
	// if err != nil {
	// 	log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
	// 	return
	// }

	result = &rulesheetDTO

	return
}

// DeleteRulesheet ...
func (rs rulesheets) Delete(ctx context.Context, id string) (deleted bool, err error) {

	deleted, err = rs.repository.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error on delete the rulesheet from repository: %v", err)
		return
	}

	return
}

func newRulesheetDTO(entity *models.Rulesheet) *dtos.Rulesheet {
	return &dtos.Rulesheet{
		ID:            entity.ID,
		Name:          entity.Name,
		Description:   entity.Description,
		HasStringRule: entity.HasStringRule,
	}
}
