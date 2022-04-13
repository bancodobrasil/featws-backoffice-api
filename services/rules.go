package services

import (
	"context"

	"github.com/bancodobrasil/featws-api/models"
	"github.com/bancodobrasil/featws-api/repository"
)

// CreateRule ...
func CreateRule(ctx context.Context, rule *models.Rule) (err error) {

	//TODO verifica unicidade do nome

	err = repository.GetRulesRepository().Create(ctx, rule)
	if err != nil {
		return
	}

	return
}

// FetchRules ...
func FetchRules(ctx context.Context, filter interface{}) (result []models.Rule, err error) {

	result, err = repository.GetRulesRepository().Find(ctx, filter)
	if err != nil {
		return
	}

	return
}

// FetchRule ...
func FetchRule(ctx context.Context, id string) (result *models.Rule, err error) {

	result, err = repository.GetRulesRepository().Get(ctx, id)
	if err != nil {
		return
	}

	return
}

// UpdateRule ...
func UpdateRule(ctx context.Context, entity models.Rule) (result *models.Rule, err error) {

	result, err = repository.GetRulesRepository().Update(ctx, entity)
	if err != nil {
		return
	}

	return
}

// DeleteRule ...
func DeleteRule(ctx context.Context, id string) (deleted bool, err error) {

	deleted, err = repository.GetRulesRepository().Delete(ctx, id)
	if err != nil {
		return
	}

	return
}
