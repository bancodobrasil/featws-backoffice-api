package services

import (
	"context"
	"errors"
	"testing"

	"github.com/bancodobrasil/featws-api/dtos"
	mocks_repository "github.com/bancodobrasil/featws-api/mocks/repository"
	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
	"github.com/bancodobrasil/featws-api/models"
	log "github.com/sirupsen/logrus"
)

func TestGetWithErrorOnFill(t *testing.T) {

	ctx := context.Background()

	expectedEntity := &dtos.Rulesheet{
		ID: 1,
	}

	repository := new(mocks_repository.Rulesheets)
	repository.On("Get", ctx, "1").Return(expectedEntity, nil)

	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Fill", expectedEntity).Return(errors.New("error on fill"))

	service := NewRulesheets(repository, gitlabService)

	_, err := service.Get(ctx, "1")

	if err == nil || err.Error() != "error on fill" {
		t.Error("expected error on fill")
	}

}

func TestGet(t *testing.T) {

	ctx := context.Background()

	expectedEntity := &dtos.Rulesheet{
		ID:   2,
		Name: "test",
	}

	repository := new(mocks_repository.Rulesheets)
	repository.On("Get", ctx, "2").Return(expectedEntity, nil)

	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Fill", expectedEntity).Return(expectedEntity)

	service := NewRulesheets(repository, gitlabService)

	result, _ := service.Get(ctx, "2")

	if result.Name != expectedEntity.Name {
		t.Error("Error on get the rulesheet")
	}
}

func TestCreateRulesheet(t *testing.T) {

	ctx := context.Background()

	dto := &dtos.Rulesheet{
		ID: 1,
	}

	rulesheet, err := models.NewRulesheetV1(*dto)
	if err != nil {
		log.Errorf("Error on create rulesheet on create model: %v", err)
		return
	}

	fakeCommitMessage := "test"

	repository := new(mocks_repository.Rulesheets)
	repository.On("Create", ctx, &rulesheet).Return(rulesheet, nil)

	gitlabService := new(mocks_services.Gitlab)
	// gitlabService.On("Fill", dto).Return(nil)
	gitlabService.On("Save", dto, fakeCommitMessage).Return(nil)

	service := NewRulesheets(repository, gitlabService)

	result := service.Create(ctx, dto)

	// assert.Nil(t, result)
	// assert.NotNil(t, result)

	if result != nil {
		t.Error("Error on create the rulesheet")
	}

}

func TestFind(t *testing.T) {

	ctx := context.Background()

	dto := &dtos.Rulesheet{
		ID:   2,
		Name: "test",
	}

	rulesheet, err := models.NewRulesheetV1(*dto)
	if err != nil {
		log.Errorf("Error on create rulesheet on create model: %v", err)
		return
	}

	repository := new(mocks_repository.Rulesheets)

	repository.On("Find", ctx, "2").Return(rulesheet, nil)

	gitlabService := new(mocks_services.Gitlab)

	service := NewRulesheets(repository, gitlabService)

	results, _ := service.Find(ctx, "2")

	for _, result := range results {
		if result.Name != rulesheet.Name {
			t.Error("Error on find the rulesheet")
		}
	}

}
