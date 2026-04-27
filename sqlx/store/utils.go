package sqlxstore

import (
	"reflect"
	"strings"
	"sync"
)

// fieldMeta holds the pre-computed column name and field index for one struct field.
type fieldMeta struct {
	col string // db column name
	idx int    // reflect field index in the struct
}

// structMeta caches the field layout for a given (type, pkCol) pair.
type structMeta struct {
	fields []fieldMeta // all exported, non-"-" db-tagged fields
	pkIdx  int         // index into fields whose col == pkCol; -1 if not found
}

type metaCacheKey struct {
	t     reflect.Type
	pkCol string
}

var metaCache sync.Map // key: metaCacheKey, value: *structMeta

func getOrBuildMeta(rt reflect.Type, pkCol string) *structMeta {
	key := metaCacheKey{t: rt, pkCol: pkCol}
	if v, ok := metaCache.Load(key); ok {
		return v.(*structMeta)
	}

	n := rt.NumField()
	fields := make([]fieldMeta, 0, n)
	pkIdx := -1
	for i := 0; i < n; i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}
		col := colFromTag(f)
		if col == "-" || col == "" {
			continue
		}
		if col == pkCol {
			pkIdx = len(fields)
		}
		fields = append(fields, fieldMeta{col: col, idx: i})
	}

	meta := &structMeta{fields: fields, pkIdx: pkIdx}
	v, _ := metaCache.LoadOrStore(key, meta)
	return v.(*structMeta)
}

// getStructColVals iterates over the exported fields of a struct (passed as any,
// pointer or value), reads `db` tags for column names, and returns the ordered
// slice of column names and their corresponding values.
//
//   - skipPK: column name to skip (the PK column during INSERT/UPDATE); empty = skip nothing.
//   - onlyNonZero: when true, fields whose value is the zero value are skipped
//     (used for PartialUpdate).
func getStructColVals(row any, skipPK string, onlyNonZero bool) (cols []string, vals []any) {
	rv := reflect.ValueOf(row)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	meta := getOrBuildMeta(rv.Type(), skipPK)
	cols = make([]string, 0, len(meta.fields))
	vals = make([]any, 0, len(meta.fields))

	for _, f := range meta.fields {
		if skipPK != "" && f.col == skipPK {
			continue
		}
		fv := rv.Field(f.idx)
		if onlyNonZero && fv.IsZero() {
			continue
		}
		cols = append(cols, f.col)
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

// initScanTarget prepares a DTO variable for sqlx scanning.
// Pass &dto; if DTO is already a pointer type, it allocates the inner struct,
// sets dto to point to it, and returns a *InnerStruct suitable for sqlx.
// If DTO is a plain struct, it returns dtoPtr unchanged (*DTO).
func initScanTarget(dtoPtr any) any {
	rv := reflect.ValueOf(dtoPtr).Elem()
	if rv.Kind() == reflect.Ptr {
		inner := reflect.New(rv.Type().Elem())
		rv.Set(inner)
		return inner.Interface()
	}
	return dtoPtr
}

// setPKField sets the primary-key column field on a struct (passed by pointer)
// to the given int64 id, handling int and uint field kinds.
func setPKField(row any, pkCol string, id int64) {
	rv := reflect.ValueOf(row)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}

	meta := getOrBuildMeta(rv.Type(), pkCol)
	if meta.pkIdx < 0 {
		return
	}

	fv := rv.Field(meta.fields[meta.pkIdx].idx)
	if !fv.CanSet() {
		return
	}
	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fv.SetInt(id)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fv.SetUint(uint64(id))
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
