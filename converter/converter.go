package converter

import "github.com/jkaveri/goflexstore/store"

// Converter is an interface that defines methods for converting between a DTO (Data Transfer Object)
// and an Entity. It is a generic interface, allowing for flexible implementation for various types.
//
// Type parameters:
//   - Entity: The type representing the Entity, typically used for database operations.
//   - DTO: The type representing the Data Transfer Object, used for data transfer between layers or systems.
//   - ID: The type of the identifier for the Entity and DTO, which must be comparable.
//
// The Converter interface is particularly useful in scenarios where data needs to be transformed
// between different layers of an application, such as from the persistence layer to the domain layer.
type Converter[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] interface {
	// ToEntity converts a DTO into an Entity. This method is used when data is received, for example,
	// from an API call, and needs to be transformed into an Entity for storage or processing.
	ToEntity(dto DTO) Entity

	// ToDTO converts an Entity into a DTO. This method is used when data from the application's
	// internal representation (Entity) needs to be transformed for external use, such as sending
	// data to a client via an API.
	ToDTO(entity Entity) DTO
}

// ToMany is a utility function that converts a slice of one type (A) to a slice of another type (B)
// using a provided conversion function.
//
// Type parameters:
//   - A: The original type of the items in the slice.
//   - B: The target type of the items in the slice.
//
// Parameters:
//   - items: A slice of type A that needs to be converted.
//   - convFn: A function that takes an item of type A and returns its equivalent in type B.
//
// Returns:
// A slice of type B with each item converted from type A using the provided conversion function.
//
// This function is useful in situations where you have a collection of items of one type that
// need to be transformed into another type, such as converting a slice of database entities
// into a slice of DTOs for API responses.
func ToMany[A any, B any](items []A, convFn func(A) B) []B {
	var result []B
	for _, item := range items {
		result = append(result, convFn(item))
	}

	return result
}
