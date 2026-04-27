// Package sqlxquery translates goflexstore query.Params into SQL fragments
// (WHERE, ORDER BY, GROUP BY, etc.) for use with jmoiron/sqlx.
//
// The builder always emits `?` placeholders; callers rebind to the
// target dialect using the driver's Rebind method or the Rebind helper.
package sqlxquery

import (
	"reflect"
	"strings"

	"github.com/infevocorp/goflexstore/query"
)

// Result holds the SQL fragments produced by Builder.Build.
type Result struct {
	Hint    string   // SQL comment hint prepended to the query, e.g. "/*+ index */"
	Where   string   // WHERE clause body (no WHERE keyword)
	GroupBy string   // GROUP BY clause body
	Having  string   // HAVING clause body
	OrderBy string   // ORDER BY clause body
	Limit   int      // 0 = no limit
	Offset  int
	Cols    []string // nil = SELECT *
	Suffix  string   // appended verbatim, e.g. "FOR UPDATE"
	Args    []any    // positional args: WHERE args then HAVING args
}

// Builder translates query.Params into SQL fragments.
type Builder struct {
	FieldToColMap map[string]string
	Dialect       Dialect
}

// NewBuilder creates a Builder, applying any provided options.
func NewBuilder(opts ...Option) *Builder {
	b := &Builder{
		FieldToColMap: make(map[string]string),
		Dialect:       DialectMySQL,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Build walks params and fills a Result.
func (b *Builder) Build(params query.Params) Result {
	var (
		r           Result
		havingParts []string
		havingArgs  []any
		groupByCols []string
		orderByParts []string
		groupByOpt  string
	)

	allParams := params.Params()
	whereParts := make([]string, 0, len(allParams))
	whereArgs  := make([]any,   0, len(allParams))

	for _, param := range allParams {
		switch p := param.(type) {
		case query.FilterParam:
			where := buildWhere(b.getColName(p.Name), p.Operator, p.Value, &whereArgs)
			whereParts = append(whereParts, where)

		case query.ORParam:
			parts := make([]string, len(p.Params))
			for i, f := range p.Params {
				parts[i] = buildWhere(b.getColName(f.Name), f.Operator, f.Value, &whereArgs)
			}
			whereParts = append(whereParts, "("+strings.Join(parts, " OR ")+")")

		case query.PaginateParam:
			r.Limit = p.Limit
			r.Offset = p.Offset

		case query.OrderByParam:
			col := b.getColName(p.Name)
			dir := "ASC"
			if p.Desc {
				dir = "DESC"
			}
			orderByParts = append(orderByParts, col+" "+dir)

		case query.SelectParam:
			cols := make([]string, len(p.Names))
			for i, name := range p.Names {
				cols[i] = b.getColName(name)
			}
			r.Cols = cols

		case query.GroupByParam:
			for _, name := range p.Names {
				groupByCols = append(groupByCols, b.getColName(name))
			}
			groupByOpt = p.Option
			for _, hf := range p.Having {
				w := buildWhere(b.getColName(hf.Name), hf.Operator, hf.Value, &havingArgs)
				havingParts = append(havingParts, w)
			}

		case query.WithLockParam:
			if p.LockType == query.LockTypeForUpdate {
				if r.Suffix != "" {
					r.Suffix += " "
				}
				r.Suffix += "FOR UPDATE"
			}

		case query.WithHintParam:
			r.Hint = "/*+ " + p.Hint + " */"

		case query.PreloadParam:
			// not supported; callers should check and return ErrPreloadNotSupported
		}
	}

	r.Where = strings.Join(whereParts, " AND ")
	r.Having = strings.Join(havingParts, " AND ")
	r.OrderBy = strings.Join(orderByParts, ", ")
	r.Args = append(whereArgs, havingArgs...)

	if len(groupByCols) > 0 {
		r.GroupBy = strings.Join(groupByCols, ", ")
		if groupByOpt != "" {
			r.GroupBy += " " + groupByOpt
		}
	}

	return r
}

func (b *Builder) getColName(name string) string {
	if col, ok := b.FieldToColMap[name]; ok {
		return col
	}
	return name
}

// buildWhere constructs a single WHERE predicate and appends the argument(s)
// to args. Returns the SQL fragment.
func buildWhere(col string, op query.Operator, value any, args *[]any) string {
	if value == nil {
		panic("filter value cannot be nil")
	}

	rv := reflect.ValueOf(value)
	if k := rv.Kind(); k == reflect.Slice || k == reflect.Array {
		n := rv.Len()
		if n > 1 {
			// Pass the whole slice as a single arg; sqlx.In() will expand it.
			*args = append(*args, value)
			return col + " " + inOperatorStr(op) + " (?)"
		}
		*args = append(*args, rv.Index(0).Interface())
		return col + " " + operatorStr(op) + " ?"
	}

	*args = append(*args, value)
	return col + " " + operatorStr(op) + " ?"
}

func operatorStr(op query.Operator) string {
	switch op {
	case query.EQ:
		return "="
	case query.NEQ:
		return "<>"
	case query.GT:
		return ">"
	case query.GTE:
		return ">="
	case query.LT:
		return "<"
	case query.LTE:
		return "<="
	default:
		return "UNKNOWN"
	}
}

func inOperatorStr(op query.Operator) string {
	switch op {
	case query.EQ:
		return "IN"
	case query.NEQ:
		return "NOT IN"
	default:
		panic("unsupported operator for IN clause")
	}
}
