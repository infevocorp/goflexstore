package query

// Param is the interface for query param
type Param interface {
	// ParamType returns the name of the param
	ParamType() string
}

// Params list of query params
type Params struct {
	params       []Param
	cachedFilter map[string]int
}

// Params returns the params
func (p Params) Params() []Param {
	return p.params
}

// Get returns the params with given paramType
func (p Params) Get(paramType string) []Param {
	params := []Param{}

	for _, param := range p.params {
		if param.ParamType() == paramType {
			params = append(params, param)
		}
	}

	return params
}

// GetFilter returns the FilterParam with given name
func (p Params) GetFilter(name string) (FilterParam, bool) {
	i, ok := p.cachedFilter[name]
	if ok {
		return p.params[i].(FilterParam), true
	}

	return FilterParam{}, false
}

// NewParams returns a new Params
//
// Example:
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

// FilterGetter create a func to get FilterParam from Params with given name.
func FilterGetter(name string) func(Params) (FilterParam, bool) {
	return func(params Params) (FilterParam, bool) {
		return params.GetFilter(name)
	}
}
