package gormstore

import (
	"github.com/jkaveri/goflexstore/converter"
	gormquery "github.com/jkaveri/goflexstore/gorm/query"
	"github.com/jkaveri/goflexstore/store"
)

type Option[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] func(*Store[Entity, DTO, ID])

func WithBatchSize[
	Entity store.Entity[ID],
	DTO store.Entity[ID],
	ID comparable,
](
	batchSize int,
) Option[Entity, DTO, ID] {
	return func(s *Store[Entity, DTO, ID]) {
		s.BatchSize = batchSize
	}
}

func WithConverter[
	Entity store.Entity[ID],
	DTO store.Entity[ID],
	ID comparable,
](
	converter converter.Converter[Entity, DTO, ID],
) Option[Entity, DTO, ID] {
	return func(s *Store[Entity, DTO, ID]) {
		s.Converter = converter
	}
}

func WithScopeBuilderOption[
	Entity store.Entity[ID],
	DTO store.Entity[ID],
	ID comparable,
](
	options ...gormquery.Option,
) Option[Entity, DTO, ID] {
	return func(s *Store[Entity, DTO, ID]) {
		s.ScopeBuilder = gormquery.NewBuilder(options...)
	}
}
