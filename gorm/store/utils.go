package gormstore

import (
	"reflect"

	"gorm.io/gorm/schema"
)

func defaultValue[T comparable](val T, defaultVal T) T {
	if val == (*new(T)) {
		return defaultVal
	}

	return val
}

func FieldToColMap(dto any) map[string]string {
	var (
		dtoTypeOf = reflect.TypeOf(dto)
		index     = map[string]string{}
		numField  = dtoTypeOf.NumField()
	)

	for i := 0; i < numField; i++ {
		field := dtoTypeOf.Field(i)
		if !field.IsExported() {
			continue
		}

		tagSettings := schema.ParseTagSetting(field.Tag.Get("gorm"), ";")
		if tagSettings["COLUMN"] != "" {
			index[field.Name] = tagSettings["COLUMN"]
		} else {
			index[field.Name] = field.Name
		}
	}

	return index
}
