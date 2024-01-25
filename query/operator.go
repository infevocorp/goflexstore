package query

import "fmt"

// Operator defines a set of constants representing operators used in filter expressions.
// These operators are used to specify the type of comparison to be performed in a query's filter condition.
type Operator uint8

const (
	// EQ represents the 'Equal' operator in a filter expression.
	EQ Operator = iota

	// NEQ represents the 'Not Equal' operator in a filter expression.
	NEQ

	// GT represents the 'Greater Than' operator in a filter expression.
	GT

	// GTE represents the 'Greater Than or Equal' operator in a filter expression.
	GTE

	// LT represents the 'Less Than' operator in a filter expression.
	LT

	// LTE represents the 'Less Than or Equal' operator in a filter expression.
	LTE
)

// String returns the string representation of the Operator.
// This method is useful for displaying or logging the operator in a human-readable format.
//
// Returns:
// A string that represents the Operator. For example, it returns "EQ" for the EQ operator.
// If the operator does not match any predefined operator, "UNKNOWN" is returned.
func (o Operator) String() string {
	switch o {
	case EQ:
		return "EQ"
	case NEQ:
		return "NEQ"
	case GT:
		return "GT"
	case GTE:
		return "GTE"
	case LT:
		return "LT"
	case LTE:
		return "LTE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", o)
	}
}
