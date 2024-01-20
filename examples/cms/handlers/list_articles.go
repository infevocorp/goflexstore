package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/jkaveri/goflexstore/examples/cms/filters"
	"github.com/jkaveri/goflexstore/query"
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
