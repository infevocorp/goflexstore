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

	"github.com/infevocorp/goflexstore/query"
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

// OnConflict struct defines the behavior to be applied during an UPSERT operation
// (a combined INSERT and UPDATE operation). This struct is used to specify how
// conflicts should be handled when attempting to create a new entity that may already exist.
//
// Fields:
//   - Columns: A slice of strings specifying the column names that should be considered for determining a conflict.
//     If a conflict is detected based on these columns, the UPSERT operation decides to update the existing row
//     rather than creating a new row.
//   - UpdateAll: A boolean flag indicating whether all fields of the entity should be updated if a conflict is
//     detected.
//     If set to true, all fields are updated with the values from the entity being upserted.
//   - DoNothing: A boolean flag that, when set to true, causes the UPSERT operation to take no action if a conflict
//     is detected.
//     This is useful when you simply want to ignore the insert if the row already exists without performing any update.
//   - Updates: A map where keys are column names and values are the new values to be used in the update.
//     This map is used to specify custom updates on specific columns when a conflict is detected.
//   - UpdateColumns: A slice of strings specifying the column names that should be updated if a conflict is detected.
//     This field allows for partial updates, where only specified columns are updated.
//   - OnConstraint: A string specifying the name of the constraint that should be considered for detecting a conflict.
//     This is used in databases that support defining and naming constraints (e.g., unique constraints).
//
// The OnConflict struct is typically used with the Upsert method of a Store interface to define custom logic for
// handling insert/update operations where there may be a conflict with existing data.
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
	//
	// This method attempts to find and return an entity that matches the criteria specified by the query parameters.
	// If successful, the found entity and a nil error are returned. If no entity is found or an error occurs, nil
	// and the error are returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - params: A variable number of query.Param, each representing a filter condition for the query.
	//
	// Returns: The found entity of type T if successful, nil and an error otherwise.
	//
	// Example:
	// Retrieving an entity by its ID:
	//
	//	entity, err := store.Get(ctx, query.Filter("id", entityID))
	//
	// Note: If no entity matches the query parameters, an error indicating "not found" is typically returned.
	Get(ctx context.Context, params ...query.Param) (T, error)

	// List retrieves a list of entities based on the provided query parameters.
	//
	// This method fetches and returns a slice of entities that match the criteria specified by the query parameters.
	// If successful, the slice of entities and a nil error are returned. If an error occurs during the retrieval,
	// nil and the error are returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - params: A variable number of query.Param, each representing a filter condition for the query.
	//
	// Returns: A slice of entities of type T if successful, nil and an error otherwise.
	//
	// Example:
	// Listing entities with a specific attribute value:
	//
	//	entities, err := store.List(ctx, query.Filter("attribute", value))
	List(ctx context.Context, params ...query.Param) ([]T, error)

	// Count returns the number of entities that match the provided query parameters.
	//
	// This method counts and returns the number of entities that satisfy the criteria specified by the
	// query parameters. If successful, the count as int64 and a nil error are returned. If an error occurs,
	// 0 and the error are returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - params: A variable number of query.Param, each representing a filter condition for the query.
	//
	// Returns: The count of matching entities as int64 if successful, 0 and an error otherwise.
	//
	// Example:
	// Counting entities with a specific condition:
	//
	//	count, err := store.Count(ctx, query.Filter("status", "active"))
	Count(ctx context.Context, params ...query.Param) (int64, error)

	// Exists checks if at least one entity exists based on the provided query parameters.
	//
	// This method determines the existence of any entity that matches the criteria specified by the query parameters.
	// It returns true if at least one entity matches, false otherwise. In case of an error, false and the error are
	// returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - params: A variable number of query.Param, each representing a filter condition for the query.
	//
	// Returns: True if at least one matching entity exists, false and an error otherwise.
	//
	// Example:
	// Checking the existence of an entity with a specific ID:
	//
	//	exists, err := store.Exists(ctx, query.Filter("id", entityID))
	Exists(ctx context.Context, params ...query.Param) (bool, error)

	// Create adds a new entity to the store.
	//
	// This method inserts a new entity into the store and returns the ID of the newly created entity if successful.
	// If an error occurs during the creation, the zero-value of ID and the error are returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - entity: The entity of type T to be added to the store.
	//
	// Returns: The ID of the newly created entity if successful, zero-value of ID and an error otherwise.
	//
	// Example:
	// Adding a new entity to the store:
	//
	//	newID, err := store.Create(ctx, newEntity)
	Create(ctx context.Context, entity T) (ID, error)

	// Upsert creates a new entity or updates an existing one based on the conflict resolution strategy defined in
	// OnConflict.
	//
	// This method either inserts a new entity or updates an existing one in the store, depending on the presence of a
	// conflict based on the specified OnConflict strategy. The OnConflict struct allows specifying how conflicts,
	// identified by duplicate
	// values in specified columns or constraints, should be resolved. Options include updating all fields, performing
	// no action, or updating specific fields of the existing entity.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - entity: The entity of type T to be created or updated in the store.
	//   - onConflict: The conflict resolution strategy, encapsulated within an OnConflict struct, defining how to
	//	 handle conflicts. The OnConflict struct includes options to specify conflict-determining columns, whether
	// 	to update all fields or just specified one, and whether to ignore the operation if a conflict is detected.
	//
	// Returns: The ID of the created or updated entity if successful, zero-value of ID and an error otherwise.
	//
	// Example:
	// Upserting an entity with conflict resolution:
	//
	//	upsertedID, err := store.Upsert(ctx, entity, OnConflict{
	//	  Columns:       []string{"column_name"}, // Columns to check for conflict
	//	  UpdateAll:     true,                    // Update all fields if conflict exists
	//	})
	//
	// Note: The OnConflict struct allows for flexible conflict resolution strategies, including updating all fields,
	// no action, custom updates, partial updates, or based on specific constraints.
	Upsert(ctx context.Context, entity T, onConflict OnConflict) (ID, error)

	// CreateMany adds multiple entities to the store in a single operation.
	//
	// This method inserts a batch of entities into the store. It returns nil if the operation is successful.
	// If an error occurs during the batch insertion, the error is returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - entities: A slice of entities of type T to be added to the store.
	//
	// Returns: Nil if successful, an error otherwise.
	//
	// Example:
	// Adding multiple entities to the store at once:
	//
	//	err := store.CreateMany(ctx, entities)
	CreateMany(ctx context.Context, entities []T) error

	// Update modifies an existing entity based on the provided query parameters or the entity's ID field.
	//
	// This method updates an entity in the store that matches the criteria specified by the query parameters. If no
	// query parameters are provided, the method uses the ID field of the entity to locate the record to be updated.
	// It returns nil if the update operation is successful. If an error occurs during the update, the error is
	// returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - entity: The modified entity of type T with updated values. The entity's ID field is used for lookup if no
	//     query parameters are provided.
	//   - params: An optional variable number of query.Param, each representing a filter condition to identify the
	//     entity to be updated. If no parameters are provided, the entity's ID field is used as the lookup criterion.
	//
	// Returns: Nil if successful, an error otherwise.
	//
	// Example:
	// Updating an existing entity in the store using query parameters:
	//
	//	err := store.Update(ctx, updatedEntity, query.Filter("id", entityID))
	//
	// Example:
	// Updating an existing entity in the store using the entity's ID field (no query parameters provided):
	//
	//	err := store.Update(ctx, updatedEntity)
	//
	// Note: Providing specific query parameters allows for more granular control over the update operation, while
	// omitting them defaults to using the entity's ID field for identification. This approach provides flexibility
	// in how entities are located for updates.
	Update(ctx context.Context, entity T, params ...query.Param) error

	// PartialUpdate modifies parts of an existing entity based on the provided query parameters or the entity's ID
	// field.
	//
	// This method allows for selective updating of fields of an existing entity in the store. Only the specified
	// fields of the entity are updated, either based on the criteria specified by the query parameters or by using
	// the entity's ID field if no parameters are provided. This method returns nil if the partial update operation is
	// successful. If an error occurs, the error is returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - entity: The entity of type T with the fields to be updated specified. If no query parameters are provided,
	//     the entity's ID field is used for lookup.
	//   - params: An optional variable number of query.Param, each representing a filter condition to identify the
	//     entity to be partially updated. If no parameters are provided, the entity's ID field is used as the
	//     lookup criterion.
	//
	// Returns: Nil if successful, an error otherwise.
	//
	// Example:
	// Partially updating an entity's specific fields using query parameters:
	//
	//	err := store.PartialUpdate(ctx, partialEntity, query.Filter("id", entityID))
	//
	// Example:
	// Partially updating an entity's specific fields using the entity's ID field (no query parameters provided):
	//
	//	err := store.PartialUpdate(ctx, partialEntity)
	//
	// Note: This method offers the flexibility to update selective fields of an entity, enhancing the efficiency of
	// data manipulation. Providing specific query parameters allows for more precise targeting of the entity to be
	//  updated, while omitting them defaults to using the entity's ID for identification.
	PartialUpdate(ctx context.Context, entity T, params ...query.Param) error

	// Delete removes an entity from the store based on the provided query parameters.
	//
	// This method deletes an existing entity from the store that matches the criteria specified by the query
	// parameters. It returns nil if the deletion is successful. If an error occurs during the deletion, the error
	// is returned.
	//
	// Parameters:
	//   - ctx: A context.Context to control the request's deadline and cancellation.
	//   - params: A variable number of query.Param, each representing a filter condition to identify the entity to
	//     be deleted.
	//
	// Returns: Nil if successful, an error otherwise.
	//
	// Example:
	// Removing an entity from the store:
	//
	//	err := store.Delete(ctx, query.Filter("id", entityID))
	Delete(ctx context.Context, params ...query.Param) error
}
