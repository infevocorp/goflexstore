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
func ClauseLockForUpdate() Param {
	return ClauseLockForUpdateParam{}
}
