package cache

import "github.com/ArturKKK/CurrencyTransfer/internal/db"

type ValuteCache interface {
	SetOne(key string, vunitRate float64) error
	GetOne(key string) (float64, error)
	Get() ([]db.Currency, error)
	Set(Currencies []db.Currency) error
}
