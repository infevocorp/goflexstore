package gormstore

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jkaveri/goflexstore/converter"
	gormopscope "github.com/jkaveri/goflexstore/gorm/opscope"
	gormquery "github.com/jkaveri/goflexstore/gorm/query"
	gormutils "github.com/jkaveri/goflexstore/gorm/utils"
	"github.com/jkaveri/goflexstore/query"
	"github.com/jkaveri/goflexstore/store"
)

// New initializes a new Store instance for handling CRUD operations on entities.
// It accepts an operation scope and a variable number of options to customize the store behavior.
// The function returns a pointer to the initialized Store.
//
// Entity and DTO are types that must implement the store.Entity interface.
// ID is the type of the identifier for the entities.
func New[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable](
	opScope *gormopscope.TransactionScope,
	options ...Option[Entity, DTO, ID],
) *Store[Entity, DTO, ID] {
	s := &Store[Entity, DTO, ID]{
		OpScope:   opScope,
		BatchSize: 50,
	}

	for _, option := range options {
		option(s)
	}

	if s.Converter == nil {
		s.Converter = converter.NewReflect[Entity, DTO, ID](nil)
	}

	if s.ScopeBuilder == nil {
		s.ScopeBuilder = gormquery.NewBuilder(
			gormquery.WithFieldToColMap(
				gormutils.FieldToColMap(*new(DTO)),
			),
		)
	}

	return s
}

// Store represents a storage mechanism using GORM for database operations.
// It supports CRUD operations and is designed to be generic for any Entity and DTO types.
//
// Entity: The domain model type.
// DTO: The data transfer object type, representing the database model.
// ID: The type of the unique identifier for the entity.
type Store[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] struct {
	OpScope      *gormopscope.TransactionScope
	Converter    converter.Converter[Entity, DTO, ID]
	ScopeBuilder *gormquery.ScopeBuilder
	BatchSize    int
}

// Get retrieves a single entity based on provided query parameters.
// It returns the entity if found, otherwise an error.
func (s *Store[Entity, DTO, ID]) Get(ctx context.Context, params ...query.Param) (Entity, error) {
	var (
		dto    DTO
		scopes = s.ScopeBuilder.Build(query.NewParams(params...))
	)

	if err := s.getTx(ctx).
		Scopes(scopes...).
		First(&dto).Error; err != nil {
		return *new(Entity), nil
	}

	return s.Converter.ToEntity(dto), nil
}

// List retrieves a list of entities matching the provided query parameters.
// Returns a slice of entities and an error if the operation fails.
func (s *Store[Entity, DTO, ID]) List(ctx context.Context, params ...query.Param) ([]Entity, error) {
	var (
		dtos   []DTO
		scopes = s.ScopeBuilder.Build(query.NewParams(params...))
	)

	if err := s.getTx(ctx).
		Scopes(scopes...).Find(&dtos).Error; err != nil {
		return nil, err
	}

	return converter.ToMany(dtos, s.Converter.ToEntity), nil
}

// Count returns the number of entities that satisfy the provided query parameters.
// The count is returned along with an error if the operation fails.
func (s *Store[Entity, DTO, ID]) Count(ctx context.Context, params ...query.Param) (int64, error) {
	var (
		count  int64
		scopes = s.ScopeBuilder.Build(query.NewParams(params...))
	)

	if err := s.getTx(ctx).
		Scopes(scopes...).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Exists checks for the existence of at least one entity that matches the query parameters.
// Returns true if such an entity exists, false otherwise.
func (s *Store[Entity, DTO, ID]) Exists(ctx context.Context, params ...query.Param) (bool, error) {
	var (
		count  int64
		scopes = s.ScopeBuilder.Build(query.NewParams(params...))
	)

	if err := s.getTx(ctx).Scopes(scopes...).
		Limit(1).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Create adds a new entity to the store and returns its ID.
// Returns an error if the creation fails.
func (s *Store[Entity, DTO, ID]) Create(ctx context.Context, entity Entity) (ID, error) {
	dto := s.Converter.ToDTO(entity)
	if err := s.getTx(ctx).Create(&dto).Error; err != nil {
		return *new(ID), err
	}

	return dto.GetID(), nil
}

// CreateMany performs batch creation of entities.
// The BatchSize field of the store determines the number of entities in each batch.
// Returns an error if the operation fails.
func (s *Store[Entity, DTO, ID]) CreateMany(ctx context.Context, entities []Entity) error {
	dtos := converter.ToMany(entities, s.Converter.ToDTO)
	batchSize := defaultValue(s.BatchSize, 50)

	return s.getTx(ctx).CreateInBatches(dtos, batchSize).Error
}

// Update modifies an existing entity in the store, including fields with zero values.
// Returns an error if the update operation fails.
func (s *Store[Entity, DTO, ID]) Update(ctx context.Context, entity Entity, params ...query.Param) error {
	dto := s.Converter.ToDTO(entity)
	id := dto.GetID()

	if id == *new(ID) && len(params) == 0 {
		return errors.New("id is required")
	}

	tx := s.getTx(ctx)

	if len(params) > 0 {
		scopes := s.ScopeBuilder.Build(query.NewParams(params...))
		tx = tx.Scopes(scopes...)
	}

	return tx.Select("*").Updates(&dto).Error
}

// PartialUpdate updates specific fields of an existing entity in the store.
// Only non-zero fields of the entity are updated.
// Returns an error if the operation fails.
func (s *Store[Entity, DTO, ID]) PartialUpdate(ctx context.Context, entity Entity, params ...query.Param) error {
	dto := s.Converter.ToDTO(entity)
	scopes := s.ScopeBuilder.Build(query.NewParams(params...))

	return s.getTx(ctx).Scopes(scopes...).Updates(dto).Error
}

// Delete removes entities from the store based on the provided query parameters.
// Returns an error if the deletion operation fails.
func (s *Store[Entity, DTO, ID]) Delete(ctx context.Context, params ...query.Param) error {
	var (
		dto    DTO
		scopes = s.ScopeBuilder.Build(query.NewParams(params...))
	)

	if err := s.getTx(ctx).
		Scopes(scopes...).
		Delete(&dto).Error; err != nil {
		return err
	}

	return nil
}

// Upsert either creates a new entity or updates an existing one based on the provided conflict resolution strategy.
// Returns the ID of the affected entity and an error if the operation fails.
func (s *Store[Entity, DTO, ID]) Upsert(ctx context.Context, entity Entity, onConflict store.OnConflict) (ID, error) {
	dto := s.Converter.ToDTO(entity)
	c := clause.OnConflict{
		Columns:      []clause.Column{},
		OnConstraint: onConflict.OnConstraint,
		DoNothing:    onConflict.DoNothing,
		UpdateAll:    onConflict.UpdateAll,
	}

	for _, col := range onConflict.Columns {
		c.Columns = append(c.Columns, clause.Column{Name: col})
	}

	if len(onConflict.Updates) > 0 {
		c.DoUpdates = clause.Assignments(onConflict.Updates)
	} else if len(onConflict.UpdateColumns) > 0 {
		c.DoUpdates = clause.AssignmentColumns(onConflict.UpdateColumns)
	}

	if err := s.getTx(ctx).Clauses(c).Create(&dto).Error; err != nil {
		return *new(ID), err
	}

	return dto.GetID(), nil
}

func (s *Store[Entity, DTO, ID]) getTx(ctx context.Context) *gorm.DB {
	return s.OpScope.Tx(ctx).WithContext(ctx).Model(new(DTO))
}
