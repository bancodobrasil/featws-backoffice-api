package services

import (
	"context"
	"errors"
	"testing"

	"github.com/bancodobrasil/featws-api/dtos"
	mocks_repository "github.com/bancodobrasil/featws-api/mocks/repository"
	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
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

	expectedEntity := &models.Rulesheet{
		Model: gorm.Model{
			ID: 2,
		},
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

// func TestCreateRulesheet(t *testing.T) {

// 	ctx :=  context.Background()

// 	expectedEntity := &models.Rulesheet{
// 		Model: gorm.Model{
// 			ID: 1,
// 		}
// 	}

// 	repository := new(mocks_repository.Rulesheets)
// 	repository.

// }
