// Package main demonstrates using goflexstore's sqlx store with a real PostgreSQL database.
//
// Usage:
//
//	docker-compose up -d
//	go run .
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	"github.com/infevocorp/goflexstore/query"
	"github.com/infevocorp/goflexstore/store"
	sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
	sqlxquery "github.com/infevocorp/goflexstore/sqlx/query"
	sqlxstore "github.com/infevocorp/goflexstore/sqlx/store"
)

// Product is the clean domain entity — no db tags.
type Product struct {
	ID       int64
	Name     string
	Price    float64
	Stock    int
	Category string
}

func (p Product) GetID() int64 { return p.ID }

// ProductRow is the DB scan target (DTO) with db struct tags for sqlx.
type ProductRow struct {
	ID       int64   `db:"id"`
	Name     string  `db:"name"`
	Price    float64 `db:"price"`
	Stock    int     `db:"stock"`
	Category string  `db:"category"`
}

func (r ProductRow) GetID() int64 { return r.ID }

const createTable = `
CREATE TABLE IF NOT EXISTS products (
	id       BIGSERIAL     PRIMARY KEY,
	name     TEXT          NOT NULL,
	price    NUMERIC(10,2) NOT NULL DEFAULT 0,
	stock    INTEGER       NOT NULL DEFAULT 0,
	category TEXT          NOT NULL DEFAULT ''
)`

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5434/goflexstore?sslmode=disable"
	}

	db, err := sqlx.Open("postgres", dsn)
	must(err, "open db")
	defer db.Close()

	must(db.Ping(), "ping db")
	fmt.Println("Connected to PostgreSQL")

	_, err = db.Exec(createTable)
	must(err, "create table")

	_, err = db.Exec("DELETE FROM products")
	must(err, "cleanup")

	// Build the store: Postgres dialect + RETURNING id for auto-increment PK.
	opScope := sqlxopscope.NewTransactionScope("main", db, nil)
	products := sqlxstore.New[Product, ProductRow, int64](
		opScope,
		sqlxstore.WithTable[Product, ProductRow, int64]("products"),
		sqlxstore.WithDialect[Product, ProductRow, int64](sqlxquery.DialectPostgres),
		sqlxstore.WithReturningID[Product, ProductRow, int64](true),
	)

	ctx := context.Background()

	// ── Create ────────────────────────────────────────────────────────────────
	section("Create")

	widgetID, err := products.Create(ctx, Product{Name: "Widget", Price: 9.99, Stock: 100, Category: "tools"})
	must(err, "create Widget")
	fmt.Printf("Created Widget       id=%d\n", widgetID)

	gadgetID, err := products.Create(ctx, Product{Name: "Gadget", Price: 24.99, Stock: 50, Category: "electronics"})
	must(err, "create Gadget")
	fmt.Printf("Created Gadget       id=%d\n", gadgetID)

	_, err = products.Create(ctx, Product{Name: "Doohickey", Price: 4.99, Stock: 200, Category: "tools"})
	must(err, "create Doohickey")
	fmt.Println("Created Doohickey")

	// ── Get ───────────────────────────────────────────────────────────────────
	section("Get")

	widget, err := products.Get(ctx, query.Filter("ID", widgetID))
	must(err, "get Widget")
	fmt.Printf("Got: id=%-3d %-15s $%.2f  stock=%d  category=%s\n",
		widget.ID, widget.Name, widget.Price, widget.Stock, widget.Category)

	// ── List ──────────────────────────────────────────────────────────────────
	section("List all")

	all, err := products.List(ctx)
	must(err, "list all")
	printProducts(all)

	section("List with filter (category=tools)")

	tools, err := products.List(ctx, query.Filter("Category", "tools"))
	must(err, "list tools")
	printProducts(tools)

	section("List ordered by price ASC with pagination")

	page, err := products.List(ctx,
		query.OrderBy("Price", false),
		query.Paginate(0, 10),
	)
	must(err, "list ordered")
	printProducts(page)

	// ── Count / Exists ────────────────────────────────────────────────────────
	section("Count / Exists")

	total, _ := products.Count(ctx)
	toolCount, _ := products.Count(ctx, query.Filter("Category", "tools"))
	fmt.Printf("Total=%d  tools=%d\n", total, toolCount)

	exists, _ := products.Exists(ctx, query.Filter("Name", "Widget"))
	fmt.Printf("Widget exists: %v\n", exists)

	// ── Update ────────────────────────────────────────────────────────────────
	section("Update")

	err = products.Update(ctx, Product{ID: widgetID, Name: "Widget Pro", Price: 14.99, Stock: 90, Category: "tools"})
	must(err, "update")
	updated, _ := products.Get(ctx, query.Filter("ID", widgetID))
	fmt.Printf("After update:  id=%-3d %-15s $%.2f  stock=%d\n",
		updated.ID, updated.Name, updated.Price, updated.Stock)

	// ── PartialUpdate ─────────────────────────────────────────────────────────
	section("PartialUpdate (price only — zero fields are skipped)")

	err = products.PartialUpdate(ctx, Product{ID: widgetID, Price: 19.99})
	must(err, "partial update")
	partial, _ := products.Get(ctx, query.Filter("ID", widgetID))
	fmt.Printf("After partial: name=%q price=%.2f stock=%d\n",
		partial.Name, partial.Price, partial.Stock)

	// ── CreateMany ────────────────────────────────────────────────────────────
	section("CreateMany")

	err = products.CreateMany(ctx, []Product{
		{Name: "Thingamajig", Price: 1.99, Stock: 500, Category: "misc"},
		{Name: "Whatchamacallit", Price: 2.99, Stock: 300, Category: "misc"},
	})
	must(err, "create many")
	total, _ = products.Count(ctx)
	fmt.Printf("Total after CreateMany: %d\n", total)

	// ── Upsert ────────────────────────────────────────────────────────────────
	section("Upsert (ON CONFLICT (id) DO UPDATE SET ...)")

	upsertID, err := products.Upsert(ctx,
		Product{ID: widgetID, Name: "Widget Pro Max", Price: 39.99, Stock: 60, Category: "tools"},
		store.OnConflict{
			Columns:   []string{"id"},
			UpdateAll: true,
		},
	)
	must(err, "upsert")
	upserted, _ := products.Get(ctx, query.Filter("ID", upsertID))
	fmt.Printf("After upsert:  id=%-3d %-15s $%.2f  stock=%d\n",
		upserted.ID, upserted.Name, upserted.Price, upserted.Stock)

	_ = gadgetID // used above

	// ── Transaction: commit ───────────────────────────────────────────────────
	section("Transaction — commit")

	txCtx, err := opScope.Begin(ctx)
	must(err, "begin tx")

	txID, err := products.Create(txCtx, Product{Name: "Premium Widget", Price: 99.99, Stock: 10, Category: "premium"})
	must(err, "create in tx")
	fmt.Printf("Created inside tx: id=%d\n", txID)

	must(opScope.End(txCtx, nil), "commit tx")
	total, _ = products.Count(ctx)
	fmt.Printf("Count after commit: %d\n", total)

	// ── Transaction: rollback ─────────────────────────────────────────────────
	section("Transaction — rollback")

	txCtx2, _ := opScope.Begin(ctx)
	_, _ = products.Create(txCtx2, Product{Name: "Ephemeral", Price: 0, Stock: 0, Category: "none"})
	_ = opScope.End(txCtx2, fmt.Errorf("intentional rollback"))
	total, _ = products.Count(ctx)
	fmt.Printf("Count after rollback (unchanged): %d\n", total)

	// ── Delete ────────────────────────────────────────────────────────────────
	section("Delete (category=misc)")

	err = products.Delete(ctx, query.Filter("Category", "misc"))
	must(err, "delete")
	total, _ = products.Count(ctx)
	fmt.Printf("Count after delete: %d\n", total)

	fmt.Println("\nAll done!")
}

func section(name string) {
	fmt.Printf("\n─── %s\n", name)
}

func printProducts(ps []Product) {
	for _, p := range ps {
		fmt.Printf("  id=%-3d %-18s $%6.2f  stock=%-4d  category=%s\n",
			p.ID, p.Name, p.Price, p.Stock, p.Category)
	}
}

func must(err error, op string) {
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}
}
