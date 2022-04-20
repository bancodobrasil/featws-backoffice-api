package services

import (
	"context"

	"github.com/bancodobrasil/featws-api/models"
	"github.com/bancodobrasil/featws-api/repository"
	log "github.com/sirupsen/logrus"
)

// CreateRulesheet ...
func CreateRulesheet(ctx context.Context, rulesheet *models.Rulesheet) (err error) {

	//TODO verifica unicidade do nome

	err = repository.GetRulesheetsRepository().Create(ctx, rulesheet)
	if err != nil {
		log.Errorf("Error on create rulesheet into repository: %v", err)
		return
	}

	err = saveInGitlab(rulesheet, "[FEATWS BOT] Create Repo")
	if err != nil {
		log.Errorf("Error on save rulesheet into repository: %v", err)
		return
	}

	err = fillWithGitlab(rulesheet)
	if err != nil {
		log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
		return
	}

	return
}

// FetchRulesheets ...
func FetchRulesheets(ctx context.Context, filter interface{}) (result []*models.Rulesheet, err error) {

	result, err = repository.GetRulesheetsRepository().Find(ctx, filter)
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
func FetchRulesheet(ctx context.Context, id string) (result *models.Rulesheet, err error) {

	result, err = repository.GetRulesheetsRepository().Get(ctx, id)
	if err != nil {
		log.Errorf("Error on fetch rulesheet(get): %v", err)
		return
	}

	if result != nil {
		err = fillWithGitlab(result)
		if err != nil {
			log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
			return
		}
	}

	return
}

// UpdateRulesheet ...
func UpdateRulesheet(ctx context.Context, entity models.Rulesheet) (result *models.Rulesheet, err error) {

	result, err = repository.GetRulesheetsRepository().Update(ctx, entity)
	if err != nil {
		log.Errorf("Error on update the rulesheet from repository: %v", err)
		return
	}

	err = saveInGitlab(&entity, "[FEATWS BOT] Update Repo")
	if err != nil {
		log.Errorf("Error on save the rulesheet into repository: %v", err)
		return
	}

	err = fillWithGitlab(result)
	if err != nil {
		log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
		return
	}

	return
}

// DeleteRulesheet ...
func DeleteRulesheet(ctx context.Context, id string) (deleted bool, err error) {

	deleted, err = repository.GetRulesheetsRepository().Delete(ctx, id)
	if err != nil {
		log.Errorf("Error on delete the rulesheet from repository: %v", err)
		return
	}

	return
}
