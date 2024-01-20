package query

// PreloadParam set the reference field to be preloaded when querying from store.
// we can specify the Params for preloading
type PreloadParam struct {
	Preload string
	Params  []Param
}

// GetName returns `preload`
func (p PreloadParam) GetName() string {
	return "preload"
}

// Preload returns a PreloadParam with the given reference field name.
//
// Example:
//
//	type Article struct {
//		ID      int
//		Title   string
//		Content string
//		Author  *Author
//	}
//
//	type Author struct {
//		ID   int
//		Name string
//	}
//
//	// preload author
//	query.NewParams(
//		query.Preload("Author"),
//	)
func Preload(preload string, params ...Param) PreloadParam {
	return PreloadParam{
		Preload: preload,
		Params:  params,
	}
}
