# Flex Stores

Welcome to Flex Store, a revolutionary suite of interfaces tailored for Golang applications. Flex Store stands out with its:

- **Flexibility and Independence:** Our interfaces are crafted to offer unparalleled flexibility and independence. This design enables dynamic query handling and seamless integration with various ORMs and data sources, making your data layer more robust and versatile.

- **Enhanced Transaction Management:** The advanced Operation Scope interface of Flex Store brings a new level of power in managing transactions. It ensures data integrity and consistency across your applications, safeguarding your critical operations.

- **Adaptability and Growth:** Flex Store grows with you. As your project expands, it adapts, supporting additional features like metric operation scopes and tracing. This adaptability makes it an ideal partner for evolving application needs.

- **Simplified Data Layer:** We understand the challenges of data layer complexity. Flex Store significantly simplifies this, reducing the implementation and maintenance burdens. This lets developers focus more on developing business logic rather than getting bogged down by boilerplate code.

Whether you are working on a small-scale project or dealing with the intricacies of a large-scale application, Flex Store is equipped to enhance your data management efficiency and effectiveness. Embrace the ease of data layer management with Flex Store, and propel your Golang applications to new heights.

- [Flex Stores](#flex-stores)
  - [Getting started](#getting-started)
  - [Features](#features)
  - [Addressing the Limitations of the Repository Pattern with Flex Store](#addressing-the-limitations-of-the-repository-pattern-with-flex-store)
  - [Contribute](#contribute)

## Getting started

1. Get module

    ```bash
    go get github.com/jkaveri/goflexstore@latest
    ```

1. define models

    ```golang
    //file: model/user.go

    type User struct {
        ID int64
        Name string
    }

    func(u *User) GetID() int64 {
        return u.ID
    }

    ```

1. define store

    ```golang
    //file: store/user.go

    import (
        "github.com/jkaveri/goflexstore/store"

        "yourmodule/model"
    )

    type UserStore interface {
        store.Store[*model.User, int64]
    }
    ```

1. implement store

    1. get gorm implementation

        ```bash
        go get github.com/jkaveri/goflexstore/gorm@latest
        ```

    1. define dto

        ```golang
        //file: store/sql/dto/user.go
        type User struct {
            ID int64 `gorm:"column:id;primaryKey"`
            Name string `gorm:"name"`
        }

        func (u *User) GetID() int64 {
            return u.ID
        }
        ```

    1. implment user store

        ```golang
        //file: store/sql/user.go
        type UserStore struct {
            *gormstore.Store[*model.User, *dto.User, int64]
        }


        func NewUserStore(tScope *gormopscope.TransactionScope) *UserStore {
            return &UserStore{
                Store: gormstore.New(tScope),
            }
        }
        ```

## Features

- [x] **Query**
  - [x] Implement true flexible and independent Store interfaces.
- [x] **Operation Scope**
  - [x] Implement transaction management with Operation Scope interface.
  - [ ] Add metric operation scope.
  - [ ] Add tracing operation scope.
- [x] **Implementation**
  - [x] Implement GORM.
  - [ ] Consider integrating bun.
  - [ ] Explore other implementation options.
- [ ] Implement Cache store with automatic caching using a simple API like `query.WithCacheKey("abc")`.
- [ ] Develop flexstore-gen to generate models, stores, and DTOs from protobuf definitions.
- [ ] Create flexstore-graphql to provide a GraphQL API for querying data.

## Addressing the Limitations of the Repository Pattern with Flex Store

The Repository pattern is a common architectural practice, but it's not without its challenges, particularly when scaling or dealing with complex queries. The following Go code snippet demonstrates typical scenarios where the Repository pattern may lead to inefficiencies:

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

In this example, we observe two primary issues:

- **Query Duplication**: Both ListByTag and ListByTagAndAuthor methods exhibit a similar query pattern for selecting articles by tags, leading to redundancy.
- **Database Coupling**: The queries are tightly bound to the MySQL syntax, making it difficult to adapt or switch to other database systems without extensive refactoring.

As applications evolve and business logic grows, these problems tend to amplify, with increasing function count in the repository leading to more duplicated queries and stronger coupling to the database syntax. Even with an ORM like GORM, the dependency on specific database syntax persists, restricting flexibility and adaptability.

Recognizing these issues, Flex Store offers a more streamlined approach as demonstrated below:

1. Define `ArticleStore` interface

    ```golang
    //store/article.go

    // ArticleStore extends the flexstore/store.Store interface.
    type ArticleStore interface {
        store.Store[model.Article, int64]
    }
    ```

2. Define `filters`

    ```golang
    //filters/filters.go

    import (
        "github.com/jkaveri/goflexstore/query"
    )

    // Tag creates a query.FilterParam for tag filtering.
    func Tag(tag ...string) query.FilterParam {
        return query.Filter("tag", tag)
    }

    var GetTag = query.FilterGetter("tag")

    ```

3. Implement API handler

    ```golang
    // handlers/handlers.go

    type ApiHandler struct {
        store ArticleStore
    }

    type ListByTagRequest struct {
        Tag string
    }

    // ListByTag handles requests to list articles by tag.
    func (h *ApiHandler) ListByTag(ctx context.Context, request ListByTagRequest) ([]model.Article, error) {
        return h.store.List(ctx, filters.Tag(request.Tag))
    }

    type ListByAuthorRequest struct {
        Tag string
        AuthorID int64
    }

    // ListByAuthor handles requests to list articles by author and tag.
    func (h *ApiHandler) ListByAuthor(ctx context.Context, request ListByAuthorRequest) ([]*model.Article, error) {
        return h.store.List(ctx,
            filters.Tag(request.Tag),
            filters.Author(request.AuthorID),
        )
    }

    ```

In this improved approach, the code is more reusable, allowing us to compose different filter sets for varied query logic.

By adopting Flex Store's approach, developers can effectively circumvent the common pitfalls associated with the Repository pattern, leading to more maintainable, adaptable, and efficient codebases.

## Contribute

[Contribute.md](./CONTRIBUTE.md)
