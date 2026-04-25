package sqlxutils

import (
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
)

// TableNamer can be implemented by a DTO to override the derived table name.
type TableNamer interface {
	TableName() string
}

// TableName derives the default table name from a DTO struct: strips any
// trailing "DTO" suffix, converts to snake_case, then pluralises.
// If the DTO implements TableNamer, that value is used instead.
func TableName(dto any) string {
	if tn, ok := dto.(TableNamer); ok {
		return tn.TableName()
	}
	name := structType(dto).Name()
	name = strings.TrimSuffix(name, "DTO")
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
