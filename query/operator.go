package query

// Operator is the operator used in filter
type Operator uint8

const (
	// EQ Equal
	EQ Operator = iota
	// NEQ Not Equal
	NEQ
	// GT Greater Than
	GT
	// GTE Greater Than or Equal
	GTE
	// LT Less Than
	LT
	// LTE Less Than or Equal
	LTE
)

// String returns the string representation of the operator
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
		return "UNKNOWN"
	}
}
