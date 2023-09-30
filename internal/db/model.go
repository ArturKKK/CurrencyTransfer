package db

type Currency struct {
	CharCode  string  `json:"char_code"`
	VunitRate float64 `json:"vunit_rate"`
}
