package query

// GroupByParam group by param, it is used to specify the field to group data when querying from store.
//
// Notes: this can make your code depends tightly on database, so use it carefully.
type GroupByParam struct {
	Names  []string
	Option string
	Having []FilterParam
}

// ParamType returns `groupby`
func (p GroupByParam) ParamType() string {
	return TypeGroupBy
}

// WithOption returns a new GroupByParam with the given option.
func (p GroupByParam) WithOption(option string) GroupByParam {
	return GroupByParam{
		Names:  p.Names,
		Option: option,
	}
}

// WithHaving returns a new GroupByParam with the given having.
func (p GroupByParam) WithHaving(params ...FilterParam) GroupByParam {
	return GroupByParam{
		Names:  p.Names,
		Having: params,
		Option: p.Option,
	}
}

// GroupBy returns a new GroupByParam with the given field name.
//
// Example:
//
//	query.NewParams(
//		query.Filter("Birthday", time.Parse("2000-01-01", "2006-01-02")).WithOP(query.GT),
//		query.GroupBy("Birthday"),
//	)
func GroupBy(names ...string) GroupByParam {
	return GroupByParam{
		Names: names,
	}
}
