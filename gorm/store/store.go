package gormstore

import (
	"context"

	"gorm.io/gorm"

	"github.com/jkaveri/goflexstore/converter"
	gormopscope "github.com/jkaveri/goflexstore/gorm/opscope"
	gormquery "github.com/jkaveri/goflexstore/gorm/query"
	"github.com/jkaveri/goflexstore/query"
	"github.com/jkaveri/goflexstore/store"
)

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
		s.Converter = converter.NewReflect[Entity, DTO](nil)
	}

	if s.ScopeBuilder == nil {
		s.ScopeBuilder = gormquery.NewBuilder(
			gormquery.WithFieldToColMap(
				FieldToColMap(new(DTO)),
			),
		)
	}

	return s
}

type Store[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] struct {
	OpScope      *gormopscope.TransactionScope
	Converter    converter.Converter[Entity, DTO, ID]
	ScopeBuilder *gormquery.ScopeBuilder
	BatchSize    int
}

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

func (s *Store[Entity, DTO, ID]) Create(ctx context.Context, entity Entity) (ID, error) {
	dto := s.Converter.ToDTO(entity)
	if err := s.getTx(ctx).Create(&dto).Error; err != nil {
		return *new(ID), err
	}

	return dto.GetID(), nil
}

func (s *Store[Entity, DTO, ID]) CreateMany(ctx context.Context, entities []Entity) error {
	dtos := converter.ToMany(entities, s.Converter.ToDTO)
	batchSize := defaultValue(s.BatchSize, 50)

	return s.getTx(ctx).CreateInBatches(dtos, batchSize).Error
}

func (s *Store[Entity, DTO, ID]) Update(ctx context.Context, entity Entity, params ...query.Param) error {
	dto := s.Converter.ToDTO(entity)
	scopes := s.ScopeBuilder.Build(query.NewParams(params...))
	id := dto.GetID()

	if id == (*new(ID)) {
		return s.getTx(ctx).Scopes(scopes...).Updates(dto).Error
	}

	return s.getTx(ctx).Scopes(scopes...).Save(&dto).Error
}

func (s *Store[Entity, DTO, ID]) PartialUpdate(ctx context.Context, entity Entity, params ...query.Param) error {
	dto := s.Converter.ToDTO(entity)
	scopes := s.ScopeBuilder.Build(query.NewParams(params...))

	return s.getTx(ctx).Scopes(scopes...).Updates(dto).Error
}

func (s *Store[Entity, DTO, ID]) getTx(ctx context.Context) *gorm.DB {
	return s.OpScope.Tx(ctx).WithContext(ctx).Model(new(DTO))
}
