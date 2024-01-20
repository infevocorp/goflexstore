package store

import (
	"context"

	"github.com/jkaveri/goflexstore/query"
)

// Entity a generic interface for an Entity model.
type Entity[ID comparable] interface {
	// GetID returns the unique identifier of this entity.
	GetID() ID
}

// Store a generic store interface for an Entity model.
type Store[T Entity[ID], ID comparable] interface {
	// Get returns a single entity based on the provided params.
	Get(ctx context.Context, params ...query.Param) (T, error)
	// List returns a list of entities based on the provided params.
	List(ctx context.Context, params ...query.Param) ([]T, error)
	// Count returns the number of entities based on the provided params.
	Count(ctx context.Context, params ...query.Param) (int, error)
	// Exist returns true if the entity exists based on the provided params.
	Exist(ctx context.Context, params ...query.Param) (bool, error)
	// Create creates a new entity.
	Create(ctx context.Context, entity T) (ID, error)
	// CreateMany creates multiple entities.
	CreateMany(ctx context.Context, entities []T) error
	// Update updates an existing entity.
	Update(ctx context.Context, entity T, params ...query.Param) error
	// PartialUpdate updates an existing entity partially.
	PartialUpdate(ctx context.Context, entity T, params ...query.Param) error
	// Upsert creates or updates an existing entity.
	Upsert(ctx context.Context, entity T, params ...query.Param) error
	// Delete deletes an existing entity.
	Delete(ctx context.Context, params ...query.Params) (T, error)
}
