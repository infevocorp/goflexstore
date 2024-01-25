package gormquery

// Option defines a function signature for options that can be applied to ScopeBuilder.
type Option func(*ScopeBuilder)

// WithCustomFilters applies custom filter functions to a ScopeBuilder.
// This function allows overriding default filter builders for specific filter parameters.
//
// Parameters:
// customFilters - A map of filter names to their corresponding custom filter functions.
//
// Example:
//
//	gormquery.WithCustomFilters(map[string]gormquery.ScopeBuilderFunc{
//	    "name": func(param query.Param) gormquery.ScopeFunc {
//	        return func(tx *gorm.DB) *gorm.DB {
//	            p := param.(query.FilterParam)
//	            return tx.Join("INNER JOIN `user_profiles` ON `user_profiles`.`user_id` = `users`.`id`").
//	                    Where("`user_profiles`.`first_name` = ?", p.Value)
//	        }
//	    },
//	})
//
// Usage of this function enables custom behavior for specific filter names, such as "name" in this example.
func WithCustomFilters(customFilters map[string]ScopeBuilderFunc) Option {
	return func(b *ScopeBuilder) {
		b.CustomFilters = customFilters
	}
}

// WithBuilder registers a new ScopeBuilderFunc under a specified name.
// This function is used to add new filter building capabilities to a ScopeBuilder.
//
// Parameters:
//
// name - The name under which the ScopeBuilderFunc will be registered.
// builder - The ScopeBuilderFunc to be registered.
//
// Example:
//
//	gormquery.WithBuilder("customFilterName", customScopeBuilderFunc)
//
// This example demonstrates how to register a custom filter function named "customFilterName".
func WithBuilder(name string, builder ScopeBuilderFunc) Option {
	return func(b *ScopeBuilder) {
		b.Registry[name] = builder
	}
}

// WithFieldToColMap configures a mapping from struct field names to database column names in ScopeBuilder.
// This function is useful when the field names in Go structs differ from the column names in the database.
//
// Parameters:
//
// fieldToColMap - A map where keys are struct field names and values are the corresponding database column names.
//
// Example:
//
//	gormquery.WithFieldToColMap(map[string]string{
//	    "FieldName": "db_column_name",
//	})
//
// This example maps the struct field "FieldName" to the database column "db_column_name".
func WithFieldToColMap(fieldToColMap map[string]string) Option {
	return func(b *ScopeBuilder) {
		b.FieldToColMap = fieldToColMap
	}
}
