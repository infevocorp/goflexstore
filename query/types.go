package query

const (
	// TypeFilter represents the type name for filter parameters in a query.
	// These parameters define conditions that data must meet to be included in the result set.
	TypeFilter = "filter"

	// TypeGroupBy represents the type name for group-by parameters in a query.
	// These parameters specify the fields that the result set should be grouped by.
	TypeGroupBy = "groupby"

	// TypeSelect represents the type name for select parameters in a query.
	// These parameters indicate the specific fields to be returned in the result set.
	TypeSelect = "select"

	// TypeOR represents the type name for OR logical operator parameters in a query.
	// These parameters are used to combine multiple conditions with OR logic, where any condition being true will result in a match.
	TypeOR = "or"

	// TypeOrderBy represents the type name for order-by parameters in a query.
	// These parameters define the sorting order of the result set based on specified fields.
	TypeOrderBy = "orderby"

	// TypePaginate represents the type name for pagination parameters in a query.
	// These parameters control the slicing of the result set into manageable segments, defining the offset and limit.
	TypePaginate = "paginate"

	// TypePreload represents the type name for preload parameters in a query.
	// These parameters specify related entities or fields that should be loaded along with the primary query results.
	TypePreload = "preload"
)
