// Package gormopscope provides tools for managing database transaction scopes
// in applications using GORM. It offers fine-grained control over transaction
// behavior, including support for nested transactions and different isolation levels.
//
// The package defines a TransactionScope struct that encapsulates transaction-related
// information and provides methods to begin, end, and manage the state of transactions.
// This allows for more robust and error-resistant transaction handling within GORM-based
// database operations.
package gormopscope
