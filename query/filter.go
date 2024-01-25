package query

// FilterParam represents a query parameter used for filtering data.
// It encapsulates the necessary information to construct a part of a query
// that filters data based on a specific field, operator, and value.
//
// Fields:
// - Name: The name of the field in the data store to apply the filter on.
// - Operator: The operator (e.g., equals, greater than) used for comparing the field's value with the provided value.
// - Value: The value to be used in comparison for filtering.
type FilterParam struct {
	Name     string
	Operator Operator
	Value    any
}

// ParamType returns the type of this parameter, which is `filter`.
// This method can be used to differentiate FilterParam from other types of query parameters
// in a system where multiple parameter types are used.
func (p FilterParam) ParamType() string {
	return TypeFilter
}

// WithOP returns a new FilterParam instance with the specified Operator, keeping the field name and value unchanged.
// This method is useful for changing the comparison operator for an existing FilterParam.
//
// Parameters:
// - op: The new Operator to be used for the filter.
//
// Returns:
// A new FilterParam with the updated operator.
func (p FilterParam) WithOP(op Operator) FilterParam {
	return FilterParam{
		Name:     p.Name,
		Operator: op,
		Value:    p.Value,
	}
}

// Filter creates a new FilterParam with the specified field name and value.
// The default operator used for the filter is EQ (equals). To use a different operator,
// chain the resulting FilterParam with the WithOP method.
//
// Parameters:
// - fieldName: The name of the field to filter on.
// - value: The value to compare against the field's value.
//
// Returns:
// A new FilterParam with the specified field name, value, and default operator EQ.
//
// Examples:
// - query.Filter("id", 1) creates a filter to check if 'id' equals 1.
// - query.Filter("id", 1).WithOP(query.GT) creates a filter to check if 'id' is greater than 1.
func Filter(fieldName string, value any) FilterParam {
	return FilterParam{
		Name:     fieldName,
		Operator: EQ,
		Value:    value,
	}
}
