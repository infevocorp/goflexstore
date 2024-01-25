// Package query offers a flexible and extensible query builder for creating complex queries.
// It is designed to facilitate the construction of database queries in a structured and
// type-safe manner. The package provides various parameter types to specify different aspects
// of a query, such as filtering, sorting, grouping, and pagination.
//
// The query builder supports the creation of queries with multiple conditions and configurations,
// allowing for precise control over the retrieved data. It is especially useful in scenarios
// where dynamic query generation based on user input or application logic is required.
//
// Key Features:
// - Type-safe query parameters: Ensures that only valid query structures are created, reducing runtime errors.
// - Extensible design: Easy to extend with new parameter types as needed.
// - Readability: The query builder syntax is designed to be readable and easy to understand.
// - Integration with stores: Seamlessly integrates with various data store implementations, providing a unified querying interface.
//
// Example Usage:
// The following example demonstrates the use of the query package to build a complex query
// with multiple parameters:
//
//	query.NewParams(
//		query.Select("ID", "Name"),
//		query.Filter("Status", "Active"),
//		query.OrderBy("CreatedAt", true),
//		query.GroupBy("Category"),
//		query.Paginate(0, 10),
//	)
//
// This creates a query that selects the 'ID' and 'Name' fields, filters records by 'Status',
// orders them by 'CreatedAt' in descending order, groups them by 'Category', and applies
// pagination to retrieve the first 10 records.
//
// The query package is versatile and can be adapted to various data retrieval needs, making it
// a valuable tool for developers working with data stores in Go.
package query
