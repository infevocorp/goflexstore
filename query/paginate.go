package query

// PaginateParam specifies the offset and limit to paginate the query result from store.
type PaginateParam struct {
	Offset int
	Limit  int
}

// ParamType returns `paginate`
func (p PaginateParam) ParamType() string {
	return TypePaginate
}

// Paginate returns a PaginateParam
func Paginate(offset, limit int) Param {
	return PaginateParam{
		Offset: offset,
		Limit:  limit,
	}
}
