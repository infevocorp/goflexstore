package sqlxutils

import (
	"reflect"
	"strings"
)

// FieldToColMap reads `db:"colname"` struct tags and returns a map from
// struct field name to database column name.
func FieldToColMap(row any) map[string]string {
	t := structType(row)
	m := make(map[string]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		col := colFromTag(f)
		if col == "-" {
			continue
		}
		m[f.Name] = col
	}

	return m
}

// colFromTag extracts the column name from a struct field's `db` tag.
// Falls back to the field name when no tag is present.
func colFromTag(f reflect.StructField) string {
	tag := f.Tag.Get("db")
	if tag == "" {
		return f.Name
	}
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	if tag == "" {
		return f.Name
	}
	return tag
}

func structType(v any) reflect.Type {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
