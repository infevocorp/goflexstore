package integration_test

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/infevocorp/goflexstore/query"
	"github.com/infevocorp/goflexstore/store"
	sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
	sqlxquery "github.com/infevocorp/goflexstore/sqlx/query"
	sqlxstore "github.com/infevocorp/goflexstore/sqlx/store"
)

// ---- fixtures ---------------------------------------------------------------

type Product struct {
	ID       int64
	Name     string
	Price    int64
	Category string
}

func (p Product) GetID() int64 { return p.ID }

type ProductRow struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Price    int64  `db:"price"`
	Category string `db:"category"`
}

func (r ProductRow) GetID() int64 { return r.ID }

// Pointer-type variants for pointer-entity tests.
type PtrProduct struct {
	ID       int64
	Name     string
	Price    int64
	Category string
}

func (p *PtrProduct) GetID() int64 { return p.ID }

type PtrProductRow struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Price    int64  `db:"price"`
	Category string `db:"category"`
}

func (r *PtrProductRow) GetID() int64 { return r.ID }

// ---- helpers ----------------------------------------------------------------

const createTable = `CREATE TABLE products (
	id       INTEGER PRIMARY KEY AUTOINCREMENT,
	name     TEXT    NOT NULL UNIQUE,
	price    INTEGER NOT NULL DEFAULT 0,
	category TEXT    NOT NULL DEFAULT ''
)`

func newDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	_, err = db.Exec(createTable)
	require.NoError(t, err)
	return db
}

func newStore(db *sqlx.DB) *sqlxstore.Store[Product, ProductRow, int64] {
	opScope := sqlxopscope.NewTransactionScope("test", db, nil)
	return sqlxstore.New[Product, ProductRow, int64](
		opScope,
		sqlxstore.WithTable[Product, ProductRow, int64]("products"),
		sqlxstore.WithDialect[Product, ProductRow, int64](sqlxquery.DialectSQLite),
	)
}

func newPtrStore(db *sqlx.DB) *sqlxstore.Store[*PtrProduct, *PtrProductRow, int64] {
	opScope := sqlxopscope.NewTransactionScope("test", db, nil)
	return sqlxstore.New[*PtrProduct, *PtrProductRow, int64](
		opScope,
		sqlxstore.WithTable[*PtrProduct, *PtrProductRow, int64]("products"),
		sqlxstore.WithDialect[*PtrProduct, *PtrProductRow, int64](sqlxquery.DialectSQLite),
	)
}

func seed(t *testing.T, s *sqlxstore.Store[Product, ProductRow, int64], products ...Product) []int64 {
	t.Helper()
	ids := make([]int64, len(products))
	for i, p := range products {
		id, err := s.Create(context.Background(), p)
		require.NoError(t, err)
		ids[i] = id
	}
	return ids
}

// ---- Filter operators (GT / GTE / LT / LTE / NEQ) --------------------------

func TestFilter_GT(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "a", Price: 50, Category: "x"},
		Product{Name: "b", Price: 100, Category: "x"},
		Product{Name: "c", Price: 150, Category: "x"},
		Product{Name: "d", Price: 200, Category: "x"},
	)

	results, err := s.List(ctx, query.Filter("Price", int64(100)).WithOP(query.GT))
	require.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.Greater(t, r.Price, int64(100))
	}
}

func TestFilter_GTE(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "a", Price: 50, Category: "x"},
		Product{Name: "b", Price: 100, Category: "x"},
		Product{Name: "c", Price: 150, Category: "x"},
		Product{Name: "d", Price: 200, Category: "x"},
	)

	results, err := s.List(ctx, query.Filter("Price", int64(100)).WithOP(query.GTE))
	require.NoError(t, err)
	assert.Len(t, results, 3)
}

func TestFilter_LT(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "a", Price: 50, Category: "x"},
		Product{Name: "b", Price: 100, Category: "x"},
		Product{Name: "c", Price: 150, Category: "x"},
	)

	results, err := s.List(ctx, query.Filter("Price", int64(100)).WithOP(query.LT))
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, int64(50), results[0].Price)
}

func TestFilter_LTE(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "a", Price: 50, Category: "x"},
		Product{Name: "b", Price: 100, Category: "x"},
		Product{Name: "c", Price: 150, Category: "x"},
	)

	results, err := s.List(ctx, query.Filter("Price", int64(100)).WithOP(query.LTE))
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestFilter_NEQ(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "a", Price: 10, Category: "A"},
		Product{Name: "b", Price: 20, Category: "B"},
		Product{Name: "c", Price: 30, Category: "A"},
	)

	results, err := s.List(ctx, query.Filter("Category", "A").WithOP(query.NEQ))
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "B", results[0].Category)
}

// ---- OR conditions ----------------------------------------------------------

func TestList_OR(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "X", Price: 10, Category: "c"},
		Product{Name: "Y", Price: 20, Category: "c"},
		Product{Name: "Z", Price: 30, Category: "c"},
	)

	results, err := s.List(ctx, query.OR(
		query.Filter("Name", "X"),
		query.Filter("Name", "Y"),
	))
	require.NoError(t, err)
	require.Len(t, results, 2)
	names := []string{results[0].Name, results[1].Name}
	assert.Contains(t, names, "X")
	assert.Contains(t, names, "Y")
}

func TestCount_OR(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "X", Price: 10, Category: "c"},
		Product{Name: "Y", Price: 20, Category: "c"},
		Product{Name: "Z", Price: 30, Category: "c"},
	)

	n, err := s.Count(ctx, query.OR(
		query.Filter("Name", "X"),
		query.Filter("Name", "Z"),
	))
	require.NoError(t, err)
	assert.Equal(t, int64(2), n)
}

func TestExists_OR(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s, Product{Name: "X", Price: 10, Category: "c"})

	ok, err := s.Exists(ctx, query.OR(
		query.Filter("Name", "X"),
		query.Filter("Name", "missing"),
	))
	require.NoError(t, err)
	assert.True(t, ok)
}

// ---- Combined AND filters ---------------------------------------------------

func TestList_CombinedFilters(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "cheap-book", Price: 100, Category: "books"},
		Product{Name: "pricey-book", Price: 300, Category: "books"},
		Product{Name: "tv", Price: 150, Category: "electronics"},
	)

	results, err := s.List(ctx,
		query.Filter("Category", "books"),
		query.Filter("Price", int64(200)).WithOP(query.LT),
	)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "cheap-book", results[0].Name)
}

// ---- GroupBy + Having -------------------------------------------------------

func TestGroupBy_Having(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "a1", Price: 10, Category: "A"},
		Product{Name: "a2", Price: 20, Category: "A"},
		Product{Name: "z1", Price: 30, Category: "Z"},
	)

	// SELECT * FROM products GROUP BY category HAVING category <> 'Z'
	results, err := s.List(ctx,
		query.GroupBy("Category").WithHaving(
			query.Filter("Category", "Z").WithOP(query.NEQ),
		),
	)
	require.NoError(t, err)
	for _, r := range results {
		assert.NotEqual(t, "Z", r.Category)
	}
}

// ---- Pointer-type entity ----------------------------------------------------

func TestPointer_Create(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newPtrStore(db)
	ctx := context.Background()

	id, err := s.Create(ctx, &PtrProduct{Name: "ptr-item", Price: 42, Category: "p"})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))
}

func TestPointer_Get(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newPtrStore(db)
	ctx := context.Background()

	_, err := s.Create(ctx, &PtrProduct{Name: "ptr-item", Price: 42, Category: "p"})
	require.NoError(t, err)

	p, err := s.Get(ctx, query.Filter("Name", "ptr-item"))
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, "ptr-item", p.Name)
	assert.Equal(t, int64(42), p.Price)
}

func TestPointer_List(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newPtrStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, &PtrProduct{Name: "p1", Price: 10, Category: "p"})
	_, _ = s.Create(ctx, &PtrProduct{Name: "p2", Price: 20, Category: "p"})

	results, err := s.List(ctx)
	require.NoError(t, err)
	require.Len(t, results, 2)
	for _, r := range results {
		require.NotNil(t, r)
	}
}

// ---- Upsert with UpdateColumns ----------------------------------------------

func TestUpsert_UpdateColumns(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	_, err := s.Create(ctx, Product{Name: "gadget", Price: 100, Category: "electronics"})
	require.NoError(t, err)

	// Conflict on "name" → update only "price"; category must stay unchanged.
	_, err = s.Upsert(ctx,
		Product{Name: "gadget", Price: 999, Category: "changed"},
		store.OnConflict{
			Columns:       []string{"name"},
			UpdateColumns: []string{"price"},
		},
	)
	require.NoError(t, err)

	updated, err := s.Get(ctx, query.Filter("Name", "gadget"))
	require.NoError(t, err)
	assert.Equal(t, int64(999), updated.Price)
	assert.Equal(t, "electronics", updated.Category)
}

// ---- Large-batch CreateMany -------------------------------------------------

func TestCreateMany_LargeBatch(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	const total = 150
	products := make([]Product, total)
	for i := range products {
		products[i] = Product{
			Name:     fmt.Sprintf("item-%03d", i),
			Price:    int64(i),
			Category: "bulk",
		}
	}

	err := s.CreateMany(ctx, products)
	require.NoError(t, err)

	n, err := s.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(total), n)
}

// ---- PartialUpdate by param -------------------------------------------------

func TestPartialUpdate_ByParam(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "book1", Price: 10, Category: "books"},
		Product{Name: "book2", Price: 20, Category: "books"},
		Product{Name: "tv", Price: 500, Category: "electronics"},
	)

	// Update only non-zero fields (Price=999) for rows where category = "books".
	err := s.PartialUpdate(ctx,
		Product{Price: 999},
		query.Filter("Category", "books"),
	)
	require.NoError(t, err)

	books, err := s.List(ctx, query.Filter("Category", "books"))
	require.NoError(t, err)
	for _, b := range books {
		assert.Equal(t, int64(999), b.Price)
	}

	tv, err := s.Get(ctx, query.Filter("Name", "tv"))
	require.NoError(t, err)
	assert.Equal(t, int64(500), tv.Price)
}

// ---- OrderBy DESC -----------------------------------------------------------

func TestList_OrderByDesc(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "cheap", Price: 10, Category: "x"},
		Product{Name: "expensive", Price: 500, Category: "x"},
		Product{Name: "mid", Price: 200, Category: "x"},
	)

	results, err := s.List(ctx, query.OrderBy("Price", true))
	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Equal(t, int64(500), results[0].Price)
	assert.Equal(t, int64(200), results[1].Price)
	assert.Equal(t, int64(10), results[2].Price)
}

// ---- Delete with multiple conditions ----------------------------------------

func TestDelete_MultipleConditions(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	s := newStore(db)
	ctx := context.Background()

	seed(t, s,
		Product{Name: "a1", Price: 100, Category: "A"},
		Product{Name: "a2", Price: 300, Category: "A"},
		Product{Name: "b1", Price: 100, Category: "B"},
	)

	// Delete where category = "A" AND price < 200; leaves a2 and b1.
	err := s.Delete(ctx,
		query.Filter("Category", "A"),
		query.Filter("Price", int64(200)).WithOP(query.LT),
	)
	require.NoError(t, err)

	n, err := s.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(2), n)
}
