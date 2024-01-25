// Package store defines the interface and implementation for data storage operations.
// It abstracts the CRUD (Create, Read, Update, Delete) functionalities and provides
// a generic way to interact with different types of data storage systems, be it SQL or NoSQL databases.
//
// This package is crucial for maintaining a clean separation of concerns between
// the data access layer and the business logic layer of an application. It allows for
// efficient data manipulation and retrieval with a variety of query parameters and
// supports advanced operations like conflict handling in upsert operations.
//
// Key components include:
// - Entity: A generic interface for models that can be stored.
// - Store: A generic interface for CRUD operations on entities.
// - OnConflict: A struct to define UPSERT operation behavior.
//
// The package aims to provide a robust and flexible way to handle data storage
// needs in a Go-based application, ensuring scalability and maintainability.
package store

import (
	"context"

	"github.com/jkaveri/goflexstore/query"
)

// Entity defines a generic interface for an Entity model in the context of data storage.
// This interface is designed to be implemented by any model that has a unique identifier.
// The ID type is a generic type constrained to types that are comparable.
type Entity[ID comparable] interface {
	// GetID returns the unique identifier of the entity. This method is intended
	// to provide a way to retrieve the unique identifier of an entity irrespective
	// of the specific type of the ID (e.g., int, string, UUID).
	GetID() ID
}

// OnConflict struct defines the behavior to be applied during an UPSERT operation (a combined INSERT and UPDATE operation).
// This struct is used to specify how conflicts should be handled when attempting to create a new entity that may already exist.
//
// Fields:
//   - Columns: A slice of strings specifying the column names that should be considered for determining a conflict.
//     If a conflict is detected based on these columns, the UPSERT operation decides to update the existing row
//     rather than creating a new row.
//   - UpdateAll: A boolean flag indicating whether all fields of the entity should be updated if a conflict is detected.
//     If set to true, all fields are updated with the values from the entity being upserted.
//   - DoNothing: A boolean flag that, when set to true, causes the UPSERT operation to take no action if a conflict is detected.
//     This is useful when you simply want to ignore the insert if the row already exists without performing any update.
//   - Updates: A map where keys are column names and values are the new values to be used in the update.
//     This map is used to specify custom updates on specific columns when a conflict is detected.
//   - UpdateColumns: A slice of strings specifying the column names that should be updated if a conflict is detected.
//     This field allows for partial updates, where only specified columns are updated.
//   - OnConstraint: A string specifying the name of the constraint that should be considered for detecting a conflict.
//     This is used in databases that support defining and naming constraints (e.g., unique constraints).
//
// The OnConflict struct is typically used with the Upsert method of a Store interface to define custom logic for handling
// insert/update operations where there may be a conflict with existing data.
type OnConflict struct {
	Columns       []string
	UpdateAll     bool
	DoNothing     bool
	Updates       map[string]any
	UpdateColumns []string
	OnConstraint  string
}

// Store defines a generic interface for CRUD (Create, Read, Update, Delete) operations
// on a specific type of Entity. This interface abstracts the data store operations allowing
// for implementation with different underlying data storage systems (e.g., SQL databases, NoSQL databases).
//
// T is a generic type that must satisfy the Entity interface, and ID is the type of the entity's identifier.
type Store[T Entity[ID], ID comparable] interface {
	// Get retrieves a single entity based on the provided query parameters.
	// Returns the found entity and nil if successful, nil and an error otherwise.
	Get(ctx context.Context, params ...query.Param) (T, error)

	// List retrieves a list of entities based on the provided query parameters.
	// Returns a slice of entities and nil if successful, nil and an error otherwise.
	List(ctx context.Context, params ...query.Param) ([]T, error)

	// Count returns the number of entities that match the provided query parameters.
	// Returns the count as int64 and nil if successful, 0 and an error otherwise.
	Count(ctx context.Context, params ...query.Param) (int64, error)

	// Exists checks if an entity exists based on the provided query parameters.
	// Returns true if at least one entity matches the parameters, false otherwise.
	Exists(ctx context.Context, params ...query.Param) (bool, error)

	// Create adds a new entity to the store.
	// Returns the ID of the newly created entity and nil if successful, zero-value ID and an error otherwise.
	Create(ctx context.Context, entity T) (ID, error)

	// Upsert creates a new entity or updates an existing one based on the conflict resolution strategy defined in OnConflict.
	// Returns the ID of the created/updated entity and nil if successful, zero-value ID and an error otherwise.
	Upsert(ctx context.Context, entity T, onConflict OnConflict) (ID, error)

	// CreateMany adds multiple entities to the store.
	// Returns nil if successful, an error otherwise.
	CreateMany(ctx context.Context, entities []T) error

	// Update modifies an existing entity based on the provided query parameters.
	// Returns nil if successful, an error otherwise.
	Update(ctx context.Context, entity T, params ...query.Param) error

	// PartialUpdate modifies parts of an existing entity based on the provided query parameters.
	// This method allows updating selective fields of an entity.
	// Returns nil if successful, an error otherwise.
	PartialUpdate(ctx context.Context, entity T, params ...query.Param) error

	// Delete removes an entity from the store based on the provided query parameters.
	// Returns nil if successful, an error otherwise.
	Delete(ctx context.Context, params ...query.Param) error
}
