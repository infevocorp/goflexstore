package sqlxstore

import (
	"github.com/infevocorp/goflexstore/converter"
	sqlxquery "github.com/infevocorp/goflexstore/sqlx/query"
	"github.com/infevocorp/goflexstore/store"
)

// Option customises a Store at construction time.
type Option[Entity store.Entity[ID], Row store.Entity[ID], ID comparable] func(*Store[Entity, Row, ID])

// WithBatchSize sets the batch size for CreateMany.
func WithBatchSize[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
	batchSize int,
) Option[Entity, Row, ID] {
	return func(s *Store[Entity, Row, ID]) {
		s.BatchSize = batchSize
	}
}

// WithConverter replaces the default reflect-based converter.
func WithConverter[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
	c converter.Converter[Entity, Row, ID],
) Option[Entity, Row, ID] {
	return func(s *Store[Entity, Row, ID]) {
		s.Converter = c
	}
}

// WithTable overrides the automatically derived table name.
func WithTable[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
	table string,
) Option[Entity, Row, ID] {
	return func(s *Store[Entity, Row, ID]) {
		s.Table = table
	}
}

// WithQueryBuilderOption replaces the default query builder with a newly
// constructed one configured with the given options.
func WithQueryBuilderOption[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
	opts ...sqlxquery.Option,
) Option[Entity, Row, ID] {
	return func(s *Store[Entity, Row, ID]) {
		s.QueryBuilder = sqlxquery.NewBuilder(opts...)
	}
}

// WithDialect sets the SQL dialect (affects Upsert SQL generation).
func WithDialect[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
	d sqlxquery.Dialect,
) Option[Entity, Row, ID] {
	return func(s *Store[Entity, Row, ID]) {
		s.Dialect = d
	}
}

// WithPKColumn overrides the primary-key column name (default: "id").
func WithPKColumn[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
	col string,
) Option[Entity, Row, ID] {
	return func(s *Store[Entity, Row, ID]) {
		s.PKColumn = col
	}
}

// WithReturningID enables the Postgres RETURNING <pk> clause for Create/Upsert.
func WithReturningID[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
	v bool,
) Option[Entity, Row, ID] {
	return func(s *Store[Entity, Row, ID]) {
		s.ReturningID = v
	}
}
