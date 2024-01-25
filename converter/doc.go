// Package converter provides a set of tools for converting between Data Transfer Objects (DTOs)
// and Entities in an application. This package is essential in scenarios where data
// needs to be transformed between different layers, such as from the persistence layer
// (e.g., database models) to the presentation layer (e.g., API models). It supports both
// manual and reflection-based conversion methods, offering flexibility to handle various
// data transformation needs.
// Example of creating a Manual converter:
// This example demonstrates how to instantiate a Manual converter. It requires
// providing custom functions for converting between the DTO and Entity types.
//
// Example:
//
//	toEntityFn := func(dto MyDTO) MyEntity {
//	    // Custom logic to convert DTO to Entity
//	}
//	toDTOFn := func(entity MyEntity) MyDTO {
//	    // Custom logic to convert Entity to DTO
//	}
//	manualConverter := converter.NewManual(toEntityFn, toDTOFn)
//
// Here, `MyEntity` and `MyDTO` are custom types representing the Entity and DTO respectively.
// Example of creating a Reflect converter:
// The Reflect converter uses reflection to automatically map fields between the DTO
// and Entity. It can be customized with a field mapping to handle cases where field names differ.
//
// Example:
//
//	fieldMapping := map[string]string{
//	    "EntityFieldName": "DTOFieldName",
//	}
//	reflectConverter := converter.NewReflect[MyEntity, MyDTO](fieldMapping)
//
// In this example, `fieldMapping` is used to define custom mappings between field names
// in the Entity (`MyEntity`) and the DTO (`MyDTO`). If field names are the same, they are
// automatically mapped without needing to be specified in `fieldMapping`.
package converter
