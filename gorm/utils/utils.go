package gormutils

import (
	"reflect"

	"gorm.io/gorm/schema"
)

// FieldToColMap create map of struct's field name to column from a dto that has gorm tags
// Example:
//
//	type User struct {
//		ID        int64     `gorm:"column:id"`
//		FirstName string    `gorm:"column:first_name"`
//		LastName  string    `gorm:"column:last_name"`
//	}
//
//	index := FieldToColMap(User{})
//	fmt.Println(index)
//	// Output
//	// map[FirstName:first_name ID:id LastName:last_name]
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
