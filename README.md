# Flex Store

Set of interfaces to implement your data layer more flexibly in Golang.

## Features

- [x] Query
  - [x] True flexible and indepent Store interfaces
- [x] Operation Scope
  - [x] Transaction management with Operation Scope interface
  - [ ] Metric operation scope
  - [ ] Tracing operation scope
- [x] Implementation
  - [x] GORM
  - [ ] bun?
  - [ ] ..?
- [ ] Cache store, auto cache with simple api like `query.WithCacheKey("abc")`

## Store

A Store is an interface that helps you abstract your data layer API, which allows the code to be easier to test and maintain.

If you are familiar with the Repository pattern, then you'll see the differences between a Store and a Repository. A Store is purely a data layer query which is another layer below the Repository. In other words, you can use the interfaces of the Store to implement your Repository with additional business logic. For example,

P.S.: Please do not take the code sample below as a getting started guide. It only has one purpose: to serve as an example. I will have another section to show how to get started later.

```golang
    type UserStore interface {
        store.Store[User, int64]
    }

    type UserRepo interface {
        GetByUsername(ctx context.Context, username string) (User, error)
    }

    type UserRepoImpl struct {
        Store UserStore
    }
    func (r *UserRepoImpl) GetByUsername(ctx context.Context, username string) (User, error) {
        return r.Store.Get(ctx, filters.Username(username))
    }
```

However, in my experience, if you already have a Store, you don't need to have a Repository layer. You can go straight to the application logic layer. This helps your application be easier to maintain.

**Why do we need to have a Store?**

I usually use GORM as an ORM, which already contains a flexible interface to communicate with databases. However, I always structure with reusable queries. For example:

When we used to implement our Repository layer like this

```golang
type ArticleRepo struct {
    db *gorm.DB
}

func (r *ArticleRepo) ListByTag(ctx context.Context, tag string) ([]model.Article, error) {
    var articles []dto.Article

    if err := r.db.
        Where("id IN (SELECT article_id FROM article_tags WHERE tag = ?)", tag).
        Find(&articles).Error; err != nil {
            return nil, err
    }

    return mapDaoToArticles(articles), nil
}

func (r *ArticleRepo) ListByTagAndAuthor(ctx context.Context, authorId int64, tag string) ([]model.Article, error) {
     var articles []dto.Article

    if err := r.db.
        Where("author_id = ?", authorId).
        Where("id IN (SELECT article_id FROM article_tags WHERE tag = ?)", tag).
        Find(&articles).Error; err != nil {
            return nil, err
    }

    return mapDaoToArticles(articles), nil
}
```

In the code above, you can see two problems:

1. The query is duplicated when selected by the tag.
2. The query is bound to the low-level infrastructure (MySQL syntax).

Over time, as the business logic grows, there are more and more functions in the repository. It means more duplicated queries and more binding to MySQL, even if you are using GORM, but it doesn't mean you can easily switch to another database.

**Now, let's see how Flex Store helps you improve the code:**

1. Define `ArticleStore` interface

    ```golang
    //store/article.go

    type ArticleStore interface {
        // extends flexstore/store.Store interface
        store.Store[model.Article, int64]
    }

    ```

1. Define `filters`

    ```golang
    //filters/filters.go

    import (
        "github.com/jkaveri/goflexstore/query"
    )

    func Tag(tag ...string) query.FilterParam {
        return query.Filter("tag", tag)
    }

    var GetTag = query.FilterGetter("tag")

    ```

1. Implement API handler

    ```golang
    // handlers/handlers.go

    type ApiHandler struct {
        store ArticleStore
    }

    type ListByTagRequest struct {
        Tag string
    }

    func (h *ApiHandler) ListByTag(ctx context.Context, request ListByTagRequest) ([]model.Article, error) {
        return h.store.List(ctx, filters.Tag(request.Tag))
    }

    type ListByAuthorRequest struct {
        Tag string
        AuthorID int64
    }

    func (h *ApiHandler) ListByAuthor(ctx context.Context, request ListByAuthorRequest) ([]*model.Article, error) {
        return h.store.List(ctx,
            filters.Tag(request.Tag),
            filters.Author(request.AuthorID),
        )
    }

    ```

You can see the code is now reusable as we can compose different filter sets to have different query logic.
