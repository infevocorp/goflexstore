package converter

import "github.com/jkaveri/goflexstore/store"

// Converter is a converter that converts between DTO and Entity
type Converter[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] interface {
	// ToEntity converts DTO to Entity
	ToEntity(dto DTO) Entity
	// ToDTO converts Entity to DTO
	ToDTO(entity Entity) DTO
}

// ToMany converts a slice of A to a slice of B
func ToMany[A any, B any](items []A, convFn func(A) B) []B {
	var result []B
	for _, item := range items {
		result = append(result, convFn(item))
	}

	return result
}
