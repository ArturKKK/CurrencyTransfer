package db

import "fmt"

var (
	ErrCharCodeNotFound = fmt.Errorf("charCode doesn't exist")

	ErrValueNotFound = fmt.Errorf("vunitRate doesn't exist")
)
