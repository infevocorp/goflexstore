// Package gormstore provides a GORM-based implementation of the store interface.
// This package leverages GORM, a popular ORM library for Golang, to offer a convenient
// and powerful way of performing CRUD operations and more on various data stores.
//
// The package defines a generic Store type that integrates with GORM's functionalities,
// offering an abstraction over GORM's native methods to align with the store.Entity interface.
// This allows for operations like Get, List, Create, Update, and Delete, to be performed
// in a type-safe and efficient manner.
//
// The Store type also supports advanced features such as transaction management through
// the TransactionScope, custom query building, batch operations, and conflict handling
// during upsert operations. It is designed to be flexible and extensible, making it
// suitable for a wide range of applications that require data persistence.
package gormstore
