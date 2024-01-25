package query

// Param is an interface representing a query parameter.
// It provides a common method to identify the type of the parameter.
type Param interface {
	// ParamType returns the name of the param, used to identify the type of the query parameter.
	ParamType() string
}

// Params is a struct that aggregates multiple query parameters.
// It also provides methods to retrieve specific types of parameters and a caching mechanism for efficient retrieval.
type Params struct {
	params       []Param
	cachedFilter map[string]int
}

// Params returns the list of all query parameters.
func (p Params) Params() []Param {
	return p.params
}

// Get returns all query parameters of a specific type.
//
// Parameters:
//   - paramType: The type of parameters to retrieve.
//
// Returns:
// A slice of Param that match the specified paramType.
func (p Params) Get(paramType string) []Param {
	params := []Param{}

	for _, param := range p.params {
		if param.ParamType() == paramType {
			params = append(params, param)
		}
	}

	return params
}

// GetFilter returns the FilterParam with the given name, if it exists.
//
// Parameters:
//   - name: The name of the filter parameter to retrieve.
//
// Returns:
// A FilterParam and a boolean indicating whether it was found.
func (p Params) GetFilter(name string) (FilterParam, bool) {
	i, ok := p.cachedFilter[name]
	if ok {
		return p.params[i].(FilterParam), true
	}

	return FilterParam{}, false
}

// NewParams creates a new Params object with the given query parameters.
// It initializes a cache for filter parameters for efficient retrieval.
//
// Parameters:
//   - params: A variable number of Param to include in the Params object.
//
// Returns:
// A new Params object containing the provided query parameters.
//
// Example:
// Creating a new Params object with various query parameters:
//
//	query.NewParams(
//		query.Select("ID", "Name"),
//		query.OrderBy("ID", true),
//		query.Filter("ID", 1),
//		query.Filter("Name", "test"),
//	)
func NewParams(params ...Param) Params {
	cachedFilter := map[string]int{}

	for i, param := range params {
		if param.ParamType() == "filter" {
			cachedFilter[param.(FilterParam).Name] = i
		}
	}

	return Params{
		params:       params,
		cachedFilter: cachedFilter,
	}
}

// FilterGetter creates a function to retrieve a FilterParam from Params by a given name.
//
// Parameters:
//   - name: The name of the filter parameter to retrieve.
//
// Returns:
// A function that takes Params and returns a FilterParam and a boolean indicating whether it was found.
func FilterGetter(name string) func(Params) (FilterParam, bool) {
	return func(params Params) (FilterParam, bool) {
		return params.GetFilter(name)
	}
}
