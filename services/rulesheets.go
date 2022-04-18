package services

import (
	"context"

	"github.com/bancodobrasil/featws-api/models"
	"github.com/bancodobrasil/featws-api/repository"
)

// CreateRulesheet ...
func CreateRulesheet(ctx context.Context, rulesheet *models.Rulesheet) (err error) {

	//TODO verifica unicidade do nome

	err = repository.GetRulesheetsRepository().Create(ctx, rulesheet)
	if err != nil {
		return
	}

	err = saveInGitlab(rulesheet, "[FEATWS BOT] Create Repo")
	if err != nil {
		return
	}

	err = fillWithGitlab(rulesheet)
	if err != nil {
		return
	}

	return
}

// FetchRulesheets ...
func FetchRulesheets(ctx context.Context, filter interface{}) (result []*models.Rulesheet, err error) {

	result, err = repository.GetRulesheetsRepository().Find(ctx, filter)
	if err != nil {
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
		return
	}

	if result != nil {
		err = fillWithGitlab(result)
		if err != nil {
			return
		}
	}

	return
}

// UpdateRulesheet ...
func UpdateRulesheet(ctx context.Context, entity models.Rulesheet) (result *models.Rulesheet, err error) {

	result, err = repository.GetRulesheetsRepository().Update(ctx, entity)
	if err != nil {
		return
	}

	err = saveInGitlab(&entity, "[FEATWS BOT] Update Repo")
	if err != nil {
		return
	}

	err = fillWithGitlab(result)
	if err != nil {
		return
	}

	return
}

// DeleteRulesheet ...
func DeleteRulesheet(ctx context.Context, id string) (deleted bool, err error) {

	deleted, err = repository.GetRulesheetsRepository().Delete(ctx, id)
	if err != nil {
		return
	}

	return
}
