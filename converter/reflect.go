package converter

import (
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/pkg/errors"

	"github.com/jkaveri/goflexstore/store"
)

// NewReflect creates a new reflection-based converter.
//
// It converts between DTO and Entity using reflection, mapping fields from one to the other.
// The `overridesMapping` argument allows specifying custom field name mappings between the Entity and DTO.
// If nil or empty, the Entity's field names are used as DTO's field names.
//
// Type parameters:
//   - Entity: The Entity type implementing store.Entity interface.
//   - DTO: The DTO type implementing store.Entity interface.
//   - ID: The type of the identifier for Entity and DTO, which must be comparable.
//
// Parameters:
//   - overridesMapping: A map where the key is the Entity's field name and the value is the DTO's field name.
//
// Returns:
// A new instance of Reflect converter with the specified field mappings.
func NewReflect[
	Entity store.Entity[ID],
	DTO store.Entity[ID],
	ID comparable,
](
	overridesMapping map[string]string,
) Converter[Entity, DTO, ID] {
	return Reflect[Entity, DTO, ID]{
		dtoFieldsMapping:   overridesMapping,
		entityFieldMapping: reverseMapping(overridesMapping),
	}
}

// Reflect is a converter that uses reflection to convert between DTO and Entity.
// It implements the Converter interface and allows for automated conversion based on field names.
//
// Type parameters:
//   - Entity: The Entity type.
//   - DTO: The DTO type.
//   - ID: The type of the identifier for Entity and DTO.
//
// Fields:
//   - dtoFieldsMapping: Map where the key is Entity's field name and the value is DTO's field name.
//   - entityFieldMapping: Map where the key is DTO's field name and the value is Entity's field name.
type Reflect[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] struct {
	// fieldMapping key is Entity's field name. value is DTO's field name.
	dtoFieldsMapping map[string]string
	// fieldMapping key is DTO's field name. value is Entity's field name.
	entityFieldMapping map[string]string
}

// ToEntity converts a DTO to an Entity using reflection.
// It creates a new instance of Entity and copies values from the DTO to the Entity based on field mappings.
//
// Parameters:
//   - dto: The DTO to be converted to Entity.
//
// Returns:
// The converted Entity.
func (c Reflect[Entity, DTO, ID]) ToEntity(dto DTO) Entity {
	entity := *new(Entity)

	reflectCopy(dto, &entity, c.entityFieldMapping)

	return entity
}

// ToDTO converts an Entity to a DTO using reflection.
// It creates a new instance of DTO and copies values from the Entity to the DTO based on field mappings.
//
// Parameters:
//   - entity: The Entity to be converted to DTO.
//
// Returns:
// The converted DTO.
func (c Reflect[Entity, DTO, ID]) ToDTO(entity Entity) DTO {
	dto := *new(DTO)

	reflectCopy(entity, &dto, c.dtoFieldsMapping)

	return dto
}

// reflectCopy performs the actual copying of values from the source to the destination.
// It iterates over the fields of the destination and sets values from the source based on the provided field mapping.
//
// Parameters:
//   - src: The source object.
//   - dst: The destination object.
//   - fieldMapping: Map where the key is the destination field name and the value is the source field name.
func reflectCopy(src any, dst any, fieldMapping map[string]string) {
	// Obtain a reflection Value of the source object.
	srcVal := reflect.ValueOf(src)

	// Unwrap the source value if it's a pointer.
	// This is to handle cases where the source is a pointer type.
	for srcVal.Kind() == reflect.Ptr {
		// If the source value is a zero value (nil), return early.
		if srcVal.IsZero() {
			return
		}

		// Get the actual value that the pointer points to.
		srcVal = srcVal.Elem()
	}

	// Obtain a reflection Value of the destination object.
	dstVal := reflect.ValueOf(dst)
	// Ensure the destination is a pointer, as we need to modify it.
	if dstVal.Kind() != reflect.Ptr {
		panic("dst must be reference type (pointer)")
	}

	// Unwrap the destination value if it's a pointer.
	// This also ensures that we're dealing with the actual value.
	for dstVal.Kind() == reflect.Ptr {
		// If the destination is a nil pointer, initialize it with a new value.
		if dstVal.IsNil() {
			dstVal.Set(reflect.New(dstVal.Type().Elem()))
		}

		// Get the actual value that the pointer points to.
		dstVal = dstVal.Elem()
	}

	// Get the type information of the destination.
	dstType := dstVal.Type()

	// Iterate over each field of the destination.
	numField := dstVal.NumField()
	for i := 0; i < numField; i++ {
		// Get the i-th field of the destination.
		dstField := dstVal.Field(i)

		// Skip if the field cannot be set (unexported private field).
		if !dstField.CanSet() {
			continue
		}

		// Get the name of the i-th field.
		dstFieldName := dstType.Field(i).Name

		// If a field mapping exists, use it to find the corresponding source field.
		if fieldMapping != nil {
			if f, ok := fieldMapping[dstFieldName]; ok && f != "" {
				dstFieldName = f
			}
		}

		// Find the field in the source object that matches the destination field.
		srcField := srcVal.FieldByName(dstFieldName)
		// Skip if the source field is not valid (doesn't exist).
		if !srcField.IsValid() {
			continue
		}

		// If the source field is a pointer but nil, skip copying.
		if (srcField.Kind() == reflect.Ptr || srcField.Kind() == reflect.Slice) && srcField.IsNil() {
			continue
		}

		// Attempt to set the destination field with the value of the source field.
		// Panic with a detailed error message if the assignment is not possible.
		if !setValue(srcField, dstField) {
			panic(errors.Errorf(
				"cannot assign src.%s(%s) to dst.%s(%s)",
				dstFieldName,
				srcField.Type().String(),
				dstFieldName,
				dstField.Type().String(),
			))
		}
	}
}

func reverseMapping[K comparable, V comparable](m map[K]V) map[V]K {
	reversed := make(map[V]K, len(m))
	for k, v := range m {
		reversed[v] = k
	}

	return reversed
}

func setValue(srcVal, dstVal reflect.Value) bool {
	// same type
	if srcVal.Type() == dstVal.Type() {
		dstVal.Set(srcVal)
		return true
	}

	if ok := tryIfTargetTypeIsScanner(srcVal, dstVal); ok {
		return true
	}

	if ok := tryIfTargetTypeIsValuer(srcVal, dstVal); ok {
		return true
	}

	if ok := tryIfStruct(srcVal, dstVal); ok {
		return true
	}

	if ok := tryIfSlice(srcVal, dstVal); ok {
		return true
	}

	return false
}

func tryIfTargetTypeIsScanner(src reflect.Value, dst reflect.Value) bool {
	// check if dst is struct so we should use pointer to dst
	// because all sql.Null* types implement sql.Scanner as pointer receiver
	if dst.Kind() == reflect.Struct && dst.CanAddr() {
		dst = dst.Addr()
	}

	// check if dst implements sql.Scanner interface
	if !dst.Type().Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) {
		return false
	}

	// check if dst is nil
	// so need to init it first
	if dst.Kind() == reflect.Ptr && dst.IsNil() {
		dst.Set(reflect.New(dst.Type().Elem()))
	}

	if results := dst.MethodByName("Scan").Call([]reflect.Value{src}); !results[0].IsNil() {
		err := results[0].Interface().(error)
		panic(errors.Errorf("cannot assign %s to %s: %v", src.String(), dst.String(), err))
	}

	return true
}

func tryIfTargetTypeIsValuer(src reflect.Value, dst reflect.Value) bool {
	// unwrap pointer
	// because all driver.Valuer types implement driver.Valuer as value receiver
	for src.Kind() == reflect.Ptr {
		if src.IsNil() {
			return false
		}

		src = src.Elem()
	}

	// check if src implements driver.Valuer interface
	if !src.Type().Implements(reflect.TypeOf((*driver.Valuer)(nil)).Elem()) {
		return false
	}

	// execute Value() method
	results := src.MethodByName("Value").Call([]reflect.Value{})

	// check if Value() method returns nil
	value := results[0].Interface()
	if value == nil {
		return true
	}

	// set value if src and dst have the same type
	if valueOf := reflect.ValueOf(value); valueOf.Type() == dst.Type() {
		dst.Set(valueOf)
	}

	return true
}

func tryIfStruct(src, dst reflect.Value) bool {
	srcType := src.Type()
	dstType := dst.Type()

	if getStructType(srcType).Kind() != reflect.Struct || getStructType(dstType).Kind() != reflect.Struct {
		return false
	}

	if dst.IsNil() {
		dst.Set(reflect.New(getStructType(dstType)))
	}

	reflectCopy(src.Interface(), dst.Interface(), nil)

	return true
}

func tryIfSlice(src, dst reflect.Value) bool {
	srcType := src.Type()
	dstType := dst.Type()

	if srcType.Kind() != reflect.Slice || dstType.Kind() != reflect.Slice {
		return false
	}

	n := src.Len()

	tmpArr := reflect.MakeSlice(dstType, n, n)

	for i := 0; i < n; i++ {
		srcElem := src.Index(i)
		dstEl := tmpArr.Index(i)

		if dstEl.Type().Kind() != reflect.Ptr {
			dstEl = dstEl.Addr()
		} else {
			dstEl.Set(reflect.New(dstEl.Type().Elem()))
		}

		reflectCopy(srcElem.Interface(), dstEl.Interface(), nil)
	}

	dst.Set(tmpArr)

	return true
}

func getStructType(src reflect.Type) reflect.Type {
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}

	return src
}
