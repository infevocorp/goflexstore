// Package gormquery provides utilities to construct GORM scopes based on query parameters
// defined in github.com/jkaveri/goflexstore/query. This package allows for flexible and reusable
// query building for GORM, enhancing code modularity and reusability.
//
// The package is mainly utilized in the github.com/jkaveri/flexstore/store/gorm package to
// create a generic, reusable store that interfaces with GORM for database operations.
package gormquery

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jkaveri/goflexstore/query"
)

// NewBuilder creates a new ScopeBuilder. It accepts various options that can modify the
// behavior of the scope builder, such as custom mappings between fields and database columns.
// This function initializes the ScopeBuilder with default handlers for different types of query
// parameters and applies any provided options to customize its behavior.
func NewBuilder(options ...Option) *ScopeBuilder {
	s := &ScopeBuilder{
		FieldToColMap: make(map[string]string),
		Registry:      make(ScopeBuilderRegistry),
		CustomFilters: make(map[string]ScopeBuilderFunc),
	}

	s.Registry = ScopeBuilderRegistry{
		query.TypeFilter:   s.Filter,
		query.TypeOR:       s.OR,
		query.TypePaginate: s.Paginate,
		query.TypeGroupBy:  s.GroupBy,
		query.TypeSelect:   s.Select,
		query.TypeOrderBy:  s.OrderBy,
		query.TypePreload:  s.Preload,
		query.TypeWithLock: s.ClauseLockUpdate,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// ScopeBuilder is a utility that constructs GORM scopes based on query parameters.
// It allows for mapping between field names and database column names, custom handling of query
// parameters, and registration of custom filter functions.
type ScopeBuilder struct {
	// FieldToColMap holds a mapping from struct field names to database column names.
	FieldToColMap map[string]string
	// Registry maps query parameter types to their corresponding scope builder functions.
	Registry ScopeBuilderRegistry
	// CustomFilters allows for the registration of custom filter functions.
	CustomFilters map[string]ScopeBuilderFunc
}

// Build constructs a slice of GORM scopes from the provided query parameters.
// It iterates through the query parameters and uses the registered scope builder functions
// to create corresponding GORM scopes.
func (b *ScopeBuilder) Build(params query.Params) []ScopeFunc {
	var scopes []ScopeFunc

	for _, param := range params.Params() {
		if builder, ok := b.Registry[param.ParamType()]; ok {
			scopes = append(scopes, builder(param))
		}
	}

	return scopes
}

// Filter constructs a GORM scope for a filter query parameter.
// It supports custom filters and converts the parameter into a GORM 'Where' clause.
func (b *ScopeBuilder) Filter(param query.Param) ScopeFunc {
	p := param.(query.FilterParam)

	// Run custom filter if available.
	if builder, ok := b.CustomFilters[p.Name]; ok {
		return builder(param)
	}

	col := b.getColName(p.Name)

	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(buildWhere(col, p.Operator, p.Value))
	}
}

// OR constructs a GORM scope for an OR query parameter.
// It creates a new GORM DB session and applies a series of 'Or' clauses based on the provided filters.
func (b *ScopeBuilder) OR(param query.Param) ScopeFunc {
	p := param.(query.ORParam)

	return func(tx *gorm.DB) *gorm.DB {
		db := tx.Session(&gorm.Session{NewDB: true})

		for i, filter := range p.Params {
			col := b.getColName(filter.Name)

			if i == 0 {
				db = db.Where(buildWhere(col, filter.Operator, filter.Value))
			} else {
				db = db.Or(buildWhere(col, filter.Operator, filter.Value))
			}
		}

		return tx.Where(db)
	}
}

// Paginate constructs a GORM scope for a paginate query parameter.
// It applies an offset and limit to the query based on the paginate parameters.
func (b *ScopeBuilder) Paginate(param query.Param) ScopeFunc {
	p := param.(query.PaginateParam)

	return func(tx *gorm.DB) *gorm.DB {
		return tx.Offset(p.Offset).Limit(p.Limit)
	}
}

// GroupBy constructs a GORM scope for a group by query parameter.
// It groups query results by specified columns and optionally applies 'Having' clauses.
func (b *ScopeBuilder) GroupBy(param query.Param) ScopeFunc {
	p := param.(query.GroupByParam)

	return func(tx *gorm.DB) *gorm.DB {
		cols := make([]string, len(p.Names))

		for i, name := range p.Names {
			cols[i] = b.getColName(name)
		}

		groupBy := strings.Join(cols, ", ")

		if p.Option != "" {
			groupBy += " " + p.Option
		}

		tx = tx.Group(groupBy)

		if len(p.Having) > 0 {
			for _, having := range p.Having {
				tx = tx.Having(buildWhere(
					b.getColName(having.Name),
					having.Operator,
					having.Value,
				))
			}
		}

		return tx
	}
}

// Select constructs a GORM scope for a select query parameter.
// It selects specific columns in the query based on the provided field names.
func (b *ScopeBuilder) Select(param query.Param) ScopeFunc {
	p := param.(query.SelectParam)

	return func(tx *gorm.DB) *gorm.DB {
		cols := make([]string, len(p.Names))

		for i, name := range p.Names {
			cols[i] = b.getColName(name)
		}

		return tx.Select(cols)
	}
}

// OrderBy constructs a GORM scope for an order by query parameter.
// It orders query results by a specified column in ascending or descending order.
func (b *ScopeBuilder) OrderBy(param query.Param) ScopeFunc {
	p := param.(query.OrderByParam)

	return func(tx *gorm.DB) *gorm.DB {
		col := b.getColName(p.Name)

		return tx.Order(clause.OrderByColumn{
			Column: clause.Column{Name: col},
			Desc:   p.Desc,
		})
	}
}

// Preload constructs a GORM scope for a preload query parameter.
// It preloads associations of the main query based on the provided field names and nested scopes.
func (b *ScopeBuilder) Preload(param query.Param) ScopeFunc {
	p := param.(query.PreloadParam)

	return func(tx *gorm.DB) *gorm.DB {
		if len(p.Params) == 0 {
			return tx.Preload(p.Name)
		}

		scopes := b.Build(query.NewParams(p.Params...))

		args := make([]any, len(scopes))

		for i := range scopes {
			args[i] = scopes[i]
		}

		return tx.Preload(p.Name, args...)
	}
}

// ClauseLockUpdate constructs a GORM scope for a locking clause query parameter.
// It adds a locking clause to the main query.
func (b *ScopeBuilder) ClauseLockUpdate(param query.Param) ScopeFunc {
	switch param.(query.WithLockParam).LockType {
	case query.LockTypeForUpdate:
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Clauses(clause.Locking{Strength: "UPDATE"})
		}
	default:
		return func(tx *gorm.DB) *gorm.DB {
			return tx
		}
	}
}

// getColName maps a field name to its corresponding column name in the database.
// If a mapping exists in FieldToColMap, it is used; otherwise, the field name itself is returned.
func (b *ScopeBuilder) getColName(name string) string {
	if col, ok := b.FieldToColMap[name]; ok {
		return col
	}

	return name
}
