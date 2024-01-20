// Package opscope is operation scope
//
// generic interface for operation scope
package opscope

import "context"

type Scope interface {
	// Begin begin operation scope
	Begin(ctx context.Context) (context.Context, error)
	// End end operation scope
	End(ctx context.Context, err error) error
	// EndWithRecover end operation scope with recovered error
	EndWithRecover(ctx context.Context, err *error)
}
