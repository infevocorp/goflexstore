package gormquery

import (
	"gorm.io/gorm"

	"github.com/jkaveri/goflexstore/query"
)

// ScopeFunc defines the type for GORM's scope function. It is a function that takes a GORM DB
// instance and returns a modified instance. It's used for applying query parameters to a GORM
// database query.
type ScopeFunc = func(*gorm.DB) *gorm.DB

// ScopeBuilderFunc is a type for functions that build a GORM's scope function from a query parameter.
// This allows for dynamic creation of query scopes based on different types of query parameters.
type ScopeBuilderFunc = func(query.Param) ScopeFunc

// ScopeBuilderRegistry is a map that acts as a registry for ScopeBuilderFuncs.
// It maps a query parameter type to its corresponding scope builder function. This registry is
// used to dynamically select the correct scope builder function based on the query parameter type.
type ScopeBuilderRegistry = map[string]ScopeBuilderFunc
