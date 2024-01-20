package query

// OrderByParam specifies the field Name to order by when querying from store.
// when Desc is true, the order is descending, otherwise ascending.
type OrderByParam struct {
	Name string
	Desc bool
}

// ParamType returns `orderby`
func (p OrderByParam) ParamType() string {
	return TypeOrderBy
}

// OrderBy returns a new OrderByParam with the given field name and order by desc or asc.
//
// Example:
//
// // other by name asc, then id desc
//
//	query.NewParams(
//		query.OrderBy("Name", false),
//		query.OrderBy("ID", true),
//	)
func OrderBy(name string, desc bool) OrderByParam {
	return OrderByParam{
		Name: name,
		Desc: desc,
	}
}
