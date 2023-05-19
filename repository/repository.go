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

// const is a block defining a set of constants that represent the different types of database
// operations that can be performed by the repository. These constants are used as labels for tracing
// and logging purposes, allowing for easier identification of the specific operation being performed.
// For example, when a `Create` operation is performed, the `create` constant is used to label the
// corresponding log or trace entry.
const (
	create = "repo-create"
	find   = "repo-find"
	count  = "repo-count"
	get    = "repo-get"
	update = "repo-update"
	delete = "repo-delete"
)

// Create it's a function thats creating a new entity of type T in the database. It first creates a new session
// using the `newSession` function, which returns a new `gorm.DB` instance with a new context. Then, it
// calls the `CreateInTransaction` function with the new context, the new session, and a pointer to the
// entity to be created. Finally, it returns the result of the `CreateInTransaction` function, which is
// an error if the creation fails.
func (r *repository[T]) Create(ctx context.Context, entity *T) error {
	db := r.newSession(ctx)
	return r.CreateInTransaction(ctx, db, entity)
}

// CreateInTransaction is a method of the repository struct that creates a new entity of type T in the
// database within a transaction. It takes a context.Context object, a *gorm.DB object representing the
// transaction, and a pointer to the entity to be created as input parameters.
func (r *repository[T]) CreateInTransaction(ctx context.Context, db *gorm.DB, entity *T) error {
	// add the span of database query on the root span of the context
	span := utils.GenerateSpanTracer(ctx, create)
	defer span()

	// This code block creates a new entity of type T in the database using the Create method of the gorm.DB object.
	// The Create method returns a *gorm.DB that provides information about the operation's result, including any errors encountered.
	result := db.Create(&entity)
	if result.Error != nil {
		log.WithContext(ctx).Errorf("error on insert the result into model: %v", result.Error)
		return result.Error
	}

	//The code verifies the success of an entity's database insertion by checking the number of affected rows.
	if result.RowsAffected != 1 {
		err := errors.New("error on insert not inserted")
		log.WithContext(ctx).Error(err.Error())
		return err
	}

	return nil
}

// Find is a method of the `repository` struct that finds a list of entities from the
// repository based on specified criteria in the `FindOptions` parameter. It takes a `context.Context`
// object, an `interface{}` object representing the entity to be found, and a pointer to a
// `FindOptions` object representing the options for the find operation as input parameters. It returns
// a slice of pointers to the type `T` and an error if any occurred during the operation.
func (r *repository[T]) Find(ctx context.Context, entity interface{}, options *FindOptions) (list []*T, err error) {
	db := r.newSession(ctx)
	return r.FindInTransaction(ctx, db, entity, options)
}

// FindInTransaction is a method of a generic repository that finds a list of entities in a database
// transaction using GORM. It takes in a context, a GORM database instance, an entity interface, and
// optional find options. It starts a new span on the context for tracing purposes, applies any
// specified find options to the database query, executes the query using GORM's Find method, and
// returns the resulting list of entities or an error if there was one.
func (r *repository[T]) FindInTransaction(ctx context.Context, db *gorm.DB, entity interface{}, options *FindOptions) (list []*T, err error) {
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-find", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	// The code snippet above checks if the options parameter is not nil. If it's not nil, it assigns a
	// default limit of 10 and verifies if the Limit field of options is not 0. If Limit is not 0, it assigns
	// the value of options.Limit to the limit variable. Finally, it applies the limit to the database query using db.Limit(limit).
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

	// The code below queries the database using the Find method, searching for parameters corresponding
	// to entity, and the results are returned in list.
	result := db.Find(&list, entity)

	err = result.Error
	if err != nil {
		log.WithContext(ctx).Errorf("Error on find: %v", err)
		return
	}

	return
}

// Count is a method that counts the number of entities in a db table. It takes a context and an entity as input parameters and returns the
// count of entities and an error if any. It creates a new database session and calls the
// CountInTransaction method passing the context, db session, and entity as parameters to count
// the number of entities in the table.
func (r *repository[T]) Count(ctx context.Context, entity interface{}) (count int64, err error) {
	db := r.newSession(ctx)
	return r.CountInTransaction(ctx, db, entity)
}

// CountInTransaction is a method that counts the number of records in a db table that match a given
// entity. It takes a context, a db connection, and an entity as input parameters. It uses the OpenTelemetry
// library to create a span for the db query and logs any errors that occur during the query. The method
// returns the count of matching records and any errors encountered during the query.
func (r *repository[T]) CountInTransaction(ctx context.Context, db *gorm.DB, entity interface{}) (count int64, err error) {
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-count", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	count = 0

	// The code is used to query a database using the Where method with a specific entity as a filter. It's
	// used to determine the number of records that match the provided entity, and the result is stored in the count variable.
	result := db.Where(entity).Count(&count)

	err = result.Error
	if err != nil {
		log.WithContext(ctx).Errorf("Error on find: %v", err)
		return
	}

	return
}

// Get method is a function that takes a context and an ID as input parameters and returns a pointer to an
// entity of type T and an error.
func (r *repository[T]) Get(ctx context.Context, id string) (entity *T, err error) {
	db := r.newSession(ctx)
	return r.GetInTransaction(ctx, db, id)
}

// GetInTransaction is a method that retrieves a single entity of type T from a database based on the provided
// ID within a transaction. The GORM library is used to execute the query in the database transaction and return
// the retrieved entity along with any errors encountered during the query. OpenTelemetry is used to add a span
// to the root span of the context, enabling tracing of the database query.
func (r *repository[T]) GetInTransaction(ctx context.Context, db *gorm.DB, id string) (entity *T, err error) {
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-get", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	// Using the "db" object to query the database and retrieve the first record that matches the given "id".
	// The result of the query is stored in the "entity" variable.
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
