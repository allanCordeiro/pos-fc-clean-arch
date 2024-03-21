package entity

import "errors"

var (
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidPrice = errors.New("invalid price")
	ErrInvalidTax   = errors.New("invalid tax")
)
