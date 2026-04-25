package sqlxquery

// Option configures a Builder.
type Option func(*Builder)

// WithFieldToColMap sets the struct-field-name → column-name mapping.
func WithFieldToColMap(m map[string]string) Option {
	return func(b *Builder) {
		b.FieldToColMap = m
	}
}

// WithDialect sets the SQL dialect used for placeholder rebinding.
func WithDialect(d Dialect) Option {
	return func(b *Builder) {
		b.Dialect = d
	}
}
