package query

// SelectParam specifies the fields to be selected in the query result.
type SelectParam struct {
	Fields []string
}

// GetName returns the name of the param
func (p SelectParam) GetName() string {
	return "select"
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
		Fields: fields,
	}
}
