package services

import (
	"context"
	"errors"
	"testing"

	mocks_repository "github.com/bancodobrasil/featws-api/mocks/repository"
	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
	"github.com/bancodobrasil/featws-api/models"
	"gorm.io/gorm"
)

func TestGetWithErrorOnFill(t *testing.T) {

	ctx := context.Background()

	expectedEntity := &models.Rulesheet{
		Model: gorm.Model{
			ID: 1,
		},
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
