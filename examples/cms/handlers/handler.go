package handlers

import (
	"github.com/labstack/echo/v4"

	"github.com/jkaveri/goflexstore/examples/cms/store"
)

type Handler struct {
	Stores store.Stores
}

func Register(stores store.Stores, e *echo.Echo) *Handler {
	h := &Handler{
		Stores: stores,
	}

	h.Register(e)

	return h
}

func (h *Handler) Register(echo *echo.Echo) {
	echo.GET("/articles", h.ListArticles)
}
