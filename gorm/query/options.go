package gormquery

type Option func(*ScopeBuilder)

func WithCustomFilters(customFilters map[string]ScopeBuilderFunc) Option {
	return func(b *ScopeBuilder) {
		b.CustomFilters = customFilters
	}
}

func WithBuilder(name string, builder ScopeBuilderFunc) Option {
	return func(b *ScopeBuilder) {
		b.Registry[name] = builder
	}
}

func WithFieldToColMap(fieldToColMap map[string]string) Option {
	return func(b *ScopeBuilder) {
		b.FieldToColMap = fieldToColMap
	}
}
