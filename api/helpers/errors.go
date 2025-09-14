package helpers

import "errors"

var (
	ErrConflict          = errors.New("record already exists")
	ErrNotFound          = errors.New("resource not found")
	ErrInternalServer    = errors.New("internal server error")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
