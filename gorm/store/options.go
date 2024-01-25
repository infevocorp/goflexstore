package gormstore

import (
	"github.com/jkaveri/goflexstore/converter"
	gormquery "github.com/jkaveri/goflexstore/gorm/query"
	"github.com/jkaveri/goflexstore/store"
)

// Option is a function that modifies the store.
// It is used to set various configuration options for the Store at the time of its creation.
type Option[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] func(*Store[Entity, DTO, ID])

// WithBatchSize sets the batch size for batch operations in the store.
// batchSize specifies the number of records to be processed in a single batch during batch operations.
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

// WithConverter sets the converter used for transforming between entity and DTO types.
// converter is an instance of a converter that can convert between the entity and DTO types.
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

// WithScopeBuilderOption sets the scope builder options for the store.
// options are a variadic list of options that configure the behavior of the scope builder.
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
