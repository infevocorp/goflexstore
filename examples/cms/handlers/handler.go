package handlers

import "github.com/jkaveri/goflexstore/examples/cms/store"

type Handler struct {
	Stores store.Stores
}

func NewHandler(stores store.Stores) *Handler {
	return &Handler{
		Stores: stores,
	}
}
