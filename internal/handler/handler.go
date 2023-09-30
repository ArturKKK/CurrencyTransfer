package handler

import (
	"github.com/ArturKKK/CurrencyTransfer/internal/cache"
	"github.com/ArturKKK/CurrencyTransfer/internal/db"

	"github.com/gorilla/mux"
)

type Handler struct {
	db *db.Database
	cache cache.ValuteCache
}

func NewHander(db *db.Database, cache cache.ValuteCache) *Handler {
	return &Handler{
		db: db,
		cache: cache,
	}
}

func (h *Handler) Router() *mux.Router {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/currencies", h.GetCurrencies).Methods("GET")
	apiRouter.HandleFunc("/convert", h.ConvertCurrency).Methods("GET")

	return apiRouter
}
