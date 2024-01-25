// Package opscope provides a unified interface for managing operation scopes within an application.
// This package defines the Scope interface, which is used to manage the lifecycle and context of operations,
// ensuring that resources are correctly allocated and released, and that error handling is consistent.
//
// The Scope interface offers methods to begin and end operation scopes, handling errors and
// providing support for error recovery mechanisms. This facilitates better control over complex
// operations that involve multiple steps, resources, or external interactions.
//
// Importantly, the opscope package is designed to be generic, allowing its integration into various
// application architectures, including those with complex transactional requirements or where
// operations span multiple service boundaries.
package opscope

import "context"

// Scope defines a generic interface for managing the lifecycle of an operation within a context.
// Implementations of this interface should provide mechanisms to begin, end, and handle errors
// within a specific operation context, ensuring proper resource management and error propagation.
type Scope interface {
	// Begin initializes an operation scope, setting up any necessary resources or context.
	// Returns an enhanced or modified context and any initialization error encountered.
	Begin(ctx context.Context) (context.Context, error)

	// End finalizes the operation scope, performing cleanup activities such as releasing resources
	// and handling any errors that occurred during the operation.
	// It accepts the current context and an error that represents the operation's outcome.
	End(ctx context.Context, err error) error

	// EndWithRecover is a specialized method to finalize the operation scope in scenarios
	// where panic recovery is required. It updates the provided error pointer with any recovered
	// error, ensuring that panic-induced errors are properly handled and reported.
	EndWithRecover(ctx context.Context, err *error)
}
