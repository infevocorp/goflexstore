package gormquery

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jkaveri/goflexstore/query"
)

// NewBuilder creates new scope builder
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
	if builder, ok := b.CustomFilters[p.Name]; ok {
		return builder(param)
	}

	col := b.getColName(p.Name)

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

func (b *ScopeBuilder) Paginate(param query.Param) ScopeFunc {
	p := param.(query.PaginateParam)

	return func(tx *gorm.DB) *gorm.DB {
		return tx.Offset(p.Offset).Limit(p.Limit)
	}
}

func (b *ScopeBuilder) GroupBy(param query.Param) ScopeFunc {
	p := param.(query.GroupByParam)

	return func(tx *gorm.DB) *gorm.DB {
		cols := make([]string, len(p.Names))

		for i, name := range p.Names {
			cols[i] = b.getColName(name)
		}

		groupBy := strings.Join(cols, ", ")

		if p.Option != "" {
			groupBy = groupBy + " " + p.Option
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

func (b *ScopeBuilder) Preload(param query.Param) ScopeFunc {
	p := param.(query.PreloadParam)

	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Preload(p.Name)

		if len(p.Params) > 0 {
			for _, param := range p.Params {
				if builder, ok := b.Registry[param.ParamType()]; ok {
					tx = builder(param)(tx)
				}
			}
		}

		return tx
	}
}

func (b *ScopeBuilder) getColName(name string) string {
	if col, ok := b.FieldToColMap[name]; ok {
		return col
	}

	return name
}
