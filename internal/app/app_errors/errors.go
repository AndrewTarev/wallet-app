package app_errors

import "errors"

var (
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrAmountMustBePositive = errors.New("amount must be greater than zero")
	ErrInvalidAmount        = errors.New("invalid amount format")
)
