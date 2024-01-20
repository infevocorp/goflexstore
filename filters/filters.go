package filters

import "github.com/jkaveri/goflexstore/query"

func IDs[T comparable](ids ...T) query.FilterParam {
	return query.Filter("ID", ids)
}

func GetIDs[T comparable](params query.Params) (query.FilterParam, bool) {
	return params.GetFilter("ID")
}
