package sqlxstore

import (
	"reflect"
	"strings"
)

// getStructColVals iterates over the exported fields of a struct (passed as any,
// pointer or value), reads `db` tags for column names, and returns the ordered
// slice of column names and their corresponding values.
//
//   - excludeCols: columns to skip (e.g. the PK column during INSERT).
//   - onlyNonZero: when true, fields whose value is the zero value are skipped
//     (used for PartialUpdate).
func getStructColVals(row any, excludeCols map[string]bool, onlyNonZero bool) (cols []string, vals []any) {
	rv := reflect.ValueOf(row)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}
		col := colFromTag(f)
		if col == "-" || col == "" {
			continue
		}
		if excludeCols[col] {
			continue
		}
		fv := rv.Field(i)
		if onlyNonZero && fv.IsZero() {
			continue
		}
		cols = append(cols, col)
		vals = append(vals, fv.Interface())
	}
	return
}

// colFromTag extracts the db column name from a struct field tag.
func colFromTag(f reflect.StructField) string {
	tag := f.Tag.Get("db")
	if tag == "" {
		return f.Name
	}
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	return tag
}

// setPKField sets the primary-key column field on a struct (passed by pointer)
// to the given int64 id, handling int and uint field kinds.
func setPKField(row any, pkCol string, id int64) {
	rv := reflect.ValueOf(row)
	if rv.Kind() != reflect.Ptr {
		return
	}
	rv = rv.Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if colFromTag(f) != pkCol {
			continue
		}
		fv := rv.Field(i)
		if !fv.CanSet() {
			return
		}
		switch fv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fv.SetInt(id)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fv.SetUint(uint64(id))
		}
		return
	}
}

// buildPlaceholders returns n "?" strings.
func buildPlaceholders(n int) []string {
	ph := make([]string, n)
	for i := range ph {
		ph[i] = "?"
	}
	return ph
}

func defaultValue[T comparable](val, def T) T {
	if val == (*new(T)) {
		return def
	}
	return val
}
