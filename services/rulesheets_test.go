package services

import (
	"context"
	"errors"
	"testing"

	"github.com/bancodobrasil/featws-api/dtos"
	mocks_repository "github.com/bancodobrasil/featws-api/mocks/repository"
	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
	"github.com/bancodobrasil/featws-api/models"
	// "github.com/bancodobrasil/featws-api/repository"
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

	services := NewRulesheets(repository, gitlabService)
	_, err = services.Get(ctx, "1")

	if err == nil || err.Error() != "error on fill" {
		t.Error("expected error on fill")
	}
}

// func TestGetSucess(t *testing.T) {

// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}
// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Get", ctx, "1").Return(&entity, nil)
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Fill", dto).Return(nil)
// 	service := NewRulesheets(repository, gitlabService)
// 	_, err = service.Get(ctx, "1")
// 	if err != nil {
// 		t.Error("unexpected error on get")
// 	}
// }

// func TestGetWithErrorOnCreateModel(t *testing.T) {

// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}

// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Get", ctx, "1").Return(&entity, errors.New("error on model creation"))
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Fill", dto).Return(nil)
// 	service := NewRulesheets(repository, gitlabService)
// 	_, err = service.Get(ctx, "1")
// 	if err == nil || err.Error() != "error on model creation" {
// 		t.Error("unexpected error on get")
// 	}
// }

// func TestCreateWithErrorOnFill(t *testing.T) {

// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}
// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Create", ctx, &entity).Return(nil)
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(nil)
// 	gitlabService.On("Fill", dto).Return(errors.New("error on fill"))
// 	service := NewRulesheets(repository, gitlabService)
// 	err = service.Create(ctx, dto)
// 	if err == nil || err.Error() != "error on fill" {
// 		t.Error("expected error on fill")
// 	}
// }

// func TestCreateSuccess(t *testing.T) {
// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}
// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Create", ctx, &entity).Return(nil)
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(nil)
// 	gitlabService.On("Fill", dto).Return(nil)
// 	service := NewRulesheets(repository, gitlabService)
// 	err = service.Create(ctx, dto)
// 	if err != nil {
// 		t.Error("unexpected error on create")
// 	}

// }

// // func TestCreateWithErrorOnModelCreation(t *testing.T) {
// // 	ctx := context.Background()
// // 	dto := &dtos.Rulesheet{
// // 		ID: 1,
// // 	}

// // 	entity, err1 := models.NewRulesheetV1(*dto)
// // 	err1 = errors.New("error on model creation")
// // 	repository := new(mocks_repository.Rulesheets)
// // 	repository.On("Create", ctx, &entity).Return(nil)
// // 	gitlabService := new(mocks_services.Gitlab)
// // 	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(nil)
// // 	gitlabService.On("Fill", dto).Return(nil)
// // 	service := NewRulesheets(repository, gitlabService)
// // 	_ = service.Create(ctx, dto)
// // 	if err1.Error() != "error on model creation" {
// // 		t.Error("expected error on model creation")
// // 	}
// // }

// func TestCreateWithError(t *testing.T) {
// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}

// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Create", ctx, &entity).Return(errors.New("error on create"))
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(nil)
// 	gitlabService.On("Fill", dto).Return(nil)
// 	service := NewRulesheets(repository, gitlabService)
// 	err = service.Create(ctx, dto)
// 	if err == nil || err.Error() != "error on create" {
// 		t.Error("expected error on create")
// 	}
// }

// func TestCreateWithErroOnSave(t *testing.T) {
// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}

// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Create", ctx, &entity).Return(nil)
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(errors.New("error on save"))
// 	gitlabService.On("Fill", dto).Return(nil)
// 	service := NewRulesheets(repository, gitlabService)
// 	err = service.Create(ctx, dto)
// 	if err == nil || err.Error() != "error on save" {
// 		t.Error("expected error on save")
// 	}
// }

// func TestUpdateWithErrorOnSave(t *testing.T) {

// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}
// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Update", ctx, entity).Return(nil, nil)
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Save", dto, "[FEATWS BOT] Update Repo").Return(errors.New("error on save"))
// 	// gitlabService.On("Fill", dto).Return(errors.New("error on fill"))
// 	service := NewRulesheets(repository, gitlabService)
// 	_, err = service.Update(ctx, *dto)
// 	if err == nil || err.Error() != "error on save" {
// 		t.Error("expected error on save")
// 	}
// }

// func TestFindSucess(t *testing.T) {
// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}
// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repositoryMocado := new(mocks_repository.Rulesheets)
// 	repositoryMocado.On("Find", ctx, dto.ID, &repository.FindOptions{}).Return(entity, nil)
// 	service := NewRulesheets(repositoryMocado, nil)
// 	_, err = service.Find(ctx, dto.ID, &FindOptions{0, 0})
// 	if err != nil {
// 		t.Error("unexpected error on find")
// 	}
// }

// func TestUpdateSuccess(t *testing.T) {
// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}
// 	entity, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Update", ctx, entity).Return(nil, nil)
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Save", dto, "[FEATWS BOT] Update Repo").Return(nil)
// 	gitlabService.On("Fill", dto).Return(nil)
// 	service := NewRulesheets(repository, gitlabService)
// 	_, err = service.Update(ctx, *dto)
// 	if err != nil {
// 		t.Error("unexpected error on update")
// 	}
// }

// func TestDeleteSuccess(t *testing.T) {
// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 	}
// 	_, err := models.NewRulesheetV1(*dto)
// 	if err != nil {
// 		t.Error("unexpected error on model creation")
// 	}
// 	repository := new(mocks_repository.Rulesheets)
// 	repository.On("Delete", ctx, "1").Return(true, nil)
// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Delete", dto).Return(true, nil)
// 	service := NewRulesheets(repository, gitlabService)
// 	_, err = service.Delete(ctx, "1")
// 	if err != nil {
// 		t.Error("unexpected error on delete")
// 	}
// }
