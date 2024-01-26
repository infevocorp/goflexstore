package query

// GroupByParam represents a parameter used to group data in a query. It specifies the fields by which the data should be grouped.
// This is useful in aggregate queries where you need to group data by certain fields before applying aggregate functions.
//
// Fields:
//   - Names: A slice of field names to group by.
//   - Option: Additional options to apply to the group by operation, such as "ROLLUP".
//   - Having: A slice of FilterParam to specify the 'HAVING' clause conditions after grouping.
//
// Note: Using GroupByParam can make your code tightly coupled to the database's implementation of grouping,
// so it should be used with care to maintain database portability.
type GroupByParam struct {
	Names  []string
	Option string
	Having []FilterParam
}

// ParamType returns the type of this parameter, which is `groupby`. This method allows distinguishing GroupByParam
// from other query parameter types in contexts where multiple parameter types are used.
func (p GroupByParam) ParamType() string {
	return TypeGroupBy
}

// WithOption returns a new GroupByParam instance with the specified option while preserving the existing group by names and having conditions.
// This method is useful for adding additional grouping options to an existing GroupByParam.
//
// Parameters:
//   - option: A string representing the additional group by option to be applied.
//
// Returns:
// A new GroupByParam with the updated option.
func (p GroupByParam) WithOption(option string) GroupByParam {
	return GroupByParam{
		Names:  p.Names,
		Option: option,
	}
}

// WithHaving returns a new GroupByParam with the specified having conditions while preserving the existing group by names and options.
// This method is useful for adding 'HAVING' clause conditions to an existing GroupByParam.
//
// Parameters:
//   - params: A variable number of FilterParam representing the conditions to be applied in the 'HAVING' clause.
//
// Returns:
// A new GroupByParam with the updated having conditions.
func (p GroupByParam) WithHaving(params ...FilterParam) GroupByParam {
	return GroupByParam{
		Names:  p.Names,
		Having: params,
		Option: p.Option,
	}
}

// GroupBy creates a new GroupByParam with the specified field names for grouping.
// This function initializes a GroupByParam to group query results by the provided field names.
//
// Parameters:
//   - names: A variable number of strings representing the names of the fields to group by.
//
// Returns:
// A new GroupByParam with the specified field names.
//
// Example:
// Using GroupBy in a query:
//
//	query.NewParams(
//		query.Filter("Birthday", time.Parse("2000-01-01", "2006-01-02")).WithOP(query.GT),
//		query.GroupBy("Birthday"),
//	)
//
// This example creates query parameters to filter records where 'Birthday' is greater than '2000-01-01' and groups the results by 'Birthday'.
func GroupBy(names ...string) GroupByParam {
	return GroupByParam{
		Names: names,
	}
}
