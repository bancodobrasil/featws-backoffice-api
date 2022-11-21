package services_test

import (
	"context"
	"errors"
	"log"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bancodobrasil/featws-api/dtos"
	mocks_repository "github.com/bancodobrasil/featws-api/mocks/repository"
	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
	"github.com/bancodobrasil/featws-api/models"
	"github.com/bancodobrasil/featws-api/repository"
	"github.com/bancodobrasil/featws-api/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	services := services.NewRulesheets(repository, gitlabService)
	_, err = services.Get(ctx, "1")

	if err == nil || err.Error() != "error on fill" {
		t.Error("expected error on fill")
	}
}

func TestGetSucess(t *testing.T) {

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
	gitlabService.On("Fill", dto).Return(nil)
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Get(ctx, "1")
	if err != nil {
		t.Error("unexpected error on get")
	}
}

func TestGetWithErrorOnCreateModel(t *testing.T) {

	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}

	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Get", ctx, "1").Return(&entity, errors.New("error on model creation"))
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Fill", dto).Return(nil)
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Get(ctx, "1")
	if err == nil || err.Error() != "error on model creation" {
		t.Error("unexpected error on get")
	}
}

func TestCreateWithErrorOnFill(t *testing.T) {

	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Create", ctx, &entity).Return(nil)
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(nil)
	gitlabService.On("Fill", dto).Return(errors.New("error on fill"))
	service := services.NewRulesheets(repository, gitlabService)
	err = service.Create(ctx, dto)
	if err == nil || err.Error() != "error on fill" {
		t.Error("expected error on fill")
	}
}

func TestCreateSuccess(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Create", ctx, &entity).Return(nil)
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(nil)
	gitlabService.On("Fill", dto).Return(nil)
	service := services.NewRulesheets(repository, gitlabService)
	err = service.Create(ctx, dto)
	if err != nil {
		t.Error("unexpected error on create")
	}

}

func TestCreateWithError(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}

	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Create", ctx, &entity).Return(errors.New("error on create"))
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(nil)
	gitlabService.On("Fill", dto).Return(nil)
	service := services.NewRulesheets(repository, gitlabService)
	err = service.Create(ctx, dto)
	if err == nil || err.Error() != "error on create" {
		t.Error("expected error on create")
	}
}

func TestCreateWithErroOnSave(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}

	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Create", ctx, &entity).Return(nil)
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Save", dto, "[FEATWS BOT] Create Repo").Return(errors.New("error on save"))
	gitlabService.On("Fill", dto).Return(nil)
	service := services.NewRulesheets(repository, gitlabService)
	err = service.Create(ctx, dto)
	if err == nil || err.Error() != "error on save" {
		t.Error("expected error on save")
	}
}

func TestUpdateWithErrorOnSave(t *testing.T) {

	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Update", ctx, entity).Return(nil, nil)
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Save", dto, "[FEATWS BOT] Update Repo").Return(errors.New("error on save"))
	// gitlabService.On("Fill", dto).Return(errors.New("error on fill"))
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Update(ctx, *dto)
	if err == nil || err.Error() != "error on save" {
		t.Error("expected error on save")
	}
}

func TestFindSuccess(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repo := new(mocks_repository.Rulesheets)
	repoFindOptions := repository.FindOptions{}
	entities := []*models.Rulesheet{&entity}
	repo.On("Find", ctx, dto.ID, &repoFindOptions).Return(entities, nil)
	service := services.NewRulesheets(repo, nil)
	serviceFindOptions := services.FindOptions{0, 0}
	_, err = service.Find(ctx, dto.ID, &serviceFindOptions)
	if err != nil {
		t.Error("unexpected error on find")
	}
}

func TestFindWithError(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repo := new(mocks_repository.Rulesheets)
	repoFindOptions := repository.FindOptions{}
	entities := []*models.Rulesheet{&entity}
	repo.On("Find", ctx, dto.ID, &repoFindOptions).Return(entities, errors.New("error on find"))
	service := services.NewRulesheets(repo, nil)
	serviceFindOptions := services.FindOptions{0, 0}
	_, err = service.Find(ctx, dto.ID, &serviceFindOptions)
	if err != nil && err.Error() != "error on find" {
		t.Error("unexpected error on find")
	}

}

func TestCountSuccess(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	_, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Count", ctx, nil).Return(int64(1), nil)
	service := services.NewRulesheets(repository, nil)
	_, err = service.Count(ctx, nil)
	if err != nil {
		t.Error("unexpected error on count")
	}
}

func TestCountWithError(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	_, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Count", ctx, nil).Return(int64(0), errors.New("error on count"))
	service := services.NewRulesheets(repository, nil)
	_, err = service.Count(ctx, nil)
	if err == nil || err.Error() != "error on count" {
		t.Error("expected error on count")
	}
}

func TestUpdateSuccess(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Update", ctx, entity).Return(nil, nil)
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Save", dto, "[FEATWS BOT] Update Repo").Return(nil)
	gitlabService.On("Fill", dto).Return(nil)
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Update(ctx, *dto)
	if err != nil {
		t.Error("unexpected error on update")
	}
}

func TestUpdateWithError(t *testing.T) {
	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID: 1,
	}
	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	repository := new(mocks_repository.Rulesheets)
	repository.On("Update", ctx, entity).Return(nil, errors.New("error on update"))
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Save", dto, "[FEATWS BOT] Update Repo").Return(nil)
	gitlabService.On("Fill", dto).Return(nil)
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Update(ctx, *dto)
	if err == nil || err.Error() != "error on update" {
		t.Error("expected error on update")
	}
}

// func (s *RepositorySuite) SetupSuite() {
// 	var (
// 		err error
// 	)
//

// 	// Dialector for mariadb
// 	dialector := mysql.New(mysql.Config{
// 		DSN:        "sqlmock_db_0", // DSN data source name
// 		DriverName: "mysql",
// 		Conn:       s.conn,
// 	})

// 	s.DB, err = gorm.Open(dialector, &gorm.Config{})
// 	assert.NoError(s.T(), err)
// }

// func (s *RepositorySuite) TestDeleteRulesheet() {
// 	ctx := context.Background()
// 	dto := &dtos.Rulesheet{
// 		ID: 1,
// 		Name: "test",
// 	}
// 	entity, err := models.NewRulesheetV1(*dto)

// 	s.mock.ExpectBegin()

// }

func TestDeleteSuccess(t *testing.T) {
	// Init fake db connection
	conn, mocks, err := sqlmock.New()
	assert.NoError(t, err)

	mocks.ExpectBegin()

	mocks.ExpectCommit()

	dialector := mysql.New(mysql.Config{
		DriverName:                "mysql",
		Conn:                      conn,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID:   1,
		Name: "test",
	}

	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	newID := strconv.Itoa(int(dto.ID))
	repository := new(mocks_repository.Rulesheets)
	repository.On("GetDB").Return(db)
	repository.On("Get", ctx, newID).Return(&entity, nil)
	entity.Name = "test-deleted-1"
	repository.On("UpdateInTransaction", ctx, mock.Anything, mock.Anything).Return(&entity, nil)
	repository.On("DeleteInTransaction", ctx, mock.Anything, "1").Return(true, nil)
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Delete", dto).Return(true, nil)
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Delete(ctx, "1")
	if err != nil {
		t.Error("unexpected error on delete")
	}
}

func TestDeleteWithError(t *testing.T) {
	// Init fake db connection
	conn, mocks, err := sqlmock.New()
	assert.NoError(t, err)

	mocks.ExpectBegin()

	mocks.ExpectRollback()

	dialector := mysql.New(mysql.Config{
		DriverName:                "mysql",
		Conn:                      conn,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID:   1,
		Name: "test",
	}

	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	newID := strconv.Itoa(int(dto.ID))
	repository := new(mocks_repository.Rulesheets)
	repository.On("GetDB").Return(db)
	repository.On("Get", ctx, newID).Return(&entity, nil)
	entity.Name = "test-deleted-1"
	repository.On("UpdateInTransaction", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("error on update"))
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Delete", dto).Return(true, nil)
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Delete(ctx, "1")
	if err == nil || err.Error() != "error on update" {
		log.Println(err)
		t.Error("expected error on update")
	}
}

func TestDeleteWithRollBackOnGet(t *testing.T) {
	// Init fake db connection
	conn, mocks, err := sqlmock.New()
	assert.NoError(t, err)

	mocks.ExpectBegin()

	mocks.ExpectRollback()

	dialector := mysql.New(mysql.Config{
		DriverName:                "mysql",
		Conn:                      conn,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID:   1,
		Name: "test",
	}

	_, err = models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	newID := strconv.Itoa(int(dto.ID))
	repository := new(mocks_repository.Rulesheets)
	repository.On("GetDB").Return(db)
	repository.On("Get", ctx, newID).Return(nil, errors.New("error on get"))
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Delete", dto).Return(true, nil)
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Delete(ctx, "1")
	if err == nil || err.Error() != "error on get" {
		t.Error("expected error on get")
	}
}

func TestDeleteWithRollBackOnDeleteInTransaction(t *testing.T) {
	// Init fake db connection
	conn, mocks, err := sqlmock.New()
	assert.NoError(t, err)

	mocks.ExpectBegin()

	mocks.ExpectRollback()

	dialector := mysql.New(mysql.Config{
		DriverName:                "mysql",
		Conn:                      conn,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	ctx := context.Background()
	dto := &dtos.Rulesheet{
		ID:   1,
		Name: "test",
	}

	entity, err := models.NewRulesheetV1(*dto)
	if err != nil {
		t.Error("unexpected error on model creation")
	}
	newID := strconv.Itoa(int(dto.ID))
	repository := new(mocks_repository.Rulesheets)
	repository.On("GetDB").Return(db)
	repository.On("Get", ctx, newID).Return(&entity, nil)
	entity.Name = "test-deleted-1"
	repository.On("UpdateInTransaction", ctx, mock.Anything, mock.Anything).Return(&entity, nil)
	repository.On("DeleteInTransaction", ctx, mock.Anything, "1").Return(false, errors.New("error on delete"))
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Delete", dto).Return(true, nil)
	service := services.NewRulesheets(repository, gitlabService)
	_, err = service.Delete(ctx, "1")
	if err == nil || err.Error() != "error on delete" {
		t.Error("expected error on delete")
	}
}
