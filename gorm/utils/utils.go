package gormutils

import (
	"reflect"

	"gorm.io/gorm/schema"
)

// FieldToColMap creates a map of struct field names to their corresponding database column names.
// This function is particularly useful for translating struct field names to database columns
// when working with GORM, especially when struct fields are tagged with GORM tags defining the column names.
//
// This function iterates over the fields of the provided struct (DTO), examines the `gorm` tag
// to find out the specified column name for each field, and then creates a mapping from the struct field name
// to the database column name. If a struct field does not have a `gorm` tag specifying a column name,
// the field name itself is used as the column name in the map.
//
// Parameter:
//
// dto - An instance of any struct type.
// This parameter is used to identify the struct fields and their corresponding GORM tags.
//
// Returns:
//
// A map where keys are struct field names and
// values are the corresponding database column names as defined by `gorm` tags.
//
// Example:
//
//	// Defining a User struct with GORM tags to specify database column names
//	type User struct {
//		ID        int64     `gorm:"column:id"`
//		FirstName string    `gorm:"column:first_name"`
//		LastName  string    `gorm:"column:last_name"`
//	}
//
//	// Creating a field-to-column map for the User struct
//	index := FieldToColMap(User{})
//	fmt.Println(index)
//	// Output:
//	// map[FirstName:first_name ID:id LastName:last_name]
//
// In this example, the User struct has fields ID, FirstName, and LastName. The `FieldToColMap` function
// creates a map where 'ID' maps to 'id', 'FirstName' maps to 'first_name', and 'LastName' maps to 'last_name'.
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
