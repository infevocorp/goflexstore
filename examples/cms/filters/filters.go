package filters

import "github.com/infevocorp/goflexstore/query"

func Tag(tag ...string) query.FilterParam {
	return query.Filter("tag", tag)
}

var GetTag = query.FilterGetter("tag")

func AuthorID(authorID ...int64) query.FilterParam {
	return query.Filter("author_id", authorID)
}

var GetAuthorID = query.FilterGetter("author_id")
