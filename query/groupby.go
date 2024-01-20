package query

// GroupByParam group by param, it is used to specify the field to group data when querying from store.
type GroupByParam struct {
	Name string
}

// ParamType returns `group_by`
func (p GroupByParam) ParamType() string {
	return "group_by"
}

// GroupBy returns a new GroupByParam with the given field name.
//
// Example:
//
//	query.NewParams(
//		query.Filter("Birthday", time.Parse("2000-01-01", "2006-01-02")).WithOP(query.GT),
//		query.GroupBy("Birthday"),
//	)
func GroupBy(name string) GroupByParam {
	return GroupByParam{
		Name: name,
	}
}

// GroupByOptionParam group by option param, it is used to specify the option of group by when querying from store.
type GroupByOptionParam struct {
	Option string
}

// GetName returns `group_by_option`
func (p GroupByOptionParam) GetName() string {
	return "group_by_option"
}

// GroupByOption returns a new GroupByOptionParam with the given option.
//
// Example:
//
//	query.NewParams(
//		query.Filter("Birthday", time.Parse("2000-01-01", "2006-01-02")).WithOP(query.GT),
//		query.GroupBy("Birthday"),
//		query.GroupByOption("WITH ROLLUP"),
//	)
func GroupByOption(option string) GroupByOptionParam {
	return GroupByOptionParam{
		Option: option,
	}
}
