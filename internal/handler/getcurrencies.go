package handler

import (
	"context"
	"encoding/json"
	"net/http"
)

func (h *Handler) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := h.db.GetCurrencies(context.TODO())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(currencies)
}