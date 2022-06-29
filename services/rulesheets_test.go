package services

import (
	"context"
	"errors"
	"testing"

	"github.com/bancodobrasil/featws-api/dtos"
	mocks_repository "github.com/bancodobrasil/featws-api/mocks/repository"
	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
	"github.com/bancodobrasil/featws-api/models"
)

func TestGetWithErrorOnFill(t *testing.T) {

	ctx := context.Background()

	dto := &dtos.Rulesheet{
		ID: 1,
	}

	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}

	repository := new(mocks_repository.Rulesheets)
	repository.On("Get", ctx, "1").Return(&entity, nil)

	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Fill", dto).Return(errors.New("error on fill"))

	service := NewRulesheets(repository, gitlabService)

	_, err = service.Get(ctx, "1")

	if err == nil || err.Error() != "error on fill" {
		t.Error("expected error on fill")
	}

}
