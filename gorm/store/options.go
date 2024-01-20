package gormstore

import (
	"github.com/jkaveri/goflexstore/converter"
	gormquery "github.com/jkaveri/goflexstore/gorm/query"
	"github.com/jkaveri/goflexstore/store"
)

// Option is a function that modifies the store
type Option[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] func(*Store[Entity, DTO, ID])

// WithBatchSize sets the batch size
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

// WithConverter sets the converter
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

// WithScopeBuilderOption sets the scope builder options
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
