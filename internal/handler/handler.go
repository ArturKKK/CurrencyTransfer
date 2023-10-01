package handler

import (
	"github.com/ArturKKK/CurrencyTransfer/internal/cache"
	"github.com/ArturKKK/CurrencyTransfer/internal/db"
	"github.com/ArturKKK/CurrencyTransfer/pkg/logging"

	"github.com/gorilla/mux"
)

type Handler struct {
	db     *db.Database
	cache  cache.ValuteCache
	logger *logging.Logger
}

func NewHander(db *db.Database, cache cache.ValuteCache, logger *logging.Logger) *Handler {
	return &Handler{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (h *Handler) Router() *mux.Router {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/currencies", h.GetCurrencies).Methods("GET")
	apiRouter.HandleFunc("/convert", h.ConvertCurrency).Methods("GET")

	return apiRouter
}
