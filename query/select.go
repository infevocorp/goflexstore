package query

// SelectParam specifies the fields to be selected in the query result.
// This struct is useful for defining which specific fields of a data model should be retrieved in a query,
// allowing for more efficient data fetching and manipulation.
//
// Fields:
//   - Names: A slice of strings representing the names of the fields to be selected.
type SelectParam struct {
	Names []string
}

// ParamType returns the type of this parameter as a string.
// In this case, it returns 'select', indicating that this parameter is used for selecting specific fields.
func (p SelectParam) ParamType() string {
	return TypeSelect
}

// Select creates and returns a new SelectParam with the specified field names.
// This function is primarily used to construct query parameters that specify which fields
// of a data model should be included in the query's result set.
//
// Parameters:
//   - fields: A variable number of string arguments, each representing a field name to be included in the selection.
//
// Returns:
// A SelectParam struct containing the provided field names.
//
// Example:
// Creating a query parameter to select specific fields 'ID' and 'Name':
//
//	query.NewParams(
//		query.Select("ID", "Name"),
//	)
func Select(fields ...string) SelectParam {
	return SelectParam{
		Names: fields,
	}
}
