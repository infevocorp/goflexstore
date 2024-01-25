package converter

import "github.com/jkaveri/goflexstore/store"

// NewManual creates a new Manual converter instance.
//
// This function allows the creation of a custom Converter by specifying
// the conversion functions directly. It is useful in cases where the conversion
// logic is not straightforward and requires custom implementation.
//
// Type parameters:
// - Entity: The type representing the Entity, typically used for database operations.
// - DTO: The type representing the Data Transfer Object, used for data transfer between layers or systems.
// - ID: The type of the identifier for the Entity and DTO, which must be comparable.
//
// Parameters:
// - toEntityFn: A function that converts a DTO to an Entity.
// - toDTOFn: A function that converts an Entity to a DTO.
//
// Returns:
// A Converter instance that uses the provided functions for conversion.
func NewManual[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable](
	toEntityFn func(dto DTO) Entity,
	toDTOFn func(entity Entity) DTO,
) Converter[Entity, DTO, ID] {
	return &Manual[Entity, DTO, ID]{
		ToEntityFn: toEntityFn,
		ToDTOFn:    toDTOFn,
	}
}

// Manual is a struct that implements the Converter interface using custom functions
// provided during its creation. This allows for flexible and custom conversion logic
// between DTOs and Entities.
//
// Type parameters:
// - Entity: The type representing the Entity.
// - DTO: The type representing the Data Transfer Object.
// - ID: The type of the identifier for the Entity and DTO.
//
// Fields:
// - ToEntityFn: A function that converts a DTO to an Entity.
// - ToDTOFn: A function that converts an Entity to a DTO.
type Manual[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] struct {
	ToEntityFn func(dto DTO) Entity
	ToDTOFn    func(entity Entity) DTO
}

// ToEntity calls the custom ToEntityFn function to convert a DTO to an Entity.
//
// This method utilizes the function provided during the creation of the Manual
// converter to transform a DTO into an Entity.
//
// Parameters:
// - dto: The DTO to convert.
//
// Returns:
// The converted Entity.
func (c *Manual[Entity, DTO, ID]) ToEntity(dto DTO) Entity {
	return c.ToEntityFn(dto)
}

// ToDTO calls the custom ToDTOFn function to convert an Entity to a DTO.
//
// This method utilizes the function provided during the creation of the Manual
// converter to transform an Entity into a DTO.
//
// Parameters:
// - entity: The Entity to convert.
//
// Returns:
// The converted DTO.
func (c *Manual[Entity, DTO, ID]) ToDTO(entity Entity) DTO {
	return c.ToDTOFn(entity)
}
