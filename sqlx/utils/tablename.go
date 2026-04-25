package sqlxutils

import (
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
)

// TableName derives the default table name from a Row struct: strips any
// trailing "Row" suffix, converts to snake_case, then pluralises.
func TableName(row any) string {
	name := structType(row).Name()
	name = strings.TrimSuffix(name, "Row")
	return inflection.Plural(toSnakeCase(name))
}

func toSnakeCase(s string) string {
	runes := []rune(s)
	var sb strings.Builder

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			// Insert underscore before uppercase when preceded by a lowercase,
			// or when the following character is lowercase (e.g. "HTTPRequest").
			if unicode.IsLower(prev) || (i+1 < len(runes) && unicode.IsLower(runes[i+1])) {
				sb.WriteRune('_')
			}
		}
		sb.WriteRune(unicode.ToLower(r))
	}

	return sb.String()
}
