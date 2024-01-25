package query

import (
	"fmt"
)

// ORParam represents a logical OR combination of multiple filter parameters.
// It is used in queries to combine multiple FilterParam instances such that
// any of the conditions being true will result in a match.
//
// Fields:
//   - Params: A slice of FilterParam representing the filter conditions to be combined with OR logic.
type ORParam struct {
	Params []FilterParam
}

// ParamType returns the type of this parameter, which is `or`.
// This method allows differentiating ORParam from other types of query parameters.
func (p ORParam) ParamType() string {
	return TypeOR
}

// OR creates a new ORParam, which is a logical OR combination of the provided filter parameters.
// This function is used to build queries where you want to match records that satisfy any one of the given filter conditions.
//
// Parameters:
//   - params: A variable number of Param, each of which should be a FilterParam.
//
// Returns:
// An ORParam that encapsulates the provided filter parameters in an OR logic.
//
// Example:
// Using OR to combine filter conditions:
//
//	query.NewParams(
//	  query.OR(
//	    query.Filter("id", 1),
//	    query.Filter("id", 2),
//	  ),
//	)
//
// This example creates query parameters that match records where 'id' is either 1 or 2.
//
// Note: The function panics if any parameter provided is not a FilterParam.
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
