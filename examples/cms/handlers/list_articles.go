package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/infevocorp/goflexstore/query"

	"github.com/infevocorp/goflexstore/examples/cms/filters"
)

type ListArticlesRequest struct {
	AuthorID int64  `query:"author_id"`
	Tag      string `query:"tag"`

	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}

func (h *Handler) ListArticles(c echo.Context) error {
	req := ListArticlesRequest{
		Offset: 0,
		Limit:  10,
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	params := []query.Param{
		query.Paginate(req.Offset, req.Limit),
		query.Preload("Author"),
		query.Preload("Tags"),
	}

	if req.AuthorID > 0 {
		params = append(params, filters.AuthorID(req.AuthorID))
	}

	if req.Tag != "" {
		params = append(params, filters.Tag(req.Tag))
	}

	articles, err := h.Stores.Article.List(c.Request().Context(), params...)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, articles)
}
