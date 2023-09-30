package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// GET /api/v1/convert?from={CharCode}&to={CharCode}&amount={amount}
func (h *Handler) ConvertCurrency(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	fromRate, err := h.cache.GetOne(from)
	if err != nil {
		fromRate, err = h.db.GetCurrencyRate(context.TODO(), from)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.cache.SetOne(from, fromRate)
	}

	toRate, err := h.cache.GetOne(to)
	if err != nil{
		toRate, err = h.db.GetCurrencyRate(context.TODO(), to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		h.cache.SetOne(to, toRate)
	}

	result := (amount / fromRate) * toRate

	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	_ = encoder.Encode(map[string]float64{"result": result})
}
