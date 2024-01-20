package query

import (
	"fmt"
)

// ORParam query param
type ORParam struct {
	Params []FilterParam
}

// ParamType returns `or`
func (p ORParam) ParamType() string {
	return TypeOR
}

// OR returns a ORParam
//
// Example:
//
//	query.NewParams(
//		query.OR(
//			query.Filter("id", 1),
//			query.Filter("id", 2),
//		),
//	)
func OR(params ...Param) Param {
	filterParams := []FilterParam{}

	for _, p := range params {
		f, ok := p.(FilterParam)
		if !ok {
			panic(fmt.Errorf("OR only accept FilterParam but got %s", p.ParamType()))
		}

		filterParams = append(filterParams, f)
	}

	return ORParam{
		Params: filterParams,
	}
}
