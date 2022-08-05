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

// FindOptions ...
type FindOptions struct {
	Limit int
	Page  int
}

// Repository ...
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Find(ctx context.Context, entity interface{}, options *FindOptions) (list []*T, err error)
	Count(ctx context.Context, entity interface{}) (count int64, err error)
	Get(ctx context.Context, id string) (entity *T, err error)
	Update(ctx context.Context, entity T) (updated *T, err error)
	Delete(ctx context.Context, id string) (deleted bool, err error)
}

// Repository ...
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
	// add the span of database query on the root span of the context
	span := utils.GenerateSpanTracer(ctx, create)
	defer span()

	db := r.newSession(ctx)

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
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-find", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	db := r.newSession(ctx)

	if options != nil {
		limit := 10
		if options.Limit != 0 {
			limit = options.Limit
		}
		db.Limit(limit)
		if options.Page != 0 {
			page := options.Page
			if page == 0 {
				page = 1
			}
			db.Offset((page - 1) * limit)
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
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-count", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	db := r.newSession(ctx)

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
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-get", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	db := r.newSession(ctx)

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
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-update", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	db := r.newSession(ctx)

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
	// add the span of database query on the root span of the context
	tracer := telemetry.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "repo-delete", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	db := r.newSession(ctx)

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
	return r.db.Session(&gorm.Session{}).Model(new(T)).WithContext(ctx)
}
