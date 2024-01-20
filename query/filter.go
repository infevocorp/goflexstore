package query

// FilterParam is a query param that represents a filter.
//
// A filter is a condition that is used to filter data from store.
type FilterParam struct {
	Name     string
	Operator Operator
	Value    any
}

// ParamType returns `filter`
func (p FilterParam) ParamType() string {
	return TypeFilter
}

// WithOP returns a new FilterParam with the given Operator.
func (p FilterParam) WithOP(op Operator) FilterParam {
	return FilterParam{
		Name:     p.Name,
		Operator: op,
		Value:    p.Value,
	}
}

// Filter returns a new FilterParam with the given field name and value.
// The default operator is EQ, use WithOP to change the operator.
//
// Example:
//
//	query.Filter("id", 1)
//	query.Filter("id", 1).WithOP(query.GT)
func Filter(fieldName string, value any) FilterParam {
	return FilterParam{
		Name:     fieldName,
		Operator: EQ,
		Value:    value,
	}
}
