package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	gormopscope "github.com/jkaveri/goflexstore/gorm/opscope"
	flexstore "github.com/jkaveri/goflexstore/store"

	"github.com/jkaveri/goflexstore/examples/cms/handlers"
	"github.com/jkaveri/goflexstore/examples/cms/model"
	"github.com/jkaveri/goflexstore/examples/cms/store"
	storesql "github.com/jkaveri/goflexstore/examples/cms/store/sql"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	stores := newStores(ctx)

	// new echo instance
	e := echo.New()

	// register handlers
	handlers.Register(stores, e)

	// Initialize the server in a goroutine so that it doesn't block.
	go func() {
		if err := e.StartServer(&http.Server{
			Addr: ":8080",
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
			ReadTimeout:  time.Duration(5) * time.Second,
			WriteTimeout: time.Duration(5) * time.Second,
		}); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("server error: %v", err)
		}
	}()

	// Block until we receive our signal.
	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	log.Println("Shutting down server...")
	if err := e.Shutdown(shutDownCtx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %+v", err)
	}

	log.Println("Server shut down gracefully")
}

func newStores(ctx context.Context) store.Stores {
	// open db
	db, err := gorm.Open(sqlite.Open("cms.db"), &gorm.Config{})
	panicIfErr(err)

	// run migrations
	err = storesql.AutoMigrate(db)
	panicIfErr(err)

	// create scope
	scope := gormopscope.NewWriteTransactionScope("write", db)

	// create stores
	stores := storesql.NewStores(scope)

	// seed test data
	seedData(ctx, stores)

	return stores
}

func seedData(ctx context.Context, stores store.Stores) {
	_, err := stores.User.Upsert(ctx, &model.User{
		ID:    1,
		Name:  "John Doe",
		Email: "jonh@email.com",
	}, flexstore.OnConflict{
		Columns:   []string{"id"},
		DoNothing: true,
	})
	panicIfErr(err)

	_, err = stores.Article.Upsert(ctx, &model.Article{
		ID:       1,
		Title:    "Article 1",
		Content:  "Content 1",
		AuthorID: 1,
		Tags: []*model.Tag{
			{
				ID:   1,
				Slug: "tag-1",
			},
		},
	}, flexstore.OnConflict{
		Columns:   []string{"id"},
		UpdateAll: true,
	})
	panicIfErr(err)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
