package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      uuid.UUID       `json:"walletId"`
	Balance decimal.Decimal `json:"balance"`
}

type WalletBalance struct {
	Balance decimal.Decimal `json:"balance"`
}
