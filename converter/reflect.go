package converter

import (
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/pkg/errors"

	"github.com/jkaveri/goflexstore/store"
)

// NewReflect creates a new converter that uses reflection to convert between DTO and Entity
//
// overridesMapping key is Entity's field name. value is DTO's field name, it is safe to be nil or empty
// by default it will use Entity's field name as DTO's field name
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

// Reflect is a converter that uses reflection to convert between DTO and Entity
type Reflect[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] struct {
	// fieldMapping key is Entity's field name. value is DTO's field name.
	dtoFieldsMapping map[string]string
	// fieldMapping key is DTO's field name. value is Entity's field name.
	entityFieldMapping map[string]string
}

func (c Reflect[Entity, DTO, ID]) ToEntity(dto DTO) Entity {
	entity := *new(Entity)

	reflectCopy(dto, &entity, c.entityFieldMapping)

	return entity
}

func (c Reflect[Entity, DTO, ID]) ToDTO(entity Entity) DTO {
	dto := *new(DTO)

	reflectCopy(entity, &dto, c.dtoFieldsMapping)

	return dto
}

func reflectCopy[SRC any, DST any](src SRC, dst DST, fieldMapping map[string]string) {
	srcVal := reflect.ValueOf(src)

	for srcVal.Kind() == reflect.Ptr {
		if srcVal.IsZero() {
			return
		}

		srcVal = srcVal.Elem()
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		panic("dst must be reference type (pointer)")
	}

	for dstVal.Kind() == reflect.Ptr {
		if dstVal.IsNil() {
			dstVal.Set(reflect.New(dstVal.Type().Elem()))
		}

		dstVal = dstVal.Elem()
	}

	dstType := dstVal.Type()

	numFiled := dstVal.NumField()
	for i := 0; i < numFiled; i++ {
		dstField := dstVal.Field(i)
		if !dstField.CanSet() {
			continue
		}

		dstFieldName := dstType.Field(i).Name
		if fieldMapping != nil {
			if f, ok := fieldMapping[dstFieldName]; ok && f != "" {
				dstFieldName = f
			}
		}

		srcField := srcVal.FieldByName(dstFieldName)
		if !srcField.IsValid() {
			continue
		}

		// no need to copy value from nil pointer
		if srcField.Kind() == reflect.Ptr && srcField.IsNil() {
			continue
		}

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
	if srcVal.Type().Kind() == dstVal.Kind() {
		dstVal.Set(srcVal)
		return true
	}

	if ok := tryIfTargetTypeIsScanner(srcVal, dstVal); ok {
		return true
	}

	if ok := tryIfTargetTypeIsValuer(srcVal, dstVal); ok {
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

	results := src.MethodByName("Value").Call([]reflect.Value{})
	value := results[0].Interface()

	switch v := value.(type) {
	case nil:
		break
	case bool:
		dst.SetBool(v)
	case int64:
		dst.SetInt(v)
	case int:
		dst.SetInt(int64(v))
	case int32:
		dst.SetInt(int64(v))
	case int16:
		dst.SetInt(int64(v))
	case int8:
		dst.SetInt(int64(v))
	case uint64:
		dst.SetUint(v)
	case uint:
		dst.SetUint(uint64(v))
	case uint32:
		dst.SetUint(uint64(v))
	case uint16:
		dst.SetUint(uint64(v))
	case uint8:
		dst.SetUint(uint64(v))
	case float64:
		dst.SetFloat(v)
	case float32:
		dst.SetFloat(float64(v))
	case string:
		dst.SetString(v)
	case []byte:
		dst.SetBytes(v)
	}

	return true
}
