package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"wallet-app/internal/app/domain"
	"wallet-app/internal/app/repository"
)

type Wallet interface {
	CreateWallet(ctx context.Context) (domain.Wallet, error)
	ProcessOperation(ctx context.Context, op domain.WalletOperation) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (decimal.Decimal, error)
}

type Service struct {
	Wallet
}

func NewService(repo *repository.WalletRepository) *Service {
	return &Service{
		Wallet: NewWalletService(repo),
	}
}
