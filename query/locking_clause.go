package query

type ClauseLockForUpdateParam struct{}

// ParamType returns the type of this parameter, which is TypeClauseLockUpdate.
// This method helps to identify ClauseLockForUpdateParam as the parameter type for pagination purposes.
func (p ClauseLockForUpdateParam) ParamType() string {
	return TypeClauseLockUpdate
}

// ClauseLockForUpdate creates a new ClauseLockForUpdateParam.
// This function is used to add a "FOR UPDATE" clause to the main query.
//
// Parameters: N/A
//
// Returns:
// A new ClauseLockForUpdateParam.
//
// Example:
// Using ClauseLockForUpdate in a query:
//
//	query.NewParams(
//		query.Filter("Birthday", time.Parse("2000-01-01", "2006-01-02")).WithOP(query.GT),
//		query.ClauseLockForUpdate(),
//	)
//
// This example creates query parameters to filter records where 'Birthday' is greater than '2000-01-01' and locks all
// the matching rows to be updated within the current transaction.
func ClauseLockForUpdate() Param {
	return ClauseLockForUpdateParam{}
}
