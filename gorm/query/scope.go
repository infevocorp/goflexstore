package gormquery

import (
	"gorm.io/gorm"

	"github.com/jkaveri/goflexstore/query"
)

// ScopeFunc is type of GORM's scope function
type ScopeFunc = func(*gorm.DB) *gorm.DB

// ScopeBuilder is a function that build a GORM's scope function from a query param
type ScopeBuilderFunc = func(query.Param) ScopeFunc

// ScopeBuilderRegistry is a registry of ScopeBuilder
type ScopeBuilderRegistry = map[string]ScopeBuilderFunc
