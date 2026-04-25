package sqlxquery_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/infevocorp/goflexstore/query"
	sqlxquery "github.com/infevocorp/goflexstore/sqlx/query"
)

func newBuilder() *sqlxquery.Builder {
	return sqlxquery.NewBuilder(
		sqlxquery.WithFieldToColMap(map[string]string{
			"ID":   "id",
			"Name": "name",
			"Age":  "age",
		}),
	)
}

func TestBuild_Filter(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(query.Filter("ID", int64(1))))

	assert.Equal(t, "id = ?", r.Where)
	assert.Equal(t, []any{int64(1)}, r.Args)
}

func TestBuild_FilterOperators(t *testing.T) {
	b := newBuilder()

	cases := []struct {
		op      query.Operator
		wantSQL string
	}{
		{query.EQ, "age = ?"},
		{query.NEQ, "age <> ?"},
		{query.GT, "age > ?"},
		{query.GTE, "age >= ?"},
		{query.LT, "age < ?"},
		{query.LTE, "age <= ?"},
	}

	for _, tc := range cases {
		r := b.Build(query.NewParams(query.Filter("Age", 30).WithOP(tc.op)))
		assert.Equal(t, tc.wantSQL, r.Where)
	}
}

func TestBuild_FilterInSlice(t *testing.T) {
	b := newBuilder()
	ids := []int64{1, 2, 3}
	r := b.Build(query.NewParams(query.Filter("ID", ids)))

	assert.Equal(t, "id IN (?)", r.Where)
	require.Len(t, r.Args, 1)
	assert.Equal(t, ids, r.Args[0])
}

func TestBuild_OR(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(
		query.OR(query.Filter("Name", "Alice"), query.Filter("Name", "Bob")),
	))

	assert.Equal(t, "(name = ? OR name = ?)", r.Where)
	assert.Equal(t, []any{"Alice", "Bob"}, r.Args)
}

func TestBuild_Paginate(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(query.Paginate(10, 5)))

	assert.Equal(t, 10, r.Offset)
	assert.Equal(t, 5, r.Limit)
}

func TestBuild_OrderBy(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(
		query.OrderBy("ID", false),
		query.OrderBy("Name", true),
	))

	assert.Equal(t, "id ASC, name DESC", r.OrderBy)
}

func TestBuild_Select(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(query.Select("ID", "Name")))

	assert.Equal(t, []string{"id", "name"}, r.Cols)
}

func TestBuild_GroupBy(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(
		query.GroupBy("Age").WithHaving(query.Filter("Age", 18).WithOP(query.GT)),
	))

	assert.Equal(t, "age", r.GroupBy)
	assert.Equal(t, "age > ?", r.Having)
	assert.Equal(t, []any{18}, r.Args)
}

func TestBuild_WithLock(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(query.WithLock(query.LockTypeForUpdate)))

	assert.Equal(t, "FOR UPDATE", r.Suffix)
}

func TestBuild_WithHint(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(query.WithHint("index_merge")))

	assert.Equal(t, "/*+ index_merge */", r.Hint)
}

func TestBuild_MultipleConditions(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(
		query.Filter("ID", int64(1)),
		query.Filter("Age", 30).WithOP(query.GT),
	))

	assert.Equal(t, "id = ? AND age > ?", r.Where)
	assert.Equal(t, []any{int64(1), 30}, r.Args)
}

func TestBuild_UnknownField_PassThrough(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(query.Filter("unknown_field", "val")))

	assert.Equal(t, "unknown_field = ?", r.Where)
}

func TestBuild_ArgsOrder_WhereBeforeHaving(t *testing.T) {
	b := newBuilder()
	r := b.Build(query.NewParams(
		query.Filter("ID", int64(1)),
		query.GroupBy("Age").WithHaving(query.Filter("Age", 18).WithOP(query.GTE)),
	))

	// WHERE arg (1) must come before HAVING arg (18)
	assert.Equal(t, []any{int64(1), 18}, r.Args)
}
