package services

import (
	"context"

	"github.com/bancodobrasil/featws-api/models"
	"github.com/bancodobrasil/featws-api/repository"
	log "github.com/sirupsen/logrus"
)

type Rulesheets interface {
	Create(context.Context, *models.Rulesheet) error
	Find(ctx context.Context, filter interface{}) ([]*models.Rulesheet, error)
	Get(ctx context.Context, id string) (*models.Rulesheet, error)
	Update(ctx context.Context, entity models.Rulesheet) (*models.Rulesheet, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type rulesheets struct {
	gitlabService Gitlab
	repository    repository.Rulesheets
}

func NewRulesheets(repository repository.Rulesheets, gitlabService Gitlab) Rulesheets {
	return rulesheets{
		repository:    repository,
		gitlabService: gitlabService,
	}
}

// CreateRulesheet ...
func (rs rulesheets) Create(ctx context.Context, rulesheet *models.Rulesheet) (err error) {

	//TODO verifica unicidade do nome

	err = rs.repository.Create(ctx, rulesheet)
	if err != nil {
		log.Errorf("Error on create rulesheet into repository: %v", err)
		return
	}

	err = rs.gitlabService.Save(rulesheet, "[FEATWS BOT] Create Repo")
	if err != nil {
		log.Errorf("Error on save rulesheet into repository: %v", err)
		return
	}

	err = rs.gitlabService.Fill(rulesheet)
	if err != nil {
		log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
		return
	}

	return
}

// FetchRulesheets ...
func (rs rulesheets) Find(ctx context.Context, filter interface{}) (result []*models.Rulesheet, err error) {

	result, err = rs.repository.Find(ctx, filter)
	if err != nil {
		log.Errorf("Error on fetch the rulesheets(find): %v", err)
		return
	}

	// for _, rulesheet := range result {
	// 	err = fillWithGitlab(rulesheet)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	return
}

// FetchRulesheet ...
func (rs rulesheets) Get(ctx context.Context, id string) (result *models.Rulesheet, err error) {

	result, err = rs.repository.Get(ctx, id)
	if err != nil {
		log.Errorf("Error on fetch rulesheet(get): %v", err)
		return
	}

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
func (rs rulesheets) Update(ctx context.Context, entity models.Rulesheet) (result *models.Rulesheet, err error) {

	result, err = rs.repository.Update(ctx, entity)
	if err != nil {
		log.Errorf("Error on update the rulesheet from repository: %v", err)
		return
	}

	err = rs.gitlabService.Save(&entity, "[FEATWS BOT] Update Repo")
	if err != nil {
		log.Errorf("Error on save the rulesheet into repository: %v", err)
		return
	}

	err = rs.gitlabService.Fill(result)
	if err != nil {
		log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
		return
	}

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
