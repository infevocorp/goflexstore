// package gormquery
package gormquery

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/jkaveri/goflexstore/query"
)

// buildWhere constructs a GORM-compatible WHERE clause based on the provided field name, operator, and value.
// It supports handling both singular and collection types and constructs the appropriate query string.
// It panics if the provided value is nil to prevent runtime errors.
func buildWhere(fieldName string, operator query.Operator, value any) (string, any) {
	if value == nil {
		panic("value cannot be nil")
	}

	var (
		valOf = reflect.ValueOf(value)
		kind  = valOf.Type().Kind()
	)

	// Handle collection types (Slice or Array) to build a WHERE IN clause if necessary.
	if kind == reflect.Slice || kind == reflect.Array {
		n := valOf.Len()

		// For multiple items, build a WHERE IN clause.
		if n > 1 {
			return buildWhereInStr(fieldName, operator), value
		}

		// For a single item, revert to standard WHERE clause.
		return buildWhereStr(fieldName, operator), valOf.Index(0).Interface()
	}

	// For non-collection types, build a standard WHERE clause.
	return buildWhereStr(fieldName, operator), value
}

// buildWhereStr constructs a standard SQL WHERE clause string using the given field name and operator.
func buildWhereStr(fieldName string, operator query.Operator) string {
	var sb strings.Builder

	// Construct the WHERE clause.
	sb.WriteString(fieldName)
	sb.WriteRune(' ')
	sb.WriteString(operatorToString(operator))
	sb.WriteString(" ?")

	return sb.String()
}

// buildWhereInStr constructs a SQL WHERE IN clause string for handling collection types.
func buildWhereInStr(fieldName string, op query.Operator) string {
	var sb strings.Builder

	// Construct the WHERE IN clause.
	sb.WriteString(fieldName)
	sb.WriteRune(' ')
	sb.WriteString(inOperatorToString(op))
	sb.WriteString(" (?)")

	return sb.String()
}

// operatorToString converts a query.Operator to its equivalent SQL operator string.
func operatorToString(op query.Operator) string {
	switch op {
	case query.EQ:
		return "="
	case query.NEQ:
		return "<>"
	case query.GT:
		return ">"
	case query.GTE:
		return ">="
	case query.LT:
		return "<"
	case query.LTE:
		return "<="
	default:
		return "UNKNOWN"
	}
}

// inOperatorToString converts a query.Operator to its equivalent SQL IN operator string.
// It supports only the EQ and NEQ operators, defaulting to "UNKNOWN" for others.
func inOperatorToString(op query.Operator) string {
	switch op {
	case query.EQ:
		return "IN"
	case query.NEQ:
		return "NOT IN"
	default:
		panic(errors.Errorf("%s is unsupported operator for IN clause", op.String()))
	}
}
