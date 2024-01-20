package gormquery

import (
	"reflect"
	"strings"

	"github.com/jkaveri/goflexstore/query"
)

func buildWhere(fieldName string, operator query.Operator, value any) (string, any) {
	if value == nil {
		panic("value cannot be nil")
	}

	var (
		valOf = reflect.ValueOf(value)
		kind  = valOf.Type().Kind()
	)

	if kind == reflect.Slice || kind == reflect.Array {
		n := valOf.Len()

		if n > 1 {
			return buildWhereInStr(fieldName, operator), value
		}

		return buildWhereStr(fieldName, operator), valOf.Index(0).Interface()
	}

	return buildWhereStr(fieldName, operator), value
}

func buildWhereStr(fieldName string, operator query.Operator) string {
	var sb strings.Builder

	sb.WriteString(fieldName)
	sb.WriteRune(' ')
	sb.WriteString(operatorToString(operator))
	sb.WriteString(" ?")

	return sb.String()
}

func buildWhereInStr(fieldName string, op query.Operator) string {
	var sb strings.Builder

	sb.WriteString(fieldName)
	sb.WriteRune(' ')
	sb.WriteString(inOperatorToString(op))
	sb.WriteString(" (?)")

	return sb.String()
}

func operatorToString(op query.Operator) string {
	switch op {
	case query.EQ:
		return "="
	case query.NEQ:
		return "!="
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

func inOperatorToString(op query.Operator) string {
	switch op {
	case query.EQ:
		return "IN"
	case query.NEQ:
		return "NOT IN"
	default:
		return "UNKNOWN"
	}
}
