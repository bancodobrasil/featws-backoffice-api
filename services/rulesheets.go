package services

import (
	"context"
	"fmt"
	"unicode"

	"github.com/bancodobrasil/featws-api/dtos"
	"github.com/bancodobrasil/featws-api/models"
	"github.com/bancodobrasil/featws-api/repository"
	"github.com/gosimple/slug"
	log "github.com/sirupsen/logrus"
)

// FindOptions type defines options for limiting and paging search results.
//
// Property
//   - Limit: property is an integer that specifies the maximum number of results to be returned by a search or query. It is often used in conjunction with the `Page` property to implement pagination. For example, if `Limit` is set to 10 and there are 50 results, the
//   - Page: property is an integer that represents the current page number in a paginated result set. It is often used in combination with the `Limit` property to determine which subset of data to return. For example, if `Limit` is set to 10 and `Page` is set
type FindOptions struct {
	Limit int
	Page  int
}

// Rulesheets defines an interface for CRUD operations on rulesheets.
// Property:
//   - Create: Create is a method that creates a new rulesheet in the database. It takes a context and a pointer to a dtos.Rulesheet object as input and returns an error if the operation fails.
//   - Find: method is used to retrieve a list of Rulesheets based on a filter and options. The filter parameter is used to specify the criteria for selecting Rulesheets, while the options parameter is used to specify additional options such as sorting and pagination. The method returns a slice of Rulesheet DTOs and
//   - Count: method is used to count the number of documents that match a given filter in the database collection. It takes a context.Context object and an entity interface{} as input parameters and returns the count of documents as an int64 and an error object. The entity parameter is used to specify the type of
//   - Get: method is used to retrieve a single Rulesheet entity by its unique identifier (id). It takes in a context.Context object and the id of the Rulesheet to be retrieved as parameters, and returns a pointer to the dtos.Rulesheet object and an error object. If the Rulesheet
//   - Update: is a method defined in the Rulesheets interface that takes a context.Context and a dtos.Rulesheet entity as input parameters and returns a pointer to a dtos.Rulesheet and an error. This method is used to update an existing rulesheet entity in the data store.
//   - Delete: method is used to delete a rulesheet from the database. It takes a context.Context and a string id as input parameters and returns a boolean value and an error. The boolean value indicates whether the deletion was successful or not. The error value indicates any error that occurred during the deletion process.
type Rulesheets interface {
	Create(context.Context, *dtos.Rulesheet) error
	Find(ctx context.Context, filter interface{}, options *FindOptions) ([]*dtos.Rulesheet, error)
	Count(ctx context.Context, entity interface{}) (count int64, err error)
	Get(ctx context.Context, id string) (*dtos.Rulesheet, error)
	Update(ctx context.Context, entity dtos.Rulesheet) (*dtos.Rulesheet, error)
	Delete(ctx context.Context, id string) (bool, error)
}

// rulesheets contains a Gitlab service and a repository for rulesheets.
//
// Property:
//   - gitlabService: It seems that `gitlabService` is a variable of type `Gitlab`, which could be a struct or an interface. It is likely used to interact with GitLab API or services related to GitLab. However, without more context or code, it's difficult to determine its exact purpose.
//   - repository: property is of type `repository.Rulesheets`. It is likely a reference to a repository object that contains information about rulesheets, such as their names, contents, and metadata. This object may be used to perform various operations on the rulesheets, such as retrieving, updating,
type rulesheets struct {
	gitlabService Gitlab
	repository    repository.Rulesheets
}

// NewRulesheets creates a new instance of a rulesheets struct with a given repository and Gitlab service.
func NewRulesheets(repository repository.Rulesheets, gitlabService Gitlab) Rulesheets {
	return rulesheets{
		repository:    repository,
		gitlabService: gitlabService,
	}
}

// CreateRulesheet is responsible for creating a new rulesheet. It takes in a `context.Context` object
// and a `*dtos.Rulesheet` object as parameters. It first converts the `*dtos.Rulesheet` object to a
// `models.Rulesheet` object using the `models.NewRulesheetV1` function. It then generates a slug for
// the rulesheet if it doesn't already have one. It creates the rulesheet in the repository using the
// `rs.repository.Create` function and saves it to GitLab using the `rs.gitlabService.Save` function.
// Finally, it fills the `*dtos.Rulesheet` object with GitLab information using the`rs.gitlabService.Fill`
// function. If any errors occur during the process, it logs the error and returns it.
func (rs rulesheets) Create(ctx context.Context, rulesheetDTO *dtos.Rulesheet) (err error) {

	rulesheet, _ := models.NewRulesheetV1(*rulesheetDTO)
	if rulesheet.Slug == "" {
		rulesheet.Slug = slug.Make(rulesheet.Name)
	}

	fmt.Print(rulesheet.Slug)

	err = rs.repository.Create(ctx, &rulesheet)
	if err != nil {
		log.Errorf("Error on create rulesheet into repository: %v", err)
		return
	}
	rulesheetDTO.ID = rulesheet.ID
	rulesheetDTO.Slug = rulesheet.Slug
	err = rs.gitlabService.Save(rulesheetDTO, "[FEATWS BOT] Create Repo")
	if err != nil {
		log.Errorf("Error on save rulesheet into repository: %v", err)
		return
	}

	err = rs.gitlabService.Fill(rulesheetDTO)
	if err != nil {
		log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
		return
	}

	return
}

// Find is responsible for finding rulesheets based on a filter and returning them as an array
// of `dtos.Rulesheet` objects. It takes in a `context.Context` object, a filter interface, and a
// `FindOptions` object as parameters. The `FindOptions` object is used to specify the limit and page
// number for pagination.
func (rs rulesheets) Find(ctx context.Context, filter interface{}, options *FindOptions) (result []*dtos.Rulesheet, err error) {

	var opts *repository.FindOptions = nil

	if options != nil {
		opts = &repository.FindOptions{
			Limit: options.Limit,
			Page:  options.Page,
		}
	}

	entities, err := rs.repository.Find(ctx, filter, opts)
	if err != nil {
		log.Errorf("Error on fetch the rulesheets(find): %v", err)
		return
	}

	result = make([]*dtos.Rulesheet, 0)

	for _, entity := range entities {
		result = append(result, newRulesheetDTO(entity))
	}

	return
}

// Count is a method of the `rulesheets` struct that implements the `Rulesheets` interface. It
// takes a `context.Context` object and an `entity` interface{} as input parameters and returns the
// count of documents as an `int64` and an error object. The `entity` parameter is used to specify the
// type of entity to count in the database collection. The function calls the `Count` method of the
// `repository` property of the `rulesheets` struct, passing in the `ctx` and `entity` parameters. If
// an error occurs during the count operation, the function logs the error and returns it. Otherwise,
// it returns the count of documents.
func (rs rulesheets) Count(ctx context.Context, entity interface{}) (count int64, err error) {

	count, err = rs.repository.Count(ctx, entity)
	if err != nil {
		log.Errorf("Error on count the entities(find): %v", err)
		return
	}

	return
}

// Get function is a method of the `rulesheets` struct that implements the `Rulesheets`
// interface. It takes a `context.Context` object and a `string` id as input parameters and returns a
// pointer to a `dtos.Rulesheet` object and an error object. The function retrieves a single rulesheet
// entity by its unique identifier (id) from the repository using the `rs.repository.Get` function. It
// then converts the `models.Rulesheet` object to a `dtos.Rulesheet` object using the `newRulesheetDTO`
// function. Finally, it fills the `*dtos.Rulesheet` object with GitLab information using the
// `rs.gitlabService.Fill` function. If any errors occur during the process, it logs the error and
// returns it.
func (rs rulesheets) Get(ctx context.Context, id string) (result *dtos.Rulesheet, err error) {

	isSlug := true
	if unicode.IsDigit(rune(id[0])) {
		isSlug = false
	}

	var entity *models.Rulesheet

	if isSlug {
		findResult, err2 := rs.repository.Find(ctx, map[string]interface{}{"slug": id}, nil)
		if err2 != nil {
			log.Errorf("Error on fetch rulesheet(get): %v", err)
			return
		}
		if len(findResult) == 0 {
			return nil, nil
		}
		entity = findResult[0]
	} else {
		entity, err = rs.repository.Get(ctx, id)
		if err != nil {
			log.Errorf("Error on fetch rulesheet(get): %v", err)
			return
		}
	}
	result = newRulesheetDTO(entity)

	if result != nil {
		err = rs.gitlabService.Fill(result)
		if err != nil {
			log.Errorf("Error on fill rulesheet with gitlab information: %v", err)
			return
		}
	}

	return
}

// UpdateRulesheet function is a method of the `rulesheets` struct that implements the `Rulesheets`
// interface. It takes a `context.Context` object and a `dtos.Rulesheet` object as input parameters and
// returns a pointer to a `dtos.Rulesheet` object and an error object.
func (rs rulesheets) Update(ctx context.Context, rulesheetDTO dtos.Rulesheet) (result *dtos.Rulesheet, err error) {

	entity, _ := models.NewRulesheetV1(rulesheetDTO)

	_, err = rs.repository.Update(ctx, entity)
	if err != nil {
		log.Errorf("Error on update the rulesheet from repository: %v", err)
		return
	}

	err = rs.gitlabService.Save(&rulesheetDTO, "[FEATWS BOT] Update Repo")
	if err != nil {
		log.Errorf("Error on save the rulesheet into repository: %v", err)
		return
	}

	result = &rulesheetDTO

	return
}

// Delete function is a method of the `rulesheets` struct that implements the `Rulesheets`
// interface. It takes a `context.Context` object and a `string` id as input parameters and returns a
// boolean value and an error object. The function is responsible for deleting a rulesheet from the db.
func (rs rulesheets) Delete(ctx context.Context, id string) (bool, error) {

	db := rs.repository.GetDB()

	tx := db.Begin()
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		tx.Rollback()
	// 	}
	// }()

	// get the specific rulesheet
	rulesheet, err := rs.repository.Get(ctx, id)
	if err != nil {
		tx.Rollback()
		log.Errorf("Error on fetch rulesheet(get): %v", err)
		return false, err
	}

	// update the ruleshet name to deleted
	rulesheet.Name = fmt.Sprintf("%s-deleted-%v", rulesheet.Name, rulesheet.ID)

	// update the rulesheet
	_, err = rs.repository.UpdateInTransaction(ctx, tx, *rulesheet)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	_, err = rs.repository.DeleteInTransaction(ctx, tx, id)
	if err != nil {
		tx.Rollback()
		log.Errorf("Error on delete the rulesheet from repository: %v", err)
		return false, err
	}

	return true, tx.Commit().Error
}

// The function creates a new DTO for a rulesheet entity
func newRulesheetDTO(entity *models.Rulesheet) *dtos.Rulesheet {
	return &dtos.Rulesheet{
		ID:            entity.ID,
		Name:          entity.Name,
		Description:   entity.Description,
		Slug:          entity.Slug,
		HasStringRule: entity.HasStringRule,
	}
}
