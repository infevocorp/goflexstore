package query

// OrderByParam specifies how to sort the results when querying from a data store.
// It defines the field by which the results should be ordered and the direction of ordering.
//
// Fields:
//   - Name: The name of the field to be used for ordering.
//   - Desc: A boolean indicating the order direction. If true, the order is descending. If false, it's ascending.
type OrderByParam struct {
	Name string
	Desc bool
}

// ParamType returns the type of this parameter, which is `orderby`.
// This method is used to distinguish OrderByParam from other types of query parameters.
func (p OrderByParam) ParamType() string {
	return TypeOrderBy
}

// OrderBy creates a new OrderByParam with the specified field name and order direction.
// This function is used in query construction to specify how the results should be sorted.
//
// Parameters:
//   - name: The name of the field to order by.
//   - desc: Boolean indicating the order direction. True for descending order, false for ascending.
//
// Returns:
// A new OrderByParam configured with the specified field and order direction.
//
// Example:
// Ordering query results by 'Name' in ascending order and then by 'ID' in descending order:
//
//	query.NewParams(
//		query.OrderBy("Name", false), // Order by 'Name' in ascending order
//		query.OrderBy("ID", true),    // Then order by 'ID' in descending order
//	)
//
// In this example, results will be first sorted by 'Name' in ascending order, and then by 'ID' in descending order.
func OrderBy(name string, desc bool) OrderByParam {
	return OrderByParam{
		Name: name,
		Desc: desc,
	}
}
