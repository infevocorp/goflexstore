// Package gormquery provides an extension for GORM, utilizing github.com/infevocorp/goflexstore/query.
// It facilitates the creation of GORM's scope functions, enhancing flexibility and reusability.
//
// The primary goal of gormquery is to seamlessly integrate complex query functionalities with GORM,
// allowing for dynamic query building based on various parameters. This integration simplifies the
// process of constructing sophisticated queries while maintaining clean and maintainable code.
//
// Key features include:
//   - Dynamic query generation: Based on parameters provided, it dynamically constructs GORM scope functions.
//   - Enhanced reusability: Encourages code reuse by abstracting common query patterns.
//   - Flexibility: Easily adapt to various querying requirements without changing the underlying database interactions.
//
// gormquery is especially useful in conjunction with the github.com/infevocorp/flexstore/store/gorm package,
// providing the necessary tools to create a generic, reusable store that leverages the power of GORM
// with enhanced query capabilities.
package gormquery
