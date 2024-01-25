package query

// PreloadParam is used to specify related entities to be preloaded when querying from a data store.
// This is particularly useful in ORM frameworks to efficiently load associated data in a single query.
//
// Fields:
//   - Name: The name of the related entity (reference field) to be preloaded.
//   - Params: Additional query parameters to apply to the preloading operation (e.g., filters, sorting).
type PreloadParam struct {
	Name   string
	Params []Param
}

// ParamType returns the type of this parameter, which is `preload`.
// This method distinguishes PreloadParam from other types of query parameters.
func (p PreloadParam) ParamType() string {
	return TypePreload
}

// Preload creates a new PreloadParam for a given reference field.
// This function is used to specify related entities that should be preloaded along with the main query results.
//
// Parameters:
//   - preload: The name of the reference field to preload.
//   - params: Optional additional query parameters to customize the preloading operation.
//
// Returns:
// A new PreloadParam configured with the specified reference field and additional parameters.
//
// Example:
// Preloading an 'Author' entity in an 'Article' query:
//
//	type Article struct {
//	    ID      int
//	    Title   string
//	    Content string
//	    Author  *Author
//	}
//
//	type Author struct {
//	    ID   int
//	    Name string
//	}
//
//	// Preload 'Author' data when querying 'Article'
//	query.NewParams(
//	    query.Preload("Author"),
//	)
//
// In this example, when querying for 'Article', the related 'Author' data is also loaded in the same query.
func Preload(preload string, params ...Param) PreloadParam {
	return PreloadParam{
		Name:   preload,
		Params: params,
	}
}
