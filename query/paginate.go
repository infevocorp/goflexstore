package query

// PaginateParam specifies the parameters for pagination when querying a data store.
// It is used to define the offset (starting point) and limit (number of items to fetch) for a query result.
//
// Fields:
//   - Offset: The number of items to skip before starting to collect the result set.
//   - Limit: The maximum number of items to return in the result set.
type PaginateParam struct {
	Offset int
	Limit  int
}

// ParamType returns the type of this parameter, which is `paginate`.
// This method helps to identify PaginateParam as the parameter type for pagination purposes.
func (p PaginateParam) ParamType() string {
	return TypePaginate
}

// Paginate creates a new PaginateParam with the specified offset and limit.
// This function is used to apply pagination to query results, controlling the portion of the result set to return.
//
// Parameters:
//   - offset: The number of items to skip in the result set.
//   - limit: The maximum number of items to include in the result set.
//
// Returns:
// A PaginateParam configured with the specified offset and limit.
//
// Example:
// Applying pagination to a query:
//
//	// Fetch the next 10 items starting from the 11th item.
//	params := query.NewParams(
//	  query.Paginate(10, 10),
//	)
//
// In this example, the query will skip the first 10 items and then fetch the next 10 items, effectively returning items 11 to 20.
func Paginate(offset, limit int) Param {
	return PaginateParam{
		Offset: offset,
		Limit:  limit,
	}
}
