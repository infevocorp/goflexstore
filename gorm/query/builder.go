package gormquery

import (
	"gorm.io/gorm"

	"github.com/jkaveri/goflexstore/query"
)

func NewBuilder(options ...Option) *ScopeBuilder {
	s := &ScopeBuilder{
		FieldToColMap: make(map[string]string),
		Registry:      make(ScopeBuilderRegistry),
		CustomFilters: make(map[string]ScopeBuilderFunc),
	}

	s.Registry = ScopeBuilderRegistry{
		"filter":   s.Filter,
		"or":       s.OR,
		"paginate": s.Paginate,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

type ScopeBuilder struct {
	FieldToColMap map[string]string
	Registry      ScopeBuilderRegistry
	CustomFilters map[string]ScopeBuilderFunc
}

func (b *ScopeBuilder) Build(params query.Params) []ScopeFunc {
	var scopes []ScopeFunc

	for _, param := range params.Params() {
		if builder, ok := b.Registry[param.ParamType()]; ok {
			scopes = append(scopes, builder(param))
		}
	}

	return scopes
}

func (b *ScopeBuilder) Filter(param query.Param) ScopeFunc {
	p := param.(query.FilterParam)

	// run custom filter
	if builder, ok := b.CustomFilters[p.FieldName]; ok {
		return builder(param)
	}

	col, ok := b.FieldToColMap[p.FieldName]
	if !ok {
		col = p.FieldName
	}

	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(buildWhere(col, p.Operator, p.Value))
	}
}

func (b *ScopeBuilder) OR(param query.Param) ScopeFunc {
	p := param.(query.ORParam)

	return func(tx *gorm.DB) *gorm.DB {
		db := tx.Session(&gorm.Session{
			NewDB: true,
		})

		for i, filter := range p.Params {
			col, ok := b.FieldToColMap[filter.FieldName]
			if !ok {
				col = filter.FieldName
			}

			if i == 0 {
				db = db.Where(buildWhere(col, filter.Operator, filter.Value))
			} else {
				db = db.Or(buildWhere(col, filter.Operator, filter.Value))
			}
		}

		return tx.Where(db)
	}
}

func (b *ScopeBuilder) Paginate(param query.Param) ScopeFunc {
	p := param.(query.PaginateParam)

	return func(tx *gorm.DB) *gorm.DB {
		return tx.Offset(p.Offset).Limit(p.Limit)
	}
}
