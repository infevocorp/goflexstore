package query

type WithHintParam struct {
	Hint string
}

// ParamType returns the type of this parameter, which is TypeWithHint.
func (p WithHintParam) ParamType() string {
	return TypeWithHint
}

// WithHint creates a new WithHintParam.
// This function is used to add a "FOR UPDATE" clause to the main query.
//
// Parameters: N/A
//
// Returns:
// A new WithHintParam.
//
// Example:
// Using WithHint in a query:
//
//	query.NewParams(
//		query.Filter("Birthday", time.Parse("2000-01-01", "2006-01-02")).WithOP(query.GT),
//		query.WithHint("INL_HASH_JOIN(user)"),
//	)
//
// This example creates query parameters to filter records where 'Birthday' is greater than '2000-01-01'
// and provides and optimizer hint
func WithHint(hintType string) Param {
	return WithHintParam{Hint: hintType}
}
