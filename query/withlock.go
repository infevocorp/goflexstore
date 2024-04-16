package query

const (
	LockTypeForUpdate LockType = iota
)

type LockType int

type WithLockParam struct {
	LockType LockType
}

// ParamType returns the type of this parameter, which is TypeWithLock.
// This method helps to identify WithLockParam as the parameter type for pagination purposes.
func (p WithLockParam) ParamType() string {
	return TypeWithLock
}

// WithLock creates a new WithLockParam.
// This function is used to add a "FOR UPDATE" clause to the main query.
//
// Parameters: N/A
//
// Returns:
// A new WithLockParam.
//
// Example:
// Using WithLock in a query:
//
//	query.NewParams(
//		query.Filter("Birthday", time.Parse("2000-01-01", "2006-01-02")).WithOP(query.GT),
//		query.WithLock(query.LockTypeForUpdate),
//	)
//
// This example creates query parameters to filter records where 'Birthday' is greater than '2000-01-01' and locks all
// the matching rows to be updated within the current transaction.
func WithLock(lockType LockType) Param {
	return WithLockParam{
		LockType: lockType,
	}
}
