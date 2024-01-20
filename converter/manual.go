package converter

import "github.com/jkaveri/goflexstore/store"

func NewManul[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable](
	toEntityFn func(dto DTO) Entity,
	toDTO func(entity Entity) DTO,
) Converter[Entity, DTO, ID] {
	return &Manual[Entity, DTO, ID]{
		ToEntityFn: toEntityFn,
		ToDTOFn:    toDTO,
	}
}

type Manual[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] struct {
	ToEntityFn func(dto DTO) Entity
	ToDTOFn    func(entity Entity) DTO
}

func (c *Manual[Entity, DTO, ID]) ToEntity(dto DTO) Entity {
	return c.ToEntityFn(dto)
}

func (c *Manual[Entity, DTO, ID]) ToDTO(dto Entity) DTO {
	return c.ToDTOFn(dto)
}
