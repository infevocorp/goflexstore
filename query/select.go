package query

// SelectParam specifies the fields to be selected in the query result.
type SelectParam struct {
	Names []string
}

// ParamType returns the name of the param
func (p SelectParam) ParamType() string {
	return TypeSelect
}

// Select returns a SelectParam
//
// Example:
//
//	query.NewParams(
//		query.Select("ID", "Name"),
//	)
func Select(fields ...string) SelectParam {
	return SelectParam{
		Names: fields,
	}
}
