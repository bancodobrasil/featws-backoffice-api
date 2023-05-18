package repository

import (
	"context"
	"errors"

	"github.com/bancodobrasil/featws-api/utils"
	telemetry "github.com/bancodobrasil/gin-telemetry"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"

	"gorm.io/gorm"
)

// FindOptions type defines options for limiting and paging search results.
//
// Property:
//
//   - Limit: an integer that specifies the maximum number of results to be returned by a search or query. It's often used in conjunction with the Page property to implement pagination.
//   - Page:  an integer that represents the current page number in a paginated result set. It's often used in combination with the Limit property to determine which subset of data to return.
type FindOptions struct {
	Limit int
	Page  int
}

// Repository defines a set of methods for interacting with a database using GORM in Go.
//
// Property:
//
//   - GetDB: returns a pointer to a gorm.DB instance that allows performing database operations.
//   - Create: creates a new entity of type T in the database.
//   - CreateInTransaction: creates a new entity in the database within a transaction. It takes a context.Context and a *gorm.DB as parameters, along with a pointer to the entity to be created. It returns an error if the creation fails.
//   - Find: retrieves a list of entities from the repository based on specified criteria in the FindOptions parameter. It returns a slice of pointers to the type T and an error if any occurred during the operation.
//   - FindInTransaction: finds multiple entities in a transactional context. It takes a context.Context object, a *gorm.DB object representing the transaction, an interface object representing the entity to be found, and a *FindOptions object representing the options for the find operation. It returns a slice of pointers to the found entities.
//   - Count: returns the number of entities in the repository.
//   - CountInTransaction: : counts the number of entities in the database within a transaction. It takes a context.Context and a *gorm.DB as parameters, along with an entity interface representing the type of entity to count. It returns the count as an int64 and an error if any occurred.
//   - Get: retrieves a single entity of type T from the repository based on the provided ID. It returns the retrieved entity and an error if any occurred during the retrieval process.
//   - GetInTransaction: retrieves a single entity of type T from the database within a transaction. It takes a context.Context and a *gorm.DB as parameters and returns a pointer to the retrieved entity of type T and an error if any.
//   - Update: updates an existing entity in the repository. It takes a context.Context object and an entity of type T as input and returns the updated entity of type T and an error. If the update is successful, the updated entity is returned; otherwise, an error is returned.
//   - UpdateInTransaction: updates an entity of type T in a transactional context. It takes a context.Context and a *gorm.DB as parameters, along with the entity to be updated. It returns the updated entity of type *T and an error if any occurred during the update process.
//   - Delete: deletes an entity from the repository based on its ID. It takes a context.Context object and the ID of the entity to be deleted as input parameters. It returns a boolean value indicating whether the entity was successfully deleted or not, along with an error object if any.
//   - DeleteInTransaction:deletes an entity with the given ID from the database within a transaction. It takes a context.Context object, a *gorm.DB object representing the transaction, and the ID of the entity to be deleted. It returns a boolean indicating whether the entity was successfully deleted and an error.
type Repository[T any] interface {
	GetDB() *gorm.DB
	Create(ctx context.Context, entity *T) error
	CreateInTransaction(ctx context.Context, db *gorm.DB, entity *T) error
	Find(ctx context.Context, entity interface{}, options *FindOptions) (list []*T, err error)
	FindInTransaction(ctx context.Context, db *gorm.DB, entity interface{}, options *FindOptions) (list []*T, err error)
	Count(ctx context.Context, entity interface{}) (count int64, err error)
	CountInTransaction(ctx context.Context, db *gorm.DB, entity interface{}) (count int64, err error)
	Get(ctx context.Context, id string) (entity *T, err error)
	GetInTransaction(ctx context.Context, db *gorm.DB, id string) (entity *T, err error)
	Update(ctx context.Context, entity T) (updated *T, err error)
	UpdateInTransaction(ctx context.Context, db *gorm.DB, entity T) (updated *T, err error)
	Delete(ctx context.Context, id string) (deleted bool, err error)
	DeleteInTransaction(ctx context.Context, db *gorm.DB, id string) (deleted bool, err error)
}

// Repository creates a generic repository struct that contains a pointer to a GORM database.
//
// Property:
//   - db: it's a pointer to a gorm.DB object, which is a database ORM library for Go. It is used to interact with a database and perform CRUD operations on the data.
type repository[T any] struct {
	db *gorm.DB
}

const (
	create = "repo-create"
	find   = "repo-find"
	count  = "repo-count"
	get    = "repo-get"
	update = "repo-update"
	delete = "repo-delete"
)

// Create ...
func (r *repository[T]) Create(ctx context.Context, entity *T) error {
	db := r.newSession(ctx)
	return r.CreateInTransaction(ctx, db, entity)
}

// CreateInTransaction ...
func (r *repository[T]) CreateInTransaction(ctx context.Context, db *gorm.DB, entity *T) error {
	// add the span of database query on the root span of the context
	span := utils.GenerateSpanTracer(ctx, create)
	defer span()

	result := db.Create(&entity)
	if result.Error != nil {
		log.WithContext(ctx).Errorf("error on insert the result into model: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected != 1 {
		err := errors.New("error on insert not inserted")
		log.WithContext(ctx).Error(err.Error())
		return err
	}

	return nil
}

// Find ...
func (r *repository[T]) Find(ctx context.Context, entity interface{}, options *FindOptions) (list []*T, err error) {
	db := r.newSession(ctx)
	return r.FindInTransaction(ctx, db, entity, options)
}

// FindInTransaction ...
func (r *repository[T]) FindInTransaction(ctx context.Context, db *gorm.DB, entity interface{}, options *FindOptions) (list []*T, err error) {
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-find", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	if options != nil {
		limit := 10
		if options.Limit != 0 {
			limit = options.Limit
		}
		db = db.Limit(limit)
		if options.Page != 0 {
			db = db.Offset((options.Page - 1) * limit)
		}
	}

	result := db.Find(&list, entity)

	err = result.Error
	if err != nil {
		log.WithContext(ctx).Errorf("Error on find: %v", err)
		return
	}

	return
}

// Count ...
func (r *repository[T]) Count(ctx context.Context, entity interface{}) (count int64, err error) {
	db := r.newSession(ctx)
	return r.CountInTransaction(ctx, db, entity)
}

// CountInTransaction ...
func (r *repository[T]) CountInTransaction(ctx context.Context, db *gorm.DB, entity interface{}) (count int64, err error) {
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-count", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	count = 0

	result := db.Where(entity).Count(&count)

	err = result.Error
	if err != nil {
		log.WithContext(ctx).Errorf("Error on find: %v", err)
		return
	}

	return
}

// Get ...
func (r *repository[T]) Get(ctx context.Context, id string) (entity *T, err error) {
	db := r.newSession(ctx)
	return r.GetInTransaction(ctx, db, id)
}

// Get ...
func (r *repository[T]) GetInTransaction(ctx context.Context, db *gorm.DB, id string) (entity *T, err error) {
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-get", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	result := db.First(&entity, id)

	err = result.Error
	if err != nil {
		log.WithContext(ctx).Errorf("Error on find one result into collection: %v", err)
		return
	}

	return
}

// Update ...
func (r *repository[T]) Update(ctx context.Context, entity T) (updated *T, err error) {
	db := r.newSession(ctx)
	return r.UpdateInTransaction(ctx, db, entity)
}

// UpdateInTransaction ...
func (r *repository[T]) UpdateInTransaction(ctx context.Context, db *gorm.DB, entity T) (updated *T, err error) {
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-update", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	result := db.Model(entity).Save(&entity)

	err = result.Error
	if err != nil {
		log.WithContext(ctx).Errorf("Error on update into collection: %v", err)
		return
	}

	updated = &entity

	return
}

// Delete ...
func (r *repository[T]) Delete(ctx context.Context, id string) (deleted bool, err error) {
	db := r.newSession(ctx)
	return r.DeleteInTransaction(ctx, db, id)
}

// DeleteInTransaction ...
func (r *repository[T]) DeleteInTransaction(ctx context.Context, db *gorm.DB, id string) (deleted bool, err error) {
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-delete", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	entity, err := r.Get(ctx, id)

	if err != nil {
		log.WithContext(ctx).Errorf("Error on get before delete: %v", err)
		return
	}

	if entity == nil {
		deleted = true
		return
	}

	result := db.Model(entity).Delete(entity)

	err = result.Error
	if err != nil {
		log.Errorf("Error on delete from collection: %v", err)
		return
	}

	deleted = result.RowsAffected == 1

	return
}

func (r *repository[T]) newSession(ctx context.Context) *gorm.DB {
	return r.GetDB().Session(&gorm.Session{}).Model(new(T)).WithContext(ctx)
}

func (r *repository[T]) GetDB() *gorm.DB {
	return r.db
}
